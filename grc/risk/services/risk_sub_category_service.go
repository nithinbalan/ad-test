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

type RiskSubCategoryServiceInterface interface {
	GetAll(params *dtos.RiskSubCategoryListQueryParamsDTO) ([]*dtos.RiskSubCategoryDTO, error)
	GetByID(id uint) (*models.RiskSubCategory, error)
	Create(currentUserID uint, riskSubCategoryFormDTO *dtos.RiskSubCategoryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, riskSubCategoryFormDTO *dtos.RiskSubCategoryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type RiskSubCategoryService struct {
	repo           repositories.RiskSubCategoryRepositoryInterface
	riskCategoryRepo repositories.RiskCategoryRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
}

func NewRiskSubCategoryService(repo repositories.RiskSubCategoryRepositoryInterface, riskCategoryRepo repositories.RiskCategoryRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface) RiskSubCategoryServiceInterface {
	return &RiskSubCategoryService{repo: repo, riskCategoryRepo: riskCategoryRepo, activityLogSvc: activityLogSvc}
}

func (s *RiskSubCategoryService) GetAll(params *dtos.RiskSubCategoryListQueryParamsDTO) ([]*dtos.RiskSubCategoryDTO, error) {
	return s.repo.GetAll(params)
}

func (s RiskSubCategoryService) GetByID(id uint) (*models.RiskSubCategory, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return item, nil
}

func (s *RiskSubCategoryService) Create(currentUserID uint, riskSubCategoryFormDTO *dtos.RiskSubCategoryFormDTO) (uint, error) {
	if err := s.validateRiskSubCategoryFormDTO(riskSubCategoryFormDTO, nil); err != nil {
		return 0, err
	}

	item := &models.RiskSubCategory{
		RiskCategoryID: riskSubCategoryFormDTO.RiskCategoryID,
		TitleEn:        riskSubCategoryFormDTO.TitleEn,
		TitleAr:        riskSubCategoryFormDTO.TitleAr,
		DescriptionEn:  riskSubCategoryFormDTO.DescriptionEn,
		DescriptionAr:  riskSubCategoryFormDTO.DescriptionAr,
		Status:         riskSubCategoryFormDTO.Status,
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
			models.RiskSubCategory{}.TableName(),
			createdItem.ID,
			reloadedItem,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskSubCategory{}.TableName(),
				"recordID":  createdItem.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return createdItem.ID, nil
}

func (s *RiskSubCategoryService) Update(id uint, currentUserID uint, riskSubCategoryFormDTO *dtos.RiskSubCategoryFormDTO) error {
	existingItem, err := s.GetByID(id)
	if err != nil {
		return err
	}

	itemForCompare, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.validateRiskSubCategoryFormDTO(riskSubCategoryFormDTO, &id); err != nil {
		return err
	}

	existingItem.RiskCategoryID = riskSubCategoryFormDTO.RiskCategoryID
	existingItem.TitleEn = riskSubCategoryFormDTO.TitleEn
	existingItem.TitleAr = riskSubCategoryFormDTO.TitleAr
	existingItem.DescriptionEn = riskSubCategoryFormDTO.DescriptionEn
	existingItem.DescriptionAr = riskSubCategoryFormDTO.DescriptionAr
	existingItem.Status = riskSubCategoryFormDTO.Status

	_, err = s.repo.Update(id, existingItem)
	if err != nil {
		return err
	}

	fieldsToCompare := []string{"RiskCategoryID", "TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "Status"}
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
			models.RiskSubCategory{}.TableName(),
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

func (s *RiskSubCategoryService) validateRiskSubCategoryFormDTO(riskSubCategoryFormDTO *dtos.RiskSubCategoryFormDTO, id *uint) error {
	if _, err := s.riskCategoryRepo.GetByID(riskSubCategoryFormDTO.RiskCategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewValidationError(fmt.Errorf("risk category does not exist"), "not_exist", "risk_category_id")
		}
		return err
	}

	if id == nil {
		if exists, err := s.repo.ExistsByTitle(riskSubCategoryFormDTO.RiskCategoryID, riskSubCategoryFormDTO.TitleEn, riskSubCategoryFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	} else {
		if exists, err := s.repo.ExistsByTitleExcludingID(*id, riskSubCategoryFormDTO.RiskCategoryID, riskSubCategoryFormDTO.TitleEn, riskSubCategoryFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}

func (s *RiskSubCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
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
				models.RiskSubCategory{}.TableName(),
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

func (s *RiskSubCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := s.repo.Activate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := s.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.RiskSubCategory{}.TableName(),
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

func (s *RiskSubCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := s.repo.Deactivate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := s.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.RiskSubCategory{}.TableName(),
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
