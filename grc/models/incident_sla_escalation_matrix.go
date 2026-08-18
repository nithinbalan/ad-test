package models

type IncidentSLAEscalationMatrix struct {
	ID                         uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	SeverityID                 uint   `json:"severity_id" gorm:"not null" binding:"required"`
	Type                       string `json:"type" binding:"required,oneof=escalation_matrix resolution_matrix" gorm:"type:varchar(20);not null"`
	NotifyWithin               *int   `json:"notify_within,omitempty"`
	NotifyWithinUnitID         *uint  `json:"notify_within_unit_id,omitempty"`
	NotifyToIncidentOwner      bool   `json:"notify_to_incident_owner" gorm:"not null;default:false"`
	NotifyToResponseTeam       bool   `json:"notify_to_response_team" gorm:"not null;default:false"`
	NotifyRoleID               *uint  `json:"notify_role_id,omitempty"`
	EscalateAfterIncidentOwner bool   `json:"escalate_after_incident_owner" gorm:"not null;default:false"`
	EscalateAfterResponseTeam  bool   `json:"escalate_after_response_team" gorm:"not null;default:false"`
	EscalateRoleID             *uint  `json:"escalate_role_id,omitempty"`
	CloseWithin                *int   `json:"close_within,omitempty"`
	CloseWithinUnitID          *uint  `json:"close_within_unit_id,omitempty"`
	CloseWithinIncidentOwner   bool   `json:"close_within_incident_owner" gorm:"not null;default:false"`
	CloseWithinResponseTeam    bool   `json:"close_within_response_team" gorm:"not null;default:false"`
	CloseWithinRoleID          *uint  `json:"close_within_role_id,omitempty"`
	Status                     string `json:"status" gorm:"type:varchar(10);not null;default:active"`

	// Relationships
	SeverityLevel    *SeverityLevel                  `json:"severity_level,omitempty" gorm:"foreignKey:SeverityID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	NotifyWithinUnit *BusinessContinuityDurationUnit `json:"notify_within_unit,omitempty" gorm:"foreignKey:NotifyWithinUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	NotifyRole       *Role                           `json:"notify_role,omitempty" gorm:"foreignKey:NotifyRoleID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	EscalateRole     *Role                           `json:"escalate_role,omitempty" gorm:"foreignKey:EscalateRoleID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CloseWithinUnit  *BusinessContinuityDurationUnit `json:"close_within_unit,omitempty" gorm:"foreignKey:CloseWithinUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CloseWithinRole  *Role                           `json:"close_within_role,omitempty" gorm:"foreignKey:CloseWithinRoleID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (IncidentSLAEscalationMatrix) TableName() string {
	return "incident_sla_escalation_matrix"
}
