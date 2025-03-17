package services

import (
	"errors"
	"fmt"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AssetPolicyCommentServiceInterface interface {
	GetAll(assetID uint, assetPolicyID uint) ([]*dtos.AssetPolicyCommentListDTO, error)
	GetFullDetailsByID(assetID uint, assetPolicyID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.AssetPolicyCommentDetailDTO, error)
	Create(assetID uint, assetPolicyID uint, currentUserID uint, dto *dtos.AssetPolicyCommentFormDTO) (uint, error)
	Update(assetID uint, assetPolicyID uint, id uint, currentUserID uint, dto *dtos.AssetPolicyCommentFormDTO) error
	Delete(assetID uint, assetPolicyID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type AssetPolicyCommentService struct {
	repo             repositories.AssetPolicyCommentRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	assetRegisterSvc AssetRegisterServiceInterface
	assetPolicySvc   AssetPolicyServiceInterface
}

func NewAssetPolicyCommentService(
	repo repositories.AssetPolicyCommentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	assetRegisterSvc AssetRegisterServiceInterface,
	assetPolicySvc AssetPolicyServiceInterface,
) *AssetPolicyCommentService {
	return &AssetPolicyCommentService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		assetRegisterSvc: assetRegisterSvc,
		assetPolicySvc:   assetPolicySvc,
	}
}

func (apcs *AssetPolicyCommentService) GetAll(assetID uint, assetPolicyID uint) ([]*dtos.AssetPolicyCommentListDTO, error) {
	if _, err := apcs.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	if _, err := apcs.assetPolicySvc.GetByID(assetPolicyID); err != nil {
		return nil, err
	}
	return apcs.repo.GetAll(assetPolicyID)
}

func (apcs *AssetPolicyCommentService) GetFullDetailsByID(assetID uint, assetPolicyID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.AssetPolicyCommentDetailDTO, error) {
	var commentDetail dtos.AssetPolicyCommentDetailDTO

	// Validate asset and policy existence
	if _, err := apcs.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	if _, err := apcs.assetPolicySvc.GetByID(assetPolicyID); err != nil {
		return nil, err
	}

	// Fetch comment details
	comment, err := apcs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Convert to DTO
	err = mapstructure.Decode(comment, &commentDetail)
	if err != nil {
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt
	return &commentDetail, nil
}

func (apcs *AssetPolicyCommentService) Create(assetID uint, assetPolicyID uint, currentUserID uint, dto *dtos.AssetPolicyCommentFormDTO) (uint, error) {
	if _, err := apcs.assetRegisterSvc.GetByID(assetID); err != nil {
		return 0, err
	}
	if _, err := apcs.assetPolicySvc.GetByID(assetPolicyID); err != nil {
		return 0, err
	}
	tx := apcs.repo.BeginTransaction()

	comment := &models.AssetPolicyComment{
		AssetPolicyID: assetPolicyID,
		Comments:      dto.Comments,
		CreatedBYID:   currentUserID,
	}

	id, err := apcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	// **Activity Logging**
	go func() {
		existingComment, err := apcs.GetFullDetailsByID(assetID, assetPolicyID, id)
		if err != nil {
			utility.LogError("Failed to fetch asset policy comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = apcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetPolicyComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (apcs *AssetPolicyCommentService) Update(assetID uint, assetPolicyID uint, id uint, currentUserID uint, dto *dtos.AssetPolicyCommentFormDTO) error {
	if _, err := apcs.assetRegisterSvc.GetByID(assetID); err != nil {
		return err
	}
	if _, err := apcs.assetPolicySvc.GetByID(assetPolicyID); err != nil {
		return err
	}
	tx := apcs.repo.BeginTransaction()

	existingComment, err := apcs.repo.GetByID(id)
	if err != nil {
		return err
	}

	oldComment, err := apcs.GetFullDetailsByID(assetID, assetPolicyID, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	now := time.Now()
	existingComment.Comments = dto.Comments
	existingComment.UpdatedAt = &now
	existingComment.UpdatedBYID = &currentUserID

	err = apcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// **Activity Logging**
	go func() {
		updatedComment, err := apcs.GetFullDetailsByID(assetID, assetPolicyID, id)
		if err != nil {
			return
		}

		changedData := utility.CompareFieldChanges(
			map[string]interface{}{"comments": oldComment.Comments},
			map[string]interface{}{"comments": updatedComment.Comments},
		)

		if len(changedData) > 0 {
			err := apcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AssetPolicyComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
	return nil
}

func (apcs *AssetPolicyCommentService) Delete(assetID uint, assetPolicyID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := apcs.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	if _, err := apcs.assetPolicySvc.GetByID(assetPolicyID); err != nil {
		return nil, err
	}
	tx := apcs.repo.BeginTransaction()

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}

	for _, id := range ids {
		if _, err := apcs.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":           id,
			"asset_policy": assetPolicyID,
		})

		err := apcs.repo.Delete(id, tx)
		if err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
		}
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()

	// **Activity Logging**
	if len(successfulIDs) > 0 {
		go func() {
			err := apcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetPolicyComment{}.TableName(),
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
