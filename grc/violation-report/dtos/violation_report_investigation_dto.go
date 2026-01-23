package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ViolationReportInvestigationListDTO struct {
	ID                                 uint                                   `json:"id"`
	ViolationReportID                  uint                                   `json:"violation_report_id"`
	InvestigationDate                  string                                 `json:"investigation_date"`
	StepsTaken                         *string                                `json:"steps_taken"`
	SuggestedRecommendation            *string                                `json:"suggested_recommendation"`
	FindingSummary                     *string                                `json:"finding_summary"`
	DisciplinaryProcedure              *string                                `json:"disciplinary_procedure"`
	IsDisciplinaryRegulatoryCompliance *bool                                  `json:"is_disciplinary_regulatory_compliance"`
	ActualRecommendation               *string                                `json:"actual_recommendation"`
	NoteAndObservations                *string                                `json:"note_and_observations"`
	ViolationReportReference           string                                 `json:"violation_report_reference"`
	ViolationTitle                     string                                 `json:"violation_title"`
	Departments                        []organizationDTOs.OrganizationUnitDTO `json:"departments"`
	CreatedUserID                      uint                                   `json:"created_user_id"`
	CreatedUserFirstNameEn             string                                 `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr             string                                 `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn              string                                 `json:"created_user_last_name_en"`
	CreatedUserLastNameAr              string                                 `json:"created_user_last_name_ar"`
	CreatedUserInitial                 string                                 `json:"created_user_initial"`
	CreatedUserColor                   string                                 `json:"created_user_color"`
	CreatedUserEmail                   string                                 `json:"created_user_email"`
	CreatedUserProfileImagePath        string                                 `json:"created_user_profile_image_path"`
	UpdatedUserID                      *uint                                  `json:"updated_user_id"`
	UpdatedUserFirstNameEn             string                                 `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr             string                                 `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn              string                                 `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr              string                                 `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                 string                                 `json:"updated_user_initial"`
	UpdatedUserColor                   string                                 `json:"updated_user_color"`
	UpdatedUserEmail                   string                                 `json:"updated_user_email"`
	UpdatedUserProfileImagePath        string                                 `json:"updated_user_profile_image_path"`
	CreatedAt                          string                                 `json:"created_at"`
	UpdatedAt                          string                                 `json:"updated_at"`
}

type ViolationReportInvestigationListQueryParamsDTO struct {
	Sort              string
	SortBy            string
	Keyword           string
	Status            string
	Page              int
	PageSize          int
	Filters           map[string][]string
	QuickFilter       []string
	InvestigationDate string
}

type ViolationReportInvestigationListResponseDTO struct {
	Items      []*ViolationReportInvestigationListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO           `json:"pagination"`
}

type ViolationReportInvestigationFormDTO struct {
	ViolationReportID                  uint                          `json:"violation_report_id"`
	InvestigationDate                  string                        `json:"investigation_date" binding:"required"`
	StepsTaken                         *string                       `json:"steps_taken"`
	SuggestedRecommendation            *string                       `json:"suggested_recommendation"`
	FindingSummary                     *string                       `json:"finding_summary"`
	DisciplinaryProcedure              *string                       `json:"disciplinary_procedure"`
	IsDisciplinaryRegulatoryCompliance *bool                         `json:"is_disciplinary_regulatory_compliance"`
	ActualRecommendation               *string                       `json:"actual_recommendation"`
	NoteAndObservations                *string                       `json:"note_and_observations"`
	Documents                          []*commonDTOs.DocumentFormDTO `json:"documents"`
	DepartmentIds                      []*uint                       `json:"department_ids"`
}

type ViolationReportInvestigationDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         string  `json:"name"`
	Token        string  `json:"token"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}

type ViolationReportInvestigationDetailDTO struct {
	ID                                 uint                                         `json:"id"`
	ViolationReport                    ViolationReportDetailDTO                     `json:"violation_report"`
	InvestigationDateString            string                                       `json:"investigation_date"`
	StepsTaken                         *string                                      `json:"steps_taken"`
	SuggestedRecommendation            *string                                      `json:"suggested_recommendation"`
	FindingSummary                     *string                                      `json:"finding_summary"`
	DisciplinaryProcedure              *string                                      `json:"disciplinary_procedure"`
	IsDisciplinaryRegulatoryCompliance *bool                                        `json:"is_disciplinary_regulatory_compliance"`
	ActualRecommendation               *string                                      `json:"actual_recommendation"`
	NoteAndObservations                *string                                      `json:"note_and_observations"`
	Documents                          []*ViolationReportInvestigationDocumentDTO   `json:"documents"`
	Departments                        []*ViolationReportInvestigationDepartmentDTO `json:"departments"`
	CreatedAt                          time.Time                                    `json:"created_at"`
	UpdatedAt                          *time.Time                                   `json:"updated_at"`
	CreatedByUser                      organizationDTOs.UserMinimalDTO              `json:"created_by"`
	UpdatedByUser                      *organizationDTOs.UserMinimalDTO             `json:"updated_by"`
}

type ViolationReportInvestigationParsedDates struct {
	InvestigationDate time.Time
}

type ViolationReportInvestigationDepartmentDTO struct {
	ID                             uint                                 `json:"id"`
	ViolationReportInvestigationID uint                                 `json:"violation_report_investigation_id"`
	DepartmentID                   uint                                 `json:"department_id"`
	Department                     organizationDTOs.OrganizationUnitDTO `json:"department"`
}
