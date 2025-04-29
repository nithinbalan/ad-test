package models

import "time"

type GRCReport struct {
	ID            uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Title         string    `json:"title" binding:"required,max=255"`
	StartDate     time.Time `json:"start_date" binding:"required"`
	EndDate       time.Time `json:"end_date" binding:"required"`
	ReviewUserID  *uint     `json:"review_user_id"`
	ApproveUserID *uint     `json:"approve_user_id"`
	CreatedBYID   uint      `json:"created_by_id" binding:"required"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	CreatedByUser *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ReviewUser    *User `gorm:"foreignKey:ReviewUserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ApproveUser   *User `gorm:"foreignKey:ApproveUserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`

	Sections []GRCReportSection `gorm:"foreignKey:GRCReportID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (GRCReport) TableName() string {
	return "grc_reports"
}
