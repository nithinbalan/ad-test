package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type KPIValueUpdateListDTO struct {
	ID                          uint                         `json:"id"`
	KPIRegisterID               uint                         `json:"kpi_register_id"`
	ReviewedDate                string                       `json:"reviewed_date"`
	Value                       int                          `json:"value"`
	Remarks                     string                       `json:"remarks"`
	ReviewUserID                *uint                        `json:"review_user_id"`
	ReviewUserFirstNameEn       *string                      `json:"review_user_first_name_en"`
	ReviewUserLastNameEn        *string                      `json:"review_user_last_name_en"`
	ReviewUserFirstNameAr       *string                      `json:"review_user_first_name_ar"`
	ReviewUserLastNameAr        *string                      `json:"review_user_last_name_ar"`
	ReviewUserInitial           *string                      `json:"review_user_initial"`
	ReviewUserEmail             *string                      `json:"review_user_email"`
	ReviewUserColor             *string                      `json:"review_user_color"`
	ReviewUserProfileImagePath  *string                      `json:"review_user_profile_image_path"`
	CreatedAt                   string                       `json:"created_at"`
	UpdatedAt                   string                       `json:"updated_at"`
	CreatedUserID               uint                         `json:"created_user_id"`
	CreatedUserFirstNameEn      string                       `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                       `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                       `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       *string                      `json:"created_user_last_name_ar"`
	CreatedUserInitial          *string                      `json:"created_user_initial"`
	CreatedUserColor            *string                      `json:"created_user_color"`
	CreatedUserEmail            *string                      `json:"created_user_email"`
	CreatedUserProfileImagePath *string                      `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint                        `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string                      `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string                      `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string                      `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string                      `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string                      `json:"updated_user_initial"`
	UpdatedUserColor            *string                      `json:"updated_user_color"`
	UpdatedUserEmail            *string                      `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string                      `json:"updated_user_profile_image_path"`
	Documents                   []*KPIValueUpdateDocumentDTO `json:"documents" gorm:"-"`
}

type KPIValueUpdateFormDTO struct {
	ReviewedDate string                           `json:"reviewed_date" binding:"required" gorm:"type:date;not null"`
	ReviewUserID uint                             `json:"review_user_id"`
	Value        float64                          `json:"value" binding:"required" gorm:"type:int;not null"`
	Remarks      *string                          `json:"remarks" gorm:"type:nvarchar(max)"`
	Documents    []*KPIValueUpdateDocumentFormDTO `json:"documents"`
	DataInputs   []*ValueUpdateInputDataFormDTO   `json:"data_inputs"`
}

type ValueUpdateInputDataFormDTO struct {
	Variable string `json:"variable"`
	Input    string `json:"input"`
	Value    int    `json:"value"`
}

type KPIValueUpdateListQueryParamsDTO struct {
	Sort             string
	SortBy           string
	Keyword          string
	Status           string
	Page             int
	PageSize         int
	ReviewedDateFrom string
	ReviewedDateTo   string
}

type KPIValueUpdateListResponseDTO struct {
	Items      []*KPIValueUpdateListDTO     `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type KPIValueUpdateDetailsDTO struct {
	ID                 uint                             `json:"id"`
	KPIRegisterID      uint                             `json:"kpi_register_id"`
	ReviewedDateString *string                          `json:"reviewed_date"`
	Value              int                              `json:"value"`
	Remarks            string                           `json:"remarks"`
	Unit               *masterDTOs.UnitDTO              `json:"unit"`
	ReviewUser         *organizationDTOs.UserMinimalDTO `json:"review_user"`
	CreatedAt          time.Time                        `json:"created_at"`
	UpdatedAt          *time.Time                       `json:"updated_at"`
	Documents          []*KPIValueUpdateDocumentDTO     `json:"documents"`
	DataInputs         []*KPIValueUpdateDataInputDTO    `json:"data_inputs"`
}

type KPIValueUpdateDataInputDTO struct {
	ID       uint   `json:"id"`
	Variable string `json:"variable"`
	Input    string `json:"input"`
	Value    int    `json:"value"`
}

type KPIValueUpdateParsedDates struct {
	ReviewedDate *time.Time
}

type KPIValueUpdateDocumentDTO struct {
	ID           *uint   `json:"id,omitempty"`
	URL          *string `json:"url"`
	Name         string  `json:"file_name"`
	OriginalName string  `json:"original_name"`
	Token        string  `json:"token"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}
type KPIValueUpdateDocumentFormDTO struct {
	ID            *uint  `json:"id,omitempty"`
	URL           string `json:"url"`
	FileName      string `json:"file_name"`
	OriginalName  string `json:"original_name"`
	Type          string `json:"type"`
	Size          int    `json:"size"`
	IsNewFile     *bool  `json:"is_new_file,omitempty"` // Optional field for indicating if the file is new or updated
	IsDeletedFile *bool  `json:"is_deleted_file,omitempty"`
}
