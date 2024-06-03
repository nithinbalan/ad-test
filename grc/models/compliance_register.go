package models

import (
	"time"
)

type ComplianceRegister struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceCategoryID     uint       `json:"compliance_category_id" binding:"required" gorm:"not null"`
	ComplianceStatusID       uint       `json:"compliance_status_id" binding:"required" gorm:"not null"`
	ResponsibleUnitID        uint       `json:"responsible_unit_id" binding:"required" gorm:"not null"`
	ResponsibleUserID        uint       `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ReviewFrequencyID        *uint      `json:"review_frequency_id" binding:"required" gorm:"default:null"`
	ReviewResponsibleUserID  *uint      `json:"review_responsible_user_id" binding:"omitempty" gorm:"default:null"`
	SeverityLevelID          *uint      `json:"severity_level_id" binding:"omitempty" gorm:"default:null"`
	ReferenceNumber          string     `json:"reference_number" binding:"omitempty,max=45" gorm:"size:45"`
	StandardCode             *string    `json:"standard_code" binding:"omitempty,max=45" gorm:"size:45"`
	TitleEn                  string     `json:"title_en" binding:"required,max=255" gorm:"type:varchar(255);not null"`
	TitleAr                  string     `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	DescriptionEn            *string    `json:"description_en" gorm:"type:varchar(max)"`
	DescriptionAr            *string    `json:"description_ar" gorm:"type:nvarchar(max)"`
	Version                  *string    `json:"version" binding:"omitempty,max=45" gorm:"size:45"`
	Deadline                 *time.Time `json:"deadline" binding:"omitempty" gorm:"type:date,default:null"`
	IssueDate                *time.Time `json:"issue_date" binding:"omitempty" gorm:"type:date"`
	ExpiryDate               *time.Time `json:"expiry_date" binding:"omitempty" gorm:"type:date"`
	NextReviewDate           *time.Time `json:"next_review_date" binding:"omitempty" gorm:"type:date"`
	LastReviewedDate         *time.Time `json:"last_reviewed_date" gorm:"type:date"`
	LastReviewedUserID       *uint      `json:"last_reviewed_user_id" gorm:"default:null"`
	LastReviewNote           *string    `json:"last_review_note" gorm:"type:varchar(max)"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBy                uint       `json:"created_by" gorm:"not null"`
	UpdatedBy                *uint      `json:"updated_by" gorm:"default:null"`
	LifecycleStatusID        uint       `json:"lifecycle_status_id" gorm:"not null"`
	ComplianceReviewStatusID *uint      `json:"compliance_review_status_id" gorm:"default:null"`

	// Relationships
	ComplianceCategory     *ComplianceCategory     `gorm:"foreignkey:ComplianceCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceStatus       *ComplianceStatus       `gorm:"foreignkey:ComplianceStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit        *OrganizationUnit       `gorm:"foreignkey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser        *User                   `gorm:"foreignkey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewFrequency        *ReviewFrequency        `gorm:"foreignkey:ReviewFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewResponsibleUser  *User                   `gorm:"foreignkey:ReviewResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel          *SeverityLevel          `gorm:"foreignkey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	LifecycleStatus        *LifecycleStatus        `gorm:"foreignkey:LifecycleStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	LastReviewedUser       *User                   `gorm:"foreignkey:LastReviewedUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceReviewStatus *ComplianceReviewStatus `gorm:"foreignkey:ComplianceReviewStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Documents              []ComplianceDocument    `gorm:"foreignKey:ComplianceRegisterID"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedBy"`
	Controls               []*ComplianceControl    `gorm:"foreignKey:ComplianceRegisterID"`
}

// TableName overrides the table name used by ComplianceRegister to `compliance_register`
func (ComplianceRegister) TableName() string {
	return "compliance_register"
}
