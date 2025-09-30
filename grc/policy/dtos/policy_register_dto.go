package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	documentsDTOs "jethur/grc/documents/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type PolicyRegisterListDTO struct {
	ID                              uint       `json:"id"`
	DocumentTypeID                  uint       `json:"document_type_id"`
	TitleEn                         string     `json:"title_en"`
	TitleAr                         string     `json:"title_ar"`
	DocumentCategoryTitleEn         string     `json:"document_category_title_en"`
	DocumentCategoryTitleAr         string     `json:"document_category_title_ar"`
	DocumentTypeTitleEn             string     `json:"document_type_title_en"`
	DocumentTypeTitleAr             string     `json:"document_type_title_ar"`
	DocumentAccessTypeTitleEn       string     `json:"document_access_type_title_en"`
	DocumentAccessTypeTitleAr       string     `json:"document_access_type_title_ar"`
	ResponsibleUnitTitleEn          string     `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string     `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint       `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string     `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string     `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string     `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string     `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string     `json:"responsible_user_initial"`
	ResponsibleUserColor            string     `json:"responsible_user_color"`
	ResponsibleUserEmail            string     `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string     `json:"responsible_user_profile_image_path"`
	DocumentStatusTitleEn           string     `json:"document_status_title_en"`
	DocumentStatusTitleAr           string     `json:"document_status_title_ar"`
	DocumentStatusColorCode         string     `json:"document_status_color_code"`
	DocumentStatusType              string     `json:"document_status_type"`
	DocumentVersion                 float64    `json:"document_version"`
	DocumentURL                     string     `json:"document_url"`
	DocumentName                    string     `json:"document_name"`
	DocumentSize                    int        `json:"document_size"`
	DocumentToken                   string     `json:"document_token"`
	DocumentOriginalName            *string    `json:"document_original_name"`
	DocumentType                    string     `json:"document_type"`
	IsFolder                        *bool      `json:"is_folder"`
	ParentDocumentID                *uint      `json:"parent_document_id"`
	DocumentCount                   *int       `json:"document_count"`
	CreatedUserID                   uint       `json:"created_user_id"`
	CreatedUserFirstNameEn          string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string     `json:"created_user_last_name_ar"`
	CreatedUserInitial              string     `json:"created_user_initial"`
	CreatedUserColor                string     `json:"created_user_color"`
	CreatedUserEmail                string     `json:"created_user_email"`
	CreatedUserProfileImagePath     string     `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn          string     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          string     `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           string     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           string     `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string     `json:"updated_user_initial"`
	UpdatedUserColor                string     `json:"updated_user_color"`
	UpdatedUserEmail                string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string     `json:"updated_user_profile_image_path"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedAt                       *time.Time `json:"updated_at"`
	Acknowledged                    *bool      `json:"acknowledged,omitempty"`
	PublishedUserFirstNameEn        string     `json:"published_user_first_name_en"`
	PublishedUserFirstNameAr        string     `json:"published_user_first_name_ar"`
	PublishedUserLastNameEn         string     `json:"published_user_last_name_en"`
	PublishedUserLastNameAr         string     `json:"published_user_last_name_ar"`
	PublishedUserInitial            string     `json:"published_user_initial"`
	PublishedUserColor              string     `json:"published_user_color"`
	PublishedUserEmail              string     `json:"published_user_email"`
	PublishedUserProfileImagePath   string     `json:"published_user_profile_image_path"`
	PublishedAt                     *time.Time `json:"published_at"`
	UserAccessTypes                 *string    `json:"user_access_types"`
}

type PolicyRegisterListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type PolicyRegisterListResponseDTO struct {
	Items      []*PolicyRegisterListDTO     `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type PolicyRegisterFormDTO struct {
	TitleEn              string  `json:"title_en" binding:"required"`
	TitleAr              string  `json:"title_ar" binding:"required"`
	DescriptionEn        *string `json:"description_en"`
	DescriptionAr        *string `json:"description_ar"`
	DocumentTypeID       uint    `json:"document_type_id" binding:"required"`
	DocumentCategoryID   uint    `json:"document_category_id" binding:"required"`
	ResponsibleUserID    uint    `json:"responsible_user_id" binding:"required"`
	DocumentAccessTypeID uint    `json:"document_access_type_id" binding:"required"`
	DocumentVersion      float64 `json:"document_version" binding:"required"`
	SharedUserIDs        []*uint `json:"shared_user_ids"`
	SharedUnitIDs        []*uint `json:"shared_unit_ids"`
	PolicyFolderID       *uint   `json:"policy_folder_id"`
	ParentDocumentID     *uint   `json:"parent_document_id"`
	ReviewFrequencyID    *uint   `json:"review_frequency_id" binding:"omitempty"`
}

type PolicyRegisterDetailsDTO struct {
	ID                    int                                    `json:"id"`
	ReferenceNumber       string                                 `json:"reference_number"`
	TitleEn               string                                 `json:"title_en"`
	TitleAr               string                                 `json:"title_ar"`
	DescriptionEn         *string                                `json:"description_en"`
	DescriptionAr         *string                                `json:"description_ar"`
	IsLocked              bool                                   `json:"is_locked"`
	DocumentType          documentsDTOs.DocumentTypeDTO          `json:"document_type"`
	DocumentCategory      documentsDTOs.DocumentCategoryDTO      `json:"document_category"`
	DocumentAccessType    documentsDTOs.DocumentAccessTypeDTO    `json:"document_access_type"`
	ResponsibleUnit       settingsDTOs.OrganizationUnitDTO       `json:"responsible_unit"`
	ResponsibleUser       organizationDTOs.UserMinimalDTO        `json:"responsible_user"`
	LatestDocumentVersion DocumentVersionDTO                     `json:"latest_document_version"`
	DocumentVersions      []DocumentVersionDTO                   `json:"document_versions"`
	DocumentStatus        documentsDTOs.DocumentStatusDTO        `json:"document_status"`
	DocumentSharedUsers   []*documentsDTOs.DocumentSharedUserDTO `json:"document_shared_users"`
	DocumentSharedUnits   []*documentsDTOs.DocumentSharedUnitDTO `json:"document_shared_units"`
	CreatedByUser         organizationDTOs.UserMinimalDTO        `json:"created_by"`
	UpdatedByUser         *organizationDTOs.UserMinimalDTO       `json:"updated_by"`
	CreatedAt             time.Time                              `json:"created_at"`
	UpdatedAt             *time.Time                             `json:"updated_at"`
	Acknowledged          *bool                                  `json:"acknowledged,omitempty"`
	ReviewFrequency       *masterDTOs.ReviewFrequencyDTO         `json:"review_frequency"`
	NextReviewDate        *time.Time                             `json:"next_review_date"`
}

type DocumentVersionDTO struct {
	ID                    uint                             `json:"id"`
	DocumentID            uint                             `json:"document_id"`
	Version               float64                          `json:"version"`
	IsLatestVersion       bool                             `json:"is_latest_version"`
	URL                   string                           `json:"url"`
	Name                  string                           `json:"name"`
	Size                  int                              `json:"size"`
	Token                 string                           `json:"token"`
	OriginalName          *string                          `json:"original_name,omitempty"`
	Type                  string                           `json:"type"`
	NextReviewLevel       *int                             `json:"next_review_level"`
	SubmittedByUser       *organizationDTOs.UserMinimalDTO `json:"submitted_by"`
	SubmittedAt           *time.Time                       `json:"submitted_at"`
	AcknowledgeDueDate    *time.Time                       `json:"acknowledge_due_date"`
	EnableAcknowledgement bool                             `json:"enable_acknowledgement"`
	CreatedByUser         organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser         organizationDTOs.UserMinimalDTO  `json:"updated_by"`
	CreatedAt             time.Time                        `json:"created_at"`
	UpdatedAt             *time.Time                       `json:"updated_at"`
	PublishedAt           *time.Time                       `json:"published_at"`
	PublishedByUser       organizationDTOs.UserMinimalDTO  `json:"published_by"`
}

type PolicyPDFData struct {
	Document          *PolicyRegisterDetailsDTO                     `json:"document"`
	Workflow          *documentsDTOs.DocumentVersionWorkflowListDTO `json:"workflow"`
	Clauses           []*PolicyClauseListDTO                        `json:"clauses"`
	CreatedAtString   string                                        `json:"created_at"`
	PublishedAtString string                                        `json:"published_at"`
}
