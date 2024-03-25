package dtos

type AuthFormDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type PasswordRecoveryFormDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordDataDTO struct {
	FirstNameEn string `json:"first_name_en"`
	LastNameEn  string `json:"last_name_en"`
	Otp         string `json:"otp"`
}
