package models

import "time"

type StrategicInitiativeStatusUpdate struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	StrategyInitiativeID        uint `json:"strategy_initiative_id" gorm:"not null"`
	StrategicInitiativeStatusID uint `json:"strategic_initiative_status_id" gorm:"not null"`

	Comment    string `json:"comment" gorm:"type:nvarchar(max);default:null"`
	Percentage uint   `json:"percentage" binding:"min=0,max=100" gorm:"not null;default:0"`

	CreatedAt time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	CreatedByID uint  `json:"created_by_id" gorm:"not null"`
	UpdatedByID *uint `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	StrategyInitiative        *StrategicIntiative        `json:"strategic_intiatives,omitempty" gorm:"foreignKey:StrategyInitiativeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	StrategicInitiativeStatus *StrategicInitiativeStatus `json:"strategic_initiative_statuses,omitempty" gorm:"foreignKey:StrategicInitiativeStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`

	CreatedBytime *User `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedByID"`
	UpdatedByUser *User `json:"updated_by_user,omitempty" gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides table name
func (StrategicInitiativeStatusUpdate) TableName() string {
	return "strategic_initiative_status_updates"
}
