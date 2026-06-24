package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type SystemRegisterListDTO struct {
	ID                           uint                              `json:"id"`
	TitleEn                      string                            `json:"title_en"`
	TitleAr                      string                            `json:"title_ar"`
	ServiceDescription           *string                           `json:"service_description"`
	ModulesComponents            *string                           `json:"modules_components"`
	ThirdPartyDependencies       *string                           `json:"third_party_dependencies"`
	ServiceImportanceID          *uint                             `json:"service_importance_id"`
	ServiceImportanceTitleEn     *string                           `json:"service_importance_title_en"`
	ServiceImportanceTitleAr     *string                           `json:"service_importance_title_ar"`
	ResponsibleITTeamID          *uint                             `json:"responsible_it_team_id"`
	ResponsibleITTeamTitleEn     *string                           `json:"responsible_it_team_title_en"`
	ResponsibleITTeamTitleAr     *string                           `json:"responsible_it_team_title_ar"`
	OutsourcingPartner           *string                           `json:"outsourcing_partner"`
	OperationLocationID          *uint                             `json:"operation_location_id"`
	OperationLocationTitleEn     *string                           `json:"operation_location_title_en"`
	OperationLocationTitleAr     *string                           `json:"operation_location_title_ar"`
	BusinessOwnerID              *uint                             `json:"business_owner_id"`
	BusinessOwnerTitleEn         *string                           `json:"business_owner_title_en"`
	BusinessOwnerTitleAr         *string                           `json:"business_owner_title_ar"`
	Customers                    *string                           `json:"customers"`
	DisruptionImpactLevelID      *uint                             `json:"disruption_impact_level_id"`
	DisruptionImpactLevelTitleEn *string                           `json:"disruption_impact_level_title_en"`
	DisruptionImpactLevelTitleAr *string                           `json:"disruption_impact_level_title_ar"`
	RTO                          *int                              `json:"rto"`
	RTOUnitID                    *uint                             `json:"rto_unit_id"`
	RTOUnitTitleEn               *string                           `json:"rto_unit_title_en"`
	RTOUnitTitleAr               *string                           `json:"rto_unit_title_ar"`
	MTPD                         *int                              `json:"mtpd"`
	MTPDUnitID                   *uint                             `json:"mtpd_unit_id"`
	MTPDUnitTitleEn              *string                           `json:"mtpd_unit_title_en"`
	MTPDUnitTitleAr              *string                           `json:"mtpd_unit_title_ar"`
	RPO                          *int                              `json:"rpo"`
	RPOUnitID                    *uint                             `json:"rpo_unit_id"`
	RPOUnitTitleEn               *string                           `json:"rpo_unit_title_en"`
	RPOUnitTitleAr               *string                           `json:"rpo_unit_title_ar"`
	ITProcesses                  *string                           `json:"it_processes"`
	ITServiceOwnerID             *uint                             `json:"it_service_owner_id"`
	ITServiceOwnerTitleEn        *string                           `json:"it_service_owner_title_en"`
	ITServiceOwnerTitleAr        *string                           `json:"it_service_owner_title_ar"`
	ProcessOperationTeamID       uint                              `json:"process_operation_team_id"`
	ProcessOperationTeamTitleEn  *string                           `json:"process_operation_team_title_en"`
	ProcessOperationTeamTitleAr  *string                           `json:"process_operation_team_title_ar"`
	ProcessRecoveryTeamID        uint                              `json:"process_recovery_team_id"`
	ProcessRecoveryTeamTitleEn   *string                           `json:"process_recovery_team_title_en"`
	ProcessRecoveryTeamTitleAr   *string                           `json:"process_recovery_team_title_ar"`
	RecoveryManagerID            *uint                             `json:"recovery_manager_id"`
	RecoveryManagerTitleEn       *string                           `json:"recovery_manager_title_en"`
	RecoveryManagerTitleAr       *string                           `json:"recovery_manager_title_ar"`
	RecoveryManagerPhone         string                            `json:"recovery_manager_phone"`
	RecoveryManagerEmail         string                            `json:"recovery_manager_email"`
	RecoveryLocation             *string                           `json:"recovery_location"`
	VendorRTO                    *int                              `json:"vendor_rto"`
	VendorRTOUnitID              *uint                             `json:"vendor_rto_unit_id"`
	VendorRTOUnitTitleEn         *string                           `json:"vendor_rto_unit_title_en"`
	VendorRTOUnitTitleAr         *string                           `json:"vendor_rto_unit_title_ar"`
	CreatedUserID                uint                              `json:"created_user_id"`
	CreatedUserFirstNameEn       string                            `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr       string                            `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn        string                            `json:"created_user_last_name_en"`
	CreatedUserLastNameAr        string                            `json:"created_user_last_name_ar"`
	CreatedUserInitial           string                            `json:"created_user_initial"`
	CreatedUserColor             string                            `json:"created_user_color"`
	CreatedUserEmail             string                            `json:"created_user_email"`
	CreatedUserProfileImagePath  string                            `json:"created_user_profile_image_path"`
	CreatedAt                    time.Time                         `json:"created_at"`
	UpdatedUserID                *uint                             `json:"updated_user_id"`
	UpdatedUserFirstNameEn       *string                           `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr       *string                           `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn        string                            `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr        string                            `json:"updated_user_last_name_ar"`
	UpdatedUserInitial           string                            `json:"updated_user_initial"`
	UpdatedUserColor             string                            `json:"updated_user_color"`
	UpdatedUserEmail             string                            `json:"updated_user_email"`
	UpdatedUserProfileImagePath  string                            `json:"updated_user_profile_image_path"`
	UpdatedAt                    *time.Time                        `json:"updated_at"`
	ImpactTypes                  []SystemRegisterImpactTypeListDTO `json:"impact_types" gorm:"-"`
}

type SystemRegisterImpactTypeListDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
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
	ImpactTypeIDs           []*uint `json:"impact_type_ids"`
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
	ID                        uint                                  `json:"id"`
	TitleEn                   string                                `json:"title_en"`
	TitleAr                   string                                `json:"title_ar"`
	ServiceDescription        *string                               `json:"service_description"`
	ModulesComponents         *string                               `json:"modules_components"`
	ThirdPartyDependencies    *string                               `json:"third_party_dependencies"`
	ServiceImportance         *SystemRegisterServiceImportanceDTO   `json:"service_importance"`
	ResponsibleITTeam         *JobTeamDetailDTO                     `json:"responsible_it_team"`
	OutsourcingPartner        *string                               `json:"outsourcing_partner"`
	OperationLocation         *masterDTOs.LocationDTO               `json:"operation_location"`
	BusinessOwner             *organizationDTOs.OrganizationUnitDTO `json:"business_owner"`
	Customers                 *string                               `json:"customers"`
	SystemRegisterImpactTypes []*SystemRegisterImpactTypeDTO        `json:"impact_types"`
	DisruptionImpactLevel     *SystemRegisterDisruptionLevelDTO     `json:"disruption_impact_level"`
	RTO                       *int                                  `json:"rto"`
	RTOUnit                   *BusinessContinuityDurationUnitDTO    `json:"rto_unit"`
	MTPD                      *int                                  `json:"mtpd"`
	MTPDUnit                  *BusinessContinuityDurationUnitDTO    `json:"mtpd_unit"`
	RPO                       *int                                  `json:"rpo"`
	RPOUnit                   *BusinessContinuityDurationUnitDTO    `json:"rpo_unit"`
	ITProcesses               *string                               `json:"it_processes"`
	ITServiceOwner            *organizationDTOs.DesignationDTO      `json:"it_service_owner"`
	ProcessOperationTeam      *JobTeamDetailDTO                     `json:"process_operation_team"`
	ProcessRecoveryTeam       *JobTeamDetailDTO                     `json:"process_recovery_team"`
	RecoveryManager           *organizationDTOs.DesignationDTO      `json:"recovery_manager"`
	RecoveryManagerPhone      string                                `json:"recovery_manager_phone"`
	RecoveryManagerEmail      string                                `json:"recovery_manager_email"`
	RecoveryLocation          *string                               `json:"recovery_location"`
	VendorRTO                 *int                                  `json:"vendor_rto"`
	VendorRTOUnit             *BusinessContinuityDurationUnitDTO    `json:"vendor_rto_unit"`
	SystemRegisterSLA         *SystemRegisterSLADTO                 `json:"system_register_sla"`
	CreatedAt                 time.Time                             `json:"created_at"`
	UpdatedAt                 *time.Time                            `json:"updated_at"`
	CreatedByUser             organizationDTOs.UserMinimalDTO       `json:"created_by"`
	UpdatedByUser             *organizationDTOs.UserMinimalDTO      `json:"updated_by"`
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

type SystemRegisterImpactTypeDTO struct {
	ID               uint              `json:"id"`
	SystemRegisterID uint              `json:"system_register_id"`
	ImpactTypeID     uint              `json:"impact_type_id"`
	ImpactType       ImpactCategoryDTO `json:"impact_type"`
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

// BusinessContinuityDurationUnitDTO represents business continuity duration units (not generic time units)
type BusinessContinuityDurationUnitDTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Minutes int    `json:"minutes"`
	Type    string `json:"type"`
	Status  string `json:"status"`
}
