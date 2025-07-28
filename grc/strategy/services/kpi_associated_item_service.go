package services

import (
	"fmt"
	assetServices "jethur/grc/asset/services"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceServices "jethur/grc/compliance/services"
	documentServices "jethur/grc/documents/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskServices "jethur/grc/risk/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyRepositories "jethur/grc/strategy/repositories"
	tprmServices "jethur/grc/tprm/services"
	"jethur/utility"
	"time"
)

type KPIAssociatedItemServiceInterface interface {
	GetAllControls(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemControlListResponseDTO, error)
	CreateKPIControls(kpiID, userID uint, dto *strategyDTOs.KPIControlMapDTO) error
	DeleteKPIControl(kpiID, controlID uint) error

	GetAllProcesses(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemProcessListResponseDTO, error)
	CreateKPIProcesses(kpiID, userID uint, dto *strategyDTOs.KPIProcessMapDTO) error
	DeleteKPIProcess(kpiID, processID uint) error

	GetAllCompliances(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemComplianceListResponseDTO, error)
	CreateKPICompliances(kpiID, userID uint, dto *strategyDTOs.KPIComplianceMapDTO) error
	DeleteKPICompliance(kpiID, compianceID uint) error

	GetAllAssets(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemAssetListResponseDTO, error)
	CreateKPIAssets(kpiID, userID uint, dto *strategyDTOs.KPIAssetMapDTO) error
	DeleteKPIAsset(kpiID, assetID uint) error

	GetAllRisks(kpiID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	CreateKPIRisks(kpiID, userID uint, dto *strategyDTOs.KPIRiskMapDTO) error
	DeleteKPIRisk(kpiID, riskID uint) error

	GetAllThirdParties(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemThirdPartyListResponseDTO, error)
	CreateKPIThirdParties(kpiID, userID uint, dto *strategyDTOs.KPIThirdPartyMapDTO) error
	DeleteKPIThirdParty(kpiID, thirdPartyRegisterID uint) error

	GetAllPolicies(kpiID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateKPIPolicies(kpiID, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteKPIPolicy(kpiID, documentID uint) error

	GetAllClauses(kpiID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	UpdateClauses(kpiID uint, currentUserID uint, ClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error
}

// Service and implementations
type KPIAssociatedItemService struct {
	repo           strategyRepositories.KPIAssociatedItemRepositoryInterface
	kpiRegisterSvc KPIRegisterServiceInterface
	controlSvc     organizationServices.ControlServiceInterface
	processSvc     organizationServices.ProcessServiceInterface
	assetSvc       assetServices.AssetRegisterServiceInterface
	riskSvc        riskServices.RiskRegisterServiceInterface
	thirdPartySvc  tprmServices.ThirdPartyRegisterServiceInterface
	complianceSvc  complianceServices.ComplianceRegisterServiceInterface
	documentSvc    documentServices.DocumentServiceInterface
	clauseSvc      complianceServices.ComplianceClauseServiceInterface
	validationSvc  commonServices.AssociatedItemValidationServiceInterface
}

func NewKPIAssociatedItemService(
	repo strategyRepositories.KPIAssociatedItemRepositoryInterface,
	kpiRegisterSvc KPIRegisterServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
	processSvc organizationServices.ProcessServiceInterface,
	assetSvc assetServices.AssetRegisterServiceInterface,
	riskSvc riskServices.RiskRegisterServiceInterface,
	thirdPartySvc tprmServices.ThirdPartyRegisterServiceInterface,
	complianceSvc complianceServices.ComplianceRegisterServiceInterface,
	documentSvc documentServices.DocumentServiceInterface,
	clauseSvc complianceServices.ComplianceClauseServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) KPIAssociatedItemServiceInterface {
	return &KPIAssociatedItemService{
		repo:           repo,
		kpiRegisterSvc: kpiRegisterSvc,
		controlSvc:     controlSvc,
		processSvc:     processSvc,
		assetSvc:       assetSvc,
		riskSvc:        riskSvc,
		thirdPartySvc:  thirdPartySvc,
		complianceSvc:  complianceSvc,
		documentSvc:    documentSvc,
		clauseSvc:      clauseSvc,
		validationSvc:  validationSvc,
	}
}

// controls
func (kais *KPIAssociatedItemService) GetAllControls(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemControlListResponseDTO, error) {
	if _, err := kais.kpiRegisterSvc.GetByID(kpiID); err != nil {
		return nil, err
	}
	return kais.repo.GetAllControls(kpiID, params)
}

func (kais *KPIAssociatedItemService) CreateKPIControls(kpiID, userID uint, dto *strategyDTOs.KPIControlMapDTO) error {
	if err := kais.validateControlIDs(kpiID, dto.ControlIDs); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()
	var mappings []*models.KPIControl
	currentTime := time.Now()

	for _, controlID := range dto.ControlIDs {
		mappings = append(mappings, &models.KPIControl{
			KPIRegisterID: kpiID,
			ControlID:     controlID,
			CreatedByID:   userID,
			CreatedAt:     currentTime,
		})
	}

	if err := kais.repo.CreateKPIControls(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (kais *KPIAssociatedItemService) DeleteKPIControl(kpiID, controlID uint) error {
	// Validate control existence
	if _, err := kais.controlSvc.GetByID(controlID); err != nil {
		return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", controlID), "not_exist", "control_id")
	}

	tx := kais.repo.BeginTransaction()

	if err := kais.repo.DeleteKPIControlMapping(kpiID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (kais *KPIAssociatedItemService) validateControlIDs(kpiID uint, controlIDs []uint) error {
	if len(controlIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no controls provided"), "required", "control_ids")
	}

	seen := make(map[uint]bool)
	for _, id := range controlIDs {
		alreadyMapped, err := kais.repo.IsControlMappedToKPI(kpiID, id)
		if err != nil {
			return fmt.Errorf("failed to check existing mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(fmt.Errorf("control ID %d already mapped to this kpi", id), "duplicate", "control_ids")
		}
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate control ID: %d", id), "duplicate", "control_ids")
		}
		if _, err := kais.controlSvc.GetByID(id); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", id), "not_exist", "control_ids")
		}
		seen[id] = true
	}
	return nil
}

// processes
func (kais *KPIAssociatedItemService) GetAllProcesses(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemProcessListResponseDTO, error) {
	if _, err := kais.kpiRegisterSvc.GetByID(kpiID); err != nil {
		return nil, err
	}
	return kais.repo.GetAllProcesses(kpiID, params)
}

func (kais *KPIAssociatedItemService) CreateKPIProcesses(kpiID, userID uint, dto *strategyDTOs.KPIProcessMapDTO) error {
	if err := kais.validateProcessIDs(kpiID, dto.ProcessIDs); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()
	var mappings []*models.KPIProcess
	currentTime := time.Now()

	for _, processID := range dto.ProcessIDs {
		mappings = append(mappings, &models.KPIProcess{
			KPIRegisterID: kpiID,
			ProcessID:     processID,
			CreatedByID:   userID,
			CreatedAt:     currentTime,
		})
	}

	if err := kais.repo.CreateKPIProcesses(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (kais *KPIAssociatedItemService) DeleteKPIProcess(kpiID, processID uint) error {
	// Validate process existence
	if _, err := kais.processSvc.GetByID(processID); err != nil {
		return utility.NewValidationError(fmt.Errorf("process does not exist for ID: %d", processID), "not_exist", "process_id")
	}

	tx := kais.repo.BeginTransaction()

	if err := kais.repo.DeleteKPIProcessMapping(kpiID, processID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (kais *KPIAssociatedItemService) validateProcessIDs(kpiID uint, processIDs []uint) error {
	if len(processIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no processes provided"), "required", "process_ids")
	}

	seen := make(map[uint]bool)
	for _, id := range processIDs {
		alreadyMapped, err := kais.repo.IsProcessMappedToKPI(kpiID, id)
		if err != nil {
			return fmt.Errorf("failed to check existing mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(fmt.Errorf("process ID %d already mapped to this kpi", id), "duplicate", "process_ids")
		}
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate process ID: %d", id), "duplicate", "process_ids")
		}
		if _, err := kais.processSvc.GetByID(id); err != nil {
			return utility.NewValidationError(fmt.Errorf("process does not exist for ID: %d", id), "not_exist", "process_ids")
		}
		seen[id] = true
	}
	return nil
}

// compliance
func (kais *KPIAssociatedItemService) GetAllCompliances(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemComplianceListResponseDTO, error) {
	if _, err := kais.kpiRegisterSvc.GetByID(kpiID); err != nil {
		return nil, err
	}
	return kais.repo.GetAllCompliances(kpiID, params)
}

func (kais *KPIAssociatedItemService) CreateKPICompliances(kpiID, userID uint, dto *strategyDTOs.KPIComplianceMapDTO) error {
	if err := kais.validateComplianceIDs(kpiID, dto.ComplianceRegisterIDs); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()
	var mappings []*models.KPICompliance
	currentTime := time.Now()

	for _, complianceID := range dto.ComplianceRegisterIDs {
		mappings = append(mappings, &models.KPICompliance{
			KPIRegisterID:        kpiID,
			ComplianceRegisterID: complianceID,
			CreatedByID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := kais.repo.CreateKPICompliances(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (kais *KPIAssociatedItemService) DeleteKPICompliance(kpiID, complianceID uint) error {
	// Validate compliance existence
	if _, err := kais.complianceSvc.GetByID(complianceID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance does not exist for ID: %d", complianceID), "not_exist", "compliance_id")
	}

	tx := kais.repo.BeginTransaction()

	if err := kais.repo.DeleteKPIComplianceMapping(kpiID, complianceID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (kais *KPIAssociatedItemService) validateComplianceIDs(kpiID uint, complianceIDs []uint) error {
	if len(complianceIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no compliances provided"), "required", "compliance_ids")
	}

	seen := make(map[uint]bool)
	for _, id := range complianceIDs {
		alreadyMapped, err := kais.repo.IsComplianceMappedToKPI(kpiID, id)
		if err != nil {
			return fmt.Errorf("failed to check existing mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(fmt.Errorf("compliance ID %d already mapped to this kpi", id), "duplicate", "compliance_ids")
		}
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate compliance ID: %d", id), "duplicate", "compliance_ids")
		}
		if _, err := kais.complianceSvc.GetByID(id); err != nil {
			return utility.NewValidationError(fmt.Errorf("compliance does not exist for ID: %d", id), "not_exist", "compliance_ids")
		}
		seen[id] = true
	}
	return nil
}

// Assets
func (kais *KPIAssociatedItemService) GetAllAssets(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemAssetListResponseDTO, error) {
	if _, err := kais.kpiRegisterSvc.GetByID(kpiID); err != nil {
		return nil, err
	}
	return kais.repo.GetAllAssets(kpiID, params)
}

func (kais *KPIAssociatedItemService) CreateKPIAssets(kpiID, userID uint, dto *strategyDTOs.KPIAssetMapDTO) error {
	if err := kais.validateAssetIDs(kpiID, dto.AssetIDs); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()
	var mappings []*models.KPIAsset
	currentTime := time.Now()

	for _, assetID := range dto.AssetIDs {
		mappings = append(mappings, &models.KPIAsset{
			KPIRegisterID:   kpiID,
			AssetRegisterID: assetID,
			CreatedByID:     userID,
			CreatedAt:       currentTime,
		})
	}

	if err := kais.repo.CreateKPIAssets(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (kais *KPIAssociatedItemService) DeleteKPIAsset(kpiID, assetID uint) error {
	// Validate asset existence
	if _, err := kais.assetSvc.GetByID(assetID); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset does not exist for ID: %d", assetID), "not_exist", "asset_id")
	}

	tx := kais.repo.BeginTransaction()

	if err := kais.repo.DeleteKPIAssetMapping(kpiID, assetID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (kais *KPIAssociatedItemService) validateAssetIDs(kpiID uint, assetIDs []uint) error {
	if len(assetIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no assets provided"), "required", "asset_ids")
	}

	seen := make(map[uint]bool)
	for _, id := range assetIDs {
		alreadyMapped, err := kais.repo.IsAssetMappedToKPI(kpiID, id)
		if err != nil {
			return fmt.Errorf("failed to check existing mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(fmt.Errorf("asset ID %d already mapped to this kpi", id), "duplicate", "asset_ids")
		}
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate asset ID: %d", id), "duplicate", "asset_ids")
		}
		if _, err := kais.assetSvc.GetByID(id); err != nil {
			return utility.NewValidationError(fmt.Errorf("asset does not exist for ID: %d", id), "not_exist", "asset_ids")
		}
		seen[id] = true
	}
	return nil
}

// Risks
func (kais *KPIAssociatedItemService) GetAllRisks(kpiID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	if _, err := kais.kpiRegisterSvc.GetByID(kpiID); err != nil {
		return nil, err
	}
	return kais.repo.GetAllRisks(kpiID, params)
}

func (kais *KPIAssociatedItemService) CreateKPIRisks(kpiID, userID uint, dto *strategyDTOs.KPIRiskMapDTO) error {
	if err := kais.validateRiskIDs(kpiID, dto.RiskIDs); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()
	var mappings []*models.KPIRisk
	currentTime := time.Now()

	for _, riskID := range dto.RiskIDs {
		mappings = append(mappings, &models.KPIRisk{
			KPIRegisterID:  kpiID,
			RiskRegisterID: riskID,
			CreatedByID:    userID,
			CreatedAt:      currentTime,
		})
	}

	if err := kais.repo.CreateKPIRisks(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (kais *KPIAssociatedItemService) DeleteKPIRisk(kpiID, riskID uint) error {
	// Validate risk existence
	if _, err := kais.riskSvc.GetByID(riskID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk does not exist for ID: %d", riskID), "not_exist", "risk_id")
	}

	tx := kais.repo.BeginTransaction()

	if err := kais.repo.DeleteKPIRiskMapping(kpiID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (kais *KPIAssociatedItemService) validateRiskIDs(kpiID uint, riskIDs []uint) error {
	if len(riskIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no risks provided"), "required", "risk_ids")
	}

	seen := make(map[uint]bool)
	for _, id := range riskIDs {
		alreadyMapped, err := kais.repo.IsRiskMappedToKPI(kpiID, id)
		if err != nil {
			return fmt.Errorf("failed to check existing mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(fmt.Errorf("risk ID %d already mapped to this kpi", id), "duplicate", "risk_ids")
		}
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate risk ID: %d", id), "duplicate", "risk_ids")
		}
		if _, err := kais.riskSvc.GetByID(id); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk does not exist for ID: %d", id), "not_exist", "risk_ids")
		}
		seen[id] = true
	}
	return nil
}

// Third party register
func (kais *KPIAssociatedItemService) GetAllThirdParties(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemThirdPartyListResponseDTO, error) {
	if _, err := kais.kpiRegisterSvc.GetByID(kpiID); err != nil {
		return nil, err
	}
	return kais.repo.GetAllThirdParties(kpiID, params)
}

func (kais *KPIAssociatedItemService) CreateKPIThirdParties(kpiID, userID uint, dto *strategyDTOs.KPIThirdPartyMapDTO) error {
	if err := kais.validateThirdPartyRegisterIDs(kpiID, dto.ThirdPartyRegisterIDs); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()
	var mappings []*models.KPIThirdParty
	currentTime := time.Now()

	for _, thirdPartyRegisterID := range dto.ThirdPartyRegisterIDs {
		mappings = append(mappings, &models.KPIThirdParty{
			KPIRegisterID:        kpiID,
			ThirdPartyRegisterID: thirdPartyRegisterID,
			CreatedByID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := kais.repo.CreateKPIThirdParties(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (kais *KPIAssociatedItemService) DeleteKPIThirdParty(kpiID, thirdPartyRegisterID uint) error {
	// Validate third party existence
	if _, err := kais.thirdPartySvc.GetByID(thirdPartyRegisterID); err != nil {
		return utility.NewValidationError(fmt.Errorf("third party register does not exist for ID: %d", thirdPartyRegisterID), "not_exist", "third_party_register_id")
	}

	tx := kais.repo.BeginTransaction()

	if err := kais.repo.DeleteKPIThirdPartyMapping(kpiID, thirdPartyRegisterID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (kais *KPIAssociatedItemService) validateThirdPartyRegisterIDs(kpiID uint, thirdPartyRegisterIDs []uint) error {
	if len(thirdPartyRegisterIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no third party register provided"), "required", "third_party_register_ids")
	}

	seen := make(map[uint]bool)
	for _, id := range thirdPartyRegisterIDs {
		alreadyMapped, err := kais.repo.IsThirdPartyMappedToKPI(kpiID, id)
		if err != nil {
			return fmt.Errorf("failed to check existing mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(fmt.Errorf("third party register ID %d already mapped to this kpi", id), "duplicate", "third_party_register_ids")
		}
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate third party register ID: %d", id), "duplicate", "third_party_register_ids")
		}
		if _, err := kais.thirdPartySvc.GetByID(id); err != nil {
			return utility.NewValidationError(fmt.Errorf("third party register does not exist for ID: %d", id), "not_exist", "third_party_register_ids")
		}
		seen[id] = true
	}
	return nil
}

// Policy
func (kais *KPIAssociatedItemService) GetAllPolicies(kpiID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := kais.kpiRegisterSvc.GetByID(kpiID); err != nil {
		return nil, err
	}
	return kais.repo.GetAllPolicies(kpiID, params)
}

func (kais *KPIAssociatedItemService) CreateKPIPolicies(kpiID, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {
	if _, err := kais.kpiRegisterSvc.GetByID(kpiID); err != nil {
		return utility.NewValidationError(fmt.Errorf("process does not exist for ID: %d", kpiID), "not_exist", "process_id")
	}
	if err := kais.validationSvc.ValidateDocumentMapping(kpiID, dto.DoumentIDs, "kpi_associated_documents", "kpi_register_id"); err != nil {
		return err
	}
	tx := kais.repo.BeginTransaction()
	var mappings []*models.KPIAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DoumentIDs {
		mappings = append(mappings, &models.KPIAssociatedDocument{
			KPIRegisterID: kpiID,
			DocumentID:    documentID,
			CreatedByID:   userID,
			CreatedAt:     currentTime,
		})
	}

	if err := kais.repo.CreateKPIPolicies(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (kais *KPIAssociatedItemService) DeleteKPIPolicy(kpiID, documentID uint) error {
	// Validate document existence
	if _, err := kais.documentSvc.GetByID(documentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document does not exist for ID: %d", documentID), "not_exist", "document_id")
	}

	tx := kais.repo.BeginTransaction()

	if err := kais.repo.DeleteKPIPolicyMapping(kpiID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// func (kais *KPIAssociatedItemService) validateDocumentIDs(kpiID uint, documentIDs []uint) error {
// 	if len(documentIDs) == 0 {
// 		return utility.NewValidationError(fmt.Errorf("no document provided"), "required", "document_ids")
// 	}

//		seen := make(map[uint]bool)
//		for _, id := range documentIDs {
//			alreadyMapped, err := kais.repo.IsPolicyMappedToKPI(kpiID, id)
//			if err != nil {
//				return fmt.Errorf("failed to check existing mapping: %v", err)
//			}
//			if alreadyMapped {
//				return utility.NewValidationError(fmt.Errorf("document ID %d already mapped to this kpi", id), "duplicate", "document_ids")
//			}
//			if seen[id] {
//				return utility.NewValidationError(fmt.Errorf("duplicate document ID: %d", id), "duplicate", "document_ids")
//			}
//			if _, err := kais.documentSvc.GetByID(id); err != nil {
//				return utility.NewValidationError(fmt.Errorf("document does not exist for ID: %d", id), "not_exist", "document_ids")
//			}
//			seen[id] = true
//		}
//		return nil
//	}
func (kais *KPIAssociatedItemService) GetAllClauses(kpiID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {
	if _, err := kais.kpiRegisterSvc.GetByID(kpiID); err != nil {
		return nil, err
	}

	response, err := kais.repo.GetAllClauses(kpiID, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (kais *KPIAssociatedItemService) UpdateClauses(kpiID uint, currentUserID uint, RiskClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error {

	if _, err := kais.kpiRegisterSvc.GetByID(kpiID); err != nil {
		return err
	}
	tx := kais.repo.BeginTransaction()
	clauses, err := kais.validateClauseIDs(RiskClauseFormDTO.ClauseIDs, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingControlClauses, err := kais.repo.GetClausesByKPIID(kpiID, tx)
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

	var clausesToAdd []*models.KPICompliance
	now := time.Now()
	for id, clause := range newClauses {
		if !existingClauses[id] {
			clausesToAdd = append(clausesToAdd, &models.KPICompliance{
				KPIRegisterID:        kpiID,
				ComplianceRegisterID: clause.ComplianceRegisterID,
				ComplianceClauseID:   &clause.ID,
				CreatedByID:          currentUserID,
				CreatedAt:            now,
			})
		}
	}

	if len(clausesToAdd) > 0 {
		if err := kais.repo.CreateKPIClauses(clausesToAdd, tx); err != nil {
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
		if err := kais.repo.DeleteKPIClauses(kpiID, deleteClauseIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}
func (kais *KPIAssociatedItemService) validateClauseIDs(clauseIDs []uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceClause, error) {
	seen := make(map[uint]bool)
	for _, id := range clauseIDs {
		if seen[id] {
			return nil, utility.NewValidationError(fmt.Errorf("duplicate clause ID: %d", id), "duplicate", "clause_id")
		}
		seen[id] = true
	}

	clauses, err := kais.clauseSvc.GetClausesByIDs(clauseIDs, tx...)
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
