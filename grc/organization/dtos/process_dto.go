package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	settingsDTOs "jethur/grc/settings/dtos"

	"time"
)

// ProcessListDTO defines the data transfer object for listing processes
type ProcessListDTO struct {
	ID                               uint       `json:"id"`
	ProcessGroupTitleEn              string     `json:"process_group_title_en"`
	ProcessGroupTitleAr              string     `json:"process_group_title_ar"`
	OwnerUnitTitleEn                 *string    `json:"owner_unit_title_en"`
	OwnerUnitTitleAr                 *string    `json:"owner_unit_title_ar"`
	OwnerID                          *uint      `json:"owner_id"`
	OwnerFirstNameEn                 *string    `json:"owner_first_name_en"`
	OwnerFirstNameAr                 *string    `json:"owner_first_name_ar"`
	OwnerLastNameEn                  *string    `json:"owner_last_name_en"`
	OwnerLastNameAr                  *string    `json:"owner_last_name_ar"`
	OwnerInitial                     string     `json:"owner_initial"`
	OwnerColor                       string     `json:"owner_color"`
	OwnerEmail                       string     `json:"owner_email"`
	OwnerProfileImagePath            string     `json:"owner_profile_image_path"`
	LifecycleStatusTitleEn           string     `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr           string     `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode         string     `json:"lifecycle_status_color_code"`
	ReferenceNumber                  string     `json:"reference_number"`
	TitleEn                          string     `json:"title_en"`
	TitleAr                          string     `json:"title_ar"`
	Description                      *string    `json:"description"`
	MarkedAsAuditableItem            bool       `json:"marked_as_auditable_item"`
	Rto                              *uint      `json:"rto"`
	RtoUnitId                        *uint      `json:"rto_uint_id"`
	RtoUnitTitleEn                   *string    `json:"rto_unit_title_en"`
	RtoUnitTitleAr                   *string    `json:"rto_unit_title_ar"`
	Mtpd                             *uint      `json:"mtpd"`
	MtpdUnitId                       *uint      `json:"mtpd_uint_id"`
	MtpdUnitTitleEn                  *string    `json:"mtpd_unit_title_en"`
	MtpdUnitTitleAr                  *string    `json:"mtpd_unit_title_ar"`
	Rpo                              *uint      `json:"rpo"`
	RpoUnitId                        *uint      `json:"rpo_uint_id"`
	RpoUnitTitleEn                   *string    `json:"rpo_unit_title_en"`
	RpoUnitTitleAr                   *string    `json:"rpo_unit_title_ar"`
	RiskLevelId                      *uint      `json:"risk_level_id"`
	RiskLevelTitleEn                 *uint      `json:"risk_level_title_en"`
	RiskLevelTitleAr                 *uint      `json:"risk_level_title_ar"`
	LastAssessedAt                   *time.Time `json:"last_assessed_at"`
	LastAssessedUserID               *uint      `json:"last_assessed_user_id"`
	LastAssessedUserFirstNameEn      *string    `json:"last_assessed_user_first_name_en"`
	LastAssessedUserFirstNameAr      *string    `json:"last_assessed_user_first_name_ar"`
	LastAssessedUserLastNameEn       *string    `json:"last_assessed_user_last_name_en"`
	LastAssessedUserLastNameAr       *string    `json:"last_assessed_user_last_name_ar"`
	LastAssessedUserInitial          *string    `json:"last_assessed_user_initial"`
	LastAssessedUserColor            *string    `json:"last_assessed_user_color"`
	LastAssessedUserEmail            *string    `json:"last_assessed_user_email"`
	LastAssessedUserProfileImagePath *string    `json:"last_assessed_user_profile_image_path"`
	IsBiaPerformed                   bool       `json:"is_bia_performed"`
	CreatedUserID                    uint       `json:"created_user_id"`
	CreatedUserFirstNameEn           string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr           string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn            string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr            string     `json:"created_user_last_name_ar"`
	CreatedUserInitial               string     `json:"created_user_initial"`
	CreatedUserColor                 string     `json:"created_user_color"`
	CreatedUserEmail                 string     `json:"created_user_email"`
	CreatedUserProfileImagePath      string     `json:"created_user_profile_image_path"`
	UpdatedUserID                    *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn           *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr           *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn            *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr            *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial               string     `json:"updated_user_initial"`
	UpdatedUserColor                 string     `json:"updated_user_color"`
	UpdatedUserEmail                 string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath      string     `json:"updated_user_profile_image_path"`
	CreatedAt                        time.Time  `json:"created_at"`
	UpdatedAt                        *time.Time `json:"updated_at"`
	RiskRatingTitleEn                *string    `json:"risk_rating_title_en"`
	RiskRatingTitleAr                *string    `json:"risk_rating_title_ar"`
	RiskRatingColorCode              *string    `json:"risk_rating_color_code"`
	RiskRatingValue                  *uint      `json:"risk_rating_value"`
}

