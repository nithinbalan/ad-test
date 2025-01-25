package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
)

type TestAndExerciseListDTO struct {
	ID                              uint    `json:"id"`
	ReferenceNumber                 string  `json:"reference_number"`
	BusinessFunctionID              uint    `json:"business_function_id"`
	BusinessFunctionTitle           string  `json:"business_function_title"`
	BcpTestAndExerciseTitle         string  `json:"bcp_test_and_exercise_title"`
	ResponsibleUserID               *uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      *string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      *string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       *string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       *string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string `json:"responsible_user_initial"`
	ResponsibleUserColor            *string `json:"responsible_user_color"`
	ResponsibleUserEmail            *string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath *string `json:"responsible_user_profile_image_path"`
	TestingDate                     string  `json:"testing_date"`
	Description                     *string `json:"description"`
	TestAndExerciseStatusTitleEn    string  `json:"test_and_exercise_status_title_en"`
	TestAndExerciseStatusTitleAr    string  `json:"test_and_exercise_status_title_ar"`
	TestAndExerciseStatusColorCode  string  `json:"test_and_exercise_status_color_code"`
	CreatedAt                       string  `json:"created_at"`
	CreatedUserID                   uint    `json:"created_user_id"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string  `json:"created_user_last_name_ar"`
	CreatedUserInitial              string  `json:"created_user_initial"`
	CreatedUserColor                string  `json:"created_user_color"`
	CreatedUserEmail                string  `json:"created_user_email"`
	CreatedUserProfileImagePath     string  `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn          *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              *string `json:"updated_user_initial"`
	UpdatedUserColor                *string `json:"updated_user_color"`
	UpdatedUserEmail                *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath     *string `json:"updated_user_profile_image_path"`
}

type TestAndExerciseListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type TestAndExerciseListResponseDTO struct {
	Items      []*TestAndExerciseListDTO    `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type TestAndExerciseFormDTO struct {
	BusinessContinuityPlanID uint    `json:"business_continuity_plan_id" binding:"required"`
	BcpTestAndExerciseID     uint    `json:"bcp_test_and_exercise_id" binding:"required"`
	ResponsibleUserID        uint    `json:"responsible_user_id" binding:"required"`
	TestingDate              string  `json:"testing_date" binding:"required"`
	Description              *string `json:"description"`
}

type TestAndExerciseDetailDTO struct {
	ID                                      int                                          `json:"id"`
	ReferenceNumber                         string                                       `json:"reference_number"`
	BusinessContinuityPlanID                uint                                         `json:"business_continuity_plan_id"`
	BusinessContinuityPlan                  BusinessContinuityPlanDetailDTO              `json:"business_continuity_plan"`
	ResponsibleUserID                       uint                                         `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ResponsibleUser                         *organizationDTOs.UserMinimalDTO             `json:"responsible_user"`
	BcpTestAndExerciseID                    uint                                         `json:"bcp_test_and_exercise_id"`
	BcpTestAndExercise                      *BcpTestAndExerciseDetailDTO                 `json:"bcp_test_and_exercise"`
	TestAndExerciseProcesses                []*TestAndExerciseProcessDTO                 `json:"test_and_exercise_processes"`
	TestAndExerciseRecoveryActionPlans      []*TestAndExerciseRecoveryActionPlanDTO      `json:"test_and_exercise_recovery_action_plans"`
	TestAndExerciseCommunicationActionPlans []*TestAndExerciseCommunicationActionPlanDTO `json:"test_and_exercise_communication_action_plans"`
	TestingDateString                       string                                       `json:"testing_date"`
	Description                             *string                                      `json:"description"`
	TestAndExerciseStatusID                 uint                                         `json:"test_and_exercise_status_id"`
	TestAndExerciseStatus                   TestAndExerciseStatusDTO                     `json:"test_and_exercise_status"`
	CreatedAt                               time.Time                                    `json:"created_at"`
	UpdatedAt                               *time.Time                                   `json:"updated_at"`
	CreatedByUser                           organizationDTOs.UserMinimalDTO              `json:"created_by_user"`
	UpdatedByUser                           *organizationDTOs.UserMinimalDTO             `json:"updated_by_user"`
}

