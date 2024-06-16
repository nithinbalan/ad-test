package models

import (
	"time"
)

type Control struct {
	ID                            uint       `json:"id" gorm:"primaryKey;autoIncrement" binding:"-"`
	ControlTypeID                 uint       `json:"control_type_id" gorm:"not null" binding:"required"`
	ControlCategoryID             uint       `json:"control_category_id" gorm:"not null" binding:"required"`
	OwnerUnitID                   *uint      `json:"owner_unit_id,omitempty" gorm:"default:null" binding:"omitempty"`
	OwnerID                       *uint      `json:"owner_id,omitempty" gorm:"default:null" binding:"omitempty"`
	ControlImplementationStatusID *uint      `json:"control_implementation_status_id,omitempty" gorm:"default:null" binding:"omitempty"`
	ControlEfficiencyStatusID     *uint      `json:"control_efficiency_status_id,omitempty" gorm:"default:null" binding:"omitempty"`
	ReferenceNumber               string     `json:"reference_number" gorm:"type:varchar(45);not null;unique" binding:"required"`
	Code                          *string    `json:"code,omitempty" gorm:"type:varchar(45);default:null" binding:"omitempty"`
	TitleEn                       string     `json:"title_en" gorm:"type:nvarchar(500);not null" binding:"required"`
	TitleAr                       string     `json:"title_ar" gorm:"type:nvarchar(500);not null" binding:"required"`
	Description                   *string    `json:"description,omitempty" gorm:"type:nvarchar(MAX);default:null" binding:"omitempty"`
	Objective                     *string    `json:"objective,omitempty" gorm:"type:nvarchar(MAX);default:null" binding:"omitempty"`
	ReviewFrequencyID             *uint      `json:"review_frequency_id,omitempty" gorm:"default:null" binding:"omitempty"`
	LastReviewDate                *time.Time `json:"last_review_date,omitempty" gorm:"type:date;default:null" binding:"omitempty"`
	NextReviewDate                *time.Time `json:"next_review_date,omitempty" gorm:"type:date;default:null" binding:"omitempty"`
	LastReviewNote                *string    `json:"last_review_note,omitempty" gorm:"type:nvarchar(MAX);default:null" binding:"omitempty"`
	TestResponsibleUserID         uint       `json:"test_responsible_user_id" gorm:"not null" binding:"required"`
	TestFrequencyID               *uint      `json:"test_frequency_id,omitempty" gorm:"default:null" binding:"omitempty"`
	LastTestDate                  *time.Time `json:"last_test_date,omitempty" gorm:"type:date;default:null" binding:"omitempty"`
	NextTestDate                  *time.Time `json:"next_test_date,omitempty" gorm:"type:date;default:null" binding:"omitempty"`
	ImplementedAt                 *time.Time `json:"implemented_at,omitempty" gorm:"type:date;default:null" binding:"omitempty"`
	CreatedByID                   uint       `json:"created_by_id" gorm:"not null" binding:"required"`
	UpdatedBy                     *uint      `json:"updated_by_id,omitempty" gorm:"default:null" binding:"omitempty"`
	CreatedAt                     time.Time  `json:"created_at" gorm:"type:datetime;not null" binding:"-"`
	UpdatedAt                     *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null" binding:"omitempty"`
	LifecycleStatusID             uint       `json:"lifecycle_status_id" gorm:"not null"`

	// Relationships
	ControlType                 *ControlType                 `gorm:"foreignKey:ControlTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ControlCategory             *ControlCategory             `gorm:"foreignKey:ControlCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OwnerUnit                   *OrganizationUnit            `gorm:"foreignKey:OwnerUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Owner                       *User                        `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ControlImplementationStatus *ControlImplementationStatus `gorm:"foreignKey:ControlImplementationStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ControlEfficiencyStatus     *ControlEfficiencyStatus     `gorm:"foreignKey:ControlEfficiencyStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewFrequency             *ReviewFrequency             `gorm:"foreignKey:ReviewFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	TestResponsibleUser         *User                        `gorm:"foreignKey:TestResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	TestFrequency               *ReviewFrequency             `gorm:"foreignKey:TestFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	LifecycleStatus             *LifecycleStatus             `gorm:"foreignKey:LifecycleStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser               *User                        `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser               *User                        `gorm:"foreignKey:UpdatedBy"`
}

func (Control) TableName() string {
	return "controls"
}
