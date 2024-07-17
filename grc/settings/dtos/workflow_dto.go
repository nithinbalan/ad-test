package dtos

import (
	"time"
)

type WorkflowListDTO struct {
	ID                     uint       `json:"id"`
	Type                   string     `json:"type"`
	DocumentTypeID         *uint      `json:"document_type_id"`
	Status                 string     `json:"status"`
	CreatedByID            uint       `json:"created_by_id"`
	UpdatedByID            *uint      `json:"updated_by_id"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              *time.Time `json:"updated_at"`
	DocumentTypeTitleEn    string     `json:"document_type_title_en"`
	DocumentTypeTitleAr    string     `json:"document_type_title_ar"`
	CreatedUserFirstNameEn string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string     `json:"created_user_first_name_ar"`
	CreatedUserEmail       string     `json:"created_user_email"`
	UpdatedUserFirstNameEn string     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr string     `json:"updated_user_first_name_ar"`
	UpdatedUserEmail       string     `json:"updated_user_email"`
}

type WorkflowLevelFormDTO struct {
	ParantLevel  *uint   `json:"parant_level"`
	WorkflowID   uint    `json:"work_flow_id"`
	RoleTypeID   uint    `json:"role_type_id" binding:"required"`
	ApprovalType *string `json:"approval_type"`
	UserIDs      []*uint `json:"user_ids"`
	RoleID       *uint   `json:"role_id"`
}

type WorkflowLevelListDTO struct {
	WorkflowLevels []WorkflowLevelDTO `json:"workflow_levels"`
}

type WorkflowLevelDTO struct {
	ID               uint                          `json:"id"`
	Level            int                           `json:"level"`
	WorkflowID       uint                          `json:"workflow_id"`
	ApprovalType     string                        `json:"approval_type"`
	CreatedByID      uint                          `json:"created_by_id"`
	UpdatedByID      *uint                         `json:"updated_by_id"`
	CreatedAt        time.Time                     `json:"created_at"`
	UpdatedAt        *time.Time                    `json:"updated_at"`
	WorkflowRoleType WorkflowRoleTypeDTO           `json:"workflow_role_type"`
	Participants     []WorkflowLevelParticipantDTO `json:"participants"`
}

type WorkflowRoleTypeDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Type    string `json:"type"`
}

type WorkflowLevelParticipantDTO struct {
	ID   uint            `json:"id"`
	User *UserMinimalDTO `json:"user"`
	Role *RoleDTO        `json:"role"`
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