type TestAndExerciseProcessDTO struct {
	ID                *uint                             `json:"id"`
	TestAndExerciseID uint                              `json:"test_and_exercise_id"`
	ProcessID         uint                              `json:"process_id"`
	Process           organizationDTOs.ProcessDetailDTO `json:"process"`
	PlannedRto        *int                              `json:"planned_rto"`
	PlannedRtoUnitId  *uint                             `json:"planned_rto_unit_id"`
	PlannedRtoUnit    *masterDTOs.TimeUnitDTO           `json:"planned_rto_unit"`
	PlannedMtpd       *int                              `json:"planned_mtpd"`
	PlannedMtpdUnitId *uint                             `json:"planned_mtpd_unit_id"`
	PlannedMtpdUnit   *masterDTOs.TimeUnitDTO           `json:"planned_mtpd_unit"`
	PlannedRpo        *int                              `json:"planned_rpo"`
	PlannedRpoUnitId  *uint                             `json:"planned_rpo_unit_id"`
	PlannedRpoUnit    *masterDTOs.TimeUnitDTO           `json:"planned_rpo_unit"`
	ActualRto         *int                              `json:"actual_rto"`
	ActualRtoUnitId   *uint                             `json:"actual_rto_unit_id"`
	ActualRtoUnit     *masterDTOs.TimeUnitDTO           `json:"actual_rto_unit"`
	ActualMtpd        *int                              `json:"actual_mtpd"`
	ActualMtpdUnitId  *uint                             `json:"actual_mtpd_unit_id"`
	ActualMtpdUnit    *masterDTOs.TimeUnitDTO           `json:"actual_mtpd_unit"`
	ActualRpo         *int                              `json:"actual_rpo"`
	ActualRpoUnitId   *uint                             `json:"actual_rpo_unit_id"`
	ActualRpoUnit     *masterDTOs.TimeUnitDTO           `json:"actual_rpo_unit"`
	TimelinessID      *uint                             `json:"timeliness_id"`
	Timeliness        TimelinessDTO                     `json:"timeliness"`
	IsRecovered       bool                              `json:"is_recovered"`
}

type TestAndExerciseRecoveryActionPlanDTO struct {
	ID                *uint                             `json:"id"`
	TestAndExerciseID uint                              `json:"test_and_exercise_id"`
	ProcessID         uint                              `json:"process_id"`
	Process           organizationDTOs.ProcessDetailDTO `json:"process"`
	BcpPhaseID        uint                              `json:"bcp_phase_id"`
	BcpPhase          BcpPhaseDTO                       `json:"bcp_phase"`
	SeverityLevel     *masterDTOs.SeverityLevelDTO      `json:"severity_level"`
	ResponsibleUserID uint                              `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ResponsibleUser   organizationDTOs.UserMinimalDTO   `json:"responsible_user"`
	ResponsibleUnitID uint                              `json:"responsible_unit_id"`
	ResponsibleUnit   settingsDTOs.OrganizationUnitDTO  `json:"responsible_unit"`
	TargetDateString  *string                           `json:"target_date"`
}

type TestAndExerciseCommunicationActionPlanDTO struct {
	ID                *uint                             `json:"id"`
	TestAndExerciseID uint                              `json:"test_and_exercise_id"`
	ProcessID         uint                              `json:"process_id"`
	Process           organizationDTOs.ProcessDetailDTO `json:"process"`
	BcpPhaseID        uint                              `json:"bcp_phase_id"`
	BcpPhase          BcpPhaseDTO                       `json:"bcp_phase"`
	ResponsibleUserID uint                              `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ResponsibleUser   organizationDTOs.UserMinimalDTO   `json:"responsible_user"`
	ResponsibleUnitID uint                              `json:"responsible_unit_id"`
	ResponsibleUnit   settingsDTOs.OrganizationUnitDTO  `json:"responsible_unit"`
	Frequency         *string                           `json:"frequency"`
	Recipient         *string                           `json:"recipient"`
	Channel           *string                           `json:"channel"`
}

type TestAndExerciseParsedDates struct {
	TestingDate *time.Time
}
