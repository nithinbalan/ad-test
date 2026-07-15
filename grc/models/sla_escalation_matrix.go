package models

type SLAEscalationMatrix struct {
	ID               uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	SeverityID       uint   `json:"severity_id" gorm:"not null" binding:"required"`
	RemediationDays  int    `json:"remediation_days" gorm:"not null" binding:"required"`
	IsEmail          bool   `json:"is_email" gorm:"not null;default:false"`
	IsAlert          bool   `json:"is_alert" gorm:"not null;default:false"`
	RoleID           uint   `json:"role_id" gorm:"not null" binding:"required"`
	NotifiyingRoleID *uint  `json:"notifiying_role_id,omitempty"`
	Status           string `json:"status" gorm:"type:varchar(10);not null;default:active"`

	// Relationships
	SeverityLevel  *SeverityLevel `json:"severity_level,omitempty" gorm:"foreignKey:SeverityID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Role           *Role          `json:"role,omitempty" gorm:"foreignKey:RoleID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	NotifiyingRole *Role          `json:"notifiying_role,omitempty" gorm:"foreignKey:NotifiyingRoleID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (SLAEscalationMatrix) TableName() string {
	return "sla_escalation_matrix"
}
