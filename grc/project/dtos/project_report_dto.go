package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ProjectClosureReportDTO struct {
	ID                    uint                             `json:"id"`
	OverallProjectOutcome *string                          `json:"overall_project_outcome"`
	ClosureReason         *string                          `json:"closure_reason"`
	ClosureComments       *string                          `json:"closure_comments"`
	LessonLearned         *string                          `json:"lesson_learned"`
	ApproverID            uint                             `json:"approver_id"`
	ApproverUser          *organizationDTOs.UserMinimalDTO `json:"approver_user"`
	CreatedByUser         organizationDTOs.UserMinimalDTO  `json:"created_by"`
	CreatedAt             time.Time                        `json:"created_at"`
	CreatedAtFormatted    string                           `json:"created_at_formatted"`
}

type ProjectReportDetailDTO struct {
	ID                      uint                                                `json:"id"`
	TitleEn                 string                                              `json:"title_en"`
	TitleAr                 string                                              `json:"title_ar"`
	DescriptionEn           *string                                             `json:"description_en"`
	DescriptionAr           *string                                             `json:"description_ar"`
	Budget                  *float64                                            `json:"budget"`
	Scope                   *string                                             `json:"scope"`
	NumberOfResources       *int                                                `json:"number_of_resources"`
	CompletionPercentage    *float64                                            `json:"completion_percentage"`
	StartDateString         string                                              `json:"start_date"`
	EndDateString           string                                              `json:"end_date"`
	StartDateFormatted      string                                              `json:"start_date_formatted"`
	EndDateFormatted        string                                              `json:"end_date_formatted"`
	ResponsibleUser         organizationDTOs.UserMinimalDTO                     `json:"responsible_user"`
	ProjectType             ProjectTypeDTO                                      `json:"project_type"`
	ProjectStatus           ProjectStatusDTO                                    `json:"project_status"`
	ProjectPhaseID          *uint                                               `json:"project_phase_id"`
	ProjectPhase            *ProjectPhaseDTO                                    `json:"project_phase"`
	ProjectClassificationID *uint                                               `json:"project_classification_id"`
	ProjectClassification   *ProjectClassificationDTO                           `json:"project_classification"`
	CreatedAt               time.Time                                           `json:"created_at"`
	UpdatedAt               *time.Time                                          `json:"updated_at"`
	CreatedByUser           organizationDTOs.UserMinimalDTO                     `json:"created_by"`
	UpdatedByUser           *organizationDTOs.UserMinimalDTO                    `json:"updated_by"`
	GeneratedAt             string                                              `json:"generated_at"`
	Milestones              []*ProjectMilestoneListDTO                          `json:"milestones"`
	Controls                []*commonDTOs.AssociatedControlListDTO              `json:"controls"`
	Risks                   []*commonDTOs.AssociatedRiskListDTO                 `json:"risks"`
	Processes               []*commonDTOs.AssociatedProcessListDTO              `json:"processes"`
	StrategicObjectives     []*commonDTOs.AssociatedStrategicObjectiveListDTO   `json:"strategic_objectives"`
	KPIs                    []*commonDTOs.AssociatedKPIRegisterListDTO          `json:"kpis"`
	AssociatedDocuments     []*commonDTOs.AssociatedDocumentListDTO             `json:"associated_documents"`
	ActionPlans             []*ProjectActionPlanListDTO                         `json:"action_plans"`
	ChangeRequests          []*ProjectChangeRequestListDTO                      `json:"change_requests"`
	DecisionLogs            []*ProjectDecisionLogListDTO                        `json:"decision_logs"`
	SupportingDocuments     []*ProjectSupportingDocumentListDTO                 `json:"supporting_documents"`
	Closure                 *ProjectClosureReportDTO                            `json:"closure"`
}
