package services

import (
	"fmt"
	slaDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	activityLogModels "jethur/grc/logs/models"
	logModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	"jethur/utility"
	"log"
)

type SLARegisterServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, page int, pageSize int) (*slaDTOs.SLARegisterListResponseDTO, error)
	GetByID(id uint) (*models.SLARegister, error)
	Create(currentUserID uint, slaRegFormDTO *slaDTOs.SLARegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, sla *slaDTOs.SLARegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type SLARegisterService struct {
	repo                repositories.SLARegisterRepositoryInterface
	slaCategoryRepo     repositories.SLACategoryRepositoryInterface
	severityLevelRepo   masterRepositories.SeverityLevelRepositoryInterface
	lifecycleStatusRepo masterRepositories.LifecycleStatusRepositoryInterface
	activityLogSvc      services.ActivityLogServiceInterface
}

func NewSLARegisterService(
	repo repositories.SLARegisterRepositoryInterface,
	slaCategoryRepo repositories.SLACategoryRepositoryInterface,
	severityLevelRepo masterRepositories.SeverityLevelRepositoryInterface,
	lifecycleStatusRepo masterRepositories.LifecycleStatusRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface) SLARegisterServiceInterface {
	return &SLARegisterService{
		repo:                repo,
		slaCategoryRepo:     slaCategoryRepo,
		severityLevelRepo:   severityLevelRepo,
		lifecycleStatusRepo: lifecycleStatusRepo,
		activityLogSvc:      activityLogSvc,
	}
}

func (srs *SLARegisterService) GetAll(sort string, sortBy string, keyword string, status string, page int, pageSize int) (*slaDTOs.SLARegisterListResponseDTO, error) {
	return srs.repo.GetAll(sort, sortBy, keyword,status, page, pageSize)
}

func (srs SLARegisterService) GetByID(id uint) (*models.SLARegister, error) {
	return srs.repo.GetByID(id)
}

func (srs *SLARegisterService) Create(currentUserID uint, slaRegFormDTO *slaDTOs.SLARegisterFormDTO) (uint, error) {

	// Validation for sla category field
	slaCategory, err := srs.slaCategoryRepo.GetByID(slaRegFormDTO.SLACategoryID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("sla category does not exist"), "not_exist", "sla category")
	}
	// Validation for Severity Level
	severityLevel, err := srs.severityLevelRepo.GetByID(slaRegFormDTO.SeverityLevelID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}
	// Validation for lifecycle status
	lifecycleStatus, err := srs.lifecycleStatusRepo.GetByID(slaRegFormDTO.LifecycleStatusID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("status does not exist"), "not_exist", "lifecycle_status")
	}
	_, _, _ = slaCategory, severityLevel, lifecycleStatus

	// id, err := srs.repo.Create(currentUserID, slaRegFormDTO)
	id, err := srs.repo.Create(currentUserID, slaRegFormDTO)
	if err != nil {
		return 0, err
	}

	go func() {
		err := srs.activityLogSvc.LogActivity(
			activityLogModels.CreateAction,
			models.SLARegister{}.TableName(),
			id,
			slaRegFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return id, nil
}

func (srs *SLARegisterService) Update(id uint, currentUserID uint, slaRegFormDTO *slaDTOs.SLARegisterFormDTO) error {
	existingSLARegister, err := srs.repo.GetByID(id)

	if err != nil {
		return err
	}
	// Validation for sla category field
	slaCategory, err := srs.slaCategoryRepo.GetByID(slaRegFormDTO.SLACategoryID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("sla category does not exist"), "not_exist", "sla category")
	}
	// Validation for Severity Level
	severityLevel, err := srs.severityLevelRepo.GetByID(slaRegFormDTO.SeverityLevelID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}
	// Validation for lifecycle status
	lifecycleStatus, err := srs.lifecycleStatusRepo.GetByID(slaRegFormDTO.LifecycleStatusID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("status does not exist"), "not_exist", "lifecycle status")
	}
	_, _, _ = slaCategory, severityLevel, lifecycleStatus

	err = srs.repo.Update(id, currentUserID, slaRegFormDTO)

	if err != nil {
		return err
	}

	fieldsToCompare := []string{
		"TitleEn",
		"TitleAr",
		"SLACategoryID",
		"SeverityLevelID",
		"ResponsibleUnitID",
		"ResponsibleUserID",
		"Description",
		"Client"}
	changes, err := utility.CompareFields(existingSLARegister, slaRegFormDTO, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := srs.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			models.SLARegister{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return nil
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
