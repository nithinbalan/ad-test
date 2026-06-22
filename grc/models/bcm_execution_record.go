package models

import (
	"time"
)

type BCMExecutionRecord struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn           string     `json:"title_en" binding:"required,max=1000" gorm:"type:varchar(1000);not null"`
	TitleAr           string     `json:"title_ar" binding:"required,max=1000" gorm:"type:varchar(1000);not null"`
	OwnerEntityID     uint       `json:"owner_entity_id" binding:"required" gorm:"not null"`
	RecordTypeID      *uint      `json:"record_type_id" gorm:"default:null"`
	IsBackUpAvailable *string    `json:"is_back_up_available" gorm:"type:varchar(20);default:null"`
	LocationID        *uint      `json:"location_id" gorm:"default:null"`
	RecoveryPlan      *string    `json:"recovery_plan" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID       uint       `json:"created_by_id"`
	UpdatedBYID       *uint      `json:"updated_by_id"`
	CreatedAt         time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	OwnerEntity   *OrganizationUnit       `gorm:"foreignKey:OwnerEntityID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RecordType    *BCMExecutionRecordType `gorm:"foreignKey:RecordTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Location      *Location               `gorm:"foreignKey:LocationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser *User                   `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BCMExecutionRecord) TableName() string {
	return "bcm_execution_records"
}
