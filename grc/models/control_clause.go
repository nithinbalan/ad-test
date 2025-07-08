package models

import (
	"time"
)

type ControlClause struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ControlID            uint      `json:"control_id" binding:"required"`
	ComplianceRegisterID uint      `json:"compliance_register_id" binding:"required"`
	ComplianceClauseID   uint      `json:"compliance_clause_id" binding:"required"`
	CreatedByID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceRegister *ComplianceRegister `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control            *Control            `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceClause   *ComplianceClause   `gorm:"foreignKey:ComplianceClauseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedByID"`
}

func (ControlClause) TableName() string {
	return "control_clauses"
}
