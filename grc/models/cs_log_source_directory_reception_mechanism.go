package models

type CSLogSourceDirectoryReceptionMechanism struct {
	ID                   uint `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	SourceDirectoryID    uint `json:"source_directory_id" gorm:"not null"`
	ReceptionMechanismID uint `json:"reception_mechanism_id" gorm:"not null"`

	// Relationships
	SourceDirectory    *CSLogSourceDirectory        `gorm:"foreignKey:SourceDirectoryID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE;"`
	ReceptionMechanism *CSLogDataReceptionMechanism `gorm:"foreignKey:ReceptionMechanismID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
}

func (CSLogSourceDirectoryReceptionMechanism) TableName() string {
	return "cs_log_source_directory_reception_mechanisms"
}
