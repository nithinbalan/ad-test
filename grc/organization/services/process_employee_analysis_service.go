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

type ProcessEmployeeAnalysisServiceInterface interface {
	GetAll(processID uint, params *orgDTOs.ProcessEmployeeAnalysisListQueryParamsDTO) (*orgDTOs.ProcessEmployeeAnalysisListResponseDTO, error)
	GetFullDetailsByID(processID uint, id uint) (*orgDTOs.ProcessEmployeeAnalysisDetailDTO, error)
	GetByID(processID uint, id uint) (*models.ProcessEmployeeAnalysis, error)
	Create(processID uint, currentUserID uint, dto *orgDTOs.ProcessEmployeeAnalysisFormDTO) (uint, error)
	Update(processID uint, id uint, currentUserID uint, dto *orgDTOs.ProcessEmployeeAnalysisFormDTO) error
	Delete(processID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type ProcessEmployeeAnalysisService struct {
	repo           repositories.ProcessEmployeeAnalysisRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
	processSvc     ProcessServiceInterface
}

func NewProcessEmployeeAnalysisService(
	repo repositories.ProcessEmployeeAnalysisRepositoryInterface,
	activityLogSvc activityLogServices.ActivityLogServiceInterface,
	processSvc ProcessServiceInterface,
) ProcessEmployeeAnalysisServiceInterface {
	return &ProcessEmployeeAnalysisService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		processSvc:     processSvc,
	}
}

func (peas *ProcessEmployeeAnalysisService) GetAll(processID uint, params *orgDTOs.ProcessEmployeeAnalysisListQueryParamsDTO) (*orgDTOs.ProcessEmployeeAnalysisListResponseDTO, error) {
	if _, err := peas.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return peas.repo.GetAll(processID, params)
}

func (peas *ProcessEmployeeAnalysisService) GetByID(processID uint, id uint) (*models.ProcessEmployeeAnalysis, error) {
	if _, err := peas.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	item, err := peas.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return item, nil
}

func (peas *ProcessEmployeeAnalysisService) GetFullDetailsByID(processID uint, id uint) (*orgDTOs.ProcessEmployeeAnalysisDetailDTO, error) {
	var detailDTO orgDTOs.ProcessEmployeeAnalysisDetailDTO

	item, err := peas.repo.GetFullDetailsByID(id)
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

func (peas *ProcessEmployeeAnalysisService) validateFormDTO(dto *orgDTOs.ProcessEmployeeAnalysisFormDTO) error {
	if dto.NormalCount != nil && *dto.NormalCount < 0 {
		return utility.NewValidationError(fmt.Errorf("normal_count must be greater than or equal to 0"), "out_of_range", "normal_count")
	}
	if dto.CriticalCount != nil && *dto.CriticalCount < 0 {
		return utility.NewValidationError(fmt.Errorf("critical_count must be greater than or equal to 0"), "out_of_range", "critical_count")
	}
	return nil
}

func (peas *ProcessEmployeeAnalysisService) Create(processID uint, currentUserID uint, dto *orgDTOs.ProcessEmployeeAnalysisFormDTO) (uint, error) {
	if _, err := peas.processSvc.GetByID(processID); err != nil {
		return 0, err
	}

	if err := peas.validateFormDTO(dto); err != nil {
		return 0, err
	}

	tx := peas.repo.BeginTransaction()

	item := &models.ProcessEmployeeAnalysis{
		ProcessID:     processID,
		DesignationID: dto.DesignationID,
		NormalCount:   dto.NormalCount,
		CriticalCount: dto.CriticalCount,
		CreatedBYID:   &currentUserID,
	}

	id, err := peas.repo.Create(item, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if dto.UserIDs != nil {
		for _, userID := range *dto.UserIDs {
			emp := &models.ProcessEmployeeAnalysisEmployee{
				ProcessEmployeeAnalysisID: id,
				UserID:                    userID,
				CreatedBYID:               &currentUserID,
			}
			if err := peas.repo.CreateEmployee(emp, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		created, err := peas.repo.GetByID(id)
		if err != nil {
			return
		}
		_ = peas.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ProcessEmployeeAnalysis{}.TableName(),
			id,
			created,
			currentUserID,
			"organization", "process", "process_employee_analysis", nil,
		)
	}()

	return id, nil
}

func (peas *ProcessEmployeeAnalysisService) Update(processID uint, id uint, currentUserID uint, dto *orgDTOs.ProcessEmployeeAnalysisFormDTO) error {
	if _, err := peas.processSvc.GetByID(processID); err != nil {
		return err
	}

	existing, err := peas.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	if err := peas.validateFormDTO(dto); err != nil {
		return err
	}

	tx := peas.repo.BeginTransaction()

	oldData := *existing

	existing.DesignationID = dto.DesignationID
	existing.NormalCount = dto.NormalCount
	existing.CriticalCount = dto.CriticalCount
	existing.UpdatedBYID = &currentUserID

	if err := peas.repo.Update(existing, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := peas.repo.DeleteEmployeesByAnalysisID(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	if dto.UserIDs != nil {
		for _, userID := range *dto.UserIDs {
			emp := &models.ProcessEmployeeAnalysisEmployee{
				ProcessEmployeeAnalysisID: id,
				UserID:                    userID,
				CreatedBYID:               &currentUserID,
			}
			if err := peas.repo.CreateEmployee(emp, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	changed, _ := utility.CompareFields(&oldData, existing, []string{
		"DesignationID", "NormalCount", "CriticalCount",
	})
	changedDataJSON, _ := utility.ChangesToJSON(changed)

	go func() {
		_ = peas.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ProcessEmployeeAnalysis{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			"organization", "process", "process_employee_analysis", nil,
		)
	}()

	return nil
}

func (peas *ProcessEmployeeAnalysisService) Delete(processID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := peas.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	var failedIds []uint
	var successfulIds []uint
	var deletedItems []map[string]interface{}

	for _, id := range ids {
		item, err := peas.repo.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		tx := peas.repo.BeginTransaction()

		if err := peas.repo.DeleteEmployeesByAnalysisID(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		if err := peas.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		if err := tx.Commit(); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
			deletedItems = append(deletedItems, map[string]interface{}{
				"id":             id,
				"designation_id": item.DesignationID,
				"normal_count":   item.NormalCount,
				"critical_count": item.CriticalCount,
			})
		}
	}

	if len(successfulIds) > 0 {
		go peas.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.ProcessEmployeeAnalysis{}.TableName(),
			successfulIds,
			deletedItems,
			currentUserID,
			nil,
		)
	}

	return failedIds, nil
}
