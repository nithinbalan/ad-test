package services

import (
	"fmt"
	activityLogModel "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"jethur/grc/settings/models"
	"jethur/grc/settings/repositories"
	"jethur/utility"
	"log"
)

type OrganizationLevelServiceInterface interface {
	GetAll() ([]*models.OrganizationLevel, error)
	Create(currentUserID uint, organizationLevels []*models.OrganizationLevel) ([]*models.OrganizationLevel, error)
}

// Service and implementations
type OrganizationLevelService struct {
	repo           repositories.OrganizationLevelRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewOrganizationLevelService(repo repositories.OrganizationLevelRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) OrganizationLevelServiceInterface {
	return &OrganizationLevelService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ols *OrganizationLevelService) GetAll() ([]*models.OrganizationLevel, error) {
	return ols.repo.GetAll()
}

func (ols *OrganizationLevelService) Create(currentUserID uint, organizationLevels []*models.OrganizationLevel) ([]*models.OrganizationLevel, error) {
	titleEnMap := make(map[string]bool)
	titleArMap := make(map[string]bool)
	isPrimaryFound := false
	updatedIds := make(map[uint]bool)

	var toCreate, toUpdate []*models.OrganizationLevel

	for _, level := range organizationLevels {
		if level.ID == 0 {
			toCreate = append(toCreate, level)
		} else {
			toUpdate = append(toUpdate, level)
		}
	}

	existingOrganizationLevels, err := ols.repo.GetAll()

	if err != nil {
		return nil, err
	}

	for _, level := range organizationLevels {
		if titleEnMap[level.TitleEn] {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
		titleEnMap[level.TitleEn] = true

		if titleArMap[level.TitleAr] {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
		titleArMap[level.TitleAr] = true

		if level.IsPrimary {
			if isPrimaryFound {
				return nil, utility.NewValidationError(fmt.Errorf("only one object can have primary"), "unique", "primary")
			}
			isPrimaryFound = true
		}
	}

	if !isPrimaryFound {
		return nil, utility.NewValidationError(fmt.Errorf("at least one object must have primary true"), "required", "primary")
	}

	updatedOrganizationLevel, err := ols.repo.Create(organizationLevels)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "IsPrimary"}

	for _, newOrgLevel := range toCreate {
		go func(level *models.OrganizationLevel) {
			err := ols.activityLogSvc.LogActivity(
				activityLogModel.CreateAction,
				models.OrganizationLevel{}.TableName(),
				level.ID,
				level,
				currentUserID,
				nil,
			)

			if err != nil {
				log.Printf("Failed to log activity: %v", err)
			}
		}(newOrgLevel)
	}

	for _, updateOrgLevel := range toUpdate {
		var existingMatch *models.OrganizationLevel
		for _, existing := range existingOrganizationLevels {
			if existing.ID == updateOrgLevel.ID {
				existingMatch = existing
				break
			}
		}

		if existingMatch != nil {
			changes, err := utility.CompareFields(existingMatch, updateOrgLevel, fieldsToCompare)
			if err != nil {
				log.Printf("Failed to compare fields: %v", err)
				continue
			}

			if len(changes) > 0 {
				go func(level *models.OrganizationLevel, changes interface{}) {
					err := ols.activityLogSvc.LogActivity(
						activityLogModel.UpdateAction,
						models.OrganizationLevel{}.TableName(),
						level.ID,
						changes,
						currentUserID,
						nil,
					)

					if err != nil {
						log.Printf("Failed to log activity for update: %v", err)
					}
				}(updateOrgLevel, changes)
			}
		} else {
			log.Printf("No existing organization level found with ID: %d", updateOrgLevel.ID)
		}

		updatedIds[updateOrgLevel.ID] = true
	}

	for _, existingLevel := range existingOrganizationLevels {
		if !updatedIds[existingLevel.ID] {
			go func(level *models.OrganizationLevel) {
				err := ols.activityLogSvc.LogActivity(
					activityLogModel.DeleteAction,
					models.OrganizationLevel{}.TableName(),
					level.ID,
					level,
					currentUserID,
					nil,
				)

				if err != nil {
					log.Printf("Failed to log delete activity: %v", err)
				}
			}(existingLevel)
		}
	}

	return updatedOrganizationLevel, nil
}
