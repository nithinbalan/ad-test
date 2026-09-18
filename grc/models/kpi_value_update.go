package models

import (
	"time"
)

type KPIValueUpdate struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KPIRegisterID     uint       `json:"kpi_register_id" binding:"required" gorm:"not null"`
	ResponsibleUserID *uint      `json:"responsible_user_id"`
	ReviewUserID      *uint      `json:"review_user_id"`
	ReviewedDate      *time.Time `json:"reviewed_date" binding:"required" gorm:"type:date;not null"`
	UnitID            uint       `json:"unit_id" binding:"required" gorm:"not null"`
	Value             *float64   `json:"value" gorm:"type:decimal(18,2);not null"`
	Remarks           *string    `json:"remarks" gorm:"type:nvarchar(max)"`
	StartDate         *time.Time `json:"start_date" gorm:"type:date;default:null"`
	EndDate           *time.Time `json:"end_date" gorm:"type:date;default:null"`
	IsLatest          bool       `json:"is_latest" gorm:"type:bit;not null"`
	PeriodLabel       *string    `json:"period_label"`
	KPIReviewStatusID *uint      `json:"kpi_review_status_id"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID       uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID       *uint      `json:"updated_by_id" gorm:"default:null"`
	CorrectiveAction  *string    `json:"corrective_action"`
	ApprovedAt        *time.Time `json:"approved_at" gorm:"type:datetime;default:null"`
	RevertedAt        *time.Time `json:"reverted_at" gorm:"type:datetime;default:null"`

	// Relationships
	KPIRegister          *KPIRegister               `gorm:"foreignKey:KPIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	KPIReviewStatus      *KPIReviewStatus           `gorm:"foreignKey:KPIReviewStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Unit                 *Unit                      `gorm:"foreignKey:UnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Documents            []KPIValueUpdateDocument   `gorm:"foreignKey:KPIValueUpdateID"`
	CreatedByUser        *User                      `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser        *User                      `gorm:"foreignKey:UpdatedBYID"`
	ResponsibleUser      *User                      `gorm:"foreignKey:ResponsibleUserID"`
	ReviewUser           *User                      `gorm:"foreignKey:ReviewUserID"`
	DataInputs           []*KPIValueUpdateDataInput `gorm:"foreignKey:KPIValueUpdateID"`
	ValueUpdateHistories []*KPIValueUpdateHistory   `gorm:"foreignKey:KPIValueUpdateID"`
	ValueUpdateReviews   []*KPIValueUpdateReview    `gorm:"foreignKey:KPIValueUpdateID"`
}

// TableName overrides the table name used by KPIValueUpdate to `kpi_value_updates`
func (KPIValueUpdate) TableName() string {
	return "kpi_value_updates"
}
