package models

import "time"

type AuditFindingStrategicObjective struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	StrategicObjectiveID uint      `json:"strategic_objective_id" gorm:"not null"`
	AuditFindingID       uint      `json:"audit_finding_id" gorm:"not null"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`
	CreatedBYID          uint      `json:"created_by_id" gorm:"not null"`

	// Relationships
	StrategicObjective *StrategicObjective `gorm:"foreignKey:StrategicObjectiveID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AuditFinding       *AuditFinding       `gorm:"foreignKey:AuditFindingID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBYID"`
}

// TableName overrides the default table name
func (AuditFindingStrategicObjective) TableName() string {
	return "audit_finding_strategic_objectives"
}
