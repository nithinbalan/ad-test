package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type InternalCommitteeListDTO struct {
	ID                             uint   `json:"id"`
	ReferenceNumber                string `json:"reference_number"`
	RequestingEntityID             uint   `json:"requesting_entity_id"`
	RequestingEntityTitleEn        string `json:"requesting_entity_title_en"`
	RequestingEntityTitleAr        string `json:"requesting_entity_title_ar"`
	CommitteeNameEn                string `json:"committee_name_en"`
	CommitteeNameAr                string `json:"committee_name_ar"`
	CommitteeTypeTitleEn           string `json:"committee_type_title_en"`
	CommitteeTypeTitleAr           string `json:"committee_type_title_ar"`
	CommitteeTermTitleEn           string `json:"committee_term_title_en"`
	CommitteeTermTitleAr           string `json:"committee_term_title_ar"`
	CommitteeClassificationTitleEn string `json:"committee_classification_title_en"`
	CommitteeClassificationTitleAr string `json:"committee_classification_title_ar"`
	CommitteeStatusTitleEn         string `json:"committee_status_title_en"`
	CommitteeStatusTitleAr         string `json:"committee_status_title_ar"`
	CommitteeStatusColorCode       string `json:"committee_status_color_code"`
	Version                        string `json:"version"`
	TermEndDate                    string `json:"term_end_date"`
	Justification                  string `json:"justification"`
	RolesAndResponsibilities       string `json:"roles_and_responsibilities"`
	CommitteeAuthorities           string `json:"committee_authorities"`
	CreatedAt                      string `json:"created_at"`
	UpdatedAt                      string `json:"updated_at"`
	CreatedUserID                  uint   `json:"created_user_id"`
	CreatedUserFirstNameEn         string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr         string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn          string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr          string `json:"created_user_last_name_ar"`
	CreatedUserInitial             string `json:"created_user_initial"`
	CreatedUserColor               string `json:"created_user_color"`
	CreatedUserEmail               string `json:"created_user_email"`
	CreatedUserProfileImagePath    string `json:"created_user_profile_image_path"`
	UpdatedUserID                  *uint  `json:"updated_user_id"`
	UpdatedUserFirstNameEn         string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr         string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn          string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr          string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial             string `json:"updated_user_initial"`
	UpdatedUserColor               string `json:"updated_user_color"`
	UpdatedUserEmail               string `json:"updated_user_email"`
	UpdatedUserProfileImagePath    string `json:"updated_user_profile_image_path"`
}

type InternalCommitteeListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	DateFrom    string
	DateTo      string
	Export      bool
}

