package services

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceAssessmentServiceInterface interface {
	GetAll(params *complianceDTOs.ComplianceAssessmentListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentDetailDTO, error)
	GetByID(id uint) (*models.ComplianceAssessment, error)
	Create(currentUserID uint, complianceAssessmentFormDTO *complianceDTOs.ComplianceAssessmentFormDTO) error
	Update(id uint, currentUserID uint, compliance *complianceDTOs.ComplianceAssessmentUpdateFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type ComplianceAssessmentService struct {
	repo           complianceRepositories.ComplianceAssessmentRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
}

func NewComplianceAssessmentService(
	repo complianceRepositories.ComplianceAssessmentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) ComplianceAssessmentServiceInterface {
	return &ComplianceAssessmentService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
	}
}

func (cras *ComplianceAssessmentService) GetAll(params *complianceDTOs.ComplianceAssessmentListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentListResponseDTO, error) {
	return cras.repo.GetAll(params)
}

func (cras ComplianceAssessmentService) GetByID(id uint) (*models.ComplianceAssessment, error) {
	compliance, err := cras.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return compliance, nil
}

func (cras ComplianceAssessmentService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentDetailDTO, error) {
	var complianceAssessmentDTO complianceDTOs.ComplianceAssessmentDetailDTO

	complianceAssessment, err := cras.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(complianceAssessment, &complianceAssessmentDTO)
	if err != nil {
		return nil, err
	}

	complianceAssessmentDTO.CreatedAt = complianceAssessment.CreatedAt
	complianceAssessmentDTO.UpdatedAt = complianceAssessment.UpdatedAt
	// todo :- also add extra fields which required in details screen

	return &complianceAssessmentDTO, nil
}

func (cras *ComplianceAssessmentService) Create(currentUserID uint, complianceAssessmentFormDTO *complianceDTOs.ComplianceAssessmentFormDTO) error {
	tx := cras.repo.BeginTransaction()
	var createdIds []uint
	var addedComplianceAssessments []map[string]interface{}

	for _, organizationUnit := range complianceAssessmentFormDTO.OrganizationUnitIDs {
		complianceAssessment := &models.ComplianceAssessment{
			ComplianceRegisterID: complianceAssessmentFormDTO.ComplianceRegisterID,
			ReviewUserID:         complianceAssessmentFormDTO.ReviewUserID,
			OrganizationUnitID:   *organizationUnit,
			Description:          complianceAssessmentFormDTO.Description,
			CreatedBYID:          currentUserID,
		}
		id, err := cras.repo.Create(complianceAssessment, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		addedComplianceAssessments = append(addedComplianceAssessments, map[string]interface{}{
			"id":                     id,
			"compliance__assessment": complianceAssessment,
		})

		createdIds = append(createdIds, id)
	}

	tx.Commit()

	go func() {

		err := cras.activityLogSvc.LogBulkActivities(
			models.CreateAction,
			models.ComplianceAssessment{}.TableName(),
			createdIds,
			addedComplianceAssessments,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceAssessment{}.TableName(),
				"recordID":  createdIds,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (cras *ComplianceAssessmentService) Update(id uint, currentUserID uint, complianceAssessmentUpdateFormDTO *complianceDTOs.ComplianceAssessmentUpdateFormDTO) error {
	tx := cras.repo.BeginTransaction()

	existingComplianceAssessment, err := cras.GetByID(id)
	if err != nil {
		return err
	}

	oldComplianceAssessment, err := cras.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComplianceAssessment.ComplianceRegisterID = complianceAssessmentUpdateFormDTO.ComplianceRegisterID
	existingComplianceAssessment.ReviewUserID = complianceAssessmentUpdateFormDTO.ReviewUserID
	existingComplianceAssessment.OrganizationUnitID = complianceAssessmentUpdateFormDTO.OrganizationUnitID
	existingComplianceAssessment.Description = complianceAssessmentUpdateFormDTO.Description
	existingComplianceAssessment.UpdatedBYID = &currentUserID

	err = cras.repo.Update(existingComplianceAssessment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingAssessment, err := cras.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"compliance_register": oldComplianceAssessment.ComplianceRegister.TitleEn,
			"organization_unit":   oldComplianceAssessment.OrganizationUnit.TitleEn,
			"description":         oldComplianceAssessment.Description,
			"review_user":         fmt.Sprintf("%s %s", oldComplianceAssessment.ReviewUser.FirstNameEN, oldComplianceAssessment.ReviewUser.LastNameEN),
		}

		newData := map[string]interface{}{
			"compliance_register": existingAssessment.ComplianceRegister.TitleEn,
			"organization_unit":   existingAssessment.OrganizationUnit.TitleEn,
			"description":         existingAssessment.Description,
			"review_user":         fmt.Sprintf("%s %s", existingAssessment.ReviewUser.FirstNameEN, existingAssessment.ReviewUser.LastNameEN),
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			// Log the update activity
			err := cras.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ComplianceAssessment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceAssessment{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}

	}()

	return nil
}

func (cras *ComplianceAssessmentService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := cras.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var deletedComplianceAssessment []map[string]interface{}

	for _, id := range ids {
		complianceAssessment, err := cras.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Delete the Compliance Register
		if err := cras.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedComplianceAssessment = append(deletedComplianceAssessment, map[string]interface{}{
			"id":                id,
			"title_en":          complianceAssessment.ComplianceRegister.TitleEn,
			"organization_unit": complianceAssessment.OrganizationUnit.TitleAr,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := cras.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ComplianceAssessment{}.TableName(),
				successfulIds,
				deletedComplianceAssessment,
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
