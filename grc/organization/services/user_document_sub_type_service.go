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

type UserDocumentSubTypeServiceInterface interface {
	GetAll(params *orgDTOs.UserDocumentSubTypeListQueryParamsDTO) ([]*orgDTOs.UserDocumentSubTypeDTO, error)
	GetByID(id uint) (*models.UserDocumentSubType, error)
	Create(currentUserID uint, dto *orgDTOs.UserDocumentSubTypeFormDTO) (uint, error)
	Update(id uint, currentUserID uint, dto *orgDTOs.UserDocumentSubTypeFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type UserDocumentSubTypeService struct {
	repo             repositories.UserDocumentSubTypeRepositoryInterface
	userDocumentType UserDocumentTypeServiceInterface
	activityLogSvc   activityLogServices.ActivityLogServiceInterface
}

func NewUserDocumentSubTypeService(
	repo repositories.UserDocumentSubTypeRepositoryInterface,
	userDocumentType UserDocumentTypeServiceInterface,
	activityLogSvc activityLogServices.ActivityLogServiceInterface,
) UserDocumentSubTypeServiceInterface {
	return &UserDocumentSubTypeService{
		repo:             repo,
		userDocumentType: userDocumentType,
		activityLogSvc:   activityLogSvc,
	}
}

func (udsts *UserDocumentSubTypeService) GetAll(params *orgDTOs.UserDocumentSubTypeListQueryParamsDTO) ([]*orgDTOs.UserDocumentSubTypeDTO, error) {
	return udsts.repo.GetAll(params)
}

func (udsts *UserDocumentSubTypeService) GetByID(id uint) (*models.UserDocumentSubType, error) {
	subType, err := udsts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return subType, nil
}

func (udsts *UserDocumentSubTypeService) Create(currentUserID uint, dto *orgDTOs.UserDocumentSubTypeFormDTO) (uint, error) {
	if err := udsts.validateForm(dto, nil); err != nil {
		return 0, err
	}

	entity := &models.UserDocumentSubType{
		UserDocumentTypeID: dto.UserDocumentTypeID,
		TitleEn:            dto.TitleEn,
		TitleAr:            dto.TitleAr,
		Status:             dto.Status,
	}

	id, err := udsts.repo.Create(entity)
	if err != nil {
		return 0, err
	}

	go func() {
		created, _ := udsts.GetByID(id)
		udsts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.UserDocumentSubType{}.TableName(),
			id,
			created,
			currentUserID,
			nil,
		)
	}()

	return id, nil
}

func (udsts *UserDocumentSubTypeService) Update(id uint, currentUserID uint, dto *orgDTOs.UserDocumentSubTypeFormDTO) error {
	oldData, err := udsts.GetByID(id)
	if err != nil {
		return err
	}

	existing, err := udsts.GetByID(id)
	if err != nil {
		return err
	}

	if err := udsts.validateForm(dto, &id); err != nil {
		return err
	}

	existing.UserDocumentTypeID = dto.UserDocumentTypeID
	existing.TitleEn = dto.TitleEn
	existing.TitleAr = dto.TitleAr
	existing.Status = dto.Status

	err = udsts.repo.Update(existing)
	if err != nil {
		return err
	}

	changed, _ := utility.CompareFields(oldData, existing, []string{
		"UserDocumentTypeID", "TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "Status",
	})
	changedDataJSON, _ := utility.ChangesToJSON(changed)

	go func() {
		udsts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.UserDocumentSubType{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			nil,
		)
	}()

	return nil
}

func (udsts *UserDocumentSubTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedItems []map[string]interface{}

	for _, id := range ids {
		item, err := udsts.repo.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		if err := udsts.repo.Delete(id); err != nil {
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
		go udsts.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.UserDocumentSubType{}.TableName(),
			successfulIds,
			deletedItems,
			currentUserID,
			nil,
		)
	}

	return failedIds, nil
}

func (udsts *UserDocumentSubTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := udsts.repo.Activate(ids)
	if len(successfulIds) > 0 {
		go udsts.activityLogSvc.LogBulkActivities(
			models.ActivateAction,
			models.UserDocumentSubType{}.TableName(),
			successfulIds,
			nil,
			currentUserID,
			nil,
		)
	}
	return err
}

func (udsts *UserDocumentSubTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := udsts.repo.Deactivate(ids)
	if len(successfulIds) > 0 {
		go udsts.activityLogSvc.LogBulkActivities(
			models.DeactivateAction,
			models.UserDocumentSubType{}.TableName(),
			successfulIds,
			nil,
			currentUserID,
			nil,
		)
	}
	return err
}

func (udsts *UserDocumentSubTypeService) validateForm(dto *orgDTOs.UserDocumentSubTypeFormDTO, id *uint) error {
	if _, err := udsts.userDocumentType.GetByID(dto.UserDocumentTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("user document type does not exist"), "not_exist", "user_document_type")
	}

	if id == nil {
		exists, err := udsts.repo.ExistsByTitle(dto.TitleEn, dto.TitleAr)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	} else {
		exists, err := udsts.repo.ExistsByTitleExcludingID(*id, dto.TitleEn, dto.TitleAr)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}
	return nil
}
