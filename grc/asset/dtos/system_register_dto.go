package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type SystemRegisterListDTO struct {
	ID                          uint       `json:"id"`
	TitleEn                     string     `json:"title_en"`
	TitleAr                     string     `json:"title_ar"`
	ServiceDescription          *string    `json:"service_description"`
	ServiceImportanceTitleEn    *string    `json:"service_importance_title_en"`
	ServiceImportanceTitleAr    *string    `json:"service_importance_title_ar"`
	ResponsibleITTeamTitleEn    *string    `json:"responsible_it_team_title_en"`
	ResponsibleITTeamTitleAr    *string    `json:"responsible_it_team_title_ar"`
	BusinessOwnerTitleEn        *string    `json:"business_owner_title_en"`
	BusinessOwnerTitleAr        *string    `json:"business_owner_title_ar"`
	ProcessOperationTeamTitleEn *string    `json:"process_operation_team_title_en"`
	ProcessOperationTeamTitleAr *string    `json:"process_operation_team_title_ar"`
	ProcessRecoveryTeamTitleEn  *string    `json:"process_recovery_team_title_en"`
	ProcessRecoveryTeamTitleAr  *string    `json:"process_recovery_team_title_ar"`
	RecoveryManagerPhone        string     `json:"recovery_manager_phone"`
	RecoveryManagerEmail        string     `json:"recovery_manager_email"`
	CreatedUserID               uint       `json:"created_user_id"`
	CreatedUserFirstNameEn      string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string     `json:"created_user_last_name_ar"`
	CreatedUserInitial          string     `json:"created_user_initial"`
	CreatedUserColor            string     `json:"created_user_color"`
	CreatedUserEmail            string     `json:"created_user_email"`
	CreatedUserProfileImagePath string     `json:"created_user_profile_image_path"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedUserID               *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string     `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string     `json:"updated_user_initial"`
	UpdatedUserColor            string     `json:"updated_user_color"`
	UpdatedUserEmail            string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath string     `json:"updated_user_profile_image_path"`
	UpdatedAt                   *time.Time `json:"updated_at"`
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
	TitleEn                 string  `json:"title_en" binding:"required,max=500"`
	TitleAr                 string  `json:"title_ar" binding:"required,max=500"`
	ServiceDescription      *string `json:"service_description"`
	ModulesComponents       *string `json:"modules_components"`
	ThirdPartyDependencies  *string `json:"third_party_dependencies"`
	ServiceImportanceID     *uint   `json:"service_importance_id"`
	ResponsibleITTeamID     *uint   `json:"responsible_it_team_id"`
	OutsourcingPartner      *string `json:"outsourcing_partner"`
	OperationLocationID     *uint   `json:"operation_location_id"`
	BusinessOwnerID         *uint   `json:"business_owner_id"`
	Customers               *string `json:"customers"`
	ImpactTypeID            *uint   `json:"impact_type_id"`
	DisruptionImpactLevelID *uint   `json:"disruption_impact_level_id"`
	RTO                     *int    `json:"rto"`
	RTOUnitID               *uint   `json:"rto_unit_id"`
	MTPD                    *int    `json:"mtpd"`
	MTPDUnitID              *uint   `json:"mtpd_unit_id"`
	RPO                     *int    `json:"rpo"`
	RPOUnitID               *uint   `json:"rpo_unit_id"`
	ITProcesses             *string `json:"it_processes"`
	ITServiceOwnerID        *uint   `json:"it_service_owner_id"`
	ProcessOperationTeamID  uint    `json:"process_operation_team_id" binding:"required"`
	ProcessRecoveryTeamID   uint    `json:"process_recovery_team_id" binding:"required"`
	RecoveryManagerID       *uint   `json:"recovery_manager_id"`
	RecoveryManagerPhone    string  `json:"recovery_manager_phone" binding:"required,max=20"`
	RecoveryManagerEmail    string  `json:"recovery_manager_email" binding:"required,max=255"`
	RecoveryLocation        *string `json:"recovery_location"`
	VendorRTO               *int    `json:"vendor_rto"`
	VendorRTOUnitID         *uint   `json:"vendor_rto_unit_id"`
	SLARegisterID           *uint   `json:"sla_register_id"`
}

