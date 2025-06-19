package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type TaskSchedularLogListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
}

type TaskSchedularLogListDTO struct {
	ID          uint       `json:"id"`
	TaskName    string     `json:"task_name"`
	Status      string     `json:"status"`
	Message     *string    `json:"message"`
	RunAt       time.Time  `json:"run_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type TaskSchedularLogListResponseDTO struct {
	Items      []*TaskSchedularLogListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
