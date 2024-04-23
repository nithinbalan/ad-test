package dtos

type KRIValueDTO struct {
	ID                     uint    `json:"id"`
	KRIRegisterID          uint    `json:"kri_register_id"`
	ReviewDate             string  `json:"review_date"`
	Value                  string  `json:"value"`
	LowThreshold           string  `json:"low_threshold"`
	MediumThreshold        string  `json:"medium_threshold"`
	HighThreshold          string  `json:"high_threshold"`
	Remarks                string  `json:"remarks"`
	CreatedAt              string  `json:"created_at"`
	UpdatedAt              string  `json:"updated_at"`
	CreatedUserFirstNameEn string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  *string `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn  *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr  *string `json:"updated_user_last_name_ar"`
}
