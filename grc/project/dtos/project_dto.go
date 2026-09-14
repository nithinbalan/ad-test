package dtos

import (
	"time"

	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type ProjectListDTO struct {
	ID                              uint     `json:"id"`
	TitleEn                         string   `json:"title_en"`
	TitleAr                         string   `json:"title_ar"`
	DescriptionEn                   *string  `json:"description_en"`
	DescriptionAr                   *string  `json:"description_ar"`
	Budget                          *float64 `json:"budget"`
	Scope                           *string  `json:"scope"`
	NumberOfResources               *int     `json:"number_of_resources"`
	CompletionPercentage            *float64 `json:"completion_percentage"`
	ProjectTypeTitleEn              string   `json:"project_type_title_en"`
	ProjectTypeTitleAr              string   `json:"project_type_title_ar"`
	ProjectStatusTitleEn            string   `json:"project_status_title_en"`
	ProjectStatusTitleAr            string   `json:"project_status_title_ar"`
	ProjectPhaseID                  *uint    `json:"project_phase_id"`
	ProjectPhaseTitleEn             *string  `json:"project_phase_title_en"`
	ProjectPhaseTitleAr             *string  `json:"project_phase_title_ar"`
	ProjectClassificationID         *uint    `json:"project_classification_id"`
	ProjectClassificationTitleEn    string   `json:"project_classification_title_en"`
	ProjectClassificationTitleAr    string   `json:"project_classification_title_ar"`
	StartDate                       string   `json:"start_date"`
	EndDate                         string   `json:"end_date"`
	ResponsibleUserID               uint     `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string   `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string   `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string   `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string   `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string   `json:"responsible_user_initial"`
	ResponsibleUserColor            string   `json:"responsible_user_color"`
	ResponsibleUserEmail            string   `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string   `json:"responsible_user_profile_image_path"`
	WorkflowStatusID                *uint    `json:"workflow_status_id"`
	WorkflowStatusTitleEn           *string  `json:"workflow_status_title_en"`
	WorkflowStatusTitleAr           *string  `json:"workflow_status_title_ar"`
	WorkflowStatusColorCode         *string  `json:"workflow_status_color_code"`
	WorkflowStatusType              *string  `json:"workflow_status_type"`
	NextReviewLevel                 *int     `json:"next_review_level"`
	IncidentActionPlanID            *uint    `json:"incident_action_plan_id"`
	IncidentActionPlanTitle         *string  `json:"incident_action_plan_title"`
	SubmittedAt                     *string  `json:"submitted_at"`
	SubmittedUserID                 *uint    `json:"submitted_user_id"`
	SubmittedUserFirstNameEn        *string  `json:"submitted_user_first_name_en"`
	SubmittedUserFirstNameAr        *string  `json:"submitted_user_first_name_ar"`
	SubmittedUserLastNameEn         *string  `json:"submitted_user_last_name_en"`
	SubmittedUserLastNameAr         *string  `json:"submitted_user_last_name_ar"`
	ApprovedAt                      *string  `json:"approved_at"`
	ApprovedUserID                  *uint    `json:"approved_user_id"`
	ApprovedUserFirstNameEn         *string  `json:"approved_user_first_name_en"`
	ApprovedUserFirstNameAr         *string  `json:"approved_user_first_name_ar"`
	ApprovedUserLastNameEn          *string  `json:"approved_user_last_name_en"`
	ApprovedUserLastNameAr          *string  `json:"approved_user_last_name_ar"`
	RejectedAt                      *string  `json:"rejected_at"`
	RejectedUserID                  *uint    `json:"rejected_user_id"`
	RejectedUserFirstNameEn         *string  `json:"rejected_user_first_name_en"`
	RejectedUserFirstNameAr         *string  `json:"rejected_user_first_name_ar"`
	RejectedUserLastNameEn          *string  `json:"rejected_user_last_name_en"`
	RejectedUserLastNameAr          *string  `json:"rejected_user_last_name_ar"`
	CreatedUserID                   uint     `json:"created_user_id"`
	CreatedUserFirstNameEn          string   `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string   `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string   `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string   `json:"created_user_last_name_ar"`
	CreatedUserInitial              string   `json:"created_user_initial"`
	CreatedUserColor                string   `json:"created_user_color"`
	CreatedUserEmail                string   `json:"created_user_email"`
	CreatedUserProfileImagePath     string   `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint    `json:"updated_user_id"`
	UpdatedUserFirstNameEn          *string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              *string  `json:"updated_user_initial"`
	UpdatedUserColor                *string  `json:"updated_user_color"`
	UpdatedUserEmail                *string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath     *string  `json:"updated_user_profile_image_path"`
}

type ProjectPaginationInfoDTO struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int   `json:"total_pages"`
	TotalCount  int64 `json:"total_count"`
}

