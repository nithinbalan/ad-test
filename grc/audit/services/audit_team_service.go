package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	auditRepositories "jethur/grc/audit/repositories"
	commonDTOs "jethur/grc/common/dtos"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"time"

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
	auditTeams, totalCount, err := ats.repo.GetAll(params)
	if err != nil {
		return nil, err
	}

	// Map results to DTOs for response
	var auditTeamDTOs []*auditDTOs.AuditTeamListDTO
	for _, auditTeam := range auditTeams {
		if auditTeam == nil {
			continue // Skip nil auditTeam
		}

		dto := &auditDTOs.AuditTeamListDTO{
			ID:                          auditTeam.ID,
			TitleEn:                     auditTeam.TitleEn,
			TitleAr:                     auditTeam.TitleAr,
			TeamLeadID:                  auditTeam.TeamLeadID,
			TeamLeadFirstNameEn:         getSafeString(auditTeam.TeamLead, func(u *models.User) string { return u.FirstNameEn }),
			TeamLeadFirstNameAr:         getSafeString(auditTeam.TeamLead, func(u *models.User) string { return u.FirstNameAr }),
			TeamLeadLastNameEn:          getSafeString(auditTeam.TeamLead, func(u *models.User) string { return u.LastNameEn }),
			TeamLeadLastNameAr:          getSafeString(auditTeam.TeamLead, func(u *models.User) string { return u.LastNameAr }),
			TeamLeadEmail:               getSafeString(auditTeam.TeamLead, func(u *models.User) string { return u.Email }),
			TeamLeadInitial:             getSafeString(auditTeam.TeamLead, func(u *models.User) string { return u.Initial }),
			TeamLeadColor:               getSafeString(auditTeam.TeamLead, func(u *models.User) string { return u.Color }),
			TeamLeadProfileImagePath:    getSafeStringPointer(auditTeam.TeamLead, func(u *models.User) *string { return u.ProfileImagePath }),
			CreatedUserFirstNameEn:      getSafeString(auditTeam.CreatedByUser, func(u *models.User) string { return u.FirstNameEn }),
			CreatedUserFirstNameAr:      getSafeString(auditTeam.CreatedByUser, func(u *models.User) string { return u.FirstNameAr }),
			CreatedUserLastNameEn:       getSafeString(auditTeam.CreatedByUser, func(u *models.User) string { return u.LastNameEn }),
			CreatedUserLastNameAr:       getSafeString(auditTeam.CreatedByUser, func(u *models.User) string { return u.LastNameAr }),
			CreatedUserEmail:            getSafeString(auditTeam.CreatedByUser, func(u *models.User) string { return u.Email }),
			CreatedUserInitial:          getSafeString(auditTeam.CreatedByUser, func(u *models.User) string { return u.Initial }),
			CreatedUserColor:            getSafeString(auditTeam.CreatedByUser, func(u *models.User) string { return u.Color }),
			CreatedUserProfileImagePath: getSafeStringPointer(auditTeam.CreatedByUser, func(u *models.User) *string { return u.ProfileImagePath }),
			UpdatedUserFirstNameEn:      getSafeString(auditTeam.UpdatedByUser, func(u *models.User) string { return u.FirstNameEn }),
			UpdatedUserFirstNameAr:      getSafeString(auditTeam.UpdatedByUser, func(u *models.User) string { return u.FirstNameAr }),
			UpdatedUserLastNameEn:       getSafeString(auditTeam.UpdatedByUser, func(u *models.User) string { return u.LastNameEn }),
			UpdatedUserLastNameAr:       getSafeString(auditTeam.UpdatedByUser, func(u *models.User) string { return u.LastNameAr }),
			UpdatedUserEmail:            getSafeString(auditTeam.UpdatedByUser, func(u *models.User) string { return u.Email }),
			UpdatedUserInitial:          getSafeString(auditTeam.UpdatedByUser, func(u *models.User) string { return u.Initial }),
			UpdatedUserColor:            getSafeString(auditTeam.UpdatedByUser, func(u *models.User) string { return u.Color }),
			UpdatedUserProfileImagePath: getSafeStringPointer(auditTeam.UpdatedByUser, func(u *models.User) *string { return u.ProfileImagePath }),
			CreatedAt:                   auditTeam.CreatedAt.String(),
			UpdatedAt:                   getSafeString(auditTeam.UpdatedAt, func(t *time.Time) string { return t.String() }),
		}

		// Map team members
		for _, member := range auditTeam.Users {
			if member == nil || member.User == nil {
				continue // Skip nil members or users
			}

			dto.Users = append(dto.Users, &auditDTOs.AuditTeamMembersDTO{
				ID:          member.ID,
				AuditTeamID: member.AuditTeamID,
				UserID:      member.UserID,
				User: organizationDTOs.UserMinimalDTO{
					ID:               int(member.User.ID),
					FirstNameEN:      member.User.FirstNameEn,
					LastNameEN:       member.User.LastNameEn,
					FirstNameAR:      member.User.FirstNameAr,
					LastNameAR:       member.User.LastNameAr,
					ProfileImagePath: member.User.ProfileImagePath,
					Email:            member.User.Email,
					Initial:          member.User.Initial,
					Color:            member.User.Color,
				},
			})
		}

		auditTeamDTOs = append(auditTeamDTOs, dto)
	}
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	// Response with pagination
	response := &auditDTOs.AuditTeamListResponseDTO{
		Items: auditTeamDTOs,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
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
		existingAuditTeam, err := ats.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch audit team", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = ats.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditTeam{}.TableName(),
			id,
			existingAuditTeam,
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

	oldAuditTeam, err := ats.repo.GetFullDetailsByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

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
	go func() {
		existingAuditTeam, err := ats.repo.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		usersAdded := []uint{}
		usersRemoved := []uint{}
		for id := range newUserIDs {
			if !existingUserIDs[id] {
				usersAdded = append(usersAdded, id)
			}
		}
		for id := range existingUserIDs {
			if !newUserIDs[id] {
				usersRemoved = append(usersRemoved, id)
			}
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})

		if len(usersAdded) > 0 || len(usersRemoved) > 0 {
			associatedItems["user_ids"] = map[string]interface{}{
				"added":   usersAdded,
				"removed": usersRemoved,
			}
		}
		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}
		//log changes
		oldData := map[string]interface{}{
			"title_en":  oldAuditTeam.TitleEn,
			"title_ar":  oldAuditTeam.TitleAr,
			"team_lead": fmt.Sprintf("%s %s", oldAuditTeam.TeamLead.FirstNameEn, oldAuditTeam.TeamLead.LastNameEn),
		}

		newData := map[string]interface{}{
			"title_en":  existingAuditTeam.TitleEn,
			"title_ar":  existingAuditTeam.TitleAr,
			"team_lead": fmt.Sprintf("%s %s", existingAuditTeam.TeamLead.FirstNameEn, existingAuditTeam.TeamLead.LastNameEn),
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := ats.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AuditTeam{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AuditTeam{}.TableName(),
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

func (ats *AuditTeamService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ats.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedAuditTeam []map[string]interface{}
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

		deletedAuditTeam = append(deletedAuditTeam, map[string]interface{}{
			"id":       id,
			"title_en": auditTeam.TitleEn,
			"title_ar": auditTeam.TitleAr,
		})
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
				deletedAuditTeam,
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

// Helper function to safely get string or empty string from a pointer
func getSafeString[T any](ptr *T, getString func(*T) string) string {
	if ptr == nil {
		return ""
	}
	return getString(ptr)
}

// Helper function to safely get *string or nil from a pointer
func getSafeStringPointer[T any](ptr *T, getString func(*T) *string) *string {
	if ptr == nil {
		return nil
	}
	return getString(ptr)
}
