package dtos

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
