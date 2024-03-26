package dtos

type AuthFormDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type PasswordRecoveryFormDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type PasswordRecoveryDataDTO struct {
	FirstNameEn string `json:"first_name_en"`
	LastNameEn  string `json:"last_name_en"`
	Otp         string `json:"otp"`
}

type PasswordRecoveryOTPVerifyFormDTO struct {
	Email string `json:"email" binding:"required,email"`
	Otp   string `json:"otp" binding:"required"`
}

type PasswordRecoverySetPasswordFormDTO struct {
	Token           string `json:"token"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}
