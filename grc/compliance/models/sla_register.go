package models

import (
	"time"

	masterModel "jethur/grc/masters/models"
	organizationModel "jethur/grc/organization/models"
	settingsModel "jethur/grc/settings/models"
)


type SLARegister struct {
    ID                 uint        `json:"id" gorm:"primary_key;auto_increment"`
    ReferenceNumber    string      `json:"reference_number" gorm:"type:varchar(45)"` 
    SLACategoryID      uint        `json:"sla_category_id" binding:"required" gorm:"column:sla_category_id"`
    SeverityLevelID    uint        `json:"severity_level_id" binding:"required" gorm:"column:severity_level_id"`
    ResponsibleUnitID  uint        `json:"responsible_unit_id" binding:"required" gorm:"column:responsible_unit_id"`
    ResponsibleUserID  uint        `json:"responsible_user_id" binding:"required" gorm:"column:responsible_user_id"`
    LifecycleStatusID  uint        `json:"lifecycle_status_id" binding:"required" gorm:"column:lifecycle_status_id"`
    TitleEn            string      `json:"title_en" binding:"required,max=250" gorm:"type:varchar(250);not null"`
    TitleAr            string      `json:"title_ar" binding:"required,max=250" gorm:"type:nvarchar(250);not null"`
	Description        string      `json:"description" gorm:"type:varchar(max)"` 
    Client             string      `json:"client" gorm:"type:varchar(100)"`
    CreatedBy          uint        `json:"created_by"`
    UpdatedBy          uint        `json:"updated_by"`
    CreatedAt          time.Time   `json:"created_at" gorm:"type:datetime;not null"`
    UpdatedAt          time.Time   `json:"updated_at" gorm:"type:datetime;not null"`

    // Additional fields for related tables
    SLACategoryTitleEn      string `json:"sla_category_title_en,omitempty" gorm:"-"`
    SLACategoryTitleAr      string `json:"sla_category_title_ar,omitempty" gorm:"-"`
    SeverityLevelTitleEn    string `json:"severity_level_title_en,omitempty" gorm:"-"`
    SeverityLevelTitleAr    string `json:"severity_level_title_ar,omitempty" gorm:"-"`
    ResponsibleUnitTitleEn  string `json:"responsible_unit_title_en,omitempty" gorm:"-"`
    ResponsibleUnitTitleAr  string `json:"responsible_unit_title_ar,omitempty" gorm:"-"`
    ResponsibleUserFirstNameEn string `json:"responsible_user_first_name_en,omitempty" gorm:"-"`
    ResponsibleUserFirstNameAr string `json:"responsible_user_first_name_ar,omitempty" gorm:"-"`
    ResponsibleUserLastNameEn  string `json:"responsible_user_last_name_en,omitempty" gorm:"-"`
    ResponsibleUserLastNameAr  string `json:"responsible_user_last_name_ar,omitempty" gorm:"-"`
    LifecycleStatusTitleEn   string `json:"lifecycle_status_title_en,omitempty" gorm:"-"`
    LifecycleStatusTitleAr   string `json:"lifecycle_status_title_ar,omitempty" gorm:"-"`
    CreatedByFirstNameEn string `json:"created_by_first_name_en,omitempty" gorm:"-"`
    CreatedByFirstNameAr string `json:"created_by_first_name_ar,omitempty" gorm:"-"`
    CreatedByLastNameEn  string `json:"created_by_last_name_en,omitempty" gorm:"-"`
    CreatedByLastNameAr  string `json:"created_by_last_name_ar,omitempty" gorm:"-"`
    UpdatedByFirstNameEn string `json:"updated_by_first_name_en,omitempty" gorm:"-"`
    UpdatedByFirstNameAr string `json:"updated_by_first_name_ar,omitempty" gorm:"-"`
    UpdatedByLastNameEn  string `json:"updated_by_last_name_en,omitempty" gorm:"-"`
    UpdatedByLastNameAr  string `json:"updated_by_last_name_ar,omitempty" gorm:"-"`

    // Relationships
    SLACategory        *SLACategory `gorm:"foreignKey:SLACategoryID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
    SeverityLevel      *masterModel.SeverityLevel `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
    ResponsibleUnit    *settingsModel.OrganizationUnit `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
    ResponsibleUser    *organizationModel.User `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
    LifecycleStatus    *masterModel.LifecycleStatus `gorm:"foreignKey:LifecycleStatusID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser      *organizationModel.User      `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	UpdatedByUser      *organizationModel.User      `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName overrides the table name used by SLARegister to `sla_registers`
func (SLARegister) TableName() string {
	return "sla_registers"
}
