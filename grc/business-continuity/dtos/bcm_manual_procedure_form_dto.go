package dtos

import (
	mastersDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type BCMManualProcedureFormDTO struct {
	ID                  uint       `json:"id"`
	TitleEn             string     `json:"title_en"`
	TitleAr             string     `json:"title_ar"`
	Code                *string    `json:"code"`
	OwnerEntityID       uint       `json:"owner_entity_id"`
	OwnerEntityTitleEn  *string    `json:"owner_entity_title_en"`
	OwnerEntityTitleAr  *string    `json:"owner_entity_title_ar"`
	LocationID          *uint      `json:"location_id"`
	LocationTitleEn     *string    `json:"location_title_en"`
	LocationTitleAr     *string    `json:"location_title_ar"`
	InternalStakeholder *string    `json:"internal_stakeholder"`
	Notes               *string    `json:"notes"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
}

type BCMManualProcedureFormFormDTO struct {
	TitleEn             string  `json:"title_en" binding:"required,max=1000" gorm:"type:varchar(1000);not null"`
	TitleAr             string  `json:"title_ar" binding:"required,max=1000" gorm:"type:varchar(1000);not null"`
	Code                *string `json:"code" gorm:"type:varchar(1000);default:null"`
	OwnerEntityID       uint    `json:"owner_entity_id" binding:"required" gorm:"not null"`
	LocationID          *uint   `json:"location_id" gorm:"default:null"`
	InternalStakeholder *string `json:"internal_stakeholder" gorm:"type:nvarchar(max);default:null"`
	Notes               *string `json:"notes" gorm:"type:nvarchar(max);default:null"`
}

type BCMManualProcedureFormListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
}

type BCMManualProcedureFormDetailDTO struct {
	ID                  uint                                 `json:"id"`
	TitleEn             string                               `json:"title_en"`
	TitleAr             string                               `json:"title_ar"`
	Code                *string                              `json:"code"`
	OwnerEntity         organizationDTOs.OrganizationUnitDTO `json:"owner_entity"`
	Location            *mastersDTOs.LocationDTO             `json:"location"`
	InternalStakeholder *string                              `json:"internal_stakeholder"`
	Notes               *string                              `json:"notes"`
	CreatedBYID         uint                                 `json:"created_by_id"`
	UpdatedBYID         *uint                                `json:"updated_by_id"`
	CreatedAt           time.Time                            `json:"created_at"`
	UpdatedAt           *time.Time                           `json:"updated_at"`
	CreatedByUser       organizationDTOs.UserMinimalDTO      `json:"created_by_user"`
	UpdatedByUser       *organizationDTOs.UserMinimalDTO     `json:"updated_by_user"`
}
