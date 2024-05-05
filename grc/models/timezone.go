package models

// Timezone represents the timezones table structure in GORM.
type Timezone struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id" binding:"required"`
	Name        string `gorm:"type:varchar(255); not null" json:"name" binding:"required"`
	Offset      string `gorm:"type:varchar(100); not null" json:"offset" binding:"required"`
	DisplayName string `gorm:"type:varchar(255); not null" json:"display_name" binding:"required"`
}

// TableName overrides the table name used by Timezone to `timezones`
func (Timezone) TableName() string {
	return "timezones"
}