// ProcessListQueryParamsDTO defines the filters and search parameters for querying processes
type ProcessListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
	DateFrom    string
	DateTo      string
}

// ProcessListResponseDTO encapsulates the paginated response for a list of processes
type ProcessListResponseDTO struct {
	Items      []*ProcessListDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

// ProcessFormDTO defines the structure for creating or updating a process
type ProcessFormDTO struct {
	ProcessGroupID   uint    `json:"process_group_id" binding:"required"`
	OwnerID          uint    `json:"owner_id" binding:"required"`
	TitleEn          string  `json:"title_en" binding:"required"`
	TitleAr          string  `json:"title_ar" binding:"required"`
	Description      *string `json:"description" binding:"omitempty"`
	ControlIds       []*uint `json:"control_ids"`
	AssetRegisterIds []*uint `json:"asset_register_ids"`
	DocumentIds      []*uint `json:"document_ids"`
}

type ProcessDetailDTO struct {
	ID                    uint                              `json:"id"`
	ReferenceNumber       string                            `json:"reference_number"`
	TitleEn               string                            `json:"title_en"`
	TitleAr               string                            `json:"title_ar"`
	Description           *string                           `json:"description"`
	ProcessGroup          *ProcessGroupDTO                  `json:"process_group"`
	OwnerUnit             *settingsDTOs.OrganizationUnitDTO `json:"owner_unit"`
	OwnerID               *uint                             `json:"owner_id"`
	Owner                 *UserMinimalDTO                   `json:"owner"`
	CreatedByUser         UserMinimalDTO                    `json:"created_by"`
	UpdatedByUser         *UserMinimalDTO                   `json:"updated_by"`
	LifecycleStatus       masterDTOs.LifecycleStatusDTO     `json:"status"`
	Rto                   *uint                             `json:"rto"`
	RtoUnitId             *uint                             `json:"rto_uint_id"`
	RtoUnit               masterDTOs.TimeUnitDTO            `json:"rto_unit"`
	Mtpd                  *uint                             `json:"mtpd"`
	MtpdUnitId            *uint                             `json:"mtpd_uint_id"`
	MtpdUnit              masterDTOs.TimeUnitDTO            `json:"mtpd_unit"`
	Rpo                   *uint                             `json:"rpo"`
	RpoUnitId             *uint                             `json:"rpo_uint_id"`
	RpoUnit               masterDTOs.TimeUnitDTO            `json:"rpo_unit"`
	RiskLevelId           *uint                             `json:"risk_level_id"`
	LastAssessedAt        *time.Time                        `json:"last_assessed_at"`
	LastAssessedByUser    *UserMinimalDTO                   `json:"last_assessed_by"`
	IsBiaPerformed        bool                              `json:"is_bia_performed"`
	CreatedAt             time.Time                         `json:"created_at"`
	UpdatedAt             *time.Time                        `json:"updated_at"`
	Controls              []*ProcessControlDTO              `json:"controls"`
	Assets                []*ProcessAssetDTO                `json:"assets"`
	Documents             []*ProcessAssociatedDocumentDTO   `json:"documents"`
	Risks                 []*ProcessRiskDTO                 `json:"risks"`
	MarkedAsAuditableItem bool                              `json:"marked_as_auditable_item"`
	RiskMatrixRating      *RiskMatrixRiskRatingDTO          `json:"risk_matrix_risk_rating"`
	RiskMatrixRatingValue *RiskMatrixRiskRatingValueDTO     `json:"risk_matrix_rating_value"`
	AuditableItem         *AuditableItemMinimalDetailDTO    `json:"auditable_item"`
}
type AuditableItemMinimalDetailDTO struct {
	ID                  uint                             `json:"id"`
	Type                string                           `json:"type"`
	TitleEn             *string                          `json:"title_en"`
	TitleAr             *string                          `json:"title_ar"`
	Comment             *string                          `json:"comment"`
	ResponsibleUnitID   uint                             `json:"responsible_unit_id"`
	ResponsibleUserID   uint                             `json:"responsible_user_id"`
	CreatedByID         uint                             `json:"created_by_id"`
	UpdatedByID         *uint                            `json:"updated_by_id"`
	AuditorRiskRatingID uint                             `json:"auditor_risk_rating_id"`
	AuditorRiskRating   *AuditorRiskRatingMinimalDTO     `json:"auditor_risk_rating"`
	ResponsibleUnit     settingsDTOs.OrganizationUnitDTO `json:"responsible_unit"`
	CreatedAt           time.Time                        `json:"created_at"`
	UpdatedAt           *time.Time                       `json:"updated_at"`
	ResponsibleUser     *UserMinimalDTO                  `json:"responsible_user"`
	CreatedByUser       *UserMinimalDTO                  `json:"created_by"`
	UpdatedByUser       *UserMinimalDTO                  `json:"updated_by"`
}
type AuditorRiskRatingMinimalDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	Status    string `json:"status"`
}
type RiskMatrixRiskRatingDTO struct {
	ID           uint   `json:"id"`
	RiskMatrixID uint   `json:"risk_matrix_id"`
	TitleEn      string `json:"title_en"`
	TitleAr      string `json:"title_ar"`
	ColorCode    string `json:"color_code"`
	Type         string `json:"type"`
	Status       string `json:"status"`
}

