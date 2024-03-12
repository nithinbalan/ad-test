package models

import (
	"time"

	masterModel "jethur/grc/masters/models"
	organizationModel "jethur/grc/organization/models"
	settingsModel "jethur/grc/settings/models"
)


type SLARegister struct {
    ID                 uint        `json:"id" gorm:"primary_key;auto_increment"`
    SLACategoryID      uint        `json:"sla_category_id" binding:"required" gorm:"column:sla_category_id"`
    SeverityLevelID    uint        `json:"severity_level_id" binding:"required" gorm:"column:severity_level_id"`
    ResponsibleUnitID  uint        `json:"responsible_unit_id" binding:"required" gorm:"column:responsible_unit_id"`
    ResponsibleUserID  uint        `json:"responsible_user_id" binding:"required" gorm:"column:responsible_user_id"`
    LifecycleStatusID  uint        `json:"lifecycle_status_id" binding:"required" gorm:"column:lifecycle_status_id"`
    TitleEn            string      `json:"title_en" binding:"required,max=100" gorm:"type:varchar(100);not null"`
    TitleAr            string      `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Description        string      `json:"description" gorm:"type:varchar(max)"` 
    Client             string      `json:"client" gorm:"type:varchar(100)"`
    CreatedBy          uint        `json:"created_by"`
    UpdatedBy          uint        `json:"updated_by"`
    CreatedAt          time.Time   `json:"created_at" gorm:"type:datetime;not null"`
    UpdatedAt          time.Time   `json:"updated_at" gorm:"type:datetime;not null"`

    // Relationships
    SLACategory        *SLACategory `gorm:"foreignKey:SLACategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
    SeverityLevel      *masterModel.SeverityLevel `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
    ResponsibleUnit    *settingsModel.OrganizationUnit `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
    ResponsibleUser    *organizationModel.User `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
    LifecycleStatus    *masterModel.LifecycleStatus `gorm:"foreignKey:LifecycleStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedByUser      *organizationModel.User                `gorm:"foreignKey:CreatedBy"`
    UpdatedByUser      *organizationModel.User                `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by SLARegister to `sla_registers`
func (SLARegister) TableName() string {
	return "sla_registers"
}
