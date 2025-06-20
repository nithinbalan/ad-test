package models

import "time"

type PolicyExceptionRisk struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	PolicyExceptionID uint      `gorm:"not null" json:"policy_exception_id"`
	RiskRegisterID    uint      `gorm:"not null" json:"risk_register_id"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	CreatedByID       uint      `gorm:"not null" json:"created_by_id"`

	PolicyException *PolicyException `gorm:"foreignKey:PolicyExceptionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskRegister    *RiskRegister    `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User            `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the default table name
func (PolicyExceptionRisk) TableName() string {
	return "policy_exception_risks"
}
