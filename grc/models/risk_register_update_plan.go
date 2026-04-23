package models

import "time"

type RiskRegisterUpdatePlan struct {
	ID                             uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DivisionID                     uint       `json:"division_id" binding:"required"`
	BusinessUnitID                 uint       `json:"business_unit_id" binding:"required"`
	RiskRegisterUpdatePlanStatusID uint       `json:"risk_register_update_plan_status_id" binding:"required"`
	NumberOfOperatingEntities      *int       `json:"number_of_operating_entities" gorm:"default:null"`
	Workdays                       *int       `json:"workdays" gorm:"default:null"`
	RiskChampionID                 uint       `json:"risk_champion_id" binding:"required"`
	StartDate                      time.Time  `json:"start_date" binding:"required" gorm:"type:datetime;not null"`
	EndDate                        time.Time  `json:"end_date" binding:"required" gorm:"type:datetime;not null"`
	Notes                          *string    `json:"notes" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID                    uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID                    *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                      time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                      *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	SubmittedBYID                  *uint      `json:"submitted_by_id" gorm:"default:null"`
	ApprovedBYID                   *uint      `json:"approved_by_id" gorm:"default:null"`
	SubmittedAt                    *time.Time `json:"submitted_at" gorm:"type:datetime;default:null"`
	ApprovedAt                     *time.Time `json:"approved_at" gorm:"type:datetime;default:null"`
	NextReviewLevel                *int       `json:"next_review_level" gorm:"default:null"`
	WorkflowStatusID               *uint      `json:"workflow_status_id" gorm:"default:null"`

	// Relationships
	Division                     *OrganizationUnit                   `gorm:"foreignKey:DivisionID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	BusinessUnit                 *OrganizationUnit                   `gorm:"foreignKey:BusinessUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskRegisterUpdatePlanStatus *RiskRegisterUpdatePlanStatus       `gorm:"foreignKey:RiskRegisterUpdatePlanStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskChampion                 *User                               `gorm:"foreignKey:RiskChampionID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	TeamMembers                  []*RiskRegisterUpdatePlanTeamMember `gorm:"foreignKey:RiskRegisterUpdatePlanID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	CreatedByUser                *User                               `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser                *User                               `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SubmittedByUser              *User                               `gorm:"foreignKey:SubmittedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ApprovedByUser               *User                               `gorm:"foreignKey:ApprovedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus               *WorkflowStatus                     `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by GORM
func (RiskRegisterUpdatePlan) TableName() string {
	return "risk_register_update_plans"
}
