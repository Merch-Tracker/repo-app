package types

// ValidationError Errors
const (
	ValidationError     = "Validation failed"
	UserExists          = "User already exists"
	UserCreateFailed    = "User creation failed"
	LoginUserReadFailed = "Login user read failed"
	JwtCreateError      = "Jwt creation failed"
)

// Errors
const (
	PasswordHashErr = "Password hashing failed"
)

// ReadBody Debug messages
const (
	Deserialized = "Deserializing"
	ReadBody     = "Read request body"
)

// LoginSuccess Success messages
const (
	LoginSuccess = "Login success"
)
