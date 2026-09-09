package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
)

type CSLogSourceDirectoryListDTO struct {
	ID                          uint    `json:"id"`
	ReferenceNumber             string  `json:"reference_number"`
	EntityName                  string  `json:"entity_name"`
	DataEntityTypeID            *uint   `json:"data_entity_type_id"`
	DataEntityTypeTitleEn       *string `json:"data_entity_type_title_en"`
	DataEntityTypeTitleAr       *string `json:"data_entity_type_title_ar"`
	ResponsibleEntityID         *uint   `json:"responsible_entity_id"`
	ResponsibleEntityTitleEn    *string `json:"responsible_entity_title_en"`
	ResponsibleEntityTitleAr    *string `json:"responsible_entity_title_ar"`
	BeneficiaryEntityID         *uint   `json:"beneficiary_entity_id"`
	BeneficiaryEntityTitleEn    *string `json:"beneficiary_entity_title_en"`
	BeneficiaryEntityTitleAr    *string `json:"beneficiary_entity_title_ar"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   *string `json:"updated_at"`
	CreatedUserID               uint    `json:"created_user_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
}

type CSLogSourceDirectoryListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
	Export   bool
}

type CSLogSourceDirectoryListResponseDTO struct {
	Items      []*CSLogSourceDirectoryListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO   `json:"pagination"`
}

type CSLogSourceDirectoryReceptionMechanismDTO struct {
	ID                   *uint                          `json:"id"`
	SourceDirectoryID    uint                           `json:"source_directory_id"`
	ReceptionMechanismID uint                           `json:"reception_mechanism_id"`
	ReceptionMechanism   CSLogDataReceptionMechanismDTO `json:"reception_mechanism"`
}

type CSLogSourceDirectoryCommunicationChannelDTO struct {
	ID                     *uint                        `json:"id"`
	SourceDirectoryID      uint                         `json:"source_directory_id"`
	CommunicationChannelID uint                         `json:"communication_channel_id"`
	CommunicationChannel   CSLogCommunicationChannelDTO `json:"communication_channel"`
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
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`

	DataEntityType    *CSLogDataEntityTypeDTO           `json:"data_entity_type"`
	ResponsibleEntity *settingsDTOs.OrganizationUnitDTO `json:"responsible_entity"`
	BeneficiaryEntity *settingsDTOs.OrganizationUnitDTO `json:"beneficiary_entity"`
	CreatedByUser     organizationDTOs.UserMinimalDTO   `json:"created_by"`
	UpdatedByUser     *organizationDTOs.UserMinimalDTO  `json:"updated_by"`

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
	ReceptionMechanismIDs   []*uint `json:"reception_mechanism_ids" binding:"required,min=1"`
	CommunicationChannelIDs []*uint `json:"communication_channel_ids" binding:"required,min=1"`
}
