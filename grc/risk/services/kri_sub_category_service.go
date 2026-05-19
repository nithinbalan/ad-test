package services

import (
	"errors"
	"fmt"
	activityLogServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type KRISubCategoryServiceInterface interface {
	GetAll(params *dtos.KRISubCategoryListQueryParamsDTO) ([]*dtos.KRISubCategoryDTO, error)
	GetByID(id uint) (*models.KRISubCategory, error)
	Create(currentUserID uint, kriSubCategoryFormDTO *dtos.KRISubCategoryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, kriSubCategoryFormDTO *dtos.KRISubCategoryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type KRISubCategoryService struct {
	repo            repositories.KRISubCategoryRepositoryInterface
	kriCategoryRepo repositories.KRICategoryRepositoryInterface
	activityLogSvc  activityLogServices.ActivityLogServiceInterface
}

func NewKRISubCategoryService(repo repositories.KRISubCategoryRepositoryInterface, kriCategoryRepo repositories.KRICategoryRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface) KRISubCategoryServiceInterface {
	return &KRISubCategoryService{repo: repo, kriCategoryRepo: kriCategoryRepo, activityLogSvc: activityLogSvc}
}

func (s *KRISubCategoryService) GetAll(params *dtos.KRISubCategoryListQueryParamsDTO) ([]*dtos.KRISubCategoryDTO, error) {
	return s.repo.GetAll(params)
}

func (s KRISubCategoryService) GetByID(id uint) (*models.KRISubCategory, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return item, nil
}

func (s *KRISubCategoryService) Create(currentUserID uint, kriSubCategoryFormDTO *dtos.KRISubCategoryFormDTO) (uint, error) {
	if err := s.validateKRISubCategoryFormDTO(kriSubCategoryFormDTO, nil); err != nil {
		return 0, err
	}

	item := &models.KRISubCategory{
		KRICategoryID: kriSubCategoryFormDTO.KRICategoryID,
		TitleEn:       kriSubCategoryFormDTO.TitleEn,
		TitleAr:       kriSubCategoryFormDTO.TitleAr,
		DescriptionEn: kriSubCategoryFormDTO.DescriptionEn,
		DescriptionAr: kriSubCategoryFormDTO.DescriptionAr,
		Status:        "active",
	}

	createdItem, err := s.repo.Create(item)
	if err != nil {
		return 0, err
	}

	go func() {
		reloadedItem, getErr := s.repo.GetByID(createdItem.ID)
		if getErr != nil {
			return
		}
		err := s.activityLogSvc.LogActivity(
			models.CreateAction,
			models.KRISubCategory{}.TableName(),
			createdItem.ID,
			reloadedItem,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.KRISubCategory{}.TableName(),
				"recordID":  createdItem.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return createdItem.ID, nil
}

func (s *KRISubCategoryService) Update(id uint, currentUserID uint, kriSubCategoryFormDTO *dtos.KRISubCategoryFormDTO) error {
	existingItem, err := s.GetByID(id)
	if err != nil {
		return err
	}

	itemForCompare, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.validateKRISubCategoryFormDTO(kriSubCategoryFormDTO, &id); err != nil {
		return err
	}

	existingItem.KRICategoryID = kriSubCategoryFormDTO.KRICategoryID
	existingItem.TitleEn = kriSubCategoryFormDTO.TitleEn
	existingItem.TitleAr = kriSubCategoryFormDTO.TitleAr
	existingItem.DescriptionEn = kriSubCategoryFormDTO.DescriptionEn
	existingItem.DescriptionAr = kriSubCategoryFormDTO.DescriptionAr

	_, err = s.repo.Update(id, existingItem)
	if err != nil {
		return err
	}

	fieldsToCompare := []string{"KRICategoryID", "TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr"}
	changes, err := utility.CompareFields(itemForCompare, existingItem, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.KRISubCategory{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

	return nil
}

func (s *KRISubCategoryService) validateKRISubCategoryFormDTO(kriSubCategoryFormDTO *dtos.KRISubCategoryFormDTO, id *uint) error {
	if _, err := s.kriCategoryRepo.GetByID(kriSubCategoryFormDTO.KRICategoryID); err != nil {
		if _, ok := err.(*utility.ResourceNotFoundError); ok {
			return utility.NewValidationError(fmt.Errorf("kri category does not exist"), "not_exist", "kri_category_id")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewValidationError(fmt.Errorf("kri category does not exist"), "not_exist", "kri_category_id")
		}
		return err
	}

	if id == nil {
		if exists, err := s.repo.ExistsByTitle(kriSubCategoryFormDTO.KRICategoryID, kriSubCategoryFormDTO.TitleEn, kriSubCategoryFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	} else {
		if exists, err := s.repo.ExistsByTitleExcludingID(*id, kriSubCategoryFormDTO.KRICategoryID, kriSubCategoryFormDTO.TitleEn, kriSubCategoryFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}

func (s *KRISubCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedItems []map[string]interface{}
	for _, id := range ids {
		item, err := s.repo.GetByID(id)
		if err != nil {
			continue
		}

		err = s.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedItems = append(deletedItems, map[string]interface{}{
			"id":       id,
			"title_en": item.TitleEn,
			"title_ar": item.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := s.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.KRISubCategory{}.TableName(),
				successfulIds,
				deletedItems,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (s *KRISubCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := s.repo.Activate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := s.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.KRISubCategory{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk activate activity: %v\n", err)
			}
		}()
	}

	return err
}

func (s *KRISubCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := s.repo.Deactivate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := s.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.KRISubCategory{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk deactivate activity: %v\n", err)
			}
		}()
	}

	return err
}
