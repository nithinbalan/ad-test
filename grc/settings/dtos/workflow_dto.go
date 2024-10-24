package dtos

import (
	"time"
)

type WorkflowListDTO struct {
	ID                          uint       `json:"id"`
	Type                        string     `json:"type"`
	TitleEn                     *string    `json:"title_an"`
	TitleAr                     *string    `json:"title_ar"`
	DocumentTypeID              *uint      `json:"document_type_id"`
	Status                      string     `json:"status"`
	CreatedByID                 uint       `json:"created_by_id"`
	UpdatedByID                 *uint      `json:"updated_by_id"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   *time.Time `json:"updated_at"`
	DocumentTypeTitleEn         string     `json:"document_type_title_en"`
	DocumentTypeTitleAr         string     `json:"document_type_title_ar"`
	CreatedUserID               uint       `json:"created_user_id"`
	CreatedUserFirstNameEn      string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string     `json:"created_user_last_name_ar"`
	CreatedUserInitial          string     `json:"created_user_initial"`
	CreatedUserColor            string     `json:"created_user_color"`
	CreatedUserEmail            string     `json:"created_user_email"`
	CreatedUserProfileImagePath string     `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string     `json:"updated_user_initial"`
	UpdatedUserColor            string     `json:"updated_user_color"`
	UpdatedUserEmail            string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath string     `json:"updated_user_profile_image_path"`
}

type WorkflowLevelFormDTO struct {
	ParentLevel  *uint   `json:"parent_level"`
	LevelTypeID  uint    `json:"level_type_id" binding:"required"`
	ApprovalType *string `json:"approval_type"`
	UserIds      []*uint `json:"user_ids"`
	RoleID       *uint   `json:"role_id"`
}

type WorkflowLevelListDTO struct {
	ID             uint               `json:"id"`
	Type           string             `json:"type"`
	DocumentTypeID *uint              `json:"document_type_id"`
	Status         string             `json:"status"`
	CreatedByID    uint               `json:"created_by_id"`
	UpdatedByID    *uint              `json:"updated_by_id"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      *time.Time         `json:"updated_at"`
	WorkflowLevels []WorkflowLevelDTO `json:"workflow_levels"`
}

type WorkflowLevelDTO struct {
	ID                uint                          `json:"id"`
	Level             int                           `json:"level"`
	WorkflowID        uint                          `json:"workflow_id"`
	ApprovalType      string                        `json:"approval_type"`
	CreatedByID       uint                          `json:"created_by_id"`
	UpdatedByID       *uint                         `json:"updated_by_id"`
	CreatedAt         time.Time                     `json:"created_at"`
	UpdatedAt         *time.Time                    `json:"updated_at"`
	WorkflowLevelType WorkflowLevelTypeDTO          `json:"workflow_level_type"`
	Participants      []WorkflowLevelParticipantDTO `json:"participants"`
}

type WorkflowLevelTypeDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Type    string `json:"type"`
}

type WorkflowLevelParticipantDTO struct {
	ID   uint            `json:"id"`
	User *UserMinimalDTO `json:"user,omitempty"`
	Role *RoleDTO        `json:"role,omitempty" `
}

type UserMinimalDTO struct {
	ID               int     `json:"id"`
	FirstNameEN      string  `json:"first_name_en"`
	LastNameEN       string  `json:"last_name_en"`
	FirstNameAR      string  `json:"first_name_ar"`
	LastNameAR       string  `json:"last_name_ar"`
	ProfileImagePath *string `json:"profile_image_path"`
	Email            string  `json:"email"`
	Initial          string  `json:"initial"`
	Color            string  `json:"color"`
}
