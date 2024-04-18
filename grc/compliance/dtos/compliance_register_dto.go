package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	"time"
)

type ComplianceRegisterListDTO struct {
	ID                               uint   `json:"id"`
	ComplianceCategoryTitleEn        string `json:"compliance_category_title_en"`
	ComplianceCategoryTitleAr        string `json:"compliance_category_title_ar"`
	ComplianceStatusTitleEn          string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr          string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode        string `json:"compliance_status_color_code"`
	ResponsibleUnitTitleEn           string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr           string `json:"responsible_unit_title_ar"`
	ResponsibleUserFirstNameEn       string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr       string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn        string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr        string `json:"responsible_user_last_name_ar"`
	ReviewFrequencyTitleEn           string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr           string `json:"review_frequency_title_ar"`
	ReviewResponsibleUserFirstNameEn string `json:"review_responsible_user_first_name_en"`
	ReviewResponsibleUserFirstNameAr string `json:"review_responsible_user_first_name_ar"`
	ReviewResponsibleUserLastNameEn  string `json:"review_responsible_user_last_name_en"`
	ReviewResponsibleUserLastNameAr  string `json:"review_responsible_user_last_name_ar"`
	SeverityLevelTitleEn             string `json:"severity_level_title_en"`
	SeverityLevelTitleAr             string `json:"severity_level_title_ar"`
	SeverityLevelColorCode           string `json:"severity_level_color_code"`
	ReferenceNumber                  string `json:"reference_number"`
	StandardCode                     string `json:"standard_code"`
	TitleEn                          string `json:"title_en"`
	TitleAr                          string `json:"title_ar"`
	DescriptionEn                    string `json:"description_en"`
	DescriptionAr                    string `json:"description_ar"`
	Version                          string `json:"version"`
	Deadline                         string `json:"deadline"`
	IssueDate                        string `json:"issue_date"`
	ExpiryDate                       string `json:"expiry_date"`
	NextReviewDate                   string `json:"next_review_date"`
	CreatedAt                        string `json:"created_at"`
	UpdatedAt                        string `json:"updated_at"`
	CreatedUserFirstNameEn           string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr           string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn            string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr            string `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn           string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr           string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn            string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr            string `json:"updated_user_last_name_ar"`
	LifecycleStatusTitleEn           string `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr           string `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode         string `json:"lifecycle_status_color_code"`
}

type ComplianceRegisterListQueryParamsDTO struct {
	Sort               string
	SortBy             string
	Keyword            string
	Status             string
	Page               int
	PageSize           int
	Filters            map[string][]string
	QuickFilter        []string
	DeadlineFrom       string
	DeadlineTo         string
	IssueDateFrom      string
	IssueDateTo        string
	ExpiryDateFrom     string
	ExpiryDateTo       string
	NextReviewDateFrom string
	NextReviewDateTo   string
}

type ComplianceRegisterListResponseDTO struct {
	Items      []*ComplianceRegisterListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ComplianceRegisterFormDTO struct {
	ComplianceCategoryID    uint                       `json:"compliance_category_id" binding:"required"`
	ComplianceStatusID      uint                       `json:"compliance_status_id" binding:"required"`
	ResponsibleUnitID       uint                       `json:"responsible_unit_id" binding:"required"`
	ResponsibleUserID       uint                       `json:"responsible_user_id" binding:"required"`
	ReviewFrequencyID       *uint                      `json:"review_frequency_id"`
	ReviewResponsibleUserID *uint                      `json:"review_responsible_user_id"`
	SeverityLevelID         *uint                      `json:"severity_level_id"`
	TitleEn                 string                     `json:"title_en" binding:"required,max=255"`
	TitleAr                 string                     `json:"title_ar" binding:"required,max=255"`
	DescriptionEn           *string                    `json:"description_en"`
	DescriptionAr           *string                    `json:"description_ar"`
	Version                 *string                    `json:"version" binding:"omitempty,max=45"`
	Deadline                string                     `json:"deadline" binding:"omitempty,datetime=2006-01-02"`
	IssueDate               string                     `json:"issue_date" binding:"omitempty,datetime=2006-01-02"`
	ExpiryDate              string                     `json:"expiry_date" binding:"omitempty,datetime=2006-01-02"`
	NextReviewDate          string                     `json:"next_review_date" binding:"omitempty,datetime=2006-01-02"`
	StandardCode            *string                    `json:"standard_code" binding:"omitempty,max=45"`
	Document                *ComplianceDocumentFormDTO `json:"document" binding:"required"`
}
type ComplianceDocumentFormDTO struct {
	ID           *uint  `json:"id,omitempty"`
	Language     string `json:"language" binding:"required"`
	URL          string `json:"url"`
	Icon         string `json:"icon" binding:"required"`
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
	Size         int    `json:"size"`
	IsNewFile    *bool  `json:"is_new_file,omitempty"` // Optional field for indicating if the file is new or updated
}
type ComplianceDocumentDTO struct {
	ID           int    `json:"id"`
	Language     string `json:"language"`
	URL          string `json:"url"`
	Icon         string `json:"icon"`
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
	Size         int    `json:"size"`
}

type ComplianceRegisterParsedDates struct {
	Deadline       *time.Time
	IssueDate      *time.Time
	ExpiryDate     *time.Time
	NextReviewDate *time.Time
}

type OrganizationUnitDTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

// @todo moved to userdto
type UserMinimalDTO struct {
	ID                int                 `json:"id"`
	FirstNameEN       string              `json:"first_name_en"`
	LastNameEN        string              `json:"last_name_en"`
	FirstNameAR       string              `json:"first_name_ar"`
	LastNameAR        string              `json:"last_name_ar"`
	OrganizationUnit1 OrganizationUnitDTO `json:"organization_unit_1"`
}

type ComplianceRegisterDetailDTO struct {
	ID                    int                           `json:"id"`
	ReferenceNumber       string                        `json:"reference_number"`
	StandardCode          string                        `json:"standard_code"`
	TitleEn               string                        `json:"title_en"`
	TitleAr               string                        `json:"title_ar"`
	DescriptionEn         string                        `json:"description_en"`
	DescriptionAr         string                        `json:"description_ar"`
	Deadline              *time.Time                    `json:"deadline"`
	IssueDate             *time.Time                    `json:"issue_date"`
	ExpiryDate            *time.Time                    `json:"expiry_date"`
	NextReviewDate        *time.Time                    `json:"next_review_date"`
	ComplianceCategory    ComplianceCategoryDTO         `json:"compliance_category"`
	ComplianceStatus      ComplianceStatusDTO           `json:"compliance_status"`
	ResponsibleUnit       OrganizationUnitDTO           `json:"responsible_unit"`
	ResponsibleUser       UserMinimalDTO                `json:"responsible_user"`
	ReviewFrequency       masterDTOs.ReviewFrequencyDTO `json:"review_frequency"`
	ReviewResponsibleUser UserMinimalDTO                `json:"review_responsible_user"`
	SeverityLevel         *masterDTOs.SeverityLevelDTO  `json:"severity_level"`
	CreatedAt             time.Time                     `json:"created_at"`
	UpdatedAt             *time.Time                    `json:"updated_at"`
	CreatedByUser         UserMinimalDTO                `json:"created_by"`
	UpdatedByUser         UserMinimalDTO                `json:"updated_by"`
	LifecycleStatus       masterDTOs.LifecycleStatusDTO `json:"status"`
	ComplianceDocuments   []ComplianceDocumentDTO       `json:"compliance_documents"`
}
