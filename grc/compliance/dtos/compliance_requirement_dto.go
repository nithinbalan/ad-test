package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type ComplianceRequirementListDTO struct {
	ID                                   uint   `json:"id"`
	ClauseID                             uint   `json:"clause_id"`
	ComplianceClauseNo                   string `json:"compliance_clause_no"`
	ComplianceClauseTitle                string `json:"compliance_clause_title"`
	Title                                string `json:"title"`
	ComplianceRequirementCategoryTitleEn string `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr string `json:"compliance_requirement_category_title_ar"`
	SeverityLevelTitleEn                 string `json:"severity_level_title_en"`
	SeverityLevelTitleAr                 string `json:"severity_level_title_ar"`
	SeverityLevelColorCode               string `json:"severity_level_color_code"`
	ResponsibleUserFirstNameEn           string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr           string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn            string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr            string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial               string `json:"responsible_user_initial"`
	ResponsibleUserColor                 string `json:"responsible_user_color"`
	ResponsibleUserEmail                 string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath      string `json:"responsible_user_profile_image_path"`
	Deadline                             string `json:"deadline"`
	ComplianceStatusTitleEn              string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr              string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode            string `json:"compliance_status_color_code"`
	CreatedAt                            string `json:"created_at"`
	UpdatedAt                            string `json:"updated_at"`
	CreatedUserFirstNameEn               string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr               string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                string `json:"created_user_last_name_ar"`
	CreatedUserInitial                   string `json:"created_user_initial"`
	CreatedUserColor                     string `json:"created_user_color"`
	CreatedUserEmail                     string `json:"created_user_email"`
	CreatedUserProfileImagePath          string `json:"created_user_profile_image_path"`
	UpdatedUserFirstNameEn               string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr               string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                   string `json:"updated_user_initial"`
	UpdatedUserColor                     string `json:"updated_user_color"`
	UpdatedUserEmail                     string `json:"updated_user_email"`
	UpdatedUserProfileImagePath          string `json:"updated_user_profile_image_path"`
}

type ComplianceRequirementListQueryParamsDTO struct {
	Sort         string
	SortBy       string
	Keyword      string
	Status       string
	Page         int
	PageSize     int
	Filters      map[string][]string
	QuickFilter  []string
	DeadlineFrom string
	DeadlineTo   string
}

type ComplianceRequirementListResponseDTO struct {
	Items      []*ComplianceRequirementListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO    `json:"pagination"`
}

type ComplianceRequirementByComplianceIDListDTO struct {
	ClauseID     uint                            `json:"clause_id"`
	ClauseNo     string                          `json:"clause_no"`
	ClauseTitle  string                          `json:"clause_title"`
	Requirements []*ComplianceRequirementListDTO `json:"requirements"`
}

type ComplianceRequirementByComplianceIDListResponseDTO struct {
	Items      []*ComplianceRequirementByComplianceIDListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO                  `json:"pagination"`
}

type ComplianceRequirementFormDTO struct {
	ComplianceClauseNo              string  `json:"compliance_clause_no" binding:"required"`
	ComplianceClauseTitle           string  `json:"compliance_clause_title" binding:"required"`
	Title                           string  `json:"title" binding:"required"`
	ComplianceRequirementCategoryID uint    `json:"compliance_requirement_category_id" binding:"required"`
	SeverityLevelID                 uint    `json:"severity_level_id" binding:"required"`
	ResponsibleUserID               uint    `json:"responsible_user_id" binding:"required"`
	Deadline                        string  `json:"deadline" binding:"omitempty,datetime=2006-01-02"`
	ComplianceStatusID              uint    `json:"compliance_status_id" binding:"required"`
	ControlIDs                      []*uint `json:"control_ids"`
	StrategicObjectiveIDs           []*uint `json:"strategic_objective_ids"`
}

type ComplianceRequirementParsedDates struct {
	Deadline *time.Time
}

type ComplianceRequirementDetailDTO struct {
	ID                                      int                                           `json:"id"`
	ComplianceClause                        ComplianceClauseDTO                           `json:"compliance_clause"`
	Title                                   string                                        `json:"title"`
	ComplianceRequirementCategory           ComplianceRequirementCategoryDTO              `json:"compliance_requirement_category"`
	SeverityLevel                           *masterDTOs.SeverityLevelDTO                  `json:"severity_level"`
	ResponsibleUser                         organizationDTOs.UserMinimalDTO               `json:"responsible_user"`
	DeadlineString                          *string                                       `json:"deadline"`
	CreatedAt                               time.Time                                     `json:"created_at"`
	UpdatedAt                               *time.Time                                    `json:"updated_at"`
	CreatedByUser                           organizationDTOs.UserMinimalDTO               `json:"created_by"`
	UpdatedByUser                           *organizationDTOs.UserMinimalDTO              `json:"updated_by"`
	LatestComplianceRequirementStatusUpdate *ComplianceRequirementStatusUpdateDTO         `json:"compliance_status"`
	Controls                                []*ComplianceRequirementControlDTO            `json:"controls"`
	StrategicObjectives                     []*ComplianceRequirementStrategicObjectiveDTO `json:"strategic_objectives"`
	ComplianceRequirementStatusUpdates      []*ComplianceRequirementStatusUpdateDTO       `json:"compliance_requirement_status_updates"`
}

type ComplianceRequirementStatusUpdateDTO struct {
	ID               uint                                            `json:"id"`
	ComplianceStatus ComplianceStatusDTO                             `json:"compliance_status"`
	Comment          string                                          `json:"comment"`
	CreatedAt        time.Time                                       `json:"created_at"`
	UpdatedAt        *time.Time                                      `json:"updated_at"`
	CreatedByUser    organizationDTOs.UserMinimalDTO                 `json:"created_by"`
	UpdatedByUser    *organizationDTOs.UserMinimalDTO                `json:"updated_by"`
	Documents        []*ComplianceRequirementStatusUpdateDocumentDTO `json:"documents"`
}

type ComplianceRequirementStatusUpdateDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}

type ComplianceRequirementControlDTO struct {
	ID                      *uint                             `json:"id"`
	ComplianceRequirementID *uint                             `json:"compliance_requirement_id"`
	ControlID               uint                              `json:"control_id"`
	Control                 organizationDTOs.ControlDetailDTO `json:"control"`
}

type ComplianceRequirementStrategicObjectiveDTO struct {
	ID                      *uint                                  `json:"id"`
	ComplianceRequirementID *uint                                  `json:"compliance_requirement_id"`
	StrategicObjectiveID    uint                                   `json:"strategic_objective_id"`
	StrategicObjective      organizationDTOs.StrategicObjectiveDTO `json:"strategic_objective"`
}
