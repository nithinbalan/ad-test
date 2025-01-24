package services

import (
	"errors"
	"fmt"

	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	businessFunctionRepositories "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BcpTeamServiceInterface interface {
	GetAll(bcpId uint, params *businessContinuityDTOs.BcpTeamsListQueryParamsDTO) ([]*businessContinuityDTOs.BcpTeamsListDTO, error)
	GetByID(id uint) (*models.BcpTeam, error)
	GetFullDetailsByID(bcpId uint, id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BcpTeamsDetailDTO, error)
	Create(bcpId uint, currentUserID uint, bcpTeamForm *businessContinuityDTOs.BcpTeamsFormDTO) (uint, error)
	Update(id uint, currentUserID uint, bcpTeamForm *businessContinuityDTOs.BcpTeamsFormDTO) error
	Delete(bcpId uint, TeamIDs []uint, currentUserID uint) ([]uint, error)
}

type BcpTeamService struct {
	repo                      businessFunctionRepositories.BcpTeamRepositoryInterface
	userSvc                   organizationServices.UserServiceInterface
	businessContinuityPlanSvc BusinessContinuityPlanServiceInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
}

func NewBcpTeamService(
	repo businessFunctionRepositories.BcpTeamRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	businessContinuityPlanSvc BusinessContinuityPlanServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) BcpTeamServiceInterface {
	return &BcpTeamService{
		repo:                      repo,
		userSvc:                   userSvc,
		businessContinuityPlanSvc: businessContinuityPlanSvc,
		activityLogSvc:            activityLogSvc,
	}
}

func (brs *BcpTeamService) GetAll(bcpId uint, params *businessContinuityDTOs.BcpTeamsListQueryParamsDTO) ([]*businessContinuityDTOs.BcpTeamsListDTO, error) {
	if _, err := brs.businessContinuityPlanSvc.GetByID(bcpId); err != nil {
		return nil, err
	}
	return brs.repo.GetAll(bcpId, params)
}

func (brs *BcpTeamService) GetByID(id uint) (*models.BcpTeam, error) {
	team, err := brs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return team, nil
}

func (brs *BcpTeamService) GetFullDetailsByID(bcpId uint, id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BcpTeamsDetailDTO, error) {
	if _, err := brs.businessContinuityPlanSvc.GetByID(bcpId); err != nil {
		return nil, err
	}
	var bcpTeamDTO businessContinuityDTOs.BcpTeamsDetailDTO

	bcpTeam, err := brs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if err := mapstructure.Decode(bcpTeam, &bcpTeamDTO); err != nil {
		return nil, err
	}

	bcpTeamDTO.CreatedAt = bcpTeam.CreatedAt
	bcpTeamDTO.UpdatedAt = bcpTeam.UpdatedAt

	return &bcpTeamDTO, nil
}

func (brs *BcpTeamService) Create(bcpId uint, currentUserID uint, bcpTeamForm *businessContinuityDTOs.BcpTeamsFormDTO) (uint, error) {
	tx := brs.repo.BeginTransaction()
	if _, err := brs.businessContinuityPlanSvc.GetByID(bcpId); err != nil {
		return 0, err
	}

	organizationUnitID, err := brs.userSvc.GetUserPrimaryUnitID(uint(bcpTeamForm.ResponsibleUserID))
	if err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible_user_unit")
	}

	bcpTeam := &models.BcpTeam{
		BusinessContinuityPlanID: bcpId,
		ResponsibleUserID:        bcpTeamForm.ResponsibleUserID,
		ResponsibleUnitID:        organizationUnitID,
		Responsibility:           bcpTeamForm.Responsibility,
		BackupUserID:             bcpTeamForm.BackupUserID,
		CreatedBYID:              currentUserID,
	}

	id, err := brs.repo.Create(bcpTeam, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		createdTeam, err := brs.GetFullDetailsByID(bcpId, id)
		if err != nil {
			utility.LogError("Failed to fetch created BCP team details", map[string]interface{}{"error": err})
			return
		}

		err = brs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.BcpTeam{}.TableName(),
			id,
			createdTeam,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity for BCP team creation", map[string]interface{}{"error": err})
		}
	}()
	return id, nil
}

