package dtos

type AuthFormDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type PasswordRecoveryFormDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordDataDTO struct {
	FirstNameEN string `json:"first_name_en"`
	LastNameEN  string `json:"last_name_en"`
	Otp         string `json:"otp"`
}
