package services

import (
	"errors"
	"time"

	"jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceAssessmentFrameworkValueServiceInterface interface {
	GetAll(complianceAssessmentFrameworkID uint, params *dtos.ComplianceAssessmentFrameworkValueListQueryParamsDTO) (*dtos.ComplianceAssessmentFrameworkValueListResponseDTO, error)
	GetByID(complianceAssessmentFrameworkID uint, id uint) (*models.ComplianceAssessmentFrameworkValue, error)
	GetFullDetailsByID(complianceAssessmentFrameworkID uint, id uint) (*dtos.ComplianceAssessmentFrameworkValueDetailDTO, error)
	Create(complianceAssessmentFrameworkID uint, currentUserID uint, form *dtos.ComplianceAssessmentFrameworkValueFormDTO) (uint, error)
	Update(complianceAssessmentFrameworkID uint, id uint, currentUserID uint, form *dtos.ComplianceAssessmentFrameworkValueFormDTO) error
	Delete(complianceAssessmentFrameworkID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type ComplianceAssessmentFrameworkValueService struct {
	repo                             repositories.ComplianceAssessmentFrameworkValueRepositoryInterface
	activityLog                      logServices.ActivityLogServiceInterface
	userSvc                          organizationServices.UserServiceInterface
	complianceAssessmentFrameworkSvc ComplianceAssessmentFrameworkServiceInterface
	complianceStatusSvc              ComplianceStatusServiceInterface
}

func NewComplianceAssessmentFrameworkValueService(
	repo repositories.ComplianceAssessmentFrameworkValueRepositoryInterface,
	activityLog logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	complianceAssessmentFrameworkSvc ComplianceAssessmentFrameworkServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
) ComplianceAssessmentFrameworkValueServiceInterface {
	return &ComplianceAssessmentFrameworkValueService{
		repo:                             repo,
		activityLog:                      activityLog,
		userSvc:                          userSvc,
		complianceAssessmentFrameworkSvc: complianceAssessmentFrameworkSvc,
		complianceStatusSvc:              complianceStatusSvc,
	}
}

func (cavs *ComplianceAssessmentFrameworkValueService) GetAll(complianceAssessmentFrameworkID uint, params *dtos.ComplianceAssessmentFrameworkValueListQueryParamsDTO) (*dtos.ComplianceAssessmentFrameworkValueListResponseDTO, error) {
	return cavs.repo.GetAll(complianceAssessmentFrameworkID, params)
}

func (cavs *ComplianceAssessmentFrameworkValueService) GetByID(complianceAssessmentFrameworkID uint, id uint) (*models.ComplianceAssessmentFrameworkValue, error) {
	if _, err := cavs.complianceAssessmentFrameworkSvc.GetByID(complianceAssessmentFrameworkID); err != nil {
		return nil, err
	}
	value, err := cavs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return value, nil
}

func (cavs *ComplianceAssessmentFrameworkValueService) GetFullDetailsByID(complianceAssessmentFrameworkID uint, id uint) (*dtos.ComplianceAssessmentFrameworkValueDetailDTO, error) {
	var complianceAssessmentFrameworkValueDTO dtos.ComplianceAssessmentFrameworkValueDetailDTO

	// Check if the compliance assessment framework exists
	if _, err := cavs.complianceAssessmentFrameworkSvc.GetByID(complianceAssessmentFrameworkID); err != nil {
		return nil, err
	}

	complianceAssessmentFrameworkValue, err := cavs.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utility.ResourceNotFoundError{}
		}
		return nil, err
	}

	err = mapstructure.Decode(complianceAssessmentFrameworkValue, &complianceAssessmentFrameworkValueDTO)
	if err != nil {
		return nil, err
	}

	complianceAssessmentFrameworkValueDTO.CreatedAt = complianceAssessmentFrameworkValue.CreatedAt
	complianceAssessmentFrameworkValueDTO.UpdatedAt = complianceAssessmentFrameworkValue.UpdatedAt

	return &complianceAssessmentFrameworkValueDTO, nil
}

