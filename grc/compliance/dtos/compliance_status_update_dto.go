package dtos

type ComplianceStatusUpdateListDTO struct {
	ID                        uint    `json:"id"`
	ComplianceTitle           string  `json:"potential_risk_title"`
	ComplianceDescription     string  `json:"potential_risk_description"`
	ComplianceStatusTitleEn   string  `json:"potential_risk_status_title_en"`
	ComplianceStatusTitleAr   string  `json:"potential_risk_status_title_ar"`
	ComplianceStatusColorCode string  `json:"potential_risk_status_color_code"`
	Comment                   string  `json:"comment"`
	CreatedAt                 string  `json:"created_at"`
	CreatedUserFirstNameEn    string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr    string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn     string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr     *string `json:"created_user_last_name_ar"`
}

type ComplianceStatusUpdateFormDTO struct {
	ComplianceStatusID uint   `json:"potential_risk_status_id" binding:"required"`
	Comment            string `json:"comment" binding:"required"`
}
