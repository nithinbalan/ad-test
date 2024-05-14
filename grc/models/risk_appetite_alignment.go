package models

import (
	"time"
)

type RiskAppetiteAlignment struct {
	ID                           uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID               uint       `json:"risk_register_id" binding:"required" gorm:"not null"`
	RiskAppetiteStatementID      uint       `json:"risk_appetite_statement_id" binding:"required" gorm:"not null"`
	RiskAppetiteToleranceLevelID uint       `json:"risk_appetite_tolerance_level_id" binding:"required" gorm:"not null"`
	CreatedAt                    *time.Time `json:"created_at" gorm:"type:datetime;default:null"`
	UpdatedAt                    *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBy                    uint       `json:"created_by" gorm:"not null"`
	UpdatedBy                    uint       `json:"updated_by" gorm:"not null"`

	// Relationships
	RiskRegister               *RiskRegister               `gorm:"foreignkey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskAppetiteStatement      *RiskAppetiteStatement      `gorm:"foreignkey:RiskAppetiteStatementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskAppetiteToleranceLevel *RiskAppetiteToleranceLevel `gorm:"foreignkey:RiskAppetiteToleranceLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser              *User                       `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser              *User                       `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by RiskAppetiteAlignment to `risk_appetite_alignments`
func (RiskAppetiteAlignment) TableName() string {
	return "risk_appetite_alignments"
}
