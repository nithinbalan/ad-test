package models

import (
	"time"
)

type RiskRegister struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber          string     `json:"reference_number" binding:"required,max=45" gorm:"type:varchar(45);not null;unique"`
	RiskCategoryID           uint       `json:"risk_category_id" binding:"required"`
	RiskSourceID             uint       `json:"risk_source_id" binding:"required"`
	OwnerUnitID              uint       `json:"owner_unit_id" binding:"required"`
	OwnerID                  uint       `json:"owner_id" binding:"required"`
	RiskStatusID             uint       `json:"risk_status_id" binding:"required"`
	Title                    string     `json:"title" binding:"required,max=255"`
	Description              *string    `json:"description" gorm:"type:nvarchar(max)"`
	ReportedAt               *time.Time `json:"reported_at"`
	CurrentLikelihoodID      *uint      `json:"current_likelihood_id" gorm:"default:null"`
	CurrentImpactID          *uint      `json:"current_impact_id" gorm:"default:null"`
	CurrentRiskRatingID      *uint      `json:"current_risk_rating_id" gorm:"default:null"`
	CurrentRiskRatingValueID *uint      `json:"current_risk_rating_value_id" gorm:"default:null"`
	ReviewFrequencyID        *uint      `json:"review_frequency_id"`
	RiskResponseStrategyID   *uint      `json:"risk_response_strategy_id"`
	NextReviewDate           *time.Time `json:"next_review_date"`
	Comments                 *string    `json:"comments" gorm:"type:nvarchar(max)"`
	CreatedByID              uint       `json:"created_by_id" binding:"required"`
	UpdatedByID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	RiskCategory              *RiskCategory                `gorm:"foreignKey:RiskCategoryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskSource                *RiskSource                  `gorm:"foreignKey:RiskSourceID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	OwnerUnit                 *OrganizationUnit            `gorm:"foreignKey:OwnerUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Owner                     *User                        `gorm:"foreignKey:OwnerID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ReviewFrequency           *ReviewFrequency             `gorm:"foreignKey:ReviewFrequencyID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskResponseStrategy      *RiskResponseStrategy        `gorm:"foreignKey:RiskResponseStrategyID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskMatrixLikelihood      *RiskMatrixLikelihood        `gorm:"foreignKey:CurrentLikelihoodID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixImpact          *RiskMatrixImpact            `gorm:"foreignKey:CurrentImpactID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixRiskRating      *RiskMatrixRiskRating        `gorm:"foreignKey:CurrentRiskRatingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskMatrixRiskRatingValue *RiskMatrixRiskRatingValue   `gorm:"foreignKey:CurrentRiskRatingValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskStatus                *RiskStatus                  `gorm:"foreignKey:RiskStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser             *User                        `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser             *User                        `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Documents                 *[]RiskDocument              `gorm:"foreignKey:RiskRegisterID"`
	ComplianceRequirements    []*RiskComplianceRequirement `gorm:"foreignKey:RiskRegisterID"`
	Controls                  []*RiskControl               `gorm:"foreignKey:RiskRegisterID"`
	Processes                 []*RiskProcess               `gorm:"foreignKey:RiskRegisterID"`
	SLARegisters              []*RiskSLARegister           `gorm:"foreignKey:RiskRegisterID"`
	StrategicObjectives       []*RiskStrategicObjective    `gorm:"foreignKey:RiskRegisterID"`
}

// TableName overrides the table name used by RiskRegister to `risk_register`
func (RiskRegister) TableName() string {
	return "risk_register"
}
