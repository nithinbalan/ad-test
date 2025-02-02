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

type ComplianceActionPlanCommentServiceInterface interface {
	GetAll(complianceID uint, actionPlanID uint) ([]*complianceDTOs.ComplianceActionPlanCommentListDTO, error)
	GetByID(complianceID uint, actionPlanID uint, id uint) (*models.ComplianceActionPlanComment, error)
	GetFullDetailsByID(complianceID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceActionPlanCommentDetailDTO, error)
	Create(complianceID uint, actionPlanID uint, currentUserID uint, dto *complianceDTOs.ComplianceActionPlanCommentFormDTO) (uint, error)
	Update(complianceID uint, actionPlanID uint, id uint, currentUserID uint, dto *complianceDTOs.ComplianceActionPlanCommentFormDTO) error
	Delete(complianceID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type ComplianceActionPlanCommentService struct {
	repo                    repositories.ComplianceActionPlanCommentRepositoryInterface
	userSvc                 organizationServices.UserServiceInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	complianceActionPlanSvc ComplianceActionPlanServiceInterface
}

func NewComplianceActionPlanCommentService(repo repositories.ComplianceActionPlanCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceActionPlanSvc ComplianceActionPlanServiceInterface,
) ComplianceActionPlanCommentServiceInterface {
	return &ComplianceActionPlanCommentService{
		repo:                    repo,
		userSvc:                 userSvc,
		activityLogSvc:          activityLogSvc,
		complianceActionPlanSvc: complianceActionPlanSvc,
	}
}

func (caps *ComplianceActionPlanCommentService) GetAll(complianceID uint, actionPlanID uint) ([]*complianceDTOs.ComplianceActionPlanCommentListDTO, error) {
	if _, err := caps.complianceActionPlanSvc.GetByID(complianceID, actionPlanID); err != nil {
		return nil, err
	}

	return caps.repo.GetAll(actionPlanID)
}

func (caps *ComplianceActionPlanCommentService) GetByID(complianceID uint, actionPlanID uint, id uint) (*models.ComplianceActionPlanComment, error) {
	if _, err := caps.complianceActionPlanSvc.GetByID(complianceID, actionPlanID); err != nil {
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

func (caps *ComplianceActionPlanCommentService) GetFullDetailsByID(complianceID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceActionPlanCommentDetailDTO, error) {
	if _, err := caps.complianceActionPlanSvc.GetByID(complianceID, actionPlanID); err != nil {
		fmt.Printf("ComplianceActionPlan not found for complianceID %d, actionPlanID %d\n", complianceID, actionPlanID)
		return nil, err
	}

	var commentDetail complianceDTOs.ComplianceActionPlanCommentDetailDTO

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

func (caps *ComplianceActionPlanCommentService) Create(complianceID uint, actionPlanID uint, currentUserID uint, complianceActionPlanCommentFormDTO *complianceDTOs.ComplianceActionPlanCommentFormDTO) (uint, error) {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.complianceActionPlanSvc.GetByID(complianceID, actionPlanID); err != nil {
		return 0, err
	}

	comment := &models.ComplianceActionPlanComment{
		ComplianceActionPlanID: actionPlanID,
		Comment:                complianceActionPlanCommentFormDTO.Comment,
		CreatedBYID:            currentUserID,
	}

	id, err := caps.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingComment, err := caps.GetFullDetailsByID(complianceID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch compliance action plan comment", map[string]interface{}{
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

func (caps *ComplianceActionPlanCommentService) Update(complianceID uint, actionPlanID uint, id uint, currentUserID uint, complianceActionPlanCommentFormDTO *complianceDTOs.ComplianceActionPlanCommentFormDTO) error {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.complianceActionPlanSvc.GetByID(complianceID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := caps.GetByID(complianceID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := caps.GetFullDetailsByID(complianceID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = complianceActionPlanCommentFormDTO.Comment
	err = caps.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingComment, err := caps.GetFullDetailsByID(complianceID, actionPlanID, id)
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
				models.ComplianceActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceActionPlanComment{}.TableName(),
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

func (caps *ComplianceActionPlanCommentService) Delete(complianceID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := caps.repo.BeginTransaction()
	if _, err := caps.complianceActionPlanSvc.GetByID(complianceID, actionPlanID); err != nil {
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
				models.ComplianceActionPlanComment{}.TableName(),
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
