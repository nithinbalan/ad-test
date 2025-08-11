package services

import (
	"fmt"
	assetDTOs "jethur/grc/asset/dtos"
	assetRepsitories "jethur/grc/asset/repositories"
	commonDTOs "jethur/grc/common/dtos"
	complianceServices "jethur/grc/compliance/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"time"

	commonServices "jethur/grc/common/services"
	"jethur/utility"
)

// Service Interface
type AssetAssociatedItemServiceInterface interface {
	GetAllControls(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) ([]map[string]interface{}, error)
	GetAllIncidents(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemIncidentListResponseDTO, error)
	GetAllRisks(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetAllVulnerabilities(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) ([]assetDTOs.VulnerabilityDTO, error)
	GetAllClauses(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)

	CreateAssetComplianceClauses(assetID uint, formDTO assetDTOs.AssetComplianceClauseFormDTO, currentUserID uint) error
	DeleteAssetComplianceClauses(assetID uint, ids []uint, currentUserID uint) ([]uint, error)

	CreateAssetVulnerabilities(assetID uint, formDTO assetDTOs.AssetVulnerabilitiesFormDTO, currentUserID uint) error
	DeleteAssetVulnerabilities(assetID uint, vulnerabilityIDs []uint, currentUserID uint) ([]uint, error)
	CreateAssetVulnerabilityThreats(vulnerabilityId uint, formDTO assetDTOs.AssetVulnerabilityThreatsFormDTO, currentUserID uint) error
	DeleteAssetVulnerabilityThreats(vulnerabilityId uint, threatIDs []uint, currentUserID uint) ([]uint, error)
	CreateAssetThreatRisks(threatId uint, formDTO assetDTOs.AssetThreatRisksFormDTO, currentUserID uint) error
	DeleteAssetThreatRisks(threatId uint, riskIDs []uint, currentUserID uint) ([]uint, error)
	GetRiskAnalysis(assetID uint) ([]*assetDTOs.VulnerabilityDTO, error)
	UpdateClauses(riskID uint, currentUserID uint, ClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error
	CreateAssetRisk(assetID uint, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error
	DeleteAssetRisk(assetID uint, riskID uint) error
	GetAllAssetDocuments(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateAssetDocument(assetID uint, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteAssetDocument(assetID uint, documentID uint) error
	CreateAssetControl(assetID uint, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error
	DeleteAssetControl(assetID uint, controlID uint) error
	GetAllAssetProcess(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	CreateAssetProcess(assetID uint, userID uint, dto *commonDTOs.AssociatedProcessMapDTO) error
	DeleteAssetProcess(assetID uint, processID uint) error
	GetAllAssetStrategicObjectives(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error)
	CreateAssetStrategicObjective(assetID uint, userID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error
	DeleteAssetStrategicObjective(assetID uint, strategicObjectiveId uint) error
}

// Service Implementation
type AssetAssociatedItemService struct {
	repo                           assetRepsitories.AssetAssociatedItemRepositoryInterface
	assetRegisterSvc               AssetRegisterServiceInterface
	controlImplementationStatusSvc organizationServices.ControlImplementationStatusServiceInterface
	controlEfficiencyStatusService organizationServices.ControlEfficiencyStatusServiceInterface
	vulnerabilitySvc               VulnerabilityServiceInterface
	threatSvc                      ThreatServiceInterface
	complianceClauseSvc            complianceServices.ComplianceClauseServiceInterface
	validationSvc                  commonServices.AssociatedItemValidationServiceInterface
}

// Constructor Function
func NewAssetAssociatedItemService(
	repo assetRepsitories.AssetAssociatedItemRepositoryInterface,
	assetRegisterSvc AssetRegisterServiceInterface,
	controlImplementationStatusSvc organizationServices.ControlImplementationStatusServiceInterface,
	controlEfficiencyStatusService organizationServices.ControlEfficiencyStatusServiceInterface,
	vulnerabilitySvc VulnerabilityServiceInterface,
	threatSvc ThreatServiceInterface,
	complianceClauseSvc complianceServices.ComplianceClauseServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) AssetAssociatedItemServiceInterface {
	return &AssetAssociatedItemService{
		repo:                           repo,
		assetRegisterSvc:               assetRegisterSvc,
		controlImplementationStatusSvc: controlImplementationStatusSvc,
		controlEfficiencyStatusService: controlEfficiencyStatusService,
		vulnerabilitySvc:               vulnerabilitySvc,
		threatSvc:                      threatSvc,
		complianceClauseSvc:            complianceClauseSvc,
		validationSvc:                  validationSvc,
	}
}

func (aais *AssetAssociatedItemService) GetAllControls(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) ([]map[string]interface{}, error) {
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

func (aais *AssetAssociatedItemService) GetAllRisks(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	// Validate asset
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}

	// Get related vulnerability IDs
	vulnerabilityIDs, err := aais.repo.GetVulnerabilityIDsByAssetID(assetID)
	if err != nil {
		return nil, err
	}
	if len(vulnerabilityIDs) == 0 {
		return emptyRiskListResponse(params), nil
	}

	// Get related threat IDs
	threatIDs, err := aais.repo.GetThreatIDsByVulnerabilityIDs(vulnerabilityIDs)
	if err != nil {
		return nil, err
	}
	if len(threatIDs) == 0 {
		return emptyRiskListResponse(params), nil
	}

	// Get related risk IDs
	riskRegisterIDs, err := aais.repo.GetRiskRegisterIDsByThreatIDs(threatIDs)
	if err != nil {
		return nil, err
	}
	if len(riskRegisterIDs) == 0 {
		return emptyRiskListResponse(params), nil
	}

	// Final call to fetch filtered risks
	return aais.repo.GetAllRisksByRiskIDs(riskRegisterIDs, params)
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
		// threatDTOs := make([]assetDTOs.ThreatDTO, 0, len(vuln.Threats))
		// for _, threat := range vuln.Threats {
		// 	// Fetch risks for the current threat
		// 	risks, err := service.repo.GetRisksByThreatID(threat.ID)
		// 	if err != nil {
		// 		return nil, err
		// 	}

		// 	// Map risks to DTOs
		// 	riskDTOs := make([]assetDTOs.RiskDTO, 0, len(risks))
		// 	for _, risk := range risks {
		// 		riskDTOs = append(riskDTOs, assetDTOs.RiskDTO{
		// 			ID:                              risk.ID,
		// 			ReferenceNumber:                 risk.ReferenceNumber,
		// 			Title:                           risk.Title,
		// 			Description:                     risk.Description,
		// 			RiskCategoryTitleEn:             risk.RiskCategoryTitleEn,
		// 			RiskCategoryTitleAr:             risk.RiskCategoryTitleAr,
		// 			RiskSourceTitleEn:               risk.RiskSourceTitleEn,
		// 			RiskSourceTitleAr:               risk.RiskSourceTitleAr,
		// 			ResponsibleUnitTitleEn:          risk.ResponsibleUnitTitleEn,
		// 			ResponsibleUnitTitleAr:          risk.ResponsibleUnitTitleAr,
		// 			ResponsibleUserID:               risk.ResponsibleUserID,
		// 			ResponsibleUserFirstNameEn:      risk.ResponsibleUserFirstNameEn,
		// 			ResponsibleUserFirstNameAr:      risk.ResponsibleUserFirstNameAr,
		// 			ResponsibleUserLastNameEn:       risk.ResponsibleUserLastNameEn,
		// 			ResponsibleUserLastNameAr:       risk.ResponsibleUserLastNameAr,
		// 			ResponsibleUserInitial:          risk.ResponsibleUserInitial,
		// 			ResponsibleUserColor:            risk.ResponsibleUserColor,
		// 			ResponsibleUserEmail:            risk.ResponsibleUserEmail,
		// 			ResponsibleUserProfileImagePath: risk.ResponsibleUserProfileImagePath,
		// 		})
		// 	}

		// 	// Append threat with its risks
		// 	threatDTOs = append(threatDTOs, assetDTOs.ThreatDTO{
		// 		ID:                              threat.ID,
		// 		TitleEn:                         threat.TitleEn,
		// 		TitleAr:                         threat.TitleAr,
		// 		ThreatTypeTitleEn:               &threat.ThreatType.TitleEn,
		// 		ThreatTypeTitleAr:               &threat.ThreatType.TitleAr,
		// 		SeverityLevelID:                 threat.SeverityLevelID,
		// 		SeverityLevelTitleEn:            threat.SeverityLevel.TitleEn,
		// 		SeverityLevelTitleAr:            threat.SeverityLevel.TitleAr,
		// 		SeverityLevelColorCode:          threat.SeverityLevel.ColorCode,
		// 		ThreatStatusID:                  threat.ThreatStatusID,
		// 		ThreatStatusTitleEn:             threat.ThreatStatus.TitleEn,
		// 		ThreatStatusTitleAr:             threat.ThreatStatus.TitleAr,
		// 		ThreatStatusColorCode:           threat.ThreatStatus.ColorCode,
		// 		ThreatStatusType:                threat.ThreatStatus.Type,
		// 		Description:                     threat.Description,
		// 		Source:                          threat.Source,
		// 		DetectionDate:                   threat.DetectionDate,
		// 		ResolutionDate:                  threat.ResolutionDate,
		// 		ResponsibleUserID:               threat.ResponsibleUserID,
		// 		ResponsibleUserFirstNameEn:      threat.ResponsibleUser.FirstNameEn,
		// 		ResponsibleUserFirstNameAr:      threat.ResponsibleUser.FirstNameAr,
		// 		ResponsibleUserLastNameEn:       threat.ResponsibleUser.LastNameEn,
		// 		ResponsibleUserLastNameAr:       threat.ResponsibleUser.LastNameAr,
		// 		ResponsibleUserInitial:          threat.ResponsibleUser.Initial,
		// 		ResponsibleUserColor:            threat.ResponsibleUser.Color,
		// 		ResponsibleUserEmail:            threat.ResponsibleUser.Email,
		// 		ResponsibleUserProfileImagePath: threat.ResponsibleUser.ProfileImagePath,
		// 		Risks:                           riskDTOs,
		// 	})
		// }

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
			// Threats:                         threatDTOs,
		})
	}

	return vulnerabilityDTOs, nil
}

func (aais *AssetAssociatedItemService) GetAllClauses(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return aais.repo.GetAllClauses(assetID, params)
}

func (aais *AssetAssociatedItemService) CreateAssetVulnerabilities(assetID uint, formDTO assetDTOs.AssetVulnerabilitiesFormDTO, currentUserID uint) error {
	// Fetch asset to ensure it exists
	_, err := aais.assetRegisterSvc.GetByID(assetID)
	if err != nil {
		return err
	}
	var existingVulnerabilities []uint

	for _, vulnerabilityID := range formDTO.VulnerabilityIds {
		// Ensure vulnerability exists
		_, err := aais.vulnerabilitySvc.GetByID(vulnerabilityID)
		if err != nil {
			return err
		}

		// ✅ Check if asset-vulnerability relationship already exists
		exists, err := aais.repo.GetVulnerabilityAssetExist(assetID, vulnerabilityID)
		if err != nil {
			return err
		}
		if exists {
			existingVulnerabilities = append(existingVulnerabilities, vulnerabilityID)
		}
	}

	// ✅ If any existing vulnerabilities were found, return validation error & exit
	if len(existingVulnerabilities) > 0 {
		return utility.NewValidationError(
			fmt.Errorf("vulnerabilities already mapped"),
			"already_exists", "vulnerability_ids",
		)
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
	tx := aais.repo.BeginTransaction()

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

func (aais *AssetAssociatedItemService) CreateAssetComplianceClauses(assetID uint, formDTO assetDTOs.AssetComplianceClauseFormDTO, currentUserID uint) error {
	// Check if asset exists
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return err
	}

	var existingClauseIDs []uint
	var records []*models.AssetComplianceRequirement

	for _, clauseID := range formDTO.ComplianceClauseIDs {
		if clauseID == nil {
			continue
		}

		// Validate clause exists
		if _, err := aais.complianceClauseSvc.GetByID(*clauseID); err != nil {
			return err
		}

		// Check if mapping already exists
		exists, err := aais.repo.GetAssetComplianceClauseExist(assetID, *clauseID)
		if err != nil {
			return err
		}
		if exists {
			existingClauseIDs = append(existingClauseIDs, *clauseID)
			continue
		}

		// Add to batch for insert
		records = append(records, &models.AssetComplianceRequirement{
			AssetRegisterID:    assetID,
			ComplianceClauseID: clauseID,
			CreatedBYID:        currentUserID,
		})
	}

	// If any existing clause mappings found, return validation error
	if len(existingClauseIDs) > 0 {
		return utility.NewValidationError(
			fmt.Errorf("compliance clauses already mapped"),
			"already_exists", "compliance_clause_ids",
		)
	}

	// Insert new records
	tx := aais.repo.BeginTransaction()
	if err := aais.repo.CreateAssetComplianceClauses(records, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (aais *AssetAssociatedItemService) DeleteAssetComplianceClauses(assetID uint, ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint

	if len(ids) == 0 {
		return nil, nil
	}

	for _, id := range ids {
		// Optionally validate record exists
		_, err := aais.repo.GetAssetComplianceClauseByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if err := aais.repo.DeleteAssetComplianceClause(assetID, id); err != nil {
			failedIDs = append(failedIDs, id)
		}
	}

	return failedIDs, nil
}

func (aais *AssetAssociatedItemService) CreateAssetVulnerabilityThreats(vulnerabilityId uint, formDTO assetDTOs.AssetVulnerabilityThreatsFormDTO, currentUserID uint) error {

	// Fetch asset to ensure it exists
	_, err := aais.vulnerabilitySvc.GetByID(vulnerabilityId)
	if err != nil {
		return err
	}
	var existingThreats []uint

	for _, threatID := range formDTO.ThreatIds {
		_, err := aais.threatSvc.GetByID(threatID)
		if err != nil {
			return err
		}

		// Check if threat is already mapped to the vulnerability
		exists, err := aais.repo.GetVulnerabilityThreatExist(vulnerabilityId, threatID)
		if err != nil {
			return err
		}
		if exists {
			existingThreats = append(existingThreats, threatID)
		}
	}

	// ✅ If any threats are already mapped, return validation error
	if len(existingThreats) > 0 {
		return utility.NewValidationError(
			fmt.Errorf("threats already mapped"),
			"already_exists", "threat_ids",
		)
	}
	// Prepare threats for insertion
	var threats []*models.VulnerabilityThreat
	for _, threatID := range formDTO.ThreatIds {
		_, err := aais.threatSvc.GetByID(threatID)
		if err != nil {
			return err
		}
		threat := &models.VulnerabilityThreat{
			VulnerabilityID: vulnerabilityId,
			ThreatID:        threatID,
			CreatedBYID:     currentUserID,
		}
		threats = append(threats, threat)
	}
	tx := aais.repo.BeginTransaction()

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
		err = aais.repo.DeleteVulnerabilityThreats(vulnerabilityId, id)
		if err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
		}
	}

	return failedIDs, nil
}

func (aais *AssetAssociatedItemService) CreateAssetThreatRisks(threatId uint, formDTO assetDTOs.AssetThreatRisksFormDTO, currentUserID uint) error {

	// Fetch asset to ensure it exists
	_, err := aais.threatSvc.GetByID(threatId)
	if err != nil {
		return err
	}

	var existingRisks []uint

	for _, riskID := range formDTO.RiskIds {
		_, err := aais.repo.GetByRiskID(riskID)
		if err != nil {
			return err
		}

		// Check if risk is already mapped to the threat
		exists, err := aais.repo.GetThreatRiskExist(threatId, riskID)
		if err != nil {
			return err
		}
		if exists {
			existingRisks = append(existingRisks, riskID)
		}
	}

	// ✅ If any risks are already mapped, return validation error
	if len(existingRisks) > 0 {
		return utility.NewValidationError(
			fmt.Errorf("risks already mapped"),
			"already_exists", "risk_ids",
		)
	}
	// Prepare threats for insertion
	var risks []*models.RiskThreat
	for _, riskID := range formDTO.RiskIds {
		risk := &models.RiskThreat{
			ThreatID:       threatId,
			RiskRegisterID: riskID,
			CreatedByID:    currentUserID,
		}
		risks = append(risks, risk)
	}
	tx := aais.repo.BeginTransaction()

	// Bulk insert risks
	if len(risks) > 0 {
		if err := aais.repo.CreateAssetThreatRisks(risks, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (aais *AssetAssociatedItemService) DeleteAssetThreatRisks(threatId uint, riskIDs []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint

	if len(riskIDs) == 0 {
		return nil, nil
	}

	// Process each risk for deletion
	for _, id := range riskIDs {
		_, err := aais.repo.GetByRiskID(id)
		if err != nil {
			fmt.Printf("Failed to fetch asset risk with ID %d: %v\n", id, err)
			continue
		}

		err = aais.repo.DeleteAssetThreatRisks(threatId, id)
		if err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
		}
	}

	return failedIDs, nil
}

func (aais *AssetAssociatedItemService) GetRiskAnalysis(assetID uint) ([]*assetDTOs.VulnerabilityDTO, error) {
	// Ensure asset exists before fetching risk analysis
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}

	// Fetch vulnerabilities, threats, and risks from the repository
	vulnerabilities, threats, risks, err := aais.repo.GetAssetRiskAnalysis(assetID)
	if err != nil {
		return nil, err
	}

	// Step 1: Create a map for Risks indexed by Threat ID
	riskMap := make(map[uint][]assetDTOs.RiskDTO)

	for _, risk := range risks {
		riskMap[risk.ThreatID] = append(riskMap[risk.ThreatID], *risk)
		// riskMap[risk.ID] = append(riskMap[risk.ID], *risk) // Store risks under each Threat ID
	}

	// Step 2: Create a map for Threats indexed by Vulnerability ID
	threatMap := make(map[uint][]assetDTOs.ThreatDTO)

	for _, threat := range threats {
		// Assign risks to each threat before adding to threatMap
		if mappedRisks, exists := riskMap[threat.ID]; exists {
			threat.Risks = mappedRisks // Attach risks to the threat
		} else {
			threat.Risks = []assetDTOs.RiskDTO{} // Ensure empty array
		}

		threatMap[threat.VulnerabilityID] = append(threatMap[threat.VulnerabilityID], *threat)
	}

	// Step 3: Assign threats (with risks) to vulnerabilities
	for _, vulnerability := range vulnerabilities {
		if mappedThreats, exists := threatMap[vulnerability.ID]; exists {
			vulnerability.Threats = mappedThreats // Attach threats to vulnerability
		} else {
			vulnerability.Threats = []assetDTOs.ThreatDTO{} // Ensure empty array
		}
	}

	// Return the structured response
	return vulnerabilities, nil
}

func emptyRiskListResponse(params *commonDTOs.AssociatedItemListQueryParamsDTO) *commonDTOs.AssociatedRiskListResponseDTO {
	return &commonDTOs.AssociatedRiskListResponseDTO{
		Items: []*commonDTOs.AssociatedRiskListDTO{},
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  0,
			TotalCount:  0,
		},
	}
}

func (aais *AssetAssociatedItemService) UpdateClauses(assetID uint, currentUserID uint, ClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error {

	tx := aais.repo.BeginTransaction()

	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return err
	}

	clauses, err := aais.validateClauseIDs(ClauseFormDTO.ClauseIDs, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingControlClauses, err := aais.repo.GetClausesByAssetId(assetID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingClauses := make(map[uint]bool) // clauseID => exists
	for _, clauseID := range existingControlClauses {
		existingClauses[clauseID] = true
	}

	newClauses := make(map[uint]*models.ComplianceClause)
	for _, clause := range clauses {
		newClauses[clause.ID] = clause
	}

	var clausesToAdd []*models.AssetComplianceRequirement
	now := time.Now()
	for id, clause := range newClauses {
		if !existingClauses[id] {
			clausesToAdd = append(clausesToAdd, &models.AssetComplianceRequirement{
				AssetRegisterID:      assetID,
				ComplianceClauseID:   &clause.ID,
				ComplianceRegisterID: clause.ComplianceRegisterID,
				CreatedBYID:          currentUserID,
				CreatedAt:            now,
			})
		}
	}

	if len(clausesToAdd) > 0 {
		if err := aais.repo.CreateAssetClauses(clausesToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	var deleteClauseIDs []uint
	for existingID := range existingClauses {
		if _, found := newClauses[existingID]; !found {
			deleteClauseIDs = append(deleteClauseIDs, existingID)
		}
	}

	if len(deleteClauseIDs) > 0 {
		if err := aais.repo.DeleteAssetClauses(assetID, deleteClauseIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (aais *AssetAssociatedItemService) validateClauseIDs(clauseIDs []uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceClause, error) {
	seen := make(map[uint]bool)
	for _, id := range clauseIDs {
		if seen[id] {
			return nil, utility.NewValidationError(fmt.Errorf("duplicate clause ID: %d", id), "duplicate", "clause_id")
		}
		seen[id] = true
	}

	clauses, err := aais.complianceClauseSvc.GetClausesByIDs(clauseIDs, tx...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch clauses: %w", err)
	}

	foundClauseIDs := make(map[uint]bool)
	for _, clause := range clauses {
		foundClauseIDs[clause.ID] = true
	}

	for _, id := range clauseIDs {
		if !foundClauseIDs[id] {
			return nil, utility.NewValidationError(
				fmt.Errorf("clause does not exist for ID: %d", id), "not_exist", "clause_id")
		}
	}

	return clauses, nil
}

func (aais *AssetAssociatedItemService) CreateAssetRisk(assetID, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error {

	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset does not exist for ID: %d", assetID), "not_exist", "asset_id")
	}

	if err := aais.validationSvc.ValidateRiskMapping(assetID, dto.RiskIDs, "risk_assets", "asset_register_id"); err != nil {
		return err
	}

	tx := aais.repo.BeginTransaction()
	var mappings []*models.RiskAsset
	currentTime := time.Now()

	for _, riskID := range dto.RiskIDs {
		mappings = append(mappings, &models.RiskAsset{
			RiskRegisterID:  riskID,
			AssetRegisterID: assetID,
			CreatedByID:     userID,
			CreatedAt:       currentTime,
		})
	}

	if err := aais.repo.CreateAssetRisk(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (aais *AssetAssociatedItemService) DeleteAssetRisk(assetID uint, riskID uint) error {
	// Validate control existence
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset does not exist for ID: %d", assetID), "not_exist", "asset_register_id")
	}

	tx := aais.repo.BeginTransaction()

	if err := aais.repo.DeleteAssetRisk(assetID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (aais *AssetAssociatedItemService) GetAllAssetDocuments(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return aais.repo.GetAllAssetDocuments(assetID, params)
}

func (aais *AssetAssociatedItemService) CreateAssetDocument(assetID, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {

	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset does not exist for ID: %d", assetID), "not_exist", "asset_id")
	}

	if err := aais.validationSvc.ValidateDocumentMapping(assetID, dto.DocumentIDs, "asset_associated_documents", "asset_register_id"); err != nil {
		return err
	}

	tx := aais.repo.BeginTransaction()
	var mappings []*models.AssetAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DocumentIDs {
		mappings = append(mappings, &models.AssetAssociatedDocument{
			DocumentID:      documentID,
			AssetRegisterID: assetID,
			CreatedByID:     userID,
			CreatedAt:       currentTime,
		})
	}

	if err := aais.repo.CreateAssetDocument(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (aais *AssetAssociatedItemService) DeleteAssetDocument(assetID uint, documentID uint) error {
	// Validate control existence
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset does not exist for ID: %d", assetID), "not_exist", "asset_register_id")
	}

	tx := aais.repo.BeginTransaction()

	if err := aais.repo.DeleteAssetDocument(assetID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (aais *AssetAssociatedItemService) CreateAssetControl(assetID, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error {

	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset does not exist for ID: %d", assetID), "not_exist", "asset_id")
	}

	if err := aais.validationSvc.ValidateControlMapping(assetID, dto.ControlIDs, "asset_controls", "asset_register_id"); err != nil {
		return err
	}

	tx := aais.repo.BeginTransaction()
	var mappings []*models.AssetControl
	currentTime := time.Now()

	for _, controlID := range dto.ControlIDs {
		mappings = append(mappings, &models.AssetControl{
			ControlID:       controlID,
			AssetRegisterID: assetID,
			CreatedBYID:     userID,
			CreatedAt:       currentTime,
		})
	}

	if err := aais.repo.CreateAssetControl(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (aais *AssetAssociatedItemService) DeleteAssetControl(assetID uint, controlID uint) error {
	// Validate control existence
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset does not exist for ID: %d", assetID), "not_exist", "asset_register_id")
	}

	tx := aais.repo.BeginTransaction()

	if err := aais.repo.DeleteAssetControl(assetID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (aais *AssetAssociatedItemService) GetAllAssetProcess(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return aais.repo.GetAllAssetProcess(assetID, params)
}

func (aais *AssetAssociatedItemService) CreateAssetProcess(assetID, userID uint, dto *commonDTOs.AssociatedProcessMapDTO) error {

	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset does not exist for ID: %d", assetID), "not_exist", "asset_id")
	}

	if err := aais.validationSvc.ValidateProcessMapping(assetID, dto.ProcessIDs, "process_assets", "asset_register_id"); err != nil {
		return err
	}

	tx := aais.repo.BeginTransaction()
	var mappings []*models.ProcessAsset
	currentTime := time.Now()

	for _, processID := range dto.ProcessIDs {
		mappings = append(mappings, &models.ProcessAsset{
			ProcessID:       processID,
			AssetRegisterID: assetID,
			CreatedByID:     userID,
			CreatedAt:       currentTime,
		})
	}

	if err := aais.repo.CreateAssetProcess(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (aais *AssetAssociatedItemService) DeleteAssetProcess(assetID uint, processID uint) error {
	// Validate control existence
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset does not exist for ID: %d", assetID), "not_exist", "asset_register_id")
	}

	tx := aais.repo.BeginTransaction()

	if err := aais.repo.DeleteAssetProcess(assetID, processID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (aais *AssetAssociatedItemService) GetAllAssetStrategicObjectives(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error) {
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return aais.repo.GetAllAssetStrategicObjectives(assetID, params)
}

func (aais *AssetAssociatedItemService) CreateAssetStrategicObjective(assetID, userID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error {

	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset does not exist for ID: %d", assetID), "not_exist", "asset_id")
	}

	if err := aais.validationSvc.ValidateStrtegicObjectiveMapping(assetID, dto.StrategicObjectiveIDs, "asset_strategic_objectives", "asset_register_id"); err != nil {
		return err
	}

	tx := aais.repo.BeginTransaction()
	var mappings []*models.AssetStrategicObjective
	currentTime := time.Now()

	for _, StrategicObjectiveID := range dto.StrategicObjectiveIDs {
		mappings = append(mappings, &models.AssetStrategicObjective{
			StrategicObjectiveID: StrategicObjectiveID,
			AssetRegisterID:      assetID,
			CreatedBYID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := aais.repo.CreateAssetStrategicObjective(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (aais *AssetAssociatedItemService) DeleteAssetStrategicObjective(assetID uint, strategicObjectiveID uint) error {
	// Validate control existence
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset does not exist for ID: %d", assetID), "not_exist", "asset_register_id")
	}

	tx := aais.repo.BeginTransaction()

	if err := aais.repo.DeleteAssetStrategicObjective(assetID, strategicObjectiveID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
