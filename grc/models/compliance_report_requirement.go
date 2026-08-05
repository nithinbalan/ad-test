package models

import (
	"time"
)

type ComplianceReportRequirement struct {
	ID                                 uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceReportID                 uint      `json:"compliance_report_id" gorm:"not null" binding:"required"`
	ComplianceRegisterID               *uint     `json:"compliance_register_id" gorm:"default:null"`
	ComplianceReportClauseID           uint      `json:"compliance_report_clause_id" gorm:"not null" binding:"required"`
	TitleEn                            string    `json:"title_en" gorm:"type:nvarchar(500);not null" binding:"required"`
	TitleAr                            string    `json:"title_ar" gorm:"type:nvarchar(500);not null" binding:"required"`
	Score                              *int      `json:"score" gorm:"default:0"`
	AverageScore                       *float64  `json:"average_score" gorm:"default:0" comment:"Calculated average score across all assessments"`
	InherentScore                      *float64  `json:"inherent_score" gorm:"type:decimal(10,2);default:null"`
	InherentRiskRatingID               *uint     `json:"inherent_risk_rating_id" gorm:"default:null"`
	ResidualScore                      *float64  `json:"residual_score" gorm:"type:decimal(10,2);default:null"`
	ResidualRiskRatingID               *uint     `json:"residual_risk_rating_id" gorm:"default:null"`
	ComplianceStatusID                 *uint     `json:"compliance_status_id" gorm:"default:null"`
	AverageComplianceStatusID          *uint     `json:"average_compliance_status_id" gorm:"default:null" comment:"Compliance status based on average score"`
	IsAcknowledged                     bool      `json:"is_acknowledged" gorm:"type:bit;not null;default:0"`
	AcknowledgeUserID                  *uint     `json:"acknowledge_user_id" gorm:"default:null"`
	AcknowledgeComment                 *string   `json:"acknowledge_comment" gorm:"type:nvarchar(MAX);default:null"`
	IsRiskBasedAssessment              *bool     `json:"is_risk_based_assessment" gorm:"default:null"`
	IsRegulationNotificationAssessment *bool     `json:"is_regulation_notification_assessment" gorm:"default:null"`
	CreatedAt                          time.Time `json:"created_at" gorm:"not null"`
	CreatedBYID                        uint      `json:"created_by_id" gorm:"not null" binding:"required"`
	UpdatedAt                          time.Time `json:"updated_at" gorm:"default:null"`
	UpdatedBYID                        uint      `json:"updated_by_id" gorm:"default:null"`

	// Associations
	ComplianceReport            ComplianceReport             `json:"compliance_report" gorm:"foreignKey:ComplianceReportID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ComplianceRegister          *ComplianceRegister          `json:"compliance_register" gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ComplianceReportClause      ComplianceReportClause       `json:"compliance_report_clause" gorm:"foreignKey:ComplianceReportClauseID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ComplianceStatus            *ComplianceStatus            `json:"compliance_status,omitempty" gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;foreignKey:ComplianceStatusID"`
	AverageComplianceStatus     *ComplianceStatus            `json:"average_compliance_status,omitempty" gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;foreignKey:AverageComplianceStatusID"`
	InherentRiskRating          *RiskMatrixRiskRating        `json:"inherent_risk_rating,omitempty" gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;foreignKey:InherentRiskRatingID"`
	ResidualRiskRating          *RiskMatrixRiskRating        `json:"residual_risk_rating,omitempty" gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;foreignKey:ResidualRiskRatingID"`
	AcknowledgeUser             *User                        `json:"acknowledge_user,omitempty" gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;foreignKey:AcknowledgeUserID"`
	ComplianceReportActionPlans []ComplianceReportActionPlan `json:"compliance_report_action_plans,omitempty" gorm:"foreignKey:ComplianceReportRequirementID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser               User                         `json:"created_by_user" gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser               User                         `json:"updated_by_user" gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (ComplianceReportRequirement) TableName() string {
	return "compliance_report_requirements"
}
