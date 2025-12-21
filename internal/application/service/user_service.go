package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/rbac"
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/exceptions"
	domainjwt "hona/backend/internal/domain/jwt"
	"hona/backend/internal/domain/ports"
	domainredis "hona/backend/internal/domain/ports/redis"
	"hona/backend/internal/infrastructure/communication/mail"
	"hona/backend/internal/infrastructure/persistence/repository/postgres"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	jwtService          domainjwt.JWTService
	unitOfWork          ports.UnitOfWork
	userCacheRepository domainredis.UserCacheRepository
	emailService        *mail.EmailService
}

func NewUserService(jwtService domainjwt.JWTService, unitOfWork ports.UnitOfWork, userCacheRepository domainredis.UserCacheRepository, emailService *mail.EmailService) *UserService {
	return &UserService{
		unitOfWork:          unitOfWork,
		userCacheRepository: userCacheRepository,
		emailService:        emailService,
		jwtService:          jwtService,
	}
}

func (us *UserService) GetRolesResponse(user *entities.User) []rbac.RoleResponse {
	r := make([]rbac.RoleResponse, len(user.Roles))
	for j, role := range user.Roles {
		p := make([]rbac.PermissionResponse, len(role.Permissions))
		for i, per := range role.Permissions {
			des := ""
			if per.Description != nil {
				des = *per.Description
			}
			p[i] = rbac.PermissionResponse{
				ID:          per.ID,
				Name:        per.Type.String(),
				Description: des,
				Category:    per.Category.String(),
			}
		}
		des := ""
		if role.Description != nil {
			des = *role.Description
		}
		r[j] = rbac.RoleResponse{
			ID:          role.ID,
			Name:        role.Type,
			Description: des,
			Permissions: p,
		}
	}
	return r
}

func (us *UserService) Login(loginInfo user.LoginRequest) (*user.LoginResponse, string, int, error) {
	foundUser, err := us.FindUserByEmail(loginInfo.Email)

	if err != nil {
		if _, ok := err.(*exceptions.NotFoundError); !ok {
			return nil, "", 0, err
		}

		invalidCredentialsErr := exceptions.NewInvalidCredentialsError("password is wrong")
		return nil, "", 0, invalidCredentialsErr
	}

	err = us.PreloadFields(foundUser, []string{"Roles.Permissions"})
	if err != nil {
		return nil, "", 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginInfo.Password)); err != nil {
		invalidCredentialsErr := exceptions.NewInvalidCredentialsError("password is wrong")
		return nil, "", 0, invalidCredentialsErr
	}
	accessToken, refreshToken, expireTime := us.jwtService.GenerateTokens(foundUser.ID, loginInfo.RememberMe)

	roles := us.GetRolesResponse(foundUser)

	return &user.LoginResponse{
		AccessToken: accessToken,
		Roles:       roles,
	}, refreshToken, expireTime, nil
}

func (us *UserService) GetUserInfosResponse(users []entities.User) []rbac.UserInfoResponse {
	r := make([]rbac.UserInfoResponse, len(users))
	for i, user := range users {
		r[i] = rbac.UserInfoResponse{
			Email: user.Email,
		}
	}
	return r
}

func (us *UserService) GetRoleUsersByID(roleID uint, options *postgres.QueryOptions) ([]entities.User, int64, error) {
	userRepo := us.unitOfWork.Factory().UserRepository()
	users, total, err := userRepo.GetRoleUsersByID(roleID, options)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
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

func (us *UserService) FindUserByEmail(email string) (*entities.User, error) {
	userRepo := us.unitOfWork.Factory().UserRepository()
	foundUser, err := userRepo.FindUserByEmail(email)
	if foundUser == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.User)
		return nil, NotFoundError
	}

	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (us *UserService) FindVerifiedUserByID(id uint) (*entities.User, error) {
	foundUser, err := us.FindUserByID(id)
	if err != nil {
		return nil, err
	}

	if !foundUser.IsEmailVerified {
		notVerifiedErr := exceptions.NewNotVerifiedError()
		return nil, notVerifiedErr
	}

	return foundUser, nil
}

func (us *UserService) FindUserByID(id uint) (*entities.User, error) {
	userRepo := us.unitOfWork.Factory().UserRepository()
	foundUser, err := userRepo.FindUserByID(id)
	if foundUser == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.User)
		return nil, NotFoundError
	}

	if err != nil {
		return nil, err
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
		return &ce
	}
	user, err := us.FindUserByEmail(email)
	if err != nil {
		if _, ok := err.(*exceptions.NotFoundError); !ok {
			return err
		}
	}
	if user != nil && user.IsEmailVerified {
		ce.Add(bootstrap.Run().Constants.ErrorFields.Email, bootstrap.Run().Constants.ErrorTags.AlreadyRegistered)
		return &ce
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
	return baseURL + "/auth/verify/email?token=" + token + "&email=" + email
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

	registrationErr := us.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		err := rf.UserRepository().DeleteUserByEmail(registerInfo.Email)
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

		err = us.SendVerificationEmail(user.SendVerificationEmailRequest{Email: newUser.Email, FirstName: newUser.FirstName, LastName: newUser.LastName})
		if err != nil {
			return err
		}

		return nil
	})

	return registrationErr
}
func (us *UserService) CreateFPLink(token string, email string) string {
	baseURL := bootstrap.Run().Env.URLs.BaseURL
	return baseURL + "auth/reset-password?token=" + token + "&email=" + email
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
	// user, err := us.FindUserByEmail(info.Email)
	// if err != nil {
	// 	return err
	// }
	token, err := us.generateRandomToken()
	if err != nil {
		return err
	}
	redisKey := bootstrap.Run().Constants.RedisKey.GenerateMLKey(info.Email)
	err = us.userCacheRepository.Set(context.Background(), redisKey, token, time.Duration(bootstrap.Run().Env.EmailVerification.ExpireMinutes))
	if err != nil {
		return err
	}

	link := us.CreateMagicLink(token, info.Email)
	data := struct {
		FirstName    string
		LastName     string
		MagicLink    string
		ExpiryMinute int
		Year         int
	}{
		FirstName:    info.FirstName,
		LastName:     info.LastName,
		MagicLink:    link,
		ExpiryMinute: bootstrap.Run().Env.EmailVerification.ExpireMinutes,
		Year:         time.Now().Year(),
	}
	return us.emailService.SendEmail(info.Email, "Email Verification", bootstrap.Run().Constants.TemplatesPath.EmailVerification, data)
}

func (us *UserService) RefreshTokens(refreshTokenInfo rbac.RefreshTokenRequest) (*rbac.RefreshTokenResponse, string, int, error) {
	accessToken, refreshToken, userID, expireTime := us.jwtService.RefreshTokens(refreshTokenInfo.RefreshToken)

	foundUser, err := us.FindUserByID(userID)
	if err != nil {
		return nil, "", 0, err
	}

	err = us.PreloadFields(foundUser, []string{"Roles.Permissions"})
	if err != nil {
		return nil, "", 0, err
	}

	roles := us.GetRolesResponse(foundUser)

	return &rbac.RefreshTokenResponse{
		AccessToken: accessToken,
		Roles:       roles,
	}, refreshToken, expireTime, nil
}

func (us *UserService) PreloadFields(user *entities.User, fields []string) error {
	userRepo := us.unitOfWork.Factory().UserRepository()
	return userRepo.PreloadFields(user, fields)
}
