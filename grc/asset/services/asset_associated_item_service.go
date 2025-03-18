package services

import (
	"fmt"
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
	GetAllVulnerabilities(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) ([]assetDTOs.VulnerabilityDTO, error)
	CreateAssetVulnerabilities(assetID uint, formDTO assetDTOs.AssetVulnerabilitiesFormDTO, currentUserID uint) error
	DeleteAssetVulnerabilities(assetID uint, vulnerabilityIDs []uint, currentUserID uint) ([]uint, error)
	CreateAssetVulnerabilityThreats(vulnerabilityId uint, formDTO assetDTOs.AssetVulnerabilityThreatsFormDTO, currentUserID uint) error
	DeleteAssetVulnerabilityThreats(vulnerabilityId uint, threatIDs []uint, currentUserID uint) ([]uint, error)
}

// Service Implementation
type AssetAssociatedItemService struct {
	repo                           assetRepsitories.AssetAssociatedItemRepositoryInterface
	assetRegisterSvc               AssetRegisterServiceInterface
	controlImplementationStatusSvc organizationServices.ControlImplementationStatusServiceInterface
	controlEfficiencyStatusService organizationServices.ControlEfficiencyStatusServiceInterface
	vulnerabilitySvc               VulnerabilityServiceInterface
	threatSvc                      ThreatServiceInterface
}

