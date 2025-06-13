package models

import "time"

type StrategicObjective struct {
	ID                           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn                      string     `json:"title_en" gorm:"type:nvarchar(MAX);not null" binding:"required"`
	TitleAr                      string     `json:"title_ar" gorm:"type:nvarchar(MAX);not null" binding:"required"`
	StrategicObjectiveCategoryID uint       `json:"strategic_objective_category_id" binding:"required" gorm:"not null"`
	ReferenceNumber              string     `json:"reference_number" gorm:"type:varchar(45);not null"`     // ✅ new
	DescriptionEn                *string    `json:"description_en" gorm:"type:nvarchar(MAX);default:null"` // ✅ new
	DescriptionAr                *string    `json:"description_ar" gorm:"type:nvarchar(MAX);default:null"` // ✅ new
	StartDate                    *time.Time `json:"start_date" gorm:"type:date"`                           // ✅ new
	EndDate                      *time.Time `json:"end_date" gorm:"type:date"`                             // ✅ new
	ResponsibleUserID            *uint      `json:"responsible_user_id" gorm:"default:null"`               // ✅ new
	ResponsibleUnitID            *uint      `json:"responsible_unit_id" gorm:"default:null"`               // ✅ new
	ReviewUserID                 *uint      `json:"review_user_id" gorm:"default:null"`                    // ✅ new
	ReviewFrequencyID            *uint      `json:"review_frequency_id" gorm:"default:null"`               // ✅ new
	Status                       string     `json:"status" gorm:"type:varchar(10);not null;default:'active'" binding:"required,oneof=active inactive"`
	CreatedAt                    time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                    *time.Time `json:"updated_at" gorm:"type:datetime" binding:"omitempty"`
	CreatedByID                  uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID                  *uint      `json:"updated_by_id" gorm:"default:null" binding:"omitempty"`

	// Relationships
	StrategicObjectiveCategory *StrategicObjectiveCategory                `gorm:"foreignKey:StrategicObjectiveCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser            *User                                      `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ResponsibleUnit            *OrganizationUnit                          `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ReviewUser                 *User                                      `gorm:"foreignKey:ReviewUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ReviewFrequency            *ReviewFrequency                           `gorm:"foreignKey:ReviewFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CreatedByUser              *User                                      `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser              *User                                      `gorm:"foreignKey:UpdatedByID"`
	ComplianceRequirements     []*ComplianceRequirementStrategicObjective `gorm:"foreignKey:StrategicObjectiveID"`
}

func (StrategicObjective) TableName() string {
	return "strategic_objectives"
}
