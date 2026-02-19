package dtos

import (
	"time"

	organizationDTOs "jethur/grc/organization/dtos"
)

type ControlConsistencyResponseDTO struct {
	ID            uint                             `json:"id"`
	TitleEn       string                           `json:"title_en"`
	TitleAr       string                           `json:"title_ar"`
	Type          string                           `json:"type"`
	Score         float64                          `json:"score"`
	Status        string                           `json:"status"`
	DescriptionEn *string                          `json:"description_en"`
	DescriptionAr *string                          `json:"description_ar"`
	CreatedAt     time.Time                        `json:"created_at"`
	CreatedByID   uint                             `json:"created_by_id"`
	CreatedByUser *organizationDTOs.UserMinimalDTO `json:"created_by_user"`
}
