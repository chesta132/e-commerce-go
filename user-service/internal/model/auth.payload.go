package model

type SignupPayload struct {
	Email      string `validate:"email,required" json:"email"`
	Password   string `validate:"required" json:"password"`
	RememberMe bool   `json:"rememberMe"`
	FullName   string `validate:"required" json:"fullName"`
}

type SigninPayload struct {
	Email      string `validate:"email,required" json:"email"`
	Password   string `validate:"required" json:"password"`
	RememberMe bool   `json:"rememberMe"`
}
