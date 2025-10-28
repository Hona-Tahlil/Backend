package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/infrastructure/mail"
	"hona/backend/internal/infrastructure/persistence"
	"hona/backend/internal/infrastructure/persistence/repository/redis"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	unitOfWork          *persistence.UnitOfWork
	userCacheRepository *redis.UserCacheRepository
	sendMLEmail         *mail.EmailService
}

func NewGeneralService(unitOfWork *persistence.UnitOfWork, userCacheRepository *redis.UserCacheRepository, sendMLEmail *mail.EmailService) *UserService {
	return &UserService{
		unitOfWork:          unitOfWork,
		userCacheRepository: userCacheRepository,
		sendMLEmail:         sendMLEmail,
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

	if !foundUser.IsVerified {
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
	if user != nil && user.IsVerified {
		ce.Add(bootstrap.Run().Constants.ErrorFields.Email, bootstrap.Run().Constants.ErrorTags.AlreadyRegistered)
		return ce
	}
	return nil
}

func (us *UserService) ValidatePasswordRegex(password string) error {
	if len(password) < 12 {
		return errors.New("password must be at least 12 characters long")
	}
	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(password) {
		return errors.New("رمز عبور باید حداقل ۸ کاراکتر و شامل حروف بزرگ، کوچک، عدد و کاراکتر خاص باشد")
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

func (us *UserService) CreateMagicLink(email string) (string, error) {
	token, err := us.generateRandomToken()
	if err != nil {
		return "", err
	}
	return token, nil
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
	err = us.unitOfWork.WithTransaction(func(rf *persistence.RepositoryFactory) error {
		err = us.unitOfWork.Factory().UserRepository().DeleteUserByEmail(registerInfo.Email)
		if err != nil {
			return err
		}
		user := &entities.User{
			Name:       registerInfo.Name,
			Email:      registerInfo.Email,
			Password:   string(hashesPasswordBytes),
			IsVerified: false,
		}
		err = us.unitOfWork.Factory().UserRepository().CreateUser(user)
		if err != nil {
			return err
		}
		token, err := us.CreateMagicLink(user.Email)
		if err != nil {
			return err
		}
		redisKey := bootstrap.Run().Constants.RedisKey.GenerateMLKey(user.Email)
		err = us.userCacheRepository.Set(context.Background(), redisKey, token, time.Duration(2)*time.Minute)
		if err != nil {
			return err
		}
		err = us.sendMLEmail.SendMLEmail(user.Email, token)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (us *UserService) ForgotPassword(forgetPasswordInfo user.ForgotPasswordRequest) error {
	return nil
}
