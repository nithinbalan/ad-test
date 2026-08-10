package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IdentityReviewUserCommentServiceInterface interface {
	GetAll(identityReviewUserID uint) ([]*models.IdentityReviewUserComment, error)
	GetByID(identityReviewUserID uint, id uint) (*models.IdentityReviewUserComment, error)
	Create(identityReviewUserID uint, currentUserID uint, formDTO *dtos.IdentityReviewUserCommentFormDTO) (uint, error)
	CreateInitial(identityReviewUserID uint, accountDecisionID uint, currentUserID uint, comment string, tx utility.DatabaseTransaction) (uint, error)
	Update(identityReviewUserID uint, id uint, currentUserID uint, formDTO *dtos.IdentityReviewUserCommentFormDTO) error
	Delete(identityReviewUserID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type IdentityReviewUserCommentService struct {
	repo                  repositories.IdentityReviewUserCommentRepositoryInterface
	activityLogSvc        logServices.ActivityLogServiceInterface
	identityReviewUserSvc IdentityReviewUserServiceInterface
	accountDecisionSvc    AccountDecisionServiceInterface
}

func NewIdentityReviewUserCommentService(
	repo repositories.IdentityReviewUserCommentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	identityReviewUserSvc IdentityReviewUserServiceInterface,
	accountDecisionSvc AccountDecisionServiceInterface,
) IdentityReviewUserCommentServiceInterface {
	return &IdentityReviewUserCommentService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		identityReviewUserSvc: identityReviewUserSvc,
		accountDecisionSvc:    accountDecisionSvc,
	}
}

func (iruscs *IdentityReviewUserCommentService) GetAll(identityReviewUserID uint) ([]*models.IdentityReviewUserComment, error) {
	if _, err := iruscs.identityReviewUserSvc.GetByID(identityReviewUserID); err != nil {
		return nil, err
	}

	return iruscs.repo.GetAll(identityReviewUserID)
}

func (iruscs *IdentityReviewUserCommentService) GetByID(identityReviewUserID uint, id uint) (*models.IdentityReviewUserComment, error) {
	if _, err := iruscs.identityReviewUserSvc.GetByID(identityReviewUserID); err != nil {
		return nil, err
	}

	comment, err := iruscs.repo.GetByID(identityReviewUserID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return comment, nil
}

func (iruscs *IdentityReviewUserCommentService) validateAccountDecision(accountDecisionID *uint) error {
	if accountDecisionID == nil {
		return nil
	}

	if _, err := iruscs.accountDecisionSvc.GetByID(*accountDecisionID); err != nil {
		return utility.NewValidationError(
			fmt.Errorf("account decision does not exist"),
			"not_exist",
			"account_decision_id",
		)
	}

	return nil
}

func (iruscs *IdentityReviewUserCommentService) Create(identityReviewUserID uint, currentUserID uint, formDTO *dtos.IdentityReviewUserCommentFormDTO) (uint, error) {
	if _, err := iruscs.identityReviewUserSvc.GetByID(identityReviewUserID); err != nil {
		return 0, err
	}

	if err := iruscs.validateAccountDecision(formDTO.AccountDecisionID); err != nil {
		return 0, err
	}

	tx := iruscs.repo.BeginTransaction()

	comment := &models.IdentityReviewUserComment{
		IdentityReviewUserID: identityReviewUserID,
		AccountDecisionID:    formDTO.AccountDecisionID,
		Comment:              formDTO.Comment,
		CreatedBYID:          currentUserID,
	}

	id, err := iruscs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		if err := iruscs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IdentityReviewUserComment{}.TableName(),
			id,
			comment,
			currentUserID,
			"identity_access_review", "identity_review_user", "comment", nil,
		); err != nil {
			utility.LogError("Failed to log identity review user comment creation", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (iruscs *IdentityReviewUserCommentService) CreateInitial(identityReviewUserID uint, accountDecisionID uint, currentUserID uint, comment string, tx utility.DatabaseTransaction) (uint, error) {
	if strings.TrimSpace(comment) == "" {
		return 0, nil
	}

	decisionID := accountDecisionID

	initialComment := &models.IdentityReviewUserComment{
		IdentityReviewUserID: identityReviewUserID,
		AccountDecisionID:    &decisionID,
		Comment:              comment,
		CreatedBYID:          currentUserID,
	}

	return iruscs.repo.Create(initialComment, tx)
}

func (iruscs *IdentityReviewUserCommentService) Update(identityReviewUserID uint, id uint, currentUserID uint, formDTO *dtos.IdentityReviewUserCommentFormDTO) error {
	if err := iruscs.validateAccountDecision(formDTO.AccountDecisionID); err != nil {
		return err
	}

	existingComment, err := iruscs.GetByID(identityReviewUserID, id)
	if err != nil {
		return err
	}

	oldComment := *existingComment

	existingComment.AccountDecisionID = formDTO.AccountDecisionID
	existingComment.Comment = formDTO.Comment
	now := time.Now()
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now

	if err := iruscs.repo.Update(existingComment); err != nil {
		return err
	}

	go func() {
		oldData := map[string]interface{}{
			"account_decision_id": oldComment.AccountDecisionID,
			"comment":             oldComment.Comment,
		}

		newData := map[string]interface{}{
			"account_decision_id": existingComment.AccountDecisionID,
			"comment":             existingComment.Comment,
		}

		changedData := utility.CompareFieldChanges(oldData, newData)
		if len(changedData) == 0 {
			return
		}

		if err := iruscs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IdentityReviewUserComment{}.TableName(),
			id,
			changedData,
			currentUserID,
			"identity_access_review", "identity_review_user", "comment", nil,
		); err != nil {
			utility.LogError("Failed to log identity review user comment update", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (iruscs *IdentityReviewUserCommentService) Delete(identityReviewUserID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := iruscs.identityReviewUserSvc.GetByID(identityReviewUserID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		comment, err := iruscs.repo.GetByID(identityReviewUserID, id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		successfulIDs = append(successfulIDs, id)
		deletedComments = append(deletedComments, map[string]interface{}{
			"id":                      comment.ID,
			"identity_review_user_id": comment.IdentityReviewUserID,
			"account_decision_id":     comment.AccountDecisionID,
			"comment":                 comment.Comment,
		})
	}

	if len(failedIDs) > 0 {
		return failedIDs, nil
	}
	if len(successfulIDs) == 0 {
		return failedIDs, nil
	}

	tx := iruscs.repo.BeginTransaction()

	if err := iruscs.repo.Delete(identityReviewUserID, successfulIDs, tx); err != nil {
		tx.Rollback()
		return failedIDs, err
	}

	if err := tx.Commit(); err != nil {
		return failedIDs, err
	}

	go func() {
		if err := iruscs.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.IdentityReviewUserComment{}.TableName(),
			successfulIDs,
			deletedComments,
			currentUserID,
			nil,
		); err != nil {
			utility.LogError("Failed to log identity review user comment deletion", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return failedIDs, nil
}
