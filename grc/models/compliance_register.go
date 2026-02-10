package models

import (
	"time"
)

type ComplianceRegister struct {
	ID                              uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceCategoryID            *uint      `json:"compliance_category_id" binding:"required" gorm:"not null"`
	ComplianceStatusID              *uint      `json:"compliance_status_id" gorm:"default:null"`
	ResponsibleUnitID               *uint      `json:"responsible_unit_id" binding:"required" gorm:"not null"`
	ResponsibleUserID               *uint      `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ReviewFrequencyID               *uint      `json:"review_frequency_id" binding:"required" gorm:"default:null"`
	ReviewResponsibleUserID         *uint      `json:"review_responsible_user_id" binding:"omitempty" gorm:"default:null"`
	SeverityLevelID                 *uint      `json:"severity_level_id" binding:"omitempty" gorm:"default:null"`
	ReferenceNumber                 string     `json:"reference_number" binding:"omitempty,max=45" gorm:"size:45"`
	StandardCode                    *string    `json:"standard_code" binding:"omitempty,max=45" gorm:"size:45"`
	TitleEn                         string     `json:"title_en" binding:"required,max=255" gorm:"type:varchar(255);not null"`
	TitleAr                         string     `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	DescriptionEn                   *string    `json:"description_en" gorm:"type:varchar(max)"`
	DescriptionAr                   *string    `json:"description_ar" gorm:"type:nvarchar(max)"`
	Version                         *string    `json:"version" binding:"omitempty,max=45" gorm:"size:45"`
	Deadline                        *time.Time `json:"deadline" binding:"omitempty" gorm:"type:date,default:null"`
	IssueDate                       *time.Time `json:"issue_date" binding:"omitempty" gorm:"type:date"`
	ExpiryDate                      *time.Time `json:"expiry_date" binding:"omitempty" gorm:"type:date"`
	NextReviewDate                  *time.Time `json:"next_review_date" binding:"omitempty" gorm:"type:date"`
	LastReviewedDate                *time.Time `json:"last_reviewed_date" gorm:"type:date"`
	LastReviewedUserID              *uint      `json:"last_reviewed_user_id" gorm:"default:null"`
	LastReviewNote                  *string    `json:"last_review_note" gorm:"type:varchar(max)"`
	IsLatest                        bool       `json:"is_latest"`
	ParentComplianceRegisterID      *uint      `json:"parent_compliance_register_id" gorm:"column:compliance_register_id"`
	CreatedAt                       time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                       *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID                     uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                     *uint      `json:"updated_by_id" gorm:"default:null"`
	LifecycleStatusID               *uint      `json:"lifecycle_status_id" gorm:"not null"`
	ComplianceReviewStatusID        *uint      `json:"compliance_review_status_id" gorm:"default:null"`
	TotalScore                      *int       `json:"total_score" gorm:"default:0"`
	Score                           *float64   `json:"score" gorm:"type:decimal(10,2)"`
	ComplianceFolderID              *uint      `json:"compliance_folder_id" binding:"omitempty" gorm:"default:null"`
	ComplianceAssessmentFrameworkID *uint      `json:"compliance_assessment_framework_id" binding:"required" gorm:"not null"`
	IsControlAssessment             bool       `json:"is_control_assessment" gorm:"column:is_control_assessment;not null;default:false"`
	IsRequirementAssessment         bool       `json:"is_requirement_assessment" gorm:"column:is_requirement_assessment;not null;default:false"`

	// New  Compliance Fields
	LegislativeBodyID                  *uint   `json:"legislative_body_id" gorm:"default:null"`
	LegislationTypeID                  *uint   `json:"legislation_type_id" gorm:"default:null"`
	LegislationNo                      *string `json:"legislation_no" gorm:"type:nvarchar(500)"`
	ImportanceClassificationID         *uint   `json:"importance_classification_id" gorm:"default:null"`
	SecrecyClassificationID            *uint   `json:"secrecy_classification_id" gorm:"default:null"`
	Subject                            *string `json:"subject" gorm:"type:nvarchar(max)"`
	IsApplicable                       bool    `json:"is_applicable default:false"`
	IsNotApplicable                    bool    `json:"is_not_applicable default:false"`
	IsRegulationNotificationAssessment bool    `json:"is_regulation_notification_assessment default:false"`
	IsRiskBasedAssessment              bool    `json:"is_risk_based_assessment default:false"`
	IsPublic                           bool    `json:"is_public default:false"`
	IsShared                           bool    `json:"is_shared default:false"`
	Note                               *string `json:"note" gorm:"type:varchar(max)"`

	// Relationships
	ComplianceCategory            *ComplianceCategory            `gorm:"foreignKey:ComplianceCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceStatus              *ComplianceStatus              `gorm:"foreignKey:ComplianceStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit               *OrganizationUnit              `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser               *User                          `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewFrequency               *ReviewFrequency               `gorm:"foreignKey:ReviewFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewResponsibleUser         *User                          `gorm:"foreignKey:ReviewResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel                 *SeverityLevel                 `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	LifecycleStatus               *LifecycleStatus               `gorm:"foreignKey:LifecycleStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	LastReviewedUser              *User                          `gorm:"foreignKey:LastReviewedUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceReviewStatus        *ComplianceReviewStatus        `gorm:"foreignKey:ComplianceReviewStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentEn                    *ComplianceDocument            `gorm:"foreignKey:ComplianceRegisterID"`
	DocumentAr                    *ComplianceDocument            `gorm:"foreignKey:ComplianceRegisterID"`
	ParentComplianceRegister      *ComplianceRegister            `gorm:"foreignKey:ParentComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ChildComplianceRegister       []*ComplianceRegister          `gorm:"foreignKey:ParentComplianceRegisterID"` // Recursive relationship
	CreatedByUser                 *User                          `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser                 *User                          `gorm:"foreignKey:UpdatedBYID"`
	Controls                      []*ComplianceControl           `gorm:"foreignKey:ComplianceRegisterID"`
	ComplianceStatusUpdates       []*ComplianceStatusUpdate      `gorm:"foreignKey:ComplianceRegisterID"`
	ComplianceClauses             []*ComplianceClause            `gorm:"foreignKey:ComplianceRegisterID"`
	ComplianceFolder              *ComplianceFolder              `gorm:"foreignKey:ComplianceFolderID"`
	ComplianceAssessmentFramework *ComplianceAssessmentFramework `gorm:"foreignKey:ComplianceAssessmentFrameworkID"`

	//New Relations
	LegislativeBody              *ComplianceLegislativeBody            `gorm:"foreignKey:LegislativeBodyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	LegislationType              *ComplianceLegislationType            `gorm:"foreignKey:LegislationTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ImportanceClassification     *ComplianceImportanceClassification   `gorm:"foreignKey:ImportanceClassificationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SecrecyClassification        *ComplianceSecrecyClassification      `gorm:"foreignKey:SecrecyClassificationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IssuingAuthorities           []*ComplianceRegisterIssuingAuthority `gorm:"foreignKey:ComplianceRegisterID"`
	ComplianceRegisterSharedUnit *[]ComplianceRegisterSharedUnit       `gorm:"foreignKey:ComplianceRegisterID"`
}

// TableName overrides the table name used by ComplianceRegister to `compliance_register`
func (ComplianceRegister) TableName() string {
	return "compliance_register"
}
