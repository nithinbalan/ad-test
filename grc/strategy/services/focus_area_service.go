package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/strategy/dtos"
	"jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type FocusAreaServiceInterface interface {
	GetAll(params *dtos.FocusAreaListQueryParamsDTO) (*dtos.FocusAreaListResponseDTO, error)
	GetByID(id uint) (*models.FocusArea, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.FocusAreaDetailDTO, error)
	Create(currentUserID uint, focusAreaFormDTO *dtos.FocusAreaFormDTO) (uint, error)
	Update(id uint, currentUserID uint, dto *dtos.FocusAreaFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type FocusAreaService struct {
	repo           repositories.FocusAreaRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewFocusAreaService(repo repositories.FocusAreaRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
) FocusAreaServiceInterface {
	return &FocusAreaService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (fas *FocusAreaService) GetAll(params *dtos.FocusAreaListQueryParamsDTO) (*dtos.FocusAreaListResponseDTO, error) {
	return fas.repo.GetAll(params)
}

func (fas FocusAreaService) GetByID(id uint) (*models.FocusArea, error) {
	focusArea, err := fas.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return focusArea, nil

}

func (fas *FocusAreaService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.FocusAreaDetailDTO, error) {
	// Initialize DTO
	var focusAreaDTO dtos.FocusAreaDetailDTO

	// Get from repo
	focusArea, err := fas.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map to DTO
	err = mapstructure.Decode(focusArea, &focusAreaDTO)
	if err != nil {
		return nil, err
	}
	focusAreaDTO.CreatedAt = focusArea.CreatedAt
	focusAreaDTO.UpdatedAt = focusArea.UpdatedAt

	return &focusAreaDTO, nil
}

func (fas *FocusAreaService) Create(currentUserID uint, dto *dtos.FocusAreaFormDTO) (uint, error) {
	tx := fas.repo.BeginTransaction()

	// Check for unique title
	if exists, err := fas.repo.ExistsByTitle(dto.TitleEn, dto.TitleAr); err != nil {
		tx.Rollback()
		return 0, err
	} else if exists {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	// Build model
	model := &models.FocusArea{
		TitleEn:       dto.TitleEn,
		TitleAr:       dto.TitleAr,
		DescriptionEn: dto.DescriptionEn,
		DescriptionAr: dto.DescriptionAr,
		CreatedBYID:   currentUserID,
	}

	// Create
	id, err := fas.repo.Create(model, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	// Async activity log
	go func() {
		fullDetails, err := fas.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch strategic objective details", map[string]interface{}{"error": err.Error()})
			return
		}

		err = fas.activityLogSvc.LogActivity(
			models.CreateAction,
			models.FocusArea{}.TableName(),
			id,
			fullDetails,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.FocusArea{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (fas *FocusAreaService) Update(id uint, currentUserID uint, dto *dtos.FocusAreaFormDTO) error {
	existing, err := fas.GetByID(id)
	if err != nil {
		return err
	}

	// Check title uniqueness if changed
	if dto.TitleEn != existing.TitleEn || dto.TitleAr != existing.TitleAr {
		if exists, err := fas.repo.ExistsByTitleExcludingID(id, dto.TitleEn, dto.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	// Update fields
	existing.TitleEn = dto.TitleEn
	existing.TitleAr = dto.TitleAr
	existing.DescriptionEn = dto.DescriptionEn
	existing.DescriptionAr = dto.DescriptionAr
	existing.Status = dto.Status
	existing.UpdatedBYID = &currentUserID

	// Save
	err = fas.repo.Update(existing)
	if err != nil {
		return err
	}

	// Log activity asynchronously
	go func() {
		fullDetails, err := fas.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch strategic objective details for activity log", map[string]interface{}{"error": err.Error()})
			return
		}

		err = fas.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.FocusArea{}.TableName(),
			id,
			fullDetails,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log update activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.FocusArea{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}

func (fas *FocusAreaService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := fas.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = fas.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := fas.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.FocusArea{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (fas *FocusAreaService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := fas.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := fas.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.FocusArea{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (fas *FocusAreaService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := fas.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := fas.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.FocusArea{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}
