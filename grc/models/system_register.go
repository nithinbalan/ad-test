package models

import (
	"time"
)

type SystemRegister struct {
	ID                      uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn                 string     `json:"title_en" gorm:"type:nvarchar(500);not null"`
	TitleAr                 string     `json:"title_ar" gorm:"type:nvarchar(500);not null"`
	ServiceDescription      *string    `json:"service_description" gorm:"type:nvarchar(max)"`
	ModulesComponents       *string    `json:"modules_components" gorm:"type:nvarchar(max)"`
	ThirdPartyDependencies  *string    `json:"third_party_dependencies" gorm:"type:nvarchar(max)"`
	ServiceImportanceID     *uint      `json:"service_importance_id" gorm:"type:int"`
	ResponsibleITTeamID     *uint      `json:"responsible_it_team_id" gorm:"type:int"`
	OutsourcingPartner      *string    `json:"outsourcing_partner" gorm:"type:nvarchar(max)"`
	OperationLocationID     *uint      `json:"operation_location_id" gorm:"type:int"`
	BusinessOwnerID         *uint      `json:"business_owner_id" gorm:"type:int"`
	Customers               *string    `json:"customers" gorm:"type:nvarchar(max)"`
	DisruptionImpactLevelID *uint      `json:"disruption_impact_level_id" gorm:"type:int"`
	RTO                     *int       `json:"rto" gorm:"type:int"`
	RTOUnitID               *uint      `json:"rto_unit_id" gorm:"type:int"`
	MTPD                    *int       `json:"mtpd" gorm:"type:int"`
	MTPDUnitID              *uint      `json:"mtpd_unit_id" gorm:"type:int"`
	RPO                     *int       `json:"rpo" gorm:"type:int"`
	RPOUnitID               *uint      `json:"rpo_unit_id" gorm:"type:int"`
	ITProcesses             *string    `json:"it_processes" gorm:"type:nvarchar(max)"`
	ITServiceOwnerID        *uint      `json:"it_service_owner_id" gorm:"type:int"`
	ProcessOperationTeamID  uint       `json:"process_operation_team_id" gorm:"not null"`
	ProcessRecoveryTeamID   uint       `json:"process_recovery_team_id" gorm:"not null"`
	RecoveryManagerID       *uint      `json:"recovery_manager_id" gorm:"type:int"`
	RecoveryManagerPhone    string     `json:"recovery_manager_phone" gorm:"type:varchar(20);not null"`
	RecoveryManagerEmail    string     `json:"recovery_manager_email" gorm:"type:varchar(255);not null"`
	RecoveryLocation        *string    `json:"recovery_location" gorm:"type:nvarchar(max)"`
	VendorRTO               *int       `json:"vendor_rto" gorm:"type:int"`
	VendorRTOUnitID         *uint      `json:"vendor_rto_unit_id" gorm:"type:int"`
	CreatedBYID             uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID             *uint      `json:"updated_by_id"`
	CreatedAt               time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt               *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	ServiceImportance         *SystemRegisterServiceImportance `gorm:"foreignKey:ServiceImportanceID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ResponsibleITTeam         *JobTeam                         `gorm:"foreignKey:ResponsibleITTeamID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	OperationLocation         *Location                        `gorm:"foreignKey:OperationLocationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	BusinessOwner             *OrganizationUnit                `gorm:"foreignKey:BusinessOwnerID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	DisruptionImpactLevel     *SystemRegisterDisruptionLevel   `gorm:"foreignKey:DisruptionImpactLevelID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RTOUnit                   *BusinessContinuityDurationUnit  `gorm:"foreignKey:RTOUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	MTPDUnit                  *BusinessContinuityDurationUnit  `gorm:"foreignKey:MTPDUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RPOUnit                   *BusinessContinuityDurationUnit  `gorm:"foreignKey:RPOUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ITServiceOwner            *Designation                     `gorm:"foreignKey:ITServiceOwnerID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ProcessOperationTeam      *JobTeam                         `gorm:"foreignKey:ProcessOperationTeamID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ProcessRecoveryTeam       *JobTeam                         `gorm:"foreignKey:ProcessRecoveryTeamID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RecoveryManager           *Designation                     `gorm:"foreignKey:RecoveryManagerID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	VendorRTOUnit             *BusinessContinuityDurationUnit  `gorm:"foreignKey:VendorRTOUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser             *User                            `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser             *User                            `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SystemRegisterSLA         *SystemRegisterSLA               `gorm:"foreignKey:SystemRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SystemRegisterImpactTypes []*SystemRegisterImpactType      `gorm:"foreignKey:SystemRegisterID"`
}

func (SystemRegister) TableName() string {
	return "system_registers"
}
