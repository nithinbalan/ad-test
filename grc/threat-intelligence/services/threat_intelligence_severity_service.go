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

type ThreatIntelligenceSeverityServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*tiDTOs.ThreatIntelligenceSeverityDTO, error)
	GetByID(id uint) (*models.ThreatIntelligenceSeverity, error)
	Create(currentUserID uint, severityFormDTO *tiDTOs.ThreatIntelligenceSeverityFormDTO) (uint, error)
	Update(id uint, currentUserID uint, severityFormDTO *tiDTOs.ThreatIntelligenceSeverityFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type ThreatIntelligenceSeverityService struct {
	repo           repositories.ThreatIntelligenceSeverityRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
}

func NewThreatIntelligenceSeverityService(repo repositories.ThreatIntelligenceSeverityRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface) ThreatIntelligenceSeverityServiceInterface {
	return &ThreatIntelligenceSeverityService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (tiss *ThreatIntelligenceSeverityService) GetAll(sort string, sortBy string, keyword string, status string) ([]*tiDTOs.ThreatIntelligenceSeverityDTO, error) {
	return tiss.repo.GetAll(sort, sortBy, keyword, status)
}

func (tiss *ThreatIntelligenceSeverityService) GetByID(id uint) (*models.ThreatIntelligenceSeverity, error) {
	severity, err := tiss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return severity, nil
}

func (tiss *ThreatIntelligenceSeverityService) Create(currentUserID uint, severityFormDTO *tiDTOs.ThreatIntelligenceSeverityFormDTO) (uint, error) {
	if err := tiss.validateThreatIntelligenceSeverityFormDTO(severityFormDTO, nil); err != nil {
		return 0, err
	}

	severity := &models.ThreatIntelligenceSeverity{
		TitleEn: severityFormDTO.TitleEn,
		TitleAr: severityFormDTO.TitleAr,
		Score:   severityFormDTO.Score,
		Status:  severityFormDTO.Status,
	}

	id, err := tiss.repo.Create(severity)
	if err != nil {
		return 0, err
	}

	go func() {
		existingSeverity, err := tiss.GetByID(id)
		if err != nil {
			utility.LogError("Failed to fetch threat intelligence severity", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = tiss.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ThreatIntelligenceSeverity{}.TableName(),
			id,
			existingSeverity,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ThreatIntelligenceSeverity{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (tiss *ThreatIntelligenceSeverityService) Update(id uint, currentUserID uint, severityFormDTO *tiDTOs.ThreatIntelligenceSeverityFormDTO) error {
	existingSeverity, err := tiss.GetByID(id)
	if err != nil {
		return err
	}

	oldSeverity, err := tiss.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := tiss.validateThreatIntelligenceSeverityFormDTO(severityFormDTO, &id); err != nil {
		return err
	}

	existingSeverity.TitleEn = severityFormDTO.TitleEn
	existingSeverity.TitleAr = severityFormDTO.TitleAr
	existingSeverity.Score = severityFormDTO.Score
	existingSeverity.Status = severityFormDTO.Status

	err = tiss.repo.Update(existingSeverity)
	if err != nil {
		return err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Score", "Status"}
	changes, err := utility.CompareFields(oldSeverity, existingSeverity, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := tiss.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ThreatIntelligenceSeverity{}.TableName(),
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

func (tiss *ThreatIntelligenceSeverityService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedSeverities []map[string]interface{}

	for _, id := range ids {
		severity, err := tiss.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = tiss.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedSeverities = append(deletedSeverities, map[string]interface{}{
			"id":       id,
			"title_en": severity.TitleEn,
			"title_ar": severity.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := tiss.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ThreatIntelligenceSeverity{}.TableName(),
				successfulIds,
				deletedSeverities,
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

func (tiss *ThreatIntelligenceSeverityService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := tiss.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := tiss.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ThreatIntelligenceSeverity{}.TableName(),
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

func (tiss *ThreatIntelligenceSeverityService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := tiss.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := tiss.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ThreatIntelligenceSeverity{}.TableName(),
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

func (tiss *ThreatIntelligenceSeverityService) validateThreatIntelligenceSeverityFormDTO(severityFormDTO *tiDTOs.ThreatIntelligenceSeverityFormDTO, id *uint) error {
	if id == nil {
		if exists, err := tiss.repo.ExistsByTitle(severityFormDTO.TitleEn, severityFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	} else {
		if exists, err := tiss.repo.ExistsByTitleExcludingID(*id, severityFormDTO.TitleEn, severityFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}
