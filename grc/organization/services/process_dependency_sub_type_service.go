package services

import (
	"errors"
	"fmt"
	activityLogServices "jethur/grc/logs/services"
	"jethur/grc/models"
	orgDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProcessDependencySubTypeServiceInterface interface {
	GetAll(params *orgDTOs.ProcessDependencySubTypeListQueryParamsDTO) ([]*orgDTOs.ProcessDependencySubTypeDTO, error)
	GetByID(id uint) (*models.ProcessDependencySubType, error)
	Create(currentUserID uint, dto *orgDTOs.ProcessDependencySubTypeFormDTO) (uint, error)
	Update(id uint, currentUserID uint, dto *orgDTOs.ProcessDependencySubTypeFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type ProcessDependencySubTypeService struct {
	repo                  repositories.ProcessDependencySubTypeRepositoryInterface
	processDependencyType ProcessDependencyTypeServiceInterface
	activityLogSvc        activityLogServices.ActivityLogServiceInterface
}

func NewProcessDependencySubTypeService(
	repo repositories.ProcessDependencySubTypeRepositoryInterface,
	processDependencyType ProcessDependencyTypeServiceInterface,
	activityLogSvc activityLogServices.ActivityLogServiceInterface,
) ProcessDependencySubTypeServiceInterface {
	return &ProcessDependencySubTypeService{
		repo:                  repo,
		processDependencyType: processDependencyType,
		activityLogSvc:        activityLogSvc,
	}
}

func (svc *ProcessDependencySubTypeService) GetAll(params *orgDTOs.ProcessDependencySubTypeListQueryParamsDTO) ([]*orgDTOs.ProcessDependencySubTypeDTO, error) {
	return svc.repo.GetAll(params)
}

func (svc *ProcessDependencySubTypeService) GetByID(id uint) (*models.ProcessDependencySubType, error) {
	subType, err := svc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return subType, nil
}

func (svc *ProcessDependencySubTypeService) Create(currentUserID uint, dto *orgDTOs.ProcessDependencySubTypeFormDTO) (uint, error) {
	if err := svc.validateForm(dto, nil); err != nil {
		return 0, err
	}

	entity := &models.ProcessDependencySubType{
		ProcessDependencyTypeID: dto.ProcessDependencyTypeID,
		TitleEn:                 dto.TitleEn,
		TitleAr:                 dto.TitleAr,
		DescriptionEn:           dto.DescriptionEn,
		DescriptionAr:           dto.DescriptionAr,
		Status:                  dto.Status,
	}

	id, err := svc.repo.Create(entity)
	if err != nil {
		return 0, err
	}

	go func() {
		created, _ := svc.GetByID(id)
		svc.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ProcessDependencySubType{}.TableName(),
			id,
			created,
			currentUserID,
			nil,
		)
	}()

	return id, nil
}

func (svc *ProcessDependencySubTypeService) Update(id uint, currentUserID uint, dto *orgDTOs.ProcessDependencySubTypeFormDTO) error {
	oldData, err := svc.GetByID(id)
	if err != nil {
		return err
	}

	existing, err := svc.GetByID(id)
	if err != nil {
		return err
	}

	if err := svc.validateForm(dto, &id); err != nil {
		return err
	}

	existing.ProcessDependencyTypeID = dto.ProcessDependencyTypeID
	existing.TitleEn = dto.TitleEn
	existing.TitleAr = dto.TitleAr
	existing.DescriptionEn = dto.DescriptionEn
	existing.DescriptionAr = dto.DescriptionAr
	existing.Status = dto.Status

	err = svc.repo.Update(existing)
	if err != nil {
		return err
	}

	changed, _ := utility.CompareFields(oldData, existing, []string{
		"ProcessDependencyTypeID", "TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "Status",
	})
	changedDataJSON, _ := utility.ChangesToJSON(changed)

	go func() {
		svc.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ProcessDependencySubType{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			nil,
		)
	}()

	return nil
}

func (svc *ProcessDependencySubTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedItems []map[string]interface{}

	for _, id := range ids {
		item, err := svc.repo.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		if err := svc.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
			deletedItems = append(deletedItems, map[string]interface{}{
				"id":       id,
				"title_en": item.TitleEn,
				"title_ar": item.TitleAr,
			})
		}
	}

	if len(successfulIds) > 0 {
		go svc.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.ProcessDependencySubType{}.TableName(),
			successfulIds,
			deletedItems,
			currentUserID,
			nil,
		)
	}

	return failedIds, nil
}

func (svc *ProcessDependencySubTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := svc.repo.Activate(ids)
	if len(successfulIds) > 0 {
		go svc.activityLogSvc.LogBulkActivities(
			models.ActivateAction,
			models.ProcessDependencySubType{}.TableName(),
			successfulIds,
			nil,
			currentUserID,
			nil,
		)
	}
	return err
}

func (svc *ProcessDependencySubTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := svc.repo.Deactivate(ids)
	if len(successfulIds) > 0 {
		go svc.activityLogSvc.LogBulkActivities(
			models.DeactivateAction,
			models.ProcessDependencySubType{}.TableName(),
			successfulIds,
			nil,
			currentUserID,
			nil,
		)
	}
	return err
}

func (svc *ProcessDependencySubTypeService) validateForm(dto *orgDTOs.ProcessDependencySubTypeFormDTO, id *uint) error {
	if _, err := svc.processDependencyType.GetByID(dto.ProcessDependencyTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("dependency type does not exist"), "not_exist", "process_dependency_type")
	}

	if id == nil {
		exists, err := svc.repo.ExistsByTitle(dto.TitleEn, dto.TitleAr)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	} else {
		exists, err := svc.repo.ExistsByTitleExcludingID(*id, dto.TitleEn, dto.TitleAr)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}
	return nil
}
