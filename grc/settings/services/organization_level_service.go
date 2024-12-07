package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/settings/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type OrganizationLevelServiceInterface interface {
	GetAll() ([]*models.OrganizationLevel, error)
	GetByID(id uint) (*models.OrganizationLevel, error)
	CreateOrganizationLevels(currentUserID uint, organizationLevels []*models.OrganizationLevel) ([]*models.OrganizationLevel, error)
	Create(currentUserID uint, orgLevel *models.OrganizationLevel) (*models.OrganizationLevel, error)
	Update(id uint, currentUserID uint, updatedOrgLevel *models.OrganizationLevel) (*models.OrganizationLevel, error)
	Delete(id uint, currentUserID uint) error

	GetPrimaryLevel() (*models.OrganizationLevel, error)
}

// Service and implementations
type OrganizationLevelService struct {
	repo           repositories.OrganizationLevelRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewOrganizationLevelService(repo repositories.OrganizationLevelRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface) OrganizationLevelServiceInterface {
	return &OrganizationLevelService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ols *OrganizationLevelService) GetAll() ([]*models.OrganizationLevel, error) {
	return ols.repo.GetAll()
}

func (ous OrganizationLevelService) GetByID(id uint) (*models.OrganizationLevel, error) {
	orgLevel, err := ous.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return orgLevel, nil
}

func (ols *OrganizationLevelService) CreateOrganizationLevels(currentUserID uint, organizationLevels []*models.OrganizationLevel) ([]*models.OrganizationLevel, error) {
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

	updatedOrganizationLevel, err := ols.repo.CreateOrganizationLevels(organizationLevels)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "IsPrimary"}

	for _, newOrgLevel := range toCreate {
		go func(level *models.OrganizationLevel) {
			err := ols.activityLogSvc.LogActivity(
				models.CreateAction,
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
						models.UpdateAction,
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
					models.DeleteAction,
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

func (ols *OrganizationLevelService) GetPrimaryLevel() (*models.OrganizationLevel, error) {
	return ols.repo.GetPrimaryLevel()
}

func (ols *OrganizationLevelService) Create(currentUserID uint, orgLevel *models.OrganizationLevel) (*models.OrganizationLevel, error) {

	if err := ols.validateOrganizationLevelForm(orgLevel, nil); err != nil {
		return nil, err
	}

	createdOrgLevel, err := ols.repo.Create(orgLevel)
	if err != nil {
		return nil, err
	}

	go func() {
		err := ols.activityLogSvc.LogActivity(
			models.CreateAction,
			models.OrganizationLevel{}.TableName(),
			orgLevel.ID,
			createdOrgLevel,
			currentUserID,
			nil,
		)

		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdOrgLevel, nil
}

func (ols *OrganizationLevelService) Update(id uint, currentUserID uint, updatedOrgLevel *models.OrganizationLevel) (*models.OrganizationLevel, error) {
	existingOrgLevel, err := ols.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := ols.validateOrganizationLevelForm(updatedOrgLevel, &id); err != nil {
		return nil, err
	}

	result, err := ols.repo.Update(id, updatedOrgLevel)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Level", "IsPrimary"}
	changes, err := utility.CompareFields(existingOrgLevel, updatedOrgLevel, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := ols.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.OrganizationLevel{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()

	return result, nil
}

func (ols *OrganizationLevelService) Delete(id uint, currentUserID uint) error {
	tx := ols.repo.BeginTransaction()

	orgLevel, err := ols.GetByID(id)
	if err != nil {
		return err
	}

	if orgLevel.IsPrimary {
		return fmt.Errorf("cannot delete a primary organization level")
	}

	var unitIDs []uint
	for _, unit := range orgLevel.OrganizationUnits {
		unitIDs = append(unitIDs, unit.ID)
	}

	// Delete all related organization units if any
	if len(unitIDs) > 0 {
		if err := ols.repo.DeleteOrganizationUnits(unitIDs, tx); err != nil {
			tx.Rollback()
			if utility.IsForeignKeyConstraintError(err) {
				return utility.NewResourceConflictError(fmt.Errorf("foreign key constraint violation while deleting organization units"))
			}
			return err
		}
	}

	err = ols.repo.Delete(id, tx)
	if err != nil {
		tx.Rollback()
		if utility.IsForeignKeyConstraintError(err) {
			return utility.NewResourceConflictError(fmt.Errorf("foreign key constraint violation while deleting organization level"))
		}
		return err
	}

	tx.Commit()

	// Log the activity
	go func() {
		err := ols.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.OrganizationLevel{}.TableName(),
			id,
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()

	return nil
}

func (ols *OrganizationLevelService) validateOrganizationLevelForm(orgLevel *models.OrganizationLevel, id *uint) error {
	if id == nil {
		if exists, err := ols.repo.ExistsByTitle(orgLevel.TitleEn, orgLevel.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

		if exists, err := ols.repo.ExistsByLevel(orgLevel.Level); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the level must be unique"), "unique", "level")
		}

	} else {
		if exists, err := ols.repo.ExistsByTitleExcludingID(*id, orgLevel.TitleEn, orgLevel.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

		if exists, err := ols.repo.ExistsByLevelExcludingID(*id, orgLevel.Level); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the level must be unique"), "unique", "level")
		}
	}
	return nil
}
