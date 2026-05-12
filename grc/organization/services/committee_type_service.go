package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type CommitteeTypeServiceInterface interface {
	GetAll(params *dtos.CommitteeTypeListQueryParamsDTO) ([]*models.CommitteeType, error)
	GetByID(id uint) (*models.CommitteeType, error)
	GetByType(typeFilter string) (*models.CommitteeType, error)
	Create(currentUserID uint, ctDTO *dtos.CommitteeTypeFormDTO) (uint, error)
	Update(id uint, currentUserID uint, ctDTO *dtos.CommitteeTypeFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type CommitteeTypeService struct {
	repo           repositories.CommitteeTypeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewCommitteeTypeService(repo repositories.CommitteeTypeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) CommitteeTypeServiceInterface {
	return &CommitteeTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (cts *CommitteeTypeService) GetAll(params *dtos.CommitteeTypeListQueryParamsDTO) ([]*models.CommitteeType, error) {
	return cts.repo.GetAll(params)
}

func (cts CommitteeTypeService) GetByID(id uint) (*models.CommitteeType, error) {
	committeeType, err := cts.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return committeeType, nil

}

func (cts *CommitteeTypeService) GetByType(typeFilter string) (*models.CommitteeType, error) {
	committeeType, err := cts.repo.GetByType(typeFilter)
	if err != nil {
		return nil, err
	}

	return committeeType, nil
}

func (cts *CommitteeTypeService) Create(currentUserID uint, ctDTO *dtos.CommitteeTypeFormDTO) (uint, error) {
	cts.applyDefaults(ctDTO)

	if err := cts.validateCommitteeTypeFormDTO(ctDTO, nil); err != nil {
		return 0, err
	}

	committeeType := &models.CommitteeType{
		TitleEn: ctDTO.TitleEn,
		TitleAr: ctDTO.TitleAr,
		Type:    ctDTO.Type,
		Status:  ctDTO.Status,
	}

	id, err := cts.repo.Create(committeeType)
	if err != nil {
		return 0, err
	}

	go func() {
		createdCommitteeType, err := cts.repo.GetByID(id)
		if err != nil {
			return
		}

		err = cts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.CommitteeType{}.TableName(),
			id,
			createdCommitteeType,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return id, nil
}

func (cts *CommitteeTypeService) Update(id uint, currentUserID uint, ctDTO *dtos.CommitteeTypeFormDTO) error {
	existingCT, err := cts.GetByID(id)
	if err != nil {
		return err
	}
	oldCT := *existingCT

	cts.applyDefaults(ctDTO)

	if err := cts.validateCommitteeTypeFormDTO(ctDTO, &id); err != nil {
		return err
	}

	existingCT.TitleEn = ctDTO.TitleEn
	existingCT.TitleAr = ctDTO.TitleAr
	existingCT.Type = ctDTO.Type
	existingCT.Status = ctDTO.Status

	err = cts.repo.Update(existingCT)
	if err != nil {
		return err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Type", "Status"}
	changes, err := utility.CompareFields(&oldCT, existingCT, fieldsToCompare)
	if err != nil {
		log.Println("compare error:", err)
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("json error:", err)
	}

	go func() {
		err := cts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.CommitteeType{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return nil
}

func (cts *CommitteeTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {

		_, err := cts.repo.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		err = cts.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func(ids []uint) {
			err := cts.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.CommitteeType{}.TableName(),
				ids,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				log.Printf("Failed to log activity: %v", err)
			}
		}(append([]uint(nil), successfulIds...))
	}

	return failedIds, nil
}

func (cts *CommitteeTypeService) validateCommitteeTypeFormDTO(ctDTO *dtos.CommitteeTypeFormDTO, excludeID *uint) error {
	if excludeID != nil {
		exists, err := cts.repo.ExistsByTitleExcludingID(*excludeID, ctDTO.TitleEn, ctDTO.TitleAr)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(fmt.Errorf("committee type with this title already exists"), "duplicate", "title")
		}
	} else {
		exists, err := cts.repo.ExistsByTitle(ctDTO.TitleEn, ctDTO.TitleAr)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(fmt.Errorf("committee type with this title already exists"), "duplicate", "title")
		}
	}

	return nil
}

func (cts *CommitteeTypeService) applyDefaults(ctDTO *dtos.CommitteeTypeFormDTO) {
	// Keep this in sync with CommitteeTypeFormDTO validation: empty values are treated as defaults.
	if ctDTO.Type == "" {
		ctDTO.Type = "other"
	}

	if ctDTO.Status == "" {
		ctDTO.Status = "active"
	}
}
