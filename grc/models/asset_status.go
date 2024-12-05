package models

type AssetStatus struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn   string `json:"title_en" binding:"required" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required" gorm:"type:nvarchar(100);not null"`
	ColorCode string `json:"color_code" binding:"required" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" binding:"required,oneof=active available maintenance disposed" gorm:"type:varchar(20);not null;default:'active';check:type IN ('active', 'available', 'maintenance', 'disposed')"`
	Status    string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active', 'inactive')"`
}

// TableName overrides the default table name used by Gorm
func (AssetStatus) TableName() string {
	return "asset_statuses"
}
