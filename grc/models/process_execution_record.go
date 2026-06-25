package models

import (
	"time"
)

type ProcessExecutionRecord struct {
	ID                   uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	ProcessID            uint      `json:"process_id" gorm:"not null"`
	BCMExecutionRecordID uint      `json:"bcm_execution_records_id" gorm:"not null;column:bcm_execution_records_id"`
	CreatedBYID          uint      `json:"created_by_id"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Process            *Process            `gorm:"foreignKey:ProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	BCMExecutionRecord *BCMExecutionRecord `gorm:"foreignKey:BCMExecutionRecordID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProcessExecutionRecord) TableName() string {
	return "process_execution_records"
}
