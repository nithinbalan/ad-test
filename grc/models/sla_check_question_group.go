package models

import (
	"time"
)

type SLACheckQuestionGroup struct {
	ID             uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	SLACheckID     uint       `json:"sla_check_id" binding:"required" gorm:"not null"`
	TitleEn        string     `json:"title_en" binding:"required,max=500" gorm:"type:varchar(500);not null"`
	TitleAr        string     `json:"title_ar" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	TotalCount     int        `json:"total_count" binding:"required" gorm:"not null;default:0"`
	CompletedCount int        `json:"completed_count" binding:"required" gorm:"not null;default:0"`
	CreatedAt      time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`
	CreatedBy      uint       `json:"created_by" binding:"required" gorm:"not null"`
	UpdatedBy      *uint      `json:"updated_by,omitempty" gorm:"default:null"`

	// Relationships
	SLACheck      *SLACheck `gorm:"foreignkey:SLACheckID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User     `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser *User     `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by SLACheckQuestionGroup to `sla_check_question_groups`
func (SLACheckQuestionGroup) TableName() string {
	return "sla_check_question_groups"
}