func (brs *BcpTeamService) Update(id uint, currentUserID uint, bcpTeamForm *businessContinuityDTOs.BcpTeamsFormDTO) error {
	tx := brs.repo.BeginTransaction()

	existingTeam, err := brs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	oldTeamDetails, err := brs.GetFullDetailsByID(existingTeam.BusinessContinuityPlanID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingTeam.ResponsibleUserID != bcpTeamForm.ResponsibleUserID {
		organizationUnitID, err := brs.userSvc.GetUserPrimaryUnitID(bcpTeamForm.ResponsibleUserID)
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible_user_unit")
		}
		existingTeam.ResponsibleUserID = bcpTeamForm.ResponsibleUserID
		existingTeam.ResponsibleUnitID = organizationUnitID
	}

	existingTeam.Responsibility = bcpTeamForm.Responsibility
	existingTeam.BackupUserID = bcpTeamForm.BackupUserID

	if err := brs.repo.Update(existingTeam, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		newTeamDetails, err := brs.GetFullDetailsByID(existingTeam.BusinessContinuityPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch updated BCP team details", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		oldData := map[string]interface{}{
			"responsibility":         oldTeamDetails.Responsibility,
			"responsible_user":       fmt.Sprintf("%s %s", oldTeamDetails.ResponsibleUser.FirstNameEN, oldTeamDetails.ResponsibleUser.LastNameEN),
			"backup_user":            fmt.Sprintf("%s %s", oldTeamDetails.BackupUser.FirstNameEN, oldTeamDetails.BackupUser.LastNameEN),
			"responsible_unit_title": oldTeamDetails.ResponsibleUnit.TitleEn,
		}

		newData := map[string]interface{}{
			"responsibility":         newTeamDetails.Responsibility,
			"responsible_user":       fmt.Sprintf("%s %s", newTeamDetails.ResponsibleUser.FirstNameEN, newTeamDetails.ResponsibleUser.LastNameEN),
			"backup_user":            fmt.Sprintf("%s %s", newTeamDetails.BackupUser.FirstNameEN, newTeamDetails.BackupUser.LastNameEN),
			"responsible_unit_title": newTeamDetails.ResponsibleUnit.TitleEn,
		}

		changedData := utility.CompareFieldChanges(oldData, newData)
		if len(changedData) > 0 {
			err := brs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.BcpTeam{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity for BCP team update", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.BcpTeam{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}
	}()

	return nil
}

func (brs *BcpTeamService) Delete(bcpId uint, teamIDs []uint, currentUserID uint) ([]uint, error) {
	tx := brs.repo.BeginTransaction()

	if _, err := brs.businessContinuityPlanSvc.GetByID(bcpId); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var deletedTeams []map[string]interface{}

	for _, teamID := range teamIDs {
		team, err := brs.repo.GetByID(teamID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, teamID)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		if team.BusinessContinuityPlanID != bcpId {
			failedIDs = append(failedIDs, teamID)
			continue
		}

		deletedTeams = append(deletedTeams, map[string]interface{}{
			"id":               team.ID,
			"responsibility":   team.Responsibility,
			"responsible_user": fmt.Sprintf("%s %s", team.ResponsibleUser.FirstNameEn, team.ResponsibleUser.FirstNameEn),
		})

		if err := brs.repo.Delete(teamID, tx); err != nil {
			failedIDs = append(failedIDs, teamID)
			continue
		}
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()

	go func() {
		if len(deletedTeams) > 0 {
			err := brs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.BcpTeam{}.TableName(),
				teamIDs,
				deletedTeams,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity for BCP team deletion", map[string]interface{}{"error": err})
			}
		}
	}()

	return failedIDs, nil
}
