package models

type IncidentStatus struct {
	ID            uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn       string `json:"title_en" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	TitleAr       string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DescriptionEn string `json:"description_en" binding:"required" gorm:"type:text"`
	DescriptionAr string `json:"description_ar" binding:"required" gorm:"type:text"`
	ColorCode     string `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Type          string `json:"type" binding:"required,max=25,oneof=open investigating in_progress resolved closed pending_review escalated re_opened deferred mitigated" gorm:"type:varchar(25);not null"`
	Status        string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by IncidentStatus to `incident_statuses`
func (IncidentStatus) TableName() string {
	return "incident_statuses"
}
