package services

import (
	assetDTOs "jethur/grc/asset/dtos"
	assetRepsitories "jethur/grc/asset/repositories"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
)

// Service Interface
type AssetAssociatedItemServiceInterface interface {
	GetAllControls(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) ([]map[string]interface{}, error)
	GetAllIncidents(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemIncidentListResponseDTO, error)
	GetAllRisks(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemRiskListResponseDTO, error)
}

// Service Implementation
type AssetAssociatedItemService struct {
	repo                           assetRepsitories.AssetAssociatedItemRepositoryInterface
	assetRegisterSvc               AssetRegisterServiceInterface
	controlImplementationStatusSvc organizationServices.ControlImplementationStatusServiceInterface
	controlEfficiencyStatusService organizationServices.ControlEfficiencyStatusServiceInterface
}

// Constructor Function
func NewAssetAssociatedItemService(
	repo assetRepsitories.AssetAssociatedItemRepositoryInterface,
	assetRegisterSvc AssetRegisterServiceInterface,
	controlImplementationStatusSvc organizationServices.ControlImplementationStatusServiceInterface,
	controlEfficiencyStatusService organizationServices.ControlEfficiencyStatusServiceInterface,
) AssetAssociatedItemServiceInterface {
	return &AssetAssociatedItemService{
		repo:                           repo,
		assetRegisterSvc:               assetRegisterSvc,
		controlImplementationStatusSvc: controlImplementationStatusSvc,
		controlEfficiencyStatusService: controlEfficiencyStatusService,
	}
}

func (aais *AssetAssociatedItemService) GetAllControls(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) ([]map[string]interface{}, error) {
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}

	controlList, err := aais.repo.GetAllControls(assetID, params)
	if err != nil {
		return nil, err
	}

	if len(controlList) == 0 {
		return []map[string]interface{}{}, nil
	}

	implementationStatuses, err := aais.controlImplementationStatusSvc.GetAll("", "", "", "")
	if err != nil {
		return nil, err
	}
	implementationStatusMap := make(map[string]models.ControlImplementationStatus)
	for _, status := range implementationStatuses {
		implementationStatusMap[status.Type] = *status
	}

	efficiencyStatuses, err := aais.controlEfficiencyStatusService.GetAll("", "", "", "")
	if err != nil {
		return nil, err
	}
	efficiencyStatusMap := make(map[string]models.ControlEfficiencyStatus)
	for _, status := range efficiencyStatuses {
		efficiencyStatusMap[status.Type] = *status
	}

	groupedControls := make(map[string]map[string][]*assetDTOs.AssetControlListDTO)
	for _, control := range controlList {
		implementationStatusType := control.ControlImplementationStatusType
		efficiencyStatusType := control.ControlEfficiencyStatusType

		if _, exists := groupedControls[implementationStatusType]; !exists {
			groupedControls[implementationStatusType] = make(map[string][]*assetDTOs.AssetControlListDTO)
		}
		if _, exists := groupedControls[implementationStatusType][efficiencyStatusType]; !exists {
			groupedControls[implementationStatusType][efficiencyStatusType] = []*assetDTOs.AssetControlListDTO{}
		}
		groupedControls[implementationStatusType][efficiencyStatusType] = append(groupedControls[implementationStatusType][efficiencyStatusType], control)
	}

	var response []map[string]interface{}
	for implStatusType, effMap := range groupedControls {
		implItem := map[string]interface{}{
			"status":   implStatusType,
			"title_en": "",
			"title_ar": "",
			"count":    0, // Count of implementation status items
			"items":    []map[string]interface{}{},
		}

		if status, exists := implementationStatusMap[implStatusType]; exists {
			implItem["title_en"] = status.TitleEn
			implItem["title_ar"] = status.TitleAr
		}

		for effStatusType, controls := range effMap {
			effItem := map[string]interface{}{
				"efficiency": effStatusType,
				"title_en":   "",
				"title_ar":   "",
				"count":      len(controls), // Count of efficiency status items
				"items":      []map[string]interface{}{},
			}

			if status, exists := efficiencyStatusMap[effStatusType]; exists {
				effItem["title_en"] = status.TitleEn
				effItem["title_ar"] = status.TitleAr
			}

			for _, control := range controls {
				controlItem := map[string]interface{}{
					"id":                                       control.ID,
					"asset_register_id":                        control.AssetRegisterID,
					"control_id":                               control.ControlID,
					"control_reference_number":                 control.ControlReferenceNumber,
					"control_code":                             control.ControlCode,
					"control_last_test_date":                   control.ControlLastTestDate,
					"control_category_title_en":                control.ControlCategoryTitleEn,
					"control_category_title_ar":                control.ControlCategoryTitleAr,
					"control_type_title_en":                    control.ControlTypeTitleEn,
					"control_type_title_ar":                    control.ControlTypeTitleAr,
					"control_title_en":                         control.ControlTitleEn,
					"control_title_ar":                         control.ControlTitleAr,
					"control_description":                      control.ControlDescription,
					"control_efficiency_status_title_en":       control.ControlEfficiencyStatusTitleEn,
					"control_efficiency_status_title_ar":       control.ControlEfficiencyStatusTitleAr,
					"control_efficiency_status_type":           control.ControlEfficiencyStatusType,
					"control_efficiency_status_color_code":     control.ControlEfficiencyStatusColorCode,
					"control_implementation_status_title_en":   control.ControlImplementationStatusTitleEn,
					"control_implementation_status_title_ar":   control.ControlImplementationStatusTitleAr,
					"control_implementation_status_type":       control.ControlImplementationStatusType,
					"control_implementation_status_color_code": control.ControlImplementationStatusColorCode,
					"control_owner_first_name_en":              control.ControlOwnerFirstNameEn,
					"control_owner_first_name_ar":              control.ControlOwnerFirstNameAr,
					"control_owner_last_name_en":               control.ControlOwnerLastNameEn,
					"control_owner_last_name_ar":               control.ControlOwnerLastNameAr,
					"control_owner_initial":                    control.ControlOwnerInitial,
					"control_owner_color":                      control.ControlOwnerColor,
					"control_owner_email":                      control.ControlOwnerEmail,
					"control_owner_profile_image_path":         control.ControlOwnerProfileImagePath,
				}
				effItem["items"] = append(effItem["items"].([]map[string]interface{}), controlItem)
			}

			implItem["items"] = append(implItem["items"].([]map[string]interface{}), effItem)
			implItem["count"] = implItem["count"].(int) + effItem["count"].(int)
		}

		response = append(response, implItem)
	}

	return response, nil
}

func (aais *AssetAssociatedItemService) GetAllIncidents(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemIncidentListResponseDTO, error) {
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return aais.repo.GetAllIncidents(assetID, params)
}

func (aais *AssetAssociatedItemService) GetAllRisks(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemRiskListResponseDTO, error) {
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return aais.repo.GetAllRisks(assetID, params)
}
