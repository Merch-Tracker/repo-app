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
