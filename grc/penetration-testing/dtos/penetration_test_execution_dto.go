package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type PenetrationTestExecutionFormDTO struct {
	TestStartDate     string `json:"test_start_date" binding:"required"`
	TestEndDate       string `json:"test_end_date" binding:"required"`
	ExecutionStatusID uint   `json:"execution_status_id" binding:"required"`
	SecurityResultID  uint   `json:"security_result_id" binding:"required"`
}

type PenetrationTestExecutionDetailDTO struct {
	ID                  uint                               `json:"id"`
	PenetrationTestID   uint                               `json:"penetration_test_id"`
	TestStartDateString *string                            `json:"test_start_date"`
	TestEndDateString   *string                            `json:"test_end_date"`
	ExecutionStatus     *PenetrationTestExecutionStatusDTO `json:"execution_status"`
	SecurityResult      *PenetrationTestSecurityResultDTO  `json:"security_result"`
	CreatedAt           time.Time                          `json:"created_at"`
	UpdatedAt           *time.Time                         `json:"updated_at"`
	CreatedByUser       organizationDTOs.UserMinimalDTO    `json:"created_by_user"`
	UpdatedByUser       *organizationDTOs.UserMinimalDTO   `json:"updated_by_user"`
}