type InternalCommitteeListResponseDTO struct {
	Items      []*InternalCommitteeListDTO  `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type InternalCommitteeFormDTO struct {
	CommitteeNameEn           string                                             `json:"committee_name_en" binding:"omitempty,max=500"`
	CommitteeNameAr           string                                             `json:"committee_name_ar" binding:"omitempty,max=500"`
	CommitteeTypeID           uint                                               `json:"committee_type_id"`
	CommitteeTermID           uint                                               `json:"committee_term_id"`
	TermEndDate               *string                                            `json:"term_end_date" binding:"omitempty,datetime=2006-01-02"`
	CommitteeClassificationID uint                                               `json:"committee_classification_id"`
	Justification             *string                                            `json:"justification"`
	RolesAndResponsibilities  *string                                            `json:"roles_and_responsibilities"`
	CommitteeAuthorities      *string                                            `json:"committee_authorities"`
	NomineeDepartments        []InternalCommitteeVersionNomineeDepartmentFormDTO `json:"nominee_departments"`
}

type InternalCommitteeVersionNomineeDepartmentFormDTO struct {
	NomineeDepartmentID uint    `json:"nominee_department_id"`
	MembershipType      *string `json:"membership_type"`
}

type InternalCommitteeParsedDates struct {
	TermEndDate *time.Time
}

type InternalCommitteeDetailDTO struct {
	ID                 uint                               `json:"id"`
	ReferenceNumber    string                             `json:"reference_number"`
	RequestingEntityID uint                               `json:"requesting_entity_id"`
	RequestingEntity   *settingsDTOs.OrganizationUnitDTO  `json:"requesting_entity"`
	CreatedByUser      *UserMinimalDTO                    `json:"created_by"`
	UpdatedByUser      *UserMinimalDTO                    `json:"updated_by"`
	CreatedAt          time.Time                          `json:"created_at"`
	UpdatedAt          *time.Time                         `json:"updated_at"`
	LatestVersion      *InternalCommitteeVersionDetailDTO `json:"latest_version"`
}

type InternalCommitteeVersionDetailDTO struct {
	ID                        uint                                                  `json:"id"`
	InternalCommitteeID       uint                                                  `json:"internal_committee_id"`
	CommitteeNameEn           string                                                `json:"committee_name_en"`
	CommitteeNameAr           string                                                `json:"committee_name_ar"`
	CommitteeTypeID           uint                                                  `json:"committee_type_id"`
	CommitteeType             *CommitteeTypeDTO                                     `json:"committee_type"`
	CommitteeTermID           uint                                                  `json:"committee_term_id"`
	CommitteeTerm             *CommitteeTermDTO                                     `json:"committee_term"`
	TermEndDateString         *string                                               `json:"term_end_date"`
	CommitteeClassificationID uint                                                  `json:"committee_classification_id"`
	CommitteeClassification   *CommitteeClassificationDetailsDTO                    `json:"committee_classification"`
	Justification             *string                                               `json:"justification"`
	RolesAndResponsibilities  *string                                               `json:"roles_and_responsibilities"`
	CommitteeAuthorities      *string                                               `json:"committee_authorities"`
	CommitteeStatusID         uint                                                  `json:"committee_status_id"`
	CommitteeStatus           *InternalCommitteeStatusDTO                           `json:"committee_status"`
	Version                   *float64                                              `json:"version"`
	IsLatestVersion           *bool                                                 `json:"is_latest_version"`
	CreatedByUser             *UserMinimalDTO                                       `json:"created_by"`
	UpdatedByUser             *UserMinimalDTO                                       `json:"updated_by"`
	CreatedAt                 time.Time                                             `json:"created_at"`
	UpdatedAt                 *time.Time                                            `json:"updated_at"`
	SubmittedByUser           *UserMinimalDTO                                       `json:"submitted_by"`
	ApprovedByUser            *UserMinimalDTO                                       `json:"approved_by"`
	SubmittedAt               *time.Time                                            `json:"submitted_at"`
	ApprovedAt                *time.Time                                            `json:"approved_at"`
	NextReviewLevel           *int                                                  `json:"next_review_level"`
	WorkflowStatusID          *uint                                                 `json:"workflow_status_id"`
	WorkflowStatus            *masterDTOs.WorkflowStatusDTO                         `json:"workflow_status"`
	NomineeDepartments        []*InternalCommitteeVersionNomineeDepartmentDetailDTO `json:"nominee_departments"`
}

type InternalCommitteeVersionNomineeDepartmentDetailDTO struct {
	ID                         uint                                                      `json:"id"`
	InternalCommitteeVersionID uint                                                      `json:"internal_committee_version_id"`
	NomineeDepartmentID        uint                                                      `json:"nominee_department_id"`
	NomineeDepartment          *settingsDTOs.OrganizationUnitDTO                         `json:"nominee_department"`
	MembershipType             *string                                                   `json:"membership_type"`
	Participants               []InternalCommitteeVersionNomineeDepartmentParticipantDTO `json:"participants"`
}

type InternalCommitteeVersionNomineeDepartmentParticipantDTO struct {
	ID         uint                              `json:"id"`
	UserID     uint                              `json:"user_id"`
	User       *UserMinimalDTO                   `json:"user"`
	Department *settingsDTOs.OrganizationUnitDTO `json:"department"`
}

type InternalCommitteeParticipantFormDTO struct {
	UserIDs []uint `json:"user_ids" binding:"required,min=1"`
}

type ReformInternalCommitteeDTO struct {
	InternalCommitteeID       uint                                               `json:"internal_committee_id" binding:"required"`
	CommitteeNameEn           string                                             `json:"committee_name_en" binding:"omitempty,max=500"`
	CommitteeNameAr           string                                             `json:"committee_name_ar" binding:"omitempty,max=500"`
	CommitteeTypeID           uint                                               `json:"committee_type_id"`
	CommitteeTermID           uint                                               `json:"committee_term_id"`
	TermEndDate               *string                                            `json:"term_end_date" binding:"omitempty,datetime=2006-01-02"`
	CommitteeClassificationID uint                                               `json:"committee_classification_id"`
	Justification             *string                                            `json:"justification"`
	RolesAndResponsibilities  *string                                            `json:"roles_and_responsibilities"`
	CommitteeAuthorities      *string                                            `json:"committee_authorities"`
	NomineeDepartments        []InternalCommitteeVersionNomineeDepartmentFormDTO `json:"nominee_departments"`
}
