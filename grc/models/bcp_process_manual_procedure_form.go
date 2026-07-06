package models

import (
	"time"
)

type BcpProcessManualProcedureForm struct {
	ID                       uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BcpProcessID             uint      `json:"bcp_process_id" gorm:"not null"`
	BCMManualProcedureFormID uint      `json:"bcm_manual_procedure_form_id" gorm:"not null;column:bcm_manual_procedure_form_id"`
	CreatedBYID              *uint     `json:"created_by_id"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	BcpProcess             *BcpProcess             `gorm:"foreignKey:BcpProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	BCMManualProcedureForm *BCMManualProcedureForm  `gorm:"foreignKey:BCMManualProcedureFormID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BcpProcessManualProcedureForm) TableName() string {
	return "bcp_process_manual_procedure_forms"
}
