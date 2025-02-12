package services

import (
	"errors"
	"fmt"

	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementCommentServiceInterface interface {
	GetAll(complianceAssessmentID uint, assessmentRequirementID uint) ([]*complianceDTOs.ComplianceAssessmentRequirementCommentListDTO, error)
	GetByID(complianceAssessmentID uint, assessmentRequirementID uint, id uint) (*models.ComplianceAssessmentRequirementComment, error)
	GetFullDetailsByID(complianceAssessmentID uint, assessmentRequirementID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentRequirementCommentDetailDTO, error)
	Create(complianceAssessmentID uint, assessmentRequirementID uint, currentUserID uint, dto *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) (uint, error)
	Update(complianceAssessmentID uint, assessmentRequirementID uint, id uint, currentUserID uint, dto *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) error
	Delete(complianceAssessmentID uint, assessmentRequirementID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type ComplianceAssessmentRequirementCommentService struct {
	repo                               repositories.ComplianceAssessmentRequirementCommentRepositoryInterface
	userSvc                            organizationServices.UserServiceInterface
	activityLogSvc                     logServices.ActivityLogServiceInterface
	complianceAssessmentRequirementSvc ComplianceAssessmentRequirementServiceInterface
}

func NewComplianceAssessmentRequirementCommentService(repo repositories.ComplianceAssessmentRequirementCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceAssessmentRequirementSvc ComplianceAssessmentRequirementServiceInterface,
) ComplianceAssessmentRequirementCommentServiceInterface {
	return &ComplianceAssessmentRequirementCommentService{
		repo:                               repo,
		userSvc:                            userSvc,
		activityLogSvc:                     activityLogSvc,
		complianceAssessmentRequirementSvc: complianceAssessmentRequirementSvc,
	}
}

func (caps *ComplianceAssessmentRequirementCommentService) GetAll(complianceAssessmentID uint, assessmentRequirementID uint) ([]*complianceDTOs.ComplianceAssessmentRequirementCommentListDTO, error) {
	if _, err := caps.complianceAssessmentRequirementSvc.GetByID(assessmentRequirementID); err != nil {
		return nil, err
	}

	return caps.repo.GetAll(assessmentRequirementID)
}

func (caps *ComplianceAssessmentRequirementCommentService) GetByID(complianceAssessmentID uint, assessmentRequirementID uint, id uint) (*models.ComplianceAssessmentRequirementComment, error) {
	if _, err := caps.complianceAssessmentRequirementSvc.GetByID(assessmentRequirementID); err != nil {
		return nil, err
	}

	comment, err := caps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (caps *ComplianceAssessmentRequirementCommentService) GetFullDetailsByID(complianceAssessmentID uint, assessmentRequirementID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentRequirementCommentDetailDTO, error) {
	if _, err := caps.complianceAssessmentRequirementSvc.GetByID(assessmentRequirementID); err != nil {
		fmt.Printf("ComplianceAssessmentRequirement not found for complianceAssessmentID %d, assessmentRequirementID %d\n", complianceAssessmentID, assessmentRequirementID)
		return nil, err
	}

	var commentDetail complianceDTOs.ComplianceAssessmentRequirementCommentDetailDTO

	comment, err := caps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(comment, &commentDetail)
	if err != nil {
		fmt.Printf("Error decoding comment: %v\n", err)
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt
	return &commentDetail, nil
}

func (caps *ComplianceAssessmentRequirementCommentService) Create(complianceAssessmentID uint, assessmentRequirementID uint, currentUserID uint, complianceAssessmentRequirementCommentFormDTO *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) (uint, error) {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.complianceAssessmentRequirementSvc.GetByID(assessmentRequirementID); err != nil {
		return 0, err
	}

	comment := &models.ComplianceAssessmentRequirementComment{
		ComplianceAssessmentRequirementID: assessmentRequirementID,
		Comment:                           complianceAssessmentRequirementCommentFormDTO.Comment,
		CreatedBYID:                       currentUserID,
	}

	id, err := caps.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingComment, err := caps.GetFullDetailsByID(complianceAssessmentID, assessmentRequirementID, id)
		if err != nil {
			utility.LogError("Failed to fetch compliance assessment requirement comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = caps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceAssessmentRequirementComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceAssessmentRequirementComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}

func (caps *ComplianceAssessmentRequirementCommentService) Update(complianceAssessmentID uint, assessmentRequirementID uint, id uint, currentUserID uint, complianceAssessmentRequirementCommentFormDTO *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) error {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.complianceAssessmentRequirementSvc.GetByID(assessmentRequirementID); err != nil {
		return err
	}

	existingComment, err := caps.GetByID(complianceAssessmentID, assessmentRequirementID, id)
	if err != nil {
		return err
	}

	oldComment, err := caps.GetFullDetailsByID(complianceAssessmentID, assessmentRequirementID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = complianceAssessmentRequirementCommentFormDTO.Comment
	err = caps.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingComment, err := caps.GetFullDetailsByID(complianceAssessmentID, assessmentRequirementID, id)
		if err != nil {
			return
		}

		changedData := make(map[string]interface{})

		oldData := map[string]interface{}{
			"comment": oldComment.Comment,
		}

		newData := map[string]interface{}{
			"comment": existingComment.Comment,
		}

		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		if len(changedData) > 0 {
			err := caps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ComplianceAssessmentRequirementComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceAssessmentRequirementComment{}.TableName(),
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

func (caps *ComplianceAssessmentRequirementCommentService) Delete(complianceAssessmentID uint, assessmentRequirementID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := caps.repo.BeginTransaction()
	if _, err := caps.complianceAssessmentRequirementSvc.GetByID(assessmentRequirementID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		if _, err := caps.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":                                id,
			"compliance_assessment_requirement": assessmentRequirementID,
		})

		err := caps.repo.Delete(id, tx)
		if err != nil {
			failedIDs = append(failedIDs, id)
		}
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()
	if len(successfulIDs) > 0 {
		go func() {
			err := caps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ComplianceAssessmentRequirementComment{}.TableName(),
				successfulIDs,
				deletedComments,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}
	return failedIDs, nil
}
