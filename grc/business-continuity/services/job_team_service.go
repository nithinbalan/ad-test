package services

import (
	"errors"
	"fmt"
	"time"

	bcDTOs "jethur/grc/business-continuity/dtos"
	bcRepositories "jethur/grc/business-continuity/repositories"
	commonDTOs "jethur/grc/common/dtos"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type JobTeamServiceInterface interface {
	GetAll(params *bcDTOs.JobTeamListQueryParamsDTO) (*bcDTOs.JobTeamListResponseDTO, error)
	GetByID(id uint) (*models.JobTeam, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*bcDTOs.JobTeamDetailDTO, error)
	Create(currentUserID uint, jobTeamFormDTO *bcDTOs.JobTeamFormDTO) (uint, error)
	Update(id uint, currentUserID uint, jobTeam *bcDTOs.JobTeamFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type JobTeamService struct {
	repo           bcRepositories.JobTeamRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	designationSvc organizationServices.DesignationServiceInterface
}

func NewJobTeamService(
	repo bcRepositories.JobTeamRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	designationSvc organizationServices.DesignationServiceInterface,
) JobTeamServiceInterface {
	return &JobTeamService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		designationSvc: designationSvc,
	}
}

func (jts *JobTeamService) GetAll(params *bcDTOs.JobTeamListQueryParamsDTO) (*bcDTOs.JobTeamListResponseDTO, error) {
	jobTeams, totalCount, err := jts.repo.GetAll(params)
	if err != nil {
		return nil, err
	}

	// Map results to DTOs for response
	var jobTeamDTOs []*bcDTOs.JobTeamListDTO
	for _, jobTeam := range jobTeams {
		if jobTeam == nil {
			continue // Skip nil jobTeam
		}

		dto := &bcDTOs.JobTeamListDTO{
			ID:                          jobTeam.ID,
			TitleEn:                     jobTeam.TitleEn,
			TitleAr:                     jobTeam.TitleAr,
			TeamLeadDesignationID:       jobTeam.TeamLeadDesignationID,
			TeamLeadDesignationTitleEn:  getSafeString(jobTeam.TeamLeadDesignation, func(d *models.Designation) string { return d.TitleEn }),
			TeamLeadDesignationTitleAr:  getSafeString(jobTeam.TeamLeadDesignation, func(d *models.Designation) string { return d.TitleAr }),
			CreatedUserFirstNameEn:      getSafeString(jobTeam.CreatedByUser, func(u *models.User) string { return u.FirstNameEn }),
			CreatedUserFirstNameAr:      getSafeString(jobTeam.CreatedByUser, func(u *models.User) string { return u.FirstNameAr }),
			CreatedUserLastNameEn:       getSafeString(jobTeam.CreatedByUser, func(u *models.User) string { return u.LastNameEn }),
			CreatedUserLastNameAr:       getSafeString(jobTeam.CreatedByUser, func(u *models.User) string { return u.LastNameAr }),
			CreatedUserEmail:            getSafeString(jobTeam.CreatedByUser, func(u *models.User) string { return u.Email }),
			CreatedUserInitial:          getSafeString(jobTeam.CreatedByUser, func(u *models.User) string { return u.Initial }),
			CreatedUserColor:            getSafeString(jobTeam.CreatedByUser, func(u *models.User) string { return u.Color }),
			CreatedUserProfileImagePath: getSafeStringPointer(jobTeam.CreatedByUser, func(u *models.User) *string { return u.ProfileImagePath }),
			UpdatedUserFirstNameEn:      getSafeString(jobTeam.UpdatedByUser, func(u *models.User) string { return u.FirstNameEn }),
			UpdatedUserFirstNameAr:      getSafeString(jobTeam.UpdatedByUser, func(u *models.User) string { return u.FirstNameAr }),
			UpdatedUserLastNameEn:       getSafeString(jobTeam.UpdatedByUser, func(u *models.User) string { return u.LastNameEn }),
			UpdatedUserLastNameAr:       getSafeString(jobTeam.UpdatedByUser, func(u *models.User) string { return u.LastNameAr }),
			UpdatedUserEmail:            getSafeString(jobTeam.UpdatedByUser, func(u *models.User) string { return u.Email }),
			UpdatedUserInitial:          getSafeString(jobTeam.UpdatedByUser, func(u *models.User) string { return u.Initial }),
			UpdatedUserColor:            getSafeString(jobTeam.UpdatedByUser, func(u *models.User) string { return u.Color }),
			UpdatedUserProfileImagePath: getSafeStringPointer(jobTeam.UpdatedByUser, func(u *models.User) *string { return u.ProfileImagePath }),
			CreatedAt:                   jobTeam.CreatedAt.String(),
			UpdatedAt:                   getSafeString(jobTeam.UpdatedAt, func(t *time.Time) string { return t.String() }),
		}

		// Map team members
		for _, member := range jobTeam.JobTeamMembers {
			if member == nil || member.Designation == nil {
				continue // Skip nil members or designations
			}

			dto.Members = append(dto.Members, &bcDTOs.JobTeamMembersDTO{
				ID:            member.ID,
				JobTeamID:     member.JobTeamID,
				DesignationID: member.DesignationID,
				Designation: organizationDTOs.DesignationDTO{
					ID:      int(member.Designation.ID),
					TitleEn: member.Designation.TitleEn,
					TitleAr: member.Designation.TitleAr,
					Status:  member.Designation.Status,
				},
			})
		}

		jobTeamDTOs = append(jobTeamDTOs, dto)
	}
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	// Response with pagination
	response := &bcDTOs.JobTeamListResponseDTO{
		Items: jobTeamDTOs,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (jts JobTeamService) GetByID(id uint) (*models.JobTeam, error) {
	jobTeam, err := jts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return jobTeam, nil
}

func (jts JobTeamService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*bcDTOs.JobTeamDetailDTO, error) {
	var jobTeamDTO bcDTOs.JobTeamDetailDTO

	// Get full details of the job team
	jobTeam, err := jts.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Decode the job team into the DTO
	err = mapstructure.Decode(jobTeam, &jobTeamDTO)
	if err != nil {
		return nil, err
	}

	jobTeamDTO.CreatedAt = jobTeam.CreatedAt
	jobTeamDTO.UpdatedAt = jobTeam.UpdatedAt

	return &jobTeamDTO, nil
}

func (jts *JobTeamService) Create(currentUserID uint, jobTeamFormDTO *bcDTOs.JobTeamFormDTO) (uint, error) {
	tx := jts.repo.BeginTransaction()
	if err := jts.validateJobTeamFormDTO(jobTeamFormDTO); err != nil {
		return 0, err
	}

	jobTeam := &models.JobTeam{
		TitleEn:               jobTeamFormDTO.TitleEn,
		TitleAr:               jobTeamFormDTO.TitleAr,
		TeamLeadDesignationID: jobTeamFormDTO.TeamLeadDesignationID,
		Scope:                 jobTeamFormDTO.Scope,
		Objective:             jobTeamFormDTO.Objective,
		CreatedBYID:           currentUserID,
	}

	id, err := jts.repo.Create(jobTeam, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Save Job Team Members
	var members []*models.JobTeamMember
	for _, designationId := range jobTeamFormDTO.DesignationIDs {
		member := &models.JobTeamMember{
			JobTeamID:     id,
			DesignationID: *designationId,
			CreatedBYID:   currentUserID,
		}
		members = append(members, member)
	}

	if len(members) > 0 {
		if err := jts.repo.CreateMembers(members, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		existingJobTeam, err := jts.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch job team", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = jts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.JobTeam{}.TableName(),
			id,
			existingJobTeam,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.JobTeam{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (jts *JobTeamService) Update(id uint, currentUserID uint, jobTeamFormDTO *bcDTOs.JobTeamFormDTO) error {
	tx := jts.repo.BeginTransaction()

	oldJobTeam, err := jts.repo.GetFullDetailsByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	existingJobTeam, err := jts.GetByID(id)
	if err != nil {
		return err
	}

	if err := jts.validateJobTeamFormDTO(jobTeamFormDTO); err != nil {
		return err
	}

	// Update Job Team
	existingJobTeam.TitleEn = jobTeamFormDTO.TitleEn
	existingJobTeam.TitleAr = jobTeamFormDTO.TitleAr
	existingJobTeam.TeamLeadDesignationID = jobTeamFormDTO.TeamLeadDesignationID
	existingJobTeam.Scope = jobTeamFormDTO.Scope
	existingJobTeam.Objective = jobTeamFormDTO.Objective
	existingJobTeam.UpdatedBYID = &currentUserID

	err = jts.repo.Update(existingJobTeam, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Collect existing designation IDs
	existingJobTeamDetails, err := jts.repo.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingDesignationIDs := make(map[uint]bool)
	for _, member := range existingJobTeamDetails.JobTeamMembers {
		existingDesignationIDs[member.ID] = true
	}
	// Collect new designation IDs
	newDesignationIDs := make(map[uint]bool)
	for _, designationId := range jobTeamFormDTO.DesignationIDs {
		newDesignationIDs[*designationId] = true
	}
	// Add new members
	var members []*models.JobTeamMember
	for _, newDesignation := range jobTeamFormDTO.DesignationIDs {
		if _, exists := existingDesignationIDs[*newDesignation]; !exists {
			member := &models.JobTeamMember{
				JobTeamID:     existingJobTeamDetails.ID,
				DesignationID: *newDesignation,
				CreatedBYID:   currentUserID,
			}
			members = append(members, member)
		}
	}
	// Save the new members
	if len(members) > 0 {
		if err := jts.repo.CreateMembers(members, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove members not in the new list
	var deleteMemberIDs []*uint
	for _, existingMember := range existingJobTeamDetails.JobTeamMembers {
		if _, exists := newDesignationIDs[existingMember.ID]; !exists {
			// Delete the existing member from the repository
			deleteMemberIDs = append(deleteMemberIDs, &existingMember.ID)
		}
	}
	// Delete the members
	if len(deleteMemberIDs) > 0 {
		if err := jts.repo.DeleteMembers(deleteMemberIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	go func() {
		existingJobTeam, err := jts.repo.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		membersAdded := []uint{}
		membersRemoved := []uint{}
		for id := range newDesignationIDs {
			if !existingDesignationIDs[id] {
				membersAdded = append(membersAdded, id)
			}
		}
		for id := range existingDesignationIDs {
			if !newDesignationIDs[id] {
				membersRemoved = append(membersRemoved, id)
			}
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})

		if len(membersAdded) > 0 || len(membersRemoved) > 0 {
			associatedItems["designation_ids"] = map[string]interface{}{
				"added":   membersAdded,
				"removed": membersRemoved,
			}
		}
		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}
		// log changes
		oldData := map[string]interface{}{
			"title_en":  oldJobTeam.TitleEn,
			"title_ar":  oldJobTeam.TitleAr,
			"team_lead": oldJobTeam.TeamLeadDesignation.TitleEn,
		}

		newData := map[string]interface{}{
			"title_en":  existingJobTeam.TitleEn,
			"title_ar":  existingJobTeam.TitleAr,
			"team_lead": existingJobTeam.TeamLeadDesignation.TitleEn,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := jts.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.JobTeam{}.TableName(),
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
					"tableName": models.JobTeam{}.TableName(),
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

func (jts *JobTeamService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := jts.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedJobTeam []map[string]interface{}
	for _, id := range ids {
		jobTeam, err := jts.GetFullDetailsByID(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if len(jobTeam.JobTeamMembers) > 0 {
			var deleteMemberIds []*uint
			for _, member := range jobTeam.JobTeamMembers {
				deleteMemberIds = append(deleteMemberIds, &member.ID)
			}
			if err := jts.repo.DeleteMembers(deleteMemberIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		err = jts.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		successfulIds = append(successfulIds, id)

		deletedJobTeam = append(deletedJobTeam, map[string]interface{}{
			"id":       id,
			"title_en": jobTeam.TitleEn,
			"title_ar": jobTeam.TitleAr,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}
	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := jts.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.JobTeam{}.TableName(),
				successfulIds,
				deletedJobTeam,
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

func (jts *JobTeamService) validateJobTeamFormDTO(jobTeamFormDTO *bcDTOs.JobTeamFormDTO) error {
	designation, err := jts.designationSvc.GetByID(jobTeamFormDTO.TeamLeadDesignationID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("job team lead designation does not exist"), "not_exist", "team lead designation")
	}
	if designation.Status != "active" {
		return utility.NewValidationError(fmt.Errorf("job team lead designation is not active"), "not_active", "team lead designation")
	}
	// Validate Designation IDs
	designationIDs := jobTeamFormDTO.DesignationIDs
	designationIDMap := make(map[uint]bool)

	for _, id := range designationIDs {
		if designationIDMap[*id] {
			return utility.NewValidationError(fmt.Errorf("duplicate designation ID: %d", *id), "duplicate", "designation")
		}

		designation, err := jts.designationSvc.GetByID(*id)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("designation does not exist for ID: %d", *id), "not_exist", "designation")
		}
		if designation.Status != "active" {
			return utility.NewValidationError(fmt.Errorf("designation is not active for ID: %d", *id), "not_active", "designation")
		}

		designationIDMap[*id] = true
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
