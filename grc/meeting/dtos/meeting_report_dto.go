package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type MeetingReportDetailDTO struct {
	ID                     uint                                              `json:"id"`
	ReferenceNumber        string                                            `json:"reference_number"`
	MeetingID              *uint                                             `json:"meeting_id"`
	MeetingTypeID          uint                                              `json:"meeting_type_id"`
	MeetingFrequencyID     *uint                                             `json:"meeting_frequency_id"`
	Meeting                *MeetingDetailDTO                                 `json:"meeting"`
	MeetingType            MeetingTypeDTO                                    `json:"meeting_type"`
	MeetingFrequency       *MeetingFrequencyDTO                              `json:"meeting_frequency"`
	Title                  string                                            `json:"title"`
	Description            *string                                           `json:"description"`
	Location               *string                                           `json:"location"`
	DateString             *string                                           `json:"date"`
	StartTime              *string                                           `json:"start_time"`
	EndTime                *string                                           `json:"end_time"`
	IsPublic               bool                                              `json:"is_public"`
	Organizer              organizationDTOs.UserMinimalDTO                   `json:"organizer"`
	CreatedByUser          organizationDTOs.UserMinimalDTO                   `json:"created_by"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO                  `json:"updated_by"`
	CreatedAt              time.Time                                         `json:"created_at"`
	UpdatedAt              *time.Time                                        `json:"updated_at"`
	MeetingStatusID        uint                                              `json:"meeting_status_id"`
	MeetingStatus          MeetingStatusDTO                                  `json:"meeting_status"`
	Agendas                []*MeetingAgendaDetailDTO                         `json:"agendas"`
	Participants           []*MeetingParticipantDTO                          `json:"participants"`
	Controls               []*commonDTOs.AssociatedControlListDTO            `json:"controls"`
	ComplianceRequirements []*commonDTOs.AssociatedClauseListDTO             `json:"compliance_requirements"`
	Risks                  []*commonDTOs.AssociatedRiskListDTO               `json:"risks"`
	Processes              []*commonDTOs.AssociatedProcessListDTO            `json:"processes"`
	Incidents              []*MeetingIncidentDTO                             `json:"incidents"`
	Documents              []*MeetingDocumentDTO                             `json:"documents"`
	GeneratedAt            string                                            `json:"generated_at"`
	MeetingMinutes         []*MeetingMinuteDetailDTO                         `json:"minutes"`
	StrategicObjectives    []*commonDTOs.AssociatedStrategicObjectiveListDTO `json:"strategic_objectives"`
	KPIs                   []*commonDTOs.AssociatedKPIRegisterListDTO        `json:"kpi"`
	AssociatedDocuments    []*commonDTOs.AssociatedDocumentListDTO           `json:"associated_documents"`
	Assets                 []*commonDTOs.AssociatedAssetListDTO              `json:"assets"`
	ActionPlans            []*MeetingActionPlanListDTO                       `json:"action_plans"`
	AppUrl                 string                                            `json:"url"`
}
