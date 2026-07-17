package services

import (
	"errors"
	"time"

	assetRiskDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	masterService "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AssetRiskReviewCommentServiceInterface interface {
	GetAll(riskID uint, currentUserID uint, params *assetRiskDTOs.AssetRiskReviewCommentListQueryParamsDTO) (*assetRiskDTOs.AssetRiskReviewCommentListResponseDTO, error)
	GetByID(riskID uint, commentID uint) (*models.RiskReviewComment, error)
	GetFullDetailsByID(riskID uint, commentID uint) (*assetRiskDTOs.AssetRiskReviewCommentDetailDTO, error)
	Create(riskID uint, currentUserID uint, dto *assetRiskDTOs.AssetRiskReviewCommentFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, dto *assetRiskDTOs.AssetRiskReviewCommentFormDTO) error
	Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(riskID uint, commentID uint, currentUserID uint) error
}

type AssetRiskReviewCommentService struct {
	repo             repositories.AssetRiskReviewCommentRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	riskSvc          AssetRiskRegisterServiceInterface
	commentStatusSvc masterService.CommentStatusServiceInterface
}

func NewAssetRiskReviewCommentService(
	repo repositories.AssetRiskReviewCommentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskSvc AssetRiskRegisterServiceInterface,
	commentStatusSvc masterService.CommentStatusServiceInterface,
) AssetRiskReviewCommentServiceInterface {
	return &AssetRiskReviewCommentService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		riskSvc:          riskSvc,
		commentStatusSvc: commentStatusSvc,
	}
}

func (arrcs *AssetRiskReviewCommentService) GetAll(riskID uint, currentUserID uint, params *assetRiskDTOs.AssetRiskReviewCommentListQueryParamsDTO) (*assetRiskDTOs.AssetRiskReviewCommentListResponseDTO, error) {
	if _, err := arrcs.riskSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	comments, err := arrcs.repo.GetAll(riskID, currentUserID, params)
	if err != nil {
		return nil, err
	}

	dtoMap := make(map[uint]*assetRiskDTOs.AssetRiskReviewCommentListDTO)
	var roots []*assetRiskDTOs.AssetRiskReviewCommentListDTO

	for _, comment := range comments.Comments {
		dto := &assetRiskDTOs.AssetRiskReviewCommentListDTO{
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
				continue
			}
		}
		roots = append(roots, dto)
	}

	return &assetRiskDTOs.AssetRiskReviewCommentListResponseDTO{
		Comments:      roots,
		OpenCount:     comments.OpenCount,
		ResolvedCount: comments.ResolvedCount,
	}, nil
}

func (arrcs *AssetRiskReviewCommentService) GetByID(riskID uint, commentID uint) (*models.RiskReviewComment, error) {
	if _, err := arrcs.riskSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	comment, err := arrcs.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, commentID)
		}
		return nil, err
	}
	return comment, nil
}

func (arrcs *AssetRiskReviewCommentService) GetFullDetailsByID(riskID uint, commentID uint) (*assetRiskDTOs.AssetRiskReviewCommentDetailDTO, error) {
	if _, err := arrcs.riskSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	comment, err := arrcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, commentID)
		}
		return nil, err
	}

	var commentDetail assetRiskDTOs.AssetRiskReviewCommentDetailDTO
	if err := mapstructure.Decode(comment, &commentDetail); err != nil {
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt
	return &commentDetail, nil
}

func (arrcs *AssetRiskReviewCommentService) Create(riskID uint, currentUserID uint, commentFormDTO *assetRiskDTOs.AssetRiskReviewCommentFormDTO) (id uint, err error) {
	tx := arrcs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	riskRegister, err := arrcs.riskSvc.GetByID(riskID, tx)
	if err != nil {
		return
	}

	if riskRegister.NextReviewLevel == nil {
		err = utility.NewValidationError(errors.New("workflow not active"), "review_comment_workflow_not_active", "risk acceptance")
		return
	}

	commentStatus, err := arrcs.commentStatusSvc.GetByType("open")
	if err != nil {
		return
	}

	comment := &models.RiskReviewComment{
		RiskRegisterID:      riskID,
		Comment:             commentFormDTO.Comment,
		CommentStatusID:     commentStatus.ID,
		Level:               *riskRegister.NextReviewLevel,
		RiskReviewCommentID: commentFormDTO.ParentCommentID,
		CreatedBYID:         currentUserID,
	}

	id, err = arrcs.repo.Create(comment, tx)
	if err != nil {
		return
	}

	go func() {
		existingComment, logErr := arrcs.GetFullDetailsByID(riskID, id)
		if logErr != nil {
			utility.LogError("Failed to fetch comment", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}

		logErr = arrcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskReviewComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"asset", "asset_risk", "risk_review_comment", nil,
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

func (arrcs *AssetRiskReviewCommentService) Update(riskID uint, id uint, currentUserID uint, commentFormDTO *assetRiskDTOs.AssetRiskReviewCommentFormDTO) (err error) {
	tx := arrcs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	if _, err = arrcs.riskSvc.GetByID(riskID, tx); err != nil {
		return
	}

	existingComment, err := arrcs.repo.GetByID(id, tx)
	if err != nil {
		return
	}

	oldComment, err := arrcs.repo.GetFullDetailsByID(id, tx)
	if err != nil {
		return
	}

	existingComment.Comment = commentFormDTO.Comment
	if err = arrcs.repo.Update(existingComment, tx); err != nil {
		return
	}

	go func() {
		newComment, logErr := arrcs.repo.GetFullDetailsByID(id)
		if logErr != nil {
			return
		}

		changes := utility.CompareFieldChanges(
			map[string]interface{}{"comment": oldComment.Comment},
			map[string]interface{}{"comment": newComment.Comment},
		)

		if len(changes) > 0 {
			logErr := arrcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskReviewComment{}.TableName(),
				id,
				changes,
				currentUserID,
				"asset", "asset_risk", "risk_review_comment", nil,
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

func (arrcs *AssetRiskReviewCommentService) Delete(riskID uint, ids []uint, currentUserID uint) (failedIDs []uint, err error) {
	tx := arrcs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil || len(failedIDs) > 0 {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	if _, err = arrcs.riskSvc.GetByID(riskID, tx); err != nil {
		return
	}

	var successfulIDs []uint
	var deletedComments []map[string]interface{}

	for _, id := range ids {
		if _, getErr := arrcs.repo.GetByID(id, tx); getErr != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":   id,
			"risk": riskID,
		})

		if deleteErr := arrcs.repo.Delete(id, tx); deleteErr != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(failedIDs) > 0 {
		return
	}

	if len(successfulIDs) > 0 {
		go func() {
			logErr := arrcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskReviewComment{}.TableName(),
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

func (arrcs *AssetRiskReviewCommentService) Resolve(riskID uint, commentID uint, currentUserID uint) (err error) {
	tx := arrcs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	if _, err = arrcs.riskSvc.GetByID(riskID, tx); err != nil {
		return
	}

	existingComment, err := arrcs.repo.GetByID(commentID, tx)
	if err != nil {
		return
	}

	commentStatus, err := arrcs.commentStatusSvc.GetByType("resolved")
	if err != nil {
		return
	}

	now := time.Now()
	existingComment.CommentStatusID = commentStatus.ID
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now

	err = arrcs.repo.Update(existingComment, tx)
	return
}
