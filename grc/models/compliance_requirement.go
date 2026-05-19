package models

import (
	"time"
)

type ComplianceRequirement struct {
	ID                              uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceClauseID              uint       `json:"compliance_clause_id" binding:"required" gorm:"not null"`
	TitleEn                         string     `json:"title_en" binding:"required,max=255" gorm:"type:varchar(255);not null"`
	TitleAr                         string     `json:"title_ar" binding:"required,max=255" gorm:"type:varchar(255);not null"`
	ComplianceRequirementCategoryID *uint      `json:"compliance_requirement_category_id" binding:"omitempty" gorm:"default:null"`
	SeverityLevelID                 *uint      `json:"severity_level_id" binding:"omitempty" gorm:"default:null"`
	Deadline                        *time.Time `json:"deadline" binding:"omitempty" gorm:"type:date"`
	ComplianceStatusID              *uint      `json:"compliance_status_id"`
	MarkedAsAuditableItem           bool       `json:"marked_as_auditable_item" gorm:"default:false" binding:"-"`
	CreatedAt                       time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                       *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID                     uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID                     *uint      `json:"updated_by_id" gorm:"default:null"`
	Score                           *float64   `json:"score" gorm:"type:decimal(10,2)"`
	InherentRiskRatingID            *uint      `json:"inherent_risk_rating_id" gorm:"default:null"`
	InherentRiskScore               *float64   `json:"inherent_risk_score" gorm:"type:decimal(10,2)"`
	ResidualRiskRatingID            *uint      `json:"residual_risk_rating_id" gorm:"default:null"`
	ResidualRiskScore               *float64   `json:"residual_risk_score" gorm:"type:decimal(10,2)"`

	// Relationships
	ComplianceClause              *ComplianceClause                             `gorm:"foreignKey:ComplianceClauseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceStatus              *ComplianceStatus                             `gorm:"foreignKey:ComplianceStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRequirementCategory *ComplianceRequirementCategory                `gorm:"foreignKey:ComplianceRequirementCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel                 *SeverityLevel                                `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                 *User                                         `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser                 *User                                         `gorm:"foreignKey:UpdatedByID"`
	Controls                      []*ComplianceControl                          `gorm:"foreignKey:ComplianceRequirementID"`
	StrategicObjectives           []*ComplianceRequirementStrategicObjective    `gorm:"foreignKey:ComplianceRequirementID"`
	Documents                     []*ComplianceRequirementAssociatedDocument    `gorm:"foreignKey:ComplianceRequirementID"`
	RequiredDocumentTitles        []*ComplianceRequirementRequiredDocumentTitle `gorm:"foreignKey:ComplianceRequirementID"`
	ResponsibleUnits              []*ComplianceRequirementResponsibleUnit       `gorm:"foreignKey:ComplianceRequirementID"`
	InherentRiskAnalysis          []*ComplianceRequirementInherentRiskAnalysis  `gorm:"foreignKey:ComplianceRequirementID"`
	InherentRiskRating            *ComplianceAssessmentRating                   `gorm:"foreignKey:InherentRiskRatingID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ResidualRiskRating            *ComplianceAssessmentRating                   `gorm:"foreignKey:ResidualRiskRatingID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by ComplianceRequirement to `compliance_requirements`
func (ComplianceRequirement) TableName() string {
	return "compliance_requirements"
}
