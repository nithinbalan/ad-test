package models

import "time"

type RiskUpdatePlanErmDeptDocument struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskUpdatePlanTimelineID uint       `json:"risk_update_plan_timeline_id" binding:"required" gorm:"not null"`
	URL                      *string    `json:"url" gorm:"type:varchar(1000);default:null"`
	Name                     *string    `json:"name" gorm:"type:nvarchar(255);default:null"`
	Size                     *int       `json:"size" gorm:"default:null"`
	Token                    *string    `json:"token" gorm:"type:nvarchar(max);default:null"`
	OriginalName             *string    `json:"original_name" gorm:"type:nvarchar(255);default:null"`
	Type                     *string    `json:"type" gorm:"type:varchar(255);default:null"`
	CreatedBYID              *uint      `json:"created_by_id" gorm:"default:null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	RiskUpdatePlanTimeline *RiskUpdatePlanTimeline `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by GORM
func (RiskUpdatePlanErmDeptDocument) TableName() string {
	return "risk_update_plan_erm_dept_documents"
}
