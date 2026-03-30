package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/models"
	"jethur/grc/strategy/repositories"
	"jethur/utility"
	"time"
)

type StrategicInitiativeAssociatedItemServiceInterface interface {
	GetAllStrategicIntiativeRisks(strategicIntiativeID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	CreateStrategicIntiativeRisk(strategicIntiativeID uint, currentUserID uint, dto *commonDTOs.AssociatedRiskMapDTO) error
	DeleteStrategicIntiativeRisk(strategicIntiativeID uint, riskID uint) error
}

type StrategicInitiativeAssociatedItemService struct {
	repo                  repositories.StrategicInitiativeAssociatedItemRepositoryInterface
	strategicIntiativeSvc StrategicIntiativeServiceInterface
	validationSvc         commonServices.AssociatedItemValidationServiceInterface
}

func NewStrategicInitiativeAssociatedItemService(
	repo repositories.StrategicInitiativeAssociatedItemRepositoryInterface,
	strategicIntiativeSvc StrategicIntiativeServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) StrategicInitiativeAssociatedItemServiceInterface {
	return &StrategicInitiativeAssociatedItemService{
		repo:                  repo,
		strategicIntiativeSvc: strategicIntiativeSvc,
		validationSvc:         validationSvc,
	}
}

func (siais *StrategicInitiativeAssociatedItemService) GetAllStrategicIntiativeRisks(strategicIntiativeID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	if _, err := siais.strategicIntiativeSvc.GetByID(strategicIntiativeID); err != nil {
		return nil, err
	}
	return siais.repo.GetAllStrategicIntiativeRisks(strategicIntiativeID, params)
}

func (siais *StrategicInitiativeAssociatedItemService) CreateStrategicIntiativeRisk(strategicIntiativeID, currentUserID uint, dto *commonDTOs.AssociatedRiskMapDTO) error {
	if _, err := siais.strategicIntiativeSvc.GetByID(strategicIntiativeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("strategic intiative does not exist for ID: %d", strategicIntiativeID), "not_exist", "strategic_intiative_id")
	}
	if err := siais.validationSvc.ValidateRiskMapping(strategicIntiativeID, dto.RiskIDs, "risks_strategic_intiatives", "strategic_intiative_id"); err != nil {
		return err
	}
	tx := siais.repo.BeginTransaction()
	var mappings []*models.RiskStrategicInitiative
	currentTime := time.Now()

	for _, riskID := range dto.RiskIDs {
		mappings = append(mappings, &models.RiskStrategicInitiative{
			RiskRegisterID:       riskID,
			StrategicIntiativeID: strategicIntiativeID,
			CreatedByID:          currentUserID,
			CreatedAt:            currentTime,
		})
	}

	if err := siais.repo.CreateStrategicIntiativeRisk(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (siais *StrategicInitiativeAssociatedItemService) DeleteStrategicIntiativeRisk(strategicIntiativeID uint, riskID uint) error {
	if _, err := siais.strategicIntiativeSvc.GetByID(strategicIntiativeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("strategic intiative does not exist for ID: %d", strategicIntiativeID), "not_exist", "strategic_intiative_id")
	}

	tx := siais.repo.BeginTransaction()

	if err := siais.repo.DeleteStrategicIntiativeRisk(strategicIntiativeID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
