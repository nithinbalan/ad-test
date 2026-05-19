package dtos

import (
	mastersDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type BCMManualProcedureFormDTO struct {
	ID                    uint       `json:"id"`
	TitleEn               string     `json:"title_en"`
	TitleAr               string     `json:"title_ar"`
	Code                  *string    `json:"code"`
	OwnerID               uint       `json:"owner_id"`
	OwnerFirstNameEn      *string    `json:"owner_first_name_en"`
	OwnerLastNameEn       *string    `json:"owner_last_name_en"`
	OwnerFirstNameAr      *string    `json:"owner_first_name_ar"`
	OwnerLastNameAr       *string    `json:"owner_last_name_ar"`
	OwnerEmail            *string    `json:"owner_email"`
	OwnerInitial          *string    `json:"owner_initial"`
	OwnerColor            *string    `json:"owner_color"`
	OwnerProfileImagePath *string    `json:"owner_profile_image_path"`
	LocationID            *uint      `json:"location_id"`
	LocationTitleEn       *string    `json:"location_title_en"`
	LocationTitleAr       *string    `json:"location_title_ar"`
	StakeholderID         *uint      `json:"stakeholder_id"`
	StakeholderTitleEn    *string    `json:"stakeholder_title_en"`
	StakeholderTitleAr    *string    `json:"stakeholder_title_ar"`
	Notes                 *string    `json:"notes"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             *time.Time `json:"updated_at"`
}

type BCMManualProcedureFormFormDTO struct {
	TitleEn       string  `json:"title_en" binding:"required,max=1000" gorm:"type:varchar(1000);not null"`
	TitleAr       string  `json:"title_ar" binding:"required,max=1000" gorm:"type:varchar(1000);not null"`
	Code          *string `json:"code" gorm:"type:varchar(1000);default:null"`
	OwnerID       uint    `json:"owner_id" binding:"required" gorm:"not null"`
	LocationID    *uint   `json:"location_id" gorm:"default:null"`
	StakeholderID *uint   `json:"stakeholder_id" gorm:"default:null"`
	Notes         *string `json:"notes" gorm:"type:nvarchar(max);default:null"`
}

type BCMManualProcedureFormListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
}

type BCMManualProcedureFormDetailDTO struct {
	ID            uint                             `json:"id"`
	TitleEn       string                           `json:"title_en"`
	TitleAr       string                           `json:"title_ar"`
	Code          *string                          `json:"code"`
	Owner         organizationDTOs.UserMinimalDTO  `json:"owner"`
	Location      *mastersDTOs.LocationDTO         `json:"location"`
	Stakeholder   *StakeholderMinimalDTO           `json:"stakeholder"`
	Notes         *string                          `json:"notes"`
	CreatedBYID   uint                             `json:"created_by_id"`
	UpdatedBYID   *uint                            `json:"updated_by_id"`
	CreatedAt     time.Time                        `json:"created_at"`
	UpdatedAt     *time.Time                       `json:"updated_at"`
	CreatedByUser organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}

type StakeholderMinimalDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}
