package types

// Errors
const (
	ValidationError     = "Validation failed"
	UserExists          = "User already exists"
	UserCreateFailed    = "User creation failed"
	LoginUserReadFailed = "Login user read failed"
	PasswordHashErr     = "Password hashing failed"
)

// Debug messages
const (
	Deserialized = "Deserializing"
	ReadBody     = "Read request body"
)

// Success messages
const (
	LoginSuccess = "Login success"
)

// JWT handling and context messages
const (
	UserDataKey    = "userData"
	TokenRecieved  = "Token recieved"
	TokenParsed    = "Token parsed"
	JwtCreateError = "Jwt creation failed"
)

// User messages
const (
	UserReadSuccess = "User read success"
	UserReadError   = "User read failed"
)

// Merchandise messages
const (
	MerchMigrationSuccess = "Merchandise migration success"
	MerchMigrationError   = "Merchandise migration failed"
	MerchCreateSuccess    = "Merchandise create successfully"
	MerchCreateError      = "Merchandise creation failed"
	MerchReadSuccess      = "Merchandise read successfully"
	MerchReadError        = "Merchandise read failed"
	MerchUpdateSuccess    = "Merchandise update successfully"
	MerchUpdateError      = "Merchandise update failed"
	MerchDeleteSuccess    = "Merchandise delete successfully"
	MerchDeleteError      = "Merchandise delete failed"
	MerchSerializeError   = "Merchandise serialize failed"
)
