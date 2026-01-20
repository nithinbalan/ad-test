package models

type ActionPlanMaster struct {
	ID                 uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	ModuleID           uint    `json:"module_id" binding:"required" gorm:"not null"`
	IncidentCategoryID *uint   `json:"incident_category_id" gorm:"default:null"`
	SeverityLevelID    *uint   `json:"severity_level_id" gorm:"default:null"`
	Title              string  `json:"title" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Description        *string `json:"description" gorm:"type:nvarchar(max)"`
	Status             string  `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active'"`

	// Relationships (reflect DB foreign keys)
	Module           *Module           `gorm:"foreignKey:ModuleID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentCategory *IncidentCategory `gorm:"foreignKey:IncidentCategoryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SeverityLevel    *SeverityLevel    `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides default table name
func (ActionPlanMaster) TableName() string {
	return "action_plan_master"
}
