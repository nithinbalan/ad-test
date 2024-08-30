package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditChecklistListDTO struct {
	ID                           uint       `json:"id"`
	AuditID                      uint       `json:"audit_id"`
	AuditPlanChecklistID         *uint      `json:"audit_plan_checklist_id"`
	Question                     string     `json:"question"`
	ConformanceID                *uint      `json:"conformance_id"`
	ConformanceTitleEn           string     `json:"conformance_title_en"`
	ConformanceTitleAr           string     `json:"conformance_title_ar"`
	ConformanceType              string     `json:"conformance_type"`
	ConformanceColorCode         string     `json:"conformance_color_code"`
	Comment                      *string    `json:"comment"`
	AnsweredByID                 *uint      `json:"answered_by_id"`
	AnsweredAt                   *time.Time `json:"answered_at"`
	CreatedByID                  uint       `json:"created_by_id"`
	CreatedAt                    time.Time  `json:"created_at"`
	UpdatedByID                  *uint      `json:"updated_by_id"`
	UpdatedAt                    *time.Time `json:"updated_at"`
	StatusID                     uint       `json:"status_id"`
	StatusTitleEn                string     `json:"status_title_en"`
	StatusTitleAr                string     `json:"status_title_ar"`
	StatusType                   string     `json:"status_type"`
	StatusColorCode              string     `json:"status_color_code"`
	AnsweredUserFirstNameEn      string     `json:"answered_user_first_name_en"`
	AnsweredUserFirstNameAr      string     `json:"answered_user_first_name_ar"`
	AnsweredUserLastNameEn       string     `json:"answered_user_last_name_en"`
	AnsweredUserLastNameAr       string     `json:"answered_user_last_name_ar"`
	AnsweredUserInitial          string     `json:"answered_user_initial"`
	AnsweredUserColor            string     `json:"answered_user_color"`
	AnsweredUserEmail            string     `json:"answered_user_email"`
	AnsweredUserProfileImagePath string     `json:"answered_user_profile_image_path"`
	CreatedUserFirstNameEn       string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr       string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn        string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr        string     `json:"created_user_last_name_ar"`
	CreatedUserInitial           string     `json:"created_user_initial"`
	CreatedUserColor             string     `json:"created_user_color"`
	CreatedUserEmail             string     `json:"created_user_email"`
	CreatedUserProfileImagePath  string     `json:"created_user_profile_image_path"`
	UpdatedUserID                *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn       *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr       *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn        *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr        *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial           *string    `json:"updated_user_initial"`
	UpdatedUserColor             *string    `json:"updated_user_color"`
	UpdatedUserEmail             *string    `json:"updated_user_email"`
	UpdatedUserProfileImagePath  *string    `json:"updated_user_profile_image_path"`
}

type AuditChecklistFormDTO struct {
	AuditPlanChecklistID *uint   `json:"audit_plan_checklist_id"`
	Question             string  `json:"question" binding:"required"`
	ConformanceID        *uint   `json:"conformance_id"`
	Comment              *string `json:"comment"`
	StatusID             uint    `json:"status_id" binding:"required"`
}

type AuditChecklistDetailDTO struct {
	ID                   uint                             `json:"id"`
	AuditID              uint                             `json:"audit_id"`
	AuditPlanChecklistID *uint                            `json:"audit_plan_checklist_id"`
	Question             string                           `json:"question"`
	ConformanceID        *uint                            `json:"conformance_id"`
	Comment              *string                          `json:"comment"`
	AnsweredByID         *uint                            `json:"answered_by_id"`
	AnsweredAt           *time.Time                       `json:"answered_at"`
	CreatedByID          uint                             `json:"created_by_id"`
	UpdatedByID          *uint                            `json:"updated_by_id"`
	StatusID             uint                             `json:"status_id"`
	AuditPlanChecklist   *AuditPlanChecklistDetailDTO     `json:"audit_plan_checklist"`
	Conformance          *masterDTOs.ConformanceDTO       `json:"conformance"`
	AnsweredBy           *organizationDTOs.UserMinimalDTO `json:"answered_by"`
	CreatedAt            time.Time                        `json:"created_at"`
	UpdatedAt            *time.Time                       `json:"updated_at"`
	CreatedByUser        organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser        *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	Status               *masterDTOs.StatusDTO            `json:"status"`
}
