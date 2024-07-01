package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type IncidentPostReviewListDTO struct {
	ID                     uint      `json:"id"`
	IncidentRegisterID     uint      `json:"incident_register_id"`
	PostReviewDate         time.Time `json:"post_review_date"`
	Summary                *string   `json:"summary"`
	CreatedAt              string    `json:"created_at"`
	UpdatedAt              string    `json:"updated_at"`
	CreatedUserFirstNameEn string    `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string    `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string    `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  string    `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn  string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr  string    `json:"updated_user_last_name_ar"`
}
type IncidentPostReviewListQueryParamsDTO struct {
	Sort           string
	SortBy         string
	Keyword        string
	Status         string
	Page           int
	PageSize       int
	Filters        map[string][]string
	QuickFilter    []string
	PostReviewDate string
}

type IncidentPostReviewListResponseDTO struct {
	Items      []*IncidentPostReviewListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type IncidentPostReviewFormDTO struct {
	IncidentRegisterID                uint      `json:"incident_register_id"`
	IncidentPostReviewEffectivenessID uint      `json:"incident_post_review_effectiveness_id"`
	PostReviewDate                    string    `json:"post_review_date" binding:"required,datetime=2006-01-02"`
	Summary                           *string   `json:"summary"`
	UserIds                           []*uint   `json:"user_ids"`
	Recommendations                   []*string `json:"recommendations" binding:"required,dive,max=500" gorm:"type:varchar(500);not null"`
	Lessons                           []*string `json:"lessons" binding:"required,dive,max=500" gorm:"type:varchar(500);not null"`
}

type IncidentPostReviewDetailDTO struct {
	ID                              uint                                   `json:"id"`
	IncidentRegister                IncidentRegisterDetailDTO              `json:"incident_register"`
	IncidentPostReviewEffectiveness IncidentPostReviewEffectivenessDTO     `json:"incident_post_review_effectiveness_id"`
	Summary                         *string                                `json:"summary"`
	PostReviewDate                  time.Time                              `json:"post_review_date"`
	Recommendations                 []*IncidentPostReviewRecommendationDTO `json:"recommendations"`
	Users                           []*IncidentPostReviewUserDTO           `json:"users"`
	Lessons                         []*IncidentPostReviewLessonDTO         `json:"lessons"`
	CreatedAt                       time.Time                              `json:"created_at"`
	UpdatedAt                       *time.Time                             `json:"updated_at"`
	CreatedByUser                   organizationDTOs.UserMinimalDTO        `json:"created_by"`
	UpdatedByUser                   organizationDTOs.UserMinimalDTO        `json:"updated_by"`
}

type IncidentPostReviewParsedDates struct {
	PostReviewDate *time.Time
}

type IncidentPostReviewLessonDTO struct {
	ID                   uint   `json:"id"`
	IncidentPostReviewID uint   `json:"incident_post_review_id"`
	Lesson               string `json:"lesson"`
}

type IncidentPostReviewUserDTO struct {
	ID                   uint                            `json:"id"`
	IncidentPostReviewID int                             `json:"incident_post_review_id"`
	UserID               int                             `json:"user_id"`
	User                 organizationDTOs.UserDetailsDTO `json:"user"`
}

type IncidentPostReviewRecommendationDTO struct {
	ID                   uint   `json:"id"`
	IncidentPostReviewID uint   `json:"incident_post_review_id"`
	Recommendation       string `json:"recommendation"`
}
