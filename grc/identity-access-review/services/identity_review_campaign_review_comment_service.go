package services

import (
	"errors"
	"time"

	identityAccessReviewDTOs "jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/repositories"
	logServices "jethur/grc/logs/services"
	masterService "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IdentityReviewCampaignReviewCommentServiceInterface interface {
	GetAll(campaignID uint, currentUserID uint, params *identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentListQueryParamsDTO) (*identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentListResponseDTO, error)
	GetByID(campaignID uint, commentID uint) (*models.IdentityReviewCampaignReviewComment, error)
	GetFullDetailsByID(campaignID uint, commentID uint) (*identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentDetailDTO, error)
	Create(campaignID uint, currentUserID uint, dto *identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentFormDTO) (uint, error)
	Update(campaignID uint, id uint, currentUserID uint, dto *identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentFormDTO) error
	Delete(campaignID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(campaignID uint, commentID uint, currentUserID uint) error
}

type IdentityReviewCampaignReviewCommentService struct {
	repo                      repositories.IdentityReviewCampaignReviewCommentRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	identityReviewCampaignSvc IdentityReviewCampaignServiceInterface
	commentStatusSvc          masterService.CommentStatusServiceInterface
}

func NewIdentityReviewCampaignReviewCommentService(
	repo repositories.IdentityReviewCampaignReviewCommentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	identityReviewCampaignSvc IdentityReviewCampaignServiceInterface,
	commentStatusSvc masterService.CommentStatusServiceInterface,
) IdentityReviewCampaignReviewCommentServiceInterface {
	return &IdentityReviewCampaignReviewCommentService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		identityReviewCampaignSvc: identityReviewCampaignSvc,
		commentStatusSvc:          commentStatusSvc,
	}
}

func (ircrcs *IdentityReviewCampaignReviewCommentService) GetAll(campaignID uint, currentUserID uint, params *identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentListQueryParamsDTO) (*identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentListResponseDTO, error) {
	if _, err := ircrcs.identityReviewCampaignSvc.GetByID(campaignID); err != nil {
		return nil, err
	}

	comments, err := ircrcs.repo.GetAll(campaignID, currentUserID, params)
	if err != nil {
		return nil, err
	}

	dtoMap := make(map[uint]*identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentListDTO)
	var roots []*identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentListDTO

	for _, comment := range comments.Comments {
		dto := &identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentListDTO{
			ID:                                    comment.ID,
			Comment:                               comment.Comment,
			CommentStatusTitleEn:                  comment.CommentStatusTitleEn,
			CommentStatusTitleAr:                  comment.CommentStatusTitleAr,
			CommentStatusType:                     comment.CommentStatusType,
			CommentStatusColorCode:                comment.CommentStatusColorCode,
			IdentityReviewCampaignID:              comment.IdentityReviewCampaignID,
			Level:                                 comment.Level,
			CreatedBYID:                           comment.CreatedBYID,
			CreatedUserFirstNameEn:                comment.CreatedUserFirstNameEn,
			CreatedUserFirstNameAr:                comment.CreatedUserFirstNameAr,
			CreatedUserLastNameEn:                 comment.CreatedUserLastNameEn,
			CreatedUserLastNameAr:                 comment.CreatedUserLastNameAr,
			CreatedUserInitial:                    comment.CreatedUserInitial,
			CreatedUserColor:                      comment.CreatedUserColor,
			CreatedUserEmail:                      comment.CreatedUserEmail,
			CreatedUserProfileImagePath:           comment.CreatedUserProfileImagePath,
			CreatedAt:                             comment.CreatedAt,
			IdentityReviewCampaignReviewCommentID: comment.IdentityReviewCampaignReviewCommentID,
		}
		dtoMap[comment.ID] = dto
	}

	for _, comment := range comments.Comments {
		dto := dtoMap[comment.ID]
		if comment.IdentityReviewCampaignReviewCommentID != nil {
			if parent, ok := dtoMap[*comment.IdentityReviewCampaignReviewCommentID]; ok {
				parent.Children = append(parent.Children, dto)
				continue
			}
		}
		roots = append(roots, dto)
	}

	return &identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentListResponseDTO{
		Comments:      roots,
		OpenCount:     comments.OpenCount,
		ResolvedCount: comments.ResolvedCount,
	}, nil
}

func (ircrcs *IdentityReviewCampaignReviewCommentService) GetByID(campaignID uint, commentID uint) (*models.IdentityReviewCampaignReviewComment, error) {
	if _, err := ircrcs.identityReviewCampaignSvc.GetByID(campaignID); err != nil {
		return nil, err
	}

	comment, err := ircrcs.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, commentID)
		}
		return nil, err
	}
	return comment, nil
}

