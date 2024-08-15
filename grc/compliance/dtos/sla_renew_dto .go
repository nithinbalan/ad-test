package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type SLARenewListDTO struct {
	ID                        uint      `json:"id"`
	SLARegisterID             uint      `json:"sla_register_id"`
	Version                   string    `json:"version"`
	IssueDate                 string    `json:"issue_date"`
	ExpiryDate                string    `json:"expiry_date"`
	SA1Date                   string    `json:"sa1_date"`
	SA2Date                   string    `json:"sa2_date"`
	Comment                   string    `json:"comment"`
	URL                       string    `json:"url"`
	Name                      string    `json:"name"`
	Size                      int       `json:"size"`
	Token                     string    `json:"token"`
	Type                      string    `json:"type"`
	OriginalName              string    `json:"original_name"`
	IsLatest                  bool      `json:"is_latest"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
	CreatedByFirstNameEn      string    `json:"created_by_first_name_en"`
	CreatedByFirstNameAr      string    `json:"created_by_first_name_ar"`
	CreatedByLastNameEn       string    `json:"created_by_last_name_en"`
	CreatedByLastNameAr       string    `json:"created_by_last_name_ar"`
	CreatedByInitial          string    `json:"created_by_initial"`
	CreatedByColor            string    `json:"created_by_color"`
	CreatedByEmail            string    `json:"created_by_email"`
	CreatedByProfileImagePath string    `json:"created_by_profile_image_path"`
	UpdatedByFirstNameEn      string    `json:"updated_by_first_name_en"`
	UpdatedByFirstNameAr      string    `json:"updated_by_first_name_ar"`
	UpdatedByLastNameEn       string    `json:"updated_by_last_name_en"`
	UpdatedByLastNameAr       string    `json:"updated_by_last_name_ar"`
	UpdatedByInitial          string    `json:"updated_by_initial"`
	UpdatedByColor            string    `json:"updated_by_color"`
	UpdatedByEmail            string    `json:"updated_by_email"`
	UpdatedByProfileImagePath string    `json:"updated_by_profile_image_path"`
}

type SLARenewListQueryParamsDTO struct {
	Sort           string
	SortBy         string
	Keyword        string
	Status         string
	Page           int
	PageSize       int
	IssueDateFrom  string
	IssueDateTo    string
	ExpiryDateFrom string
	ExpiryDateTo   string
}

type SLARenewListResponseDTO struct {
	Items      []*SLARenewListDTO           `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type SLARenewFormDTO struct {
	IssueDate  string                  `json:"issue_date" binding:"omitempty,datetime=2006-01-02"`
	ExpiryDate string                  `json:"expiry_date" binding:"omitempty,datetime=2006-01-02"`
	SA1Date    string                  `json:"sa1_date" binding:"omitempty,datetime=2006-01-02"`
	SA2Date    string                  `json:"sa2_date" binding:"omitempty,datetime=2006-01-02"`
	Comment    *string                 `json:"comment"`
	Document   SLARenewDocumentFormDTO `json:"document" binding:"required"`
}

type SLARenewDocumentFormDTO struct {
	ID           *uint  `json:"id,omitempty"`
	URL          string `json:"url"`
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
	Size         int    `json:"size"`
}

type ParsedSLARenewDates struct {
	IssueDate  *time.Time
	ExpiryDate *time.Time
	SA1Date    *time.Time
	SA2Date    *time.Time
}

//renew
