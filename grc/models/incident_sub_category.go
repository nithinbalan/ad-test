package models

type IncidentSubCategory struct {
	ID                 uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentCategoryID uint   `json:"incident_category_id" binding:"required" gorm:"not null"`
	TitleEn            string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr            string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Status             string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`

	// Relationships
	IncidentCategory *IncidentCategory `gorm:"foreignkey:IncidentCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by IncidentSubCategory to `incident_sub_categories`
func (IncidentSubCategory) TableName() string {
	return "incident_sub_categories"
}
