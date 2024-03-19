package services

import (
	"fmt"
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	activityLogModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	"jethur/utility"
	"log"
)

type SLAServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.SLARegister, error)
	GetByID(id uint) (*models.SLARegister, error)
	Create(currentUserID uint, slaRegister *models.SLARegister) (*models.SLARegister, error)
	Update(id uint, currentUserID uint, slaRegister *models.SLARegister) (*models.SLARegister, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type SLAService struct {
	repo           repositories.SLARepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
	slaCategoryRepo   repositories.SLACategoryRepositoryInterface
	severityLevelRepo masterRepositories.SeverityLevelRepositoryInterface
	lifecycleStatusRepo masterRepositories.LifecycleStatusRepositoryInterface
	
}

func NewSLAService(repo repositories.SLARepositoryInterface,slaCategoryRepo repositories.SLACategoryRepositoryInterface,severityLevelRepo masterRepositories.SeverityLevelRepositoryInterface,lifecycleStatusRepo masterRepositories.LifecycleStatusRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) SLAServiceInterface {
	return &SLAService{
		repo:           repo,
		slaCategoryRepo: slaCategoryRepo,
		severityLevelRepo: severityLevelRepo,
		lifecycleStatusRepo: lifecycleStatusRepo,
		activityLogSvc: activityLogSvc,
	}
}

func (srs *SLAService) GetAll(sort string, sortBy string, keyword string) ([]*models.SLARegister, error) {
	return srs.repo.GetAll(sort, sortBy, keyword)
}

func (srs SLAService) GetByID(id uint) (*models.SLARegister, error) {
	return srs.repo.GetByID(id)
}

func (srs *SLAService) Create(currentUserID uint, slaRegister *models.SLARegister) (*models.SLARegister, error) {
	
	// Validation for sla category field
	if _, err := srs.slaCategoryRepo.GetByID(slaRegister.SLACategoryID); err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("sla category does not exist"), "not_exist", "sla category")
	}
	// Validation for Severity Level
    if _, err := srs.severityLevelRepo.GetByID(slaRegister.SeverityLevelID); err != nil {
        return nil, utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
    }
	// Validation for lifecycle status 
    if _, err := srs.lifecycleStatusRepo.GetByID(slaRegister.LifecycleStatusID); err != nil {
        return nil, utility.NewValidationError(fmt.Errorf("status does not exist"), "not_exist", "lifecycle_status")
    }
	createdSLA, err := srs.repo.Create(slaRegister)
	if err != nil {
		return nil, err
	}

	
	go func() {
		err := srs.activityLogSvc.LogActivity(
			activityLogModels.CreateAction,
			models.SLARegister{}.TableName(),
			slaRegister.ID,
			createdSLA,
			currentUserID,
			nil,
		)

		if err != nil {
			//save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdSLA, nil
}

func (srs *SLAService) Update(id uint, currentUserID uint, updatedSLA *models.SLARegister) (*models.SLARegister, error) {
	existingSLA, err := srs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	// Validation for sla category field
	if _, err := srs.slaCategoryRepo.GetByID(updatedSLA.SLACategoryID); err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("sla category does not exist"), "not_exist", "sla category")
	}
	// Validation for Severity Level
    if _, err := srs.severityLevelRepo.GetByID(updatedSLA.SeverityLevelID); err != nil {
        return nil, utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
    }
	// Validation for lifecycle status 
    if _, err := srs.lifecycleStatusRepo.GetByID(updatedSLA.LifecycleStatusID); err != nil {
        return nil, utility.NewValidationError(fmt.Errorf("status does not exist"), "not_exist", "lifecycle_status")
    }
	result, err := srs.repo.Update(id, updatedSLA)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{
		"TitleEn",
		"TitleAr",
		"SLACategoryID",   
		"SeverityLevelID", 
		"ResponsibleUnitID",
		"ResponsibleUserID",
		"Description",
		"Client",
	}
	changes, err := utility.CompareFields(existingSLA, updatedSLA, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := srs.activityLogSvc.LogActivity(
			activityLogModels.UpdateAction,
			models.SLARegister{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()
	return result, err
}

func (srs *SLAService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := srs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = srs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := srs.activityLogSvc.LogBulkActivities(
				activityLogModels.DeleteAction,
				models.SLARegister{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}