func (cavs *ComplianceAssessmentFrameworkValueService) Create(complianceAssessmentFrameworkID uint, currentUserID uint, form *dtos.ComplianceAssessmentFrameworkValueFormDTO) (uint, error) {
	complianceAssessmentFramework, err := cavs.complianceAssessmentFrameworkSvc.GetByID(complianceAssessmentFrameworkID)
	if err != nil {
		return 0, err
	}

	// Handle IsApplicable with proper default value
	var isApplicable *bool
	if form.IsApplicable != nil {
		isApplicable = form.IsApplicable
	} else {
		// Default to true when not specified
		defaultValue := true
		isApplicable = &defaultValue
	}

	complianceAssessmentFrameworkValue := &models.ComplianceAssessmentFrameworkValue{
		ComplianceAssessmentFrameworkID: complianceAssessmentFramework.ID,
		TitleEn:                         form.TitleEn,
		TitleAr:                         form.TitleAr,
		Description:                     form.Description,
		Value:                           form.Value,
		ScoreMin:                        form.ScoreMin,
		ScoreMax:                        form.ScoreMax,
		IsApplicable:                    isApplicable,
		ComplianceStatusID:              form.ComplianceStatusID,
		CreatedBYID:                     currentUserID,
		CreatedAt:                       time.Now(),
	}

	id, err := cavs.repo.Create(complianceAssessmentFrameworkValue)
	if err != nil {
		return 0, err
	}

	go func() {
		activityData := map[string]interface{}{
			"compliance_assessment_framework_id": complianceAssessmentFrameworkID,
			"title_en":                           form.TitleEn,
			"title_ar":                           form.TitleAr,
			"description":                        form.Description,
			"value":                              form.Value,
			"score_min":                          form.ScoreMin,
			"score_max":                          form.ScoreMax,
			"is_applicable":                      isApplicable,
			"compliance_status_id":               form.ComplianceStatusID,
		}
		err := cavs.activityLog.LogActivity(models.CreateAction, models.ComplianceAssessmentFrameworkValue{}.TableName(), id, activityData, currentUserID, nil)
		if err != nil {
			utility.LogError("Failed to log create compliance assessment framework value activity", map[string]interface{}{
				"recordID": id,
				"error":    err.Error(),
			})
		}
	}()

	return id, nil
}

func (cavs *ComplianceAssessmentFrameworkValueService) Update(complianceAssessmentFrameworkID uint, id uint, currentUserID uint, form *dtos.ComplianceAssessmentFrameworkValueFormDTO) error {
	complianceAssessmentFrameworkValue, err := cavs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utility.ResourceNotFoundError{}
		}
		return err
	}

	// Store old data for activity logging
	oldData := map[string]interface{}{
		"title_en":             complianceAssessmentFrameworkValue.TitleEn,
		"title_ar":             complianceAssessmentFrameworkValue.TitleAr,
		"description":          complianceAssessmentFrameworkValue.Description,
		"value":                complianceAssessmentFrameworkValue.Value,
		"score_min":            complianceAssessmentFrameworkValue.ScoreMin,
		"score_max":            complianceAssessmentFrameworkValue.ScoreMax,
		"is_applicable":        complianceAssessmentFrameworkValue.IsApplicable,
		"compliance_status_id": complianceAssessmentFrameworkValue.ComplianceStatusID,
	}

	// Handle IsApplicable with proper default value
	var isApplicable *bool
	if form.IsApplicable != nil {
		isApplicable = form.IsApplicable
	} else {
		// Default to true when not specified
		defaultValue := true
		isApplicable = &defaultValue
	}

	// Update fields
	complianceAssessmentFrameworkValue.TitleEn = form.TitleEn
	complianceAssessmentFrameworkValue.TitleAr = form.TitleAr
	complianceAssessmentFrameworkValue.Description = form.Description
	complianceAssessmentFrameworkValue.Value = form.Value
	complianceAssessmentFrameworkValue.ScoreMin = form.ScoreMin
	complianceAssessmentFrameworkValue.ScoreMax = form.ScoreMax
	complianceAssessmentFrameworkValue.IsApplicable = isApplicable
	complianceAssessmentFrameworkValue.ComplianceStatusID = form.ComplianceStatusID
	complianceAssessmentFrameworkValue.UpdatedBYID = &currentUserID
	complianceAssessmentFrameworkValue.UpdatedAt = timePtr(time.Now())

	err = cavs.repo.Update(complianceAssessmentFrameworkValue)
	if err != nil {
		return err
	}

	go func() {
		newData := map[string]interface{}{
			"title_en":             form.TitleEn,
			"title_ar":             form.TitleAr,
			"description":          form.Description,
			"value":                form.Value,
			"score_min":            form.ScoreMin,
			"score_max":            form.ScoreMax,
			"is_applicable":        isApplicable,
			"compliance_status_id": form.ComplianceStatusID,
		}
		changes := utility.CompareFieldChanges(oldData, newData)
		if len(changes) > 0 {
			err := cavs.activityLog.LogActivity(models.UpdateAction, models.ComplianceAssessmentFrameworkValue{}.TableName(), id, changes, currentUserID, nil)
			if err != nil {
				utility.LogError("Failed to log update compliance assessment framework value activity", map[string]interface{}{
					"recordID": id,
					"error":    err.Error(),
				})
			}
		}
	}()

	return nil
}

func (cavs *ComplianceAssessmentFrameworkValueService) Delete(complianceAssessmentFrameworkID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := cavs.complianceAssessmentFrameworkSvc.GetByID(complianceAssessmentFrameworkID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedValues []map[string]interface{}

	for _, id := range ids {
		value, err := cavs.repo.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		err = cavs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedValues = append(deletedValues, map[string]interface{}{
			"id":       id,
			"title_en": value.TitleEn,
			"title_ar": value.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := cavs.activityLog.LogBulkActivities(
				models.DeleteAction,
				models.ComplianceAssessmentFrameworkValue{}.TableName(),
				successfulIds,
				deletedValues,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete compliance assessment framework value activity", map[string]interface{}{
					"recordIDs": successfulIds,
					"error":     err.Error(),
				})
			}
		}()
	}

	return failedIds, nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}
