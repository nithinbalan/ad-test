package models

import (
	"time"
)

type RiskUpdatePlanIntroductoryMeetingDocument struct {
	ID                       uint       `json:"id" gorm:"primary_key;auto_increment"`
	RiskRegisterUpdatePlanID uint       `json:"risk_register_update_plan_id" gorm:"not null"`
	URL                      *string    `json:"url" gorm:"type:varchar(1000);default:null"`
	Name                     *string    `json:"name" gorm:"type:nvarchar(255);default:null"`
	Size                     *int       `json:"size" gorm:"default:null"`
	Token                    *string    `json:"token" gorm:"type:varchar(255);default:null"`
	OriginalName             *string    `json:"original_name" gorm:"type:nvarchar(255);default:null"`
	Type                     *string    `json:"type" gorm:"type:varchar(255);default:null"`
	CreatedBYID              *uint      `json:"created_by_id" gorm:"default:null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	RiskRegisterUpdatePlan *RiskRegisterUpdatePlan `gorm:"foreignKey:RiskRegisterUpdatePlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (RiskUpdatePlanIntroductoryMeetingDocument) TableName() string {
	return "risk_update_plan_introductory_meeting_documents"
}
