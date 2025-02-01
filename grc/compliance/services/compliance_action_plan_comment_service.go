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

	"gorm.io/gorm"
)

type ComplianceActionPlanCommentServiceInterface interface {
	GetAll(actionPlanID uint) ([]*complianceDTOs.ComplianceActionPlanCommentListDTO, error)
	GetByID(actionPlanID uint, id uint) (*models.ComplianceActionPlanComment, error)
	GetFullDetailsByID(id uint) (*models.ComplianceActionPlanComment, error)
	Create(actionPlanID uint, currentUserID uint, dto *complianceDTOs.ComplianceActionPlanCommentFormDTO) (uint, error)
	Update(actionPlanID uint, id uint, currentUserID uint, dto *complianceDTOs.ComplianceActionPlanCommentFormDTO) error
	Delete(actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type ComplianceActionPlanCommentService struct {
	repo           repositories.ComplianceActionPlanCommentRepositoryInterface
	userSvc        organizationServices.UserServiceInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewComplianceActionPlanCommentService(repo repositories.ComplianceActionPlanCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) ComplianceActionPlanCommentServiceInterface {
	return &ComplianceActionPlanCommentService{
		repo:           repo,
		userSvc:        userSvc,
		activityLogSvc: activityLogSvc,
	}
}

func (caps *ComplianceActionPlanCommentService) GetAll(actionPlanID uint) ([]*complianceDTOs.ComplianceActionPlanCommentListDTO, error) {
	return caps.repo.GetAll(actionPlanID)
}

func (caps *ComplianceActionPlanCommentService) GetByID(actionPlanID uint, id uint) (*models.ComplianceActionPlanComment, error) {
	comment, err := caps.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (caps *ComplianceActionPlanCommentService) GetFullDetailsByID(id uint) (*models.ComplianceActionPlanComment, error) {
	return caps.repo.GetFullDetailsByID(id)
}

func (caps *ComplianceActionPlanCommentService) Create(actionPlanID uint, currentUserID uint, dto *complianceDTOs.ComplianceActionPlanCommentFormDTO) (uint, error) {
	tx := caps.repo.BeginTransaction()

	comment := &models.ComplianceActionPlanComment{
		ComplianceActionPlanID: actionPlanID,
		Comment:                dto.Comment,
		CreatedBYID:            currentUserID,
	}

	id, err := caps.repo.Create(comment)
	if err != nil {
		return 0, err
	}

	tx.Commit()
	go func() {
		existingComment, err := caps.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch test and exercise", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = caps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}

func (caps *ComplianceActionPlanCommentService) Update(actionPlanID uint, id uint, currentUserID uint, dto *complianceDTOs.ComplianceActionPlanCommentFormDTO) error {
	tx := caps.repo.BeginTransaction()

	existingComment, err := caps.repo.GetFullDetailsByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldComment, err := caps.GetFullDetailsByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = dto.Comment
	err = caps.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingComment, err := caps.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		changedData := make(map[string]interface{})

		// log changes
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
				models.TestAndExerciseRecoveryActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.TestAndExerciseRecoveryActionPlan{}.TableName(),
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

func (caps *ComplianceActionPlanCommentService) Delete(actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := caps.repo.BeginTransaction()

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		if _, err := caps.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":                id,
			"compliance_action": actionPlanID,
		})

		err := caps.repo.Delete(id, tx)
		if err != nil {
			failedIDs = append(failedIDs, id)
		}
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, errors.New("some deletions failed")
	}

	tx.Commit()
	if len(successfulIDs) > 0 {
		go func() {
			err := caps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.TestAndExerciseRecoveryActionPlan{}.TableName(),
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
