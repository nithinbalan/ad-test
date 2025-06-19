package services

import (
	"errors"
	"fmt"
	logservice "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type CompetencyServiceInterface interface {
	GetAll(params *dtos.CompetencyListQueryParamsDTO) ([]*dtos.CompetencyWithGroupDTO, error)
	GetByID(id uint) (*models.Competency, error)
	Create(currentUserID uint, competency *dtos.CompetencyFormDTO) (*models.Competency, error)
	Update(id uint, currentUserID uint, competency *models.Competency) (*models.Competency, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type CompetencyService struct {
	repo           repositories.CompetencyRepositoryInterface
	activityLogSvc logservice.ActivityLogServiceInterface
}

func NewCompetencyService(repo repositories.CompetencyRepositoryInterface, activityLogSvc logservice.ActivityLogServiceInterface) CompetencyServiceInterface {
	return &CompetencyService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (cs *CompetencyService) GetAll(params *dtos.CompetencyListQueryParamsDTO) ([]*dtos.CompetencyWithGroupDTO, error) {
	return cs.repo.GetAll(params)
}

func (cs *CompetencyService) GetByID(id uint) (*models.Competency, error) {
	competency, err := cs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return competency, nil
}

func (cs *CompetencyService) Create(currentUserID uint, dto *dtos.CompetencyFormDTO) (*models.Competency, error) {
	if exists, err := cs.repo.ExistsByTitle(dto.TitleEn, dto.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	competency := &models.Competency{
		CompetencyGroupID: dto.CompetencyGroupID,
		TitleEn:           dto.TitleEn,
		TitleAr:           dto.TitleAr,
		Status:            dto.Status,
	}

	createdCompetency, err := cs.repo.Create(competency)
	if err != nil {
		return nil, err
	}

	go func() {
		err := cs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Competency{}.TableName(),
			createdCompetency.ID,
			createdCompetency,
			currentUserID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdCompetency, nil
}

func (cs *CompetencyService) Update(id uint, currentUserID uint, updatedCompetency *models.Competency) (*models.Competency, error) {
	existing, err := cs.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedCompetency.TitleEn != existing.TitleEn || updatedCompetency.TitleAr != existing.TitleAr {
		if exists, err := cs.repo.ExistsByTitleExcludingID(id, updatedCompetency.TitleEn, updatedCompetency.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := cs.repo.Update(id, updatedCompetency)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status", "CompetencyGroupID"}
	changes, err := utility.CompareFields(existing, updatedCompetency, fieldsToCompare)
	if err != nil {
		log.Println("error comparing fields:", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error generating JSON:", err.Error())
	}

	go func() {
		err := cs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.Competency{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to log update activity: %v", err)
		}
	}()

	return result, nil
}

func (cs *CompetencyService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds, successfulIds []uint
	for _, id := range ids {
		_, err := cs.repo.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if err := cs.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := cs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.Competency{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				log.Printf("Failed to log bulk delete activity: %v", err)
			}
		}()
	}

	return failedIds, nil
}

func (cs *CompetencyService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := cs.repo.Activate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := cs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.Competency{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				log.Printf("Failed to log activation activity: %v", err)
			}
		}()
	}
	return err
}

func (cs *CompetencyService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := cs.repo.Deactivate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := cs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.Competency{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				log.Printf("Failed to log deactivation activity: %v", err)
			}
		}()
	}
	return err
}
