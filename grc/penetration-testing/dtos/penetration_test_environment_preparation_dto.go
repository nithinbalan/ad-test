package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type PenetrationTestEnvironmentPreparationFormDTO struct {
	EnvironmentNotifyDate  string `json:"environment_notify_date" binding:"required"`
	ExpectedReadyDate      string `json:"expected_ready_date" binding:"required"`
	EnvironmentReadinessID uint   `json:"environment_readiness_id" binding:"required"`
}

type PenetrationTestEnvironmentPreparationDetailDTO struct {
	ID                          uint                                    `json:"id"`
	PenetrationTestID           uint                                    `json:"penetration_test_id"`
	EnvironmentNotifyDateString *string                                 `json:"environment_notify_date"`
	ExpectedReadyDateString     *string                                 `json:"expected_ready_date"`
	EnvironmentReadiness        *PenetrationTestEnvironmentReadinessDTO `json:"environment_readiness"`
	CreatedAt                   time.Time                               `json:"created_at"`
	UpdatedAt                   *time.Time                              `json:"updated_at"`
	CreatedByUser               organizationDTOs.UserMinimalDTO         `json:"created_by_user"`
	UpdatedByUser               *organizationDTOs.UserMinimalDTO        `json:"updated_by_user"`
}
