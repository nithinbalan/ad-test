package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
)

type ComplianceRegisterListDTO struct {
	ID                                    uint   `json:"id"`
	ComplianceCategoryTitleEn             string `json:"compliance_category_title_en"`
	ComplianceCategoryTitleAr             string `json:"compliance_category_title_ar"`
	ComplianceStatusTitleEn               string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr               string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode             string `json:"compliance_status_color_code"`
	ResponsibleUnitTitleEn                string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr                string `json:"responsible_unit_title_ar"`
	ResponsibleUserID                     uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn            string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr            string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn             string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr             string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial                string `json:"responsible_user_initial"`
	ResponsibleUserColor                  string `json:"responsible_user_color"`
	ResponsibleUserEmail                  string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath       string `json:"responsible_user_profile_image_path"`
	ReviewFrequencyTitleEn                string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr                string `json:"review_frequency_title_ar"`
	ReviewResponsibleUserID               uint   `json:"review_responsible_user_id"`
	ReviewResponsibleUserFirstNameEn      string `json:"review_responsible_user_first_name_en"`
	ReviewResponsibleUserFirstNameAr      string `json:"review_responsible_user_first_name_ar"`
	ReviewResponsibleUserLastNameEn       string `json:"review_responsible_user_last_name_en"`
	ReviewResponsibleUserLastNameAr       string `json:"review_responsible_user_last_name_ar"`
	ReviewResponsibleUserInitial          string `json:"review_responsible_user_initial"`
	ReviewResponsibleUserColor            string `json:"review_responsible_user_color"`
	ReviewResponsibleUserEmail            string `json:"review_responsible_user_email"`
	ReviewResponsibleUserProfileImagePath string `json:"review_responsible_user_profile_image_path"`
	SeverityLevelTitleEn                  string `json:"severity_level_title_en"`
	SeverityLevelTitleAr                  string `json:"severity_level_title_ar"`
	SeverityLevelColorCode                string `json:"severity_level_color_code"`
	ReferenceNumber                       string `json:"reference_number"`
	StandardCode                          string `json:"standard_code"`
	TitleEn                               string `json:"title_en"`
	TitleAr                               string `json:"title_ar"`
	DescriptionEn                         string `json:"description_en"`
	DescriptionAr                         string `json:"description_ar"`
	Version                               string `json:"version"`
	Deadline                              string `json:"deadline"`
	IssueDate                             string `json:"issue_date"`
	ExpiryDate                            string `json:"expiry_date"`
	NextReviewDate                        string `json:"next_review_date"`
	LastReviewedDate                      string `json:"last_reviewed_date"`
	CreatedAt                             string `json:"created_at"`
	UpdatedAt                             string `json:"updated_at"`
	CreatedUserID                         uint   `json:"created_user_id"`
	CreatedUserFirstNameEn                string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr                string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                 string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                 string `json:"created_user_last_name_ar"`
	CreatedUserInitial                    string `json:"created_user_initial"`
	CreatedUserColor                      string `json:"created_user_color"`
	CreatedUserEmail                      string `json:"created_user_email"`
	CreatedUserProfileImagePath           string `json:"created_user_profile_image_path"`
	UpdatedUserID                         *uint  `json:"updated_user_id"`
	UpdatedUserFirstNameEn                string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr                string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                 string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                 string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                    string `json:"updated_user_initial"`
	UpdatedUserColor                      string `json:"updated_user_color"`
	UpdatedUserEmail                      string `json:"updated_user_email"`
	UpdatedUserProfileImagePath           string `json:"updated_user_profile_image_path"`
	LifecycleStatusTitleEn                string `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr                string `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode              string `json:"lifecycle_status_color_code"`
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
	DateFrom           string
	DateTo             string
	NextReviewDateFrom string
	NextReviewDateTo   string
	Export             bool
}

