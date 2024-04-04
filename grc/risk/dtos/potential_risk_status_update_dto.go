package dtos

type PotentialRiskStatusUpdateListDTO struct {
	ID                           uint    `json:"id"`
	PotentialRiskTitle           string  `json:"potential_risk_title"`
	PotentialRiskDescription     string  `json:"potential_risk_description"`
	PotentialRiskStatusTitleEn   string  `json:"potential_risk_status_title_en"`
	PotentialRiskStatusTitleAr   string  `json:"potential_risk_status_title_ar"`
	PotentialRiskStatusColorCode string  `json:"potential_risk_status_color_code"`
	Comment                      string  `json:"comment"`
	CreatedAt                    string  `json:"created_at"`
	CreatedUserFirstNameEn       string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr       string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn        string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr        *string `json:"created_user_last_name_ar"`
}

type PotentialRiskStatusUpdateFormDTO struct {
	PotentialRiskStatusID uint   `json:"potential_risk_status_id" binding:"required"`
	Comment               string `json:"comment" binding:"required"`
}
