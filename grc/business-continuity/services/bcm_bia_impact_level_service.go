package services

import (
	"errors"
	"fmt"
	"log"
	"strings"

	bcmDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/repositories"
	logsServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BCMBIAImpactLevelServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*bcmDTOs.BCMBIAImpactLevelDTO, error)
	GetByID(id uint) (*models.BCMBIAImpactLevel, error)
	GetFullDetailsByID(id uint) (*bcmDTOs.BCMBIAImpactLevelDetailDTO, error)
	GetByImpactAnalysisLevelID(impactAnalysisLevelID uint) (*models.BCMBIAImpactLevel, error)
	GetByOutageBucketID(outageBucketID uint) (*models.BCMBIAImpactLevel, error)
	UpdateImpactLevels(currentUserID uint, req *bcmDTOs.BCMBIAImpactLevelBulkUpdateDTO) ([]*models.BCMBIAImpactLevel, error)
}

// Service and implementations
type BCMBIAImpactLevelService struct {
	repo           repositories.BCMBIAImpactLevelRepositoryInterface
	activityLogSvc logsServices.ActivityLogServiceInterface
}

func NewBCMBIAImpactLevelService(repo repositories.BCMBIAImpactLevelRepositoryInterface, activityLogSvc logsServices.ActivityLogServiceInterface) BCMBIAImpactLevelServiceInterface {
	return &BCMBIAImpactLevelService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (bbils *BCMBIAImpactLevelService) GetAll(sort string, sortBy string, keyword string, status string) ([]*bcmDTOs.BCMBIAImpactLevelDTO, error) {
	return bbils.repo.GetAll(sort, sortBy, keyword, status)
}

func (bbils *BCMBIAImpactLevelService) GetByID(id uint) (*models.BCMBIAImpactLevel, error) {
	bcmBiaImpactLevel, err := bbils.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return bcmBiaImpactLevel, nil
}

func (bbils *BCMBIAImpactLevelService) GetFullDetailsByID(id uint) (*bcmDTOs.BCMBIAImpactLevelDetailDTO, error) {
	var detailDTO bcmDTOs.BCMBIAImpactLevelDetailDTO

	bcmBiaImpactLevel, err := bbils.repo.GetFullDetailsByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if err := mapstructure.Decode(bcmBiaImpactLevel, &detailDTO); err != nil {
		return nil, err
	}

	return &detailDTO, nil
}

func (bbils *BCMBIAImpactLevelService) GetByImpactAnalysisLevelID(impactAnalysisLevelID uint) (*models.BCMBIAImpactLevel, error) {
	bcmBiaImpactLevel, err := bbils.repo.GetByImpactAnalysisLevelID(impactAnalysisLevelID)

	if err != nil {
		return nil, err
	}

	return bcmBiaImpactLevel, nil
}

func (bbils *BCMBIAImpactLevelService) GetByOutageBucketID(outageBucketID uint) (*models.BCMBIAImpactLevel, error) {
	bcmBiaImpactLevel, err := bbils.repo.GetByOutageBucketID(outageBucketID)

	if err != nil {
		return nil, err
	}

	return bcmBiaImpactLevel, nil
}

func (bbils *BCMBIAImpactLevelService) UpdateImpactLevels(currentUserID uint, req *bcmDTOs.BCMBIAImpactLevelBulkUpdateDTO) ([]*models.BCMBIAImpactLevel, error) {
	fieldsToCompare := []string{
		"ImpactAnalysisLevelID", "ClassificationEn", "ClassificationAr", "IsCritical",
		"RtoFromValue", "RtoToValue", "RtoCondition", "RtoUnitID",
		"MtpdFromValue", "MtpdToValue", "MtpdCondition", "MtpdUnitID",
		"RpoFromValue", "RpoToValue", "RpoCondition", "RpoUnitID",
		"RspFromValue", "RspToValue", "RspCondition", "RspUnitID",
	}

	var updatedLevels []*models.BCMBIAImpactLevel
	var updatedIDs []uint
	var allChanges []utility.ChangeDetail

	for _, form := range req.Settings {
		existingLevel, err := bbils.repo.GetByID(form.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, utility.NewResourceNotFoundError(err, form.ID)
			}
			return nil, err
		}

		if form.RtoCondition, err = normalizeBCMBIAImpactLevelCondition(form.RtoCondition, "rto_condition"); err != nil {
			return nil, err
		}
		if form.MtpdCondition, err = normalizeBCMBIAImpactLevelCondition(form.MtpdCondition, "mtpd_condition"); err != nil {
			return nil, err
		}
		if form.RpoCondition, err = normalizeBCMBIAImpactLevelCondition(form.RpoCondition, "rpo_condition"); err != nil {
			return nil, err
		}
		if form.RspCondition, err = normalizeBCMBIAImpactLevelCondition(form.RspCondition, "rsp_condition"); err != nil {
			return nil, err
		}

		changes, err := utility.CompareFields(existingLevel, form, fieldsToCompare)
		if err != nil {
			log.Println("error :", err.Error())
		}
		allChanges = append(allChanges, changes...)

		existingLevel.ImpactAnalysisLevelID = form.ImpactAnalysisLevelID
		existingLevel.ClassificationEn = form.ClassificationEn
		existingLevel.ClassificationAr = form.ClassificationAr
		existingLevel.IsCritical = form.IsCritical
		existingLevel.RtoFromValue = form.RtoFromValue
		existingLevel.RtoToValue = form.RtoToValue
		existingLevel.RtoCondition = form.RtoCondition
		existingLevel.RtoUnitID = form.RtoUnitID
		existingLevel.MtpdFromValue = form.MtpdFromValue
		existingLevel.MtpdToValue = form.MtpdToValue
		existingLevel.MtpdCondition = form.MtpdCondition
		existingLevel.MtpdUnitID = form.MtpdUnitID
		existingLevel.RpoFromValue = form.RpoFromValue
		existingLevel.RpoToValue = form.RpoToValue
		existingLevel.RpoCondition = form.RpoCondition
		existingLevel.RpoUnitID = form.RpoUnitID
		existingLevel.RspFromValue = form.RspFromValue
		existingLevel.RspToValue = form.RspToValue
		existingLevel.RspCondition = form.RspCondition
		existingLevel.RspUnitID = form.RspUnitID

		if err := bbils.repo.Update(existingLevel); err != nil {
			return nil, err
		}

		updatedLevels = append(updatedLevels, existingLevel)
		updatedIDs = append(updatedIDs, existingLevel.ID)
	}

	if len(updatedIDs) > 0 {
		go func() {
			err := bbils.activityLogSvc.LogBulkActivities(
				models.UpdateAction,
				models.BCMBIAImpactLevel{}.TableName(),
				updatedIDs,
				allChanges,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk update activity: %v\n", err)
			}
		}()
	}

	return updatedLevels, nil
}

var validBCMBIAImpactLevelConditions = map[string]string{
	"exact":        "exact",
	"lessthan":     "less_than",
	"less_than":    "less_than",
	"range":        "range",
	"greaterthan":  "greater_than",
	"greater_than": "greater_than",
}

func normalizeBCMBIAImpactLevelCondition(condition *string, field string) (*string, error) {
	if condition == nil {
		return nil, nil
	}

	key := strings.ToLower(strings.TrimSpace(*condition))
	key = strings.ReplaceAll(key, " ", "_")
	key = strings.ReplaceAll(key, "-", "_")

	normalized, ok := validBCMBIAImpactLevelConditions[strings.ReplaceAll(key, "_", "")]
	if !ok {
		normalized, ok = validBCMBIAImpactLevelConditions[key]
	}
	if !ok {
		return nil, utility.NewValidationError(
			fmt.Errorf("invalid %s value: %s, allowed values are exact, less_than, range, greater_than", field, *condition),
			"invalid",
			field,
		)
	}

	return &normalized, nil
}
