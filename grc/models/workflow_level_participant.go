package models

type WorkflowLevelParticipant struct {
	ID              uint  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	WorkflowLevelID uint  `json:"workflow_level_id" binding:"required" gorm:"not null"`
	UserID          *uint `json:"user_id" gorm:"default:null"`
	RoleID          *uint `json:"role_id" gorm:"default:null"`

	// Relationships
	WorkflowLevel *WorkflowLevel `gorm:"foreignKey:WorkflowLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User          *User          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Role          *Role          `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by GORM to `workflow_level_participants`
func (WorkflowLevelParticipant) TableName() string {
	return "workflow_level_participants"
}
