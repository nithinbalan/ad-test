package models

import (
	"time"
)

type SystemRegister struct {
	ID                            uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ReferenceNumber               string     `json:"reference_number" gorm:"type:varchar(100);uniqueIndex;not null"`
	TitleEn                       string     `json:"title_en" gorm:"type:nvarchar(1000);not null"`
	TitleAr                       string     `json:"title_ar" gorm:"type:nvarchar(1000);not null"`
	NumberOfPeopleNecessary       *int       `json:"number_of_people_necessary" gorm:"type:int"`
	BusinessOwnerUnitID           *uint      `json:"business_owner_unit_id" gorm:"type:int"`
	ITResponsibleUnitID           *uint      `json:"it_responsible_unit_id" gorm:"type:int"`
	ThirdPartyDependencies        *string    `json:"third_party_dependencies" gorm:"type:nvarchar(max)"`
	OutsourcingPartners           *string    `json:"outsourcing_partners" gorm:"type:nvarchar(max)"`
	CustomerUser                  *string    `json:"customer_user" gorm:"type:nvarchar(max)"`
	ProcessOperationTeamID        *uint      `json:"process_operation_team_id" gorm:"type:int"`
	ProcessRecoveryTeamID         *uint      `json:"process_recovery_team_id" gorm:"type:int"`
	RecoveryManagerID             *uint      `json:"recovery_manager_id" gorm:"type:int"`
	RecoveryManagerContactDetails *string    `json:"recovery_manager_contact_details" gorm:"type:nvarchar(max)"`
	CreatedBYID                   uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID                   *uint      `json:"updated_by_id"`
	CreatedAt                     time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                     *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	BusinessOwnerUnit                  *OrganizationUnit                    `gorm:"foreignKey:BusinessOwnerUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ITResponsibleUnit                  *OrganizationUnit                    `gorm:"foreignKey:ITResponsibleUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ProcessOperationTeam               *JobTeam                             `gorm:"foreignKey:ProcessOperationTeamID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ProcessRecoveryTeam                *JobTeam                             `gorm:"foreignKey:ProcessRecoveryTeamID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RecoveryManager                    *User                                `gorm:"foreignKey:RecoveryManagerID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser                      *User                                `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser                      *User                                `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SystemRegisterJobTitles            []*SystemRegisterJobTitle            `gorm:"foreignKey:SystemRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SystemRegisterOperationalLocations []*SystemRegisterOperationalLocation `gorm:"foreignKey:SystemRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SystemRegisterRecoveryLocations    []*SystemRegisterRecoveryLocation    `gorm:"foreignKey:SystemRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SystemRegisterRecoveryNeededTools  []*SystemRegisterRecoveryNeededTool  `gorm:"foreignKey:SystemRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (SystemRegister) TableName() string {
	return "system_registers"
}