type SystemRegisterDetailDTO struct {
	ID                     uint                                  `json:"id"`
	TitleEn                string                                `json:"title_en"`
	TitleAr                string                                `json:"title_ar"`
	ServiceDescription     *string                               `json:"service_description"`
	ModulesComponents      *string                               `json:"modules_components"`
	ThirdPartyDependencies *string                               `json:"third_party_dependencies"`
	ServiceImportance      *SystemRegisterServiceImportanceDTO   `json:"service_importance"`
	ResponsibleITTeam      *JobTeamDetailDTO                     `json:"responsible_it_team"`
	OutsourcingPartner     *string                               `json:"outsourcing_partner"`
	OperationLocation      *masterDTOs.LocationDTO               `json:"operation_location"`
	BusinessOwner          *organizationDTOs.OrganizationUnitDTO `json:"business_owner"`
	Customers              *string                               `json:"customers"`
	ImpactType             *ImpactCategoryDTO                    `json:"impact_type"`
	DisruptionImpactLevel  *SystemRegisterDisruptionLevelDTO     `json:"disruption_impact_level"`
	RTO                    *int                                  `json:"rto"`
	RTOUnit                *masterDTOs.TimeUnitDTO               `json:"rto_unit"`
	MTPD                   *int                                  `json:"mtpd"`
	MTPDUnit               *masterDTOs.TimeUnitDTO               `json:"mtpd_unit"`
	RPO                    *int                                  `json:"rpo"`
	RPOUnit                *masterDTOs.TimeUnitDTO               `json:"rpo_unit"`
	ITProcesses            *string                               `json:"it_processes"`
	ITServiceOwner         *organizationDTOs.DesignationDTO      `json:"it_service_owner"`
	ProcessOperationTeam   *JobTeamDetailDTO                     `json:"process_operation_team"`
	ProcessRecoveryTeam    *JobTeamDetailDTO                     `json:"process_recovery_team"`
	RecoveryManager        *organizationDTOs.DesignationDTO      `json:"recovery_manager"`
	RecoveryManagerPhone   string                                `json:"recovery_manager_phone"`
	RecoveryManagerEmail   string                                `json:"recovery_manager_email"`
	RecoveryLocation       *string                               `json:"recovery_location"`
	VendorRTO              *int                                  `json:"vendor_rto"`
	VendorRTOUnit          *masterDTOs.TimeUnitDTO               `json:"vendor_rto_unit"`
	SLA                    *SystemRegisterSLADTO                 `json:"sla"`
	CreatedAt              time.Time                             `json:"created_at"`
	UpdatedAt              *time.Time                            `json:"updated_at"`
	CreatedByUser          organizationDTOs.UserMinimalDTO       `json:"created_by"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO      `json:"updated_by"`
}

type SystemRegisterSLADTO struct {
	ID          uint         `json:"id"`
	SLARegister SLADetailDTO `json:"sla_register"`
}

type ImpactCategoryDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

type JobTeamDetailDTO struct {
	ID                  uint                             `json:"id"`
	TitleEn             string                           `json:"title_en"`
	TitleAr             string                           `json:"title_ar"`
	Scope               *string                          `json:"scope"`
	Objective           *string                          `json:"objective"`
	TeamLeadDesignation *organizationDTOs.DesignationDTO `json:"team_lead_designation"`
	TeamLeadUser        *organizationDTOs.UserMinimalDTO `json:"team_lead_user"`
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

type SLADetailDTO struct {
	ID              uint    `json:"id"`
	ReferenceNumber string  `json:"reference_number"`
	TitleEn         string  `json:"title_en"`
	TitleAr         string  `json:"title_ar"`
	Description     *string `json:"description,omitempty"`
}
