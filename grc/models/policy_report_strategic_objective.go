package models

import "time"

type PolicyReportStrategicObjective struct {
	ID                                uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PolicyReportID                    uint      `json:"policy_report_id" gorm:"not null;index"`
	StrategicObjectiveTitleEn         string    `json:"strategic_objective_title_en" gorm:"type:nvarchar(255);not null"`
	StrategicObjectiveTitleAr         string    `json:"strategic_objective_title_ar" gorm:"type:nvarchar(255)"`
	DescriptionEn                     *string   `json:"description_en" gorm:"type:nvarchar(max)"`
	DescriptionAr                     *string   `json:"description_ar" gorm:"type:nvarchar(max)"`
	StartDate                         *string   `json:"start_date" gorm:"type:varchar(50)"`
	EndDate                           *string   `json:"end_date" gorm:"type:varchar(50)"`
	StrategicObjectiveCategoryTitleEn string    `json:"strategic_objective_category_title_en" gorm:"type:nvarchar(255)"`
	StrategicObjectiveCategoryTitleAr string    `json:"strategic_objective_category_title_ar" gorm:"type:nvarchar(255)"`
	ResponsibleUserNameEn             string    `json:"responsible_user_name_en" gorm:"type:nvarchar(255)"`
	ResponsibleUserNameAr             string    `json:"responsible_user_name_ar" gorm:"type:nvarchar(255)"`
	ResponsibleUnitTitleEn            string    `json:"responsible_unit_title_en" gorm:"type:nvarchar(255)"`
	ResponsibleUnitTitleAr            string    `json:"responsible_unit_title_ar" gorm:"type:nvarchar(255)"`
	CreatedByID                       uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt                         time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	PolicyReport  *PolicyReport `gorm:"foreignKey:PolicyReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PolicyReportStrategicObjective) TableName() string {
	return "policy_report_strategic_objectives"
}
