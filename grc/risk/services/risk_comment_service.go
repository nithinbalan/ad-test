package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskCommentServiceInterface interface {
	GetAll(riskID uint) ([]*dtos.RiskCommentListDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskComment, error)
	GetFullDetailsByID(riskID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.RiskCommentDetailDTO, error)
	Create(riskID uint, currentUserID uint, dto *dtos.RiskCommentFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, dto *dtos.RiskCommentFormDTO) error
	Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error)

	SubmitRisk(riskID uint, currentUserID uint) error
	RevertRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error
	RejectRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error
	ApproveRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error
	ResubmitRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error
}

type RiskCommentService struct {
	repo            repositories.RiskCommentRepositoryInterface
	riskRegisterSvc RiskRegisterServiceInterface
	riskStatusSvc   RiskStatusServiceInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	userSvc         organizationServices.UserServiceInterface
	riskSettingsSvc RiskSettingsServiceInterface
}

func NewRiskCommentService(repo repositories.RiskCommentRepositoryInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	riskStatusSvc RiskStatusServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	riskSettingsSvc RiskSettingsServiceInterface,
) *RiskCommentService {
	return &RiskCommentService{repo: repo,
		riskRegisterSvc: riskRegisterSvc,
		riskStatusSvc:   riskStatusSvc,
		activityLogSvc:  activityLogSvc,
		userSvc:         userSvc,
		riskSettingsSvc: riskSettingsSvc,
	}
}

func (rcs *RiskCommentService) GetAll(riskID uint) ([]*dtos.RiskCommentListDTO, error) {
	return rcs.repo.GetAll(riskID)
}

