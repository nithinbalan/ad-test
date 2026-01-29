package models

import "time"

type ConsultationRequest struct {
	ID                          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title                       string     `json:"title" gorm:"type:nvarchar(255);not null"`
	ReferenceNumber             string     `json:"reference_number" gorm:"type:varchar(45);not null;uniqueIndex"` // CONS-12323
	DateRaised                  time.Time  `json:"date_raised" gorm:"type:datetime;not null"`
	RequestingDepartmentUnitID  uint       `json:"requesting_department_unit_id" gorm:"not null"`
	ConsultedDepartmentUnitID   uint       `json:"consulted_department_unit_id" gorm:"not null"`
	ResponsibleUserID           *uint      `json:"responsible_user_id" gorm:"default:null"`
	RequestDetails              *string    `json:"request_details" gorm:"type:nvarchar(MAX);not null"`
	ConsultationRequestStatusID *uint      `json:"consultation_request_status_id" gorm:"not null"`
	Reply                       *string    `json:"reply" gorm:"type:nvarchar(MAX);default:null"`
	ReplyDate                   *time.Time `json:"reply_date" gorm:"type:date;default:null"`
	Notes                       *string    `json:"notes" gorm:"type:nvarchar(MAX);default:null"`

	// Audit Fields
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	CreatedBYID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	UpdatedByID *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	ResponsibleUser           *User                      `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	RequestingDepartmentUnit  *OrganizationUnit          `gorm:"foreignKey:RequestingDepartmentUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ConsultedDepartmentUnit   *OrganizationUnit          `gorm:"foreignKey:ConsultedDepartmentUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ConsultationRequestStatus *ConsultationRequestStatus `gorm:"foreignKey:ConsultationRequestStatusID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CreatedByUser             *User                      `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser             *User                      `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Comments []*ConsultationRequestComment `gorm:"foreignKey:ConsultationRequestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName overrides the default table name
func (ConsultationRequest) TableName() string {
	return "consultation_requests"
}
