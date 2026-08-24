package models

import (
	"time"
)

type KPIRegister struct {
	ID                    uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Identifier            *string `json:"identifier" gorm:"type:varchar(100);default:null"`
	Type                  string  `json:"type" gorm:"type:varchar(20);not null;default:'manual'" binding:"oneof=automated manual"`
	IsKRI                 bool    `json:"is_kri" gorm:"column:is_kri;type:bit;not null;default:false"`
	ReferenceNumber       *string `json:"reference_number" gorm:"type:varchar(45);not null;unique" binding:"required"`
	Title                 string  `json:"title" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	Description           *string `json:"description" gorm:"type:nvarchar(max)"`
	KPICategoryID         uint    `json:"kpi_category_id"`
	ResponsibleUserID     *uint   `json:"responsible_user_id"`
	KPIStatusID           uint    `json:"kpi_status_id"`
	StatusID              *uint   `json:"status_id"`
	KPIAchievementLogicID *uint   `json:"kpi_achievement_logic_id"`

	ReviewFrequencyID *uint    `json:"review_frequency_id"`
	CalculationTypeID *uint    `json:"calculation_type_id"`
	Formula           *string  `json:"formula"`
	UnitID            *uint    `json:"unit_id"`
	Source            *string  `json:"source"`
	ActualValue       *float64 `json:"actual_value" gorm:"type:decimal(18,2)"`
	PreviousValue     *float64 `json:"previous_value" gorm:"type:decimal(18,2)"`
	KPIReviewStatusID *uint    `json:"kpi_review_status_id"`

	TargetValue       *float64 `json:"target_value"`
	TargetDescription *string  `json:"target_description"`

	StartDate  *time.Time `json:"start_date" gorm:"type:date;default:null"`
	EndDate    *time.Time `json:"end_date" gorm:"type:date;default:null"`
	TargetDate *time.Time `json:"target_date" gorm:"type:date"`

	Remarks *string `json:"remarks" gorm:"type:nvarchar(max)"`

	Achievement *float64 `json:"achievement" gorm:"type:decimal(18,2)"`

	EntryTypeID      *uint      `json:"entry_type_id"`
	DataSourceID     *uint      `json:"data_source_id"`
	LastReviewedDate *time.Time `json:"last_reviewed_date" gorm:"type:date"`
	NextReviewDate   *time.Time `json:"next_review_date" gorm:"type:date;not null"`
	LastReviewNote   *string    `json:"last_review_note" gorm:"type:nvarchar(max)"`
	ReviewUserID     *uint      `json:"review_user_id"`

	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	KPICategory         *KPICategory            `gorm:"foreignKey:KPICategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser     *User                   `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	KPIStatus           *KPIStatus              `gorm:"foreignKey:KPIStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	KPIReviewStatus     *KPIReviewStatus        `gorm:"foreignKey:KPIReviewStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewFrequency     *ReviewFrequency        `gorm:"foreignKey:ReviewFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CalculationType     *CalculationType        `gorm:"foreignKey:CalculationTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Unit                *Unit                   `gorm:"foreignKey:UnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser       *User                   `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser       *User                   `gorm:"foreignKey:UpdatedBYID"`
	DataInputs          []*KPIRegisterDataInput `gorm:"foreignKey:KPIRegisterID"`
	Status              *LifecycleStatus        `gorm:"foreignKey:StatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	KPIAchievementLogic *KPIAchievementLogic    `gorm:"foreignKey:KPIAchievementLogicID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewUser          *User                   `gorm:"foreignKey:ReviewUserID"`
}

// TableName overrides the table name used by KPIRegister to `kpi_register`
func (KPIRegister) TableName() string {
	return "kpi_register"
}
