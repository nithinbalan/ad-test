package models

import (
	"time"
)

type SLARegister struct {
	ID                 uint       `json:"id" gorm:"primary_key;auto_increment"`
	ReferenceNumber    string     `json:"reference_number" binding:"omitempty,max=45" gorm:"size:45"`
	SLACategoryID      uint       `json:"sla_category_id" binding:"required" gorm:"column:sla_category_id"`
	SeverityLevelID    uint       `json:"severity_level_id" binding:"required" gorm:"column:severity_level_id"`
	ResponsibleUnitID  uint       `json:"responsible_unit_id" binding:"required" gorm:"column:responsible_unit_id"`
	ResponsibleUserID  uint       `json:"responsible_user_id" binding:"required" gorm:"column:responsible_user_id"`
	LifecycleStatusID  uint       `json:"lifecycle_status_id" binding:"required" gorm:"column:lifecycle_status_id"`
	ComplianceStatusID *uint      `json:"compliance_status_id" gorm:"default:null"`
	TitleEn            string     `json:"title_en" binding:"required,max=250" gorm:"type:varchar(250);not null"`
	TitleAr            string     `json:"title_ar" binding:"required,max=250" gorm:"type:nvarchar(250);not null"`
	Description        *string    `json:"description" binding:"max=250" gorm:"type:varchar(max)"`
	Client             string     `json:"client" gorm:"type:varchar(250)"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID        uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID        uint       `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	SLACategory      *SLACategory      `gorm:"foreignKey:SLACategoryID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	SeverityLevel    *SeverityLevel    `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	ResponsibleUnit  *OrganizationUnit `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	ResponsibleUser  *User             `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	LifecycleStatus  *LifecycleStatus  `gorm:"foreignKey:LifecycleStatusID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	ComplianceStatus *ComplianceStatus `gorm:"foreignKey:ComplianceStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Documents        []*SLADocument    `gorm:"foreignKey:SLARegisterID"`
	Checks           []*SLACheck       `gorm:"foreignKey:SLARegisterID"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser    *User             `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by SLARegister to `sla_register`
func (SLARegister) TableName() string {
	return "sla_register"
}
