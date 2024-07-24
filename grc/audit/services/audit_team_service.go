package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	auditRepositories "jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditTeamServiceInterface interface {
	GetAll(params *auditDTOs.AuditTeamListQueryParamsDTO) (*auditDTOs.AuditTeamListResponseDTO, error)
	GetByID(id uint) (*models.AuditTeam, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditTeamDetailDTO, error)
	Create(currentUserID uint, auditTeamFormDTO *auditDTOs.AuditTeamFormDTO) (uint, error)
	Update(id uint, currentUserID uint, audit *auditDTOs.AuditTeamFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type AuditTeamService struct {
	repo           auditRepositories.AuditTeamRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
}

func NewAuditTeamService(
	repo auditRepositories.AuditTeamRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface) AuditTeamServiceInterface {
	return &AuditTeamService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
	}
}

func (ats *AuditTeamService) GetAll(params *auditDTOs.AuditTeamListQueryParamsDTO) (*auditDTOs.AuditTeamListResponseDTO, error) {
	return ats.repo.GetAll(params)
}

func (ats AuditTeamService) GetByID(id uint) (*models.AuditTeam, error) {
	auditTeam, err := ats.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditTeam, nil
}

func (ats AuditTeamService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditTeamDetailDTO, error) {
	var auditTeamDTO auditDTOs.AuditTeamDetailDTO

	// Get full details of the audit action plan
	auditTeam, err := ats.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Decode the audit action plan into the DTO
	err = mapstructure.Decode(auditTeam, &auditTeamDTO)
	if err != nil {
		return nil, err
	}

	auditTeamDTO.CreatedAt = auditTeam.CreatedAt
	auditTeamDTO.UpdatedAt = auditTeam.UpdatedAt

	// TODO: Add any additional fields required for the details screen

	return &auditTeamDTO, nil
}

func (ats *AuditTeamService) Create(currentUserID uint, auditTeamFormDTO *auditDTOs.AuditTeamFormDTO) (uint, error) {
	tx := ats.repo.BeginTransaction()
	if err := ats.validateAuditTeamFormDTO(auditTeamFormDTO); err != nil {
		return 0, err
	}

	auditTeam := &models.AuditTeam{
		TitleEn:     auditTeamFormDTO.TitleEn,
		TitleAr:     auditTeamFormDTO.TitleAr,
		TeamLeadID:  auditTeamFormDTO.TeamLeadID,
		CreatedBYID: currentUserID,
	}

	id, err := ats.repo.Create(auditTeam, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Save Audit Team Members
	var users []*models.AuditTeamMember
	for _, userId := range auditTeamFormDTO.UserIDs {
		user := &models.AuditTeamMember{
			AuditTeamID: id,
			UserID:      *userId,
			CreatedBYID: currentUserID,
		}
		users = append(users, user)
	}

	if len(users) > 0 {
		if err := ats.repo.CreateMembers(users, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := ats.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditTeam{}.TableName(),
			id,
			auditTeamFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditTeam{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ats *AuditTeamService) Update(id uint, currentUserID uint, auditTeamFormDTO *auditDTOs.AuditTeamFormDTO) error {
	tx := ats.repo.BeginTransaction()

	existingAuditTeam, err := ats.GetByID(id)

	if err != nil {
		return err
	}

	if err := ats.validateAuditTeamFormDTO(auditTeamFormDTO); err != nil {
		return err
	}

	// Update Compliance ActionPlan
	existingAuditTeam.TitleEn = auditTeamFormDTO.TitleEn
	existingAuditTeam.TitleAr = auditTeamFormDTO.TitleAr
	existingAuditTeam.TeamLeadID = auditTeamFormDTO.TeamLeadID
	existingAuditTeam.UpdatedBYID = &currentUserID

	err = ats.repo.Update(existingAuditTeam, tx)

	if err != nil {
		tx.Rollback()
		return err
	}
	// Collect existing user IDs
	existingAuditTeamDetails, err := ats.repo.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingUserIDs := make(map[uint]bool)
	for _, user := range existingAuditTeamDetails.Users {
		existingUserIDs[user.ID] = true
	}
	// Collect new user IDs
	newUserIDs := make(map[uint]bool)
	for _, user := range auditTeamFormDTO.UserIDs {
		newUserIDs[*user] = true
	}
	// Add new users
	var users []*models.AuditTeamMember
	for _, newUser := range auditTeamFormDTO.UserIDs {
		if _, exists := existingUserIDs[*newUser]; !exists {
			user := &models.AuditTeamMember{
				AuditTeamID: existingAuditTeamDetails.ID,
				UserID:      *newUser,
				CreatedBYID: currentUserID,
			}
			users = append(users, user)
		}
	}
	// Save the new user
	if len(users) > 0 {
		if err := ats.repo.CreateMembers(users, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove users not in the new list
	var deleteUserIDs []*uint
	for _, existingUser := range existingAuditTeamDetails.Users {
		if _, exists := newUserIDs[existingUser.ID]; !exists {
			// Delete the existing user from the repository
			deleteUserIDs = append(deleteUserIDs, &existingUser.ID)
		}
	}
	// Delete the user
	if len(deleteUserIDs) > 0 {
		if err := ats.repo.DeleteMembers(deleteUserIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	return nil
}

func (ats *AuditTeamService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ats.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		auditTeam, err := ats.GetFullDetailsByID(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if len(auditTeam.Users) > 0 {
			var deleteUserIds []*uint
			for _, user := range auditTeam.Users {
				deleteUserIds = append(deleteUserIds, &user.ID)
			}
			if err := ats.repo.DeleteMembers(deleteUserIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		err = ats.repo.Delete(id, tx) // Added tx here
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		successfulIds = append(successfulIds, id)
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}
	tx.Commit()

	if len(successfulIds) > 0 {

		go func() {
			err := ats.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AuditTeam{}.TableName(),
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

func (ats *AuditTeamService) validateAuditTeamFormDTO(auditTeamFormDTO *auditDTOs.AuditTeamFormDTO) error {
	if _, err := ats.userSvc.GetByID(auditTeamFormDTO.TeamLeadID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit team lead does not exist"), "not_exist", "team lead")
	}
	// Validate User IDs
	userIDs := auditTeamFormDTO.UserIDs
	userIDMap := make(map[uint]bool)

	for _, id := range userIDs {
		if userIDMap[*id] {
			return utility.NewValidationError(fmt.Errorf("duplicate user ID: %d", *id), "duplicate", "user")
		}

		if _, err := ats.userSvc.GetByID(*id); err != nil {
			return utility.NewValidationError(fmt.Errorf("user does not exist for ID: %d", *id), "not_exist", "user")
		}

		userIDMap[*id] = true
	}

	return nil
}