// Constructor Function
func NewAssetAssociatedItemService(
	repo assetRepsitories.AssetAssociatedItemRepositoryInterface,
	assetRegisterSvc AssetRegisterServiceInterface,
	controlImplementationStatusSvc organizationServices.ControlImplementationStatusServiceInterface,
	controlEfficiencyStatusService organizationServices.ControlEfficiencyStatusServiceInterface,
	vulnerabilitySvc VulnerabilityServiceInterface,
	threatSvc ThreatServiceInterface,
) AssetAssociatedItemServiceInterface {
	return &AssetAssociatedItemService{
		repo:                           repo,
		assetRegisterSvc:               assetRegisterSvc,
		controlImplementationStatusSvc: controlImplementationStatusSvc,
		controlEfficiencyStatusService: controlEfficiencyStatusService,
		vulnerabilitySvc:               vulnerabilitySvc,
		threatSvc:                      threatSvc,
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

func (service *AssetAssociatedItemService) GetAllVulnerabilities(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) ([]assetDTOs.VulnerabilityDTO, error) {
	// Ensure the asset exists
	if _, err := service.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}

	// Fetch vulnerabilities from the repository
	vulnerabilities, err := service.repo.GetAllVulnerabilities(assetID)
	if err != nil {
		return nil, err
	}

	// Map vulnerabilities to DTOs
	vulnerabilityDTOs := make([]assetDTOs.VulnerabilityDTO, 0, len(vulnerabilities))
	for _, vuln := range vulnerabilities {
		// Map threats to DTOs
		threatDTOs := make([]assetDTOs.ThreatDTO, 0, len(vuln.Threats))
		for _, threat := range vuln.Threats {
			// Fetch risks for the current threat
			risks, err := service.repo.GetRisksByThreatID(threat.ID)
			if err != nil {
				return nil, err
			}

			// Map risks to DTOs
			riskDTOs := make([]assetDTOs.RiskDTO, 0, len(risks))
			for _, risk := range risks {
				riskDTOs = append(riskDTOs, assetDTOs.RiskDTO{
					ID:                              risk.ID,
					ReferenceNumber:                 risk.ReferenceNumber,
					Title:                           risk.Title,
					Description:                     risk.Description,
					RiskCategoryTitleEn:             risk.RiskCategoryTitleEn,
					RiskCategoryTitleAr:             risk.RiskCategoryTitleAr,
					RiskSourceTitleEn:               risk.RiskSourceTitleEn,
					RiskSourceTitleAr:               risk.RiskSourceTitleAr,
					ResponsibleUnitTitleEn:          risk.ResponsibleUnitTitleEn,
					ResponsibleUnitTitleAr:          risk.ResponsibleUnitTitleAr,
					ResponsibleUserID:               risk.ResponsibleUserID,
					ResponsibleUserFirstNameEn:      risk.ResponsibleUserFirstNameEn,
					ResponsibleUserFirstNameAr:      risk.ResponsibleUserFirstNameAr,
					ResponsibleUserLastNameEn:       risk.ResponsibleUserLastNameEn,
					ResponsibleUserLastNameAr:       risk.ResponsibleUserLastNameAr,
					ResponsibleUserInitial:          risk.ResponsibleUserInitial,
					ResponsibleUserColor:            risk.ResponsibleUserColor,
					ResponsibleUserEmail:            risk.ResponsibleUserEmail,
					ResponsibleUserProfileImagePath: risk.ResponsibleUserProfileImagePath,
				})
			}

			// Append threat with its risks
			threatDTOs = append(threatDTOs, assetDTOs.ThreatDTO{
				ID:                              threat.ID,
				TitleEn:                         threat.TitleEn,
				TitleAr:                         threat.TitleAr,
				ThreatTypeTitleEn:               &threat.ThreatType.TitleEn,
				ThreatTypeTitleAr:               &threat.ThreatType.TitleAr,
				SeverityLevelID:                 threat.SeverityLevelID,
				SeverityLevelTitleEn:            threat.SeverityLevel.TitleEn,
				SeverityLevelTitleAr:            threat.SeverityLevel.TitleAr,
				SeverityLevelColorCode:          threat.SeverityLevel.ColorCode,
				ThreatStatusID:                  threat.ThreatStatusID,
				ThreatStatusTitleEn:             threat.ThreatStatus.TitleEn,
				ThreatStatusTitleAr:             threat.ThreatStatus.TitleAr,
				ThreatStatusColorCode:           threat.ThreatStatus.ColorCode,
				ThreatStatusType:                threat.ThreatStatus.Type,
				Description:                     threat.Description,
				Source:                          threat.Source,
				DetectionDate:                   threat.DetectionDate,
				ResolutionDate:                  threat.ResolutionDate,
				ResponsibleUserID:               threat.ResponsibleUserID,
				ResponsibleUserFirstNameEn:      threat.ResponsibleUser.FirstNameEn,
				ResponsibleUserFirstNameAr:      threat.ResponsibleUser.FirstNameAr,
				ResponsibleUserLastNameEn:       threat.ResponsibleUser.LastNameEn,
				ResponsibleUserLastNameAr:       threat.ResponsibleUser.LastNameAr,
				ResponsibleUserInitial:          threat.ResponsibleUser.Initial,
				ResponsibleUserColor:            threat.ResponsibleUser.Color,
				ResponsibleUserEmail:            threat.ResponsibleUser.Email,
				ResponsibleUserProfileImagePath: threat.ResponsibleUser.ProfileImagePath,
				Risks:                           riskDTOs,
			})
		}

		var severityLevelTitleEn, severityLevelTitleAr, severityLevelColorCode *string
		if vuln.SeverityLevel != nil {
			severityLevelTitleEn = &vuln.SeverityLevel.TitleEn
			severityLevelTitleAr = &vuln.SeverityLevel.TitleAr
			severityLevelColorCode = &vuln.SeverityLevel.ColorCode
		}

		// Map vulnerability to DTO
		vulnerabilityDTOs = append(vulnerabilityDTOs, assetDTOs.VulnerabilityDTO{
			ID:                              vuln.ID,
			ReferenceNumber:                 vuln.ReferenceNumber,
			TitleEn:                         vuln.TitleEn,
			TitleAr:                         vuln.TitleAr,
			Description:                     vuln.Description,
			Source:                          vuln.Source,
			SeverityLevelID:                 vuln.SeverityLevelID,
			SeverityLevelTitleEn:            severityLevelTitleEn,
			SeverityLevelTitleAr:            severityLevelTitleAr,
			SeverityLevelColorCode:          severityLevelColorCode,
			VulnerabilityStatusID:           vuln.VulnerabilityStatusID,
			VulnerabilityStatusTitleEn:      vuln.VulnerabilityStatus.TitleEn,
			VulnerabilityStatusTitleAr:      vuln.VulnerabilityStatus.TitleAr,
			VulnerabilityStatusColorCode:    vuln.VulnerabilityStatus.ColorCode,
			VulnerabilityStatusType:         vuln.VulnerabilityStatus.Type,
			VulnerabilityTypeTitleEn:        vuln.VulnerabilityType.TitleEn,
			VulnerabilityTypeTitleAr:        vuln.VulnerabilityType.TitleAr,
			ResponsibleUserID:               vuln.ResponsibleUser.ID,
			ResponsibleUserFirstNameEn:      vuln.ResponsibleUser.FirstNameEn,
			ResponsibleUserFirstNameAr:      vuln.ResponsibleUser.FirstNameAr,
			ResponsibleUserLastNameEn:       vuln.ResponsibleUser.LastNameEn,
			ResponsibleUserLastNameAr:       vuln.ResponsibleUser.LastNameAr,
			ResponsibleUserInitial:          vuln.ResponsibleUser.Initial,
			ResponsibleUserColor:            vuln.ResponsibleUser.Color,
			ResponsibleUserEmail:            vuln.ResponsibleUser.Email,
			ResponsibleUserProfileImagePath: vuln.ResponsibleUser.ProfileImagePath,
			Threats:                         threatDTOs,
		})
	}

	return vulnerabilityDTOs, nil
}

func (aais *AssetAssociatedItemService) CreateAssetVulnerabilities(assetID uint, formDTO assetDTOs.AssetVulnerabilitiesFormDTO, currentUserID uint) error {
	tx := aais.repo.BeginTransaction()

	// Fetch asset to ensure it exists
	_, err := aais.assetRegisterSvc.GetByID(assetID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Prepare vulnerabilities for insertion
	var vulnerabilities []*models.VulnerabilityAsset
	for _, vulnerabilityID := range formDTO.VulnerabilityIds {
		vulnerability := &models.VulnerabilityAsset{
			AssetRegisterID: assetID,
			VulnerabilityID: vulnerabilityID,
			CreatedByID:     currentUserID,
		}
		vulnerabilities = append(vulnerabilities, vulnerability)
	}

	// Bulk insert vulnerabilities
	if len(vulnerabilities) > 0 {
		if err := aais.repo.CreateAssetVulnerabilities(vulnerabilities, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (aais *AssetAssociatedItemService) DeleteAssetVulnerabilities(assetID uint, vulnerabilityIDs []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint

	if len(vulnerabilityIDs) == 0 {
		return nil, nil
	}

	// Process each vulnerability for deletion
	for _, id := range vulnerabilityIDs {
		_, err := aais.vulnerabilitySvc.GetByID(id)
		if err != nil {
			fmt.Printf("Failed to fetch asset vulnerability with ID %d: %v\n", id, err)
			continue
		}

		err = aais.repo.DeleteAssetVulnerability(assetID, id)
		if err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
		}
	}

	return failedIDs, nil
}

func (aais *AssetAssociatedItemService) CreateAssetVulnerabilityThreats(vulnerabilityId uint, formDTO assetDTOs.AssetVulnerabilityThreatsFormDTO, currentUserID uint) error {
	tx := aais.repo.BeginTransaction()

	// Fetch asset to ensure it exists
	_, err := aais.vulnerabilitySvc.GetByID(vulnerabilityId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Prepare threats for insertion
	var threats []*models.VulnerabilityThreat
	for _, threatID := range formDTO.ThreatIds {
		threat := &models.VulnerabilityThreat{
			VulnerabilityID: vulnerabilityId,
			ThreatID:        threatID,
			CreatedBYID:     currentUserID,
		}
		threats = append(threats, threat)
	}

	// Bulk insert threats
	if len(threats) > 0 {
		if err := aais.repo.CreateAssetVulnerabilityThreats(threats, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (aais *AssetAssociatedItemService) DeleteAssetVulnerabilityThreats(vulnerabilityId uint, threatIDs []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint

	if len(threatIDs) == 0 {
		return nil, nil
	}

	// Process each threat for deletion
	for _, id := range threatIDs {
		_, err := aais.threatSvc.GetByID(id)
		if err != nil {
			fmt.Printf("Failed to fetch asset threat with ID %d: %v\n", id, err)
			continue
		}

		err = aais.repo.DeleteVulnerabilityThreats(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
		}
	}

	return failedIDs, nil
}
