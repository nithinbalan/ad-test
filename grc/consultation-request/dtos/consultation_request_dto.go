package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type ConsultationRequestListDTO struct {
	ID                              uint   `json:"id"`
	ReferenceNumber                 string `json:"reference_number"`
	DateRaised                      string `json:"date_raised"`
	Title                           string `json:"title"`
	RequestingDepartmentUnitID      uint   `json:"requesting_department_unit_id"`
	RequestingDepartmentUnitTitleEn string `json:"requesting_department_unit_title_en"`
	RequestingDepartmentUnitTitleAr string `json:"requesting_department_unit_title_ar"`

	ConsultedDepartmentUnitID      uint   `json:"consulted_department_unit_id"`
	ConsultedDepartmentUnitTitleEn string `json:"consulted_department_unit_title_en"`
	ConsultedDepartmentUnitTitleAr string `json:"consulted_department_unit_title_ar"`

	RequestDetails string `json:"request_details"`

	ResponsibleUserID      *uint  `json:"responsible_user_id"`
	ResponsibleFirstNameEn string `json:"responsible_user_first_name_en"`
	ResponsibleLastNameEn  string `json:"responsible_user_last_name_en"`
	ResponsibleEmail       string `json:"responsible_user_email"`
	ResponsibleFirstNameAr string `json:"responsible_user_first_name_ar"` // Arabic first name
	ResponsibleLastNameAr  string `json:"responsible_user_last_name_ar"`  // Arabic last name
	ResponsibleInitial     string `json:"responsible_user_initial"`

	ConsultationRequestStatusID uint   `json:"consultation_request_status_id"`
	StatusTitleEn               string `json:"status_title_en"`
	StatusTitleAr               string `json:"status_title_ar"`

	Reply     *string `json:"reply"`
	ReplyDate *string `json:"reply_date"`
	Notes     *string `json:"notes"`

	CreatedAt             time.Time `json:"created_at"`
	CreatedByID           uint      `json:"created_by_id"`
	CreatedByFirstNameEn  string    `json:"created_by_first_name_en"`
	CreatedByLastNameEn   string    `json:"created_by_last_name_en"`
	CreatedByFirstNameAr  string    `json:"created_by_first_name_ar"`
	CreatedByLastNameAr   string    `json:"created_by_last_name_ar"`
	CreatedByInitial      string    `json:"created_by_initial"` // initials
	CreatedByEmail        string    `json:"created_by_email"`
	CreatedByColor        string    `json:"created_by_color"`
	CreatedByProfileImage string    `json:"created_by_profile_image_path"`

	UpdatedByFirstNameEn  string `json:"updated_by_first_name_en"`
	UpdatedByLastNameEn   string `json:"updated_by_last_name_en"`
	UpdatedByFirstNameAr  string `json:"updated_by_first_name_ar"` // Arabic first name
	UpdatedByLastNameAr   string `json:"updated_by_last_name_ar"`  // Arabic last name
	UpdatedByInitial      string `json:"updated_by_initial"`       // initials
	UpdatedByEmail        string `json:"updated_by_email"`
	UpdatedByColor        string `json:"updated_by_color"`
	UpdatedByProfileImage string `json:"updated_by_profile_image_path"`
}
type ConsultationRequestListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Export   bool
}
type ConsultationRequestListResponseDTO struct {
	Items      []*ConsultationRequestListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO  `json:"pagination"`
}

// ConsultationRequestFormDTO is used for creating/updating consultation requests via API
type ConsultationRequestFormDTO struct {
	ConsultedDepartmentID uint    `json:"consulted_department_id" binding:"required"`
	RequestDetails        *string `json:"request_details" binding:"required"`
	Title                 string  `json:"title" binding:"required"`
	// ConsultationRequestStatusID uint    `json:"consultation_request_status_id"`
	// Reply                       *string `json:"reply,omitempty"`
	// Notes                       *string `json:"notes,omitempty"`
}
type ConsultationRequestCommentFormDTO struct {
	ConsultationRequestStatusID *uint  `json:"consultation_request_status_id"`
	Comment                     string `json:"comment" binding:"required"`
	Documents					[]*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type ConsultationRequestCommentUpdateFormDTO struct {
	Comment   string                        `json:"comment" binding:"required"`
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type ConsultationRequestDetailDTO struct {
	ID              uint   `json:"id"`
	ReferenceNumber string `json:"reference_number"`
	RaisedDate      string `json:"date_raised"`
	Title           string `json:"title"`

	RequestDetails *string `json:"request_details"`
	Reply          *string `json:"reply"`
	ReplyedDate    string  `json:"reply_date"`
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
	CreatedDate               string                        `json:"created_at"`
}
type ConsultationRequestStatusDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Type    string `json:"type"`
}

type ConsultationRequestCommentListDTO struct {
	ID                    uint   `json:"id"`                      // Comment ID
	ConsultationRequestID uint   `json:"consultation_request_id"` // Parent request ID
	Comment               string `json:"comment"`                 // The actual comment text

	// Status fields (flattened into the comment DTO)
	ConsultationRequestStatusID uint       `json:"consultation_request_status_id,omitempty"`
	StatusTitleEn               string     `json:"status_title_en,omitempty"`
	StatusTitleAr               string     `json:"status_title_ar,omitempty"`
	StatusDescriptionEn         string     `json:"status_description_en,omitempty"`
	StatusDescriptionAr         string     `json:"status_description_ar,omitempty"`
	StatusType                  string     `json:"status_type,omitempty"`
	Status                      string     `json:"status,omitempty"`
	StatusCreatedAt             *time.Time `json:"status_created_at,omitempty"`

	CreatedAt            time.Time `json:"created_at"`
	CreatedByID          uint      `json:"created_by_id"`
	CreatedByFirstNameEn string    `json:"created_by_first_name_en"`
	CreatedByLastNameEn  string    `json:"created_by_last_name_en"`
	CreatedByFirstNameAr string    `json:"created_by_first_name_ar"`
	CreatedByLastNameAr  string    `json:"created_by_last_name_ar"`

	CreatedByInitial string `json:"created_by_initial"`

	Documents []*ConsultationRequestCommentDocumentDTO `json:"documents" gorm:"-"`

}

type ConsultationRequestCommentDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}