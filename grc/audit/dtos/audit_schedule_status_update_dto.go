package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditScheduleStatusUpdateListDTO struct {
	ID                           uint                                        `json:"id"`
	AuditScheduleTitle           string                                      `json:"audit_schedule_title"`
	AuditScheduleStatusTitleEn   string                                      `json:"audit_schedule_status_title_en"`
	AuditScheduleStatusTitleAr   string                                      `json:"audit_schedule_status_title_ar"`
	AuditScheduleStatusColorCode string                                      `json:"audit_schedule_status_color_code"`
	Comment                      string                                      `json:"comment"`
	Percentage                   string                                      `json:"percentage"`
	CreatedAt                    string                                      `json:"created_at"`
	UpdatedAt                    string                                      `json:"updated_at"`
	CreatedUserID                uint                                        `json:"created_user_id"`
	CreatedUserFirstNameEn       string                                      `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr       string                                      `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn        string                                      `json:"created_user_last_name_en"`
	CreatedUserLastNameAr        *string                                     `json:"created_user_last_name_ar"`
	CreatedUserInitial           string                                      `json:"created_user_initial"`
	CreatedUserColor             string                                      `json:"created_user_color"`
	CreatedUserEmail             string                                      `json:"created_user_email"`
	CreatedUserProfileImagePath  string                                      `json:"created_user_profile_image_path"`
	UpdatedUserID                *uint                                       `json:"updated_user_id"`
	UpdatedUserFirstNameEn       *string                                     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr       *string                                     `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn        *string                                     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr        *string                                     `json:"updated_user_last_name_ar"`
	UpdatedUserInitial           *string                                     `json:"updated_user_initial"`
	UpdatedUserColor             *string                                     `json:"updated_user_color"`
	UpdatedUserEmail             *string                                     `json:"updated_user_email"`
	UpdatedUserProfileImagePath  *string                                     `json:"updated_user_profile_image_path"`
	Documents                    []*models.AuditScheduleStatusUpdateDocument `json:"documents" gorm:"-"`
}

type AuditScheduleStatusUpdateFormDTO struct {
	AuditScheduleStatusID uint                          `json:"audit_schedule_status_id" binding:"required"`
	Comment               string                        `json:"comment" binding:"required"`
	Percentage            *uint                         `json:"percentage" binding:"required"`
	Documents             []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type AuditScheduleStatusUpdateDTO struct {
	ID                  uint                                    `json:"id"`
	AuditScheduleStatus AuditScheduleStatusDTO                  `json:"audit_schedule_status"`
	Comment             string                                  `json:"comment"`
	Percentage          *uint                                   `json:"percentage"`
	CreatedAt           time.Time                               `json:"created_at"`
	UpdatedAt           *time.Time                              `json:"updated_at"`
	CreatedByUser       organizationDTOs.UserMinimalDTO         `json:"created_by_user"`
	UpdatedByUser       *organizationDTOs.UserMinimalDTO        `json:"updated_by_user"`
	Documents           []*AuditScheduleStatusUpdateDocumentDTO `json:"documents"`
}

type AuditScheduleStatusUpdateDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}
