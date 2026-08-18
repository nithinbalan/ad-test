package models

import (
	"time"
)

type CSLogSourceDirectory struct {
	ID                  uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber     string  `json:"reference_number" gorm:"type:nvarchar(45);not null;unique"`
	EntityName          string  `json:"entity_name" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	DataEntityTypeID    *uint   `json:"data_entity_type_id" gorm:"column:data_entity_type_id"`
	Location            *string `json:"location" gorm:"type:nvarchar(500)"`
	ResponsibleEntityID *uint   `json:"responsible_entity_id" gorm:"column:responsible_entity_id"`
	BeneficiaryEntityID *uint   `json:"beneficiary_entity_id" gorm:"column:beneficiary_entity_id"`
	ContactName         *string `json:"contact_name" gorm:"type:nvarchar(255)"`
	ContactMobile       *string `json:"contact_mobile" gorm:"type:varchar(50)"`
	ContactEmail        *string `json:"contact_email" gorm:"type:varchar(255)"`
	Website             *string `json:"website" gorm:"type:nvarchar(500)"`
	Notes               *string `json:"notes" gorm:"type:nvarchar(max)"`

	CreatedBYID uint       `json:"created_by_id" gorm:"column:created_by_id;not null"`
	UpdatedBYID *uint      `json:"updated_by_id" gorm:"column:updated_by_id"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	DataEntityType    *CSLogDataEntityType `gorm:"foreignKey:DataEntityTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ResponsibleEntity *OrganizationUnit    `gorm:"foreignKey:ResponsibleEntityID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	BeneficiaryEntity *OrganizationUnit    `gorm:"foreignKey:BeneficiaryEntityID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser     *User                `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser     *User                `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`

	ReceptionMechanisms   []*CSLogSourceDirectoryReceptionMechanism   `gorm:"foreignKey:SourceDirectoryID"`
	CommunicationChannels []*CSLogSourceDirectoryCommunicationChannel `gorm:"foreignKey:SourceDirectoryID"`
}

func (CSLogSourceDirectory) TableName() string {
	return "cs_log_source_directory"
}