type ComplianceRegisterListResponseDTO struct {
	Items      []*ComplianceRegisterListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ComplianceRegisterFormDTO struct {
	ComplianceCategoryID     uint                      `json:"compliance_category_id" binding:"required"`
	ComplianceStatusID       uint                      `json:"compliance_status_id" binding:"required"`
	ResponsibleUnitID        uint                      `json:"responsible_unit_id" binding:"required"`
	ResponsibleUserID        uint                      `json:"responsible_user_id" binding:"required"`
	ReviewFrequencyID        *uint                     `json:"review_frequency_id"`
	ReviewResponsibleUserID  *uint                     `json:"review_responsible_user_id"`
	SeverityLevelID          *uint                     `json:"severity_level_id"`
	TitleEn                  string                    `json:"title_en" binding:"required,max=255"`
	TitleAr                  string                    `json:"title_ar" binding:"required,max=255"`
	DescriptionEn            *string                   `json:"description_en"`
	DescriptionAr            *string                   `json:"description_ar"`
	Version                  *string                   `json:"version" binding:"omitempty,max=45"`
	Deadline                 *string                   `json:"deadline" binding:"omitempty,datetime=2006-01-02"`
	IssueDate                *string                   `json:"issue_date" binding:"omitempty,datetime=2006-01-02"`
	ExpiryDate               *string                   `json:"expiry_date" binding:"omitempty,datetime=2006-01-02"`
	NextReviewDate           *string                   `json:"next_review_date" binding:"omitempty,datetime=2006-01-02"`
	LastReviewedDate         *string                   `json:"last_reviewed_date" binding:"omitempty,datetime=2006-01-02"`
	LastReviewedUserID       *uint                     `json:"last_reviewed_user_id" binding:"omitempty"`
	ComplianceReviewStatusID *uint                     `json:"compliance_review_status_id"`
	LastReviewNote           *string                   `json:"last_review_note"`
	StandardCode             *string                   `json:"standard_code" binding:"omitempty,max=45"`
	DocumentEn               ComplianceDocumentFormDTO `json:"document_en" binding:"required"`
	DocumentAr               ComplianceDocumentFormDTO `json:"document_ar" binding:"required"`
	ControlIDs               []*uint                   `json:"control_ids"`
}

type ComplianceDocumentFormDTO struct {
	ID           *uint  `json:"id,omitempty"`
	Language     string `json:"language" binding:"required"`
	URL          string `json:"url"`
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
	Size         int    `json:"size"`
	IsNewFile    *bool  `json:"is_new_file,omitempty"` // Optional field for indicating if the file is new or updated
}

type ComplianceDocumentDTO struct {
	ID           *uint  `json:"id,omitempty"`
	Language     string `json:"language"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
	Size         int    `json:"size"`
}

type ComplianceRegisterParsedDates struct {
	Deadline         *time.Time
	IssueDate        *time.Time
	ExpiryDate       *time.Time
	NextReviewDate   *time.Time
	LastReviewedDate *time.Time
}

type ComplianceRegisterDetailDTO struct {
	ID                      int                                `json:"id"`
	ReferenceNumber         string                             `json:"reference_number"`
	StandardCode            string                             `json:"standard_code"`
	TitleEn                 string                             `json:"title_en"`
	TitleAr                 string                             `json:"title_ar"`
	DescriptionEn           *string                            `json:"description_en"`
	DescriptionAr           *string                            `json:"description_ar"`
	DeadlineString          *string                            `json:"deadline"`
	IssueDateString         *string                            `json:"issue_date"`
	ExpiryDateString        *string                            `json:"expiry_date"`
	NextReviewDateString    *string                            `json:"next_review_date"`
	LastReviewedDateString  *string                            `json:"last_reviewed_date"`
	LastReviewNote          *string                            `json:"last_review_note"`
	LastReviewedUser        *organizationDTOs.UserMinimalDTO   `json:"last_reviewed_user"`
	ComplianceReviewStatus  *ComplianceReviewStatusDTO         `json:"compliance_review_status"`
	ComplianceCategory      ComplianceCategoryDTO              `json:"compliance_category"`
	ComplianceStatus        ComplianceStatusDTO                `json:"compliance_status"`
	ResponsibleUnit         settingsDTOs.OrganizationUnitDTO   `json:"responsible_unit"`
	ResponsibleUser         organizationDTOs.UserMinimalDTO    `json:"responsible_user"`
	ResponsibleUserID       uint                               `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ReviewFrequency         *masterDTOs.ReviewFrequencyDTO     `json:"review_frequency"`
	ReviewResponsibleUser   *organizationDTOs.UserMinimalDTO   `json:"review_responsible_user"`
	SeverityLevel           *masterDTOs.SeverityLevelDTO       `json:"severity_level"`
	Version                 *string                            `json:"version"`
	CreatedAt               time.Time                          `json:"created_at"`
	UpdatedAt               *time.Time                         `json:"updated_at"`
	CreatedByUser           organizationDTOs.UserMinimalDTO    `json:"created_by"`
	UpdatedByUser           *organizationDTOs.UserMinimalDTO   `json:"updated_by"`
	LifecycleStatus         masterDTOs.LifecycleStatusDTO      `json:"status"`
	DocumentEn              *ComplianceDocumentDTO             `json:"document_en"`
	DocumentAr              *ComplianceDocumentDTO             `json:"document_ar"`
	Controls                []*ComplianceControlDTO            `json:"controls"`
	ComplianceStatusUpdates []*ComplianceStatusUpdateDetailDTO `json:"compliance_status_updates"`
}

type ComplianceControlDTO struct {
	ID                   *uint                             `json:"id"`
	ComplianceRegisterID uint                              `json:"compliance_register_id"`
	ControlID            uint                              `json:"control_id"`
	Control              organizationDTOs.ControlDetailDTO `json:"control"`
}
