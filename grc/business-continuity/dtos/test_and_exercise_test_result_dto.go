package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type TestAndExerciseTestResultFormDTO struct {
	TestCompletionDate            string                             `json:"test_completion_date" binding:"required"`
	CommunicationFeedbackOptionID uint                               `json:"communication_feedback_option_id" binding:"required"`
	RecoveryFeedbackOptionID      uint                               `json:"recovery_feedback_option_id" binding:"required"`
	TestResultOptionID            uint                               `json:"test_result_option_id" binding:"required"`
	TestObservations              *string                            `json:"test_observations"`
	Recommendations               *string                            `json:"recommendations"`
	Findings                      *string                            `json:"findings"`
	TestAndExerciseProcesses      []*TestAndExerciseProcessesFormDTO `json:"test_and_exercise_processes" binding:"required"`
	Documents                     []*commonDTOs.DocumentFormDTO      `json:"documents" binding:"omitempty"`
}

type TestAndExerciseProcessesFormDTO struct {
	ID               uint  `json:"id"`
	IsRecovered      bool  `json:"is_recovered"`
	TimelinessID     *uint `json:"timeliness_id"`
	ActualRto        *int  `json:"actual_rto"`
	ActualRtoUnitId  *uint `json:"actual_rto_unit_id"`
	ActualMtpd       *int  `json:"actual_mtpd"`
	ActualMtpdUnitId *uint `json:"actual_mtpd_unit_id"`
	ActualRpo        *int  `json:"actual_rpo"`
	ActualRpoUnitId  *uint `json:"actual_rpo_unit_id"`
}

type TestAndExerciseTestResultDetailDTO struct {
	ID                            uint                                    `json:"id"`
	TestAndExercise               TestAndExerciseDetailDTO                `json:"test_and_exercise"`
	TestCompletionDateString      string                                  `json:"test_completion_date"`
	CommunicationFeedbackOptionID uint                                    `json:"communication_feedback_option_id"`
	RecoveryFeedbackOptionID      uint                                    `json:"recovery_feedback_option_id"`
	TestResultOptionID            uint                                    `json:"test_result_option_id"`
	TestObservations              *string                                 `json:"test_observations"`
	Recommendations               *string                                 `json:"recommendations"`
	Findings                      *string                                 `json:"findings"`
	CommunicationFeedbackOption   FeedbackOptionDTO                       `json:"communication_feedback_option"`
	RecoveryFeedbackOption        FeedbackOptionDTO                       `json:"recovery_feedback_option"`
	TestResultOption              TestResultOptionDTO                     `json:"test_result_option"`
	Documents                     []*TestAndExerciseTestResultDocumentDTO `json:"documents" gorm:"-"`
	CreatedAt                     time.Time                               `json:"created_at"`
	UpdatedAt                     *time.Time                              `json:"updated_at"`
	CreatedByUser                 organizationDTOs.UserMinimalDTO         `json:"created_by_user"`
	UpdatedByUser                 *organizationDTOs.UserMinimalDTO        `json:"updated_by_user"`
}

type TestAndExerciseProcessDetailsDTO struct {
	ID                uint                    `json:"id"`
	Process           ProcessDetailDTO        `json:"process"`
	IsRecovered       bool                    `json:"is_recovered"`
	TimelinessID      uint                    `json:"timeliness_id"`
	Timeliness        TimelinessDTO           `json:"timeliness"`
	ActualRto         *int                    `json:"actual_rto"`
	ActualRtoUnitId   *uint                   `json:"actual_rto_unit_id"`
	ActualRtoUnit     *masterDTOs.TimeUnitDTO `json:"actual_rto_unit"`
	ActualMtpd        *int                    `json:"actual_mtpd"`
	ActualMtpdUnitId  *uint                   `json:"actual_mtpd_unit_id"`
	ActualMtpdUnit    *masterDTOs.TimeUnitDTO `json:"actual_mtpd_unit"`
	ActualRpo         *int                    `json:"actual_rpo"`
	ActualRpoUnitId   *uint                   `json:"actual_rpo_unit_id"`
	ActualRpoUnit     *masterDTOs.TimeUnitDTO `json:"actual_rpo_unit"`
	PlannedRto        *int                    `json:"planned_rto"`
	PlannedRtoUnitId  *uint                   `json:"planned_rto_unit_id"`
	PlannedRtoUnit    *masterDTOs.TimeUnitDTO `json:"planned_rto_unit"`
	PlannedMtpd       *int                    `json:"planned_mtpd"`
	PlannedMtpdUnitId *uint                   `json:"planned_mtpd_unit_id"`
	PlannedMtpdUnit   *masterDTOs.TimeUnitDTO `json:"planned_mtpd_unit"`
	PlannedRpo        *int                    `json:"planned_rpo"`
	PlannedRpoUnitId  *uint                   `json:"planned_rpo_unit_id"`
	PlannedRpoUnit    *masterDTOs.TimeUnitDTO `json:"planned_rpo_unit"`
}

type TestAndExerciseTestResultParsedDates struct {
	TestCompletionDate *time.Time
}

type TestAndExerciseTestResultDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}
