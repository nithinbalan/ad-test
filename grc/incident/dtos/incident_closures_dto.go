package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type IncidentClosureListDTO struct {
	ID                 uint       `json:"id"`
	IncidentRegisterID uint       `json:"incident_register_id"`
	CompletionDate     string     `json:"completion_date"`
	ResolutionSummary  string     `json:"resolution_summary"`
	Notes              string     `json:"notes"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
}

type IncidentClosureListQueryParamsDTO struct {
	Sort           string
	SortBy         string
	Keyword        string
	Status         string
	Page           int
	PageSize       int
	Filters        map[string][]string
	QuickFilter    []string
	CompletionDate string
}

type IncidentClosureListResponseDTO struct {
	Items      []*IncidentClosureListDTO    `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type IncidentClosureFormDTO struct {
	CompletionDate    string    `json:"completion_date" binding:"required,datetime=2006-01-02"` // Date format: "2006-01-02"
	ResolutionSummary *string   `json:"resolution_summary" binding:"omitempty"`
	Notes             *string   `json:"notes" binding:"omitempty"`
	Lessons           []*string `json:"lessons" binding:"required,dive,max=500" gorm:"type:varchar(500);not null"`
}

type IncidentClosureParsedDates struct {
	CompletionDate *time.Time
}

type IncidentClosureDetailDTO struct {
	ID                   int                             `json:"id"`
	IncidentRegister     IncidentRegisterDetailDTO       `json:"incident_register"`
	CompletionDateString *string                         `json:"completion_date"`
	ResolutionSummary    *string                         `json:"resolution_summary"`
	Notes                *string                         `json:"notes"`
	Lessons              []*IncidentClosureLessonDTO     `json:"lessons"`
	CreatedAt            time.Time                       `json:"created_at"`
	UpdatedAt            *time.Time                      `json:"updated_at"`
	CreatedByUser        organizationDTOs.UserMinimalDTO `json:"created_by"`
	UpdatedByUser        organizationDTOs.UserMinimalDTO `json:"updated_by"`
}
