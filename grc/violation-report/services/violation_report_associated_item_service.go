package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/models"
	violationReportRepositories "jethur/grc/violation-report/repositories"
	"jethur/utility"
	"time"
)

type ViolationReportAssociatedItemServiceInterface interface {
	GetAllControls(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	GetAllComplianceRegisters(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedComplianceListResponseDTO, error)
	GetAllRisks(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetAllDocuments(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	GetAllAssets(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error)
	CreateViolationReportRisk(violationReportID uint, currentUserID uint, dto *commonDTOs.AssociatedRiskMapDTO) error
	DeleteViolationReportRisk(violationReportID uint, riskID uint) error
	CreateViolationReportDocument(violationReportID uint, currentUserID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteViolationReportDocument(violationReportID uint, documentID uint) error
	CreateViolationReportControl(violationReportID uint, currentUserID uint, dto *commonDTOs.AssociatedControlMapDTO) error
	DeleteViolationReportControl(violationReportID uint, controlID uint) error
	CreateViolationReportAsset(violationReportID uint, currentUserID uint, dto *commonDTOs.AssociatedAssetMapDTO) error
	DeleteViolationReportAsset(violationReportID uint, assetID uint) error
	CreateViolationReportComplianceRegister(violationReportID uint, currentUserID uint, dto *commonDTOs.AssociatedComplianceMapDTO) error
	DeleteViolationReportComplianceRegister(violationReportID uint, complianceRegisterID uint) error
}

// Service and implementations
type ViolationReportAssociatedItemService struct {
	repo               violationReportRepositories.ViolationReportAssociatedItemRepositoryInterface
	violationReportSvc ViolationReportServiceInterface
	validationSvc      commonServices.AssociatedItemValidationServiceInterface
}

func NewViolationReportAssociatedItemService(
	repo violationReportRepositories.ViolationReportAssociatedItemRepositoryInterface,
	violationReportSvc ViolationReportServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,

) ViolationReportAssociatedItemServiceInterface {
	return &ViolationReportAssociatedItemService{
		repo:               repo,
		violationReportSvc: violationReportSvc,
		validationSvc:      validationSvc,
	}
}

func (vias *ViolationReportAssociatedItemService) GetAllControls(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
	if _, err := vias.violationReportSvc.GetByID(violationReportID); err != nil {
		return nil, err
	}
	return vias.repo.GetAllControls(violationReportID, params)
}

func (vias *ViolationReportAssociatedItemService) GetAllComplianceRegisters(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedComplianceListResponseDTO, error) {
	if _, err := vias.violationReportSvc.GetByID(violationReportID); err != nil {
		return nil, err
	}
	return vias.repo.GetAllComplianceRegisters(violationReportID, params)
}

func (vias *ViolationReportAssociatedItemService) GetAllRisks(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	if _, err := vias.violationReportSvc.GetByID(violationReportID); err != nil {
		return nil, err
	}
	return vias.repo.GetAllRisks(violationReportID, params)
}

func (vias *ViolationReportAssociatedItemService) GetAllDocuments(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := vias.violationReportSvc.GetByID(violationReportID); err != nil {
		return nil, err
	}
	return vias.repo.GetAllDocuments(violationReportID, params)
}

func (vias *ViolationReportAssociatedItemService) GetAllAssets(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error) {
	if _, err := vias.violationReportSvc.GetByID(violationReportID); err != nil {
		return nil, err
	}
	return vias.repo.GetAllAssets(violationReportID, params)
}

func (vias *ViolationReportAssociatedItemService) CreateViolationReportRisk(violationReportID uint, currentUserID uint, dto *commonDTOs.AssociatedRiskMapDTO) error {
	if _, err := vias.violationReportSvc.GetByID(violationReportID); err != nil {
		return utility.NewValidationError(fmt.Errorf("violation report does not exist for ID: %d", violationReportID), "not_exist", "violation_report_id")
	}
	if err := vias.validationSvc.ValidateRiskMapping(violationReportID, dto.RiskIDs, "violation_report_risks", "violation_report_id"); err != nil {
		return err
	}
	tx := vias.repo.BeginTransaction()
	var mappings []*models.ViolationReportRisk
	currentTime := time.Now()

	for _, riskID := range dto.RiskIDs {
		mappings = append(mappings, &models.ViolationReportRisk{
			RiskRegisterID:    riskID,
			ViolationReportID: violationReportID,
			CreatedByID:       currentUserID,
			CreatedAt:         currentTime,
		})
	}

	if err := vias.repo.CreateViolationReportRisk(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (vias *ViolationReportAssociatedItemService) DeleteViolationReportRisk(violationReportID uint, riskID uint) error {
	// Validate control existence
	if _, err := vias.violationReportSvc.GetByID(violationReportID); err != nil {
		return utility.NewValidationError(fmt.Errorf("violation report does not exist for ID: %d", violationReportID), "not_exist", "violation_report_id")
	}

	tx := vias.repo.BeginTransaction()

	if err := vias.repo.DeleteViolationReportRisk(violationReportID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (vias *ViolationReportAssociatedItemService) CreateViolationReportDocument(violationReportID, currentUserID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {
	if _, err := vias.violationReportSvc.GetByID(violationReportID); err != nil {
		return utility.NewValidationError(fmt.Errorf("violation report does not exist for ID: %d", violationReportID), "not_exist", "violation_report_id")
	}
	if err := vias.validationSvc.ValidateDocumentMapping(violationReportID, dto.DocumentIDs, "violation_report_associated_documents", "violation_report_id"); err != nil {
		return err
	}
	tx := vias.repo.BeginTransaction()
	var mappings []*models.ViolationReportAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DocumentIDs {
		mappings = append(mappings, &models.ViolationReportAssociatedDocument{
			DocumentID:        documentID,
			ViolationReportID: violationReportID,
			CreatedByID:       currentUserID,
			CreatedAt:         currentTime,
		})
	}

	if err := vias.repo.CreateViolationReportDocument(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (vias *ViolationReportAssociatedItemService) DeleteViolationReportDocument(violationReportID uint, documentID uint) error {
	// Validate control existence
	if _, err := vias.violationReportSvc.GetByID(violationReportID); err != nil {
		return utility.NewValidationError(fmt.Errorf("violation report does not exist for ID: %d", violationReportID), "not_exist", "violation_report_id")
	}

	tx := vias.repo.BeginTransaction()

	if err := vias.repo.DeleteViolationReportDocument(violationReportID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
func (vias *ViolationReportAssociatedItemService) CreateViolationReportControl(violationReportID, currentUserID uint, dto *commonDTOs.AssociatedControlMapDTO) error {
	if _, err := vias.violationReportSvc.GetByID(violationReportID); err != nil {
		return utility.NewValidationError(fmt.Errorf("violation report does not exist for ID: %d", violationReportID), "not_exist", "violation_report_id")
	}
	if err := vias.validationSvc.ValidateControlMapping(violationReportID, dto.ControlIDs, "violation_report_controls", "violation_report_id"); err != nil {
		return err
	}
	tx := vias.repo.BeginTransaction()
	var mappings []*models.ViolationReportControl
	currentTime := time.Now()

	for _, controlID := range dto.ControlIDs {
		mappings = append(mappings, &models.ViolationReportControl{
			ControlID:         controlID,
			ViolationReportID: violationReportID,
			CreatedByID:       currentUserID,
			CreatedAt:         currentTime,
		})
	}

	if err := vias.repo.CreateViolationReportControl(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (vias *ViolationReportAssociatedItemService) DeleteViolationReportControl(violationReportID uint, controlID uint) error {
	// Validate control existence
	if _, err := vias.violationReportSvc.GetByID(violationReportID); err != nil {
		return utility.NewValidationError(fmt.Errorf("violation report does not exist for ID: %d", violationReportID), "not_exist", "violation_report_id")
	}

	tx := vias.repo.BeginTransaction()

	if err := vias.repo.DeleteViolationReportControl(violationReportID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (vias *ViolationReportAssociatedItemService) CreateViolationReportAsset(violationReportID, currentUserID uint, dto *commonDTOs.AssociatedAssetMapDTO) error {
	if _, err := vias.violationReportSvc.GetByID(violationReportID); err != nil {
		return utility.NewValidationError(fmt.Errorf("violation report does not exist for ID: %d", violationReportID), "not_exist", "violation_report_id")
	}
	if err := vias.validationSvc.ValidateAssetMapping(violationReportID, dto.AssetIDs, "violation_report_assets", "violation_report_id"); err != nil {
		return err
	}
	tx := vias.repo.BeginTransaction()
	var mappings []*models.ViolationReportAsset
	currentTime := time.Now()

	for _, assetID := range dto.AssetIDs {
		mappings = append(mappings, &models.ViolationReportAsset{
			AssetRegisterID:   assetID,
			ViolationReportID: violationReportID,
			CreatedBYID:       currentUserID,
			CreatedAt:         currentTime,
		})
	}

	if err := vias.repo.CreateViolationReportAsset(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (vias *ViolationReportAssociatedItemService) DeleteViolationReportAsset(violationReportID uint, assetID uint) error {
	// Validate control existence
	if _, err := vias.violationReportSvc.GetByID(violationReportID); err != nil {
		return err
	}

	tx := vias.repo.BeginTransaction()

	if err := vias.repo.DeleteViolationReportAsset(violationReportID, assetID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (vias *ViolationReportAssociatedItemService) CreateViolationReportComplianceRegister(violationReportID, currentUserID uint, dto *commonDTOs.AssociatedComplianceMapDTO) error {
	if _, err := vias.violationReportSvc.GetByID(violationReportID); err != nil {
		return utility.NewValidationError(fmt.Errorf("violation report does not exist for ID: %d", violationReportID), "not_exist", "violation_report_id")
	}
	if err := vias.validationSvc.ValidateComplianceMapping(violationReportID, dto.ComplianceRegisterIDs, "violation_report_compliance_registers", "violation_report_id"); err != nil {
		return err
	}
	tx := vias.repo.BeginTransaction()
	var mappings []*models.ViolationReportComplianceRegister
	currentTime := time.Now()

	for _, complianceRegisterID := range dto.ComplianceRegisterIDs {
		mappings = append(mappings, &models.ViolationReportComplianceRegister{
			ComplianceRegisterID: complianceRegisterID,
			ViolationReportID:    violationReportID,
			CreatedByID:          currentUserID,
			CreatedAt:            currentTime,
		})
	}

	if err := vias.repo.CreateViolationReportComplianceRegister(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (vias *ViolationReportAssociatedItemService) DeleteViolationReportComplianceRegister(violationReportID uint, complianceRegisterID uint) error {
	if _, err := vias.violationReportSvc.GetByID(violationReportID); err != nil {
		return utility.NewValidationError(fmt.Errorf("violation report does not exist for ID: %d", violationReportID), "not_exist", "violation_report_id")
	}

	tx := vias.repo.BeginTransaction()

	if err := vias.repo.DeleteViolationReportComplianceRegister(violationReportID, complianceRegisterID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
