package dtos

import "time"

type UserDelegationFormDTO struct {
	UserID          uint    `json:"user_id" binding:"required"`
	DelegateeUserID uint    `json:"delegatee_user_id" binding:"required"`
	Comments        *string `json:"comments"`
	TaskIDs         []uint  `json:"task_ids" binding:"required"`
	ModuleIDs       []uint  `json:"module_ids" binding:"required"`
}

type UserDelegationHistoryListDTO struct {
	DelegateeUserID               uint       `json:"delegatee_user_id"`
	DelegateeUserFirstNameEn      string     `json:"delegatee_user_first_name_en"`
	DelegateeUserFirstNameAr      string     `json:"delegatee_user_first_name_ar"`
	DelegateeUserLastNameEn       string     `json:"delegatee_user_last_name_en"`
	DelegateeUserLastNameAr       string     `json:"delegatee_user_last_name_ar"`
	DelegateeUserInitial          string     `json:"delegatee_user_initial"`
	DelegateeUserColor            string     `json:"delegatee_user_color"`
	DelegateeUserEmail            string     `json:"delegatee_user_email"`
	DelegateeUserProfileImagePath string     `json:"delegatee_user_profile_image_path"`
	ModuleCount                   int        `json:"module_count"`
	TaskCount                     int        `json:"task_count"`
	Comments                      *string    `json:"comments"`
	AssignedDate                  time.Time  `json:"assigned_date"`
	CreatedUserID                 uint       `json:"created_user_id"`
	CreatedUserFirstNameEn        string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr        string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn         string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr         string     `json:"created_user_last_name_ar"`
	CreatedUserInitial            string     `json:"created_user_initial"`
	CreatedUserColor              string     `json:"created_user_color"`
	CreatedUserEmail              string     `json:"created_user_email"`
	CreatedUserProfileImagePath   string     `json:"created_user_profile_image_path"`
	UpdatedUserID                 *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn        *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr        *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn         *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr         *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial            *string    `json:"updated_user_initial"`
	UpdatedUserColor              *string    `json:"updated_user_color"`
	UpdatedUserEmail              *string    `json:"updated_user_email"`
	UpdatedUserProfileImagePath   *string    `json:"updated_user_profile_image_path"`
	CreatedAt                     time.Time  `json:"created_at"`
	UpdatedAt                     *time.Time `json:"updated_at"`
}

type UserDelegationHistoryListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Value       int
	QuickFilter []string
}
