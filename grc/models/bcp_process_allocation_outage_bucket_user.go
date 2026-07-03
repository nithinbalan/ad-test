package models

import (
	"time"
)

type BcpProcessAllocationOutageBucketUser struct {
	ID                               uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BcpProcessAllocationOutageBucketID uint    `json:"bcp_process_allocation_outage_bucket_id" gorm:"not null"`
	UserID                           *uint     `json:"user_id" gorm:"default:null"`
	IsRemote                         bool      `json:"is_remote" gorm:"type:bit;not null;default:0"`
	CreatedBYID                      *uint     `json:"created_by_id" gorm:"default:null"`
	CreatedAt                        time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	BcpProcessAllocationOutageBucket *BcpProcessAllocationOutageBucket `gorm:"foreignKey:BcpProcessAllocationOutageBucketID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	User                             *User                             `gorm:"foreignKey:UserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser                    *User                             `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BcpProcessAllocationOutageBucketUser) TableName() string {
	return "bcp_process_allocation_outage_bucket_users"
}
