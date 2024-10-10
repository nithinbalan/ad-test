package services

import (
	"errors"
	"fmt"
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditDraftFindingConversationServiceInterface interface {
	GetAll(auditDraftFindingID uint, queryParams *dtos.AuditDraftFindingConversationListQueryParamsDTO) ([]*dtos.AuditDraftFindingConversationListDTO, error)
	GetByID(auditDraftFindingID uint, id uint) (*models.AuditDraftFindingConversation, error)
	GetFullDetailsByID(auditDraftFindingID uint, id uint) (*dtos.AuditDraftFindingConversationDetailDTO, error)
	Create(auditDraftFindingID uint, currentUserID uint, formDTO *dtos.AuditDraftFindingConversationFormDTO) (uint, error)
	Update(auditDraftFindingID uint, id uint, currentUserID uint, auditFinding *dtos.AuditDraftFindingConversationFormDTO) error
	Delete(auditDraftFindingID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type AuditDraftFindingConversationService struct {
	repo                       repositories.AuditDraftFindingConversationRepositoryInterface
	activityLogSvc             logServices.ActivityLogServiceInterface
	auditDraftFindingSvc       AuditDraftFindingServiceInterface
	auditDraftFindingStatusSvc AuditDraftFindingStatusServiceInterface
	auditSvc                   AuditServiceInterface
}

func NewAuditDraftFindingConversationService(
	repo repositories.AuditDraftFindingConversationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditDraftFindingSvc AuditDraftFindingServiceInterface,
	auditDraftFindingStatusSvc AuditDraftFindingStatusServiceInterface,
	auditSvc AuditServiceInterface,
) AuditDraftFindingConversationServiceInterface {
	return &AuditDraftFindingConversationService{
		repo:                       repo,
		activityLogSvc:             activityLogSvc,
		auditDraftFindingSvc:       auditDraftFindingSvc,
		auditDraftFindingStatusSvc: auditDraftFindingStatusSvc,
		auditSvc:                   auditSvc,
	}
}

func (adfcs *AuditDraftFindingConversationService) GetAll(auditDraftFindingID uint, queryParams *dtos.AuditDraftFindingConversationListQueryParamsDTO) ([]*dtos.AuditDraftFindingConversationListDTO, error) {
	if _, err := adfcs.auditDraftFindingSvc.GetByID(auditDraftFindingID); err != nil {
		return nil, err
	}

	return adfcs.repo.GetAll(auditDraftFindingID, queryParams)
}

func (adfcs AuditDraftFindingConversationService) GetByID(auditDraftFindingID uint, id uint) (*models.AuditDraftFindingConversation, error) {
	if _, err := adfcs.auditDraftFindingSvc.GetByID(auditDraftFindingID); err != nil {
		return nil, err
	}

	auditDraftFindingConversation, err := adfcs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditDraftFindingConversation, nil
}

func (adfcs AuditDraftFindingConversationService) GetFullDetailsByID(auditDraftFindingID uint, id uint) (*dtos.AuditDraftFindingConversationDetailDTO, error) {
	var auditDraftFindingConversationDTO dtos.AuditDraftFindingConversationDetailDTO

	if _, err := adfcs.auditDraftFindingSvc.GetByID(auditDraftFindingID); err != nil {
		return nil, err
	}

	auditDraftFindingConversation, err := adfcs.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditDraftFindingConversation, &auditDraftFindingConversationDTO)
	if err != nil {
		return nil, err
	}

	auditDraftFindingConversationDTO.CreatedAt = auditDraftFindingConversation.CreatedAt
	auditDraftFindingConversationDTO.UpdatedAt = auditDraftFindingConversation.UpdatedAt

	return &auditDraftFindingConversationDTO, nil
}
func (adfcs *AuditDraftFindingConversationService) Create(auditDraftFindingID uint, currentUserID uint, formDTO *dtos.AuditDraftFindingConversationFormDTO) (uint, error) {
	tx := adfcs.repo.BeginTransaction()

	existingAuditDraftFinding, err := adfcs.auditDraftFindingSvc.GetByID(auditDraftFindingID)
	if err != nil {
		return 0, err
	}
	// Check if the type is 'AR', restrict to auditees or audit members
	if formDTO.Type == "AR" {
		// Fetch the related audit by AuditDraftFindingID
		auditID := existingAuditDraftFinding.AuditID
		audit, err := adfcs.auditSvc.GetFullDetailsByID(auditID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Check if the current user is the TeamLead
		isAuthorizedUser := false
		if uint(audit.AuditPlan.TeamLead.ID) == currentUserID {
			isAuthorizedUser = true
		}

		// Check if the current user is one of the AuditPlan Members
		if !isAuthorizedUser {
			for _, member := range audit.AuditPlan.Members {
				if member.UserID == currentUserID {
					isAuthorizedUser = true
					break
				}
			}
		}

		// If the current user is not authorized, rollback and return forbidden error
		if !isAuthorizedUser {
			tx.Rollback()
			return 0, utility.NewForbiddenError("The current user is not authorized to perform this action.")
		}
	}

	conversation := &models.AuditDraftFindingConversation{
		AuditDraftFindingID: auditDraftFindingID,
		Comment:         formDTO.Comment,
		Type:                formDTO.Type,
		IsAccepted:          *formDTO.IsAccepted,
		IsDismissed:         *formDTO.IsDismissed,
		CreatedBYID:         currentUserID,
	}

	id, err := adfcs.repo.Create(conversation, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Check if the form indicates that the finding should be accepted
	if formDTO.IsAccepted != nil && *formDTO.IsAccepted {
		// Check if the finding is already accepted
		auditDraftFindingStatus, err := adfcs.auditDraftFindingStatusSvc.GetByType("accepted")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if existingAuditDraftFinding.AuditDraftFindingStatusID == auditDraftFindingStatus.ID {
			tx.Rollback()
			return 0, utility.NewForbiddenError("The finding has already been accepted.")
		}
		// Update Audit Draft Finding Status
		if err := adfcs.auditDraftFindingSvc.StatusUpdate(auditDraftFindingID, currentUserID, auditDraftFindingStatus.ID); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Check if the form indicates that the finding should be IsDismissed
	if formDTO.IsDismissed != nil && *formDTO.IsDismissed {
		// Check if the finding is already IsDismissed
		auditDraftFindingStatusDismissed, err := adfcs.auditDraftFindingStatusSvc.GetByType("dismissed")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		auditDraftFindingStatusOpen, err := adfcs.auditDraftFindingStatusSvc.GetByType("open")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if existingAuditDraftFinding.AuditDraftFindingStatusID != auditDraftFindingStatusOpen.ID {
			tx.Rollback()
			return 0, utility.NewForbiddenError("Action not allowed. The audit finding must be in 'Open' status before any further actions can be taken.")
		}
		// Update Audit Draft Finding Status
		if err := adfcs.auditDraftFindingSvc.StatusUpdate(auditDraftFindingID, currentUserID, auditDraftFindingStatusDismissed.ID); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := adfcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditDraftFindingConversation{}.TableName(),
			id,
			conversation,
			currentUserID,
			nil,
		)

		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditDraftFindingConversation{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (adfcs *AuditDraftFindingConversationService) Update(auditDraftFindingID uint, id uint, currentUserID uint, AuditDraftFindingConversationFormDTO *dtos.AuditDraftFindingConversationFormDTO) error {
	tx := adfcs.repo.BeginTransaction()

	if _, err := adfcs.auditDraftFindingSvc.GetByID(auditDraftFindingID); err != nil {
		return err
	}

	existingAuditDraftFindingConversation, err := adfcs.GetByID(auditDraftFindingID, id)

	if err != nil {
		return err
	}

	// Update Audit Finding ActionPlan
	existingAuditDraftFindingConversation.Comment = AuditDraftFindingConversationFormDTO.Comment
	existingAuditDraftFindingConversation.Type = AuditDraftFindingConversationFormDTO.Type
	existingAuditDraftFindingConversation.UpdatedBYID = &currentUserID

	err = adfcs.repo.Update(existingAuditDraftFindingConversation, tx)

	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (adfcs *AuditDraftFindingConversationService) Delete(auditDraftFindingID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := adfcs.repo.BeginTransaction()

	if _, err := adfcs.auditDraftFindingSvc.GetByID(auditDraftFindingID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := adfcs.GetFullDetailsByID(auditDraftFindingID, id)

		if err != nil {
			continue
		}

		err = adfcs.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	tx.Commit()

	if len(successfulIds) > 0 {

		go func() {
			err := adfcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AuditFinding{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}
