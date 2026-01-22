package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	settingServices "jethur/grc/settings/services"
	"jethur/grc/violation-report/dtos"
	"jethur/grc/violation-report/repositories"
	"jethur/utility"
	"log"
)

// Interface
type ViolationReportCategoryServiceInterface interface {
	GetAll(params *dtos.ViolationReportCategoryListQueryParamsDTO) ([]*dtos.ViolationReportCategoryListDTO, error)
	GetByID(id uint) (*models.ViolationReportCategory, error)
	Create(currentUserID uint, vrcDTO *dtos.ViolationReportCategoryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, vrcDTO *dtos.ViolationReportCategoryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Implementation
type ViolationReportCategoryService struct {
	repo                repositories.ViolationReportCategoryRepositoryInterface
	activityLogSvc      services.ActivityLogServiceInterface
	organizationUnitSvc settingServices.OrganizationUnitServiceInterface
}

func NewViolationReportCategoryService(repo repositories.ViolationReportCategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface, organizationUnitSvc settingServices.OrganizationUnitServiceInterface) ViolationReportCategoryServiceInterface {
	return &ViolationReportCategoryService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		organizationUnitSvc: organizationUnitSvc,
	}
}

func (vrcs *ViolationReportCategoryService) GetAll(params *dtos.ViolationReportCategoryListQueryParamsDTO) ([]*dtos.ViolationReportCategoryListDTO, error) {
	return vrcs.repo.GetAll(params)
}

func (vrcs *ViolationReportCategoryService) GetByID(id uint) (*models.ViolationReportCategory, error) {
	return vrcs.repo.GetByID(id)
}

func (vrcs *ViolationReportCategoryService) Create(currentUserID uint, vrcDTO *dtos.ViolationReportCategoryFormDTO) (uint, error) {
	if err := vrcs.validateViolationReportCategoryFormDTO(vrcDTO, nil); err != nil {
		return 0, err
	}

	vrc := &models.ViolationReportCategory{
		DepartmentID:  vrcDTO.DepartmentID,
		TitleEn:       vrcDTO.TitleEn,
		TitleAr:       vrcDTO.TitleAr,
		DescriptionEn: vrcDTO.DescriptionEn,
		DescriptionAr: vrcDTO.DescriptionAr,
		Status:        vrcDTO.Status,
	}

	id, err := vrcs.repo.Create(vrc)
	if err != nil {
		return 0, err
	}

	go func() {
		createdCategory, err := vrcs.repo.GetByID(id)
		if err != nil {
			return
		}
		err = vrcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ViolationReportCategory{}.TableName(),
			id,
			createdCategory,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return id, nil
}

func (vrcs *ViolationReportCategoryService) Update(id uint, currentUserID uint, vrcDTO *dtos.ViolationReportCategoryFormDTO) error {
	existing, err := vrcs.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := vrcs.validateViolationReportCategoryFormDTO(vrcDTO, &id); err != nil {
		return err
	}

	updatedCategory := &models.ViolationReportCategory{
		DepartmentID:  vrcDTO.DepartmentID,
		TitleEn:       vrcDTO.TitleEn,
		TitleAr:       vrcDTO.TitleAr,
		DescriptionEn: vrcDTO.DescriptionEn,
		DescriptionAr: vrcDTO.DescriptionAr,
		Status:        vrcDTO.Status,
	}

	existing.DepartmentID = vrcDTO.DepartmentID
	existing.TitleEn = vrcDTO.TitleEn
	existing.TitleAr = vrcDTO.TitleAr
	existing.DescriptionEn = vrcDTO.DescriptionEn
	existing.DescriptionAr = vrcDTO.DescriptionAr
	existing.Status = vrcDTO.Status

	err = vrcs.repo.Update(existing)
	if err != nil {
		return err
	}

	fieldsToCompare := []string{"DepartmentID", "TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "Status"}
	changes, err := utility.CompareFields(existing, updatedCategory, fieldsToCompare)
	if err != nil {
		log.Println("compare error:", err)
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("json error:", err)
	}

	go func() {
		err := vrcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ViolationReportCategory{}.TableName(),
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

func (vrcs *ViolationReportCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deletedViolationReportCategories []map[string]interface{}

	for _, id := range ids {
		vrc, err := vrcs.repo.GetByID(id)
		if err != nil {
			continue
		}

		err = vrcs.repo.Delete(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
			deletedViolationReportCategories = append(deletedViolationReportCategories, map[string]interface{}{
				"id":       id,
				"title_en": vrc.TitleEn,
				"title_ar": vrc.TitleAr,
			})
		}
	}

	if len(successfulIDs) > 0 {
		go func() {
			err := vrcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ViolationReportCategory{}.TableName(),
				successfulIDs,
				deletedViolationReportCategories,
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

func (vrcs *ViolationReportCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIDs, err := vrcs.repo.Activate(ids)

	if len(successfulIDs) > 0 {
		go func() {
			err := vrcs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ViolationReportCategory{}.TableName(),
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

func (vrcs *ViolationReportCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIDs, err := vrcs.repo.Deactivate(ids)

	if len(successfulIDs) > 0 {
		go func() {
			err := vrcs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ViolationReportCategory{}.TableName(),
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

func (vrcs *ViolationReportCategoryService) validateViolationReportCategoryFormDTO(vrcDTO *dtos.ViolationReportCategoryFormDTO, id *uint) error {
	if _, err := vrcs.organizationUnitSvc.GetByID(vrcDTO.DepartmentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("department does not exist"), "not_exist", "department")
	}

	if id == nil {
		if exists, err := vrcs.repo.ExistsByTitle(vrcDTO.TitleEn, vrcDTO.TitleAr, vrcDTO.DepartmentID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	} else {
		if exists, err := vrcs.repo.ExistsByTitleExcludingID(*id, vrcDTO.TitleEn, vrcDTO.TitleAr, vrcDTO.DepartmentID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}
