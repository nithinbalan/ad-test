package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type KPIRegisterListDTO struct {
	ID                              int      `json:"id"`
	ReferenceNumber                 *string  `json:"reference_number"`
	Title                           string   `json:"title"`
	Description                     *string  `json:"description"`
	KPICategoryID                   uint     `json:"kpi_category_id"`
	KPICategoryTitleEn              string   `json:"kpi_category_title_en"`
	KPICategoryTitleAr              string   `json:"kpi_category_title_ar"`
	ResponsibleUserID               *uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      *string  `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn       *string  `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr      *string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr       *string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string  `json:"responsible_user_initial"`
	ResponsibleUserEmail            *string  `json:"responsible_user_email"`
	ResponsibleUserColor            *string  `json:"responsible_user_color"`
	ResponsibleUserProfileImagePath *string  `json:"responsible_user_profile_image_path"`
	KPIStatusID                     uint     `json:"kpi_status_id"`
	KPIStatusTitleEn                string   `json:"kpi_status_title_en"`
	KPIStatusTitleAr                string   `json:"kpi_status_title_ar"`
	KPIStatusColorCode              string   `json:"kpi_status_color_code"`
	KPIStatusType                   string   `json:"kpi_status_type"`
	ReviewFrequencyTitleEn          *string  `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr          *string  `json:"review_frequency_title_ar"`
	ReviewFrequencyType             *string  `json:"review_frequency_type"`
	CalculationTypeTitleEn          *string  `json:"calculation_type_title_en"`
	CalculationTypeTitleAr          *string  `json:"calculation_type_title_ar"`
	CalculationTypeType             *string  `json:"calculation_type_type"`
	TargetValue                     *float64 `json:"target_value"`
	UnitTitleEn                     *string  `json:"unit_title_en"`
	UnitTitleAr                     *string  `json:"unit_title_ar"`
	UnitCodeEn                      *string  `json:"unit_code_en"`
	UnitCodeAr                      *string  `json:"unit_code_ar"`
	UnitColorCode                   *string  `json:"unit_color_code"`
	StartDate                       *string  `json:"start_date"`
	EndDate                         *string  `json:"end_date"`
	Remarks                         *string  `json:"remarks"`
	Achievement                     *float64 `json:"achievement"`
	CreatedUserID                   uint     `json:"created_user_id"`
	CreatedUserFirstNameEn          string   `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string   `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string   `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string   `json:"created_user_last_name_ar"`
	CreatedUserInitial              string   `json:"created_user_initial"`
	CreatedUserColor                string   `json:"created_user_color"`
	CreatedUserEmail                string   `json:"created_user_email"`
	CreatedUserProfileImagePath     string   `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint    `json:"updated_user_id"`
	UpdatedUserFirstNameEn          *string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              *string  `json:"updated_user_initial"`
	UpdatedUserColor                *string  `json:"updated_user_color"`
	UpdatedUserEmail                *string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath     *string  `json:"updated_user_profile_image_path"`
}

type KPIRegisterFormDTO struct {
	Title             string  `json:"title" binding:"required"`
	Description       *string `json:"description"`
	KPICategoryID     uint    `json:"kpi_category_id" binding:"required"`
	ResponsibleUserID *uint   `json:"responsible_user_id" binding:"required"`

	ReviewFrequencyID *uint               `json:"review_frequency_id" binding:"required"`
	CalculationTypeID *uint               `json:"calculation_type_id" binding:"required"`
	DataInputs        []*InputDataFormDTO `json:"data_inputs"`
	Formula           *string             `json:"formula"`
	UnitID            *uint               `json:"unit_id" binding:"required"`

	TargetValue *float64 `json:"target_value" binding:"required"`

	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`

	Remarks *string `json:"remarks" gorm:"type:nvarchar(max)"`
}

type InputDataFormDTO struct {
	Variable string `json:"variable"`
	Input    string `json:"input"`
}

type KPIRegisterListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type KPIRegisterListResponseDTO struct {
	Items      []*KPIRegisterListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type KPIRegisterDetailDTO struct {
	ID                uint                             `json:"id"`
	ReferenceNumber   *string                          `json:"reference_number"`
	Title             string                           `json:"title"`
	Description       *string                          `json:"description"`
	KPICategoryID     uint                             `json:"kpi_category_id"`
	KPICategory       *KPICategoryDTO                  `json:"kpi_category"`
	ResponsibleUserID *uint                            `json:"responsible_user_id"`
	ResponsibleUser   *organizationDTOs.UserMinimalDTO `json:"responsible_user"`
	KPIStatusID       uint                             `json:"kpi_status_id"`
	KPIStatus         KPIStatusDTO                     `json:"kpi_status"`
	ReviewFrequencyID *uint                            `json:"review_frequency_id"`
	ReviewFrequency   *masterDTOs.ReviewFrequencyDTO   `json:"review_frequency"`
	CalculationTypeID *uint                            `json:"calculation_type_id"`
	CalculationType   *masterDTOs.CalculationTypeDTO   `json:"calculation_type"`
	Formula           *string                          `json:"formula"`
	UnitID            *uint                            `json:"unit_id"`
	Unit              *masterDTOs.UnitDTO              `json:"unit"`
	ActualValue       *float64                         `json:"actual_value"`
	TargetValue       *float64                         `json:"target_value"`
	Achievement       *float64                         `json:"achievement"`
	StartDateString   *string                          `json:"start_date"`
	EndDateString     *string                          `json:"end_date"`
	Remarks           *string                          `json:"remarks"`
	DataInputs        []*KPIRegisterDataInputDTO       `json:"data_inputs"`
	CreatedAt         time.Time                        `json:"created_at"`
	UpdatedAt         *time.Time                       `json:"updated_at"`
	CreatedByUser     organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser     *organizationDTOs.UserMinimalDTO `json:"updated_by"`
}

type KPIRegisterParsedDates struct {
	StartDate *time.Time
	EndDate   *time.Time
}

type KPIRegisterDataInputDTO struct {
	ID       uint   `json:"id"`
	Variable string `json:"variable"`
	Input    string `json:"input"`
}

type KPITrendAnalysisDTO struct {
	ReviewedAtString string  `json:"reviewed_at"`
	Value            float64 `json:"value"`
}
