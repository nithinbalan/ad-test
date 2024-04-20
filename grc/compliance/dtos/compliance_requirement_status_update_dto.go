package dtos

type ComplianceRequirementStatusUpdateListDTO struct {
	ID                         uint    `json:"id"`
	ComplianceRequirementTitle string  `json:"compliance_requirement_title"`
	ComplianceStatusTitleEn    string  `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr    string  `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode  string  `json:"compliance_status_color_code"`
	Comment                    string  `json:"comment"`
	CreatedAt                  string  `json:"created_at"`
	UpdatedAt                  string  `json:"updated_at"`
	CreatedUserFirstNameEn     string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr     string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn      string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr      *string `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn     string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr     string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn      string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr      string  `json:"updated_user_last_name_ar"`
}

type ComplianceRequirementStatusUpdateFormDTO struct {
	ComplianceStatusID uint                                                `json:"compliance_status_id" binding:"required"`
	Comment            string                                              `json:"comment" binding:"required"`
	Documents          []*ComplianceRequirementStatusUpdateDocumentFormDTO `json:"documents" binding:"required"`
}

type ComplianceRequirementStatusUpdateDocumentFormDTO struct {
	ID           *uint  `json:"id,omitempty"`
	URL          string `json:"url"`
	Icon         string `json:"icon" binding:"required"`
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
	Size         int    `json:"size"`
	IsNewFile    *bool  `json:"is_new_file,omitempty"` // Optional field for indicating if the file is new or updated
}
