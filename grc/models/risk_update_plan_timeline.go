package models

import "time"

type RiskUpdatePlanTimeline struct {
	ID                           uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskUpdatePlanID             uint       `json:"risk_update_plan_id" binding:"required" gorm:"not null"`
	RiskUpdatePlanTimelineTypeID uint       `json:"risk_update_plan_timeline_type_id" binding:"required" gorm:"not null"`
	IsBusinessContextUpdated     bool       `json:"is_business_context_updated" gorm:"default:false;not null"`
	IsNotApplicable              bool       `json:"is_not_applicable" gorm:"default:false;not null"`
	IsCompleted                  bool       `json:"is_completed" gorm:"default:false;not null"`
	Justification                *string    `json:"justification" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID                  uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID                  *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                    time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                    *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	RiskUpdatePlan                  *RiskRegisterUpdatePlan                         `gorm:"foreignKey:RiskUpdatePlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskUpdatePlanTimelineType      *RiskUpdatePlanTimelineType                     `gorm:"foreignKey:RiskUpdatePlanTimelineTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ErmDeptDocuments                []*RiskUpdatePlanErmDeptDocument                `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IntroductoryMeetingDocuments    []*RiskUpdatePlanIntroductoryMeetingDocument    `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	BusinessContexts                []*RiskUpdatePlanBusinessContext                `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Risks                           []*RiskUpdatePlanRisk                           `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Controls                        []*RiskUpdatePlanControl                        `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Processes                       []*RiskUpdatePlanProcess                        `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	LeadershipAssessments           []*RiskUpdatePlanLeadershipAssessment           `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	LeadershipDocuments             []*RiskUpdatePlanLeadershipDocument             `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskAssessmentWorkshopDocuments []*RiskUpdatePlanRiskAssessmentWorkshopDocument `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskTeamAssessments             []*RiskUpdatePlanRiskTeamAssessment             `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskTeamDocuments               []*RiskUpdatePlanRiskTeamDocument               `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskWorkshopAssessmentDocuments []*RiskUpdatePlanRiskWorkshopAssessmentDocument `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SelfAwarenessTrainings          []*RiskUpdatePlanSelfAwarenessTraining          `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser                   *User                                           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser                   *User                                           `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by RiskUpdatePlanTimeline to `risk_update_plan_timelines`
func (RiskUpdatePlanTimeline) TableName() string {
	return "risk_update_plan_timelines"
}
