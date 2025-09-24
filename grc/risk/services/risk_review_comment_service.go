package services

import (
	"errors"
	"fmt"
	"time"

	logServices "jethur/grc/logs/services"
	masterService "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskReviewCommentServiceInterface interface {
	GetAll(riskID uint, currentUserID uint, params *riskDTOs.RiskReviewCommentListQueryParamsDTO) (*riskDTOs.RiskReviewCommentListResponseDTO, error)
	GetByID(riskID uint, commentID uint) (*models.RiskReviewComment, error)
	GetFullDetailsByID(riskID uint, commentID uint) (*riskDTOs.RiskReviewCommentDetailDTO, error)
	Create(riskID uint, currentUserID uint, dto *riskDTOs.RiskReviewCommentFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, dto *riskDTOs.RiskReviewCommentFormDTO) error
	Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(riskID uint, commentId uint, currentUserID uint) error
}

type RiskReviewCommentService struct {
	repo             repositories.RiskReviewCommentRepositoryInterface
	userSvc          organizationServices.UserServiceInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	riskRegisterSvc  RiskRegisterServiceInterface
	commentStatusSvc masterService.CommentStatusServiceInterface
}

func NewRiskReviewCommentService(repo repositories.RiskReviewCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	commentStatusSvc masterService.CommentStatusServiceInterface,
) RiskReviewCommentServiceInterface {
	return &RiskReviewCommentService{
		repo:             repo,
		userSvc:          userSvc,
		activityLogSvc:   activityLogSvc,
		riskRegisterSvc:  riskRegisterSvc,
		commentStatusSvc: commentStatusSvc,
	}
}

func (pvccs *RiskReviewCommentService) GetAll(riskID uint, currentUserID uint, params *riskDTOs.RiskReviewCommentListQueryParamsDTO) (*riskDTOs.RiskReviewCommentListResponseDTO, error) {
	_, err := pvccs.riskRegisterSvc.GetFullDetailsByID(riskID)
	if err != nil {
		return nil, err
	}

	comments, err := pvccs.repo.GetAll(riskID, currentUserID, params)
	dtoMap := make(map[uint]*riskDTOs.RiskReviewCommentListDTO)
	var roots []*riskDTOs.RiskReviewCommentListDTO

	for _, comment := range comments.Comments {
		dto := &riskDTOs.RiskReviewCommentListDTO{
			ID:                          comment.ID,
			Comment:                     comment.Comment,
			CommentStatusTitleEn:        comment.CommentStatusTitleEn,
			CommentStatusTitleAr:        comment.CommentStatusTitleAr,
			CommentStatusType:           comment.CommentStatusType,
			CommentStatusColorCode:      comment.CommentStatusColorCode,
			RiskRegisterID:              comment.RiskRegisterID,
			Level:                       comment.Level,
			CreatedBYID:                 comment.CreatedBYID,
			CreatedUserFirstNameEn:      comment.CreatedUserFirstNameEn,
			CreatedUserFirstNameAr:      comment.CreatedUserFirstNameAr,
			CreatedUserLastNameEn:       comment.CreatedUserLastNameEn,
			CreatedUserLastNameAr:       comment.CreatedUserLastNameAr,
			CreatedUserInitial:          comment.CreatedUserInitial,
			CreatedUserColor:            comment.CreatedUserColor,
			CreatedUserEmail:            comment.CreatedUserEmail,
			CreatedUserProfileImagePath: comment.CreatedUserProfileImagePath,
			CreatedAt:                   comment.CreatedAt,
			RiskReviewCommentID:         comment.RiskReviewCommentID,
		}
		dtoMap[comment.ID] = dto
	}

	for _, comment := range comments.Comments {
		dto := dtoMap[comment.ID]
		if comment.RiskReviewCommentID != nil {
			if parent, ok := dtoMap[*comment.RiskReviewCommentID]; ok {
				parent.Children = append(parent.Children, dto)
			}
		} else {
			roots = append(roots, dto)
		}
	}

	response := &riskDTOs.RiskReviewCommentListResponseDTO{
		Comments:      roots,
		OpenCount:     comments.OpenCount,
		ResolvedCount: comments.ResolvedCount,
	}

	return response, nil
}

func (pvccs *RiskReviewCommentService) GetByID(riskID uint, commentID uint) (*models.RiskReviewComment, error) {
	if _, err := pvccs.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	comment, err := pvccs.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, riskID)
		}
		return nil, err
	}
	return comment, nil
}

func (pvccs *RiskReviewCommentService) GetFullDetailsByID(riskID uint, commentID uint) (*riskDTOs.RiskReviewCommentDetailDTO, error) {
	if _, err := pvccs.riskRegisterSvc.GetByID(riskID); err != nil {
		fmt.Printf("Risk Register not found for riskID %d", riskID)
		return nil, err
	}

	var commentDetail riskDTOs.RiskReviewCommentDetailDTO

	comment, err := pvccs.repo.GetFullDetailsByID(riskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, riskID)
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

func (pvccs *RiskReviewCommentService) Create(riskID uint, currentUserID uint, commentFormDTO *riskDTOs.RiskReviewCommentFormDTO) (uint, error) {
	tx := pvccs.repo.BeginTransaction()

	_, err := pvccs.riskRegisterSvc.GetFullDetailsByID(riskID)
	if err != nil {
		return 0, err
	}

	commentStatus, err := pvccs.commentStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}

	comment := &models.RiskReviewComment{
		RiskRegisterID:  riskID,
		Comment:         commentFormDTO.Comment,
		CommentStatusID: commentStatus.ID,
		// Level:               *riskRegister.NextReviewLevel,
		RiskReviewCommentID: commentFormDTO.ParentCommentID,
		CreatedBYID:         currentUserID,
	}

	id, err := pvccs.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingComment, err := pvccs.GetFullDetailsByID(riskID, id)
		if err != nil {
			utility.LogError("Failed to fetch comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = pvccs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskReviewComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskReviewComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

	}()
	return id, nil
}

func (pvccs *RiskReviewCommentService) Update(riskID uint, id uint, currentUserID uint, commentFormDTO *riskDTOs.RiskReviewCommentFormDTO) error {
	tx := pvccs.repo.BeginTransaction()

	if _, err := pvccs.riskRegisterSvc.GetByID(riskID); err != nil {
		return err
	}

	existingComment, err := pvccs.GetByID(riskID, id)
	if err != nil {
		return err
	}

	oldComment, err := pvccs.GetFullDetailsByID(riskID, id)
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
		existingComment, err := pvccs.GetFullDetailsByID(riskID, id)
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
				models.RiskReviewComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.RiskReviewComment{}.TableName(),
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

func (pvccs *RiskReviewCommentService) Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := pvccs.repo.BeginTransaction()
	if _, err := pvccs.riskRegisterSvc.GetByID(riskID); err != nil {
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
			"id":   id,
			"risk": riskID,
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
				models.RiskReviewComment{}.TableName(),
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

func (pvccs *RiskReviewCommentService) Resolve(riskID uint, commentId uint, currentUserID uint) error {

	tx := pvccs.repo.BeginTransaction()

	if _, err := pvccs.riskRegisterSvc.GetByID(riskID); err != nil {
		return err
	}

	existingComment, err := pvccs.GetByID(riskID, commentId)
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
