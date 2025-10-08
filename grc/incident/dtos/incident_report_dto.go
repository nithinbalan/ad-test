package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type IncidentReportDetailDTO struct {
	ID                     uint                                              `json:"id"`
	ReferenceNumber        string                                            `json:"reference_number"`
	Title                  string                                            `json:"title"`
	Description            string                                            `json:"description"`
	IncidentAtFormatted    string                                            `json:"incident_at_formatted"`
	ReportedAtFormatted    string                                            `json:"reported_at_formatted"`
	Location               *string                                           `json:"location"`
	IncidentCategory       IncidentCategoryDTO                               `json:"incident_category"`
	IncidentType           IncidentTypeDTO                                   `json:"incident_type"`
	SeverityLevel          *masterDTOs.SeverityLevelDTO                      `json:"severity_level"`
	IncidentStatus         IncidentStatusDTO                                 `json:"incident_status"`
	ReportedByUser         organizationDTOs.UserMinimalDTO                   `json:"reported_by"`
	Documents              []*IncidentDocumentDTO                            `json:"documents"`
	Controls               []*commonDTOs.AssociatedControlListDTO            `json:"controls"`
	Risks                  []*commonDTOs.AssociatedRiskListDTO               `json:"risks"`
	Assets                 []*commonDTOs.AssociatedAssetListDTO              `json:"assets"`
	Processes              []*commonDTOs.AssociatedProcessListDTO            `json:"processes"`
	Timelines              []*IncidentTimelineListDTO                        `json:"timelines"`
	ResponsibleUnit        settingsDTOs.OrganizationUnitDTO                  `json:"responsible_unit"`
	ResponsibleUser        organizationDTOs.UserMinimalDTO                   `json:"responsible_user"`
	ResponsibleUserID      uint                                              `json:"responsible_user_id" binding:"required" gorm:"not null"`
	SLARegisters           []*IncidentSLARegisterDTO                         `json:"sla_registers"`
	AssociatedDocuments    []*commonDTOs.AssociatedDocumentListDTO           `json:"associated_documents"`
	CreatedAt              time.Time                                         `json:"created_at"`
	UpdatedAt              *time.Time                                        `json:"updated_at"`
	CreatedByUser          organizationDTOs.UserMinimalDTO                   `json:"created_by"`
	UpdatedByUser          organizationDTOs.UserMinimalDTO                   `json:"updated_by"`
	GeneratedAt            string                                            `json:"generated_at"`
	CommunicationLogs      []*IncidentCommunicationLogDetailDTO              `json:"communication_log"`
	ActionPlans            []*IncidentActionPlanListDTO                      `json:"action_plan"`
	IncidentClosure        *IncidentClosureDetailDTO                         `json:"incident_closure"`
	IncidentPostReview     *IncidentPostReviewDetailDTO                      `json:"incident_post_review"`
	RootCauseAnalyses      []*IncidentRootCauseAnalysisListDTO               `json:"root_cause_analysis"`
	StrategicObjectives    []*commonDTOs.AssociatedStrategicObjectiveListDTO `json:"strategic_objectives"`
	KPIs                   []*commonDTOs.AssociatedKPIRegisterListDTO        `json:"kpi"`
	ComplianceRequirements []*commonDTOs.AssociatedClauseListDTO             `json:"compliance_clause"`
	Investigation          *IncidentInvestigationDetailDTO                   `json:"investigation"`
	AppUrl                 string                                            `json:"url"`
}
