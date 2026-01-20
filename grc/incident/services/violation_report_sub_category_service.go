package services

import (
	"fmt"
	"jethur/grc/incident/dtos"
	"jethur/grc/incident/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"
)

// Interface
type ViolationReportSubCategoryServiceInterface interface {
	GetAll(params *dtos.ViolationReportSubCategoryListQueryParamsDTO) ([]*dtos.ViolationReportSubCategoryListDTO, error)
	GetByID(id uint) (*models.ViolationReportSubCategory, error)
	Create(currentUserID uint, vrscDTO *dtos.ViolationReportSubCategoryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, vrscDTO *dtos.ViolationReportSubCategoryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Implementation
type ViolationReportSubCategoryService struct {
	repo                       repositories.ViolationReportSubCategoryRepositoryInterface
	activityLogSvc             services.ActivityLogServiceInterface
	violationReportCategorySvc ViolationReportCategoryServiceInterface
}

func NewViolationReportSubCategoryService(repo repositories.ViolationReportSubCategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface, violationReportCategoryServiceInterface ViolationReportCategoryServiceInterface) ViolationReportSubCategoryServiceInterface {
	return &ViolationReportSubCategoryService{
		repo:                       repo,
		activityLogSvc:             activityLogSvc,
		violationReportCategorySvc: violationReportCategoryServiceInterface,
	}
}

func (vrscs *ViolationReportSubCategoryService) GetAll(params *dtos.ViolationReportSubCategoryListQueryParamsDTO) ([]*dtos.ViolationReportSubCategoryListDTO, error) {
	return vrscs.repo.GetAll(params)
}

func (vrscs *ViolationReportSubCategoryService) GetByID(id uint) (*models.ViolationReportSubCategory, error) {
	return vrscs.repo.GetByID(id)
}

func (vrscs *ViolationReportSubCategoryService) Create(currentUserID uint, vrscDTO *dtos.ViolationReportSubCategoryFormDTO) (uint, error) {
	if err := vrscs.validateViolationReportSubCategoryFormDTO(vrscDTO, nil); err != nil {
		return 0, err
	}

	vrsc := &models.ViolationReportSubCategory{
		ViolationReportCategoryID: vrscDTO.ViolationReportCategoryID,
		TitleEn:                   vrscDTO.TitleEn,
		TitleAr:                   vrscDTO.TitleAr,
		DescriptionEn:             vrscDTO.DescriptionEn,
		DescriptionAr:             vrscDTO.DescriptionAr,
		Status:                    vrscDTO.Status,
	}

	id, err := vrscs.repo.Create(vrsc)
	if err != nil {
		return 0, err
	}

	go func() {
		createdSubCategory, err := vrscs.repo.GetByID(id)
		if err != nil {
			return
		}
		err = vrscs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ViolationReportSubCategory{}.TableName(),
			id,
			createdSubCategory,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return id, nil
}

func (vrscs *ViolationReportSubCategoryService) Update(id uint, currentUserID uint, vrscDTO *dtos.ViolationReportSubCategoryFormDTO) error {
	existing, err := vrscs.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := vrscs.validateViolationReportSubCategoryFormDTO(vrscDTO, &id); err != nil {
		return err
	}

	updatedSubCategory := &models.ViolationReportSubCategory{
		ViolationReportCategoryID: vrscDTO.ViolationReportCategoryID,
		TitleEn:                   vrscDTO.TitleEn,
		TitleAr:                   vrscDTO.TitleAr,
		DescriptionEn:             vrscDTO.DescriptionEn,
		DescriptionAr:             vrscDTO.DescriptionAr,
		Status:                    vrscDTO.Status,
	}

	existing.ViolationReportCategoryID = vrscDTO.ViolationReportCategoryID
	existing.TitleEn = vrscDTO.TitleEn
	existing.TitleAr = vrscDTO.TitleAr
	existing.DescriptionEn = vrscDTO.DescriptionEn
	existing.DescriptionAr = vrscDTO.DescriptionAr
	existing.Status = vrscDTO.Status

	err = vrscs.repo.Update(existing)
	if err != nil {
		return err
	}

	fieldsToCompare := []string{"ViolationReportCategoryID", "TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "Status"}
	changes, err := utility.CompareFields(existing, updatedSubCategory, fieldsToCompare)
	if err != nil {
		log.Println("compare error:", err)
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("json error:", err)
	}

	go func() {
		err := vrscs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ViolationReportSubCategory{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log update activity: %v", err)
		}
	}()

	return nil
}

func (vrscs *ViolationReportSubCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deletedViolationReportSubCategories []map[string]interface{}

	for _, id := range ids {
		vrsc, err := vrscs.repo.GetByID(id)
		if err != nil {
			continue
		}

		err = vrscs.repo.Delete(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
			deletedViolationReportSubCategories = append(deletedViolationReportSubCategories, map[string]interface{}{
				"id":       id,
				"title_en": vrsc.TitleEn,
				"title_ar": vrsc.TitleAr,
			})
		}
	}

	if len(successfulIDs) > 0 {
		go func() {
			err := vrscs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ViolationReportSubCategory{}.TableName(),
				successfulIDs,
				deletedViolationReportSubCategories,
				currentUserID,
				nil,
			)
			if err != nil {
				log.Printf("Failed to log delete activity: %v", err)
			}
		}()
	}

	return failedIDs, nil
}

func (vrscs *ViolationReportSubCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIDs, err := vrscs.repo.Activate(ids)

	if len(successfulIDs) > 0 {
		go func() {
			err := vrscs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ViolationReportSubCategory{}.TableName(),
				successfulIDs,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				log.Printf("Failed to log activate activity: %v", err)
			}
		}()
	}

	return err
}

func (vrscs *ViolationReportSubCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIDs, err := vrscs.repo.Deactivate(ids)

	if len(successfulIDs) > 0 {
		go func() {
			err := vrscs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ViolationReportSubCategory{}.TableName(),
				successfulIDs,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				log.Printf("Failed to log deactivate activity: %v", err)
			}
		}()
	}

	return err
}

func (vrscs *ViolationReportSubCategoryService) validateViolationReportSubCategoryFormDTO(vrscDTO *dtos.ViolationReportSubCategoryFormDTO, id *uint) error {
	if _, err := vrscs.violationReportCategorySvc.GetByID(vrscDTO.ViolationReportCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("violation report category does not exist"), "not_exist", "violation report category")
	}

	if id == nil {
		if exists, err := vrscs.repo.ExistsByTitle(vrscDTO.TitleEn, vrscDTO.TitleAr, vrscDTO.ViolationReportCategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	} else {
		if exists, err := vrscs.repo.ExistsByTitleExcludingID(*id, vrscDTO.TitleEn, vrscDTO.TitleAr, vrscDTO.ViolationReportCategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}
