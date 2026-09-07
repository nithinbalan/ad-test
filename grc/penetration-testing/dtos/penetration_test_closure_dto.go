package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type PenetrationTestClosureFormDTO struct {
	ClosureOwnerID uint   `json:"closure_owner_id" binding:"required"`
	ClosureDate    string `json:"closure_date" binding:"required"`
	ClosureRemarks string `json:"closure_remarks" binding:"required"`
}

type PenetrationTestClosureDetailDTO struct {
	ID                uint                             `json:"id"`
	PenetrationTestID uint                             `json:"penetration_test_id"`
	ClosureDateString *string                          `json:"closure_date"`
	ClosureRemarks    string                           `json:"closure_remarks"`
	ClosureOwner      organizationDTOs.UserMinimalDTO  `json:"closure_owner"`
	CreatedAt         time.Time                        `json:"created_at"`
	UpdatedAt         *time.Time                       `json:"updated_at"`
	CreatedByUser     organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser     *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}