func (rcs *RiskCommentService) GetByID(riskID uint, id uint) (*models.RiskComment, error) {
	comment, err := rcs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (rcs *RiskCommentService) GetFullDetailsByID(riskID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.RiskCommentDetailDTO, error) {
	var commentDetail dtos.RiskCommentDetailDTO

	comment, err := rcs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(comment, &commentDetail)
	if err != nil {
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt
	return &commentDetail, nil
}

func (rcs *RiskCommentService) Create(riskID uint, currentUserID uint, riskCommentFormDTO *dtos.RiskCommentFormDTO) (uint, error) {
	tx := rcs.repo.BeginTransaction()

	comment := &models.RiskComment{
		RiskRegisterID: riskID,
		Comments:       riskCommentFormDTO.Comment,
		CreatedBYID:    currentUserID,
	}

	id, err := rcs.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingComment, err := rcs.GetFullDetailsByID(riskID, id)
		if err != nil {
			utility.LogError("Failed to fetch risk comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = rcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskComment{}.TableName(),
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

func (rcs *RiskCommentService) Update(riskID uint, id uint, currentUserID uint, riskCommentFormDTO *dtos.RiskCommentFormDTO) error {
	tx := rcs.repo.BeginTransaction()

	existingComment, err := rcs.GetByID(riskID, id)
	if err != nil {
		return err
	}

	oldComment, err := rcs.GetFullDetailsByID(riskID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comments = riskCommentFormDTO.Comment
	err = rcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingComment, err := rcs.GetFullDetailsByID(riskID, id)
		if err != nil {
			return
		}

		changedData := utility.CompareFieldChanges(
			map[string]interface{}{"comment": oldComment.Comment},
			map[string]interface{}{"comment": existingComment.Comment},
		)

		if len(changedData) > 0 {
			err := rcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskComment{}.TableName(),
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

func (rcs *RiskCommentService) Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := rcs.repo.BeginTransaction()

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}

	for _, id := range ids {
		if _, err := rcs.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":      id,
			"risk_id": riskID,
		})

		err := rcs.repo.Delete(id, tx)
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
	if len(successfulIDs) > 0 {
		go func() {
			err := rcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskComment{}.TableName(),
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

func (rcs *RiskCommentService) SubmitRisk(riskID uint, currentUserID uint) error {
	tx := rcs.repo.BeginTransaction()

	submittedStatus, err := rcs.riskStatusSvc.GetByType("submitted")
	if err != nil {
		tx.Rollback()
		return err
	}

	risk, err := rcs.riskRegisterSvc.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Ensure only the responsible user can submit
	if risk.ResponsibleUserID != currentUserID {
		tx.Rollback()
		return utility.NewForbiddenError("Only the responsible user can submit this risk.")
	}

	// Update risk status to "submitted"
	risk.RiskStatusID = submittedStatus.ID
	err = rcs.repo.UpdateRiskStatus(risk, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the action
	go func() {
		err := rcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskRegister{}.TableName(),
			riskID,
			map[string]interface{}{
				"status": "submitted",
				"user":   currentUserID,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  riskID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}

func (rcs *RiskCommentService) RevertRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error {
	tx := rcs.repo.BeginTransaction()

	riskSetting, err := rcs.riskSettingsSvc.Get()
	if err != nil {
		return err
	}

	// Get Review User IDs
	reviewUserIDs, err := rcs.userSvc.GetUserIdsByRole(riskSetting.ReviewRoleID)
	if err != nil {
		return err
	}

	if !contains(reviewUserIDs, currentUserID) {
		return utility.NewForbiddenError("Only the review user can revert this risk.")
	}

	revertedStatus, err := rcs.riskStatusSvc.GetByType("reverted")
	if err != nil {
		tx.Rollback()
		return err
	}

	risk, err := rcs.riskRegisterSvc.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	comment := &models.RiskComment{
		RiskRegisterID:     riskID,
		RiskTimelineTypeID: commentForm.RiskTimelineTypeID,
		Comments:           commentForm.Comment,
		RiskStatusID:       &revertedStatus.ID,
		CreatedBYID:        currentUserID,
	}

	_, err = rcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update risk status to "reverted"
	risk.RiskStatusID = revertedStatus.ID
	err = rcs.repo.UpdateRiskStatus(risk, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the action
	go func() {
		err := rcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskRegister{}.TableName(),
			riskID,
			map[string]interface{}{
				"status":  "reverted",
				"comment": commentForm.Comment,
				"user":    currentUserID,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  riskID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}

func (rcs *RiskCommentService) RejectRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error {
	tx := rcs.repo.BeginTransaction()

	riskSetting, err := rcs.riskSettingsSvc.Get()
	if err != nil {
		return err
	}

	// Get Review User IDs
	reviewUserIDs, err := rcs.userSvc.GetUserIdsByRole(riskSetting.ReviewRoleID)
	if err != nil {
		return err
	}

	if !contains(reviewUserIDs, currentUserID) {
		return utility.NewForbiddenError("Only the review user can reject this risk.")
	}

	rejectedStatus, err := rcs.riskStatusSvc.GetByType("rejected")
	if err != nil {
		tx.Rollback()
		return err
	}

	risk, err := rcs.riskRegisterSvc.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	comment := &models.RiskComment{
		RiskRegisterID:     riskID,
		RiskTimelineTypeID: commentForm.RiskTimelineTypeID,
		Comments:           commentForm.Comment,
		RiskStatusID:       &rejectedStatus.ID,
		CreatedBYID:        currentUserID,
	}

	_, err = rcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update risk status to "rejected"
	risk.RiskStatusID = rejectedStatus.ID
	err = rcs.repo.UpdateRiskStatus(risk, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the action
	go func() {
		err := rcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskRegister{}.TableName(),
			riskID,
			map[string]interface{}{
				"status":  "rejected",
				"comment": commentForm.Comment,
				"user":    currentUserID,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  riskID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}

func (rcs *RiskCommentService) ApproveRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error {
	tx := rcs.repo.BeginTransaction()

	riskSetting, err := rcs.riskSettingsSvc.Get()
	if err != nil {
		return err
	}

	// Get Review User IDs
	reviewUserIDs, err := rcs.userSvc.GetUserIdsByRole(riskSetting.ReviewRoleID)
	if err != nil {
		return err
	}

	if !contains(reviewUserIDs, currentUserID) {
		return utility.NewForbiddenError("Only the review user can approve this risk.")
	}

	approvedStatus, err := rcs.riskStatusSvc.GetByType("approved")
	if err != nil {
		tx.Rollback()
		return err
	}

	risk, err := rcs.riskRegisterSvc.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Validate if the current status allows approval
	if risk.RiskStatusID == approvedStatus.ID {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("risk is already approved"), "invalid_state", "risk")
	}

	comment := &models.RiskComment{
		RiskRegisterID:     riskID,
		RiskTimelineTypeID: commentForm.RiskTimelineTypeID,
		Comments:           commentForm.Comment,
		RiskStatusID:       &approvedStatus.ID,
		CreatedBYID:        currentUserID,
	}

	_, err = rcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update risk status to "approved"
	risk.RiskStatusID = approvedStatus.ID
	err = rcs.repo.UpdateRiskStatus(risk, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the action
	go func() {
		err := rcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskRegister{}.TableName(),
			riskID,
			map[string]interface{}{
				"status":  "approved",
				"comment": commentForm.Comment,
				"user":    currentUserID,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  riskID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}

func (rcs *RiskCommentService) ResubmitRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error {
	tx := rcs.repo.BeginTransaction()

	resubmittedStatus, err := rcs.riskStatusSvc.GetByType("resubmitted")
	if err != nil {
		tx.Rollback()
		return err
	}

	risk, err := rcs.riskRegisterSvc.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Ensure only the responsible user can submit
	if risk.ResponsibleUserID != currentUserID {
		tx.Rollback()
		return utility.NewForbiddenError("Only the responsible user can re-submit this risk.")
	}

	// Validate if the current status allows resubmission
	if risk.RiskStatusID == resubmittedStatus.ID {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("risk is already resubmitted"), "invalid_state", "risk")
	}

	comment := &models.RiskComment{
		RiskRegisterID:     riskID,
		RiskTimelineTypeID: commentForm.RiskTimelineTypeID,
		Comments:           commentForm.Comment,
		RiskStatusID:       &resubmittedStatus.ID,
		CreatedBYID:        currentUserID,
	}

	_, err = rcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update risk status to "resubmitted"
	risk.RiskStatusID = resubmittedStatus.ID
	err = rcs.repo.UpdateRiskStatus(risk, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the action
	go func() {
		err := rcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskRegister{}.TableName(),
			riskID,
			map[string]interface{}{
				"status":  "resubmitted",
				"comment": commentForm.Comment,
				"user":    currentUserID,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  riskID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}

func contains(slice []uint, value uint) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
