package services

import (
	"fmt"
	libraryDTOs "jethur/grc/libraries/dtos"
	"jethur/grc/libraries/models"
	"jethur/grc/libraries/repositories"

	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/utility"
)

type KRILibraryServiceInterface interface {
	GetAll(params *libraryDTOs.KRILibraryListQueryParamsDTO) (*libraryDTOs.KRILibraryListResponseDTO, error)
	GetByID(id uint) (*models.KRILibrary, error)
	Create(currentUserID uint, libraryFormDTO *libraryDTOs.KRILibraryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, kriLibrary *libraryDTOs.KRILibraryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type KRILibraryService struct {
	repo               repositories.KRILibraryRepositoryInterface
	kriCategorySvc     KRICategoryServiceInterface
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
}

func NewKRILibraryService(
	repo repositories.KRILibraryRepositoryInterface,
	kriCategorySvc KRICategoryServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) KRILibraryServiceInterface {
	return &KRILibraryService{
		repo:               repo,
		kriCategorySvc:     kriCategorySvc,
		reviewFrequencySvc: reviewFrequencySvc,
		activityLogSvc:     activityLogSvc,
	}
}

func (kls *KRILibraryService) GetAll(params *libraryDTOs.KRILibraryListQueryParamsDTO) (*libraryDTOs.KRILibraryListResponseDTO, error) {
	return kls.repo.GetAll(params)
}

func (kls KRILibraryService) GetByID(id uint) (*models.KRILibrary, error) {
	return kls.repo.GetByID(id)
}

func (kls *KRILibraryService) Create(currentUserID uint, libraryFormDTO *libraryDTOs.KRILibraryFormDTO) (uint, error) {

	if err := kls.validateKRILibraryFormDTO(libraryFormDTO, nil); err != nil {
		return 0, err
	}

	kriLibrary := &models.KRILibrary{
		KRICategoryID:     libraryFormDTO.KRICategoryID,
		ReviewFrequencyID: libraryFormDTO.ReviewFrequencyID,
		Title:             libraryFormDTO.Title,
		Description:       libraryFormDTO.Description,
		CalculationType:   libraryFormDTO.CalculationType,
		Source:            libraryFormDTO.Source,
		Status:            libraryFormDTO.Status,
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
			libraryFormDTO,
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

func (kls *KRILibraryService) Update(id uint, currentUserID uint, libraryFormDTO *libraryDTOs.KRILibraryFormDTO) error {

	existingKRILibrary, err := kls.repo.GetByID(id)

	if err != nil {
		return err
	}

	if err := kls.validateKRILibraryFormDTO(libraryFormDTO, &id); err != nil {
		return err
	}

	existingKRILibrary.KRICategoryID = libraryFormDTO.KRICategoryID
	existingKRILibrary.ReviewFrequencyID = libraryFormDTO.ReviewFrequencyID
	existingKRILibrary.Title = libraryFormDTO.Title
	existingKRILibrary.Description = libraryFormDTO.Description
	existingKRILibrary.CalculationType = libraryFormDTO.CalculationType
	existingKRILibrary.Source = libraryFormDTO.Source
	existingKRILibrary.Status = libraryFormDTO.Status

	err = kls.repo.Update(existingKRILibrary)

	if err != nil {
		return err
	}

	go func() {
		err := kls.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			models.KRILibrary{}.TableName(),
			id,
			libraryFormDTO,
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

func (kls *KRILibraryService) validateKRILibraryFormDTO(libraryFormDTO *libraryDTOs.KRILibraryFormDTO, id *uint) error {
	if _, err := kls.kriCategorySvc.GetByID(libraryFormDTO.KRICategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("KRI category does not exist"), "not_exist", "kri category")
	}

	if libraryFormDTO.ReviewFrequencyID != nil {
		if _, err := kls.reviewFrequencySvc.GetByID(*libraryFormDTO.ReviewFrequencyID); err != nil {
			return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
		}
	}

	if id == nil {
		if exists, err := kls.repo.ExistsByTitle(libraryFormDTO.Title, libraryFormDTO.KRICategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	} else {
		fmt.Println("DI", &id)
		fmt.Println("libraryFormDTO.ID", libraryFormDTO.ID)
		if exists, err := kls.repo.ExistsByTitleExcludingID(libraryFormDTO.ID, libraryFormDTO.Title, libraryFormDTO.KRICategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}
