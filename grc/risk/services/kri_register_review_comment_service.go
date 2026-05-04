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

type KRIRegisterReviewCommentServiceInterface interface {
	GetAll(kriRegisterID uint, currentUserID uint, params *riskDTOs.KRIRegisterReviewCommentListQueryParamsDTO) (*riskDTOs.KRIRegisterReviewCommentListResponseDTO, error)
	GetByID(kriRegisterID uint, commentID uint) (*models.KRIRegisterReviewComment, error)
	GetFullDetailsByID(kriRegisterID uint, commentID uint) (*riskDTOs.KRIRegisterReviewCommentDetailDTO, error)
	Create(kriRegisterID uint, currentUserID uint, dto *riskDTOs.KRIRegisterReviewCommentFormDTO) (uint, error)
	Update(kriRegisterID uint, id uint, currentUserID uint, dto *riskDTOs.KRIRegisterReviewCommentFormDTO) error
	Delete(kriRegisterID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(kriRegisterID uint, commentID uint, currentUserID uint) error
}

type KRIRegisterReviewCommentService struct {
	repo             repositories.KRIRegisterReviewCommentRepositoryInterface
	userSvc          organizationServices.UserServiceInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	kriRegisterSvc   KRIRegisterServiceInterface
	commentStatusSvc masterService.CommentStatusServiceInterface
}

func NewKRIRegisterReviewCommentService(
	repo repositories.KRIRegisterReviewCommentRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	kriRegisterSvc KRIRegisterServiceInterface,
	commentStatusSvc masterService.CommentStatusServiceInterface,
) KRIRegisterReviewCommentServiceInterface {
	return &KRIRegisterReviewCommentService{
		repo:             repo,
		userSvc:          userSvc,
		activityLogSvc:   activityLogSvc,
		kriRegisterSvc:   kriRegisterSvc,
		commentStatusSvc: commentStatusSvc,
	}
}

func (krcs *KRIRegisterReviewCommentService) GetAll(kriRegisterID uint, currentUserID uint, params *riskDTOs.KRIRegisterReviewCommentListQueryParamsDTO) (*riskDTOs.KRIRegisterReviewCommentListResponseDTO, error) {
	_, err := krcs.kriRegisterSvc.GetFullDetailsByID(kriRegisterID)
	if err != nil {
		return nil, err
	}

	comments, err := krcs.repo.GetAll(kriRegisterID, currentUserID, params)
	if err != nil {
		return nil, err
	}
	dtoMap := make(map[uint]*riskDTOs.KRIRegisterReviewCommentListDTO)
	var roots []*riskDTOs.KRIRegisterReviewCommentListDTO

	for _, comment := range comments.Comments {
		dto := &riskDTOs.KRIRegisterReviewCommentListDTO{
			ID:                          comment.ID,
			Comment:                     comment.Comment,
			IsEdited:                    comment.IsEdited,
			CommentStatusTitleEn:        comment.CommentStatusTitleEn,
			CommentStatusTitleAr:        comment.CommentStatusTitleAr,
			CommentStatusType:           comment.CommentStatusType,
			CommentStatusColorCode:      comment.CommentStatusColorCode,
			KRIRegisterID:               comment.KRIRegisterID,
			Level:                       comment.Level,
			CreatedUserID:               comment.CreatedUserID,
			CreatedUserFirstNameEn:      comment.CreatedUserFirstNameEn,
			CreatedUserFirstNameAr:      comment.CreatedUserFirstNameAr,
			CreatedUserLastNameEn:       comment.CreatedUserLastNameEn,
			CreatedUserLastNameAr:       comment.CreatedUserLastNameAr,
			CreatedUserInitial:          comment.CreatedUserInitial,
			CreatedUserColor:            comment.CreatedUserColor,
			CreatedUserEmail:            comment.CreatedUserEmail,
			CreatedUserProfileImagePath: comment.CreatedUserProfileImagePath,
			CreatedAt:                   comment.CreatedAt,
			ParentCommentID:             comment.ParentCommentID,
		}
		dtoMap[comment.ID] = dto
	}

	for _, comment := range comments.Comments {
		dto := dtoMap[comment.ID]
		if comment.ParentCommentID != nil {
			if parent, ok := dtoMap[*comment.ParentCommentID]; ok {
				parent.Children = append(parent.Children, dto)
			}
		} else {
			roots = append(roots, dto)
		}
	}

	response := &riskDTOs.KRIRegisterReviewCommentListResponseDTO{
		Comments:      roots,
		OpenCount:     comments.OpenCount,
		ResolvedCount: comments.ResolvedCount,
	}

	return response, nil
}

func (krcs *KRIRegisterReviewCommentService) GetByID(kriRegisterID uint, commentID uint) (*models.KRIRegisterReviewComment, error) {
	if _, err := krcs.kriRegisterSvc.GetByID(kriRegisterID); err != nil {
		return nil, err
	}

	comment, err := krcs.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, kriRegisterID)
		}
		return nil, err
	}
	return comment, nil
}

