package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

// ConsultationRequestFormDTO is used for creating/updating consultation requests via API
type ConsultationRequestFormDTO struct {
	ConsultedDepartmentID       uint    `json:"consulted_department_id" binding:"required"`
	RequestDetails              *string `json:"request_details" binding:"required"`
	ConsultationRequestStatusID *uint   `json:"consultation_request_status_id"` // optional, default = New
	Reply                       *string `json:"reply,omitempty"`
	Notes                       *string `json:"notes,omitempty"`
}
type ConsultationRequestListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
}
type ConsultationRequestListResponseDTO struct {
	Items      []*ConsultationRequestListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO  `json:"pagination"`
}
type ConsultationRequestListDTO struct {
	ID              uint   `json:"id"`
	ReferenceNumber string `json:"reference_number"`
	DateRaised      string `json:"date_raised"`

	RequestingDepartmentUnitID      uint   `json:"requesting_department_unit_id"`
	RequestingDepartmentUnitTitleEn string `json:"requesting_department_unit_title_en"`
	RequestingDepartmentUnitTitleAr string `json:"requesting_department_unit_title_ar"`

	ConsultedDepartmentUnitID      uint   `json:"consulted_department_unit_id"`
	ConsultedDepartmentUnitTitleEn string `json:"consulted_department_unit_title_en"`
	ConsultedDepartmentUnitTitleAr string `json:"consulted_department_unit_title_ar"`

	RequestDetails string `json:"request_details"`

	ResponsibleUserID      *uint  `json:"responsible_user_id"`
	ResponsibleFirstNameEn string `json:"responsible_first_name_en"`
	ResponsibleLastNameEn  string `json:"responsible_last_name_en"`
	ResponsibleEmail       string `json:"responsible_email"`

	ConsultationRequestStatusID uint   `json:"consultation_request_status_id"`
	StatusTitleEn               string `json:"status_title_en"`
	StatusTitleAr               string `json:"status_title_ar"`
	StatusColor                 string `json:"status_color"`

	Reply     *string `json:"reply"`
	ReplyDate *string `json:"reply_date"`
	Notes     *string `json:"notes"`

	CreatedAt            time.Time `json:"created_at"`
	CreatedByID          uint      `json:"created_by_id"`
	CreatedByFirstNameEn string    `json:"created_by_first_name_en"`
	CreatedByLastNameEn  string    `json:"created_by_last_name_en"`

	UpdatedByID          *uint   `json:"updated_by_id"`
	UpdatedByFirstNameEn string  `json:"updated_by_first_name_en"`
	UpdatedByLastNameEn  string  `json:"updated_by_last_name_en"`
	UpdatedAt            *string `json:"updated_at"`
}

type ConsultationRequestDetailDTO struct {
	ID              uint    `json:"id"`
	ReferenceNumber string  `json:"reference_number"`
	DateRaised      string  `json:"date_raised"`
	DueDate         *string `json:"due_date"`

	RequestDetails *string `json:"request_details"`
	Reply          *string `json:"reply"`
	ReplyDate      *string `json:"reply_date"`
	Notes          *string `json:"notes"`

	// Department Units
	RequestingDepartmentUnit settingsDTOs.OrganizationUnitDTO `json:"requesting_department_unit"`
	ConsultedDepartmentUnit  settingsDTOs.OrganizationUnitDTO `json:"consulted_department_unit"`

	// Users
	ResponsibleUser *organizationDTOs.UserMinimalDTO `json:"responsible_user"`
	CreatedByUser   organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser   *organizationDTOs.UserMinimalDTO `json:"updated_by"`

	// Status
	ConsultationRequestStatus *ConsultationRequestStatusDTO `json:"consultation_request_status"`

	// Audit
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
type ConsultationRequestStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
}
