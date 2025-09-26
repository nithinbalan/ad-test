package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type IncidentReportPDFDetailDTO struct {
	Incident            IncidentRegisterDetailDTO                         `json:"incident"`
	Timelines           []*IncidentTimelineListDTO                        `json:"timelines"`
	StrategicObjectives []*commonDTOs.AssociatedStrategicObjectiveListDTO `json:"strategic_objectives"`
	KPIs                []*commonDTOs.AssociatedKPIRegisterListDTO        `json:"kpis"`
	Controls            []*commonDTOs.AssociatedControlListDTO            `json:"controls"`
	Compliance          []*IncidentComplianceRequirementListDTO            `json:"compliance"`
	Processes           []*commonDTOs.AssociatedProcessListDTO            `json:"processes"`
	Risks               []*IncidentRiskDTO                                `json:"risks"`
	Documents           []*commonDTOs.AssociatedDocumentListDTO           `json:"documents"`
	Assets              []*commonDTOs.AssociatedAssetListDTO              `json:"assets"`
	Investigation       *IncidentInvestigationDetailDTO                   `json:"investigation"`
	RootCauses          []*IncidentRootCauseAnalysisListDTO               `json:"root_cause_analysis"`
	CommunicationLogs   []*IncidentCommunicationLogListDTO                `json:"communication_logs"`
	ActionPlans         []*IncidentActionPlanListDTO                      `json:"action_plans"`
	Closure             *IncidentClosureDetailDTO                         `json:"closure"`
	PostReview          *IncidentPostReviewDetailDTO                      `json:"post_review"`
	GeneratedAt         string                                            `json:"generated_at"`
	IncidentAt          string                                            `json:"incident_at_formatted"`
	ReportedAt          string                                            `json:"reported_at_formatted"`
}
