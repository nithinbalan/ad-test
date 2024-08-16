package models

import (
	"time"
)

type AuditableItem struct {
	ID                      uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	TitleEn                 *string    `json:"title_en" gorm:"type:varchar(500);default:null"`
	TitleAr                 *string    `json:"title_ar" gorm:"type:nvarchar(500);default:null"`
	ProcessID               *uint      `json:"process_id" gorm:"default:null"`
	ControlID               *uint      `json:"control_id" gorm:"default:null"`
	DocumentID              *uint      `json:"document_id" gorm:"default:null"`
	IncidentRegisterID      *uint      `json:"incident_register_id" gorm:"default:null"`
	ResponsibleUnitID       uint       `json:"responsible_unit_id" binding:"required" gorm:"not null"`
	ComplianceRequirementID *uint      `json:"compliance_requirement_id" gorm:"default:null"`
	RiskRegisterID          *uint      `json:"risk_register_id" gorm:"default:null"`
	UserID                  *uint      `json:"user_id" gorm:"default:null"`
	Comment                 string     `json:"comment" binding:"required" gorm:"type:nvarchar(MAX);not null"`
	AuditorRiskRatingID     uint       `json:"auditor_risk_rating_id" gorm:"not null"`
	AuditStatusID           uint       `json:"audit_status_id" gorm:"not null"`
	CreatedAt               time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt               *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID             uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID             *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	Process               *Process               `gorm:"foreignKey:ProcessID"`
	Control               *Control               `gorm:"foreignKey:ControlID"`
	Document              *Document              `gorm:"foreignKey:DocumentID"`
	IncidentRegister      *IncidentRegister      `gorm:"foreignKey:IncidentRegisterID"`
	ComplianceRequirement *ComplianceRequirement `gorm:"foreignKey:ComplianceRequirementID"`
	RiskRegister          *RiskRegister          `gorm:"foreignKey:RiskRegisterID"`
	User                  *User                  `gorm:"foreignKey:UserID"`
	AuditorRiskRating     *AuditorRiskRating     `gorm:"foreignKey:AuditorRiskRatingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AuditStatus           *AuditStatus           `gorm:"foreignKey:AuditStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit       *OrganizationUnit      `gorm:"foreignKey:ResponsibleUnitID"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser         *User                  `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by AuditableItem to `auditable_items`
func (AuditableItem) TableName() string {
	return "auditable_items"
}
