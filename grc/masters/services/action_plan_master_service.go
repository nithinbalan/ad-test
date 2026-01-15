package services

import (
	"errors"
	"fmt"
	activityLogServices "jethur/grc/logs/services"
	"jethur/grc/masters/dtos"
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ActionPlanMasterServiceInterface interface {
	GetAll(params *dtos.ActionPlanMasterListQueryParamsDTO) (*dtos.ActionPlanMasterListResponseDTO, error)
	GetByID(id uint) (*models.ActionPlanMaster, error)
	GetFullDetailsByID(id uint) (*dtos.ActionPlanMasterDetailsDTO, error)
	GetModuleByID(id uint) (*models.Module, error)
	GetIncidentCategoryByID(id uint) (*models.IncidentCategory, error)
	Create(currentUserID uint, formDTO *dtos.ActionPlanMasterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *dtos.ActionPlanMasterFormDTO) error
	Delete(ids []uint,currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type ActionPlanMasterService struct {
	repo                 repositories.ActionPlanMasterRepositoryInterface
	activityLogSvc       activityLogServices.ActivityLogServiceInterface
	severityLevelService SeverityLevelServiceInterface
}

func NewActionPlanMasterService(
	repo repositories.ActionPlanMasterRepositoryInterface,
	activityLogSvc activityLogServices.ActivityLogServiceInterface,
	severityLevelService SeverityLevelServiceInterface,
) ActionPlanMasterServiceInterface {
	return &ActionPlanMasterService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		severityLevelService: severityLevelService,
	}
}

func (apms *ActionPlanMasterService) GetAll(params *dtos.ActionPlanMasterListQueryParamsDTO) (*dtos.ActionPlanMasterListResponseDTO, error) {
	return apms.repo.GetAll(params)
}

func (apms *ActionPlanMasterService) GetByID(id uint) (*models.ActionPlanMaster, error) {
	actionPlanMaster, err := apms.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return actionPlanMaster, nil
}

func (apms *ActionPlanMasterService) GetFullDetailsByID(id uint) (*dtos.ActionPlanMasterDetailsDTO, error) {
	var actionPlanMasterDetailsDTO dtos.ActionPlanMasterDetailsDTO

	actionPlanMaster, err := apms.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(actionPlanMaster, &actionPlanMasterDetailsDTO)
	if err != nil {
		return nil, err
	}

	return &actionPlanMasterDetailsDTO, nil
}

func (apms *ActionPlanMasterService) GetModuleByID(id uint) (*models.Module, error) {
	module, err := apms.repo.GetModuleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return module, nil
}

func (apms *ActionPlanMasterService) GetIncidentCategoryByID(id uint) (*models.IncidentCategory, error) {
	incidentCategory, err := apms.repo.GetIncidentCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incidentCategory, nil
}

func (apms *ActionPlanMasterService) Create(currentUserID uint, formDTO *dtos.ActionPlanMasterFormDTO) (uint, error) {

	if err := apms.validateActionPlanMasterFormDTO(formDTO, nil); err != nil {
		return 0, err
	}

	actionPlanMaster := &models.ActionPlanMaster{
		ModuleID:           formDTO.ModuleID,
		IncidentCategoryID: formDTO.IncidentCategoryID,
		SeverityLevelID:    formDTO.SeverityLevelID,
		Title:              formDTO.Title,
		Description:        formDTO.Description,
	}

	id, err := apms.repo.Create(actionPlanMaster)
	if err != nil {
		return 0, err
	}

	go func() {
		createdActionPlan, err := apms.repo.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		err = apms.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ActionPlanMaster{}.TableName(),
			id,
			createdActionPlan,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ActionPlanMaster{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (apms *ActionPlanMasterService) Update(id uint, currentUserID uint, form *dtos.ActionPlanMasterFormDTO) error {
	existingActionPlanMaster, err := apms.GetByID(id)
	if err != nil {
		return err
	}
	if err := apms.validateActionPlanMasterFormDTO(form, &id); err != nil {
		return err
	}

	existingActionPlanMaster.ModuleID = form.ModuleID
	existingActionPlanMaster.IncidentCategoryID = form.IncidentCategoryID
	existingActionPlanMaster.SeverityLevelID = form.SeverityLevelID
	existingActionPlanMaster.Title = form.Title
	existingActionPlanMaster.Description = form.Description

	if err := apms.repo.Update(existingActionPlanMaster); err != nil {
		return err
	}

	go func() {
		updatedActionPlan, err := apms.repo.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		err = apms.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ActionPlanMaster{}.TableName(),
			id,
			updatedActionPlan,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.ActionPlanMaster{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (apms *ActionPlanMasterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := apms.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var deletedActionPlans []map[string]interface{}

	for _, id := range ids {

		actionPlanMaster, err := apms.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if err := apms.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		successfulIds = append(successfulIds, id)

		deletedActionPlans = append(deletedActionPlans, map[string]interface{}{
			"id":    id,
			"title": actionPlanMaster.Title,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := apms.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ActionPlanMaster{}.TableName(),
				successfulIds,
				deletedActionPlans,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log ActionPlanMaster delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}



func (apms *ActionPlanMasterService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := apms.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := apms.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ActionPlanMaster{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (apms *ActionPlanMasterService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := apms.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := apms.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ActionPlanMaster{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}


func (apms *ActionPlanMasterService) validateActionPlanMasterFormDTO(form *dtos.ActionPlanMasterFormDTO, id *uint) error {

	_, err := apms.GetModuleByID(uint(form.ModuleID))
	if err != nil {
			return utility.NewValidationError(
				fmt.Errorf("module not found"),"not_exist","module_id")
	}

	if form.IncidentCategoryID != nil {
		_, err := apms.GetIncidentCategoryByID(uint(*form.IncidentCategoryID))
		if err != nil {
				return utility.NewValidationError(
					fmt.Errorf("incident category not found"),"not_exist","incident_category_id")
		}
	}

	if form.SeverityLevelID != nil {
		_, err := apms.severityLevelService.GetByID(uint(*form.SeverityLevelID))
		if err != nil {
				return utility.NewValidationError(
					fmt.Errorf("severity level not found"),"not_exist","severity_level_id")
		}
	}

	if id == nil {
		exists, err := apms.repo.ExistsByTitle(form.Title)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(
				fmt.Errorf("action plan title must be unique"),"unique","title")
		}
	} else {
		exists, err := apms.repo.ExistsByTitleExcludingID(*id, form.Title, form.ModuleID)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(fmt.Errorf("action plan title must be unique"),"unique","title")
		}
	}

	return nil
}