type ProcessControlDTO struct {
	ID        *uint            `json:"id"`
	ProcessID uint             `json:"process_id"`
	ControlID uint             `json:"control_id"`
	Control   ControlDetailDTO `json:"control"`
}

type ProcessRiskDTO struct {
	ID             *uint          `json:"id"`
	ProcessID      uint           `json:"process_id"`
	RiskRegisterID uint           `json:"risk_register_id"`
	Risk           RiskMinimalDTO `json:"risk"`
}

type RiskMinimalDTO struct {
	ID                       uint                          `json:"id"`
	Title                    string                        `json:"title"`
	Description              *string                       `json:"description"`
	CurrentRiskRatingID      *uint                         `json:"current_risk_rating_id"`
	CurrentRiskRatingValueID *uint                         `json:"current_risk_rating_value_id"`
	CurrentRiskRatingValue   *RiskMatrixRiskRatingValueDTO `json:"current_risk_rating_value"`
}

type RiskMatrixRiskRatingValueDTO struct {
	ID                     uint `json:"id"`
	Value                  int  `json:"value"` // This is used to find the maximum rating
	RiskMatrixRiskRatingID uint `json:"risk_matrix_risk_rating_id"`
}

type ProcessAssociatedDocumentDTO struct {
	ID         *uint              `json:"id"`
	ProcessID  uint               `json:"process_id"`
	DocumentID uint               `json:"document_id"`
	Document   DocumentDetailsDTO `json:"document"`
}

type DocumentDetailsDTO struct {
	ID            int     `json:"id"`
	TitleEn       string  `json:"title_en"`
	TitleAr       string  `json:"title_ar"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
	IsLocked      bool    `json:"is_locked"`
	Acknowledged  *bool   `json:"acknowledged,omitempty"`
	Rto           *uint   `json:"rto"`
	RtoUnitId     *uint   `json:"rto_uint_id"`
	Mtpd          *uint   `json:"mtpd"`
	MtpdUnitId    *uint   `json:"mtpd_uint_id"`
	Rpo           *uint   `json:"rpo"`
	RpoUnitId     *uint   `json:"rpo_uint_id"`
}

type ProcessAssetDTO struct {
	ID              *uint           `json:"id"`
	ProcessID       uint            `json:"process_id"`
	AssetRegisterID uint            `json:"asset_register_id"`
	Asset           AssetMinimalDTO `json:"asset"`
}

type AssetMinimalDTO struct {
	ID          uint    `json:"id"`
	TitleEn     string  `json:"title_en"`
	TitleAr     string  `json:"title_ar"`
	Description *string `json:"description"`
	Rto         *uint   `json:"rto"`
	RtoUnitId   *uint   `json:"rto_uint_id"`
	Mtpd        *uint   `json:"mtpd"`
	MtpdUnitId  *uint   `json:"mtpd_uint_id"`
	Rpo         *uint   `json:"rpo"`
	RpoUnitId   *uint   `json:"rpo_uint_id"`
}
