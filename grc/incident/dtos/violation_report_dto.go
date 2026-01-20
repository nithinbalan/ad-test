package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"

	"time"
)

type ViolationReportListDTO struct {
	ID                            uint      `json:"id"`
	ReferenceNumber               string    `json:"reference_number"`
	ReportedAt                    time.Time `json:"reported_at"`
	ViolationReportChannelID      uint      `json:"violation_report_channel_id"`
	ViolationReportChannelTitleEn string    `json:"violation_report_channel_title_en"`
	ViolationReportChannelTitleAr string    `json:"violation_report_channel_title_ar"`
	// IDOfReporter                      *uint      `json:"id_of_reporter"`
	// NameOfReporter                    *string    `json:"name_of_reporter"`
	// PhoneOfReporter                   *string    `json:"phone_of_reporter"`
	TitleEn                           string     `json:"title_en"`
	TitleAr                           string     `json:"title_ar"`
	ReportedEntityID                  uint       `json:"reported_entity_id"`
	ReportedEntityTitleEn             string     `json:"reported_entity_title_en"`
	ReportedEntityTitleAr             string     `json:"reported_entity_title_ar"`
	DescriptionEn                     *string    `json:"description_en"`
	DescriptionAr                     *string    `json:"description_ar"`
	ViolationReportCategoryID         *uint      `json:"violation_report_category_id"`
	ViolationReportCategoryTitleEn    *string    `json:"violation_report_category_title_en"`
	ViolationReportCategoryTitleAr    *string    `json:"violation_report_category_title_ar"`
	ViolationReportSubCategoryID      *uint      `json:"violation_report_sub_category_id"`
	ViolationReportSubCategoryTitleEn *string    `json:"violation_report_sub_category_title_en"`
	ViolationReportSubCategoryTitleAr *string    `json:"violation_report_sub_category_title_ar"`
	ViolationReportReasonID           uint       `json:"violation_report_reason_id"`
	ViolationReportReasonTitleEn      string     `json:"violation_report_reason_title_en"`
	ViolationReportReasonTitleAr      string     `json:"violation_report_reason_title_ar"`
	SeverityLevelID                   *uint      `json:"severity_level_id"`
	SeverityLevelTitleEn              *string    `json:"severity_level_title_en"`
	SeverityLevelTitleAr              *string    `json:"severity_level_title_ar"`
	SeverityLevelColor                *string    `json:"severity_level_color_code"`
	ViolationReportStatusID           uint       `json:"violation_report_status_id"`
	ViolationReportStatusTitleEn      string     `json:"violation_report_status_title_en"`
	ViolationReportStatusTitleAr      string     `json:"violation_report_status_title_ar"`
	ViolationReportStatusColorCode    string     `json:"violation_report_status_color_code"`
	CreatedAt                         time.Time  `json:"created_at"`
	UpdatedAt                         *time.Time `json:"updated_at"`
	CreatedByID                       uint       `json:"created_by_id"`
	CreatedUserFirstNameEn            string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr            string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn             string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr             string     `json:"created_user_last_name_ar"`
	CreatedUserInitial                string     `json:"created_user_initial"`
	CreatedUserColor                  string     `json:"created_user_color"`
	CreatedUserEmail                  string     `json:"created_user_email"`
	CreatedUserProfileImage           string     `json:"created_user_profile_image_path"`
	UpdatedByID                       *uint      `json:"updated_by_id"`
	UpdatedUserFirstNameEn            string     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr            string     `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn             string     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr             string     `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                string     `json:"updated_user_initial"`
	UpdatedUserColor                  string     `json:"updated_user_color"`
	UpdatedUserEmail                  string     `json:"updated_user_email"`
	UpdatedUserProfileImage           string     `json:"updated_user_profile_image_path"`
}

type ViolationReportFormDTO struct {
	ViolationReportChannelID     uint                              `json:"violation_report_channel_id" binding:"required"`
	IDOfReporter                 *uint                             `json:"id_of_reporter"`
	NameOfReporter               *string                           `json:"name_of_reporter"`
	PhoneOfReporter              *string                           `json:"phone_of_reporter"`
	TitleEn                      string                            `json:"title_en" binding:"required,max=255"`
	TitleAr                      string                            `json:"title_ar" binding:"required,max=255"`
	ReportedEntityID             uint                              `json:"reported_entity_id" binding:"required"`
	DescriptionEn                *string                           `json:"description_en"`
	DescriptionAr                *string                           `json:"description_ar"`
	EvidenceDescription          *string                           `json:"evidence_description"`
	ViolationReportCategoryID    *uint                             `json:"violation_report_category_id"`
	ViolationReportSubCategoryID *uint                             `json:"violation_report_sub_category_id"`
	ViolationReportReasonID      uint                              `json:"violation_report_reason_id" binding:"required"`
	SeverityLevelID              *uint                             `json:"severity_level_id"`
	ReportedAt                   string                            `json:"reported_at"`
	Documents                    []*ViolationReportDocumentFormDTO `json:"documents"`
}

type ViolationReportDocumentFormDTO struct {
	ID            *uint  `json:"id,omitempty"`
	URL           string `json:"url"`
	FileName      string `json:"file_name"`
	OriginalName  string `json:"original_name"`
	Type          string `json:"type"`
	Size          int    `json:"size"`
	IsNewFile     *bool  `json:"is_new_file,omitempty"`
	IsDeletedFile *bool  `json:"is_deleted_file,omitempty"`
}

type ViolationReportDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         string  `json:"name"`
	Token        string  `json:"token"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}

type ViolationReportListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	ReportedAt  string
	Export      bool
}

type ViolationReportListResponseDTO struct {
	Items      []*ViolationReportListDTO    `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ViolationReportDetailDTO struct {
	ID                         uint                                  `json:"id"`
	ReferenceNumber            string                                `json:"reference_number"`
	TitleEn                    string                                `json:"title_en"`
	TitleAr                    string                                `json:"title_ar"`
	DescriptionEn              string                                `json:"description_en"`
	DescriptionAr              string                                `json:"description_ar"`
	ReportedAt                 time.Time                             `json:"reported_at"`
	ViolationReportChannel     ViolationReportChannelDTO             `json:"violation_report_channel"`
	ViolationReportCategory    *ViolationReportCategoryDetailsDTO    `json:"violation_report_category"`
	ViolationReportSubCategory *ViolationReportSubCategoryDetailsDTO `json:"violation_report_sub_category"`
	ViolationReportReason      ViolationReportReasonDTO              `json:"violation_report_reason"`
	SeverityLevel              *masterDTOs.SeverityLevelDTO          `json:"severity_level"`
	Documents                  []*ViolationReportDocumentDTO         `json:"documents"`
	ReportedEntity             settingsDTOs.OrganizationUnitDTO      `json:"reported_entity"`
	EvidenceDescription        *string                               `json:"evidence_description"`
	CreatedAt                  time.Time                             `json:"created_at"`
	UpdatedAt                  *time.Time                            `json:"updated_at"`
	CreatedByUser              organizationDTOs.UserMinimalDTO       `json:"created_by"`
	UpdatedByUser              organizationDTOs.UserMinimalDTO       `json:"updated_by"`
}

type ViolationReportParsedDates struct {
	ReportedAt time.Time
}
