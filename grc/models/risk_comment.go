package models

import "time"

type RiskComment struct {
	ID                 uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID     uint      `json:"risk_register_id" gorm:"not null"`
	RiskTimelineTypeID *uint     `json:"risk_timeline_type_id" gorm:"not null"`
	Comments           string    `json:"comments" gorm:"type:text;not null"`
	RiskStatusID       *uint     `json:"risk_status_id" gorm:"not null"`
	CreatedBYID        uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt          time.Time `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	RiskRegister     *RiskRegister     `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskTimelineType *RiskTimelineType `gorm:"foreignKey:RiskTimelineTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskStatus       *RiskStatus       `gorm:"foreignKey:RiskStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`

	CreatedByUser *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (RiskComment) TableName() string {
	return "risk_comments"
}
