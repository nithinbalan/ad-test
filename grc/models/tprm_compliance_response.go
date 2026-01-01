package models

import (
	"time"
)

type TPRMComplianceResponse struct {
	ID          uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn     string     `json:"title_en" binding:"required,max=1000" gorm:"type:nvarchar(1000);not null"`
	TitleAr     string     `json:"title_ar" binding:"required,max=1000" gorm:"type:nvarchar(1000);not null"`
	Score       int        `json:"score"`
	Status      string     `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active', 'inactive')"`
	ColorCode   *string    `json:"color_code"`
	CreatedBYID uint       `json:"created_by_id" gorm:"not null" binding:"required"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	CreatedByUser *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by TPRMComplianceResponse to `tprm_compliance_responses`
func (TPRMComplianceResponse) TableName() string {
	return "tprm_compliance_responses"
}
