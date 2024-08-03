package dtos

import (
	commonDTOs "jethur/grc/common/dtos"

	"time"
)

type KRIValueUpdateListDTO struct {
	ID                          uint                         `json:"id"`
	KRIRegisterID               uint                         `json:"kri_register_id"`
	ReviewedDate                string                       `json:"reviewed_date"`
	Value                       int                          `json:"value"`
	LowThreshold                int                          `json:"low_threshold"`
	MediumThreshold             int                          `json:"medium_threshold"`
	HighThreshold               int                          `json:"high_threshold"`
	Remarks                     string                       `json:"remarks"`
	CreatedAt                   string                       `json:"created_at"`
	UpdatedAt                   string                       `json:"updated_at"`
	CreatedUserFirstNameEn      string                       `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                       `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                       `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       *string                      `json:"created_user_last_name_ar"`
	CreatedUserInitial          *string                      `json:"created_user_initial"`
	CreatedUserColor            *string                      `json:"created_user_color"`
	CreatedUserEmail            *string                      `json:"created_user_email"`
	CreatedUserProfileImagePath *string                      `json:"created_user_profile_image_path"`
	UpdatedUserFirstNameEn      *string                      `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string                      `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string                      `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string                      `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string                      `json:"updated_user_initial"`
	UpdatedUserColor            *string                      `json:"updated_user_color"`
	UpdatedUserEmail            *string                      `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string                      `json:"updated_user_profile_image_path"`
	Documents                   []*KRIValueUpdateDocumentDTO `json:"documents" gorm:"-"`
}

type KRIValueUpdateFormDTO struct {
	ReviewedDate string                           `json:"reviewed_date" binding:"required" gorm:"type:date;not null"`
	Value        int                              `json:"value" binding:"required" gorm:"type:int;not null"`
	Remarks      *string                          `json:"remarks" gorm:"type:nvarchar(max)"`
	Documents    []*KRIValueUpdateDocumentFormDTO `json:"documents"`
}

type KRIValueUpdateListQueryParamsDTO struct {
	Sort             string
	SortBy           string
	Keyword          string
	Status           string
	Page             int
	PageSize         int
	ReviewedDateFrom string
	ReviewedDateTo   string
}

type KRIValueUpdateListResponseDTO struct {
	Items      []*KRIValueUpdateListDTO     `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type KRIValueUpdateDetailsDTO struct {
	ID            uint `json:"id"`
	KRIRegisterID uint `json:"kri_register_id"`
	// KRIRegister     KRIRegisterDetailsDTO        `json:"kri_register"`
	// ReviewedDate    time.Time                    `json:"reviewed_date"`
	ReviewedDateString *string                      `json:"reviewed_date"`
	Value              int                          `json:"value"`
	LowThreshold       int                          `json:"low_threshold"`
	MediumThreshold    int                          `json:"medium_threshold"`
	HighThreshold      int                          `json:"high_threshold"`
	Remarks            string                       `json:"remarks"`
	CreatedAt          time.Time                    `json:"created_at"`
	UpdatedAt          *time.Time                   `json:"updated_at"`
	Documents          []*KRIValueUpdateDocumentDTO `json:"documents"`
}

type KRIValueUpdateParsedDates struct {
	ReviewedDate *time.Time
}

type KRIValueUpdateDocumentDTO struct {
	ID           *uint   `json:"id,omitempty"`
	URL          *string `json:"url"`
	Name         string  `json:"file_name"`
	OriginalName string  `json:"original_name"`
	Token        string  `json:"token"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}
type KRIValueUpdateDocumentFormDTO struct {
	ID            *uint  `json:"id,omitempty"`
	URL           string `json:"url"`
	FileName      string `json:"file_name"`
	OriginalName  string `json:"original_name"`
	Type          string `json:"type"`
	Size          int    `json:"size"`
	IsNewFile     *bool  `json:"is_new_file,omitempty"` // Optional field for indicating if the file is new or updated
	IsDeletedFile *bool  `json:"is_deleted_file,omitempty"`
}
