package models

import (
	"time"
)

type BcpProcessAllocationOutageBucket struct {
	ID                     uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BcpProcessAllocationID uint      `json:"bcp_process_allocation_id" gorm:"not null"`
	OutageBucketID         *uint     `json:"outage_bucket_id" gorm:"default:null"`
	IsSkipped              bool      `json:"is_skipped" gorm:"type:bit;not null;default:0"`
	CreatedBYID            *uint     `json:"created_by_id" gorm:"default:null"`
	CreatedAt              time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	BcpProcessStaffAllocation *BcpProcessStaffAllocation             `gorm:"foreignKey:BcpProcessAllocationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	OutageBucket              *BCMBIAOutageBucket                    `gorm:"foreignKey:OutageBucketID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser             *User                                  `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Users                     []BcpProcessAllocationOutageBucketUser `gorm:"foreignKey:BcpProcessAllocationOutageBucketID;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BcpProcessAllocationOutageBucket) TableName() string {
	return "bcp_process_allocation_outage_buckets"
}
