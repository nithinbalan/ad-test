package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	settingServices "jethur/grc/settings/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyProfileRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileServiceInterface interface {
	GetAll(params *strategyDTOs.StrategyProfileListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileListDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileDetailDTO, error)
	GetByID(id uint) (*models.StrategyProfile, error)
	Create(currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileFormDTO) (uint, error)
	Update(id uint, currentUserID uint, strategyProfile *strategyDTOs.StrategyProfileFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type StrategyProfileService struct {
	repo                strategyProfileRepositories.StrategyProfileRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	userSvc             organizationServices.UserServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
	organizationUnitSvc settingServices.OrganizationUnitServiceInterface
}

func NewStrategyProfileService(
	repo strategyProfileRepositories.StrategyProfileRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	organizationUnitSvc settingServices.OrganizationUnitServiceInterface,
) StrategyProfileServiceInterface {
	return &StrategyProfileService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		userSvc:             userSvc,
		userSettingSvc:      userSettingSvc,
		organizationUnitSvc: organizationUnitSvc,
	}
}

func (sps *StrategyProfileService) GetAll(params *strategyDTOs.StrategyProfileListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileListDTO, error) {
	return sps.repo.GetAll(params)
}

func (sps StrategyProfileService) GetByID(id uint) (*models.StrategyProfile, error) {
	strategyProfile, err := sps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategyProfile, nil
}

func (sps StrategyProfileService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileDetailDTO, error) {
	// Initialize your DTO
	var strategyProfileDTO strategyDTOs.StrategyProfileDetailDTO

	strategyProfile, err := sps.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(strategyProfile, &strategyProfileDTO)
	if err != nil {
		return nil, err
	}

	strategyProfileDTO.StartDateString = strategyProfile.StartDate.Format("2006-01-02")
	strategyProfileDTO.EndDateString = strategyProfile.EndDate.Format("2006-01-02")

	strategyProfileDTO.CreatedAt = strategyProfile.CreatedAt
	strategyProfileDTO.UpdatedAt = strategyProfile.UpdatedAt

	return &strategyProfileDTO, nil
}

func (sps *StrategyProfileService) Create(currentUserID uint, strategyProfileFormDTO *strategyDTOs.StrategyProfileFormDTO) (uint, error) {

	tx := sps.repo.BeginTransaction()

	dates, err := sps.parseStrategyProfileDates(strategyProfileFormDTO)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := sps.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	strategyProfile := &models.StrategyProfile{
		ReferenceNumber:   referenceNumber,
		Title:             strategyProfileFormDTO.Title,
		Description:       strategyProfileFormDTO.Description,
		Vision:            strategyProfileFormDTO.Vision,
		Mission:           strategyProfileFormDTO.Mission,
		ResponsibleUserID: strategyProfileFormDTO.ResponsibleUserID,
		StartDate:         dates.StartDate,
		EndDate:           dates.EndDate,
		CreatedBYID:       currentUserID,
	}
	id, err := sps.repo.Create(strategyProfile, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingStrategyProfile, err := sps.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch strategy profile", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = sps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StrategyProfile{}.TableName(),
			id,
			existingStrategyProfile,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StrategyProfile{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (sps *StrategyProfileService) Update(id uint, currentUserID uint, strategyProfileFormDTO *strategyDTOs.StrategyProfileFormDTO) error {
	tx := sps.repo.BeginTransaction()

	existingStrategyProfile, err := sps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldStrategyProfile, err := sps.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	dates, err := sps.parseStrategyProfileDates(strategyProfileFormDTO)
	if err != nil {
		return err
	}

	existingStrategyProfile.Title = strategyProfileFormDTO.Title
	existingStrategyProfile.Description = strategyProfileFormDTO.Description
	existingStrategyProfile.Vision = strategyProfileFormDTO.Vision
	existingStrategyProfile.Mission = strategyProfileFormDTO.Mission
	existingStrategyProfile.ResponsibleUserID = strategyProfileFormDTO.ResponsibleUserID
	existingStrategyProfile.StartDate = dates.StartDate
	existingStrategyProfile.EndDate = dates.EndDate
	existingStrategyProfile.UpdatedBYID = &currentUserID

	err = sps.repo.Update(existingStrategyProfile, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingStrategyProfile, err := sps.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"title":            oldStrategyProfile.Title,
			"description":      oldStrategyProfile.Description,
			"vision":           oldStrategyProfile.Vision,
			"mission":          oldStrategyProfile.Mission,
			"responsible_user": utility.GetFieldValue(oldStrategyProfile.ResponsibleUser),
			"start_date":       oldStrategyProfile.StartDateString,
			"end_date":         oldStrategyProfile.EndDateString,
		}

		newData := map[string]interface{}{
			"title":            existingStrategyProfile.Title,
			"description":      existingStrategyProfile.Description,
			"vision":           existingStrategyProfile.Vision,
			"mission":          existingStrategyProfile.Mission,
			"responsible_user": utility.GetFieldValue(existingStrategyProfile.ResponsibleUser),
			"start_date":       existingStrategyProfile.StartDateString,
			"end_date":         existingStrategyProfile.EndDateString,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := sps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.StrategyProfile{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StrategyProfile{}.TableName(),
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

func (sps *StrategyProfileService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := sps.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedStrategyProfile []map[string]interface{}

	for _, id := range ids {
		strategyProfile, err := sps.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedStrategyProfile = append(deletedStrategyProfile, map[string]interface{}{
			"id":       id,
			"title_en": strategyProfile.Title,
		})

		// Delete the business
		if err := sps.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := sps.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.StrategyProfile{}.TableName(),
				successfulIds,
				deletedStrategyProfile,
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

func (sps *StrategyProfileService) parseStrategyProfileDates(dto *strategyDTOs.StrategyProfileFormDTO) (*strategyDTOs.StrategyParsedDates, error) {
	parsed := &strategyDTOs.StrategyParsedDates{}

	var err error
	if dto.StartDate != "" {
		parsed.StartDate, err = utility.ParseDate(dto.StartDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
		}
	}

	if dto.EndDate != "" {
		parsed.EndDate, err = utility.ParseDate(dto.EndDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid end date"), "invalid_date", "end date")
		}
	}

	return parsed, nil
}
