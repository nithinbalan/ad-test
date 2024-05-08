package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type SLARegisterListDTO struct {
	ID                         uint      `json:"id"`
	ReferenceNumber            string    `json:"reference_number"`
	ResponsibleUserID          uint      `json:"responsible_user_id"`
	LifecycleStatusID          uint      `json:"lifecycle_status_id"`
	TitleEn                    string    `json:"title_en"`
	TitleAr                    string    `json:"title_ar"`
	Description                string    `json:"description"`
	Client                     string    `json:"client"`
	CreatedBy                  uint      `json:"created_by"`
	UpdatedBy                  uint      `json:"updated_by"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
	SLACategoryTitleEn         string    `json:"sla_category_title_en"`
	SLACategoryTitleAr         string    `json:"sla_category_title_ar"`
	SeverityLevelTitleEn       string    `json:"severity_level_title_en"`
	SeverityLevelTitleAr       string    `json:"severity_level_title_ar"`
	ResponsibleUnitTitleEn     string    `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr     string    `json:"responsible_unit_title_ar"`
	ResponsibleUserFirstNameEn string    `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr string    `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn  string    `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr  string    `json:"responsible_user_last_name_ar"`
	LifecycleStatusTitleEn     string    `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr     string    `json:"lifecycle_status_title_ar"`
	ComplianceStatusTitleEn    string    `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr    string    `json:"compliance_status_title_ar"`
	CreatedByFirstNameEn       string    `json:"created_by_first_name_en"`
	CreatedByFirstNameAr       string    `json:"created_by_first_name_ar"`
	CreatedByLastNameEn        string    `json:"created_by_last_name_en"`
	CreatedByLastNameAr        string    `json:"created_by_last_name_ar"`
	UpdatedByFirstNameEn       string    `json:"updated_by_first_name_en"`
	UpdatedByFirstNameAr       string    `json:"updated_by_first_name_ar"`
	UpdatedByLastNameEn        string    `json:"updated_by_last_name_en"`
	UpdatedByLastNameAr        string    `json:"updated_by_last_name_ar"`
	// Fields from sla_documents table
	Version    string `json:"version"`
	IssueDate  string `json:"issue_date"`
	ExpiryDate string `json:"expiry_date"`
	SA1Date    string `json:"sa1_date"`
	SA2Date    string `json:"sa2_date"`
	Comment    string `json:"comment"`
	URL        string `json:"url"`
	Name       string `json:"name"`
	Size       int    `json:"size"`
	Token      string `json:"token"`
	Type       string `json:"type"`
	Icon       string `json:"icon"`
	IsLatest   bool   `json:"is_latest"`
}

type SLARegisterListQueryParamsDTO struct {
	Sort           string
	SortBy         string
	Keyword        string
	Status         string
	Page           int
	PageSize       int
	Filters        map[string][]string
	QuickFilter    []string
	IssueDateFrom  string
	IssueDateTo    string
	ExpiryDateFrom string
	ExpiryDateTo   string
}

type SLARegisterListResponseDTO struct {
	Items      []*SLARegisterListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type SLARegisterFormDTO struct {
	SLACategoryID     uint               `json:"sla_category_id" binding:"required"`
	SeverityLevelID   uint               `json:"severity_level_id"`
	ResponsibleUnitID uint               `json:"responsible_unit_id"`
	ResponsibleUserID uint               `json:"responsible_user_id"`
	LifecycleStatusID uint               `json:"lifecycle_status_id"`
	TitleEn           string             `json:"title_en" binding:"required,max=250"`
	TitleAr           string             `json:"title_ar" binding:"required,max=250"`
	Description       *string            `json:"description"`
	Client            string             `json:"client" binding:"max=250"`
	Version           string             `json:"version" binding:"required"`
	IssueDate         string             `json:"issue_date" binding:"omitempty,datetime=2006-01-02"`
	ExpiryDate        string             `json:"expiry_date" binding:"omitempty,datetime=2006-01-02"`
	SA1Date           string             `json:"sa1_date" binding:"omitempty,datetime=2006-01-02"`
	SA2Date           string             `json:"sa2_date" binding:"omitempty,datetime=2006-01-02"`
	Comment           *string            `json:"comment"`
	Document          SLADocumentFormDTO `json:"document" binding:"required"`
}

type SLADocumentFormDTO struct {
	ID           *uint  `json:"id,omitempty"`
	URL          string `json:"url"`
	Icon         string `json:"icon" binding:"required"`
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
	Size         int    `json:"size"`
	IsNewFile    *bool  `json:"is_new_file,omitempty"` // Optional field for indicating if the file is new or updated
}
type SLADocumentDTO struct {
	ID         *uint      `json:"id"`
	Version    string     `json:"version"`
	IssueDate  *time.Time `json:"issue_date,omitempty"`
	ExpiryDate *time.Time `json:"expiry_date,omitempty"`
	Comment    *string    `json:"comment,omitempty"`
	URL        string     `json:"url"`
	Name       string     `json:"name"`
	Size       int        `json:"size"`
	Token      string     `json:"token"`
	Type       string     `json:"type"`
	Icon       string     `json:"icon"`
	IsLatest   bool       `json:"is_latest"`
	SA1Date    *time.Time `json:"sa1_date,omitempty"`
	SA2Date    *time.Time `json:"sa2_date,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	CreatedBy  uint       `json:"created_by"`
	UpdatedBy  uint       `json:"updated_by,omitempty"`
}

type ParsedSLADates struct {
	IssueDate  *time.Time
	ExpiryDate *time.Time
	SA1Date    *time.Time
	SA2Date    *time.Time
}

type SLARegisterDetailDTO struct {
	ID               uint                                 `json:"id"`
	ReferenceNumber  string                               `json:"reference_number"`
	SLACategory      SLACategoryDTO                       `json:"sla_category"`
	SeverityLevel    *masterDTOs.SeverityLevelDTO         `json:"severity_level"`
	ResponsibleUnit  organizationDTOs.OrganizationUnitDTO `json:"responsible_unit"`
	ResponsibleUser  organizationDTOs.UserMinimalDTO      `json:"responsible_user"`
	LifecycleStatus  masterDTOs.LifecycleStatusDTO        `json:"lifecycle_status"`
	ComplianceStatus *ComplianceStatusDTO                 `json:"compliance_status,omitempty"`
	TitleEn          string                               `json:"title_en"`
	TitleAr          string                               `json:"title_ar"`
	Description      *string                              `json:"description,omitempty"`
	Client           string                               `json:"client"`
	CreatedAt        time.Time                            `json:"created_at"`
	UpdatedAt        time.Time                            `json:"updated_at"`
	CreatedByUser    organizationDTOs.UserMinimalDTO      `json:"created_by"`
	UpdatedByUser    *organizationDTOs.UserMinimalDTO     `json:"updated_by,omitempty"`
	Documents        []*SLADocumentDTO                    `json:"documents"`
}
