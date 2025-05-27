package services

import (
	"errors"
	"fmt"
	"strings"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	processDTOs "jethur/grc/organization/dtos"
	processRepositories "jethur/grc/organization/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProcessActivityServiceInterface interface {
	GetAll(processID uint, params *processDTOs.ProcessActivityListQueryParamsDTO) (*processDTOs.ProcessActivityListResponseDTO, error)
	GetByID(processID uint, id uint) (*models.ProcessActivity, error)
	GetFullDetailsByID(processID uint, id uint) (*processDTOs.ProcessActivityDetailDTO, error)
	Create(processID uint, currentUserID uint, form *processDTOs.ProcessActivityFormDTO) (uint, error)
	Update(processID uint, id uint, currentUserID uint, form *processDTOs.ProcessActivityFormDTO) error
	Delete(processID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type ProcessActivityService struct {
	repo           processRepositories.ProcessActivityRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	processSvc     ProcessServiceInterface
	userSvc        UserServiceInterface
}

func NewProcessActivityService(
	repo processRepositories.ProcessActivityRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	processSvc ProcessServiceInterface,
	userSvc UserServiceInterface,
) ProcessActivityServiceInterface {
	return &ProcessActivityService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		processSvc:     processSvc,
		userSvc:        userSvc,
	}
}

func (pas *ProcessActivityService) GetAll(processID uint, params *processDTOs.ProcessActivityListQueryParamsDTO) (*processDTOs.ProcessActivityListResponseDTO, error) {
	if _, err := pas.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pas.repo.GetAll(processID, params)
}

func (pas *ProcessActivityService) GetByID(processID uint, id uint) (*models.ProcessActivity, error) {
	if _, err := pas.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	activity, err := pas.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return activity, nil
}

func (pas *ProcessActivityService) GetFullDetailsByID(processID uint, id uint) (*processDTOs.ProcessActivityDetailDTO, error) {
	var dto processDTOs.ProcessActivityDetailDTO
	if _, err := pas.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	activity, err := pas.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	err = mapstructure.Decode(activity, &dto)
	if err != nil {
		return nil, err
	}
	dto.CreatedAt = activity.CreatedAt
	dto.UpdatedAt = activity.UpdatedAt
	return &dto, nil
}

func (pas *ProcessActivityService) Create(processID uint, currentUserID uint, form *processDTOs.ProcessActivityFormDTO) (uint, error) {
	if err := pas.validateProcessActivityFormDTO(form); err != nil {
		return 0, err
	}
	tx := pas.repo.BeginTransaction()

	if _, err := pas.processSvc.GetByID(processID); err != nil {
		tx.Rollback()
		return 0, err
	}

	activity := &models.ProcessActivity{
		ProcessID:         processID,
		TitleEn:           form.TitleEn,
		TitleAr:           form.TitleAr,
		DescriptionEn:     form.DescriptionEn,
		DescriptionAr:     form.DescriptionAr,
		Sequence:          form.Sequence,
		EstimatedDuration: form.EstimatedDuration,
		ExpectedOutput:    form.ExpectedOutput,
		CreatedBYID:       currentUserID,
	}

	id, err := pas.repo.Create(activity, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var userMappings []*models.ProcessActivityUser
	for _, responsibleID := range form.ResponsibleUserIDs {
		userMappings = append(userMappings, &models.ProcessActivityUser{
			ProcessActivityID: id,
			UserID:            responsibleID,
			RoleType:          "responsible",
		})
	}
	for _, accountableID := range form.AccountableUserIDs {
		userMappings = append(userMappings, &models.ProcessActivityUser{
			ProcessActivityID: id,
			UserID:            accountableID,
			RoleType:          "accountable",
		})
	}
	for _, consultedID := range form.ConsultedUserIDs {
		userMappings = append(userMappings, &models.ProcessActivityUser{
			ProcessActivityID: id,
			UserID:            consultedID,
			RoleType:          "consulted",
		})
	}
	if len(userMappings) > 0 {
		if err := pas.repo.CreateProcessActivityUsers(userMappings, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	go func() {
		existingActivity, err := pas.GetFullDetailsByID(processID, id)
		if err != nil {
			utility.LogError("Failed to fetch created activity for logging", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = pas.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ProcessActivity{}.TableName(),
			id,
			existingActivity,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ProcessActivity{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (pas *ProcessActivityService) Update(processID uint, id uint, currentUserID uint, form *processDTOs.ProcessActivityFormDTO) error {
	if err := pas.validateProcessActivityFormDTO(form); err != nil {
		return err
	}

	oldActivity, err := pas.GetFullDetailsByID(processID, id)
	if err != nil {
		return err
	}
	oldUserMap := groupUserDTOsByRole(oldActivity.Users)

	tx := pas.repo.BeginTransaction()

	existing, err := pas.GetByID(processID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existing.TitleEn = form.TitleEn
	existing.TitleAr = form.TitleAr
	existing.DescriptionEn = form.DescriptionEn
	existing.DescriptionAr = form.DescriptionAr
	existing.Sequence = form.Sequence
	existing.EstimatedDuration = form.EstimatedDuration
	existing.ExpectedOutput = form.ExpectedOutput
	existing.UpdatedBYID = &currentUserID

	if err := pas.repo.Update(existing, tx); err != nil {
		tx.Rollback()
		return err
	}

	newUserMap := map[string]map[uint]bool{
		"responsible": make(map[uint]bool),
		"accountable": make(map[uint]bool),
		"consulted":   make(map[uint]bool),
	}
	for _, id := range form.ResponsibleUserIDs {
		newUserMap["responsible"][id] = true
	}
	for _, id := range form.AccountableUserIDs {
		newUserMap["accountable"][id] = true
	}
	for _, id := range form.ConsultedUserIDs {
		newUserMap["consulted"][id] = true
	}
	existingMappings, err := pas.repo.GetActivityUsers(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingUserMap := make(map[string]map[uint]*uint)
	for _, u := range existingMappings {
		if existingUserMap[u.RoleType] == nil {
			existingUserMap[u.RoleType] = make(map[uint]*uint)
		}
		existingUserMap[u.RoleType][u.UserID] = &u.ID
	}

	var toInsert []*models.ProcessActivityUser
	var toDelete []*uint

	for role, newIDs := range newUserMap {
		existingIDs := existingUserMap[role]

		// Add new users
		for newID := range newIDs {
			if _, exists := existingIDs[newID]; !exists {
				toInsert = append(toInsert, &models.ProcessActivityUser{
					ProcessActivityID: id,
					UserID:            newID,
					RoleType:          role,
				})
			}
		}

		// Delete removed users
		for existingID, recordID := range existingIDs {
			if _, stillExists := newIDs[existingID]; !stillExists {
				toDelete = append(toDelete, recordID)
			}
		}
	}
	// Convert []*uint to []uint
	var deleteIDs []uint
	for _, idPtr := range toDelete {
		if idPtr != nil {
			deleteIDs = append(deleteIDs, *idPtr)
		}
	}

	if len(deleteIDs) > 0 {
		if err := pas.repo.DeleteProcessActivityUsers(deleteIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(toInsert) > 0 {
		if err := pas.repo.CreateProcessActivityUsers(toInsert, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingActivity, err := pas.GetFullDetailsByID(processID, id)
		if err != nil {
			utility.LogError("Failed to fetch updated activity for logging", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		newUserMap := groupUserDTOsByRole(existingActivity.Users)
		// Build oldData and newData for comparison
		oldData := map[string]interface{}{
			"title_en":        oldActivity.TitleEn,
			"title_ar":        oldActivity.TitleAr,
			"description_en":  oldActivity.DescriptionEn,
			"description_ar":  oldActivity.DescriptionAr,
			"sequence":        oldActivity.Sequence,
			"expected_output": oldActivity.ExpectedOutput,
			"estimated_duration": func() string {
				if oldActivity.EstimatedDuration != nil {
					return *oldActivity.EstimatedDuration
				}
				return ""
			}(),
		}

		newData := map[string]interface{}{
			"title_en":        existingActivity.TitleEn,
			"title_ar":        existingActivity.TitleAr,
			"description_en":  existingActivity.DescriptionEn,
			"description_ar":  existingActivity.DescriptionAr,
			"sequence":        existingActivity.Sequence,
			"expected_output": existingActivity.ExpectedOutput,
			"estimated_duration": func() string {
				if existingActivity.EstimatedDuration != nil {
					return *existingActivity.EstimatedDuration
				}
				return ""
			}(),
		}

		// Compare and collect changed fields
		changedData := utility.CompareFieldChanges(oldData, newData)

		// Compare user mappings (if fetched)
		if oldUserMap != nil && newUserMap != nil {
			roleTypes := []string{"responsible", "accountable", "consulted"}

			for _, role := range roleTypes {
				added := difference(newUserMap[role], oldUserMap[role])
				removed := difference(oldUserMap[role], newUserMap[role])

				if len(added) > 0 || len(removed) > 0 {
					changedData[role+"_user_ids"] = map[string]interface{}{
						"added":   added,
						"removed": removed,
					}
				}
			}
		}

		// Log only if any changes exist
		if len(changedData) > 0 {
			err := pas.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ProcessActivity{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log update activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"table":    models.ProcessActivity{}.TableName(),
					"recordID": id,
					"error":    err.Error(),
				})
			}
		}
	}()

	return nil
}

func (pas *ProcessActivityService) Delete(processID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := pas.processSvc.GetByID(processID); err != nil {
		return nil, err
	}

	var failed []uint
	var successful []uint
	var deletedActivities []*models.ProcessActivity

	for _, id := range ids {
		// Fetch full detail for log before delete
		activity, err := pas.repo.GetFullDetailsByID(id)
		if err != nil {
			failed = append(failed, id)
			continue
		}
		deletedActivities = append(deletedActivities, activity)

		if err := pas.repo.DeleteProcessActivityUsersByActivityID(id); err != nil {
			failed = append(failed, id)
			continue
		}
		if err := pas.repo.Delete(id); err != nil {
			failed = append(failed, id)
		} else {
			successful = append(successful, id)
		}
	}

	if len(successful) > 0 {
		go func() {
			err := pas.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ProcessActivity{}.TableName(),
				successful,
				deletedActivities,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete activity for process activities", map[string]interface{}{
					"action":    models.DeleteAction,
					"recordIDs": successful,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}()
	}

	return failed, nil
}

func groupUserDTOsByRole(users []*processDTOs.ProcessActivityUserDTO) map[string][]uint {
	roleMap := map[string][]uint{}
	for _, u := range users {
		role := strings.ToLower(u.RoleType)
		roleMap[role] = append(roleMap[role], u.UserID)
	}
	return roleMap
}

func difference(slice1, slice2 []uint) []uint {
	diff := []uint{}
	exists := make(map[uint]bool)

	for _, val := range slice2 {
		exists[val] = true
	}

	for _, val := range slice1 {
		if !exists[val] {
			diff = append(diff, val)
		}
	}

	return diff
}

func (pas *ProcessActivityService) validateProcessActivityFormDTO(form *processDTOs.ProcessActivityFormDTO) error {
	userIDSet := make(map[uint]bool)
	allUserIDs := append(append(form.ResponsibleUserIDs, form.AccountableUserIDs...), form.ConsultedUserIDs...)

	for _, userID := range allUserIDs {
		if userIDSet[userID] {
			return utility.NewValidationError(fmt.Errorf("duplicate user ID: %d", userID), "duplicate", "user_ids")
		}
		if _, err := pas.userSvc.GetByID(userID); err != nil {
			return utility.NewValidationError(fmt.Errorf("user does not exist for ID: %d", userID), "not_exist", "user_ids")
		}
		userIDSet[userID] = true
	}

	return nil
}
