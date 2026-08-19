package models

import (
	"time"
)

type IncidentRegister struct {
	ID                                  uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber                     string     `json:"reference_number" binding:"omitempty,max=45" gorm:"size:45"`
	IncidentCategoryID                  uint       `json:"incident_category_id" binding:"required"`
	IncidentTypeID                      *uint      `json:"incident_type_id"`
	SeverityLevelID                     uint       `json:"severity_level_id" binding:"required"`
	IncidentStatusID                    *uint      `json:"incident_status_id"`
	ResponsibleUnitID                   uint       `json:"responsible_unit_id" binding:"required" gorm:"not null"`
	ResponsibleUserID                   uint       `json:"responsible_user_id" binding:"required" gorm:"not null"`
	BusinessOwnerID                     *uint      `json:"business_owner_id"`
	IncidentResponseTeamID              *uint      `json:"incident_response_team_id"`
	ResponseTeamLeadID                  *uint      `json:"response_team_lead_id"`
	IncidentScopeID                     *uint      `json:"incident_scope_id"`
	IncidentSensitivityClassificationID *uint      `json:"incident_sensitivity_classification_id"`
	ConfidentialitySeverityLevelID      *uint      `json:"confidentiality_severity_level_id"`
	IntegritySeverityLevelID            *uint      `json:"integrity_severity_level_id"`
	AvailabilitySeverityLevelID         *uint      `json:"availability_severity_level_id"`
	Title                               string     `json:"title" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	Description                         *string    `json:"description" gorm:"type:nvarchar(max)"`
	IncidentAt                          *time.Time `json:"incident_at" binding:"required" gorm:"type:date;not null"`
	Location                            *string    `json:"location" gorm:"type:varchar(100)"`
	ReportedBy                          uint       `json:"reported_by" binding:"required"`
	ReportedAt                          *time.Time `json:"reported_at" gorm:"type:date"`
	QuickAction                         *string    `json:"quick_action" gorm:"type:nvarchar(max)"`
	MarkedAsAuditableItem               bool       `json:"marked_as_auditable_item" gorm:"default:false" binding:"-"`
	IsAcknowledge                       bool       `json:"is_acknowledge" gorm:"default:false" binding:"-"`
	AcknowledgedAt                      *time.Time `json:"acknowledged_at" gorm:"type:datetime"`
	AcknowledgedByID                    *uint      `json:"acknowledged_by_id"`
	SchedulerLastRunAt                  *time.Time `json:"scheduler_last_run_at" gorm:"type:datetime"`
	IsSLATimeAchieved                   bool       `json:"is_sla_time_achieved" gorm:"default:false" binding:"-"`
	IsSLATimeExceeded                   bool       `json:"is_sla_time_exceeded" gorm:"default:false" binding:"-"`
	CreatedBYID                         uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID                         *uint      `json:"updated_by_id"`
	CreatedAt                           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                           *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	IncidentCategory                  *IncidentCategory                  `gorm:"foreignKey:IncidentCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentType                      *IncidentType                      `gorm:"foreignKey:IncidentTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel                     *SeverityLevel                     `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentStatus                    *IncidentStatus                    `gorm:"foreignKey:IncidentStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit                   *OrganizationUnit                  `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser                   *User                              `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	BusinessOwner                     *User                              `gorm:"foreignKey:BusinessOwnerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentResponseTeam              *IncidentResponseTeam              `gorm:"foreignKey:IncidentResponseTeamID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponseTeamLead                  *User                              `gorm:"foreignKey:ResponseTeamLeadID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentScope                     *IncidentScope                     `gorm:"foreignKey:IncidentScopeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentSensitivityClassification *IncidentSensitivityClassification `gorm:"foreignKey:IncidentSensitivityClassificationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ConfidentialitySeverityLevel      *SeverityLevel                     `gorm:"foreignKey:ConfidentialitySeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IntegritySeverityLevel            *SeverityLevel                     `gorm:"foreignKey:IntegritySeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AvailabilitySeverityLevel         *SeverityLevel                     `gorm:"foreignKey:AvailabilitySeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Documents                         []*IncidentDocument                `gorm:"foreignKey:IncidentRegisterID"`
	Controls                          []*IncidentControl                 `gorm:"foreignKey:IncidentRegisterID"`
	ComplianceRequirements            []*IncidentComplianceRequirement   `gorm:"foreignKey:IncidentRegisterID"`
	Risks                             []*IncidentRisk                    `gorm:"foreignKey:IncidentRegisterID"`
	Processes                         []*IncidentProcess                 `gorm:"foreignKey:IncidentRegisterID"`
	SLARegisters                      []*IncidentSLA                     `gorm:"foreignKey:IncidentRegisterID"`
	AssociatedDocuments               []*IncidentAssociatedDocument      `gorm:"foreignKey:IncidentRegisterID"`
	Assets                            []*IncidentAsset                   `gorm:"foreignKey:IncidentRegisterID"`
	Timelines                         []*IncidentTimeline                `gorm:"foreignKey:IncidentRegisterID"`
	CommunicationLogs                 []*IncidentCommunicationLog        `gorm:"foreignKey:IncidentRegisterID"`
	ActionPlans                       []*IncidentActionPlan              `gorm:"foreignKey:IncidentRegisterID"`
	IncidentClosure                   *IncidentClosure                   `gorm:"foreignKey:IncidentRegisterID"`
	IncidentPostReview                *IncidentPostReview                `gorm:"foreignKey:IncidentRegisterID"`
	IncidentRootCauseAnalyses         []*IncidentRootCauseAnalysis       `gorm:"foreignKey:IncidentRegisterID"`
	StrategicObjectives               []*IncidentStrategicObjective      `gorm:"foreignKey:IncidentRegisterID"`
	KPIs                              []*IncidentKPI                     `gorm:"foreignKey:IncidentRegisterID"`
	Investigation                     *IncidentInvestigation             `gorm:"foreignKey:IncidentRegisterID"`
	ResponseTeamMembers               []*IncidentRegisterResponseTeam    `gorm:"foreignKey:IncidentRegisterID"`
	ReportedByUser                    *User                              `gorm:"foreignKey:ReportedBy"`
	AcknowledgedByUser                *User                              `gorm:"foreignKey:AcknowledgedByID"`
	CreatedByUser                     *User                              `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser                     *User                              `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by IncidentRegister to `incident_register`
func (IncidentRegister) TableName() string {
	return "incident_register"
}
