package models

import (
	"time"
)

type ProcessExternalDependency struct {
	ID                       uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ProcessID                uint       `json:"process_id" gorm:"not null"`
	EntityName               string     `json:"entity_name" gorm:"type:nvarchar(1000);not null"`
	EntityTypeID             *uint      `json:"entity_type_id" gorm:"default:null"`
	NeededSupportDescription *string    `json:"needed_support_description" gorm:"type:nvarchar(MAX);default:null"`
	DependencyTierID         *uint      `json:"dependency_tier_id" gorm:"default:null"`
	SLARegisterID            *uint      `json:"sla_register_id" gorm:"default:null"`
	VendorRTO                *int       `json:"vendor_rto" gorm:"default:null"`
	VendorRTOUnitID          *uint      `json:"vendor_rto_unit_id" gorm:"default:null"`
	InternalStakeholder      *string    `json:"internal_stakeholder" gorm:"type:nvarchar(MAX);default:null"`
	ContactNumber            *string    `json:"contact_number" gorm:"type:nvarchar(100);default:null"`
	ContactEmail             *string    `json:"contact_email" gorm:"type:nvarchar(255);default:null"`
	CreatedBYID              *uint      `json:"created_by_id" gorm:"default:null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	Process        *Process                             `gorm:"foreignKey:ProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	EntityType     *ProcessExternalDependencyEntityType `gorm:"foreignKey:EntityTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	DependencyTier *ProcessDependencyTier               `gorm:"foreignKey:DependencyTierID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	VendorRTOUnit  *BusinessContinuityDurationUnit      `gorm:"foreignKey:VendorRTOUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser  *User                                `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser  *User                                `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProcessExternalDependency) TableName() string {
	return "process_external_dependencies"
}
