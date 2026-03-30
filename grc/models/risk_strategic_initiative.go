package models

import "time"

type RiskStrategicInitiative struct {
	ID                   uint      `gorm:"column:id;primaryKey;autoIncrement"`
	StrategicIntiativeID uint      `gorm:"column:strategic_intiative_id;not null"`
	RiskRegisterID       uint      `gorm:"column:risk_register_id;not null"`
	CreatedByID          uint      `gorm:"column:created_by_id;not null"`
	CreatedAt            time.Time `gorm:"column:created_at;not null"`

	// Relationships
	StrategicIntiative *StrategicIntiative `gorm:"foreignKey:StrategicIntiativeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskRegister       *RiskRegister       `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedByID"`
}

func (RiskStrategicInitiative) TableName() string {
	return "risks_strategic_intiatives"
}
