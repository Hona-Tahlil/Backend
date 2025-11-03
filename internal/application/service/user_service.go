package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
	domainredis "hona/backend/internal/domain/ports/redis"
	"hona/backend/internal/infrastructure/mail"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	unitOfWork          ports.UnitOfWork
	userCacheRepository domainredis.UserCacheRepository
	emailService        *mail.EmailService
}

func NewGeneralService(unitOfWork ports.UnitOfWork, userCacheRepository domainredis.UserCacheRepository, emailService *mail.EmailService) *UserService {
	return &UserService{
		unitOfWork:          unitOfWork,
		userCacheRepository: userCacheRepository,
		emailService:        emailService,
	}
}

func (us *UserService) Login(loginInfo user.LoginRequest) user.LoginResponse {
	token := loginInfo.Email + "/" + loginInfo.Password

	return user.LoginResponse{
		JWTToken: token,
	}
}
func (us *UserService) FindUserByEmail(email string) (*entities.User, error) {
	foundUser, err := us.unitOfWork.Factory().UserRepository().FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if foundUser == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.User)
		return nil, NotFoundError
	}

	return foundUser, nil
}

func (us *UserService) FindVerifiedUserByEmail(email string) (*entities.User, error) {
	foundUser, err := us.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if !foundUser.IsEmailVerified {
		notVerifiedErr := exceptions.NewNotVerifiedError()
		return nil, notVerifiedErr
	}

	return foundUser, nil
}

func (us *UserService) validateDuplicateEmail(email string) error {
	var ce exceptions.ConflictErrors
	redisKey := bootstrap.Run().Constants.RedisKey.GenerateMLKey(email)
	data, err := us.userCacheRepository.Get(context.Background(), redisKey)
	if err != nil {
		return err
	}
	if data != nil {
		ce.Add(bootstrap.Run().Constants.ErrorFields.Email, bootstrap.Run().Constants.ErrorTags.AlreadyRegistered)
		return ce
	}
	user, err := us.FindUserByEmail(email)
	if err != nil {
		if _, ok := err.(*exceptions.NotFoundError); !ok {
			return err
		}
	}
	if user != nil && user.IsEmailVerified {
		ce.Add(bootstrap.Run().Constants.ErrorFields.Email, bootstrap.Run().Constants.ErrorTags.AlreadyRegistered)
		return ce
	}
	return nil
}

func (us *UserService) ValidatePasswordRegex(password string) error {
	var ve exceptions.ValidationErrors
	if len(password) < 8 {
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Password, bootstrap.Run().Constants.ErrorTags.MinimumLength)
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Password, bootstrap.Run().Constants.ErrorTags.ContainsLowercase)
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Password, bootstrap.Run().Constants.ErrorTags.ContainsUppercase)
	}
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	if !hasDigit {
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Password, bootstrap.Run().Constants.ErrorTags.ContainsNumber)
	}
	hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(password)
	if !hasSpecial {
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Password, bootstrap.Run().Constants.ErrorTags.ContainsSpecialChar)
	}

	if !hasLower || !hasUpper || !hasDigit || !hasSpecial {
		return &ve
	}

	return nil
}

func (us *UserService) generateRandomToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

func (us *UserService) CreateMagicLink(token, email string) string {
	baseURL := bootstrap.Run().Env.URLs.BaseURL
	return baseURL + "/auth/verify/email?token=" + token + "?email=" + email
}

