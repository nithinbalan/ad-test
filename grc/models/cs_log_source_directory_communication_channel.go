package models

type CSLogSourceDirectoryCommunicationChannel struct {
	ID                     uint `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	SourceDirectoryID      uint `json:"source_directory_id" gorm:"not null"`
	CommunicationChannelID uint `json:"communication_channel_id" gorm:"not null"`

	// Relationships
	SourceDirectory      *CSLogSourceDirectory      `gorm:"foreignKey:SourceDirectoryID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	CommunicationChannel *CSLogCommunicationChannel `gorm:"foreignKey:CommunicationChannelID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (CSLogSourceDirectoryCommunicationChannel) TableName() string {
	return "cs_log_source_directory_communication_channels"
}
