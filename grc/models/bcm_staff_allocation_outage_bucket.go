package models

import (
	"time"
)

type BCMStaffAllocationOutageBucket struct {
	ID                   uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	BCMStaffAllocationID uint       `json:"bcm_staff_allocation_id" gorm:"not null"`
	OutageBucketID       *uint      `json:"outage_bucket_id" gorm:"default:null"`
	IsSkipped            bool       `json:"is_skipped" gorm:"type:bit;default:false"`
	UpdatedBYID          *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedBYID          uint       `json:"created_by_id"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	BCMStaffAllocation *BCMStaffAllocation                  `gorm:"foreignKey:BCMStaffAllocationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	OutageBucket       *BCMBIAOutageBucket                  `gorm:"foreignKey:OutageBucketID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser      *User                                `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser      *User                                `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Users              []BCMStaffAllocationOutageBucketUser `gorm:"foreignKey:BCMStaffAllocationOutageBucketID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BCMStaffAllocationOutageBucket) TableName() string {
	return "bcm_staff_allocation_outage_buckets"
}
