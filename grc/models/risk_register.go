package models

import (
	"time"
)

type RiskRegister struct {
	ID                        uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber           *string    `json:"reference_number" binding:"omitempty,max=45" gorm:"type:varchar(45);unique"`
	RiskCategoryID            *uint      `json:"risk_category_id"`
	RiskSourceID              *uint      `json:"risk_source_id" gorm:"default:null"`
	RiskMatrixID              uint       `json:"risk_matrix_id" binding:"required"`
	ResponsibleUnitID         *uint      `json:"responsible_unit_id"`
	ResponsibleUserID         *uint      `json:"responsible_user_id"`
	RiskStatusID              uint       `json:"risk_status_id" binding:"required"`
	Title                     string     `json:"title" binding:"required,max=255"`
	Description               *string    `json:"description" gorm:"type:nvarchar(max)"`
	ReportedAt                *time.Time `json:"reported_at"`
	CurrentLikelihoodID       *uint      `json:"current_likelihood_id" gorm:"default:null"`
	CurrentImpactID           *uint      `json:"current_impact_id" gorm:"default:null"`
	CurrentRiskRatingID       *uint      `json:"current_risk_rating_id" gorm:"default:null"`
	CurrentRiskRatingValueID  *uint      `json:"current_risk_rating_value_id" gorm:"default:null"`
	InherentLikelihoodID      *uint      `json:"inherent_likelihood_id" gorm:"default:null"`
	InherentImpactID          *uint      `json:"inherent_impact_id" gorm:"default:null"`
	InherentRiskRatingID      *uint      `json:"inherent_risk_rating_id" gorm:"default:null"`
	InherentRiskRatingValueID *uint      `json:"inherent_risk_rating_value_id" gorm:"default:null"`
	ResidualLikelihoodID      *uint      `json:"residual_likelihood_id" gorm:"default:null"`
	ResidualImpactID          *uint      `json:"residual_impact_id" gorm:"default:null"`
	ResidualRiskRatingID      *uint      `json:"residual_risk_rating_id" gorm:"default:null"`
	ResidualRiskRatingValueID *uint      `json:"residual_risk_rating_value_id" gorm:"default:null"`
	ReviewFrequencyID         *uint      `json:"review_frequency_id"`
	LastReviewNote            *string    `json:"last_review_note,omitempty" gorm:"type:nvarchar(MAX);default:null" binding:"omitempty"`
	NextReviewDate            *time.Time `json:"next_review_date"`
	RiskResponseStrategyID    *uint      `json:"risk_response_strategy_id"`
	MarkedAsAuditableItem     bool       `json:"marked_as_auditable_item" gorm:"default:false" binding:"-"`
	SubmittedBy               *uint      `json:"submitted_by" gorm:"default:null"`
	SubmittedAt               *time.Time `json:"submitted_at" gorm:"default:null"`
	NextReviewLevel           *int       `json:"next_review_level" gorm:"default:null"`
	ReviewedBy                *uint      `json:"reviewed_by" gorm:"default:null"`
	SeverityLevelID           *uint      `json:"severity_level_id"`
	CreatedBYID               uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID               *uint      `json:"updated_by_id"`
	CreatedAt                 time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                 *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	RiskCode                 *string    `json:"risk_code" gorm:"type:nvarchar(500);default:null"`
	RegistrationDate         *time.Time `json:"registration_date" gorm:"type:datetime;default:null"`
	DivisionID               uint       `json:"division_id" binding:"required"`
	SectionID                *uint      `json:"section_id" gorm:"default:null"`
	IsStrategic              bool       `json:"is_strategic" gorm:"not null;default:false"`
	IsOperational            bool       `json:"is_operational" gorm:"not null;default:false"`
	ExpectedCause            *string    `json:"expected_cause" gorm:"default:null"`
	Consequence              *string    `json:"consequence" gorm:"default:null"`
	IsInternal               bool       `json:"is_internal" gorm:"not null;default:false"`
	IsExternal               bool       `json:"is_external" gorm:"not null;default:false"`
	RiskRegisterUpdatePlanID *uint      `json:"risk_register_update_plan_id" gorm:"default:null"`

	// Relationships
	RiskCategory            *RiskCategory              `gorm:"foreignKey:RiskCategoryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskSource              *RiskSource                `gorm:"foreignKey:RiskSourceID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ResponsibleUnit         *OrganizationUnit          `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ResponsibleUser         *User                      `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ReviewFrequency         *ReviewFrequency           `gorm:"foreignKey:ReviewFrequencyID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskResponseStrategy    *RiskResponseStrategy      `gorm:"foreignKey:RiskResponseStrategyID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SubmittedByUser         *User                      `gorm:"foreignKey:SubmittedBy;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ReviewedByUser          *User                      `gorm:"foreignKey:ReviewedBy;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CurrentLikelihood       *RiskMatrixLikelihood      `gorm:"foreignKey:CurrentLikelihoodID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CurrentImpact           *RiskMatrixImpact          `gorm:"foreignKey:CurrentImpactID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CurrentRiskRating       *RiskMatrixRiskRating      `gorm:"foreignKey:CurrentRiskRatingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CurrentRiskRatingValue  *RiskMatrixRiskRatingValue `gorm:"foreignKey:CurrentRiskRatingValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	InherentLikelihood      *RiskMatrixLikelihood      `gorm:"foreignKey:InherentLikelihoodID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	InherentImpact          *RiskMatrixImpact          `gorm:"foreignKey:InherentImpactID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	InherentRiskRating      *RiskMatrixRiskRating      `gorm:"foreignKey:InherentRiskRatingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	InherentRiskRatingValue *RiskMatrixRiskRatingValue `gorm:"foreignKey:InherentRiskRatingValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResidualLikelihood      *RiskMatrixLikelihood      `gorm:"foreignKey:ResidualLikelihoodID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResidualImpact          *RiskMatrixImpact          `gorm:"foreignKey:ResidualImpactID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResidualRiskRating      *RiskMatrixRiskRating      `gorm:"foreignKey:ResidualRiskRatingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResidualRiskRatingValue *RiskMatrixRiskRatingValue `gorm:"foreignKey:ResidualRiskRatingValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskStatus              *RiskStatus                `gorm:"foreignKey:RiskStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SeverityLevel           *SeverityLevel             `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser           *User                      `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser           *User                      `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Section                 *OrganizationUnit          `gorm:"foreignKey:SectionID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Division                *OrganizationUnit          `gorm:"foreignKey:DivisionID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskRegisterUpdatePlan  *RiskRegisterUpdatePlan    `gorm:"foreignKey:RiskRegisterUpdatePlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`

	Documents              *[]RiskDocument              `gorm:"foreignKey:RiskRegisterID"`
	ComplianceRequirements []*RiskComplianceRequirement `gorm:"foreignKey:RiskRegisterID"`
	Controls               []*RiskControl               `gorm:"foreignKey:RiskRegisterID"`
	Processes              []*RiskProcess               `gorm:"foreignKey:RiskRegisterID"`
	SLARegisters           []*RiskSLA                   `gorm:"foreignKey:RiskRegisterID"`
	StrategicObjectives    []*RiskStrategicObjective    `gorm:"foreignKey:RiskRegisterID"`
	KRIs                   []*RiskKri                   `gorm:"foreignKey:RiskRegisterID"`
	AssociatedDocuments    []*RiskAssociatedDocument    `gorm:"foreignKey:RiskRegisterID"`
	Assets                 []*RiskAsset                 `gorm:"foreignKey:RiskRegisterID"`
	Threats                []*RiskThreat                `gorm:"foreignKey:RiskRegisterID"`
	AuditableItems         []*AuditableItem             `gorm:"foreignKey:RiskRegisterID"`
	RiskWorkflows          []*RiskWorkflow              `gorm:"foreignKey:RiskRegisterID"`
}

// TableName overrides the table name used by RiskRegister to `risk_register`
func (RiskRegister) TableName() string {
	return "risk_register"
}
