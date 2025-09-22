package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/organization/dtos"
)

type ActivityLogListDTO struct {
	ID                   uint      `json:"id"`
	Action               string    `json:"action"`
	TableName            string    `json:"table_name"`
	RecordID             uint      `json:"record_id"`
	ChangedData          string    `json:"changed_data"`
	UserID               uint      `json:"user_id"`
	UserFirstNameEn      string    `json:"user_first_name_en"`
	UserFirstNameAr      string    `json:"user_first_name_ar"`
	UserLastNameEn       string    `json:"user_last_name_en"`
	UserLastNameAr       string    `json:"user_last_name_ar"`
	UserEmail            string    `json:"user_email"`
	UserProfileImagePath string    `json:"user_profile_image_path"`
	UserInitial          string    `json:"user_initial"`
	UserColor            string    `json:"user_color"`
	ModuleID             *uint     `json:"module_id"`
	ModuleIdentifier     *string   `json:"module_identifier"`
	ModuleTitleEn        *string   `json:"module_title_en"`
	ModuleTitleAr        *string   `json:"module_title_ar"`
	SubModuleID          *uint     `json:"sub_module_id"`
	SubModuleIdentifier  *string   `json:"sub_module_identifier"`
	SubModuleTitleEn     *string   `json:"sub_module_title_en"`
	SubModuleTitleAr     *string   `json:"sub_module_title_ar"`
	Section              *string   `json:"section"`
	CreatedAt            time.Time `json:"created_at"`
}

type ActivityLogListResponseDTO struct {
	Items      []*ActivityLogListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ModuleMinimalDTO struct {
	ID         uint   `json:"id"`
	Identifier string `json:"identifier"`
	TitleEn    string `json:"title_en"`
	TitleAr    string `json:"title_ar"`
}

type ActivityLogDetailDTO struct {
	ID          uint                `json:"id"`
	Action      string              `json:"action"`
	TableName   string              `json:"table_name"`
	RecordID    uint                `json:"record_id"`
	ChangedData string              `json:"changed_data"`
	UserID      uint                `json:"user_id"`
	User        dtos.UserMinimalDTO `json:"user"`
	ModuleID    *uint               `json:"module_id"`
	Module      *ModuleMinimalDTO   `json:"module"`
	SubModuleID *uint               `json:"sub_module_id"`
	SubModule   *ModuleMinimalDTO   `json:"sub_module"`
	Section     *string             `json:"section"`
	CreatedAt   time.Time           `json:"created_at"`
}

type ActivityLogListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}
