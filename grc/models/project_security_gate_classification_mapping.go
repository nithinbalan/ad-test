package models

type ProjectSecurityGateClassificationMapping struct {
	ID               uint `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	SecurityGateID   uint `json:"security_gate_id" gorm:"not null"`
	ClassificationID uint `json:"classification_id" gorm:"not null"`

	// Relationships
	SecurityGate   *ProjectSecurityGate               `json:"security_gate" gorm:"foreignKey:SecurityGateID"`
	Classification *ProjectSecurityGateClassification `json:"classification" gorm:"foreignKey:ClassificationID"`
}

func (ProjectSecurityGateClassificationMapping) TableName() string {
	return "project_security_gate_classification_mappings"
}
