package services

import (
	"errors"
	"fmt"
	"time"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyInitiativeStatusUpdateServiceInterface interface {
	GetAll(strategyInitiativeID uint, sort string, sortBy string, keyword string) ([]*strategyDTOs.StrategyInitiativeStatusUpdateListDTO, error)
	Create(strategyInitiativeID uint, currentUserID uint, dto *strategyDTOs.StrategyInitiativeStatusUpdateFormDTO) (uint, error)
	Update(strategyInitiativeID uint, id uint, currentUserID uint, dto *strategyDTOs.StrategyInitiativeStatusUpdateFormDTO) error
	GetLatestStrategyInitiativeStatusUpdate(strategyInitiativeID uint) (*strategyDTOs.StrategyInitiativeStatusUpdateListDTO, error)
}

type StrategyInitiativeStatusUpdateService struct {
	repo                          strategyRepositories.StrategyInitiativeStatusUpdateRepositoryInterface
	strategicIntiativeRepo        strategyRepositories.StrategicIntiativeRepositoryInterface
	strategicInitiativeStatusRepo strategyRepositories.StrategicInitiativeStatusRepositoryInterface
	activityLogSvc                logServices.ActivityLogServiceInterface
}

func NewStrategyInitiativeStatusUpdateService(
	repo strategyRepositories.StrategyInitiativeStatusUpdateRepositoryInterface,
	strategicIntiativeRepo strategyRepositories.StrategicIntiativeRepositoryInterface,
	strategicInitiativeStatusRepo strategyRepositories.StrategicInitiativeStatusRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) StrategyInitiativeStatusUpdateServiceInterface {
	return &StrategyInitiativeStatusUpdateService{
		repo:                          repo,
		strategicIntiativeRepo:        strategicIntiativeRepo,
		strategicInitiativeStatusRepo: strategicInitiativeStatusRepo,
		activityLogSvc:                activityLogSvc,
	}
}

func (s *StrategyInitiativeStatusUpdateService) GetAll(strategyInitiativeID uint, sort string, sortBy string, keyword string) ([]*strategyDTOs.StrategyInitiativeStatusUpdateListDTO, error) {
	if _, err := s.strategicIntiativeRepo.GetByID(strategyInitiativeID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, strategyInitiativeID)
		}
		return nil, err
	}

	return s.repo.GetAll(strategyInitiativeID, sort, sortBy, keyword)
}

func (s *StrategyInitiativeStatusUpdateService) Create(strategyInitiativeID uint, currentUserID uint, dto *strategyDTOs.StrategyInitiativeStatusUpdateFormDTO) (uint, error) {
	tx := s.repo.BeginTransaction()

	strategicInitiative, err := s.strategicIntiativeRepo.GetByID(strategyInitiativeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, utility.NewResourceNotFoundError(err, strategyInitiativeID)
		}
		return 0, err
	}

	if _, err := s.strategicInitiativeStatusRepo.GetByID(dto.StrategyInitiativeStatusID); err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return 0, utility.NewResourceNotFoundError(err, dto.StrategyInitiativeStatusID)
		// }
		return 0, err
	}

	if dto.Percentage < 0 || dto.Percentage > 100 {
		return 0, utility.NewValidationError(fmt.Errorf("percentage must be between 0 and 100"), "invalid_range", "percentage")
	}

	statusUpdate := &models.StrategicInitiativeStatusUpdate{
		StrategyInitiativeID:        strategyInitiativeID,
		StrategicInitiativeStatusID: dto.StrategyInitiativeStatusID,
		Comment:                     dto.Comment,
		Percentage:                  dto.Percentage,
		CreatedByID:                 currentUserID,
	}

	id, err := s.repo.Create(statusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Keep existing single PUT API intact; sync completion and status from status-update APIs too.
	cp := float64(dto.Percentage)
	now := time.Now()
	strategicInitiative.CompletionPercentage = &cp
	strategicInitiative.StrategicInitiativeStatusID = dto.StrategyInitiativeStatusID
	strategicInitiative.UpdatedBYID = &currentUserID
	strategicInitiative.UpdatedAt = &now

	if err := s.strategicIntiativeRepo.Update(strategicInitiative, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		err := s.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StrategicInitiativeStatusUpdate{}.TableName(),
			id,
			statusUpdate,
			currentUserID,
			"organization", "strategy_intiative", "status_update", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StrategicInitiativeStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (s *StrategyInitiativeStatusUpdateService) Update(strategyInitiativeID uint, id uint, currentUserID uint, dto *strategyDTOs.StrategyInitiativeStatusUpdateFormDTO) error {
	tx := s.repo.BeginTransaction()

	strategicInitiative, err := s.strategicIntiativeRepo.GetByID(strategyInitiativeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, strategyInitiativeID)
		}
		return err
	}

	statusUpdate, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}
	if statusUpdate.StrategyInitiativeID != strategyInitiativeID {
		return utility.NewResourceNotFoundError(fmt.Errorf("status update not found"), id)
	}

	if _, err := s.strategicInitiativeStatusRepo.GetByID(dto.StrategyInitiativeStatusID); err != nil {
		return err
	}

	statusUpdate.StrategicInitiativeStatusID = dto.StrategyInitiativeStatusID
	statusUpdate.Comment = dto.Comment
	statusUpdate.Percentage = dto.Percentage
	statusUpdate.UpdatedByID = &currentUserID

	if err := s.repo.Update(statusUpdate, tx); err != nil {
		tx.Rollback()
		return err
	}

	cp := float64(dto.Percentage)
	now := time.Now()
	strategicInitiative.CompletionPercentage = &cp
	strategicInitiative.StrategicInitiativeStatusID = dto.StrategyInitiativeStatusID
	strategicInitiative.UpdatedBYID = &currentUserID
	strategicInitiative.UpdatedAt = &now

	if err := s.strategicIntiativeRepo.Update(strategicInitiative, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		changedData := map[string]interface{}{
			"strategic_initiative_status_id": dto.StrategyInitiativeStatusID,
			"percentage":                     dto.Percentage,
			"comment":                        dto.Comment,
		}
		err = s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.StrategicInitiativeStatusUpdate{}.TableName(),
			id,
			changedData,
			currentUserID,
			"organization", "strategy_intiative", "status_update", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.StrategicInitiativeStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (s *StrategyInitiativeStatusUpdateService) GetLatestStrategyInitiativeStatusUpdate(strategyInitiativeID uint) (*strategyDTOs.StrategyInitiativeStatusUpdateListDTO, error) {
	if _, err := s.strategicIntiativeRepo.GetByID(strategyInitiativeID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, strategyInitiativeID)
		}
		return nil, err
	}

	return s.repo.GetLatestStrategyInitiativeStatusUpdate(strategyInitiativeID)
}
