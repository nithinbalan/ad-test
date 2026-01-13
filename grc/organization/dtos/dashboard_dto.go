package dtos

type DashboardDTO struct {
	TasksAssigned    int64 `json:"tasks_assigned"`
	TasksCompleted   int64 `json:"tasks_completed"`
	OverdueItems     int64 `json:"overdue_items"`
	NewNotifications int64 `json:"new_notifications"`
}

type UserTaskDTO struct {
	ID                              uint   `json:"id"`
	ReferenceNumber                 string `json:"reference_number"`
	ModuleID                        uint   `json:"module_id"`
	Title                           string `json:"title"`
	Description                     string `json:"description"`
	SeverityLevelTitleEn            string `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string `json:"severity_level_title_ar"`
	SeverityLevelColorCode          string `json:"severity_level_color_code"`
	ResponsibleUserID               uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	DelegateeUserID                 *uint  `json:"delegatee_user_id"`
	DelegateeUserFirstNameEn        string `json:"delegatee_user_first_name_en"`
	DelegateeUserFirstNameAr        string `json:"delegatee_user_first_name_ar"`
	DelegateeUserLastNameEn         string `json:"delegatee_user_last_name_en"`
	DelegateeUserLastNameAr         string `json:"delegatee_user_last_name_ar"`
	DelegateeUserInitial            string `json:"delegatee_user_initial"`
	DelegateeUserColor              string `json:"delegatee_user_color"`
	DelegateeUserEmail              string `json:"delegatee_user_email"`
	DelegateeUserProfileImagePath   string `json:"delegatee_user_profile_image_path"`
	TargetDate                      string `json:"target_date"`
	URLParam                        string `json:"url_param"`
	TaskStatusTitleEn               string `json:"task_status_title_en"`
	TaskStatusTitleAr               string `json:"task_status_title_ar"`
	TaskStatusColorCode             string `json:"task_status_color_code"`
	CreatedAt                       string `json:"created_at"`
	UpdatedAt                       *string `json:"updated_at"`
	CreatedUserID                   uint   `json:"created_user_id"`
	CreatedUserFirstNameEn          string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string `json:"created_user_last_name_ar"`
	CreatedUserInitial              string `json:"created_user_initial"`
	CreatedUserColor                string `json:"created_user_color"`
	CreatedUserEmail                string `json:"created_user_email"`
	CreatedUserProfileImagePath     string `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint  `json:"updated_user_id"`
	UpdatedUserFirstNameEn          string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string `json:"updated_user_initial"`
	UpdatedUserColor                string `json:"updated_user_color"`
	UpdatedUserEmail                string `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string `json:"updated_user_profile_image_path"`
}

type AssignedRiskDTO struct {
	ID             uint            `json:"id"`
	Title          string          `json:"title"`
	ReferenceNumber string         `json:"reference_number"`
	Severity       SeverityLevelDTO `json:"severity"`
	LastUpdated    string          `json:"last_updated"`
}

type SeverityLevelDTO struct {
	Title string `json:"title"`
	Color string `json:"color"`
}

type AssignedRiskFlatDTO struct {
	ID             uint   `gorm:"column:id"`
	ReferenceNumber string `gorm:"column:reference_number"`
	Title          string `gorm:"column:title"`
	SeverityTitle  string `gorm:"column:severity_title"`
	SeverityColor  string `gorm:"column:severity_color"`
	LastUpdated    string `gorm:"column:last_updated"`
}

type TaskCountsDTO struct {
	TotalAssigned  int64 `gorm:"column:total_assigned"`
	TotalCompleted int64 `gorm:"column:total_completed"`
	TotalOverdue   int64 `gorm:"column:total_overdue"`
}

type RecentActivityDTO struct {
	Title   string `json:"title"`
	TimeAgo string `json:"time_ago"`
}
