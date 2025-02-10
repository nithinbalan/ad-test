package services

import (
	"errors"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementServiceInterface interface {
	GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentRequirementListQueryParamsDTO) ([]*complianceDTOs.ComplianceAssessmentRequirementListDTO, error)
	GetFullDetailsByID(assessmentId uint, requirementId uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentRequirementDetailDTO, error)
	Update(complianceAssessmentId uint, requirementId uint, currentUserID uint, complianceRequirement *complianceDTOs.ComplianceAssessmentRequirementUpdateFormDTO) error
}

// Service and implementations
type ComplianceAssessmentRequirementService struct {
	repo                    complianceRepositories.ComplianceAssessmentRequirementRepositoryInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	userSvc                 organizationServices.UserServiceInterface
	complianceRegisterSvc   ComplianceRegisterServiceInterface
	complianceAssessmentSvc ComplianceAssessmentServiceInterface
}

func NewComplianceAssessmentRequirementService(
	repo complianceRepositories.ComplianceAssessmentRequirementRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	complianceAssessmentSvc ComplianceAssessmentServiceInterface,
) ComplianceAssessmentRequirementServiceInterface {
	return &ComplianceAssessmentRequirementService{
		repo:                    repo,
		activityLogSvc:          activityLogSvc,
		userSvc:                 userSvc,
		complianceRegisterSvc:   complianceRegisterSvc,
		complianceAssessmentSvc: complianceAssessmentSvc,
	}
}

func (cras *ComplianceAssessmentRequirementService) GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentRequirementListQueryParamsDTO) ([]*complianceDTOs.ComplianceAssessmentRequirementListDTO, error) {

	if _, err := cras.complianceAssessmentSvc.GetByID(assessmentID); err != nil {
		return nil, err
	}

	return cras.repo.GetAll(assessmentID, params)
}

func (cras ComplianceAssessmentRequirementService) GetByID(id uint) (*models.ComplianceAssessmentRequirement, error) {
	complianceAssessmentRequirement, err := cras.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceAssessmentRequirement, nil
}

func (cras ComplianceAssessmentRequirementService) GetFullDetailsByID(assessmentId uint, requirementId uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentRequirementDetailDTO, error) {
	var complianceAssessmentRequirementDTO complianceDTOs.ComplianceAssessmentRequirementDetailDTO

	complianceAssessment, err := cras.repo.GetFullDetailsByID(requirementId, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, requirementId)
		}
		return nil, err
	}

	err = mapstructure.Decode(complianceAssessment, &complianceAssessmentRequirementDTO)
	if err != nil {
		return nil, err
	}

	complianceAssessmentRequirementDTO.CreatedAt = complianceAssessment.CreatedAt
	complianceAssessmentRequirementDTO.UpdatedAt = complianceAssessment.UpdatedAt
	// todo :- also add extra fields which required in details screen

	return &complianceAssessmentRequirementDTO, nil
}

func (cras ComplianceAssessmentRequirementService) Update(complianceAssessmentId uint, requirementId uint, currentUserID uint, complianceRequirementUpdateFormDTO *complianceDTOs.ComplianceAssessmentRequirementUpdateFormDTO) error {
	tx := cras.repo.BeginTransaction()

	existingComplianceAssessmentRequirement, err := cras.GetByID(requirementId)
	if err != nil {
		return err
	}

	oldComplianceAssessmentRequirement, err := cras.GetFullDetailsByID(complianceAssessmentId, requirementId)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComplianceAssessmentRequirement.Score = &complianceRequirementUpdateFormDTO.Score
	existingComplianceAssessmentRequirement.Comments = complianceRequirementUpdateFormDTO.Comments
	existingComplianceAssessmentRequirement.ComplianceStatusID = &complianceRequirementUpdateFormDTO.ComplianceStatusID
	existingComplianceAssessmentRequirement.UpdatedBYID = &currentUserID

	err = cras.repo.Update(existingComplianceAssessmentRequirement, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingAssessmentRequirement, err := cras.GetFullDetailsByID(complianceAssessmentId, requirementId)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"score":             oldComplianceAssessmentRequirement.Score,
			"comments":          oldComplianceAssessmentRequirement.Comments,
			"compliance_status": utility.GetFieldValue(oldComplianceAssessmentRequirement.ComplianceStatus),
		}

		newData := map[string]interface{}{
			"score":             existingAssessmentRequirement.Score,
			"comments":          existingAssessmentRequirement.Comments,
			"compliance_status": utility.GetFieldValue(existingAssessmentRequirement.ComplianceStatus),
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
				models.ComplianceAssessmentRequirement{}.TableName(),
				requirementId,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceAssessmentRequirement{}.TableName(),
					"recordID":  requirementId,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}

	}()

	return nil
}
