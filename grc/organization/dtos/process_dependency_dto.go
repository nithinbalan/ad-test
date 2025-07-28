package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
)

type ProcessDependencyListDTO struct {
	ID             uint    `json:"id"`
	ProcessID      uint    `json:"process_id"`
	ProcessTitleEn *string `json:"process_title_en"`
	ProcessTitleAr *string `json:"process_title_ar"`

	ProcessDependencySubTypeID      uint    `json:"process_dependency_sub_type_id"`
	ProcessDependencySubTypeType    *string `json:"process_dependency_sub_type_type"`
	ProcessDependencySubTypeTitleEn *string `json:"process_dependency_sub_type_title_en"`
	ProcessDependencySubTypeTitleAr *string `json:"process_dependency_sub_type_title_ar"`
	ProcessDependencyTypeTitleEn    *string `json:"process_dependency_type_title_en"`
	ProcessDependencyTypeTitleAr    *string `json:"process_dependency_type_title_ar"`

	DependentProcessID      *uint   `json:"dependent_process_id"`
	DependentProcessTitleEn *string `json:"dependent_process_title_en"`
	DependentProcessTitleAr *string `json:"dependent_process_title_ar"`

	OrganizationUnitID      *uint   `json:"organization_unit_id"`
	OrganizationUnitTitleEn *string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr *string `json:"organization_unit_title_ar"`

	AssetID      *uint   `json:"asset_id"`
	AssetTitleEn *string `json:"asset_title_en"`
	AssetTitleAr *string `json:"asset_title_ar"`

	ThirdPartyRegisterID    *uint   `json:"third_party_register_id"`
	ThirdPartyRegisterTitle *string `json:"third_party_register_title"`

	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type ProcessDependencyDetailDTO struct {
	ID                         uint                              `json:"id"`
	ProcessID                  uint                              `json:"process_id"`
	ProcessDependencySubTypeID uint                              `json:"process_dependency_sub_type_id"`
	ProcessDependencySubType   ProcessDependencySubTypeDetailDTO `json:"process_dependency_sub_type"`
	DependentProcessID         *uint                             `json:"dependent_process_id"`
	DependentProcess           *ProcessDetailDTO                 `json:"dependent_process"`
	OrganizationUnitID         *uint                             `json:"organization_unit_id"`
	OrganizationUnit           *OrganizationUnit1DTO             `json:"organization_unit"`
	AssetID                    *uint                             `json:"asset_id"`
	Asset                      *AssetMinimalDTO                  `json:"asset"`
	ThirdPartyRegisterID       *uint                             `json:"third_party_register_id"`
	ThirdPartyRegister         *ThirdPartyMinimalDTO             `json:"third_party_register"`
	Title                      string                            `json:"title"`
	Description                *string                           `json:"description"`
	Rto                        *int                              `json:"rto"`
	RtoUnitId                  *uint                             `json:"rto_uint_id"`
	RtoUnit                    masterDTOs.TimeUnitDTO            `json:"rto_unit"`
	Mtpd                       *int                              `json:"mtpd"`
	MtpdUnitId                 *uint                             `json:"mtpd_uint_id"`
	MtpdUnit                   masterDTOs.TimeUnitDTO            `json:"mtpd_unit"`
	Rpo                        *int                              `json:"rpo"`
	RpoUnitId                  *uint                             `json:"rpo_uint_id"`
	RpoUnit                    masterDTOs.TimeUnitDTO            `json:"rpo_unit"`
}

type ProcessDependencyFormDTO struct {
	ProcessDependencySubTypeID uint    `json:"process_dependency_sub_type_id" binding:"required"`
	DependentProcessIDs        []*uint `json:"dependent_process_ids"`
	OrganizationUnitIDs        []*uint `json:"organization_unit_ids"`
	AssetIDs                   []*uint `json:"asset_ids"`
	ThirdPartyRegisterIDs      []*uint `json:"third_party_register_ids"`
	Title                      *string `json:"title"`
	Description                *string `json:"description"`
}

type ProcessDependencyUpdateFormDTO struct {
	ProcessDependencySubTypeID uint    `json:"process_dependency_sub_type_id" binding:"required"`
	DependentProcessID         *uint   `json:"dependent_process_id"`
	OrganizationUnitID         *uint   `json:"organization_unit_id"`
	AssetID                    *uint   `json:"asset_id"`
	ThirdPartyRegisterID       *uint   `json:"third_party_register_id"`
	Title                      *string `json:"title" binding:"max=255"`
	Description                *string `json:"description"`
}

type ProcessDependencyListQueryParamsDTO struct {
	Sort    string              `json:"sort"`
	SortBy  string              `json:"sort_by"`
	Keyword string              `json:"keyword"`
	Filters map[string][]string `json:"filters"`
}

type ThirdPartyMinimalDTO struct {
	ID                 uint                             `json:"id"`
	ReferenceNumber    string                           `json:"reference_number"`
	Title              string                           `json:"title"`
	Country            masterDTOs.CountryDTO            `json:"country"`
	Address            *string                          `json:"address"`
	Website            *string                          `json:"website"`
	OrganizationUnit   settingsDTOs.OrganizationUnitDTO `json:"organization_unit"`
	RepresentativeName *string                          `json:"representative_name"`
	RepresentativeRole *string                          `json:"representative_role"`
	Email              *string                          `json:"email"`
	PhoneNumber        *string                          `json:"phone_number"`
}
