package services

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceAssessmentFrameworkServiceInterface interface {
	GetAll(params *complianceDTOs.ComplianceAssessmentFrameworkListQueryParamsDTO) ([]*complianceDTOs.ComplianceAssessmentFrameworkListDTO, error)
	GetByID(id uint) (*models.ComplianceAssessmentFramework, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentFrameworkDetailResponseDTO, error)
	Create(currentUserID uint, complianceAssessmentFrameworkFormDTO *complianceDTOs.ComplianceAssessmentFrameworkFormDTO) (uint, error)
	Update(id uint, currentUserID uint, complianceAssessmentFrameworkFormDTO *complianceDTOs.ComplianceAssessmentFrameworkFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	GetByFrameworkValueID(id uint) (*models.ComplianceAssessmentFrameworkValue, error)
}

// Service and implementations
type ComplianceAssessmentFrameworkService struct {
	repo           complianceRepositories.ComplianceAssessmentFrameworkRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewComplianceAssessmentFrameworkService(
	repo complianceRepositories.ComplianceAssessmentFrameworkRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) ComplianceAssessmentFrameworkServiceInterface {
	return &ComplianceAssessmentFrameworkService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (casf *ComplianceAssessmentFrameworkService) GetAll(params *complianceDTOs.ComplianceAssessmentFrameworkListQueryParamsDTO) ([]*complianceDTOs.ComplianceAssessmentFrameworkListDTO, error) {
	return casf.repo.GetAll(params)
}

func (casf *ComplianceAssessmentFrameworkService) GetByID(id uint) (*models.ComplianceAssessmentFramework, error) {
	complianceAssessmentFramework, err := casf.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceAssessmentFramework, nil
}

func (casf *ComplianceAssessmentFrameworkService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentFrameworkDetailResponseDTO, error) {
	// Initialize your DTO
	var complianceAssessmentFrameworkDTO complianceDTOs.ComplianceAssessmentFrameworkDetailResponseDTO

	complianceAssessmentFramework, err := casf.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(complianceAssessmentFramework, &complianceAssessmentFrameworkDTO)
	if err != nil {
		return nil, err
	}

	complianceAssessmentFrameworkDTO.CreatedAt = complianceAssessmentFramework.CreatedAt
	complianceAssessmentFrameworkDTO.UpdatedAt = complianceAssessmentFramework.UpdatedAt

	return &complianceAssessmentFrameworkDTO, nil
}

func (casf *ComplianceAssessmentFrameworkService) Create(currentUserID uint, complianceAssessmentFrameworkFormDTO *complianceDTOs.ComplianceAssessmentFrameworkFormDTO) (uint, error) {
	tx := casf.repo.BeginTransaction()

	complianceAssessmentFramework := &models.ComplianceAssessmentFramework{
		TitleEn:                   complianceAssessmentFrameworkFormDTO.TitleEn,
		TitleAr:                   complianceAssessmentFrameworkFormDTO.TitleAr,
		Description:               complianceAssessmentFrameworkFormDTO.Description,
		AssessmentFrameworkTypeID: complianceAssessmentFrameworkFormDTO.AssessmentFrameworkTypeID,
		CreatedBYID:               currentUserID,
	}

	id, err := casf.repo.Create(complianceAssessmentFramework, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	// Log activity asynchronously
	go func() {
		err := casf.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceAssessmentFramework{}.TableName(),
			id,
			map[string]interface{}{
				"title_en":                     complianceAssessmentFramework.TitleEn,
				"title_ar":                     complianceAssessmentFramework.TitleAr,
				"description":                  complianceAssessmentFramework.Description,
				"assessment_framework_type_id": complianceAssessmentFramework.AssessmentFrameworkTypeID,
			},
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceAssessmentFramework{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (casf *ComplianceAssessmentFrameworkService) Update(id uint, currentUserID uint, complianceAssessmentFrameworkFormDTO *complianceDTOs.ComplianceAssessmentFrameworkFormDTO) error {
	tx := casf.repo.BeginTransaction()

	existingComplianceAssessmentFramework, err := casf.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	oldComplianceAssessmentFramework, err := casf.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update only provided fields

	existingComplianceAssessmentFramework.TitleEn = complianceAssessmentFrameworkFormDTO.TitleEn
	existingComplianceAssessmentFramework.TitleAr = complianceAssessmentFrameworkFormDTO.TitleAr
	existingComplianceAssessmentFramework.Description = complianceAssessmentFrameworkFormDTO.Description
	existingComplianceAssessmentFramework.AssessmentFrameworkTypeID = complianceAssessmentFrameworkFormDTO.AssessmentFrameworkTypeID
	existingComplianceAssessmentFramework.UpdatedBYID = &currentUserID

	err = casf.repo.Update(existingComplianceAssessmentFramework, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity asynchronously
	go func() {
		newComplianceAssessmentFramework, err := casf.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"title_en":                     oldComplianceAssessmentFramework.TitleEn,
			"title_ar":                     oldComplianceAssessmentFramework.TitleAr,
			"description":                  oldComplianceAssessmentFramework.Description,
			"assessment_framework_type_id": oldComplianceAssessmentFramework.AssessmentFrameworkType,
		}

		newData := map[string]interface{}{
			"title_en":                     newComplianceAssessmentFramework.TitleEn,
			"title_ar":                     newComplianceAssessmentFramework.TitleAr,
			"description":                  newComplianceAssessmentFramework.Description,
			"assessment_framework_type_id": newComplianceAssessmentFramework.AssessmentFrameworkType,
		}

		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := casf.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ComplianceAssessmentFramework{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceAssessmentFramework{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}

func (casf *ComplianceAssessmentFrameworkService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := casf.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedComplianceAssessmentFramework []map[string]interface{}

	for _, id := range ids {
		complianceAssessmentFramework, err := casf.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedComplianceAssessmentFramework = append(deletedComplianceAssessmentFramework, map[string]interface{}{
			"id":                           id,
			"title_en":                     complianceAssessmentFramework.TitleEn,
			"title_ar":                     complianceAssessmentFramework.TitleAr,
			"description":                  complianceAssessmentFramework.Description,
			"assessment_framework_type_id": complianceAssessmentFramework.AssessmentFrameworkType,
		})

		// Delete the business
		if err := casf.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := casf.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ComplianceAssessmentFramework{}.TableName(),
				successfulIds,
				deletedComplianceAssessmentFramework,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

		}()
	}

	return failedIds, nil
}
func (casf *ComplianceAssessmentFrameworkService) GetByFrameworkValueID(id uint) (*models.ComplianceAssessmentFrameworkValue, error) {
	return casf.repo.GetByFrameworkValueID(id)
}
