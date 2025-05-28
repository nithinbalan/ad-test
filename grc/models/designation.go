package models

type Designation struct {
	ID               uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn          string `json:"title_en" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	TitleAr          string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Status           string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
	JobDescriptionID *uint  `json:"job_description_id" gorm:"default:null"`

	JobDescription *JobDescription `gorm:"foreignKey:JobDescriptionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by Designations to `designation`
func (Designation) TableName() string {
	return "designations"
}
