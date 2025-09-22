package services

import (
	"errors"
	"fmt"
	"time"

	logServices "jethur/grc/logs/services"
	masterService "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/grc/policy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PolicyVersionClauseCommentServiceInterface interface {
	GetAll(documentID uint, clauseID uint, params *policyDTOs.PolicyVersionClauseCommentListQueryParamsDTO) ([]*policyDTOs.PolicyVersionClauseCommentListDTO, error)
	GetAllDocumentComments(documentID uint, currentUserID uint, params *policyDTOs.PolicyVersionClauseCommentListQueryParamsDTO) (*policyDTOs.PolicyVersionClauseCommentListResponseDTO, error)
	GetByID(documentID uint, clauseID uint, id uint) (*models.DocumentVersionClauseComment, error)
	GetFullDetailsByID(documentID uint, clauseID uint, id uint, tx ...utility.DatabaseTransaction) (*policyDTOs.PolicyVersionClauseCommentDetailDTO, error)
	Create(documentID uint, clauseID uint, currentUserID uint, dto *policyDTOs.PolicyVersionClauseCommentFormDTO) (uint, error)
	Update(documentID uint, clauseID uint, id uint, currentUserID uint, dto *policyDTOs.PolicyVersionClauseCommentFormDTO) error
	Delete(documentID uint, clauseID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(documentID uint, clauseID uint, commentId uint, currentUserID uint) error
}

type PolicyVersionClauseCommentService struct {
	repo              repositories.PolicyVersionClauseCommentRepositoryInterface
	userSvc           organizationServices.UserServiceInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
	policyRegisterSvc PolicyRegisterServiceInterface
	commentStatusSvc  masterService.CommentStatusServiceInterface
}

func NewPolicyVersionClauseCommentService(repo repositories.PolicyVersionClauseCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	policyRegisterSvc PolicyRegisterServiceInterface,
	commentStatusSvc masterService.CommentStatusServiceInterface,
) PolicyVersionClauseCommentServiceInterface {
	return &PolicyVersionClauseCommentService{
		repo:              repo,
		userSvc:           userSvc,
		activityLogSvc:    activityLogSvc,
		policyRegisterSvc: policyRegisterSvc,
		commentStatusSvc:  commentStatusSvc,
	}
}

func (pvccs *PolicyVersionClauseCommentService) GetAll(documentID uint, clauseID uint, params *policyDTOs.PolicyVersionClauseCommentListQueryParamsDTO) ([]*policyDTOs.PolicyVersionClauseCommentListDTO, error) {
	if _, err := pvccs.policyRegisterSvc.GetByID(documentID); err != nil {
		return nil, err
	}

	return pvccs.repo.GetAll(clauseID, params)
}

func (pvccs *PolicyVersionClauseCommentService) GetAllDocumentComments(documentID uint, currentUserID uint, params *policyDTOs.PolicyVersionClauseCommentListQueryParamsDTO) (*policyDTOs.PolicyVersionClauseCommentListResponseDTO, error) {
	policyRegister, err := pvccs.policyRegisterSvc.GetFullDetailsByID(documentID, &currentUserID)
	if err != nil {
		return nil, err
	}

	return pvccs.repo.GetAllDocumentComments(documentID, policyRegister.LatestDocumentVersion.ID, params)
}

func (pvccs *PolicyVersionClauseCommentService) GetByID(documentID uint, clauseID uint, id uint) (*models.DocumentVersionClauseComment, error) {
	if _, err := pvccs.policyRegisterSvc.GetByID(documentID); err != nil {
		return nil, err
	}

	comment, err := pvccs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (pvccs *PolicyVersionClauseCommentService) GetFullDetailsByID(documentID uint, clauseID uint, id uint, tx ...utility.DatabaseTransaction) (*policyDTOs.PolicyVersionClauseCommentDetailDTO, error) {
	if _, err := pvccs.policyRegisterSvc.GetByID(documentID); err != nil {
		fmt.Printf("PolicyVersionClause not found for documentID %d, versionID %d\n", documentID, clauseID)
		return nil, err
	}

	var commentDetail policyDTOs.PolicyVersionClauseCommentDetailDTO

	comment, err := pvccs.repo.GetFullDetailsByID(id, tx...)
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

func (pvccs *PolicyVersionClauseCommentService) Create(documentID uint, clauseID uint, currentUserID uint, commentFormDTO *policyDTOs.PolicyVersionClauseCommentFormDTO) (uint, error) {
	tx := pvccs.repo.BeginTransaction()

	policyRegister, err := pvccs.policyRegisterSvc.GetFullDetailsByID(documentID, &currentUserID)
	if err != nil {
		return 0, err
	}

	commentStatus, err := pvccs.commentStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}

	comment := &models.DocumentVersionClauseComment{
		DocumentVersionClauseID: clauseID,
		Comment:                 commentFormDTO.Comment,
		CommentStatusID:         commentStatus.ID,
		Level:                   *policyRegister.LatestDocumentVersion.NextReviewLevel,
		CreatedBYID:             currentUserID,
	}

	id, err := pvccs.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingComment, err := pvccs.GetFullDetailsByID(documentID, clauseID, id)
		if err != nil {
			utility.LogError("Failed to fetch comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = pvccs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.DocumentVersionClauseComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.DocumentVersionClauseComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

	}()
	return id, nil
}

func (pvccs *PolicyVersionClauseCommentService) Update(documentID uint, clauseID uint, id uint, currentUserID uint, commentFormDTO *policyDTOs.PolicyVersionClauseCommentFormDTO) error {
	tx := pvccs.repo.BeginTransaction()

	if _, err := pvccs.policyRegisterSvc.GetByID(documentID); err != nil {
		return err
	}

	existingComment, err := pvccs.GetByID(documentID, clauseID, id)
	if err != nil {
		return err
	}

	oldComment, err := pvccs.GetFullDetailsByID(documentID, clauseID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = commentFormDTO.Comment
	err = pvccs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingComment, err := pvccs.GetFullDetailsByID(documentID, clauseID, id)
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
			err := pvccs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.DocumentVersionClauseComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				"", "", "", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.DocumentVersionClauseComment{}.TableName(),
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

func (pvccs *PolicyVersionClauseCommentService) Delete(documentID uint, clauseID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := pvccs.repo.BeginTransaction()
	if _, err := pvccs.policyRegisterSvc.GetByID(documentID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		if _, err := pvccs.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":     id,
			"clause": clauseID,
		})

		err := pvccs.repo.Delete(id, tx)
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
			err := pvccs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.DocumentVersionClauseComment{}.TableName(),
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

func (pvccs *PolicyVersionClauseCommentService) Resolve(documentID uint, clauseID uint, commentId uint, currentUserID uint) error {

	tx := pvccs.repo.BeginTransaction()

	if _, err := pvccs.policyRegisterSvc.GetByID(documentID); err != nil {
		return err
	}

	existingComment, err := pvccs.GetByID(documentID, clauseID, commentId)
	if err != nil {
		return err
	}

	commentStatus, err := pvccs.commentStatusSvc.GetByType("resolved")
	if err != nil {
		return err
	}

	now := time.Now()

	existingComment.CommentStatusID = commentStatus.ID
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now
	err = pvccs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
