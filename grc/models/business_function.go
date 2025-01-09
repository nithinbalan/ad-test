package models

import (
	"time"
)

type BusinessFunction struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber   string     `json:"reference_number" binding:"omitempty,max=45" gorm:"size:45"`
	ResponsibleUserID uint       `json:"responsible_user_id" binding:"required" gorm:"not null"`
	SeverityLevelID   *uint      `json:"severity_level_id" binding:"required"`
	Title             string     `json:"title" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	Description       *string    `json:"description" gorm:"type:nvarchar(max)"`
	BcpStatus         bool       `json:"bcp_status"`
	ReviewFrequencyID *uint      `json:"review_frequency_id"`
	NextReviewDate    *time.Time `json:"next_review_date"`
	LastReviewDate    *time.Time `json:"last_review_date"`
	Rto               *uint      `json:"rto"`
	RtoUnitId         *uint      `json:"rto_unit_id"`
	Mtpd              *uint      `json:"mtpd"`
	MtpdUnitId        *uint      `json:"mtpd_unit_id"`
	Rpo               *uint      `json:"rpo"`
	RpoUnitId         *uint      `json:"rpo_unit_id"`
	CreatedBYID       uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID       *uint      `json:"updated_by_id"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	SeverityLevel   *SeverityLevel             `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewFrequency *ReviewFrequency           `gorm:"foreignKey:ReviewFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser *User                      `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Processes       []*BusinessFunctionProcess `gorm:"foreignKey:BusinessFunctionID"`
	CreatedByUser   *User                      `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser   *User                      `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by BusinessFunction to `business_functions`
func (BusinessFunction) TableName() string {
	return "business_functions"
}
