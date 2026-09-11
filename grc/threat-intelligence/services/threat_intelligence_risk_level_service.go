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

type ThreatIntelligenceRiskLevelServiceInterface interface {
	GetAll(params *tiDTOs.ThreatIntelligenceRiskLevelListQueryParamsDTO) ([]*tiDTOs.ThreatIntelligenceRiskLevelDTO, error)
	GetByID(id uint) (*models.ThreatIntelligenceRiskLevel, error)
	Create(currentUserID uint, riskLevelFormDTO *tiDTOs.ThreatIntelligenceRiskLevelFormDTO) (uint, error)
	Update(id uint, currentUserID uint, riskLevelFormDTO *tiDTOs.ThreatIntelligenceRiskLevelFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type ThreatIntelligenceRiskLevelService struct {
	repo           repositories.ThreatIntelligenceRiskLevelRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
}

func NewThreatIntelligenceRiskLevelService(repo repositories.ThreatIntelligenceRiskLevelRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface) ThreatIntelligenceRiskLevelServiceInterface {
	return &ThreatIntelligenceRiskLevelService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (tirls *ThreatIntelligenceRiskLevelService) GetAll(params *tiDTOs.ThreatIntelligenceRiskLevelListQueryParamsDTO) ([]*tiDTOs.ThreatIntelligenceRiskLevelDTO, error) {
	return tirls.repo.GetAll(params)
}

func (tirls *ThreatIntelligenceRiskLevelService) GetByID(id uint) (*models.ThreatIntelligenceRiskLevel, error) {
	riskLevel, err := tirls.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskLevel, nil
}

func (tirls *ThreatIntelligenceRiskLevelService) Create(currentUserID uint, riskLevelFormDTO *tiDTOs.ThreatIntelligenceRiskLevelFormDTO) (uint, error) {
	if err := tirls.validateThreatIntelligenceRiskLevelFormDTO(riskLevelFormDTO, nil); err != nil {
		return 0, err
	}

	riskLevel := &models.ThreatIntelligenceRiskLevel{
		TitleEn:   riskLevelFormDTO.TitleEn,
		TitleAr:   riskLevelFormDTO.TitleAr,
		Type:      riskLevelFormDTO.Type,
		ColorCode: riskLevelFormDTO.ColorCode,
		Status:    riskLevelFormDTO.Status,
	}

	id, err := tirls.repo.Create(riskLevel)
	if err != nil {
		return 0, err
	}

	go func() {
		existingRiskLevel, err := tirls.GetByID(id)
		if err != nil {
			utility.LogError("Failed to fetch threat intelligence risk level", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = tirls.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ThreatIntelligenceRiskLevel{}.TableName(),
			id,
			existingRiskLevel,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ThreatIntelligenceRiskLevel{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (tirls *ThreatIntelligenceRiskLevelService) Update(id uint, currentUserID uint, riskLevelFormDTO *tiDTOs.ThreatIntelligenceRiskLevelFormDTO) error {
	existingRiskLevel, err := tirls.GetByID(id)
	if err != nil {
		return err
	}

	oldRiskLevel, err := tirls.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := tirls.validateThreatIntelligenceRiskLevelFormDTO(riskLevelFormDTO, &id); err != nil {
		return err
	}

	existingRiskLevel.TitleEn = riskLevelFormDTO.TitleEn
	existingRiskLevel.TitleAr = riskLevelFormDTO.TitleAr
	existingRiskLevel.Type = riskLevelFormDTO.Type
	existingRiskLevel.ColorCode = riskLevelFormDTO.ColorCode
	existingRiskLevel.Status = riskLevelFormDTO.Status

	err = tirls.repo.Update(existingRiskLevel)
	if err != nil {
		return err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Type", "ColorCode", "Status"}
	changes, err := utility.CompareFields(oldRiskLevel, existingRiskLevel, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := tirls.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ThreatIntelligenceRiskLevel{}.TableName(),
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

func (tirls *ThreatIntelligenceRiskLevelService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedRiskLevels []map[string]interface{}

	for _, id := range ids {
		riskLevel, err := tirls.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = tirls.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedRiskLevels = append(deletedRiskLevels, map[string]interface{}{
			"id":       id,
			"title_en": riskLevel.TitleEn,
			"title_ar": riskLevel.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := tirls.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ThreatIntelligenceRiskLevel{}.TableName(),
				successfulIds,
				deletedRiskLevels,
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

func (tirls *ThreatIntelligenceRiskLevelService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := tirls.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := tirls.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ThreatIntelligenceRiskLevel{}.TableName(),
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

func (tirls *ThreatIntelligenceRiskLevelService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := tirls.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := tirls.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ThreatIntelligenceRiskLevel{}.TableName(),
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

func (tirls *ThreatIntelligenceRiskLevelService) validateThreatIntelligenceRiskLevelFormDTO(riskLevelFormDTO *tiDTOs.ThreatIntelligenceRiskLevelFormDTO, id *uint) error {
	if id == nil {
		if exists, err := tirls.repo.ExistsByTitle(riskLevelFormDTO.TitleEn, riskLevelFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

		if exists, err := tirls.repo.ExistsByType(riskLevelFormDTO.Type); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the type must be unique"), "unique", "type")
		}
	} else {
		if exists, err := tirls.repo.ExistsByTitleExcludingID(*id, riskLevelFormDTO.TitleEn, riskLevelFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

		if exists, err := tirls.repo.ExistsByTypeExcludingID(*id, riskLevelFormDTO.Type); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the type must be unique"), "unique", "type")
		}
	}

	return nil
}
