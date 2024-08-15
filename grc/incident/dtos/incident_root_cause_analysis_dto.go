package dtos

type IncidentRootCauseAnalysisListDTO struct {
	ID                     uint    `json:"id"`
	IncidentRegisterTitle  string  `json:"incident_register_title" binding:"required"`
	Why                    uint    `json:"why" binding:"required"`
	Question               string  `json:"question" binding:"required"`
	Answer                 string  `json:"answer,omitempty"`
	CreatedAt              string  `json:"created_at"`
	UpdatedAt              *string `json:"updated_at"`
	CreatedUserFirstNameEn string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  string  `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn  *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr  *string `json:"updated_user_last_name_ar"`
}

type IncidentRootCauseAnalysisListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
}
type IncidentRootCauseAnalysisFormDTO struct {
	Question string `json:"question" binding:"required"`
	Answer   string `json:"answer" binding:"required"`
}
