package models

import (
	"time"
)

type BCMStaffAllocationOutageBucketUser struct {
	ID                               uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	BCMStaffAllocationOutageBucketID uint      `json:"bcm_staff_allocation_outage_bucket_id" gorm:"not null"`
	UserID                           uint      `json:"user_id" gorm:"not null"`
	IsRemote                         bool      `json:"is_remote" gorm:"type:bit;default:false"`
	CreatedBYID                      uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                        time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	BCMStaffAllocationOutageBucket *BCMStaffAllocationOutageBucket `gorm:"foreignKey:BCMStaffAllocationOutageBucketID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	User                           *User                           `gorm:"foreignKey:UserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser                  *User                           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BCMStaffAllocationOutageBucketUser) TableName() string {
	return "bcm_staff_allocation_outage_bucket_users"
}
