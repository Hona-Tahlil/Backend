package service

import (
	"errors"
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/infrastructure/persistence"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	unitOfWork *persistence.UnitOfWork
}

func NewGeneralService(unitOfWork *persistence.UnitOfWork) *UserService {
	return &UserService{
		unitOfWork: unitOfWork,
	}
}

func (us *UserService) Login(loginInfo user.LoginRequest) user.LoginResponse {
	token := loginInfo.Email + "/" + loginInfo.Password

	return user.LoginResponse{
		JWTToken: token,
	}
}

func (us *UserService) validateDuplicateEmail(email string) error {
	var ce exceptions.ConflictErrors

	user, err := us.unitOfWork.Factory().UserRepository().FindUserByEmail(email)
	if err != nil {
		return err
	}
	if user != nil && user.EmailVerified {
		return ce
	}
	return nil
}


func (us *UserService)ValidatePasswordRegex(password string) error {
	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(password) {
		return errors.New("رمز عبور باید حداقل ۸ کاراکتر و شامل حروف بزرگ، کوچک، عدد و کاراکتر خاص باشد")
	}
	return nil
}


func (us *UserService) Register(registerInfo user.RegisterRequest) error {
	//validate Duplicate Email
	err := us.validateDuplicateEmail(registerInfo.Email)
	if err != nil {
		return err
	}
	//password Validation
	err = us.ValidatePasswordRegex(registerInfo.Password)
	if err != nil {
		return err
	}
	//Hash Password
	hashesPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(registerInfo.Password), 14)
	if err != nil {
		return err
	}
	//Insertion
	err = us.unitOfWork.WithTransaction(func(rf *persistence.RepositoryFactory) error {
		err = us.unitOfWork.Factory().UserRepository().DeleteUserByEmail(registerInfo.Email)
		if err != nil {
			return err
		}
		user := &entities.User{
			Name:          registerInfo.Name,
			Email:         registerInfo.Email,
			Password:      string(hashesPasswordBytes),
			EmailVerified: false,
		}
		err = us.unitOfWork.Factory().UserRepository().CreateUser(user)
		if err != nil {
			return err
		}
		return nil
	})

	return nil
}

func (us *UserService) VerifyEmail(verifyEmailInfo user.VerifyEmailRequest) error {
	return nil
}

func (us *UserService) ForgotPassword(forgetPasswordInfo user.ForgotPasswordRequest) error {
	return nil
}
