package dtos

import (
	commonDTOs "jethur/grc/common/dtos"

	"time"
)

type KRIValueUpdateListDTO struct {
	ID                     uint    `json:"id"`
	KRIRegisterID          uint    `json:"kri_register_id"`
	ReviewedDate           string  `json:"reviewed_date"`
	Value                  string  `json:"value"`
	LowThreshold           string  `json:"low_threshold"`
	MediumThreshold        string  `json:"medium_threshold"`
	HighThreshold          string  `json:"high_threshold"`
	Remarks                string  `json:"remarks"`
	CreatedAt              string  `json:"created_at"`
	UpdatedAt              string  `json:"updated_at"`
	CreatedUserFirstNameEn string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  *string `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn  *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr  *string `json:"updated_user_last_name_ar"`
}

type KRIValueUpdateFormDTO struct {
	ReviewedDate string                           `json:"reviewed_date" binding:"required" gorm:"type:date;not null"`
	Value        string                           `json:"value" binding:"required" gorm:"type:varchar(455);not null"`
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
	ID              uint                  `json:"id"`
	KRIRegisterID   uint                  `json:"kri_register_id"`
	KRIRegister     KRIRegisterDetailsDTO `json:"kri_register"`
	ReviewedDate    string                `json:"reviewed_date"`
	Value           string                `json:"value"`
	LowThreshold    string                `json:"low_threshold"`
	MediumThreshold string                `json:"medium_threshold"`
	HighThreshold   string                `json:"high_threshold"`
	Remarks         string                `json:"remarks"`
	CreatedAt       string                `json:"created_at"`
	UpdatedAt       string                `json:"updated_at"`
}

type KRIValueUpdateParsedDates struct {
	ReviewedDate *time.Time
}

type KRIValueUpdateDocumentFormDTO struct {
	ID            *uint  `json:"id,omitempty"`
	URL           string `json:"url"`
	Icon          string `json:"icon"`
	FileName      string `json:"file_name"`
	OriginalName  string `json:"original_name"`
	Type          string `json:"type"`
	Size          int    `json:"size"`
	IsNewFile     *bool  `json:"is_new_file,omitempty"` // Optional field for indicating if the file is new or updated
	IsDeletedFile *bool  `json:"is_deleted_file,omitempty"`
}
