package user

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Personal struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
