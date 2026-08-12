package services

import (
	"errors"
	"fmt"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IncidentResponseTeamServiceInterface interface {
	GetAll(params *incidentDTOs.IncidentResponseTeamListQueryParamsDTO) (*incidentDTOs.IncidentResponseTeamListResponseDTO, error)
	GetByID(id uint) (*models.IncidentResponseTeam, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentResponseTeamDetailDTO, error)
	Create(currentUserID uint, teamFormDTO *incidentDTOs.IncidentResponseTeamFormDTO) (uint, error)
	Update(id uint, currentUserID uint, team *incidentDTOs.IncidentResponseTeamFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type IncidentResponseTeamService struct {
	repo           incidentRepositories.IncidentResponseTeamRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
}

func NewIncidentResponseTeamService(
	repo incidentRepositories.IncidentResponseTeamRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) IncidentResponseTeamServiceInterface {
	return &IncidentResponseTeamService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
	}
}

func (irts *IncidentResponseTeamService) GetAll(params *incidentDTOs.IncidentResponseTeamListQueryParamsDTO) (*incidentDTOs.IncidentResponseTeamListResponseDTO, error) {
	teams, totalCount, err := irts.repo.GetAll(params)
	if err != nil {
		return nil, err
	}

	var teamDTOs []*incidentDTOs.IncidentResponseTeamListDTO
	for _, team := range teams {
		if team == nil {
			continue
		}

		dto := &incidentDTOs.IncidentResponseTeamListDTO{
			ID:                           team.ID,
			TitleEn:                      team.TitleEn,
			TitleAr:                      team.TitleAr,
			Scope:                        team.Scope,
			Objective:                    team.Objective,
			TeamLeadUserID:               team.TeamLeadUserID,
			TeamLeadUserFirstNameEn:      getSafeString(team.TeamLeadUser, func(u *models.User) string { return u.FirstNameEn }),
			TeamLeadUserLastNameEn:       getSafeString(team.TeamLeadUser, func(u *models.User) string { return u.LastNameEn }),
			TeamLeadUserFirstNameAr:      getSafeString(team.TeamLeadUser, func(u *models.User) string { return u.FirstNameAr }),
			TeamLeadUserLastNameAr:       getSafeString(team.TeamLeadUser, func(u *models.User) string { return u.LastNameAr }),
			TeamLeadUserEmail:            getSafeString(team.TeamLeadUser, func(u *models.User) string { return u.Email }),
			TeamLeadUserInitial:          getSafeString(team.TeamLeadUser, func(u *models.User) string { return u.Initial }),
			TeamLeadUserColor:            getSafeString(team.TeamLeadUser, func(u *models.User) string { return u.Color }),
			TeamLeadUserProfileImagePath: getSafeStringPointer(team.TeamLeadUser, func(u *models.User) *string { return u.ProfileImagePath }),
			CreatedUserID:                team.CreatedBYID,
			CreatedUserFirstNameEn:       getSafeString(team.CreatedByUser, func(u *models.User) string { return u.FirstNameEn }),
			CreatedUserFirstNameAr:       getSafeString(team.CreatedByUser, func(u *models.User) string { return u.FirstNameAr }),
			CreatedUserLastNameEn:        getSafeString(team.CreatedByUser, func(u *models.User) string { return u.LastNameEn }),
			CreatedUserLastNameAr:        getSafeString(team.CreatedByUser, func(u *models.User) string { return u.LastNameAr }),
			CreatedUserEmail:             getSafeString(team.CreatedByUser, func(u *models.User) string { return u.Email }),
			CreatedUserInitial:           getSafeString(team.CreatedByUser, func(u *models.User) string { return u.Initial }),
			CreatedUserColor:             getSafeString(team.CreatedByUser, func(u *models.User) string { return u.Color }),
			CreatedUserProfileImagePath:  getSafeStringPointer(team.CreatedByUser, func(u *models.User) *string { return u.ProfileImagePath }),
			UpdatedUserID:                team.UpdatedBYID,
			UpdatedUserFirstNameEn:       getSafeString(team.UpdatedByUser, func(u *models.User) string { return u.FirstNameEn }),
			UpdatedUserFirstNameAr:       getSafeString(team.UpdatedByUser, func(u *models.User) string { return u.FirstNameAr }),
			UpdatedUserLastNameEn:        getSafeString(team.UpdatedByUser, func(u *models.User) string { return u.LastNameEn }),
			UpdatedUserLastNameAr:        getSafeString(team.UpdatedByUser, func(u *models.User) string { return u.LastNameAr }),
			UpdatedUserEmail:             getSafeString(team.UpdatedByUser, func(u *models.User) string { return u.Email }),
			UpdatedUserInitial:           getSafeString(team.UpdatedByUser, func(u *models.User) string { return u.Initial }),
			UpdatedUserColor:             getSafeString(team.UpdatedByUser, func(u *models.User) string { return u.Color }),
			UpdatedUserProfileImagePath:  getSafeStringPointer(team.UpdatedByUser, func(u *models.User) *string { return u.ProfileImagePath }),
			CreatedAt:                    team.CreatedAt.String(),
			UpdatedAt:                    getSafeString(team.UpdatedAt, func(t *time.Time) string { return t.String() }),
		}

		for _, member := range team.IncidentResponseTeamMembers {
			if member == nil || member.User == nil {
				continue
			}

			dto.Members = append(dto.Members, &incidentDTOs.IncidentResponseTeamMembersDTO{
				ID:                     member.ID,
				IncidentResponseTeamID: member.IncidentResponseTeamID,
				UserID:                 member.UserID,
				User: organizationDTOs.UserMinimalDTO{
					ID:               int(member.User.ID),
					FirstNameEN:      member.User.FirstNameEn,
					LastNameEN:       member.User.LastNameEn,
					FirstNameAR:      member.User.FirstNameAr,
					LastNameAR:       member.User.LastNameAr,
					Initial:          member.User.Initial,
					Color:            member.User.Color,
					Email:            member.User.Email,
					MobileNumber:     member.User.MobileNumber,
					ProfileImagePath: member.User.ProfileImagePath,
				},
			})
		}

		teamDTOs = append(teamDTOs, dto)
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	response := &incidentDTOs.IncidentResponseTeamListResponseDTO{
		Items: teamDTOs,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (irts IncidentResponseTeamService) GetByID(id uint) (*models.IncidentResponseTeam, error) {
	team, err := irts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return team, nil
}

func (irts IncidentResponseTeamService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentResponseTeamDetailDTO, error) {
	var teamDTO incidentDTOs.IncidentResponseTeamDetailDTO

	team, err := irts.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(team, &teamDTO)
	if err != nil {
		return nil, err
	}

	teamDTO.CreatedAt = team.CreatedAt
	teamDTO.UpdatedAt = team.UpdatedAt

	return &teamDTO, nil
}

func (irts *IncidentResponseTeamService) Create(currentUserID uint, teamFormDTO *incidentDTOs.IncidentResponseTeamFormDTO) (uint, error) {
	tx := irts.repo.BeginTransaction()
	if err := irts.validateIncidentResponseTeamFormDTO(teamFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	team := &models.IncidentResponseTeam{
		TitleEn:        teamFormDTO.TitleEn,
		TitleAr:        teamFormDTO.TitleAr,
		TeamLeadUserID: teamFormDTO.TeamLeadUserID,
		Scope:          teamFormDTO.Scope,
		Objective:      teamFormDTO.Objective,
		CreatedBYID:    currentUserID,
	}

	id, err := irts.repo.Create(team, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var members []*models.IncidentResponseTeamMember
	for _, userId := range teamFormDTO.UserIDs {
		member := &models.IncidentResponseTeamMember{
			IncidentResponseTeamID: id,
			UserID:                 *userId,
			CreatedBYID:            currentUserID,
		}
		members = append(members, member)
	}

	if len(members) > 0 {
		if err := irts.repo.CreateMembers(members, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		existingTeam, err := irts.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch incident response team", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = irts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentResponseTeam{}.TableName(),
			id,
			existingTeam,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentResponseTeam{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (irts *IncidentResponseTeamService) Update(id uint, currentUserID uint, teamFormDTO *incidentDTOs.IncidentResponseTeamFormDTO) error {
	tx := irts.repo.BeginTransaction()

	oldTeam, err := irts.repo.GetFullDetailsByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	existingTeam, err := irts.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := irts.validateIncidentResponseTeamFormDTO(teamFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	existingTeam.TitleEn = teamFormDTO.TitleEn
	existingTeam.TitleAr = teamFormDTO.TitleAr
	existingTeam.TeamLeadUserID = teamFormDTO.TeamLeadUserID
	existingTeam.Scope = teamFormDTO.Scope
	existingTeam.Objective = teamFormDTO.Objective
	existingTeam.UpdatedBYID = &currentUserID

	err = irts.repo.Update(existingTeam, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingTeamDetails, err := irts.repo.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingUserIDs := make(map[uint]bool)
	for _, member := range existingTeamDetails.IncidentResponseTeamMembers {
		existingUserIDs[member.ID] = true
	}

	newUserIDs := make(map[uint]bool)
	for _, userId := range teamFormDTO.UserIDs {
		newUserIDs[*userId] = true
	}

	var members []*models.IncidentResponseTeamMember
	for _, newUserId := range teamFormDTO.UserIDs {
		if _, exists := existingUserIDs[*newUserId]; !exists {
			member := &models.IncidentResponseTeamMember{
				IncidentResponseTeamID: existingTeamDetails.ID,
				UserID:                 *newUserId,
				CreatedBYID:            currentUserID,
			}
			members = append(members, member)
		}
	}

	if len(members) > 0 {
		if err := irts.repo.CreateMembers(members, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	var deleteMemberIDs []*uint
	for _, existingMember := range existingTeamDetails.IncidentResponseTeamMembers {
		if _, exists := newUserIDs[existingMember.ID]; !exists {
			deleteMemberIDs = append(deleteMemberIDs, &existingMember.ID)
		}
	}

	if len(deleteMemberIDs) > 0 {
		if err := irts.repo.DeleteMembers(deleteMemberIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		existingTeam, err := irts.repo.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		membersAdded := []uint{}
		membersRemoved := []uint{}
		for id := range newUserIDs {
			if !existingUserIDs[id] {
				membersAdded = append(membersAdded, id)
			}
		}
		for id := range existingUserIDs {
			if !newUserIDs[id] {
				membersRemoved = append(membersRemoved, id)
			}
		}

		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})

		if len(membersAdded) > 0 || len(membersRemoved) > 0 {
			associatedItems["user_ids"] = map[string]interface{}{
				"added":   membersAdded,
				"removed": membersRemoved,
			}
		}
		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}

		oldData := map[string]interface{}{
			"title_en": oldTeam.TitleEn,
			"title_ar": oldTeam.TitleAr,
		}

		newData := map[string]interface{}{
			"title_en": existingTeam.TitleEn,
			"title_ar": existingTeam.TitleAr,
		}

		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		if len(changedData) > 0 {
			err := irts.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.IncidentResponseTeam{}.TableName(),
				id,
				changedData,
				currentUserID,
				"",
				"",
				"",
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.IncidentResponseTeam{}.TableName(),
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

func (irts *IncidentResponseTeamService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := irts.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedTeam []map[string]interface{}

	for _, id := range ids {
		team, err := irts.GetFullDetailsByID(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if len(team.IncidentResponseTeamMembers) > 0 {
			var deleteMemberIds []*uint
			for _, member := range team.IncidentResponseTeamMembers {
				deleteMemberIds = append(deleteMemberIds, &member.ID)
			}
			if err := irts.repo.DeleteMembers(deleteMemberIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		err = irts.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		successfulIds = append(successfulIds, id)

		deletedTeam = append(deletedTeam, map[string]interface{}{
			"id":       id,
			"title_en": team.TitleEn,
			"title_ar": team.TitleAr,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}
	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := irts.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IncidentResponseTeam{}.TableName(),
				successfulIds,
				deletedTeam,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (irts *IncidentResponseTeamService) validateIncidentResponseTeamFormDTO(teamFormDTO *incidentDTOs.IncidentResponseTeamFormDTO) error {
	if teamFormDTO.TeamLeadUserID != nil {
		_, err := irts.userSvc.GetByID(*teamFormDTO.TeamLeadUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("team lead user does not exist"), "not_exist", "team lead user")
		}
	}

	userIDs := teamFormDTO.UserIDs
	userIDMap := make(map[uint]bool)

	for _, id := range userIDs {
		if userIDMap[*id] {
			return utility.NewValidationError(fmt.Errorf("duplicate user ID: %d", *id), "duplicate", "user")
		}

		_, err := irts.userSvc.GetByID(*id)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("user does not exist for ID: %d", *id), "not_exist", "user")
		}

		userIDMap[*id] = true
	}

	return nil
}

func getSafeString[T any](ptr *T, getString func(*T) string) string {
	if ptr == nil {
		return ""
	}
	return getString(ptr)
}

func getSafeStringPointer[T any](ptr *T, getString func(*T) *string) *string {
	if ptr == nil {
		return nil
	}
	return getString(ptr)
}
