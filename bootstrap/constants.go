package bootstrap

type Constants struct {
	Context         Context
	JWTKeysPath     JWTKeysPath
	ErrorFields     ErrorFields
	ErrorTags       ErrorTags
	JWTConstants    JWTConstants
	SuccessMessages SuccessMessages
	EntityConstants EntityConstants
	TemplatesPath   TemplatesPath
}

type EntityConstants struct {
	Request string
}

type SuccessMessages struct {
	Register          string
	PhoneVerification string
	Login             string
	AddAddress        string
	EditAddress       string
	DeleteAddress     string
	ChangePassword    string
	ForgotPassword    string
	RefreshToken      string
	EmailVerification string
	UpdateProfile     string
	CreateRole        string
	UpdateRole        string
	DeleteRole        string
	UpdateUserRole    string
	Generic           string
	AddPet            string
	RemovePet         string
	UpdatePet         string
}

type JWTConstants struct {
	AccessTokenType  string
	RefreshTokenType string
}

type ErrorFields struct {
	User         string
	Phone        string
	Email        string
	Password     string
	MagicLink    string
	Address      string
	Name         string
	Province     string
	City         string
	Role         string
	Permission   string
	BirthDate    string
	IsAdult      string
	Pet          string
	Species      string
	PetSitter    string
	Request      string
	CalendarSlot string
	Service      string
}

type ErrorTags struct {
	AlreadyRegistered      string
	MinimumLength          string
	ContainsLowercase      string
	ContainsUppercase      string
	ContainsNumber         string
	ContainsSpecialChar    string
	Expired                string
	Invalid                string
	NotRegistered          string
	NotVerified            string
	InvalidAuthCredentials string
	ExpiredAuthToken       string
	InvalidAuthToken       string
	Unauthorized           string
	AwaitingApproval       string
	Rejected               string
	NotExist               string
	AlreadyExist           string
	ForbiddenStatus        string
	Pending                string
	NotAccepted            string
	InvalidRecaptcha       string
	Required               string
	Numeric                string
	AccessDenied           string
	Binding                string
	Generic                string
	NotFound               string
	UnacceptableInput      string
	DuplicateName          string
	CalendarConflict       string
	OldInfo                string
}

type JWTKeysPath struct {
	PublicKey  string
	PrivateKey string
}

type Context struct {
	Translator     string
	ID             string
	RefreshToken   string
	AcceptLanguage string
}

type TemplatesPath struct {
	Path                   string
	EmailVerification      string
	NewRequest             string
	PetOwnerRequestCancel  string
	PetSitterRequestCancel string
	RequestAccepted        string
	RequestDeclined        string
	RequestEdited          string
}

func NewConstants() *Constants {
	return &Constants{
		Context: Context{
			Translator:     "translator",
			ID:             "id",
			RefreshToken:   "refreshToken",
			AcceptLanguage: "Accept-Language",
		},
		JWTKeysPath: JWTKeysPath{
			PublicKey:  "./internal/infrastructure/jwt/public_key.pem",
			PrivateKey: "./internal/infrastructure/jwt/private_key.pem",
		},
		ErrorFields: ErrorFields{
			User:         "user",
			Phone:        "phone",
			Email:        "email",
			Password:     "password",
			MagicLink:    "magicLink",
			Address:      "address",
			Name:         "name",
			Province:     "province",
			City:         "city",
			Role:         "role",
			Permission:   "permission",
			BirthDate:    "birthDate",
			IsAdult:      "isAdult",
			Pet:          "pet",
			Species:      "species",
			PetSitter:    "petSitter",
			Request:      "request",
			CalendarSlot: "calendarSlot",
			Service:      "service",
		},
		ErrorTags: ErrorTags{
			AlreadyRegistered:      "errors.alreadyRegistered",
			MinimumLength:          "errors.minimumLength",
			ContainsLowercase:      "errors.containsLowercase",
			ContainsUppercase:      "errors.containsUppercase",
			ContainsNumber:         "errors.containsNumber",
			ContainsSpecialChar:    "errors.containsSpecialChar",
			Expired:                "errors.Expired",
			Invalid:                "errors.invalid",
			NotRegistered:          "errors.notRegistered",
			NotVerified:            "errors.notVerified",
			InvalidAuthCredentials: "errors.invalidAuthCredentials",
			ExpiredAuthToken:       "errors.expiredAuthToken",
			InvalidAuthToken:       "errors.invalidAuthToken",
			Unauthorized:           "errors.unauthorized",
			AwaitingApproval:       "errors.awaitingApproval",
			Rejected:               "errors.rejected",
			NotExist:               "errors.notExist",
			AlreadyExist:           "errors.alreadyExist",
			ForbiddenStatus:        "errors.forbiddenStatus",
			Pending:                "errors.pending",
			NotAccepted:            "errors.notAccepted",
			InvalidRecaptcha:       "errors.invalidRecaptcha",
			Required:               "errors.required",
			Numeric:                "errors.numeric",
			AccessDenied:           "errors.accessDenied",
			Binding:                "errors.binding",
			Generic:                "errors.generic",
			NotFound:               "errors.notFound",
			UnacceptableInput:      "errors.unacceptableInput",
			DuplicateName:          "errors.duplicateName",
			CalendarConflict:       "errors.calendarConflict",
			OldInfo:                "errors.oldInfo",
		},
		JWTConstants: JWTConstants{
			AccessTokenType:  "access",
			RefreshTokenType: "refresh",
		},
		SuccessMessages: SuccessMessages{
			Register:          "successMessage.userRegister",
			PhoneVerification: "successMessage.phoneVerification",
			Login:             "successMessage.login",
			AddAddress:        "successMessage.addAddress",
			EditAddress:       "successMessage.editAddress",
			DeleteAddress:     "successMessage.deleteAddress",
			ChangePassword:    "successMessage.changePassword",
			ForgotPassword:    "successMessage.forgotPassword",
			RefreshToken:      "successMessage.refreshToken",
			EmailVerification: "successMessage.emailVerification",
			UpdateProfile:     "successMessage.updateProfile",
			CreateRole:        "successMessage.createRole",
			UpdateRole:        "successMessage.updateRole",
			DeleteRole:        "successMessage.deleteRole",
			UpdateUserRole:    "successMessage.updateUserRoles",
			Generic:           "successMessage.generic",
			AddPet:            "successMessage.addPet",
			RemovePet:         "successMessage.removePet",
			UpdatePet:         "successMessage.updatePet",
		},
		EntityConstants: EntityConstants{
			Request: "request",
		},
		TemplatesPath: TemplatesPath{
			Path:                   "./internal/infrastructure/mail/",
			EmailVerification:      "email_verification.html",
			NewRequest:             "new_request.html",
			PetOwnerRequestCancel:  "pet_owner_request_cancel.html",
			PetSitterRequestCancel: "pet_sitter_request_cancel.html",
			RequestAccepted:        "request_accepted.html",
			RequestDeclined:        "request_declined.html",
			RequestEdited:          "request_edited.html",
		},
	}
}
