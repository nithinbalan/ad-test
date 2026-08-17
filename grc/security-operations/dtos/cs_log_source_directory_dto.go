package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type CSLogSourceDirectoryListDTO struct {
	ID                       uint       `json:"id"`
	ReferenceNumber          string     `json:"reference_number"`
	EntityName               string     `json:"entity_name"`
	DataEntityTypeID         *uint      `json:"data_entity_type_id"`
	DataEntityTypeTitleEn    *string    `json:"data_entity_type_title_en"`
	DataEntityTypeTitleAr    *string    `json:"data_entity_type_title_ar"`
	ResponsibleEntityID      *uint      `json:"responsible_entity_id"`
	ResponsibleEntityTitleEn *string    `json:"responsible_entity_title_en"`
	ResponsibleEntityTitleAr *string    `json:"responsible_entity_title_ar"`
	BeneficiaryEntityID      *uint      `json:"beneficiary_entity_id"`
	BeneficiaryEntityTitleEn *string    `json:"beneficiary_entity_title_en"`
	BeneficiaryEntityTitleAr *string    `json:"beneficiary_entity_title_ar"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
	CreatedByID              uint       `json:"created_by_id"`
	CreatedByFirstNameEn     string     `json:"created_by_first_name_en"`
	CreatedByFirstNameAr     string     `json:"created_by_first_name_ar"`
	CreatedByLastNameEn      string     `json:"created_by_last_name_en"`
	CreatedByLastNameAr      string     `json:"created_by_last_name_ar"`
	CreatedByEmail           string     `json:"created_by_email"`
	UpdatedByID              *uint      `json:"updated_by_id"`
	UpdatedByFirstNameEn     string     `json:"updated_by_first_name_en"`
	UpdatedByFirstNameAr     string     `json:"updated_by_first_name_ar"`
	UpdatedByLastNameEn      string     `json:"updated_by_last_name_en"`
	UpdatedByLastNameAr      string     `json:"updated_by_last_name_ar"`
	UpdatedByEmail           string     `json:"updated_by_email"`
}

type CSLogSourceDirectoryListQueryParamsDTO struct {
	Sort                string
	SortBy              string
	Keyword             string
	Page                int
	PageSize            int
	DataEntityTypeID    *uint
	ResponsibleEntityID *uint
	BeneficiaryEntityID *uint
}

type CSLogSourceDirectoryListResponseDTO struct {
	Items      []*CSLogSourceDirectoryListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO   `json:"pagination"`
}

type CSLogSourceDirectoryLookupDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

type CSLogSourceDirectoryOrganizationUnitDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

type CSLogSourceDirectoryUserDTO struct {
	ID          uint   `json:"id"`
	FirstNameEn string `json:"first_name_en"`
	FirstNameAr string `json:"first_name_ar"`
	LastNameEn  string `json:"last_name_en"`
	LastNameAr  string `json:"last_name_ar"`
	Email       string `json:"email"`
}

type CSLogSourceDirectoryReceptionMechanismDTO struct {
	ID                   uint                           `json:"id"`
	SourceDirectoryID    uint                           `json:"source_directory_id"`
	ReceptionMechanismID uint                           `json:"reception_mechanism_id"`
	ReceptionMechanism   *CSLogSourceDirectoryLookupDTO `json:"reception_mechanism"`
}

type CSLogSourceDirectoryCommunicationChannelDTO struct {
	ID                     uint                           `json:"id"`
	SourceDirectoryID      uint                           `json:"source_directory_id"`
	CommunicationChannelID uint                           `json:"communication_channel_id"`
	CommunicationChannel   *CSLogSourceDirectoryLookupDTO `json:"communication_channel"`
}

type CSLogSourceDirectoryDetailDTO struct {
	ID                  uint       `json:"id"`
	ReferenceNumber     string     `json:"reference_number"`
	EntityName          string     `json:"entity_name"`
	DataEntityTypeID    *uint      `json:"data_entity_type_id"`
	Location            *string    `json:"location"`
	ResponsibleEntityID *uint      `json:"responsible_entity_id"`
	BeneficiaryEntityID *uint      `json:"beneficiary_entity_id"`
	ContactName         *string    `json:"contact_name"`
	ContactMobile       *string    `json:"contact_mobile"`
	ContactEmail        *string    `json:"contact_email"`
	Website             *string    `json:"website"`
	Notes               *string    `json:"notes"`
	CreatedByID         uint       `json:"created_by_id"`
	UpdatedByID         *uint      `json:"updated_by_id"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`

	DataEntityType    *CSLogSourceDirectoryLookupDTO           `json:"data_entity_type"`
	ResponsibleEntity *CSLogSourceDirectoryOrganizationUnitDTO `json:"responsible_entity"`
	BeneficiaryEntity *CSLogSourceDirectoryOrganizationUnitDTO `json:"beneficiary_entity"`
	CreatedBy         *CSLogSourceDirectoryUserDTO             `json:"created_by"`
	UpdatedBy         *CSLogSourceDirectoryUserDTO             `json:"updated_by"`

	ReceptionMechanisms   []*CSLogSourceDirectoryReceptionMechanismDTO   `json:"reception_mechanisms"`
	CommunicationChannels []*CSLogSourceDirectoryCommunicationChannelDTO `json:"communication_channels"`
}

type CSLogSourceDirectoryFormDTO struct {
	EntityName              string  `json:"entity_name" binding:"required,max=255"`
	DataEntityTypeID        *uint   `json:"data_entity_type_id"`
	Location                *string `json:"location"`
	ResponsibleEntityID     *uint   `json:"responsible_entity_id"`
	BeneficiaryEntityID     *uint   `json:"beneficiary_entity_id"`
	ContactName             *string `json:"contact_name"`
	ContactMobile           *string `json:"contact_mobile"`
	ContactEmail            *string `json:"contact_email"`
	Website                 *string `json:"website"`
	Notes                   *string `json:"notes"`
	ReceptionMechanismIDs   []uint  `json:"reception_mechanism_ids" binding:"required,min=1"`
	CommunicationChannelIDs []uint  `json:"communication_channel_ids" binding:"required,min=1"`
}
