package services

import (
	"fmt"
	"jethur/grc/business-continuity/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/models"
	"jethur/utility"
	"time"
)

type BcpAssociatedItemServiceInterface interface {
	GetAllBcpDocuments(BcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateBcpDocument(BcpID uint, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteBcpDocument(BcpID uint, documentID uint) error
	GetAllBcpControl(BcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	CreateBcpControl(BcpID uint, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error
	DeleteBcpControl(BcpID uint, controlID uint) error
	GetAllBcpKPI(BcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error)
	CreateBcpKPI(BcpID uint, userID uint, dto *commonDTOs.AssociatedKPIMapDTO) error
	DeleteBcpKPI(BcpID uint, kpiRegisterID uint) error
	GetAllBcpRisk(BcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	CreateBcpRisk(BcpID uint, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error
	DeleteBcpRisk(BcpID uint, riskID uint) error
	GetAllBcpAssets(BcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error)
	CreateBcpAsset(BcpID uint, userID uint, dto *commonDTOs.AssociatedAssetMapDTO) error
	DeleteBcpAsset(BcpID uint, assetID uint) error
	GetAllBcpIncidents(BcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error)
	CreateBcpIncident(BcpID uint, userID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error
	DeleteBcpIncident(BcpID uint, incidentID uint) error
	GetAllBcpStrategicObjectives(BcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error)
	CreateBcpStrategicObjectives(BcpID uint, userID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error
	DeleteBcpStrategicObjectives(BcpID uint, strategicObjectiveID uint) error
	GetAllClauses(BcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	UpdateClauses(BcpID uint, currentUserID uint, ClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error
}

// Service and implementations
type BcpAssociatedItemService struct {
	repo               repositories.BcpAssociatedItemRepositoryInterface
	BcpSvc             BusinessContinuityPlanServiceInterface
	validationSvc      commonServices.AssociatedItemValidationServiceInterface
	bcpRiskAnalysisSvc BcpRiskAnalysisServiceInterface
}

func NewBcpAssociatedItemService(repo repositories.BcpAssociatedItemRepositoryInterface,
	BcpSvc BusinessContinuityPlanServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
	bcpRiskAnalysisSvc BcpRiskAnalysisServiceInterface,
) BcpAssociatedItemServiceInterface {
	return &BcpAssociatedItemService{
		repo:               repo,
		BcpSvc:             BcpSvc,
		validationSvc:      validationSvc,
		bcpRiskAnalysisSvc: bcpRiskAnalysisSvc,
	}
}

func (bais *BcpAssociatedItemService) GetAllBcpDocuments(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := bais.BcpSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return bais.repo.GetAllBcpDocuments(assetID, params)
}

func (bais *BcpAssociatedItemService) CreateBcpDocument(BcpID, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {

	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return utility.NewValidationError(fmt.Errorf("bcp does not exist for ID: %d", BcpID), "not_exist", "business_continuity_plan_id")
	}

	if err := bais.validationSvc.ValidateDocumentMapping(BcpID, dto.DocumentIDs, "bcp_associated_documents", "business_continuity_plan_id"); err != nil {
		return err
	}

	tx := bais.repo.BeginTransaction()
	var mappings []*models.BCPAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DocumentIDs {
		mappings = append(mappings, &models.BCPAssociatedDocument{
			DocumentID:               documentID,
			BusinessContinuityPlanID: BcpID,
			CreatedByID:              userID,
			CreatedAt:                currentTime,
		})
	}

	if err := bais.repo.CreateBcpDocument(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (bais *BcpAssociatedItemService) DeleteBcpDocument(BcpID uint, documentID uint) error {
	// Validate control existence
	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return utility.NewValidationError(fmt.Errorf("bcp does not exist for ID: %d", BcpID), "not_exist", "business_continuity_plan_id")
	}

	tx := bais.repo.BeginTransaction()

	if err := bais.repo.DeleteBcpDocument(BcpID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (bais *BcpAssociatedItemService) GetAllBcpControl(bcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
	if _, err := bais.BcpSvc.GetByID(bcpID); err != nil {
		return nil, err
	}
	return bais.repo.GetAllBcpControl(bcpID, params)
}

func (bais *BcpAssociatedItemService) CreateBcpControl(BcpID, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error {

	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return utility.NewValidationError(fmt.Errorf("bcp does not exist for ID: %d", BcpID), "not_exist", "business_continuity_plan_id")
	}

	if err := bais.validationSvc.ValidateControlMapping(BcpID, dto.ControlIDs, "bcp_controls", "business_continuity_plan_id"); err != nil {
		return err
	}

	tx := bais.repo.BeginTransaction()
	var mappings []*models.BCPControl
	currentTime := time.Now()

	for _, ControlID := range dto.ControlIDs {
		mappings = append(mappings, &models.BCPControl{
			ControlID:                ControlID,
			BusinessContinuityPlanID: BcpID,
			CreatedBYID:              userID,
			CreatedAt:                currentTime,
		})
	}

	if err := bais.repo.CreateBcpControl(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (bais *BcpAssociatedItemService) DeleteBcpControl(BcpID uint, controlID uint) error {
	// Validate control existence
	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return utility.NewValidationError(fmt.Errorf("bcp does not exist for ID: %d", BcpID), "not_exist", "business_continuity_plan_id")
	}

	tx := bais.repo.BeginTransaction()

	if err := bais.repo.DeleteBcpControl(BcpID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (bais *BcpAssociatedItemService) GetAllBcpKPI(bcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error) {
	if _, err := bais.BcpSvc.GetByID(bcpID); err != nil {
		return nil, err
	}
	return bais.repo.GetAllBcpKPI(bcpID, params)
}

func (bais *BcpAssociatedItemService) CreateBcpKPI(BcpID, userID uint, dto *commonDTOs.AssociatedKPIMapDTO) error {

	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return err
	}

	if err := bais.validationSvc.ValidateKPIMapping(BcpID, dto.KpiIDs, "bcp_kpis", "business_continuity_plan_id"); err != nil {
		return err
	}

	tx := bais.repo.BeginTransaction()
	var mappings []*models.BCPKPI
	currentTime := time.Now()

	for _, KpiID := range dto.KpiIDs {
		mappings = append(mappings, &models.BCPKPI{
			KPIRegisterID:            KpiID,
			BusinessContinuityPlanID: BcpID,
			CreatedByID:              userID,
			CreatedAt:                currentTime,
		})
	}

	if err := bais.repo.CreateBcpKPI(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (bais *BcpAssociatedItemService) DeleteBcpKPI(BcpID uint, kpiRegisterID uint) error {
	// Validate control existence
	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return err
	}

	tx := bais.repo.BeginTransaction()

	if err := bais.repo.DeleteBcpKPI(BcpID, kpiRegisterID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
func (bais *BcpAssociatedItemService) GetAllBcpRisk(bcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	if _, err := bais.BcpSvc.GetByID(bcpID); err != nil {
		return nil, err
	}
	return bais.bcpRiskAnalysisSvc.GetAllRisksByBcpId(bcpID, params)
}

func (bais *BcpAssociatedItemService) CreateBcpRisk(BcpID, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error {

	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return err
	}

	if err := bais.validationSvc.ValidateRiskMapping(BcpID, dto.RiskIDs, "bcp_risks", "business_continuity_plan_id"); err != nil {
		return err
	}

	tx := bais.repo.BeginTransaction()
	var mappings []*models.BcpRisk
	currentTime := time.Now()

	for _, RiskID := range dto.RiskIDs {
		mappings = append(mappings, &models.BcpRisk{
			RiskRegisterID:           RiskID,
			BusinessContinuityPlanID: BcpID,
			CreatedBYID:              userID,
			CreatedAt:                currentTime,
		})
	}

	if err := bais.repo.CreateBcpRisk(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (bais *BcpAssociatedItemService) DeleteBcpRisk(BcpID uint, riskID uint) error {
	// Validate control existence
	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return err
	}

	tx := bais.repo.BeginTransaction()

	if err := bais.repo.DeleteBcpRisk(BcpID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (bais *BcpAssociatedItemService) GetAllBcpAssets(bcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error) {
	if _, err := bais.BcpSvc.GetByID(bcpID); err != nil {
		return nil, err
	}
	return bais.repo.GetAllBcpAssets(bcpID, params)
}

func (bais *BcpAssociatedItemService) CreateBcpAsset(BcpID, userID uint, dto *commonDTOs.AssociatedAssetMapDTO) error {

	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return err
	}

	if err := bais.validationSvc.ValidateAssetMapping(BcpID, dto.AssetIDs, "bcp_assets", "business_continuity_plan_id"); err != nil {
		return err
	}

	tx := bais.repo.BeginTransaction()
	var mappings []*models.BcpAsset
	currentTime := time.Now()
	for _, AssetID := range dto.AssetIDs {
		mappings = append(mappings, &models.BcpAsset{
			AssetRegisterID:          AssetID,
			BusinessContinuityPlanID: BcpID,
			CreatedBYID:              userID,
			CreatedAt:                currentTime,
		})
	}

	if err := bais.repo.CreateBcpAsset(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (bais *BcpAssociatedItemService) DeleteBcpAsset(BcpID uint, assetID uint) error {
	// Validate control existence
	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return err
	}

	tx := bais.repo.BeginTransaction()

	if err := bais.repo.DeleteBcpAsset(BcpID, assetID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (bais *BcpAssociatedItemService) GetAllBcpIncidents(bcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error) {
	if _, err := bais.BcpSvc.GetByID(bcpID); err != nil {
		return nil, err
	}
	return bais.repo.GetAllBcpIncidents(bcpID, params)
}

func (bais *BcpAssociatedItemService) CreateBcpIncident(BcpID, userID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error {

	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return err
	}

	if err := bais.validationSvc.ValidateIncidentMapping(BcpID, dto.IncidentIDs, "bcp_incidents", "business_continuity_plan_id"); err != nil {
		return err
	}

	tx := bais.repo.BeginTransaction()
	var mappings []*models.BcpIncident
	currentTime := time.Now()
	for _, incidentID := range dto.IncidentIDs {
		mappings = append(mappings, &models.BcpIncident{
			IncidentRegisterID:       incidentID,
			BusinessContinuityPlanID: BcpID,
			CreatedBYID:              userID,
			CreatedAt:                currentTime,
		})
	}

	if err := bais.repo.CreateBcpIncident(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (bais *BcpAssociatedItemService) DeleteBcpIncident(BcpID uint, incidentID uint) error {
	// Validate control existence
	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return err
	}

	tx := bais.repo.BeginTransaction()

	if err := bais.repo.DeleteBcpIncident(BcpID, incidentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (bais *BcpAssociatedItemService) GetAllBcpStrategicObjectives(bcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error) {
	if _, err := bais.BcpSvc.GetByID(bcpID); err != nil {
		return nil, err
	}
	return bais.repo.GetAllBcpStrategicObjectives(bcpID, params)
}

func (bais *BcpAssociatedItemService) CreateBcpStrategicObjectives(BcpID, userID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error {

	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return err
	}

	if err := bais.validationSvc.ValidateStrtegicObjectiveMapping(BcpID, dto.StrategicObjectiveIDs, "bcp_strategic_objectives", "business_continuity_plan_id"); err != nil {
		return err
	}

	tx := bais.repo.BeginTransaction()
	var mappings []*models.BcpStrategicObjective
	currentTime := time.Now()
	for _, strategicObjectiveID := range dto.StrategicObjectiveIDs {
		mappings = append(mappings, &models.BcpStrategicObjective{
			StrategicObjectiveID:     strategicObjectiveID,
			BusinessContinuityPlanID: BcpID,
			CreatedBYID:              userID,
			CreatedAt:                currentTime,
		})
	}

	if err := bais.repo.CreateBcpStrategicObjectives(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (bais *BcpAssociatedItemService) DeleteBcpStrategicObjectives(BcpID uint, strategicObjectiveID uint) error {
	// Validate control existence
	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return err
	}

	tx := bais.repo.BeginTransaction()

	if err := bais.repo.DeleteBcpStrategicObjectives(BcpID, strategicObjectiveID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
func (bais *BcpAssociatedItemService) GetAllClauses(BcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {
	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return nil, err
	}

	response, err := bais.repo.GetAllClauses(BcpID, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (bais *BcpAssociatedItemService) UpdateClauses(BcpID uint, currentUserID uint, ClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error {

	if _, err := bais.BcpSvc.GetByID(BcpID); err != nil {
		return err
	}
	tx := bais.repo.BeginTransaction()
	clauses, err := bais.validationSvc.ValidateClauseIDs(ClauseFormDTO.ClauseIDs)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingControlClauses, err := bais.repo.GetClausesByBcpId(BcpID, tx)
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

	var clausesToAdd []*models.BcpComplianceRequirement
	now := time.Now()
	for id, clause := range newClauses {
		if !existingClauses[id] {
			clausesToAdd = append(clausesToAdd, &models.BcpComplianceRequirement{
				BusinessContinuityPlanID: BcpID,
				ComplianceRegisterID:     clause.ComplianceRegisterID,
				ComplianceClauseID:       clause.ID,
				CreatedByID:              currentUserID,
				CreatedAt:                now,
			})
		}
	}

	if len(clausesToAdd) > 0 {
		if err := bais.repo.CreateBcpClauses(clausesToAdd, tx); err != nil {
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
		if err := bais.repo.DeleteBcpClauses(BcpID, deleteClauseIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}
