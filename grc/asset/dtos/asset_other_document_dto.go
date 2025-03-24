package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AssetOtherDocumentListDTO struct {
	ID                           uint                         `json:"id"`
	AssetRegisterID              uint                         `json:"asset_register_id"`
	AssetDocumentCategoryID      uint                         `json:"asset_document_category_id"`
	AssetDocumentCategoryTitleEn string                       `json:"asset_document_category_title_en"`
	AssetDocumentCategoryTitleAr string                       `json:"asset_document_category_title_ar"`
	Documents                    []*AssetOtherDocumentListDTO `json:"documents"`
	CreatedByID                  uint                         `json:"created_by_id"`
	CreatedByFirstNameEn         string                       `json:"created_by_first_name_en"`
	CreatedByFirstNameAr         string                       `json:"created_by_first_name_ar"`
	CreatedByLastNameEn          string                       `json:"created_by_last_name_en"`
	CreatedByLastNameAr          string                       `json:"created_by_last_name_ar"`
	CreatedByInitial             string                       `json:"created_by_initial"`
	CreatedByColor               string                       `json:"created_by_color"`
	CreatedByEmail               string                       `json:"created_by_email"`
	CreatedByProfileImagePath    *string                      `json:"created_by_profile_image_path"`
	CreatedAt                    time.Time                    `json:"created_at"`
	UpdatedByID                  *uint                        `json:"updated_by_id"`
	UpdatedByFirstNameEn         *string                      `json:"updated_by_first_name_en"`
	UpdatedByFirstNameAr         *string                      `json:"updated_by_first_name_ar"`
	UpdatedByLastNameEn          *string                      `json:"updated_by_last_name_en"`
	UpdatedByLastNameAr          *string                      `json:"updated_by_last_name_ar"`
	UpdatedByInitial             *string                      `json:"updated_by_initial"`
	UpdatedByColor               *string                      `json:"updated_by_color"`
	UpdatedByEmail               *string                      `json:"updated_by_email"`
	UpdatedByProfileImagePath    *string                      `json:"updated_by_profile_image_path"`
	UpdatedAt                    *time.Time                   `json:"updated_at"`
}

type AssetOtherDocumentFormDTO struct {
	AssetDocumentCategoryID uint                          `json:"asset_document_category_id" binding:"required"`
	Documents               []*commonDTOs.DocumentFormDTO `json:"documents" binding:"required"` // At least one document expected
}

type AssetOtherDocumentDetailDTO struct {
	ID                    uint                             `json:"id"`
	AssetRegisterID       uint                             `json:"asset_register_id"`
	AssetDocumentCategory AssetDocumentCategoryDTO         `json:"asset_document_category"`
	URL                   string                           `json:"url"`
	Name                  string                           `json:"name"`
	Size                  int                              `json:"size"`
	Token                 string                           `json:"token"`
	OriginalName          string                           `json:"original_name"`
	Type                  string                           `json:"type"`
	CreatedAt             time.Time                        `json:"created_at"`
	UpdatedAt             *time.Time                       `json:"updated_at"`
	CreatedByUser         organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser         *organizationDTOs.UserMinimalDTO `json:"updated_by"`
}
