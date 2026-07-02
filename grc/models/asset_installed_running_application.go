package models

import "time"

type AssetInstalledRunningApplication struct {
	ID                              uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AssetRegisterID                 uint      `json:"asset_register_id" gorm:"not null"`
	InstalledRunningApplicationID   uint      `json:"installed_running_application_id" gorm:"not null"`
	CreatedBYID                     uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                       time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AssetRegister                 *AssetRegister                 `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	InstalledRunningApplication   *InstalledRunningApplication   `json:"installed_running_application,omitempty" gorm:"foreignKey:InstalledRunningApplicationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser                  *User                         `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (AssetInstalledRunningApplication) TableName() string {
	return "asset_installed_running_applications"
}