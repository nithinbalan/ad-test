package models

import (
	"time"

	masterModel "jethur/grc/masters/models"
	organizationModel "jethur/grc/organization/models"
	settingsModel "jethur/grc/settings/models"
)


type SLARegister struct {
    ID                 uint        `json:"id" gorm:"primary_key;auto_increment"`
    ReferenceID        string      `json:"reference_id" gorm:"type:varchar(45)"` 
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
	CreatedByUser   *organizationModel.User            `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser   *organizationModel.User            `gorm:"foreignKey:UpdatedBy"`
}

type SLARegisterWithAllField struct {
    ID                 uint        `json:"id"`
    ReferenceID        string      `json:"reference_id"` 
    SLACategoryID      uint        `json:"sla_category_id"`
    SeverityLevelID    uint        `json:"severity_level_id"`
    ResponsibleUnitID  uint        `json:"responsible_unit_id"`
    ResponsibleUserID  uint        `json:"responsible_user_id"`
    LifecycleStatusID  uint        `json:"lifecycle_status_id"`
    TitleEn            string      `json:"title_en"`
    TitleAr            string      `json:"title_ar"`
    Description        string      `json:"description"` 
    Client             string      `json:"client"`
    CreatedBy          uint        `json:"created_by"`
    UpdatedBy          uint        `json:"updated_by"`
    CreatedAt          time.Time   `json:"created_at"`
    UpdatedAt          time.Time   `json:"updated_at"`

    // Additional fields for related tables
    SLACategoryTitleEn      string `json:"sla_category_title_en"`
    SLACategoryTitleAr      string `json:"sla_category_title_ar"`
    SeverityLevelTitleEn    string `json:"severity_level_title_en"`
    SeverityLevelTitleAr    string `json:"severity_level_title_ar"`
    ResponsibleUnitTitleEn  string `json:"responsible_unit_title_en"`
    ResponsibleUnitTitleAr  string `json:"responsible_unit_title_ar"`
    ResponsibleUserFirstNameEn string `json:"responsible_user_first_name_en"`
    ResponsibleUserFirstNameAr string `json:"responsible_user_first_name_ar"`
    ResponsibleUserLastNameEn  string `json:"responsible_user_last_name_en"`
    ResponsibleUserLastNameAr  string `json:"responsible_user_last_name_ar"`
    LifecycleStatusTitleEn   string `json:"lifecycle_status_title_en"`
    LifecycleStatusTitleAr   string `json:"lifecycle_status_title_ar"`
    CreatedByFirstNameEn string `json:"created_by_first_name_en"`
    CreatedByFirstNameAr string `json:"created_by_first_name_ar"`
    CreatedByLastNameEn  string `json:"created_by_last_name_en"`
    CreatedByLastNameAr  string `json:"created_by_last_name_ar"`
    UpdatedByFirstNameEn string `json:"updated_by_first_name_en"`
    UpdatedByFirstNameAr string `json:"updated_by_first_name_ar"`
    UpdatedByLastNameEn  string `json:"updated_by_last_name_en"`
    UpdatedByLastNameAr  string `json:"updated_by_last_name_ar"`

}

// TableName overrides the table name used by SLARegister to `sla_registers`
func (SLARegister) TableName() string {
	return "sla_registers"
}
