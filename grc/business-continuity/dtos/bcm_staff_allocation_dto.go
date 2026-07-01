package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
)

type BCMStaffAllocationDetailDTO struct {
	ID                      uint                                               `json:"id"`
	ProcessID               uint                                               `json:"process_id"`
	ProcessEmployeeAnalysis *organizationDTOs.ProcessEmployeeAnalysisDetailDTO `json:"process_employee_analysis"`
	OutageBuckets           []BCMStaffAllocationOutageBucketDetailDTO          `json:"outage_buckets"`
}

type BCMStaffAllocationOutageBucketDetailDTO struct {
	ID             uint                        `json:"id"`
	OutageBucketID uint                        `json:"outage_bucket_id"`
	OutageBucket   BCMBIAOutageBucketDTO       `json:"outage_bucket"`
	IsSkipped      bool                        `json:"is_skipped"`
	Users          []BCMStaffAllocationUserDTO `json:"users"`
}

type BCMStaffAllocationUserDTO struct {
	ID       uint                             `json:"id"`
	UserID   uint                             `json:"user_id"`
	User     *organizationDTOs.UserMinimalDTO `json:"user"`
	IsRemote bool                             `json:"is_remote"`
}

type BCMStaffAllocationUpdateFormDTO struct {
	BCMStaffAllocationOutageBucketID uint                              `json:"bcm_staff_allocation_outage_bucket_id" binding:"required"`
	Users                            []BCMStaffAllocationUpdateUserDTO `json:"users" binding:"required"`
}

type BCMStaffAllocationUpdateUserDTO struct {
	UserID   uint `json:"user_id" binding:"required"`
	IsRemote bool `json:"is_remote"`
}
