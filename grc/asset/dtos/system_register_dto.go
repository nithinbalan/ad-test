package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type SystemRegisterListDTO struct {
	ID                            uint       `json:"id"`
	ReferenceNumber               string     `json:"reference_number"`
	NumberOfPeopleNecessary       *int       `json:"number_of_people_necessary"`
	BusinessOwnerUnitTitleEn      *string    `json:"business_owner_unit_title_en"`
	BusinessOwnerUnitTitleAr      *string    `json:"business_owner_unit_title_ar"`
	ITResponsibleUnitTitleEn      *string    `json:"it_responsible_unit_title_en"`
	ITResponsibleUnitTitleAr      *string    `json:"it_responsible_unit_title_ar"`
	ThirdPartyDependencies        *string    `json:"third_party_dependencies"`
	OutsourcingPartners           *string    `json:"outsourcing_partners"`
	CustomerUser                  *string    `json:"customer_user"`
	ProcessOperationTeamID        *uint      `json:"process_operation_team_id"`
	ProcessOperationTeamTitleEn   *string    `json:"process_operation_team_title_en"`
	ProcessOperationTeamTitleAr   *string    `json:"process_operation_team_title_ar"`
	ProcessRecoveryTeamID         *uint      `json:"process_recovery_team_id"`
	ProcessRecoveryTeamTitleEn    *string    `json:"process_recovery_team_title_en"`
	ProcessRecoveryTeamTitleAr    *string    `json:"process_recovery_team_title_ar"`
	RecoveryManagerFirstNameEn    *string    `json:"recovery_manager_first_name_en"`
	RecoveryManagerFirstNameAr    *string    `json:"recovery_manager_first_name_ar"`
	RecoveryManagerLastNameEn     *string    `json:"recovery_manager_last_name_en"`
	RecoveryManagerLastNameAr     *string    `json:"recovery_manager_last_name_ar"`
	RecoveryManagerContactDetails *string    `json:"recovery_manager_contact_details"`
	CreatedUserID                 uint       `json:"created_user_id"`
	CreatedUserFirstNameEn        string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr        string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn         string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr         string     `json:"created_user_last_name_ar"`
	CreatedUserInitial            string     `json:"created_user_initial"`
	CreatedUserColor              string     `json:"created_user_color"`
	CreatedUserEmail              string     `json:"created_user_email"`
	CreatedUserProfileImagePath   string     `json:"created_user_profile_image_path"`
	CreatedAt                     time.Time  `json:"created_at"`
	UpdatedUserID                 *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn        *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr        *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn         *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr         *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial            string     `json:"updated_user_initial"`
	UpdatedUserColor              string     `json:"updated_user_color"`
	UpdatedUserEmail              string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath   string     `json:"updated_user_profile_image_path"`
	UpdatedAt                     *time.Time `json:"updated_at"`
}

type SystemRegisterListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
	Export   bool
}

