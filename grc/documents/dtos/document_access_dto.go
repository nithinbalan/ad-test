package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type DocumentSharedUserListResponseDTO struct {
	Items      []*organizationDTOs.UserMinimalDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO       `json:"pagination"`
}
