package models

import "time"

type PolicyException struct {
	ID                         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	DocumentID                 uint       `json:"document_id" binding:"required" gorm:"not null"`
	DocumentVersionID          uint       `json:"document_version_id" binding:"required" gorm:"not null"`
	DocumentVersionClauseID    *uint      `json:"document_version_clause_id" gorm:"default:null"`
	DocumentRequirementGroupID *uint      `json:"document_requirement_group_id" gorm:"default:null"`
	ReferenceNumber            string     `json:"reference_number" binding:"required" gorm:"type:varchar(45);not null"`
	TitleEn                    string     `json:"title_en" binding:"required" gorm:"type:nvarchar(500);not null"`
	TitleAr                    string     `json:"title_ar" binding:"required" gorm:"type:nvarchar(500);not null"`
	Requirements               *string    `json:"requirements" gorm:"type:nvarchar(MAX);default:null"`
	Justification              *string    `json:"justification" gorm:"type:nvarchar(MAX);default:null"`
	StartDate                  time.Time  `json:"start_date" binding:"required" gorm:"type:date;not null"`
	EndDate                    *time.Time `json:"end_date" gorm:"type:date;default:null"`
	RequestedByID              uint       `json:"requested_by_id" binding:"required" gorm:"not null"`
	OrganizationUnitID         uint       `json:"organization_unit_id" binding:"required" gorm:"not null"`
	WorkflowStatusID           *uint      `json:"workflow_status_id" gorm:"default:null"`
	CreatedAt                  time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	CreatedByID                uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedAt                  *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	UpdatedByID                *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	Document                 *Document                 `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentVersion          *DocumentVersion          `gorm:"foreignKey:DocumentVersionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentVersionClause    *DocumentVersionClause    `gorm:"foreignKey:DocumentVersionClauseID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	DocumentRequirementGroup *DocumentRequirementGroup `gorm:"foreignKey:DocumentRequirementGroupID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	RequestedBy              *User                     `gorm:"foreignKey:RequestedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrganizationUnit         *OrganizationUnit         `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowStatus           *WorkflowStatus           `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CreatedByUser            *User                     `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser            *User                     `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the default table name
func (PolicyException) TableName() string {
	return "policy_exceptions"
}
