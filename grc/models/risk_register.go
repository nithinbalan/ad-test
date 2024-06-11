package models

import (
	"time"
)

type RiskRegister struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber   string     `json:"reference_number" binding:"required,max=45" gorm:"type:varchar(45);not null;unique"`
	RiskCategoryID    uint       `json:"risk_category_id" binding:"required"`
	RiskSourceID      uint       `json:"risk_source_id" binding:"required"`
	OwnerUnitID       uint       `json:"owner_unit_id" binding:"required"`
	OwnerID           uint       `json:"owner_id" binding:"required"`
	RiskStatusID      uint       `json:"risk_status_id" binding:"required"`
	Title             string     `json:"title" binding:"required,max=255"`
	Description       *string    `json:"description" gorm:"type:nvarchar(max)"`
	ReportedAt        *time.Time `json:"reported_at"`
	ReviewFrequencyID *uint      `json:"review_frequency_id"`
	NextReviewDate    *time.Time `json:"next_review_date"`
	Comments          *string    `json:"comments" gorm:"type:nvarchar(max)"`
	CreatedBy         uint       `json:"created_by" binding:"required"`
	UpdatedBy         *uint      `json:"updated_by"`
	CreatedAt         time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	RiskCategory           *RiskCategory                `gorm:"foreignkey:RiskCategoryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskSource             *RiskSource                  `gorm:"foreignkey:RiskSourceID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	OwnerUnit              *OrganizationUnit            `gorm:"foreignkey:OwnerUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Owner                  *User                        `gorm:"foreignkey:OwnerID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ReviewFrequency        *ReviewFrequency             `gorm:"foreignkey:ReviewFrequencyID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskStatus             *RiskStatus                  `gorm:"foreignkey:RiskStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                        `gorm:"foreignkey:CreatedBy;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser          *User                        `gorm:"foreignkey:UpdatedBy;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Documents              *[]RiskDocument              `gorm:"foreignKey:RiskRegisterID"`
	ComplianceRequirements []*RiskComplianceRequirement `gorm:"foreignKey:RiskRegisterID"`
	Controls               []*RiskControl               `gorm:"foreignKey:RiskRegisterID"`
	Processes              []*RiskProcess               `gorm:"foreignKey:RiskRegisterID"`
}

// TableName overrides the table name used by RiskRegister to `risk_register`
func (RiskRegister) TableName() string {
	return "risk_register"
}
