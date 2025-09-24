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

type PolicyClauseCommentServiceInterface interface {
	GetAllDocumentComments(documentID uint, currentUserID uint, params *policyDTOs.PolicyClauseCommentListQueryParamsDTO) (*policyDTOs.PolicyClauseCommentListResponseDTO, error)
	GetByID(documentID uint, id uint) (*models.DocumentVersionClauseComment, error)
	GetFullDetailsByID(documentID uint, id uint, tx ...utility.DatabaseTransaction) (*policyDTOs.PolicyClauseCommentDetailDTO, error)
	Create(documentID uint, currentUserID uint, dto *policyDTOs.PolicyClauseCommentFormDTO) (uint, error)
	Update(documentID uint, id uint, currentUserID uint, dto *policyDTOs.PolicyClauseCommentFormDTO) error
	Delete(documentID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(documentID uint, commentId uint, currentUserID uint) error
}

type PolicyClauseCommentService struct {
	repo              repositories.PolicyClauseCommentRepositoryInterface
	userSvc           organizationServices.UserServiceInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
	policyRegisterSvc PolicyRegisterServiceInterface
	commentStatusSvc  masterService.CommentStatusServiceInterface
}

func NewPolicyClauseCommentService(repo repositories.PolicyClauseCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	policyRegisterSvc PolicyRegisterServiceInterface,
	commentStatusSvc masterService.CommentStatusServiceInterface,
) PolicyClauseCommentServiceInterface {
	return &PolicyClauseCommentService{
		repo:              repo,
		userSvc:           userSvc,
		activityLogSvc:    activityLogSvc,
		policyRegisterSvc: policyRegisterSvc,
		commentStatusSvc:  commentStatusSvc,
	}
}

func (pvccs *PolicyClauseCommentService) GetAllDocumentComments(documentID uint, currentUserID uint, params *policyDTOs.PolicyClauseCommentListQueryParamsDTO) (*policyDTOs.PolicyClauseCommentListResponseDTO, error) {
	policyRegister, err := pvccs.policyRegisterSvc.GetFullDetailsByID(documentID, &currentUserID)
	if err != nil {
		return nil, err
	}

	comments, err := pvccs.repo.GetAllDocumentComments(documentID, policyRegister.LatestDocumentVersion.ID, currentUserID, params)
	dtoMap := make(map[uint]*policyDTOs.PolicyClauseCommentListDTO)
	var roots []*policyDTOs.PolicyClauseCommentListDTO

	for _, comment := range comments.Comments {
		dto := &policyDTOs.PolicyClauseCommentListDTO{
			ID:                             comment.ID,
			Comment:                        comment.Comment,
			CommentStatusTitleEn:           comment.CommentStatusTitleEn,
			CommentStatusTitleAr:           comment.CommentStatusTitleAr,
			CommentStatusType:              comment.CommentStatusType,
			CommentStatusColorCode:         comment.CommentStatusColorCode,
			DocumentVersionClauseID:        comment.DocumentVersionClauseID,
			DocumentVersionClauseNumber:    comment.DocumentVersionClauseNumber,
			DocumentVersionID:              comment.DocumentVersionID,
			Level:                          comment.Level,
			CreatedBYID:                    comment.CreatedBYID,
			CreatedUserFirstNameEn:         comment.CreatedUserFirstNameEn,
			CreatedUserFirstNameAr:         comment.CreatedUserFirstNameAr,
			CreatedUserLastNameEn:          comment.CreatedUserLastNameEn,
			CreatedUserLastNameAr:          comment.CreatedUserLastNameAr,
			CreatedUserInitial:             comment.CreatedUserInitial,
			CreatedUserColor:               comment.CreatedUserColor,
			CreatedUserEmail:               comment.CreatedUserEmail,
			CreatedUserProfileImagePath:    comment.CreatedUserProfileImagePath,
			CreatedAt:                      comment.CreatedAt,
			DocumentVersionClauseCommentID: comment.DocumentVersionClauseCommentID,
		}
		dtoMap[comment.ID] = dto
	}

	for _, comment := range comments.Comments {
		dto := dtoMap[comment.ID]
		if comment.DocumentVersionClauseCommentID != nil {
			if parent, ok := dtoMap[*comment.DocumentVersionClauseCommentID]; ok {
				parent.Children = append(parent.Children, dto)
			}
		} else {
			roots = append(roots, dto)
		}
	}

	response := &policyDTOs.PolicyClauseCommentListResponseDTO{
		Comments:      roots,
		OpenCount:     comments.OpenCount,
		ResolvedCount: comments.ResolvedCount,
	}

	return response, nil
}

func (pvccs *PolicyClauseCommentService) GetByID(documentID uint, id uint) (*models.DocumentVersionClauseComment, error) {
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

func (pvccs *PolicyClauseCommentService) GetFullDetailsByID(documentID uint, id uint, tx ...utility.DatabaseTransaction) (*policyDTOs.PolicyClauseCommentDetailDTO, error) {
	if _, err := pvccs.policyRegisterSvc.GetByID(documentID); err != nil {
		fmt.Printf("Policy not found for documentID %d", documentID)
		return nil, err
	}

	var commentDetail policyDTOs.PolicyClauseCommentDetailDTO

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

func (pvccs *PolicyClauseCommentService) Create(documentID uint, currentUserID uint, commentFormDTO *policyDTOs.PolicyClauseCommentFormDTO) (uint, error) {
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
		DocumentVersionClauseID:        commentFormDTO.ClauseID,
		DocumentVersionID:              &policyRegister.LatestDocumentVersion.ID,
		Comment:                        commentFormDTO.Comment,
		CommentStatusID:                commentStatus.ID,
		Level:                          *policyRegister.LatestDocumentVersion.NextReviewLevel,
		DocumentVersionClauseCommentID: commentFormDTO.ParentCommentID,
		CreatedBYID:                    currentUserID,
	}

	id, err := pvccs.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingComment, err := pvccs.GetFullDetailsByID(documentID, id)
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

func (pvccs *PolicyClauseCommentService) Update(documentID uint, id uint, currentUserID uint, commentFormDTO *policyDTOs.PolicyClauseCommentFormDTO) error {
	tx := pvccs.repo.BeginTransaction()

	if _, err := pvccs.policyRegisterSvc.GetByID(documentID); err != nil {
		return err
	}

	existingComment, err := pvccs.GetByID(documentID, id)
	if err != nil {
		return err
	}

	oldComment, err := pvccs.GetFullDetailsByID(documentID, id)
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
		existingComment, err := pvccs.GetFullDetailsByID(documentID, id)
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

func (pvccs *PolicyClauseCommentService) Delete(documentID uint, ids []uint, currentUserID uint) ([]uint, error) {
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
			"id":       id,
			"document": documentID,
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

func (pvccs *PolicyClauseCommentService) Resolve(documentID uint, commentId uint, currentUserID uint) error {

	tx := pvccs.repo.BeginTransaction()

	if _, err := pvccs.policyRegisterSvc.GetByID(documentID); err != nil {
		return err
	}

	existingComment, err := pvccs.GetByID(documentID, commentId)
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
