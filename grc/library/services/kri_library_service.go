package services

import (
	"fmt"

	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	riskServices "jethur/grc/risk/services"

	"jethur/utility"
)

type KRILibraryServiceInterface interface {
	SetKriCategoryService(svc riskServices.KRICategoryServiceInterface)

	GetAll(params *libraryDTOs.KRILibraryListQueryParamsDTO) (*libraryDTOs.KRILibraryListResponseDTO, error)
	GetByID(id uint) (*models.KRILibrary, error)
	Create(currentUserID uint, kriLibraryFormDTO *libraryDTOs.KRILibraryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, kriLibrary *libraryDTOs.KRILibraryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
	GetKRICategoryCounts() ([]libraryDTOs.KRICategoryCountDTO, error)
}

// Service and implementations
type KRILibraryService struct {
	repo               repositories.KRILibraryRepositoryInterface
	kriCategorySvc     riskServices.KRICategoryServiceInterface
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
}

func NewKRILibraryService(
	repo repositories.KRILibraryRepositoryInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) KRILibraryServiceInterface {
	return &KRILibraryService{
		repo:               repo,
		reviewFrequencySvc: reviewFrequencySvc,
		activityLogSvc:     activityLogSvc,
	}
}

func (kls *KRILibraryService) SetKriCategoryService(svc riskServices.KRICategoryServiceInterface) {
	kls.kriCategorySvc = svc
}

func (kls *KRILibraryService) GetAll(params *libraryDTOs.KRILibraryListQueryParamsDTO) (*libraryDTOs.KRILibraryListResponseDTO, error) {
	return kls.repo.GetAll(params)
}

func (kls KRILibraryService) GetByID(id uint) (*models.KRILibrary, error) {
	return kls.repo.GetByID(id)
}

func (kls *KRILibraryService) Create(currentUserID uint, kriLibraryFormDTO *libraryDTOs.KRILibraryFormDTO) (uint, error) {
	// if err := kls.validateKRILibraryFormDTO(kriLibraryFormDTO, nil); err != nil {
	// 	return 0, err
	// }

	kriLibrary := &models.KRILibrary{
		KRICategory: kriLibraryFormDTO.KRICategory,
		Title:       kriLibraryFormDTO.Title,
		Description: kriLibraryFormDTO.Description,
	}

	id, err := kls.repo.Create(kriLibrary)
	if err != nil {
		return 0, err
	}

	go func() {
		err := kls.activityLogSvc.LogActivity(
			models.CreateAction,
			models.KRILibrary{}.TableName(),
			id,
			kriLibraryFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
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

func (kls *KRILibraryService) Update(id uint, currentUserID uint, kriLibraryFormDTO *libraryDTOs.KRILibraryFormDTO) error {
	existingKRILibrary, err := kls.repo.GetByID(id)
	if err != nil {
		return err
	}

	// if err := kls.validateKRILibraryFormDTO(kriLibraryFormDTO, &id); err != nil {
	// 	return err
	// }

	existingKRILibrary.KRICategory = kriLibraryFormDTO.KRICategory
	existingKRILibrary.Title = kriLibraryFormDTO.Title
	existingKRILibrary.Description = kriLibraryFormDTO.Description

	err = kls.repo.Update(existingKRILibrary)
	if err != nil {
		return err
	}

	go func() {
		err := kls.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.KRILibrary{}.TableName(),
			id,
			kriLibraryFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
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
				models.DeleteAction, // Make sure this action constant is defined appropriately
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
				models.ActivateAction,
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
				models.DeactivateAction,
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

func (svc *KRILibraryService) GetKRICategoryCounts() ([]libraryDTOs.KRICategoryCountDTO, error) {
	return svc.repo.GetKRICategoryCounts()
}

// func (kls *KRILibraryService) validateKRILibraryFormDTO(kriLibraryFormDTO *libraryDTOs.KRILibraryFormDTO, id *uint) error {
// 	if _, err := kls.kriCategorySvc.GetByID(kriLibraryFormDTO.KRICategoryID); err != nil {
// 		return utility.NewValidationError(fmt.Errorf("KRI category does not exist"), "not_exist", "kri category")
// 	}

// 	if kriLibraryFormDTO.ReviewFrequencyID != nil {
// 		if _, err := kls.reviewFrequencySvc.GetByID(*kriLibraryFormDTO.ReviewFrequencyID); err != nil {
// 			return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
// 		}
// 	}

// 	if id == nil {
// 		if exists, err := kls.repo.ExistsByTitle(kriLibraryFormDTO.Title, kriLibraryFormDTO.KRICategoryID); err != nil {
// 			return err
// 		} else if exists {
// 			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
// 		}
// 	} else {
// 		if exists, err := kls.repo.ExistsByTitleExcludingID(*id, kriLibraryFormDTO.Title, kriLibraryFormDTO.KRICategoryID); err != nil {
// 			return err
// 		} else if exists {
// 			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
// 		}
// 	}

// 	return nil
// }
