package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditPlanCommencementLetterListDTO struct {
	ID                          uint    `json:"id"`
	AuditPlanID                 uint    `json:"audit_plan_id"`
	ResponsibleUnitTitleEn      string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr      string  `json:"responsible_unit_title_ar"`
	ToEmailIDs                  *string `json:"to_email_ids"`
	FromUserID                  uint    `json:"from_user_id"`
	FromUserFirstNameEn         string  `json:"from_user_first_name_en"`
	FromUserFirstNameAr         string  `json:"from_user_first_name_ar"`
	FromUserLastNameEn          string  `json:"from_user_last_name_en"`
	FromUserLastNameAr          string  `json:"from_user_last_name_ar"`
	FromUserInitial             string  `json:"from_user_initial"`
	FromUserColor               string  `json:"from_user_color"`
	FromUserEmail               string  `json:"from_user_email"`
	FromUserProfileImagePath    string  `json:"from_user_profile_image_path"`
	ToUserEmails                string  `json:"to_user_emails"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   *string `json:"updated_at"`
	CreatedUserID               uint    `json:"created_user_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
}

type AuditPlanCommencementLetterFormDTO struct {
	SubjectEn  string `json:"subject_en" binding:"required"`
	SubjectAr  string `json:"subject_ar" binding:"required"`
	BodyEn     string `json:"body_en" binding:"required"`
	BodyAr     string `json:"body_ar" binding:"required"`
	FromUserID uint   `json:"from_user_id" binding:"required"`
}

type AuditPlanCommencementLetterDetailDTO struct {
	ID            uint                                    `json:"id"`
	AuditPlanID   uint                                    `json:"audit_plan_id"`
	ToEmailIDs    string                                  `json:"to_email_ids"`
	SubjectEn     string                                  `json:"subject_en"`
	SubjectAr     string                                  `json:"subject_ar"`
	BodyEn        string                                  `json:"body_en"`
	BodyAr        string                                  `json:"body_ar"`
	CreatedAt     time.Time                               `json:"created_at"`
	UpdatedAt     *time.Time                              `json:"updated_at,omitempty"`
	CreatedByUser organizationDTOs.UserMinimalDTO         `json:"created_by"`
	UpdatedByUser *organizationDTOs.UserMinimalDTO        `json:"updated_by"`
	AuditPlan     *AuditPlanDetailDTO                     `json:"audit_plan"`
	ToUser        AuditPlanCommencementLetterToUserDTO    `json:"to_user"`
	CCUsers       *[]AuditPlanCommencementLetterCCUserDTO `json:"cc_users"`
}

type AuditPlanCommencementLetterToUserDTO struct {
	ID     uint                            `json:"id"`
	UserID uint                            `json:"user_id"`
	User   organizationDTOs.UserMinimalDTO `json:"user"`
}

type AuditPlanCommencementLetterCCUserDTO struct {
	ID     uint                            `json:"id"`
	UserID uint                            `json:"user_id"`
	User   organizationDTOs.UserMinimalDTO `json:"user"`
}

type AuditPlanCommencementLetterSendFormDTO struct {
	ToUserID  uint    `json:"to_user_id" binding:"required"`
	CCUserIDs []*uint `json:"cc_user_ids"`
}
