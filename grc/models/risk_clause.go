package models

import (
	"time"
)

type RiskClause struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID       uint      `json:"risk_register_id" binding:"required"`
	ComplianceRegisterID uint      `json:"compliance_register_id" binding:"required"`
	ComplianceClauseID   uint      `json:"compliance_clause_id" binding:"required"`
	CreatedByID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Risk               *RiskRegister       `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRegister *ComplianceRegister `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceClause   *ComplianceClause   `gorm:"foreignKey:ComplianceClauseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedByID"`
}

func (RiskClause) TableName() string {
	return "risk_clauses"
}
