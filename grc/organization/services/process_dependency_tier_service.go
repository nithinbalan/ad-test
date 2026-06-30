package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type ProcessDependencyTierServiceInterface interface {
	GetAll(params *dtos.ProcessDependencyTierListQueryParamsDTO) ([]*models.ProcessDependencyTier, error)
	GetByID(id uint) (*models.ProcessDependencyTier, error)
	Create(currentUserID uint, formDTO *dtos.ProcessDependencyTierFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *dtos.ProcessDependencyTierFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type ProcessDependencyTierService struct {
	repo           repositories.ProcessDependencyTierRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewProcessDependencyTierService(repo repositories.ProcessDependencyTierRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) ProcessDependencyTierServiceInterface {
	return &ProcessDependencyTierService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (s *ProcessDependencyTierService) GetAll(params *dtos.ProcessDependencyTierListQueryParamsDTO) ([]*models.ProcessDependencyTier, error) {
	return s.repo.GetAll(params)
}

func (s *ProcessDependencyTierService) GetByID(id uint) (*models.ProcessDependencyTier, error) {
	tier, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return tier, nil
}

func (s *ProcessDependencyTierService) Create(currentUserID uint, formDTO *dtos.ProcessDependencyTierFormDTO) (uint, error) {
	s.applyDefaults(formDTO)

	if err := s.validateFormDTO(formDTO, nil); err != nil {
		return 0, err
	}

	tier := &models.ProcessDependencyTier{
		TitleEn:     formDTO.TitleEn,
		TitleAr:     formDTO.TitleAr,
		Status:      formDTO.Status,
		CreatedBYID: &currentUserID,
	}

	id, err := s.repo.Create(tier)
	if err != nil {
		return 0, err
	}

	go func() {
		created, err := s.repo.GetByID(id)
		if err != nil {
			return
		}

		err = s.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ProcessDependencyTier{}.TableName(),
			id,
			created,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return id, nil
}

func (s *ProcessDependencyTierService) Update(id uint, currentUserID uint, formDTO *dtos.ProcessDependencyTierFormDTO) error {
	existing, err := s.GetByID(id)
	if err != nil {
		return err
	}
	old := *existing

	s.applyDefaults(formDTO)

	if err := s.validateFormDTO(formDTO, &id); err != nil {
		return err
	}

	existing.TitleEn = formDTO.TitleEn
	existing.TitleAr = formDTO.TitleAr
	existing.Status = formDTO.Status
	existing.UpdatedBYID = &currentUserID

	if err := s.repo.Update(existing); err != nil {
		return err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(&old, existing, fieldsToCompare)
	if err != nil {
		log.Println("compare error:", err)
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("json error:", err)
	}

	go func() {
		err := s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ProcessDependencyTier{}.TableName(),
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

func (s *ProcessDependencyTierService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		if _, err := s.repo.GetByID(id); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if err := s.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func(ids []uint) {
			err := s.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ProcessDependencyTier{}.TableName(),
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

func (s *ProcessDependencyTierService) validateFormDTO(formDTO *dtos.ProcessDependencyTierFormDTO, excludeID *uint) error {
	if excludeID != nil {
		exists, err := s.repo.ExistsByTitleExcludingID(*excludeID, formDTO.TitleEn, formDTO.TitleAr)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(fmt.Errorf("process dependency tier with this title already exists"), "duplicate", "title")
		}
	} else {
		exists, err := s.repo.ExistsByTitle(formDTO.TitleEn, formDTO.TitleAr)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(fmt.Errorf("process dependency tier with this title already exists"), "duplicate", "title")
		}
	}

	return nil
}

func (s *ProcessDependencyTierService) applyDefaults(formDTO *dtos.ProcessDependencyTierFormDTO) {
	if formDTO.Status == "" {
		formDTO.Status = "active"
	}
}
