package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type AuditPlanBudgetListDTO struct {
	ID                           uint    `json:"id"`
	AuditPlanID                  uint    `json:"audit_plan_id"`
	ExpenseCategoryID            uint    `json:"expense_category_id"`
	ExpenseCategoryTitleEn       string  `json:"expense_category_title_en"`
	ExpenseCategoryDescriptionEn *string `json:"expense_category_description_en"`
	ExpenseCategoryDescriptionAr *string `json:"expense_category_description_ar"`
	Title                        string  `json:"title"`
	EstimatedCost                float64 `json:"estimated_cost"`
	Comments                     *string `json:"comments"`
	AuditPlanTitle               string  `json:"audit_plan_title"`
	CreatedUserFirstNameEn       string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr       string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn        string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr        string  `json:"created_user_last_name_ar"`
	CreatedUserInitial           string  `json:"created_user_initial"`
	CreatedUserColor             string  `json:"created_user_color"`
	CreatedUserEmail             string  `json:"created_user_email"`
	CreatedUserProfileImagePath  string  `json:"created_user_profile_image_path"`
	UpdatedUserFirstNameEn       string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr       string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn        string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr        string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial           string  `json:"updated_user_initial"`
	UpdatedUserColor             string  `json:"updated_user_color"`
	UpdatedUserEmail             string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath  string  `json:"updated_user_profile_image_path"`
	CreatedAt                    string  `json:"created_at"`
	UpdatedAt                    *string `json:"updated_at"`
}

type AuditPlanBudgetFormDTO struct {
	ExpenseCategoryID uint    `json:"expense_category_id" binding:"required"`
	Title             string  `json:"title" binding:"required,max=255"`
	EstimatedCost     float64 `json:"estimated_cost" binding:"required"`
	Comments          *string `json:"comments"`
}

type AuditPlanBudgetDetailDTO struct {
	ID              uint                             `json:"id"`
	AuditPlanID     uint                             `json:"audit_plan_id"`
	ExpenseCategory ExpenseCategoryDTO               `json:"expense_category"`
	AuditPlan       AuditPlanDetailDTO               `json:"audit_plan"`
	Title           string                           `json:"title"`
	EstimatedCost   float64                          `json:"estimated_cost"`
	Comments        *string                          `json:"comments"`
	CreatedAt       time.Time                        `json:"created_at"`
	UpdatedAt       *time.Time                       `json:"updated_at"`
	CreatedByUser   organizationDTOs.UserMinimalDTO  `json:"created_by_id"`
	UpdatedByUser   *organizationDTOs.UserMinimalDTO `json:"updated_by_id"`
}

type AuditPlanBudgetListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Filters     map[string][]string
	QuickFilter []string
}

type ExpenseCategoryDTO struct {
	ID            uint    `json:"id"`
	TitleEn       string  `json:"title_en"`
	TitleAr       string  `json:"title_ar"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
}
