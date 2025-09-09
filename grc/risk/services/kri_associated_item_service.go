package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/models"
	"jethur/grc/risk/repositories"
	"jethur/utility"
	"time"
)

type KRIAssociatedItemServiceInterface interface {
	GetAllKRIDocuments(kriID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateKRIDocument(kriID uint, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteKRIDocument(kriID uint, documentID uint) error
	GetAllKRIControls(kriID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	CreateKRIControl(kriID uint, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error
	DeleteKRIControl(kriID uint, controlID uint) error
	GetAllKRIAssets(kriID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error)
	CreateKRIAsset(kriID uint, userID uint, dto *commonDTOs.AssociatedAssetMapDTO) error
	DeleteKRIAsset(kriID uint, assetID uint) error
	GetAllKRIIncidents(kriID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error)
	CreateKRIIncident(kriID uint, userID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error
	DeleteKRIIncident(kriID uint, incidentID uint) error
	GetAllKRIRisks(kriID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	CreateKRIRisk(kriID uint, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error
	DeleteKRIRisk(kriID uint, riskID uint) error
	GetAllKRIKPI(kriID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error)
	CreateKRIKPI(kriID uint, userID uint, dto *commonDTOs.AssociatedKPIMapDTO) error
	DeleteKRIKPI(kriID uint, kpiID uint) error
	GetAllKRIProcesses(kriID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	CreateKRIProcess(kriID uint, userID uint, dto *commonDTOs.AssociatedProcessMapDTO) error
	DeleteKRIProcess(kriID uint, processID uint) error
	GetAllClauses(kriID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	UpdateClauses(kriID uint, currentUserID uint, ClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error
}

// Service and implementations
type KRIAssociatedItemService struct {
	repo          repositories.KRIAssociatedItemRepositoryInterface
	KRISvc        KRIRegisterServiceInterface
	validationSvc commonServices.AssociatedItemValidationServiceInterface
}

func NewKRIAssociatedItemService(repo repositories.KRIAssociatedItemRepositoryInterface,
	KRISvc KRIRegisterServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) KRIAssociatedItemServiceInterface {
	return &KRIAssociatedItemService{
		repo:          repo,
		KRISvc:        KRISvc,
		validationSvc: validationSvc,
	}
}

func (kais *KRIAssociatedItemService) GetAllKRIDocuments(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := kais.KRISvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return kais.repo.GetAllKRIDocuments(assetID, params)
}

func (kais *KRIAssociatedItemService) CreateKRIDocument(kriID, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {

	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return utility.NewValidationError(fmt.Errorf("KRI does not exist for ID: %d", kriID), "not_exist", "kri_register_id")
	}

	if err := kais.validationSvc.ValidateDocumentMapping(kriID, dto.DocumentIDs, "kri_associated_documents", "kri_register_id"); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()
	var mappings []*models.KRIAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DocumentIDs {
		mappings = append(mappings, &models.KRIAssociatedDocument{
			DocumentID:    documentID,
			KRIRegisterID: kriID,
			CreatedByID:   userID,
			CreatedAt:     currentTime,
		})
	}

	if err := kais.repo.CreateKRIDocument(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (kais *KRIAssociatedItemService) DeleteKRIDocument(kriID uint, documentID uint) error {
	// Validate control existence
	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return utility.NewValidationError(fmt.Errorf("KRI does not exist for ID: %d", kriID), "not_exist", "kri_register_id")
	}

	tx := kais.repo.BeginTransaction()

	if err := kais.repo.DeleteKRIDocument(kriID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (kais *KRIAssociatedItemService) GetAllKRIControls(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
	if _, err := kais.KRISvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return kais.repo.GetAllKRIControls(assetID, params)
}

func (kais *KRIAssociatedItemService) CreateKRIControl(kriID, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error {

	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return utility.NewValidationError(fmt.Errorf("KRI does not exist for ID: %d", kriID), "not_exist", "kri_register_id")
	}

	if err := kais.validationSvc.ValidateControlMapping(kriID, dto.ControlIDs, "kri_controls", "kri_register_id"); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()
	var mappings []*models.KRIControl
	currentTime := time.Now()

	for _, controlID := range dto.ControlIDs {
		mappings = append(mappings, &models.KRIControl{
			ControlID:     controlID,
			KRIRegisterID: kriID,
			CreatedBYID:   userID,
			CreatedAt:     currentTime,
		})
	}

	if err := kais.repo.CreateKRIControl(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (kais *KRIAssociatedItemService) DeleteKRIControl(kriID uint, controlID uint) error {
	// Validate control existence
	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return utility.NewValidationError(fmt.Errorf("KRI does not exist for ID: %d", kriID), "not_exist", "kri_register_id")
	}

	tx := kais.repo.BeginTransaction()

	if err := kais.repo.DeleteKRIControl(kriID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (kais *KRIAssociatedItemService) GetAllKRIAssets(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error) {
	if _, err := kais.KRISvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return kais.repo.GetAllKRIAssets(assetID, params)
}

func (kais *KRIAssociatedItemService) CreateKRIAsset(kriID, userID uint, dto *commonDTOs.AssociatedAssetMapDTO) error {

	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return err
	}

	if err := kais.validationSvc.ValidateAssetMapping(kriID, dto.AssetIDs, "asset_kris", "kri_register_id"); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()
	var mappings []*models.AssetKri
	currentTime := time.Now()

	for _, AssetID := range dto.AssetIDs {
		mappings = append(mappings, &models.AssetKri{
			AssetRegisterID: AssetID,
			KriRegisterID:   kriID,
			CreatedBYID:     userID,
			CreatedAt:       currentTime,
		})
	}

	if err := kais.repo.CreateKRIAsset(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (kais *KRIAssociatedItemService) DeleteKRIAsset(kriID uint, assetID uint) error {

	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()

	if err := kais.repo.DeleteKRIAsset(kriID, assetID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (kais *KRIAssociatedItemService) GetAllKRIIncidents(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error) {
	if _, err := kais.KRISvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return kais.repo.GetAllKRIIncidents(assetID, params)
}

func (kais *KRIAssociatedItemService) CreateKRIIncident(kriID, userID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error {

	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return err
	}

	if err := kais.validationSvc.ValidateIncidentMapping(kriID, dto.IncidentIDs, "kri_incidents", "kri_register_id"); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()
	var mappings []*models.KriIncident
	currentTime := time.Now()

	for _, incidentID := range dto.IncidentIDs {
		mappings = append(mappings, &models.KriIncident{
			IncidentRegisterID: incidentID,
			KriRegisterID:      kriID,
			CreatedBYID:        userID,
			CreatedAt:          currentTime,
		})
	}

	if err := kais.repo.CreateKRIIncident(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (kais *KRIAssociatedItemService) DeleteKRIIncident(kriID uint, incidentID uint) error {

	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()

	if err := kais.repo.DeleteKRIIncident(kriID, incidentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (kais *KRIAssociatedItemService) GetAllKRIRisks(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	if _, err := kais.KRISvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return kais.repo.GetAllKRIRisks(assetID, params)
}

func (kais *KRIAssociatedItemService) CreateKRIRisk(kriID, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error {

	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return err
	}

	if err := kais.validationSvc.ValidateRiskMapping(kriID, dto.RiskIDs, "risk_kris", "kri_register_id"); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()
	var mappings []*models.RiskKri
	currentTime := time.Now()

	for _, riskID := range dto.RiskIDs {
		mappings = append(mappings, &models.RiskKri{
			RiskRegisterID: riskID,
			KriRegisterID:  kriID,
			CreatedByID:    userID,
			CreatedAt:      currentTime,
		})
	}

	if err := kais.repo.CreateKRIRisk(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (kais *KRIAssociatedItemService) DeleteKRIRisk(kriID uint, riskID uint) error {

	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()

	if err := kais.repo.DeleteKRIRisk(kriID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (kais *KRIAssociatedItemService) GetAllKRIKPI(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error) {
	if _, err := kais.KRISvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return kais.repo.GetAllKRIKPI(assetID, params)
}

func (kais *KRIAssociatedItemService) CreateKRIKPI(kriID, userID uint, dto *commonDTOs.AssociatedKPIMapDTO) error {

	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return err
	}

	if err := kais.validationSvc.ValidateKPIMapping(kriID, dto.KpiIDs, "kri_kpis", "kri_register_id"); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()
	var mappings []*models.KriKpi
	currentTime := time.Now()

	for _, kpiID := range dto.KpiIDs {
		mappings = append(mappings, &models.KriKpi{
			KPIRegisterID: kpiID,
			KRIRegisterID: kriID,
			CreatedByID:   userID,
			CreatedAt:     currentTime,
		})
	}

	if err := kais.repo.CreateKRIKPI(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (kais *KRIAssociatedItemService) DeleteKRIKPI(kriID uint, kpiID uint) error {

	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()

	if err := kais.repo.DeleteKRIKPI(kriID, kpiID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (kais *KRIAssociatedItemService) GetAllKRIProcesses(kriID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return nil, err
	}
	return kais.repo.GetAllKRIProcesses(kriID, params)
}

func (kais *KRIAssociatedItemService) CreateKRIProcess(kriID, userID uint, dto *commonDTOs.AssociatedProcessMapDTO) error {

	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return err
	}

	if err := kais.validationSvc.ValidateProcessMapping(kriID, dto.ProcessIDs, "kri_processes", "kri_register_id"); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()
	var mappings []*models.KriProcess
	currentTime := time.Now()

	for _, processID := range dto.ProcessIDs {
		mappings = append(mappings, &models.KriProcess{
			ProcessID:     processID,
			KRIRegisterID: kriID,
			CreatedByID:   userID,
			CreatedAt:     currentTime,
		})
	}

	if err := kais.repo.CreateKRIProcess(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (kais *KRIAssociatedItemService) DeleteKRIProcess(kriID uint, processID uint) error {

	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return err
	}

	tx := kais.repo.BeginTransaction()

	if err := kais.repo.DeleteKRIProcess(kriID, processID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (kais *KRIAssociatedItemService) GetAllClauses(kriID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {
	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return nil, err
	}

	response, err := kais.repo.GetAllClauses(kriID, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (kais *KRIAssociatedItemService) UpdateClauses(kriID uint, currentUserID uint, ClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error {

	if _, err := kais.KRISvc.GetByID(kriID); err != nil {
		return err
	}
	tx := kais.repo.BeginTransaction()
	clauses, err := kais.validationSvc.ValidateClauseIDs(ClauseFormDTO.ClauseIDs)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingControlClauses, err := kais.repo.GetClausesByKriId(kriID, tx)
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

	var clausesToAdd []*models.KriComplianceRequirement
	now := time.Now()
	for id, clause := range newClauses {
		if !existingClauses[id] {
			clausesToAdd = append(clausesToAdd, &models.KriComplianceRequirement{
				KriRegisterID:        kriID,
				ComplianceRegisterID: clause.ComplianceRegisterID,
				ComplianceClauseID:   clause.ID,
				CreatedByID:          currentUserID,
				CreatedAt:            now,
			})
		}
	}

	if len(clausesToAdd) > 0 {
		if err := kais.repo.CreateKriClauses(clausesToAdd, tx); err != nil {
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
		if err := kais.repo.DeleteKriClauses(kriID, deleteClauseIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}
