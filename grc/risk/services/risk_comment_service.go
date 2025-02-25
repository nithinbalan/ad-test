package services

import (
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	"jethur/utility"
)

type RiskCommentServiceInterface interface {
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
}

func NewRiskCommentService(repo repositories.RiskCommentRepositoryInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	riskStatusSvc RiskStatusServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface) *RiskCommentService {
	return &RiskCommentService{repo: repo,
		riskRegisterSvc: riskRegisterSvc,
		riskStatusSvc:   riskStatusSvc,
		activityLogSvc:  activityLogSvc,
	}
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
		RiskStatusID:       revertedStatus.ID,
		CreatedByID:        currentUserID,
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
		RiskStatusID:       rejectedStatus.ID,
		CreatedByID:        currentUserID,
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
		RiskStatusID:       approvedStatus.ID,
		CreatedByID:        currentUserID,
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

	// Validate if the current status allows resubmission
	if risk.RiskStatusID == resubmittedStatus.ID {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("risk is already resubmitted"), "invalid_state", "risk")
	}

	comment := &models.RiskComment{
		RiskRegisterID:     riskID,
		RiskTimelineTypeID: commentForm.RiskTimelineTypeID,
		Comments:           commentForm.Comment,
		RiskStatusID:       resubmittedStatus.ID,
		CreatedByID:        currentUserID,
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
