package models

import "time"

type ComplianceStatusUpdate struct {
	ID                 uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceID       uint      `json:"compliance_id" gorm:"not null"`
	ComplianceStatusID uint      `json:"compliance_status_id" gorm:"not null"`
	Comment            string    `json:"comment" binding:"required,max=100" gorm:"type:nvarchar(max);not null"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`
	CreatedBy          uint      `json:"created_by" gorm:"not null"`

	//Relationships
	ComplianceRegister *ComplianceRegister `gorm:"foreignkey:ComplianceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceStatus   *ComplianceStatus   `gorm:"foreignkey:ComplianceStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by ComplianceStatusUpdate to `compliance_status_updates`
func (ComplianceStatusUpdate) TableName() string {
	return "compliance_status_updates"
}
