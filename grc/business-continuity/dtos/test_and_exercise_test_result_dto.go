package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type TestAndExerciseTestResultFormDTO struct {
	TestObservations               *string                                  `json:"test_observations"`
	Recommendations                *string                                  `json:"recommendations"`
	Findings                       *string                                  `json:"findings"`
	TestAndExerciseProcesses       []*TestAndExerciseProcessesFormDTO       `json:"test_and_exercise_processes" binding:"required"`
	TestAndExerciseSystemRegisters []*TestAndExerciseSystemRegistersFormDTO `json:"test_and_exercise_system_registers" binding:"required"`
	Documents                      []*commonDTOs.DocumentFormDTO            `json:"documents" binding:"omitempty"`
}

type TestResultSystemRegisterFormDTO struct {
	TestCompletionDate            string `json:"test_completion_date"`
	CommunicationFeedbackOptionID uint   `json:"communication_feedback_option_id"`
	RecoveryFeedbackOptionID      uint   `json:"recovery_feedback_option_id"`
	TestResultOptionID            uint   `json:"test_result_option_id"`
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
	ID                        uint                                     `json:"id"`
	TestAndExercise           TestAndExerciseDetailDTO                 `json:"test_and_exercise"`
	TestObservations          *string                                  `json:"test_observations"`
	Recommendations           *string                                  `json:"recommendations"`
	Findings                  *string                                  `json:"findings"`
	Documents                 []*TestAndExerciseTestResultDocumentDTO  `json:"documents" gorm:"-"`
	TestResultSystemRegisters []*TestResultSystemRegisterDetailDTO     `json:"test_result_system_registers"`
	PerformanceMeasures       []*TestResultPerformanceMeasureDetailDTO `json:"performance_measures"`
	CreatedAt                 time.Time                                `json:"created_at"`
	UpdatedAt                 *time.Time                               `json:"updated_at"`
	CreatedByUser             organizationDTOs.UserMinimalDTO          `json:"created_by_user"`
}

type TestResultSystemRegisterDetailDTO struct {
	ID                            uint                `json:"id"`
	SystemRegister                SystemDTO           `json:"system_register"`
	TestCompletionDateString      string              `json:"test_completion_date"`
	CommunicationFeedbackOptionID uint                `json:"communication_feedback_option_id"`
	RecoveryFeedbackOptionID      uint                `json:"recovery_feedback_option_id"`
	TestResultOptionID            uint                `json:"test_result_option_id"`
	CommunicationFeedbackOption   FeedbackOptionDTO   `json:"communication_feedback_option"`
	RecoveryFeedbackOption        FeedbackOptionDTO   `json:"recovery_feedback_option"`
	TestResultOption              TestResultOptionDTO `json:"test_result_option"`
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

type TestAndExerciseSystemRegistersFormDTO struct {
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

type SystemDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

type TestResultPerformanceMeasureFormDTO struct {
	Criterion      *string `json:"criterion"`
	Target         *string `json:"target"`
	AchievedResult *string `json:"achieved_result"`
	Remarks        *string `json:"remarks"`
}

type TestResultPerformanceMeasureBulkCreateDTO struct {
	PerformanceMeasures []*TestResultPerformanceMeasureFormDTO `json:"performance_measures" binding:"required,dive"`
}

type TestResultPerformanceMeasureDetailDTO struct {
	ID             uint    `json:"id"`
	TestResultID   uint    `json:"test_result_id"`
	Criterion      *string `json:"criterion"`
	Target         *string `json:"target"`
	AchievedResult *string `json:"achieved_result"`
	Remarks        *string `json:"remarks"`
}
