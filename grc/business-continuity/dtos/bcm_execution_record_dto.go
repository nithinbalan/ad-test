package dtos

import (
	mastersDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type BCMExecutionRecordDTO struct {
	ID                 uint       `json:"id"`
	TitleEn            string     `json:"title_en"`
	TitleAr            string     `json:"title_ar"`
	OwnerEntityID      uint       `json:"owner_entity_id"`
	OwnerEntityTitleEn *string    `json:"owner_entity_title_en"`
	OwnerEntityTitleAr *string    `json:"owner_entity_title_ar"`
	RecordTypeID       *uint      `json:"record_type_id"`
	RecordTypeTitleEn  *string    `json:"record_type_title_en"`
	RecordTypeTitleAr  *string    `json:"record_type_title_ar"`
	IsBackUpAvailable  bool       `json:"is_back_up_available"`
	LocationID         *uint      `json:"location_id"`
	LocationTitleEn    *string    `json:"location_title_en"`
	LocationTitleAr    *string    `json:"location_title_ar"`
	RecoveryPlan       *string    `json:"recovery_plan"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
}

type BCMExecutionRecordFormDTO struct {
	TitleEn           string  `json:"title_en" binding:"required,max=1000" gorm:"type:varchar(1000);not null"`
	TitleAr           string  `json:"title_ar" binding:"required,max=1000" gorm:"type:varchar(1000);not null"`
	OwnerEntityID     uint    `json:"owner_entity_id" binding:"required" gorm:"not null"`
	RecordTypeID      *uint   `json:"record_type_id" gorm:"default:null"`
	IsBackUpAvailable bool    `json:"is_back_up_available" gorm:"default:0"`
	LocationID        *uint   `json:"location_id" gorm:"default:null"`
	RecoveryPlan      *string `json:"recovery_plan" gorm:"type:nvarchar(max);default:null"`
}

type BCMExecutionRecordListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
}

type BCMExecutionRecordDetailDTO struct {
	ID                uint                                 `json:"id"`
	TitleEn           string                               `json:"title_en"`
	TitleAr           string                               `json:"title_ar"`
	OwnerEntity       organizationDTOs.OrganizationUnitDTO `json:"owner_entity"`
	RecordType        *BCMExecutionRecordTypeDTO           `json:"record_type"`
	IsBackUpAvailable bool                                 `json:"is_back_up_available"`
	Location          *mastersDTOs.LocationDTO             `json:"location"`
	RecoveryPlan      *string                              `json:"recovery_plan"`
	CreatedAt         time.Time                            `json:"created_at"`
	UpdatedAt         *time.Time                           `json:"updated_at"`
	CreatedByUser     organizationDTOs.UserMinimalDTO      `json:"created_by_user"`
	UpdatedByUser     *organizationDTOs.UserMinimalDTO     `json:"updated_by_user"`
}
