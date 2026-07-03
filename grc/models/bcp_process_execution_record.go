package models

import (
	"time"
)

type BcpProcessExecutionRecord struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BcpProcessID         uint      `json:"bcp_process_id" gorm:"not null"`
	BCMExecutionRecordID uint      `json:"bcm_execution_records_id" gorm:"not null;column:bcm_execution_records_id"`
	CreatedBYID          *uint     `json:"created_by_id"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	BcpProcess         *BcpProcess         `gorm:"foreignKey:BcpProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	BCMExecutionRecord *BCMExecutionRecord  `gorm:"foreignKey:BCMExecutionRecordID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BcpProcessExecutionRecord) TableName() string {
	return "bcp_process_execution_records"
}
