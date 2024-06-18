package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type IncidentCommunicationLogListDTO struct {
	ID                  uint   `json:"id"`
	IncidentRegisterID  uint   `json:"incident_register_id"`
	CommunicatedAt      string `json:"communicated_at"`
	CommunicatedBy      string `json:"communicated_by"`
	CommunicatedTo      string `json:"communicated_to"`
	CommunicationMethod string `json:"communication_method"`
	Summary             string `json:"summary"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

type IncidentCommunicationLogListQueryParamsDTO struct {
	Sort           string
	SortBy         string
	Keyword        string
	Status         string
	Page           int
	PageSize       int
	Filters        map[string][]string
	QuickFilter    []string
	CommunicatedAt string
}

type IncidentCommunicationLogListResponseDTO struct {
	Items      []*IncidentCommunicationLogListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO       `json:"pagination"`
}

type IncidentCommunicationLogFormDTO struct {
	CommunicatedAt      string `json:"communicated_at" binding:"required,datetime=2006-01-02T15:04:05"`
	CommunicatedBy      string `json:"communicated_by" binding:"required"`
	CommunicatedTo      string `json:"communicated_to" binding:"required"`
	CommunicationMethod string `json:"communication_method" binding:"required"`
	Summary             string `json:"summary" binding:"omitempty"`
}

type IncidentCommunicationLogParsedDates struct {
	CommunicatedAt *time.Time
}

type IncidentCommunicationLogDetailDTO struct {
	ID                  int                             `json:"id"`
	CommunicatedAt      time.Time                       `json:"communicated_at"`
	CommunicatedBy      string                          `json:"communicated_by"`
	CommunicatedTo      string                          `json:"communicated_to"`
	CommunicationMethod string                          `json:"communication_method"`
	Summary             *string                         `json:"summary"`
	IncidentRegister    IncidentRegisterDetailDTO       `json:"incident_register"`
	CreatedAt           time.Time                       `json:"created_at"`
	UpdatedAt           *time.Time                      `json:"updated_at"`
	CreatedByUser       organizationDTOs.UserMinimalDTO `json:"created_by"`
	UpdatedByUser       organizationDTOs.UserMinimalDTO `json:"updated_by"`
}
