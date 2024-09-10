package services

import (
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
)

type ProcessAssociatedItemServiceInterface interface {
	GetAllControls(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]map[string]interface{}, error)
	GetAllRisks(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemRiskListResponseDTO, error)
	GetAllIncidents(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemIncidentListResponseDTO, error)
	GetAllDocuments(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemDocumentListResponseDTO, error)
	GetAllAssociatedCounts(processID uint) (*organizationDTOs.ProcessAssociatedItemsCountResponseDTO, error)
}

type ProcessAssociatedItemService struct {
	repo                           organizationRepositories.ProcessAssociatedItemRepositoryInterface
	processSvc                     ProcessServiceInterface
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface
	controlEfficiencyStatusService ControlEfficiencyStatusServiceInterface
}

func NewProcessAssociatedItemService(
	repo organizationRepositories.ProcessAssociatedItemRepositoryInterface,
	processSvc ProcessServiceInterface,
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface,
	controlEfficiencyStatusService ControlEfficiencyStatusServiceInterface,
) ProcessAssociatedItemServiceInterface {
	return &ProcessAssociatedItemService{
		repo:                           repo,
		processSvc:                     processSvc,
		controlImplementationStatusSvc: controlImplementationStatusSvc,
		controlEfficiencyStatusService: controlEfficiencyStatusService,
	}
}

func (pais *ProcessAssociatedItemService) GetAllControls(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]map[string]interface{}, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}

	controlList, err := pais.repo.GetAllControls(processID, params)
	if err != nil {
		return nil, err
	}

	if len(controlList) == 0 {
		return []map[string]interface{}{}, nil
	}

	implementationStatuses, err := pais.controlImplementationStatusSvc.GetAll("", "", "", "")
	if err != nil {
		return nil, err
	}
	implementationStatusMap := make(map[string]models.ControlImplementationStatus)
	for _, status := range implementationStatuses {
		implementationStatusMap[status.Type] = *status
	}

	efficiencyStatuses, err := pais.controlEfficiencyStatusService.GetAll("", "", "", "")
	if err != nil {
		return nil, err
	}
	efficiencyStatusMap := make(map[string]models.ControlEfficiencyStatus)
	for _, status := range efficiencyStatuses {
		efficiencyStatusMap[status.Type] = *status
	}

	groupedControls := make(map[string]map[string][]*organizationDTOs.ProcessControlListDTO)
	for _, control := range controlList {
		implementationStatusType := control.ControlImplementationStatusType
		efficiencyStatusType := control.ControlEfficiencyStatusType

		if _, exists := groupedControls[implementationStatusType]; !exists {
			groupedControls[implementationStatusType] = make(map[string][]*organizationDTOs.ProcessControlListDTO)
		}
		if _, exists := groupedControls[implementationStatusType][efficiencyStatusType]; !exists {
			groupedControls[implementationStatusType][efficiencyStatusType] = []*organizationDTOs.ProcessControlListDTO{}
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
					"process_id":                               control.ProcessID,
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

func (pais *ProcessAssociatedItemService) GetAllRisks(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemRiskListResponseDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllRisks(processID, params)
}

func (pais *ProcessAssociatedItemService) GetAllIncidents(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemIncidentListResponseDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllIncidents(processID, params)
}

func (pais *ProcessAssociatedItemService) GetAllDocuments(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemDocumentListResponseDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllDocuments(processID, params)
}

func (pais *ProcessAssociatedItemService) GetAllAssociatedCounts(processID uint) (*organizationDTOs.ProcessAssociatedItemsCountResponseDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	response, err := pais.repo.GetAllAssociatedCounts(processID)
	if err != nil {
		return nil, err
	}

	return response, nil
}
