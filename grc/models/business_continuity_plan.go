package models

import (
	"time"
)

type BusinessContinuityPlan struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber    string     `json:"reference_number" binding:"omitempty,max=45" gorm:"size:45"`
	BusinessFunctionID uint       `json:"business_function_id" binding:"required"`
	Version            *int       `json:"version"`
	IsLatestVersion    *bool      `json:"is_latest_version" gorm:"default:true"`
	Purpose            *string    `json:"purpose"`
	Scope              *string    `json:"scope"`
	Objectives         *string    `json:"objectives"`
	ResponsibleUserID  uint       `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ResponsibleUnitID  *uint      `json:"responsible_unit_id" gorm:"not null" binding:"required"`
	BcpStatusID        uint       `json:"bcp_status_id" binding:"required" gorm:"not null"`
	Rto                *int       `json:"rto"`
	RtoUnitId          *uint      `json:"rto_unit_id"`
	Mtpd               *int       `json:"mtpd"`
	MtpdUnitId         *uint      `json:"mtpd_unit_id"`
	Rpo                *int       `json:"rpo"`
	RpoUnitId          *uint      `json:"rpo_unit_id"`
	Assumptions        *string    `json:"assumptions"`
	CreatedBYID        uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID        *uint      `json:"updated_by_id"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	BusinessFunction      *BusinessFunction      `gorm:"foreignKey:BusinessFunctionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Processes             []BcpProcess           `gorm:"foreignKey:BusinessContinuityPlanID"`
	BcpStatus             *BcpStatus             `gorm:"foreignKey:BcpStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser       *User                  `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit       *OrganizationUnit      `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RtoUnit               *TimeUnit              `gorm:"foreignKey:RtoUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	MtpdUnit              *TimeUnit              `gorm:"foreignKey:MtpdUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RpoUnit               *TimeUnit              `gorm:"foreignKey:RpoUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Documents             []BcpDocument          `gorm:"foreignKey:BusinessContinuityPlanID"`
	BcpRecoveryPlans      []BcpRecoveryPlan      `gorm:"foreignKey:BusinessContinuityPlanID"`
	BcpCommunicationPlans []BcpCommunicationPlan `gorm:"foreignKey:BusinessContinuityPlanID"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser         *User                  `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by BusinessContinuityPlan to `business_continuity_plans`
func (BusinessContinuityPlan) TableName() string {
	return "business_continuity_plans"
}
