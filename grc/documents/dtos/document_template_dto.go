package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type DocumentTemplateListDTO struct {
	ID                          int    `json:"id"`
	TitleEn                     string `json:"title_en"`
	TitleAr                     string `json:"title_ar"`
	DocumentTypeTitleEn         string `json:"document_type_title_en"`
	DocumentTypeTitleAr         string `json:"document_type_title_ar"`
	Status                      string `json:"status"`
	URL                         string `json:"url"`
	Name                        string `json:"file_name"`
	OriginalName                string `json:"original_name"`
	Type                        string `json:"type"`
	Size                        int    `json:"size"`
	CreatedUserID               uint   `json:"created_user_id"`
	CreatedUserFirstNameEn      string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string `json:"created_user_last_name_ar"`
	CreatedUserInitial          string `json:"created_user_initial"`
	CreatedUserColor            string `json:"created_user_color"`
	CreatedUserEmail            string `json:"created_user_email"`
	CreatedUserProfileImagePath string `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint  `json:"updated_user_id"`
	UpdatedUserFirstNameEn      string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string `json:"updated_user_initial"`
	UpdatedUserColor            string `json:"updated_user_color"`
	UpdatedUserEmail            string `json:"updated_user_email"`
}

type DocumentTemplateFormDTO struct {
	TitleEn        string                      `json:"title_en" binding:"required"`
	TitleAr        string                      `json:"title_ar" binding:"required"`
	DocumentTypeID uint                        `json:"document_type_id" binding:"required"`
	Status         string                      `json:"status" binding:"required,oneof=active inactive"`
	Document       *commonDTOs.DocumentFormDTO `json:"document" binding:"required"`
}

type DocumentTemplateDetailsDTO struct {
	ID            int                              `json:"id"`
	TitleEn       string                           `json:"title_en"`
	TitleAr       string                           `json:"title_ar"`
	DocumentType  DocumentTypeDTO                  `json:"document_type"`
	Status        string                           `json:"status"`
	Document      *DocumentTemplateDocumentDTO     `json:"document"`
	CreatedByUser *organizationDTOs.UserMinimalDTO `json:"created_by_user"`
	UpdatedByUser *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
	CreatedAt     time.Time                        `json:"created_at"`
	UpdatedAt     *time.Time                       `json:"updated_at"`
}

type DocumentTemplateDocumentDTO struct {
	ID           uint   `json:"id,omitempty"`
	URL          string `json:"url"`
	Name         string `json:"file_name"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
	Size         int    `json:"size"`
}