func (krcs *KRIRegisterReviewCommentService) GetFullDetailsByID(kriRegisterID uint, commentID uint) (*riskDTOs.KRIRegisterReviewCommentDetailDTO, error) {
	if _, err := krcs.kriRegisterSvc.GetByID(kriRegisterID); err != nil {
		fmt.Printf("KRI Register not found for kriRegisterID %d", kriRegisterID)
		return nil, err
	}

	var commentDetail riskDTOs.KRIRegisterReviewCommentDetailDTO

	comment, err := krcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, kriRegisterID)
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

func (krcs *KRIRegisterReviewCommentService) Create(kriRegisterID uint, currentUserID uint, commentFormDTO *riskDTOs.KRIRegisterReviewCommentFormDTO) (uint, error) {
	tx := krcs.repo.BeginTransaction()

	kriRegister, err := krcs.kriRegisterSvc.GetByID(kriRegisterID)
	if err != nil {
		return 0, err
	}

	commentStatus, err := krcs.commentStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}

	comment := &models.KRIRegisterReviewComment{
		KRIRegisterID:   kriRegisterID,
		Comment:         commentFormDTO.Comment,
		CommentStatusID: commentStatus.ID,
		Level:           kriRegister.NextReviewLevel,
		ParentCommentID: commentFormDTO.ParentCommentID,
		CreatedBYID:     currentUserID,
	}

	id, err := krcs.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingComment, err := krcs.GetFullDetailsByID(kriRegisterID, id)
		if err != nil {
			utility.LogError("Failed to fetch comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = krcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.KRIRegisterReviewComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.KRIRegisterReviewComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}

func (krcs *KRIRegisterReviewCommentService) Update(kriRegisterID uint, id uint, currentUserID uint, commentFormDTO *riskDTOs.KRIRegisterReviewCommentFormDTO) error {
	tx := krcs.repo.BeginTransaction()

	if _, err := krcs.kriRegisterSvc.GetByID(kriRegisterID); err != nil {
		return err
	}

	existingComment, err := krcs.GetByID(kriRegisterID, id)
	if err != nil {
		return err
	}

	oldComment, err := krcs.GetFullDetailsByID(kriRegisterID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = commentFormDTO.Comment
	existingComment.IsEdited = true
	now := time.Now()
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now
	err = krcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingComment, err := krcs.GetFullDetailsByID(kriRegisterID, id)
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
			err := krcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.KRIRegisterReviewComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				"", "", "", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.KRIRegisterReviewComment{}.TableName(),
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

func (krcs *KRIRegisterReviewCommentService) Delete(kriRegisterID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := krcs.repo.BeginTransaction()
	if _, err := krcs.kriRegisterSvc.GetByID(kriRegisterID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		if _, err := krcs.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":  id,
			"kri": kriRegisterID,
		})

		err := krcs.repo.Delete(id, tx)
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
			err := krcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.KRIRegisterReviewComment{}.TableName(),
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

func (krcs *KRIRegisterReviewCommentService) Resolve(kriRegisterID uint, commentID uint, currentUserID uint) error {
	tx := krcs.repo.BeginTransaction()

	if _, err := krcs.kriRegisterSvc.GetByID(kriRegisterID); err != nil {
		return err
	}

	existingComment, err := krcs.GetByID(kriRegisterID, commentID)
	if err != nil {
		return err
	}

	commentStatus, err := krcs.commentStatusSvc.GetByType("resolved")
	if err != nil {
		return err
	}

	now := time.Now()

	existingComment.CommentStatusID = commentStatus.ID
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now
	err = krcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
