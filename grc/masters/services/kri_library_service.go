package services

import (
	"fmt"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	"jethur/grc/masters/models"
	"jethur/grc/masters/repositories"
	"jethur/utility"
)

type KRILibraryServiceInterface interface {
	GetAll(params *masterDTOs.KRILibraryListQueryParamsDTO) (*masterDTOs.KRILibraryListResponseDTO, error)
	GetByID(id uint) (*models.KRILibrary, error)
	Create(currentUserID uint, kriLibraryFormDTO *masterDTOs.KRILibraryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, kriLibrary *masterDTOs.KRILibraryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type KRILibraryService struct {
	repo               repositories.KRILibraryRepositoryInterface
	kriCategorySvc     KRICategoryServiceInterface
	reviewFrequencySvc ReviewFrequencyServiceInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
}

func NewKRILibraryService(
	repo repositories.KRILibraryRepositoryInterface,
	kriCategorySvc KRICategoryServiceInterface,
	reviewFrequencySvc ReviewFrequencyServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) KRILibraryServiceInterface {
	return &KRILibraryService{
		repo:               repo,
		kriCategorySvc:     kriCategorySvc,
		reviewFrequencySvc: reviewFrequencySvc,
		activityLogSvc:     activityLogSvc,
	}
}

func (kls *KRILibraryService) GetAll(params *masterDTOs.KRILibraryListQueryParamsDTO) (*masterDTOs.KRILibraryListResponseDTO, error) {
	return kls.repo.GetAll(params)
}

func (kls KRILibraryService) GetByID(id uint) (*models.KRILibrary, error) {
	return kls.repo.GetByID(id)
}

func (kls *KRILibraryService) Create(currentUserID uint, kriLibraryFormDTO *masterDTOs.KRILibraryFormDTO) (uint, error) {

	if err := kls.validateKRILibraryFormDTO(kriLibraryFormDTO, nil); err != nil {
		return 0, err
	}

	kriLibrary := &models.KRILibrary{
		KRICategoryID:     kriLibraryFormDTO.KRICategoryID,
		ReviewFrequencyID: kriLibraryFormDTO.ReviewFrequencyID,
		Title:             kriLibraryFormDTO.Title,
		Description:       kriLibraryFormDTO.Description,
		CalculationType:   kriLibraryFormDTO.CalculationType,
		Source:            kriLibraryFormDTO.Source,
		Status:            kriLibraryFormDTO.Status,
	}

	id, err := kls.repo.Create(kriLibrary)
	if err != nil {
		return 0, err
	}

	go func() {
		err := kls.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.KRILibrary{}.TableName(),
			id,
			kriLibraryFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": models.KRILibrary{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (kls *KRILibraryService) Update(id uint, currentUserID uint, kriLibraryFormDTO *masterDTOs.KRILibraryFormDTO) error {

	existingKRILibrary, err := kls.repo.GetByID(id)

	if err != nil {
		return err
	}

	if err := kls.validateKRILibraryFormDTO(kriLibraryFormDTO, &id); err != nil {
		return err
	}

	existingKRILibrary.KRICategoryID = kriLibraryFormDTO.KRICategoryID
	existingKRILibrary.ReviewFrequencyID = kriLibraryFormDTO.ReviewFrequencyID
	existingKRILibrary.Title = kriLibraryFormDTO.Title
	existingKRILibrary.Description = kriLibraryFormDTO.Description
	existingKRILibrary.CalculationType = kriLibraryFormDTO.CalculationType
	existingKRILibrary.Source = kriLibraryFormDTO.Source
	existingKRILibrary.Status = kriLibraryFormDTO.Status

	err = kls.repo.Update(existingKRILibrary)

	if err != nil {
		return err
	}

	go func() {
		err := kls.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			models.KRILibrary{}.TableName(),
			id,
			kriLibraryFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.UpdateAction,
				"tableName": models.KRILibrary{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (kls *KRILibraryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := kls.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = kls.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := kls.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.KRILibrary{}.TableName(),
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

func (kls *KRILibraryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := kls.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := kls.activityLogSvc.LogBulkActivities(
				logModels.ActivateAction,
				models.KRILibrary{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (kls *KRILibraryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := kls.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := kls.activityLogSvc.LogBulkActivities(
				logModels.DeactivateAction,
				models.KRILibrary{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (kls *KRILibraryService) validateKRILibraryFormDTO(kriLibraryFormDTO *masterDTOs.KRILibraryFormDTO, id *uint) error {
	if _, err := kls.kriCategorySvc.GetByID(kriLibraryFormDTO.KRICategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("KRI category does not exist"), "not_exist", "kri category")
	}

	if kriLibraryFormDTO.ReviewFrequencyID != nil {
		if _, err := kls.reviewFrequencySvc.GetByID(*kriLibraryFormDTO.ReviewFrequencyID); err != nil {
			return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
		}
	}

	if id == nil {
		if exists, err := kls.repo.ExistsByTitle(kriLibraryFormDTO.Title, kriLibraryFormDTO.KRICategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	} else {
		if exists, err := kls.repo.ExistsByTitleExcludingID(*id, kriLibraryFormDTO.Title, kriLibraryFormDTO.KRICategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}
