package bootstrap

type Constants struct {
	Context      Context
	JWTKeysPath  JWTKeysPath
	ErrorFields  ErrorFields
	ErrorTags    ErrorTags
	JWTConstants JWTConstants
}

type JWTConstants struct {
	AccessTokenType  string
	RefreshTokenType string
}

type ErrorFields struct {
	User       string
	Phone      string
	Email      string
	Password   string
	OTP        string
	Address    string
	Name       string
	Province   string
	City       string
	Role       string
	Permission string
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
	NotActive              string
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
	AlreadyBlocked         string
	AlreadyActive          string
	AlreadyResolved        string
	AlreadyArchived        string
	AlreadyCompleted       string
	NotAccepted            string
	StatusNotChange        string
	AlreadyCanceled        string
	AlreadyRejected        string
	AlreadyAccepted        string
	AlreadyDraft           string
	InvalidRecaptcha       string
	Required               string
	Numeric                string
	AccessDenied           string
	Binding                string
	Generic                string
	NotFound               string
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
			User:       "user",
			Phone:      "phone",
			Email:      "email",
			Password:   "password",
			OTP:        "otp",
			Address:    "address",
			Name:       "name",
			Province:   "province",
			City:       "city",
			Role:       "role",
			Permission: "permission",
		},
		ErrorTags: ErrorTags{
			AlreadyRegistered:      "alreadyRegistered",
			MinimumLength:          "minimumLength",
			ContainsLowercase:      "containsLowercase",
			ContainsUppercase:      "containsUppercase",
			ContainsNumber:         "containsNumber",
			ContainsSpecialChar:    "containsSpecialChar",
			Expired:                "Expired",
			Invalid:                "invalid",
			NotRegistered:          "notRegistered",
			NotVerified:            "notVerified",
			NotActive:              "notActive",
			InvalidAuthCredentials: "invalidAuthCredentials",
			ExpiredAuthToken:       "expiredAuthToken",
			InvalidAuthToken:       "invalidAuthToken",
			Unauthorized:           "unauthorized",
			AwaitingApproval:       "awaitingApproval",
			Rejected:               "rejected",
			NotExist:               "notExist",
			AlreadyExist:           "alreadyExist",
			ForbiddenStatus:        "forbiddenStatus",
			Pending:                "pending",
			AlreadyBlocked:         "alreadyBlocked",
			AlreadyActive:          "alreadyActive",
			AlreadyResolved:        "alreadyResolved",
			AlreadyArchived:        "alreadyArchived",
			AlreadyCompleted:       "alreadyCompleted",
			NotAccepted:            "notAccepted",
			StatusNotChange:        "statusNotChange",
			AlreadyCanceled:        "alreadyCanceled",
			AlreadyRejected:        "alreadyRejected",
			AlreadyAccepted:        "alreadyAccepted",
			AlreadyDraft:           "alreadyDraft",
			InvalidRecaptcha:       "invalidRecaptcha",
			Required:               "required",
			Numeric:                "numeric",
			AccessDenied:           "accessDenied",
			Binding:                "binding",
			Generic:                "generic",
			NotFound:               "notFound",
		},
		JWTConstants: JWTConstants{
			AccessTokenType:  "access",
			RefreshTokenType: "refresh",
		},
	}
}
