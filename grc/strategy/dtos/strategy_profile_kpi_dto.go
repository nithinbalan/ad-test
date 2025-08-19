package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type StrategyProfileKpiListDTO struct {
	ID                                       int      `json:"id"`
	StrategyProfileID                        uint     `json:"strategy_profile_id"`
	StrategyProfileTitle                     string   `json:"strategy_profile_title"`
	StrategyProfileStrategicObjectiveID      *uint    `json:"strategy_profile_strategic_objective_id"`
	StrategyProfileStrategicObjectiveTitleEn *string  `json:"strategy_profile_strategic_objective_title_en"`
	StrategyProfileStrategicObjectiveTitleAr *string  `json:"strategy_profile_strategic_objective_title_ar"`
	KpiID                                    *uint    `json:"kpi_id"`
	KpiTitle                                 *string  `json:"kpi_title"`
	Weightage                                *float64 `json:"weightage"`
	CreatedUserID                            uint     `json:"created_user_id"`
	CreatedUserFirstNameEn                   string   `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr                   string   `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                    string   `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                    string   `json:"created_user_last_name_ar"`
	CreatedUserInitial                       string   `json:"created_user_initial"`
	CreatedUserColor                         string   `json:"created_user_color"`
	CreatedUserEmail                         string   `json:"created_user_email"`
	CreatedUserProfileImagePath              string   `json:"created_user_profile_image_path"`
	UpdatedUserID                            *uint    `json:"updated_user_id"`
	UpdatedUserFirstNameEn                   *string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr                   *string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                    *string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                    *string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                       *string  `json:"updated_user_initial"`
	UpdatedUserColor                         *string  `json:"updated_user_color"`
	UpdatedUserEmail                         *string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath              *string  `json:"updated_user_profile_image_path"`
}

type StrategyProfileKpiFormDTO struct {
	StrategyProfileStrategicObjectiveID *uint  `json:"strategy_profile_strategic_objective_id"`
	KpiIDs                              []uint `json:"kpi_ids" binding:"required,min=1"`
}

type StrategyProfileKpiListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type KpiRegisterMinimalDetails struct {
	ID          uint     `json:"id"`
	Achievement *float64 `json:"achievement,omitempty"`
	Weightage   *float64 `json:"weightage,omitempty"`
}

type StrategyProfileKpiDetailDTO struct {
	ID                                  uint                                        `json:"id"`
	StrategyProfileStrategicObjectiveID *uint                                       `json:"strategy_profile_strategic_objective_id"`
	StrategyProfileStrategicObjective   *StrategyProfileStrategicObjectiveDetailDTO `json:"strategy_profile_strategic_objective"`
	KpiID                               *uint                                       `json:"kpi_id"`
	Kpi                                 *KPIRegisterDetailDTO                       `json:"kpi"`
	Weightage                           *float64                                    `json:"weightage"`
	CreatedByID                         uint                                        `json:"created_by_id"`
	CreatedByUser                       organizationDTOs.UserMinimalDTO             `json:"created_by_user"`
	UpdatedByID                         *uint                                       `json:"updated_by_id"`
	UpdatedByUser                       *organizationDTOs.UserMinimalDTO            `json:"updated_by_user"`
	CreatedAt                           time.Time                                   `json:"created_at"`
	UpdatedAt                           *time.Time                                  `json:"updated_at"`
}

type StrategyProfileKpiWeightageUpdateDTO struct {
	Weightage *float64 `json:"weightage" binding:"required,min=0,max=100"`
}
