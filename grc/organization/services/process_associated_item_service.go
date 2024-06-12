package services

import (
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
)

type ProcessAssociatedItemServiceInterface interface {
	GetAllControls(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]map[string]interface{}, error)
	GetAllRisks(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessRiskListDTO, error)
	GetAllIncidents(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessIncidentListDTO, error)
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
	// Check if the process exists
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}

	// Fetch the controls from the repository
	controlList, err := pais.repo.GetAllControls(processID, params)
	if err != nil {
		return nil, err
	}

	// Fetching implementation status titles in both English and Arabic
	implementationStatuses, err := pais.controlImplementationStatusSvc.GetAll("", "", "", "")
	if err != nil {
		return nil, err
	}
	implementationStatusMap := make(map[string]models.ControlImplementationStatus)
	for _, status := range implementationStatuses {
		implementationStatusMap[status.Type] = *status
	}

	// Fetching efficiency status titles in both English and Arabic
	efficiencyStatuses, err := pais.controlEfficiencyStatusService.GetAll("", "", "", "")
	if err != nil {
		return nil, err
	}
	efficiencyStatusMap := make(map[string]models.ControlEfficiencyStatus)
	for _, status := range efficiencyStatuses {
		efficiencyStatusMap[status.Type] = *status
	}

	// Grouping controls by implementation status type and then by efficiency status type
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

	// Transforming the grouped controls to the desired response format
	var response []map[string]interface{}
	for implStatusType, effMap := range groupedControls {
		implItem := map[string]interface{}{
			"status":   implStatusType,
			"title_en": "",
			"title_ar": "", // These will be filled in the next step
			"count":    0,  // Count of implementation status items
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
				"title_ar":   "",            // These will be filled in the next step
				"count":      len(controls), // Count of efficiency status items
				"items":      []map[string]interface{}{},
			}

			if status, exists := efficiencyStatusMap[effStatusType]; exists {
				effItem["title_en"] = status.TitleEn
				effItem["title_ar"] = status.TitleAr
			}

			for _, control := range controls {
				controlItem := map[string]interface{}{
					"control_id":                       control.ControlID,
					"control_category_title_en":        control.ControlCategoryTitleEn,
					"control_category_title_ar":        control.ControlCategoryTitleAr,
					"control_type_title_en":            control.ControlTypeTitleEn,
					"control_type_title_ar":            control.ControlTypeTitleAr,
					"control_title_en":                 control.ControlTitleEn,
					"control_title_ar":                 control.ControlTitleAr,
					"control_description":              control.ControlDescription,
					"control_efficiency_status_en":     control.ControlEfficiencyStatusEn,
					"control_efficiency_status_ar":     control.ControlEfficiencyStatusAr,
					"control_implementation_status_en": control.ControlImplementationStatusEn,
					"control_implementation_status_ar": control.ControlImplementationStatusAr,
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

func (pais *ProcessAssociatedItemService) GetAllRisks(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessRiskListDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllRisks(processID, params)
}

func (pais *ProcessAssociatedItemService) GetAllIncidents(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessIncidentListDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllIncidents(processID, params)
}
