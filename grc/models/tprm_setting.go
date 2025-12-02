package models

import "time"

type TPRMSetting struct {
	ID               uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	IsExport         bool       `json:"is_export"`
	IsExternalVendor bool       `json:"is_external_vendor"`
	UpdatedBYID      *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt        *time.Time `json:"updated_at" gorm:"type:datetime"`
}

func (TPRMSetting) TableName() string {
	return "tprm_settings"
}
