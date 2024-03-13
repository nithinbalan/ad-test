package services

import (
	"fmt"
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	activityLogModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"jethur/utility"
	"log"
)

type SLAServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.SLARegisterWithAllField, error)
	GetByID(id uint) (*models.SLARegister, error)
	Create(currentUserID uint, slaRegister *models.SLARegister) (*models.SLARegister, error)
	Update(id uint, currentUserID uint, slaRegister *models.SLARegister) (*models.SLARegister, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type SLAService struct {
	repo           repositories.SLARepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewSLAService(repo repositories.SLARepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) SLAServiceInterface {
	return &SLAService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (srs *SLAService) GetAll(sort string, sortBy string, keyword string) ([]*models.SLARegisterWithAllField, error) {
	return srs.repo.GetAll(sort, sortBy, keyword)
}

func (srs SLAService) GetByID(id uint) (*models.SLARegister, error) {
	return srs.repo.GetByID(id)
}

func (srs *SLAService) Create(currentUserID uint, slaRegister *models.SLARegister) (*models.SLARegister, error) {

	if exists, err := srs.repo.ExistsByTitle(slaRegister.TitleEn, slaRegister.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
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

	if updatedSLA.TitleEn != existingSLA.TitleEn || updatedSLA.TitleAr != existingSLA.TitleAr {
		if exists, err := srs.repo.ExistsByTitleExcludingID(id, updatedSLA.TitleEn, updatedSLA.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

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



