package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/organization/dtos"
	"time"
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
	CreatedAt            time.Time `json:"created_at"`
}

type ActivityLogListResponseDTO struct {
	Items      []*ActivityLogListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
type ActivityLogDetailDTO struct {
	ID          uint                `json:"id"`
	Action      string              `json:"action"`
	TableName   string              `json:"table_name"`
	RecordID    uint                `json:"record_id"`
	ChangedData string              `json:"changed_data"`
	UserID      uint                `json:"user_id"`
	User        dtos.UserMinimalDTO `json:"user"`
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
