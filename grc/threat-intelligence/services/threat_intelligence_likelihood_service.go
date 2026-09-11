package services

import (
	"errors"
	"fmt"
	"log"

	activityLogServices "jethur/grc/logs/services"
	"jethur/grc/models"
	tiDTOs "jethur/grc/threat-intelligence/dtos"
	"jethur/grc/threat-intelligence/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type ThreatIntelligenceLikelihoodServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*tiDTOs.ThreatIntelligenceLikelihoodDTO, error)
	GetByID(id uint) (*models.ThreatIntelligenceLikelihood, error)
	Create(currentUserID uint, likelihoodFormDTO *tiDTOs.ThreatIntelligenceLikelihoodFormDTO) (uint, error)
	Update(id uint, currentUserID uint, likelihoodFormDTO *tiDTOs.ThreatIntelligenceLikelihoodFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type ThreatIntelligenceLikelihoodService struct {
	repo           repositories.ThreatIntelligenceLikelihoodRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
}

func NewThreatIntelligenceLikelihoodService(repo repositories.ThreatIntelligenceLikelihoodRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface) ThreatIntelligenceLikelihoodServiceInterface {
	return &ThreatIntelligenceLikelihoodService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (tils *ThreatIntelligenceLikelihoodService) GetAll(sort string, sortBy string, keyword string, status string) ([]*tiDTOs.ThreatIntelligenceLikelihoodDTO, error) {
	return tils.repo.GetAll(sort, sortBy, keyword, status)
}

func (tils *ThreatIntelligenceLikelihoodService) GetByID(id uint) (*models.ThreatIntelligenceLikelihood, error) {
	likelihood, err := tils.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return likelihood, nil
}

func (tils *ThreatIntelligenceLikelihoodService) Create(currentUserID uint, likelihoodFormDTO *tiDTOs.ThreatIntelligenceLikelihoodFormDTO) (uint, error) {
	if err := tils.validateThreatIntelligenceLikelihoodFormDTO(likelihoodFormDTO, nil); err != nil {
		return 0, err
	}

	likelihood := &models.ThreatIntelligenceLikelihood{
		TitleEn: likelihoodFormDTO.TitleEn,
		TitleAr: likelihoodFormDTO.TitleAr,
		Score:   likelihoodFormDTO.Score,
		Status:  likelihoodFormDTO.Status,
	}

	id, err := tils.repo.Create(likelihood)
	if err != nil {
		return 0, err
	}

	go func() {
		existingLikelihood, err := tils.GetByID(id)
		if err != nil {
			utility.LogError("Failed to fetch threat intelligence likelihood", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = tils.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ThreatIntelligenceLikelihood{}.TableName(),
			id,
			existingLikelihood,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ThreatIntelligenceLikelihood{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (tils *ThreatIntelligenceLikelihoodService) Update(id uint, currentUserID uint, likelihoodFormDTO *tiDTOs.ThreatIntelligenceLikelihoodFormDTO) error {
	existingLikelihood, err := tils.GetByID(id)
	if err != nil {
		return err
	}

	oldLikelihood, err := tils.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := tils.validateThreatIntelligenceLikelihoodFormDTO(likelihoodFormDTO, &id); err != nil {
		return err
	}

	existingLikelihood.TitleEn = likelihoodFormDTO.TitleEn
	existingLikelihood.TitleAr = likelihoodFormDTO.TitleAr
	existingLikelihood.Score = likelihoodFormDTO.Score
	existingLikelihood.Status = likelihoodFormDTO.Status

	err = tils.repo.Update(existingLikelihood)
	if err != nil {
		return err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Score", "Status"}
	changes, err := utility.CompareFields(oldLikelihood, existingLikelihood, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := tils.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ThreatIntelligenceLikelihood{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()

	return nil
}

func (tils *ThreatIntelligenceLikelihoodService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedLikelihoods []map[string]interface{}

	for _, id := range ids {
		likelihood, err := tils.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = tils.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedLikelihoods = append(deletedLikelihoods, map[string]interface{}{
			"id":       id,
			"title_en": likelihood.TitleEn,
			"title_ar": likelihood.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := tils.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ThreatIntelligenceLikelihood{}.TableName(),
				successfulIds,
				deletedLikelihoods,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (tils *ThreatIntelligenceLikelihoodService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := tils.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := tils.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ThreatIntelligenceLikelihood{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk activate activity: %v\n", err)
			}
		}()
	}

	return err
}

func (tils *ThreatIntelligenceLikelihoodService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := tils.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := tils.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ThreatIntelligenceLikelihood{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk deactivate activity: %v\n", err)
			}
		}()
	}

	return err
}

func (tils *ThreatIntelligenceLikelihoodService) validateThreatIntelligenceLikelihoodFormDTO(likelihoodFormDTO *tiDTOs.ThreatIntelligenceLikelihoodFormDTO, id *uint) error {
	if id == nil {
		if exists, err := tils.repo.ExistsByTitle(likelihoodFormDTO.TitleEn, likelihoodFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	} else {
		if exists, err := tils.repo.ExistsByTitleExcludingID(*id, likelihoodFormDTO.TitleEn, likelihoodFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}
