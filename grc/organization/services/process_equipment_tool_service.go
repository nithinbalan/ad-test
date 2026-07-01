package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProcessEquipmentToolServiceInterface interface {
	GetAll(params *dtos.ProcessEquipmentToolListQueryParamsDTO) ([]*models.ProcessEquipmentTool, error)
	GetByID(id uint) (*models.ProcessEquipmentTool, error)
	Create(currentUserID uint, formDTO *dtos.ProcessEquipmentToolFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *dtos.ProcessEquipmentToolFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type ProcessEquipmentToolService struct {
	repo           repositories.ProcessEquipmentToolRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewProcessEquipmentToolService(repo repositories.ProcessEquipmentToolRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) ProcessEquipmentToolServiceInterface {
	return &ProcessEquipmentToolService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (pets *ProcessEquipmentToolService) GetAll(params *dtos.ProcessEquipmentToolListQueryParamsDTO) ([]*models.ProcessEquipmentTool, error) {
	return pets.repo.GetAll(params)
}

func (pets *ProcessEquipmentToolService) GetByID(id uint) (*models.ProcessEquipmentTool, error) {
	tool, err := pets.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return tool, nil
}

func (pets *ProcessEquipmentToolService) Create(currentUserID uint, formDTO *dtos.ProcessEquipmentToolFormDTO) (uint, error) {
	pets.applyDefaults(formDTO)

	if err := pets.validateFormDTO(formDTO, nil); err != nil {
		return 0, err
	}

	now := time.Now()
	tool := &models.ProcessEquipmentTool{
		TitleEn:     formDTO.TitleEn,
		TitleAr:     formDTO.TitleAr,
		Status:      formDTO.Status,
		CreatedBYID: &currentUserID,
		CreatedAt:   now,
	}

	id, err := pets.repo.Create(tool)
	if err != nil {
		return 0, err
	}

	go func() {
		created, err := pets.repo.GetByID(id)
		if err != nil {
			return
		}
		if err := pets.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ProcessEquipmentTool{}.TableName(),
			id,
			created,
			currentUserID,
			"", "", "", nil,
		); err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return id, nil
}

func (pets *ProcessEquipmentToolService) Update(id uint, currentUserID uint, formDTO *dtos.ProcessEquipmentToolFormDTO) error {
	existing, err := pets.GetByID(id)
	if err != nil {
		return err
	}
	old := *existing

	pets.applyDefaults(formDTO)

	if err := pets.validateFormDTO(formDTO, &id); err != nil {
		return err
	}

	now := time.Now()
	existing.TitleEn = formDTO.TitleEn
	existing.TitleAr = formDTO.TitleAr
	existing.Status = formDTO.Status
	existing.UpdatedBYID = &currentUserID
	existing.UpdatedAt = &now

	if err := pets.repo.Update(existing); err != nil {
		return err
	}

	changes, err := utility.CompareFields(&old, existing, []string{"TitleEn", "TitleAr", "Status"})
	if err != nil {
		log.Println("compare error:", err)
	}
	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("json error:", err)
	}

	go func() {
		if err := pets.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ProcessEquipmentTool{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		); err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return nil
}

func (pets *ProcessEquipmentToolService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		if _, err := pets.repo.GetByID(id); err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		if err := pets.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func(ids []uint) {
			if err := pets.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ProcessEquipmentTool{}.TableName(),
				ids,
				nil,
				currentUserID,
				nil,
			); err != nil {
				log.Printf("Failed to log activity: %v", err)
			}
		}(append([]uint(nil), successfulIds...))
	}

	return failedIds, nil
}

func (pets *ProcessEquipmentToolService) validateFormDTO(formDTO *dtos.ProcessEquipmentToolFormDTO, excludeID *uint) error {
	if excludeID != nil {
		exists, err := pets.repo.ExistsByTitleExcludingID(*excludeID, formDTO.TitleEn, formDTO.TitleAr)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(fmt.Errorf("process equipment tool with this title already exists"), "duplicate", "title")
		}
	} else {
		exists, err := pets.repo.ExistsByTitle(formDTO.TitleEn, formDTO.TitleAr)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(fmt.Errorf("process equipment tool with this title already exists"), "duplicate", "title")
		}
	}
	return nil
}

func (pets *ProcessEquipmentToolService) applyDefaults(formDTO *dtos.ProcessEquipmentToolFormDTO) {
	if formDTO.Status == "" {
		formDTO.Status = "active"
	}
}
