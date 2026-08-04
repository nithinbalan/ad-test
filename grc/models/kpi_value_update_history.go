package models

import "time"

type KPIValueUpdateHistory struct {
	ID               uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	KPIValueUpdateID uint      `json:"kpi_value_update_id" gorm:"not null"`
	Value            *float64  `json:"value" gorm:"type:decimal(18,2);default:null"`
	UnitID           *uint     `json:"unit_id" gorm:"default:null"`
	Remarks          *string   `json:"remarks" gorm:"type:nvarchar(max);default:null"`
	CreatedAt        time.Time `json:"created_at" gorm:"type:datetime;not null"`
	CreatedBYID      *uint     `json:"created_by_id" gorm:"default:null"`
	CorrectiveAction *string   `json:"corrective_action"`

	KPIValueUpdate              *KPIValueUpdate               `gorm:"foreignKey:KPIValueUpdateID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Unit                        *Unit                         `gorm:"foreignKey:UnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser               *User                         `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	KPIValueUpdateHistoryInputs []*KPIValueUpdateHistoryInput `gorm:"foreignKey:KPIValueUpdateHistoryID"`
}

func (KPIValueUpdateHistory) TableName() string {
	return "kpi_value_update_histories"
}