type ProjectListResponseDTO struct {
	Items      []*ProjectListDTO        `json:"items"`
	Pagination ProjectPaginationInfoDTO `json:"pagination"`
}

type ProjectListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type ProjectFormDTO struct {
	TitleEn                 string   `json:"title_en" binding:"required,max=255"`
	TitleAr                 string   `json:"title_ar" binding:"required,max=255"`
	DescriptionEn           *string  `json:"description_en"`
	DescriptionAr           *string  `json:"description_ar"`
	Budget                  *float64 `json:"budget"`
	Scope                   *string  `json:"scope"`
	NumberOfResources       *int     `json:"number_of_resources"`
	ResponsibleUserID       uint     `json:"responsible_user_id" binding:"required"`
	StartDate               string   `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate                 string   `json:"end_date" binding:"required,datetime=2006-01-02"`
	ProjectTypeID           uint     `json:"project_type_id" binding:"required"`
	ProjectStatusID         uint     `json:"project_status_id" binding:"required"`
	ProjectPhaseID          *uint    `json:"project_phase_id"`
	ProjectClassificationID uint     `json:"project_classification_id" binding:"required"`
	IncidentActionPlanID    *uint    `json:"incident_action_plan_id"`
}

type ProjectIncidentActionPlanDTO struct {
	ID                      uint                             `json:"id"`
	IncidentRegisterID      uint                             `json:"incident_register_id"`
	IncidentReferenceNumber *string                          `json:"incident_reference_number"`
	IncidentTitle           *string                          `json:"incident_title"`
	Title                   string                           `json:"title"`
	Description             *string                          `json:"description"`
	TargetDate              *string                          `json:"target_date"`
	Percentage              *uint                            `json:"percentage"`
	ActionPlanStatusID      uint                             `json:"action_plan_status_id"`
	ActionPlanStatusTitleEn *string                          `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr *string                          `json:"action_plan_status_title_ar"`
	ActionPlanStatusColor   *string                          `json:"action_plan_status_color_code"`
	ResponsibleUserID       uint                             `json:"responsible_user_id"`
	ResponsibleUser         *organizationDTOs.UserMinimalDTO `json:"responsible_user"`
}

type ProjectDetailDTO struct {
	ID                      uint                             `json:"id"`
	TitleEn                 string                           `json:"title_en"`
	TitleAr                 string                           `json:"title_ar"`
	DescriptionEn           *string                          `json:"description_en"`
	DescriptionAr           *string                          `json:"description_ar"`
	Budget                  *float64                         `json:"budget"`
	Scope                   *string                          `json:"scope"`
	NumberOfResources       *int                             `json:"number_of_resources"`
	CompletionPercentage    *float64                         `json:"completion_percentage"`
	StartDateString         string                           `json:"start_date"`
	EndDateString           string                           `json:"end_date"`
	ResponsibleUser         organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
	ProjectType             ProjectTypeDTO                   `json:"project_type"`
	ProjectStatus           ProjectStatusDTO                 `json:"project_status"`
	ProjectPhaseID          *uint                            `json:"project_phase_id"`
	ProjectPhase            *ProjectPhaseDTO                 `json:"project_phase"`
	ProjectClassificationID *uint                            `json:"project_classification_id"`
	ProjectClassification   *ProjectClassificationDTO        `json:"project_classification"`
	WorkflowStatusID        *uint                            `json:"workflow_status_id"`
	WorkflowStatus          *masterDTOs.WorkflowStatusDTO    `json:"workflow_status"`
	NextReviewLevel         *int                             `json:"next_review_level"`
	IncidentActionPlanID    *uint                            `json:"incident_action_plan_id"`
	IncidentActionPlan      *ProjectIncidentActionPlanDTO    `json:"incident_action_plan" mapstructure:"-"`
	SubmittedAt             *time.Time                       `json:"submitted_at"`
	SubmittedByUser         *organizationDTOs.UserMinimalDTO `json:"submitted_by_user"`
	ApprovedAt              *time.Time                       `json:"approved_at"`
	ApprovedByUser          *organizationDTOs.UserMinimalDTO `json:"approved_by_user"`
	RejectedAt              *time.Time                       `json:"rejected_at"`
	RejectedByUser          *organizationDTOs.UserMinimalDTO `json:"rejected_by_user"`
	CreatedAt               time.Time                        `json:"created_at"`
	UpdatedAt               *time.Time                       `json:"updated_at"`
	CreatedByUser           organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser           *organizationDTOs.UserMinimalDTO `json:"updated_by"`
}

type ProjectStatusUpdateDTO struct {
	ProjectStatusID uint `json:"project_status_id" binding:"required"`
}
