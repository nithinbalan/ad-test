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
