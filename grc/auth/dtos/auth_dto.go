package dtos

type AuthFormDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ForgotPasswordFormDTO struct {
	Email string `json:"email" binding:"required,email"`
}