func (us *UserService) Register(registerInfo user.RegisterRequest) error {
	err := us.validateDuplicateEmail(registerInfo.Email)
	if err != nil {
		return err
	}
	err = us.ValidatePasswordRegex(registerInfo.Password)
	if err != nil {
		return err
	}

	hashesPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(registerInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = us.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		err = rf.UserRepository().DeleteUserByEmail(registerInfo.Email)
		if err != nil {
			return err
		}
		newUser := &entities.User{
			FirstName:       registerInfo.FirstName,
			LastName:        registerInfo.LastName,
			Email:           registerInfo.Email,
			Password:        string(hashesPasswordBytes),
			IsEmailVerified: false,
		}
		err = rf.UserRepository().CreateUser(newUser)
		if err != nil {
			return err
		}

		err = us.SendVerificationEmail(user.SendVerificationEmailRequest{Email: newUser.Email})
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
func (us *UserService) CreateFPLink(token string, email string) string {
	baseURL := bootstrap.Run().Env.URLs.BaseURL
	return baseURL + "/auth/reset-password?token=" + token + "?email=" + email
}

func (us *UserService) SendRestPassEmail(email string) error {
	user, err := us.FindUserByEmail(email)
	if err != nil {
		return err
	}
	token, err := us.generateRandomToken()
	if err != nil {
		return err
	}
	redisKey := bootstrap.Run().Constants.RedisKey.GenerateFPKey(user.Email)
	err = us.userCacheRepository.Set(context.Background(), redisKey, token, time.Duration(bootstrap.Run().Env.EmailVerification.ExpireMinutes))
	if err != nil {
		return err
	}
	link := us.CreateFPLink(token, email)
	data := struct {
		FirstName    string
		LastName     string
		ResetLink    string
		ExpiryMinute int
		Year         int
	}{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		ResetLink:    link,
		ExpiryMinute: bootstrap.Run().Env.EmailVerification.ExpireMinutes,
		Year:         time.Now().Year(),
	}
	us.emailService.SendEmail(user.Email, "Reset Password", "forget_password.html", data)

	return nil
}

func (us *UserService) ForgotPassword(forgetPasswordInfo user.ForgotPasswordRequest) error {
	_, err := us.FindUserByEmail(forgetPasswordInfo.Email)
	if err != nil {
		return err
	}
	err = us.SendRestPassEmail(forgetPasswordInfo.Email)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) ResetPassword(resetPasswordInfo user.ResetPasswordRequest) error {
	RPData, err := us.userCacheRepository.Get(context.Background(), bootstrap.Run().Constants.RedisKey.GenerateMLKey(resetPasswordInfo.Email))
	if err != nil {
		return err
	}
	if RPData == nil {
		return exceptions.NewNotFoundError("token")
	}
	if RPData.Token != resetPasswordInfo.Token {
		var ve exceptions.ValidationErrors
		ve.AddError("token", "wrongToken")
		return &ve
	}
	user, err := us.FindUserByEmail(resetPasswordInfo.Email)
	if err != nil {
		return err
	}
	err = us.ValidatePasswordRegex(resetPasswordInfo.Password)
	if err != nil {
		return err
	}
	hashesPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(resetPasswordInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashesPasswordBytes)
	err = us.unitOfWork.Factory().UserRepository().SaveUser(user)
	if err != nil {
		return err
	}
	return nil

}

func (us *UserService) VerifyEmail(info user.VerifyEmailRequest) error {
	mlData, err := us.userCacheRepository.Get(context.Background(), bootstrap.Run().Constants.RedisKey.GenerateMLKey(info.Email))
	if err != nil {
		return err
	}
	if mlData == nil {
		return exceptions.NewNotFoundError("token")
	}
	if mlData.Token != info.Token {
		var ve exceptions.ValidationErrors
		ve.AddError("token", "wrongToken")
		return &ve
	}

	user, err := us.FindUserByEmail(info.Email)
	if err != nil {
		return err
	}
	user.IsEmailVerified = true
	us.unitOfWork.Factory().UserRepository().SaveUser(user)

	return nil
}

func (us *UserService) SendVerificationEmail(info user.SendVerificationEmailRequest) error {
	user, err := us.FindUserByEmail(info.Email)
	if err != nil {
		return err
	}
	token, err := us.generateRandomToken()
	if err != nil {
		return err
	}
	redisKey := bootstrap.Run().Constants.RedisKey.GenerateMLKey(user.Email)
	err = us.userCacheRepository.Set(context.Background(), redisKey, token, time.Duration(bootstrap.Run().Env.EmailVerification.ExpireMinutes))
	if err != nil {
		return err
	}

	link := us.CreateMagicLink(token, user.Email)
	data := struct {
		FirstName    string
		LastName     string
		MagicLink    string
		ExpiryMinute int
		Year         int
	}{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		MagicLink:    link,
		ExpiryMinute: bootstrap.Run().Env.EmailVerification.ExpireMinutes,
		Year:         time.Now().Year(),
	}
	us.emailService.SendEmail(user.Email, "Email Verification", bootstrap.Run().Constants.TemplatesPath.EmailVerification, data)

	return nil
}
