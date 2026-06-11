package dtos

import "time"

type GovernanceSettingDTO struct {
	ID                  uint                `json:"id"`
	ApproveDepartmentID uint                `json:"approve_department_id"`
	ApproveDepartment   OrganizationUnitDTO `json:"approve_department"`

	CreatedByID   uint           `json:"created_by_id"`
	CreatedByUser UserMinimalDTO `json:"created_by"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedByID   *uint           `json:"updated_by_id"`
	UpdatedByUser *UserMinimalDTO `json:"updated_by"`

	UpdatedAt *time.Time `json:"updated_at"`
}
