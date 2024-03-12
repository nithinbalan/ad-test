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

type SLARegisterServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.SLARegister, error)
	GetByID(id uint) (*models.SLARegister, error)
	Create(currentUserID uint, slaRegister *models.SLARegister) (*models.SLARegister, error)
	Update(id uint, currentUserID uint, slaRegister *models.SLARegister) (*models.SLARegister, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type SLARegisterService struct {
	repo           repositories.SLARegisterRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewSLARegisterService(repo repositories.SLARegisterRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) SLARegisterServiceInterface {
	return &SLARegisterService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (srs *SLARegisterService) GetAll(sort string, sortBy string, keyword string) ([]*models.SLARegister, error) {
	return srs.repo.GetAll(sort, sortBy, keyword)
}

func (srs SLARegisterService) GetByID(id uint) (*models.SLARegister, error) {
	return srs.repo.GetByID(id)
}

func (srs *SLARegisterService) Create(currentUserID uint, slaRegister *models.SLARegister) (*models.SLARegister, error) {

	if exists, err := srs.repo.ExistsByTitle(slaRegister.TitleEn, slaRegister.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	} 
	
	createdSLARegister, err := srs.repo.Create(slaRegister)
	if err != nil {
		return nil, err
	}

	go func() {
		err := srs.activityLogSvc.LogActivity(
			activityLogModels.CreateAction,
			models.SLARegister{}.TableName(),
			slaRegister.ID,
			createdSLARegister,
			currentUserID,
			nil,
		)

		if err != nil {
			//save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdSLARegister, nil
}

func (srs *SLARegisterService) Update(id uint, currentUserID uint, updatedSLARegister *models.SLARegister) (*models.SLARegister, error) {
	existingSLARegister, err := srs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedSLARegister.TitleEn != existingSLARegister.TitleEn || updatedSLARegister.TitleAr != existingSLARegister.TitleAr {
		if exists, err := srs.repo.ExistsByTitleExcludingID(id, updatedSLARegister.TitleEn, updatedSLARegister.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := srs.repo.Update(id, updatedSLARegister)

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
	changes, err := utility.CompareFields(existingSLARegister, updatedSLARegister, fieldsToCompare)

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

func (srs *SLARegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
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



