package models

type WorkflowLevelParticipant struct {
	ID                   uint  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	WorkflowLevelID      uint  `json:"workflow_level_id" binding:"required" gorm:"not null"`
	UserID               *uint `json:"user_id" gorm:"default:null"`
	RoleID               *uint `json:"role_id" gorm:"default:null"`
	WorkflowCustomRoleID *uint `json:"workflow_custom_role_id" gorm:"column:workflow_custom_role_id;default:null"`
	UnitID               *uint `json:"unit_id" gorm:"column:unit_id;default:null"`

	// Relationships
	WorkflowLevel      *WorkflowLevel      `gorm:"foreignKey:WorkflowLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User               *User               `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Role               *Role               `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowCustomRole *WorkflowCustomRole `gorm:"foreignKey:WorkflowCustomRoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Unit *OrganizationUnit `gorm:"foreignKey:UnitID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (WorkflowLevelParticipant) TableName() string {
	return "workflow_level_participants"
}
