package models

type PlannedAdvisoryInternalStatus struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required" gorm:"type:nvarchar(max);not null"`
	TitleAr string `json:"title_ar" binding:"required" gorm:"type:nvarchar(max);not null"`
	Status  string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by PlannedAdvisoryInternalStatus to `planned_advisory_internal_statuses`
func (PlannedAdvisoryInternalStatus) TableName() string {
	return "planned_advisory_internal_statuses"
}