func (ircrcs *IdentityReviewCampaignReviewCommentService) GetFullDetailsByID(campaignID uint, commentID uint) (*identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentDetailDTO, error) {
	if _, err := ircrcs.identityReviewCampaignSvc.GetByID(campaignID); err != nil {
		return nil, err
	}

	comment, err := ircrcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, commentID)
		}
		return nil, err
	}

	var commentDetail identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentDetailDTO
	if err := mapstructure.Decode(comment, &commentDetail); err != nil {
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt
	return &commentDetail, nil
}

func (ircrcs *IdentityReviewCampaignReviewCommentService) Create(campaignID uint, currentUserID uint, commentFormDTO *identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentFormDTO) (id uint, err error) {
	tx := ircrcs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	campaign, err := ircrcs.identityReviewCampaignSvc.GetByID(campaignID, tx)
	if err != nil {
		return
	}

	if campaign.NextReviewLevel == nil {
		err = utility.NewValidationError(errors.New("workflow not active"), "review_comment_workflow_not_active", "identity review campaign")
		return
	}

	commentStatus, err := ircrcs.commentStatusSvc.GetByType("open")
	if err != nil {
		return
	}

	comment := &models.IdentityReviewCampaignReviewComment{
		IdentityReviewCampaignID:              campaignID,
		Comment:                               commentFormDTO.Comment,
		CommentStatusID:                       commentStatus.ID,
		Level:                                 *campaign.NextReviewLevel,
		IdentityReviewCampaignReviewCommentID: commentFormDTO.ParentCommentID,
		CreatedBYID:                           currentUserID,
	}

	id, err = ircrcs.repo.Create(comment, tx)
	if err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		existingComment, logErr := ircrcs.GetFullDetailsByID(campaignID, id)
		if logErr != nil {
			utility.LogError("Failed to fetch comment", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}

		logErr = ircrcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IdentityReviewCampaignReviewComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"identity_access_review", "review_campaign", "review_comment", nil,
		)
		if logErr != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":   models.CreateAction,
				"recordID": id,
				"error":    logErr.Error(),
			})
		}
	}()

	return
}

func (ircrcs *IdentityReviewCampaignReviewCommentService) Update(campaignID uint, id uint, currentUserID uint, commentFormDTO *identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentFormDTO) (err error) {
	tx := ircrcs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	if _, err = ircrcs.identityReviewCampaignSvc.GetByID(campaignID, tx); err != nil {
		return
	}

	existingComment, err := ircrcs.repo.GetByID(id, tx)
	if err != nil {
		return
	}

	oldComment, err := ircrcs.repo.GetFullDetailsByID(id, tx)
	if err != nil {
		return
	}

	existingComment.Comment = commentFormDTO.Comment
	if err = ircrcs.repo.Update(existingComment, tx); err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	go func() {
		newComment, logErr := ircrcs.repo.GetFullDetailsByID(id)
		if logErr != nil {
			return
		}

		changes := utility.CompareFieldChanges(
			map[string]interface{}{"comment": oldComment.Comment},
			map[string]interface{}{"comment": newComment.Comment},
		)

		if len(changes) > 0 {
			logErr := ircrcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.IdentityReviewCampaignReviewComment{}.TableName(),
				id,
				changes,
				currentUserID,
				"identity_access_review", "review_campaign", "review_comment", nil,
			)
			if logErr != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"recordID": id,
					"error":    logErr.Error(),
				})
			}
		}
	}()

	return
}

func (ircrcs *IdentityReviewCampaignReviewCommentService) Delete(campaignID uint, ids []uint, currentUserID uint) (failedIDs []uint, err error) {
	tx := ircrcs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil || len(failedIDs) > 0 {
			tx.Rollback()
			return
		}
	}()

	if _, err = ircrcs.identityReviewCampaignSvc.GetByID(campaignID, tx); err != nil {
		return
	}

	var successfulIDs []uint
	var deletedComments []map[string]interface{}

	for _, id := range ids {
		if _, getErr := ircrcs.repo.GetByID(id, tx); getErr != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":       id,
			"campaign": campaignID,
		})

		if deleteErr := ircrcs.repo.Delete(id, tx); deleteErr != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(failedIDs) > 0 {
		return
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	if len(successfulIDs) > 0 {
		go func() {
			logErr := ircrcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IdentityReviewCampaignReviewComment{}.TableName(),
				successfulIDs,
				deletedComments,
				currentUserID,
				nil,
			)
			if logErr != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{
					"error": logErr.Error(),
				})
			}
		}()
	}

	return
}

func (ircrcs *IdentityReviewCampaignReviewCommentService) Resolve(campaignID uint, commentID uint, currentUserID uint) (err error) {
	tx := ircrcs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	if _, err = ircrcs.identityReviewCampaignSvc.GetByID(campaignID, tx); err != nil {
		return
	}

	existingComment, err := ircrcs.repo.GetByID(commentID, tx)
	if err != nil {
		return
	}

	commentStatus, err := ircrcs.commentStatusSvc.GetByType("resolved")
	if err != nil {
		return
	}

	now := time.Now()
	existingComment.CommentStatusID = commentStatus.ID
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now

	if err = ircrcs.repo.Update(existingComment, tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return
}
