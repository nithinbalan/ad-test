package models

import (
	"time"
)

type ProcessManualProcedureForm struct {
	ID                       uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	ProcessID                uint      `json:"process_id" gorm:"not null"`
	BCMManualProcedureFormID uint      `json:"bcm_manual_procedure_form_id" gorm:"not null;column:bcm_manual_procedure_form_id"`
	CreatedBYID              uint      `json:"created_by_id"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Process                *Process               `gorm:"foreignKey:ProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	BCMManualProcedureForm *BCMManualProcedureForm `gorm:"foreignKey:BCMManualProcedureFormID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                  `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProcessManualProcedureForm) TableName() string {
	return "process_manual_procedure_forms"
}
