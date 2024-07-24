package models

import (
	"time"
)

type IncidentRegister struct {
	ID                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber       string     `json:"reference_number" binding:"required,max=45" gorm:"type:varchar(45);not null;unique"`
	IncidentCategoryID    uint       `json:"incident_category_id" binding:"required"`
	IncidentSubCategoryID uint       `json:"incident_sub_category_id" binding:"required"`
	IncidentTypeID        uint       `json:"incident_type_id" binding:"required"`
	SeverityLevelID       uint       `json:"severity_level_id" binding:"required"`
	IncidentStatusID      *uint      `json:"incident_status_id"`
	Title                 string     `json:"title" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	Description           *string    `json:"description" gorm:"type:nvarchar(max)"`
	IncidentAt            *time.Time `json:"incident_at" binding:"required" gorm:"type:date;not null"`
	Location              *string    `json:"location" gorm:"type:varchar(100)"`
	ReportedBy            uint       `json:"reported_by" binding:"required"`
	ReportedAt            *time.Time `json:"reported_at" gorm:"type:date"`
	QuickAction           *string    `json:"quick_action" gorm:"type:nvarchar(max)"`
	MarkedAsAuditableItem bool       `json:"marked_as_auditable_item" gorm:"default:false" binding:"-"`
	CreatedByID           uint       `json:"created_by_id" binding:"required"`
	UpdatedByID           *uint      `json:"updated_by_id"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	IncidentCategory       *IncidentCategory                `gorm:"foreignKey:IncidentCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentSubCategory    *IncidentSubCategory             `gorm:"foreignKey:IncidentSubCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentType           *IncidentType                    `gorm:"foreignKey:IncidentTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel          *SeverityLevel                   `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentStatus         *IncidentStatus                  `gorm:"foreignKey:IncidentStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Documents              []*IncidentDocument              `gorm:"foreignKey:IncidentRegisterID"`
	Controls               []*IncidentControl               `gorm:"foreignKey:IncidentRegisterID"`
	ComplianceRequirements []*IncidentComplianceRequirement `gorm:"foreignKey:IncidentRegisterID"`
	Risks                  []*IncidentRisk                  `gorm:"foreignKey:IncidentRegisterID"`
	Processes              []*IncidentProcess               `gorm:"foreignKey:IncidentRegisterID"`
	SLARegisters           []*IncidentSLA                   `gorm:"foreignKey:IncidentRegisterID"`
	ReportedByUser         *User                            `gorm:"foreignKey:ReportedBy"`
	CreatedByUser          *User                            `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser          *User                            `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by IncidentRegister to `incident_register`
func (IncidentRegister) TableName() string {
	return "incident_register"
}
