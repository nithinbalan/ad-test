package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type PenetrationTestManagementSharingFormDTO struct {
	ResultsSharedWithManagement bool    `json:"results_shared_with_management"`
	ManagementNotifyDate        *string `json:"management_notify_date,omitempty"`
	ManagementResponseDate      *string `json:"management_response_date,omitempty"`
	RemediationNotes            *string `json:"remediation_notes,omitempty"`
}

type PenetrationTestManagementSharingDetailDTO struct {
	ID                           uint                             `json:"id"`
	PenetrationTestID            uint                             `json:"penetration_test_id"`
	ResultsSharedWithManagement  bool                             `json:"results_shared_with_management"`
	ManagementNotifyDateString   *string                          `json:"management_notify_date"`
	ManagementResponseDateString *string                          `json:"management_response_date"`
	RemediationNotes             *string                          `json:"remediation_notes"`
	CreatedAt                    time.Time                        `json:"created_at"`
	UpdatedAt                    *time.Time                       `json:"updated_at"`
	CreatedByUser                organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser                *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}
