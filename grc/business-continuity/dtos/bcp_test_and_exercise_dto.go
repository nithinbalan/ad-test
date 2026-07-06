package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type BcpTestAndExerciseListDTO struct {
	ID                              uint    `json:"id"`
	Title                           string  `json:"title"`
	Objective                       *string `json:"objective"`
	Scenario                        *string `json:"scenario"`
	SuccessCriteria                 *string `json:"success_criteria"`
	ReviewFrequencyTitleEn          *string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr          *string `json:"review_frequency_title_ar"`
	ResponsibleUserID               *uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      *string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      *string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       *string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       *string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          *string `json:"responsible_user_initial"`
	ResponsibleUserColor            *string `json:"responsible_user_color"`
	ResponsibleUserEmail            *string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath *string `json:"responsible_user_profile_image_path"`
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

type BcpTestAndExerciseListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type BcpTestAndExerciseListResponseDTO struct {
	Items      []*BcpTestAndExerciseListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type BcpTestAndExerciseFormDTO struct {
	Title             string  `json:"title" binding:"required"`
	ResponsibleUserID *uint   `json:"responsible_user_id" binding:"required"`
	ReviewFrequencyID *uint   `json:"review_frequency_id"`
	Objective         *string `json:"objective"`
	Scenario          *string `json:"scenario"`
	SuccessCriteria   *string `json:"success_criteria"`
	SystemRegisterIDs []*uint `json:"system_register_ids"`
}

type BcpTestAndExerciseDetailDTO struct {
	ID                                int                                 `json:"id"`
	Title                             string                              `json:"title"`
	Objective                         *string                             `json:"objective"`
	Scenario                          *string                             `json:"scenario"`
	SuccessCriteria                   *string                             `json:"success_criteria"`
	BusinessContinuityPlan            BusinessContinuityPlanDetailDTO     `json:"business_continuity_plan"`
	ResponsibleUserID                 *uint                               `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ResponsibleUser                   *organizationDTOs.UserMinimalDTO    `json:"responsible_user"`
	ReviewFrequencyID                 *uint                               `json:"review_frequency_id"`
	ReviewFrequency                   *masterDTOs.ReviewFrequencyDTO      `json:"review_frequency"`
	BcpTestAndExerciseBcpProcesses    []*BcpTestAndExerciseBcpProcess     `json:"bcp_test_and_exercise_bcp_processes"`
	BcpTestAndExerciseSystemRegisters []*BcpTestAndExerciseSystemRegister `json:"bcp_test_and_exercise_system_registers"`
	CreatedAt                         time.Time                           `json:"created_at"`
	UpdatedAt                         *time.Time                          `json:"updated_at"`
	CreatedByUser                     organizationDTOs.UserMinimalDTO     `json:"created_by_user"`
	UpdatedByUser                     *organizationDTOs.UserMinimalDTO    `json:"updated_by_user"`
}

type BcpTestAndExerciseBcpProcess struct {
	ID                   *uint         `json:"id"`
	BcpTestAndExerciseID uint          `json:"bcp_test_and_exercise_id"`
	BcpProcessID         uint          `json:"bcp_process_id"`
	BcpProcess           BcpProcessDTO `json:"bcp_process"`
}

type BcpTestAndExerciseSystemRegister struct {
	ID                   *uint                                     `json:"id"`
	BcpTestAndExerciseID uint                                      `json:"bcp_test_and_exercise_id"`
	SystemRegisterID     uint                                      `json:"system_register_id"`
	SystemRegister       BcpTestAndExerciseSystemRegisterDetailDTO `json:"system_register"`
}

type BcpTestAndExerciseSystemRegisterDetailDTO struct {
	ID         uint   `json:"id"`
	TitleEn    string `json:"title_en"`
	TitleAr    string `json:"title_ar"`
	RTO        *int   `json:"rto"`
	RTOUnitID  *uint  `json:"rto_unit_id"`
	MTPD       *int   `json:"mtpd"`
	MTPDUnitID *uint  `json:"mtpd_unit_id"`
	RPO        *int   `json:"rpo"`
	RPOUnitID  *uint  `json:"rpo_unit_id"`
}
