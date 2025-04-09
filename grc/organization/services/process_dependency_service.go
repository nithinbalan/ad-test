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

type ProcessDependencyServiceInterface interface {
	GetAll(params *orgDTOs.ProcessDependencyListQueryParamsDTO) ([]*orgDTOs.ProcessDependencyDTO, error)
	GetByID(id uint) (*models.ProcessDependency, error)
	Create(currentUserID uint, dto *orgDTOs.ProcessDependencyFormDTO) (uint, error)
	Update(id uint, currentUserID uint, dto *orgDTOs.ProcessDependencyFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type ProcessDependencyService struct {
	repo                     repositories.ProcessDependencyRepositoryInterface
	processDependencySubType ProcessDependencySubTypeServiceInterface
	activityLogSvc           activityLogServices.ActivityLogServiceInterface
}

func NewProcessDependencyService(
	repo repositories.ProcessDependencyRepositoryInterface,
	processDependencySubType ProcessDependencySubTypeServiceInterface,
	activityLogSvc activityLogServices.ActivityLogServiceInterface,
) ProcessDependencyServiceInterface {
	return &ProcessDependencyService{
		repo:                     repo,
		processDependencySubType: processDependencySubType,
		activityLogSvc:           activityLogSvc,
	}
}

func (svc *ProcessDependencyService) GetAll(params *orgDTOs.ProcessDependencyListQueryParamsDTO) ([]*orgDTOs.ProcessDependencyDTO, error) {
	return svc.repo.GetAll(params)
}

func (svc *ProcessDependencyService) GetByID(id uint) (*models.ProcessDependency, error) {
	processDependency, err := svc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return processDependency, nil
}

func (svc *ProcessDependencyService) Create(currentUserID uint, dto *orgDTOs.ProcessDependencyFormDTO) (uint, error) {
	if err := svc.validateForm(dto, nil); err != nil {
		return 0, err
	}

	entity := &models.ProcessDependency{
		ProcessID:                  *dto.DependentProcessID,
		ProcessDependencySubTypeID: dto.ProcessDependencySubTypeID,
		DependentProcessID:         dto.DependentProcessID,
		AssetID:                    dto.AssetID,
		OrganizationUnitID:         dto.OrganizationUnitID,
		Title:                      dto.Title,
		Description:                *dto.Description,
		CreatedBYID:                currentUserID,
	}

	id, err := svc.repo.Create(entity)
	if err != nil {
		return 0, err
	}

	go func() {
		created, _ := svc.GetByID(id)
		svc.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ProcessDependency{}.TableName(),
			id,
			created,
			currentUserID,
			nil,
		)
	}()

	return id, nil
}

func (svc *ProcessDependencyService) Update(id uint, currentUserID uint, dto *orgDTOs.ProcessDependencyFormDTO) error {
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

	existing.ProcessDependencySubTypeID = dto.ProcessDependencySubTypeID
	existing.DependentProcessID = dto.DependentProcessID
	existing.AssetID = dto.AssetID
	existing.OrganizationUnitID = dto.OrganizationUnitID
	existing.Title = dto.Title
	existing.Description = *dto.Description
	existing.UpdatedBYID = &currentUserID

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
			models.ProcessDependency{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			nil,
		)
	}()

	return nil
}

func (svc *ProcessDependencyService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
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
				"id":    id,
				"title": item.Title,
			})
		}
	}

	if len(successfulIds) > 0 {
		go svc.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.ProcessDependency{}.TableName(),
			successfulIds,
			deletedItems,
			currentUserID,
			nil,
		)
	}

	return failedIds, nil
}

func (svc *ProcessDependencyService) validateForm(dto *orgDTOs.ProcessDependencyFormDTO, id *uint) error {
	if _, err := svc.processDependencySubType.GetByID(dto.ProcessDependencySubTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("dependency sub type does not exist"), "not_exist", "process dependency sub type")
	}

	return nil
}
