package models

import (
	"time"
)

type BCMManualProcedureForm struct {
	ID            uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn       string     `json:"title_en" binding:"required,max=1000" gorm:"type:varchar(1000);not null"`
	TitleAr       string     `json:"title_ar" binding:"required,max=1000" gorm:"type:varchar(1000);not null"`
	Code          *string    `json:"code" gorm:"type:varchar(1000);default:null"`
	OwnerID       uint       `json:"owner_id" binding:"required" gorm:"not null"`
	LocationID    *uint      `json:"location_id" gorm:"default:null"`
	StakeholderID *uint      `json:"stakeholder_id" gorm:"default:null"`
	Notes         *string    `json:"notes" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID   uint       `json:"created_by_id"`
	UpdatedBYID   *uint      `json:"updated_by_id"`
	CreatedAt     time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	Owner         *User        `gorm:"foreignKey:OwnerID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Stakeholder   *Stakeholder `gorm:"foreignKey:StakeholderID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Location      *Location    `gorm:"foreignKey:LocationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser *User        `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BCMManualProcedureForm) TableName() string {
	return "bcm_manual_procedure_forms"
}
