package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type SLARegisterListDTO struct {
	ID                 		   uint   	   `json:"id"`
	ReferenceNumber            string `json:"reference_number"`
	ResponsibleUserID  		   uint        `json:"responsible_user_id"`
    LifecycleStatusID  		   uint        `json:"lifecycle_status_id"`
    TitleEn            		   string      `json:"title_en"`
    TitleAr            		   string      `json:"title_ar"`
	Description        		   string      `json:"description"` 
    Client             		   string      `json:"client"`
    CreatedBy          		   uint        `json:"created_by"`
    UpdatedBy          		   uint        `json:"updated_by"`
    CreatedAt          		   time.Time   `json:"created_at"`
    UpdatedAt          		   time.Time   `json:"updated_at"`
    SLACategoryTitleEn         string 	   `json:"sla_category_title_en"`
    SLACategoryTitleAr         string 	   `json:"sla_category_title_ar"`
    SeverityLevelTitleEn       string 	   `json:"severity_level_title_en"`
    SeverityLevelTitleAr       string 	   `json:"severity_level_title_ar"`
    ResponsibleUnitTitleEn     string 	   `json:"responsible_unit_title_en"`
    ResponsibleUnitTitleAr     string 	   `json:"responsible_unit_title_ar"`
    ResponsibleUserFirstNameEn string 	   `json:"responsible_user_first_name_en"`
    ResponsibleUserFirstNameAr string 	   `json:"responsible_user_first_name_ar"`
    ResponsibleUserLastNameEn  string 	   `json:"responsible_user_last_name_en"`
    ResponsibleUserLastNameAr  string 	   `json:"responsible_user_last_name_ar"`
    LifecycleStatusTitleEn     string 	   `json:"lifecycle_status_title_en"`
    LifecycleStatusTitleAr     string 	   `json:"lifecycle_status_title_ar"`
    CreatedByFirstNameEn       string 	   `json:"created_by_first_name_en"`
    CreatedByFirstNameAr       string 	   `json:"created_by_first_name_ar"`
    CreatedByLastNameEn        string 	   `json:"created_by_last_name_en"`
    CreatedByLastNameAr        string 	   `json:"created_by_last_name_ar"`
    UpdatedByFirstNameEn       string 	   `json:"updated_by_first_name_en"`
    UpdatedByFirstNameAr       string 	   `json:"updated_by_first_name_ar"`
    UpdatedByLastNameEn        string 	   `json:"updated_by_last_name_en"`
    UpdatedByLastNameAr        string 	   `json:"updated_by_last_name_ar"`
}

type SLARegisterListResponseDTO struct {
	Items      []*SLARegisterListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type SLARegisterFormDTO struct {

    SLACategoryID      uint        `json:"sla_category_id" binding:"required"`
    SeverityLevelID    uint        `json:"severity_level_id" binding:"required"`
    ResponsibleUnitID  uint        `json:"responsible_unit_id" binding:"required"`
    ResponsibleUserID  uint        `json:"responsible_user_id" binding:"required"`
    LifecycleStatusID  uint        `json:"lifecycle_status_id" binding:"required"`
    TitleEn            string      `json:"title_en" binding:"required,max=250"`
    TitleAr            string      `json:"title_ar" binding:"required,max=250"`
	Description        string      `json:"description"` 
    Client             string      `json:"client"`
}