type SystemRegisterListResponseDTO struct {
	Items      []*SystemRegisterListDTO     `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type SystemRegisterFormDTO struct {
	NumberOfPeopleNecessary       *int    `json:"number_of_people_necessary"`
	BusinessOwnerUnitID           *uint   `json:"business_owner_unit_id"`
	ITResponsibleUnitID           *uint   `json:"it_responsible_unit_id"`
	ThirdPartyDependencies        *string `json:"third_party_dependencies"`
	OutsourcingPartners           *string `json:"outsourcing_partners"`
	CustomerUser                  *string `json:"customer_user"`
	ProcessOperationTeamID        *uint   `json:"process_operation_team_id"`
	ProcessRecoveryTeamID         *uint   `json:"process_recovery_team_id"`
	RecoveryManagerID             *uint   `json:"recovery_manager_id"`
	RecoveryManagerContactDetails *string `json:"recovery_manager_contact_details"`
	JobTitleIDs                   []uint  `json:"job_title_ids"`
	OperationalLocationIDs        []uint  `json:"operational_location_ids"`
	RecoveryLocationIDs           []uint  `json:"recovery_location_ids"`
	RecoveryNeededToolIDs         []uint  `json:"recovery_needed_tool_ids"`
}

type SystemRegisterDetailDTO struct {
	ID                                 uint                                  `json:"id"`
	ReferenceNumber                    string                                `json:"reference_number"`
	NumberOfPeopleNecessary            *int                                  `json:"number_of_people_necessary"`
	BusinessOwnerUnit                  *organizationDTOs.OrganizationUnitDTO `json:"business_owner_unit"`
	ITResponsibleUnit                  *organizationDTOs.OrganizationUnitDTO `json:"it_responsible_unit"`
	ThirdPartyDependencies             *string                               `json:"third_party_dependencies"`
	OutsourcingPartners                *string                               `json:"outsourcing_partners"`
	CustomerUser                       *string                               `json:"customer_user"`
	ProcessOperationTeam               *JobTeamDetailDTO                     `json:"process_operation_team"`
	ProcessRecoveryTeam                *JobTeamDetailDTO                     `json:"process_recovery_team"`
	RecoveryManager                    *organizationDTOs.UserMinimalDTO      `json:"recovery_manager"`
	RecoveryManagerContactDetails      *string                               `json:"recovery_manager_contact_details"`
	SystemRegisterJobTitles            []*JobTitleDTO                        `json:"job_titles"`
	SystemRegisterOperationalLocations []*OperationalLocationDTO             `json:"operational_locations"`
	SystemRegisterRecoveryLocations    []*RecoveryLocationDTO                `json:"recovery_locations"`
	SystemRegisterRecoveryNeededTools  []*RecoveryToolDTO                    `json:"recovery_needed_tools"`
	CreatedAt                          time.Time                             `json:"created_at"`
	UpdatedAt                          *time.Time                            `json:"updated_at"`
	CreatedByUser                      organizationDTOs.UserMinimalDTO       `json:"created_by"`
	UpdatedByUser                      *organizationDTOs.UserMinimalDTO      `json:"updated_by"`
}

type JobTitleDTO struct {
	ID               uint                            `json:"id"`
	SystemRegisterID uint                            `json:"system_register_id"`
	DesignationID    uint                            `json:"designation_id"`
	Designation      organizationDTOs.DesignationDTO `json:"designation"`
}

type OperationalLocationDTO struct {
	ID                    uint                    `json:"id"`
	SystemRegisterID      uint                    `json:"system_register_id"`
	OperationalLocationID uint                    `json:"operational_location_id"`
	OperationalLocation   *masterDTOs.LocationDTO `json:"operational_location"`
}

type RecoveryLocationDTO struct {
	ID                 uint                    `json:"id"`
	SystemRegisterID   uint                    `json:"system_register_id"`
	RecoveryLocationID uint                    `json:"recovery_location_id"`
	RecoveryLocation   *masterDTOs.LocationDTO `json:"recovery_location"`
}

type RecoveryToolDTO struct {
	ID                   uint                          `json:"id"`
	SystemRegisterID     uint                          `json:"system_register_id"`
	RecoveryNeededToolID uint                          `json:"recovery_need_tool_id"`
	RecoveryNeededTool   SystemRegisterRecoveryToolDTO `json:"recovery_tool"`
}

type JobTeamDetailDTO struct {
	ID                  uint                             `json:"id"`
	TitleEn             string                           `json:"title_en"`
	TitleAr             string                           `json:"title_ar"`
	Scope               *string                          `json:"scope"`
	Objective           *string                          `json:"objective"`
	TeamLeadDesignation *organizationDTOs.DesignationDTO `json:"team_lead_designation"`
	JobTeamMembers      []*JobTeamMembersDTO             `json:"members"`
	CreatedAt           time.Time                        `json:"created_at"`
	UpdatedAt           *time.Time                       `json:"updated_at"`
	CreatedByUser       *organizationDTOs.UserMinimalDTO `json:"created_by"`
	UpdatedByUser       *organizationDTOs.UserMinimalDTO `json:"updated_by"`
}

type JobTeamMembersDTO struct {
	ID            uint                            `json:"id"`
	JobTeamID     uint                            `json:"job_team_id"`
	DesignationID uint                            `json:"designation_id"`
	Designation   organizationDTOs.DesignationDTO `json:"designation"`
}
