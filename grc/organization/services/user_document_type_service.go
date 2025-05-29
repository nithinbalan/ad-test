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

type UserDocumentTypeServiceInterface interface {
	GetAll(params *orgDTOs.UserDocumentTypeListQueryParamsDTO) ([]*orgDTOs.UserDocumentTypeDTO, error)
	GetByID(id uint) (*models.UserDocumentType, error)
	Create(currentUserID uint, dto *orgDTOs.UserDocumentTypeFormDTO) (uint, error)
	Update(id uint, currentUserID uint, dto *orgDTOs.UserDocumentTypeFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type UserDocumentTypeService struct {
	repo           repositories.UserDocumentTypeRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
}

func NewUserDocumentTypeService(
	repo repositories.UserDocumentTypeRepositoryInterface,
	activityLogSvc activityLogServices.ActivityLogServiceInterface,
) UserDocumentTypeServiceInterface {
	return &UserDocumentTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (udts *UserDocumentTypeService) GetAll(params *orgDTOs.UserDocumentTypeListQueryParamsDTO) ([]*orgDTOs.UserDocumentTypeDTO, error) {
	return udts.repo.GetAll(params)
}

func (udts *UserDocumentTypeService) GetByID(id uint) (*models.UserDocumentType, error) {
	subType, err := udts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return subType, nil
}

func (udts *UserDocumentTypeService) Create(currentUserID uint, dto *orgDTOs.UserDocumentTypeFormDTO) (uint, error) {
	if err := udts.validateForm(dto, nil); err != nil {
		return 0, err
	}

	entity := &models.UserDocumentType{
		TitleEn: dto.TitleEn,
		TitleAr: dto.TitleAr,
		Status:  dto.Status,
	}

	id, err := udts.repo.Create(entity)
	if err != nil {
		return 0, err
	}

	go func() {
		created, _ := udts.GetByID(id)
		udts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.UserDocumentType{}.TableName(),
			id,
			created,
			currentUserID,
			nil,
		)
	}()

	return id, nil
}

func (udts *UserDocumentTypeService) Update(id uint, currentUserID uint, dto *orgDTOs.UserDocumentTypeFormDTO) error {
	oldData, err := udts.GetByID(id)
	if err != nil {
		return err
	}

	existing, err := udts.GetByID(id)
	if err != nil {
		return err
	}

	if err := udts.validateForm(dto, &id); err != nil {
		return err
	}

	existing.TitleEn = dto.TitleEn
	existing.TitleAr = dto.TitleAr
	existing.Status = dto.Status

	err = udts.repo.Update(existing)
	if err != nil {
		return err
	}

	changed, _ := utility.CompareFields(oldData, existing, []string{
		"TitleEn", "TitleAr", "Status",
	})
	changedDataJSON, _ := utility.ChangesToJSON(changed)

	go func() {
		udts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.UserDocumentType{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			nil,
		)
	}()

	return nil
}

func (udts *UserDocumentTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedItems []map[string]interface{}

	for _, id := range ids {
		item, err := udts.repo.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		if err := udts.repo.Delete(id); err != nil {
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
		go udts.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.UserDocumentType{}.TableName(),
			successfulIds,
			deletedItems,
			currentUserID,
			nil,
		)
	}

	return failedIds, nil
}

func (udts *UserDocumentTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := udts.repo.Activate(ids)
	if len(successfulIds) > 0 {
		go udts.activityLogSvc.LogBulkActivities(
			models.ActivateAction,
			models.UserDocumentType{}.TableName(),
			successfulIds,
			nil,
			currentUserID,
			nil,
		)
	}
	return err
}

func (udts *UserDocumentTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := udts.repo.Deactivate(ids)
	if len(successfulIds) > 0 {
		go udts.activityLogSvc.LogBulkActivities(
			models.DeactivateAction,
			models.UserDocumentType{}.TableName(),
			successfulIds,
			nil,
			currentUserID,
			nil,
		)
	}
	return err
}

func (udts *UserDocumentTypeService) validateForm(dto *orgDTOs.UserDocumentTypeFormDTO, id *uint) error {

	if id == nil {
		exists, err := udts.repo.ExistsByTitle(dto.TitleEn, dto.TitleAr)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	} else {
		exists, err := udts.repo.ExistsByTitleExcludingID(*id, dto.TitleEn, dto.TitleAr)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}
	return nil
}
