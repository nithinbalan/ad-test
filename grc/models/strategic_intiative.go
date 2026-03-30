package models

import "time"

type StrategicIntiative struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn           string     `json:"title_en" gorm:"type:nvarchar(MAX);not null" binding:"required"`
	TitleAr           string     `json:"title_ar" gorm:"type:nvarchar(MAX);not null" binding:"required"`
	DescriptionEn     *string    `json:"description_en" gorm:"type:nvarchar(MAX);default:null"`
	DescriptionAr     *string    `json:"description_ar" gorm:"type:nvarchar(MAX);default:null"`
	StartDate         *time.Time `json:"start_date" gorm:"type:date;default:null"`
	EndDate           *time.Time `json:"end_date" gorm:"type:date;default:null"`
	ResponsibleUserID uint       `json:"responsible_user_id" gorm:"default:null"`
	ResponsibleUnitID uint       `json:"responsible_unit_id" gorm:"default:null"`

	ObjectivesEn      *string `json:"objectives_en" gorm:"type:nvarchar(MAX);default:null"`
	ObjectivesAr      *string `json:"objectives_ar" gorm:"type:nvarchar(MAX);default:null"`
	RecommendationsEn *string `json:"recommendations_en" gorm:"type:nvarchar(MAX);default:null"`
	RecommendationsAr *string `json:"recommendations_ar" gorm:"type:nvarchar(MAX);default:null"`
	OwnerDepartmentID uint    `json:"owner_department_id" gorm:"default:null"`

	StrategicInitiativeStatusID uint `json:"strategic_initiative_status_id" gorm:"column:strategic_initiative_status_id;not null" binding:"required"`
	CompletionPercentage        *float64 `json:"completion_percentage" gorm:"type:decimal(5,2);default:null"`
	KPIs 						*string `json:"kpis" gorm:"column:kpis;type:nvarchar(MAX);default:null"`
	ExpectedBenefitsEn          *string  `json:"expected_benefits_en" gorm:"type:nvarchar(MAX);default:null"`
	ExpectedBenefitsAr          *string  `json:"expected_benefits_ar" gorm:"type:nvarchar(MAX);default:null"`
	Budget                      *float64 `json:"budget" gorm:"type:decimal(18,2);default:null"`
	ResourceRequirementsEn      *string  `json:"resource_requirements_en" gorm:"type:nvarchar(MAX);default:null"`
	ResourceRequirementsAr      *string  `json:"resource_requirements_ar" gorm:"type:nvarchar(MAX);default:null"`
	DependenciesEn              *string  `json:"dependencies_en" gorm:"type:nvarchar(MAX);default:null"`
	DependenciesAr              *string  `json:"dependencies_ar" gorm:"type:nvarchar(MAX);default:null"`

	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime" binding:"omitempty"`
	CreatedBYID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID *uint      `json:"updated_by_id" gorm:"default:null" binding:"omitempty"`

	ResponsibleUser           User                      `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ResponsibleUnit           OrganizationUnit          `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	OwnerDepartment           OrganizationUnit          `gorm:"foreignKey:OwnerDepartmentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	StrategicInitiativeStatus StrategicInitiativeStatus `gorm:"foreignKey:StrategicInitiativeStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser             User                      `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser             *User                     `gorm:"foreignKey:UpdatedBYID"`
}

func (StrategicIntiative) TableName() string {
	return "strategic_intiatives"
}
