package services

import (
	"errors"
	"fmt"
	activityLogServices "jethur/grc/logs/services"
	"jethur/grc/models"
	orgDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProcessInternalDependencyServiceInterface interface {
	GetAll(processID uint, params *orgDTOs.ProcessInternalDependencyListQueryParamsDTO) (*orgDTOs.ProcessInternalDependencyListResponseDTO, error)
	GetFullDetailsByID(processID uint, id uint) (*orgDTOs.ProcessInternalDependencyDetailDTO, error)
	GetByID(processID uint, id uint) (*models.ProcessInternalDependency, error)
	Create(processID uint, currentUserID uint, dto *orgDTOs.ProcessInternalDependencyFormDTO) (uint, error)
	Update(processID uint, id uint, currentUserID uint, dto *orgDTOs.ProcessInternalDependencyFormDTO) error
	Delete(processID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type ProcessInternalDependencyService struct {
	repo           repositories.ProcessInternalDependencyRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
	processSvc     ProcessServiceInterface
}

func NewProcessInternalDependencyService(
	repo repositories.ProcessInternalDependencyRepositoryInterface,
	activityLogSvc activityLogServices.ActivityLogServiceInterface,
	processSvc ProcessServiceInterface,
) ProcessInternalDependencyServiceInterface {
	return &ProcessInternalDependencyService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		processSvc:     processSvc,
	}
}

func (pids *ProcessInternalDependencyService) GetAll(processID uint, params *orgDTOs.ProcessInternalDependencyListQueryParamsDTO) (*orgDTOs.ProcessInternalDependencyListResponseDTO, error) {
	if _, err := pids.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pids.repo.GetAll(processID, params)
}

func (pids *ProcessInternalDependencyService) GetByID(processID uint, id uint) (*models.ProcessInternalDependency, error) {
	if _, err := pids.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	item, err := pids.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return item, nil
}

func (pids *ProcessInternalDependencyService) GetFullDetailsByID(processID uint, id uint) (*orgDTOs.ProcessInternalDependencyDetailDTO, error) {
	var detailDTO orgDTOs.ProcessInternalDependencyDetailDTO

	item, err := pids.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(item, &detailDTO)
	if err != nil {
		return nil, err
	}

	detailDTO.CreatedAt = item.CreatedAt
	detailDTO.UpdatedAt = item.UpdatedAt

	return &detailDTO, nil
}

func (pids *ProcessInternalDependencyService) validateFormDTO(dto *orgDTOs.ProcessInternalDependencyFormDTO) error {
	if dto.DependencyPercentageFrom != nil {
		if *dto.DependencyPercentageFrom < 0 || *dto.DependencyPercentageFrom > 100 {
			return utility.NewValidationError(fmt.Errorf("dependency_percentage_from must be between 0 and 100"), "out_of_range", "dependency_percentage_from")
		}
	}
	if dto.DependencyPercentageTo != nil {
		if *dto.DependencyPercentageTo < 0 || *dto.DependencyPercentageTo > 100 {
			return utility.NewValidationError(fmt.Errorf("dependency_percentage_to must be between 0 and 100"), "out_of_range", "dependency_percentage_to")
		}
	}
	return nil
}

func (pids *ProcessInternalDependencyService) Create(processID uint, currentUserID uint, dto *orgDTOs.ProcessInternalDependencyFormDTO) (uint, error) {
	if _, err := pids.processSvc.GetByID(processID); err != nil {
		return 0, err
	}

	if err := pids.validateFormDTO(dto); err != nil {
		return 0, err
	}

	item := &models.ProcessInternalDependency{
		ProcessID:                processID,
		SupportingUnitID:         dto.SupportingUnitID,
		DependencyTierID:         dto.DependencyTierID,
		DependencyPercentageFrom: dto.DependencyPercentageFrom,
		DependencyPercentageTo:   dto.DependencyPercentageTo,
		RequiredSupport:          dto.RequiredSupport,
		Notes:                    dto.Notes,
		CreatedBYID:              &currentUserID,
	}

	id, err := pids.repo.Create(item)
	if err != nil {
		return 0, err
	}

	go func() {
		created, err := pids.repo.GetByID(id)
		if err != nil {
			return
		}
		_ = pids.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ProcessInternalDependency{}.TableName(),
			id,
			created,
			currentUserID,
			"organization", "process", "process_internal_dependency", nil,
		)
	}()

	return id, nil
}

func (pids *ProcessInternalDependencyService) Update(processID uint, id uint, currentUserID uint, dto *orgDTOs.ProcessInternalDependencyFormDTO) error {
	if _, err := pids.processSvc.GetByID(processID); err != nil {
		return err
	}

	existing, err := pids.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	oldData := *existing

	if err := pids.validateFormDTO(dto); err != nil {
		return err
	}

	existing.SupportingUnitID = dto.SupportingUnitID
	existing.DependencyTierID = dto.DependencyTierID
	existing.DependencyPercentageFrom = dto.DependencyPercentageFrom
	existing.DependencyPercentageTo = dto.DependencyPercentageTo
	existing.RequiredSupport = dto.RequiredSupport
	existing.Notes = dto.Notes
	existing.UpdatedBYID = &currentUserID

	if err := pids.repo.Update(existing); err != nil {
		return err
	}

	go func() {
		oldFields := map[string]interface{}{
			"supporting_unit_id":         oldData.SupportingUnitID,
			"dependency_tier_id":         oldData.DependencyTierID,
			"dependency_percentage_from": oldData.DependencyPercentageFrom,
			"dependency_percentage_to":   oldData.DependencyPercentageTo,
			"required_support":           oldData.RequiredSupport,
			"notes":                      oldData.Notes,
		}
		newFields := map[string]interface{}{
			"supporting_unit_id":         existing.SupportingUnitID,
			"dependency_tier_id":         existing.DependencyTierID,
			"dependency_percentage_from": existing.DependencyPercentageFrom,
			"dependency_percentage_to":   existing.DependencyPercentageTo,
			"required_support":           existing.RequiredSupport,
			"notes":                      existing.Notes,
		}

		changedData := utility.CompareFieldChanges(oldFields, newFields)
		if len(changedData) > 0 {
			_ = pids.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ProcessInternalDependency{}.TableName(),
				id,
				changedData,
				currentUserID,
				"organization", "process", "process_internal_dependency", nil,
			)
		}
	}()

	return nil
}

func (pids *ProcessInternalDependencyService) Delete(processID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := pids.processSvc.GetByID(processID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		if _, err := pids.repo.GetByID(id); err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		if err := pids.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go pids.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.ProcessInternalDependency{}.TableName(),
			successfulIds,
			nil,
			currentUserID,
			nil,
		)
	}

	return failedIds, nil
}
