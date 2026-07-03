package models

import (
	"time"
)

type BcpProcessInternalDependency struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BcpProcessID             uint       `json:"bcp_process_id" gorm:"not null"`
	SupportingUnitID         uint       `json:"supporting_unit_id" gorm:"not null"`
	DependencyTierID         uint       `json:"dependency_tier_id" gorm:"not null"`
	DependencyPercentageFrom *int       `json:"dependency_percentage_from" gorm:"default:null"`
	DependencyPercentageTo   *int       `json:"dependency_percentage_to" gorm:"default:null"`
	RequiredSupport          *string    `json:"required_support" gorm:"type:nvarchar(MAX);default:null"`
	Notes                    *string    `json:"notes" gorm:"type:nvarchar(MAX);default:null"`
	CreatedBYID              *uint      `json:"created_by_id" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	BcpProcess     *BcpProcess            `gorm:"foreignKey:BcpProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SupportingUnit *OrganizationUnit      `gorm:"foreignKey:SupportingUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	DependencyTier *ProcessDependencyTier `gorm:"foreignKey:DependencyTierID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser  *User                  `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BcpProcessInternalDependency) TableName() string {
	return "bcp_process_internal_dependencies"
}
