package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"log"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type JobDescriptionServiceInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*models.JobDescription, error)
	GetByID(id uint) (*models.JobDescription, error)
	GetFullDetailsByID(id uint) (*dtos.JobDescriptionDetailDTO, error)
	Create(currentUserID uint, jobDesc *models.JobDescription) (*models.JobDescription, error)
	Update(id uint, currentUserID uint, jobDesc *models.JobDescription) (*models.JobDescription, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type JobDescriptionService struct {
	repo           repositories.JobDescriptionRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewJobDescriptionService(repo repositories.JobDescriptionRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) JobDescriptionServiceInterface {
	return &JobDescriptionService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (s *JobDescriptionService) GetAll(sort, sortBy, keyword, status string) ([]*models.JobDescription, error) {
	return s.repo.GetAll(sort, sortBy, keyword, status)
}

func (s *JobDescriptionService) GetByID(id uint) (*models.JobDescription, error) {
	jobDesc, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return jobDesc, nil
}

func (s *JobDescriptionService) GetFullDetailsByID(id uint) (*dtos.JobDescriptionDetailDTO, error) {
	var dto dtos.JobDescriptionDetailDTO

	jd, err := s.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	err = mapstructure.Decode(jd, &dto)
	if err != nil {
		return nil, err
	}
	return &dto, nil
}

func (s *JobDescriptionService) Create(currentUserID uint, jobDesc *models.JobDescription) (*models.JobDescription, error) {
	if exists, err := s.repo.ExistsByTitle(jobDesc.TitleEn, jobDesc.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	created, err := s.repo.Create(jobDesc)
	if err != nil {
		return nil, err
	}

	go func() {
		err := s.activityLogSvc.LogActivity(models.CreateAction, jobDesc.TableName(), created.ID, created, currentUserID, nil)
		if err != nil {
			log.Printf("Failed to log create activity: %v", err)
		}
	}()

	return created, nil
}

func (s *JobDescriptionService) Update(id uint, currentUserID uint, updated *models.JobDescription) (*models.JobDescription, error) {
	existing, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updated.TitleEn != existing.TitleEn || updated.TitleAr != existing.TitleAr {
		if exists, err := s.repo.ExistsByTitleExcludingID(id, updated.TitleEn, updated.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := s.repo.Update(id, updated)
	if err != nil {
		return nil, err
	}

	fields := []string{"TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "RolesAndResponsibilities", "Status"}
	changes, err := utility.CompareFields(existing, updated, fields)
	if err != nil {
		log.Println("compare error:", err)
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("json conversion error:", err)
	}

	go func() {
		err := s.activityLogSvc.LogActivity(models.UpdateAction, updated.TableName(), result.ID, changedDataJSON, currentUserID, nil)
		if err != nil {
			log.Printf("Failed to log update activity: %v", err)
		}
	}()

	return result, nil
}

func (s *JobDescriptionService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failed, successful []uint

	for _, id := range ids {
		if _, err := s.repo.GetByID(id); err != nil {
			continue
		}
		if err := s.repo.Delete(id); err != nil {
			failed = append(failed, id)
		} else {
			successful = append(successful, id)
		}
	}

	if len(successful) > 0 {
		go func() {
			err := s.activityLogSvc.LogBulkActivities(models.DeleteAction, models.JobDescription{}.TableName(), successful, nil, currentUserID, nil)
			if err != nil {
				log.Printf("Failed to log delete activity: %v", err)
			}
		}()
	}

	return failed, nil
}

func (s *JobDescriptionService) Activate(ids []uint, currentUserID uint) error {
	successful, err := s.repo.Activate(ids)
	if len(successful) > 0 {
		go func() {
			_ = s.activityLogSvc.LogBulkActivities(models.ActivateAction, models.JobDescription{}.TableName(), successful, nil, currentUserID, nil)
		}()
	}
	return err
}

func (s *JobDescriptionService) Deactivate(ids []uint, currentUserID uint) error {
	successful, err := s.repo.Deactivate(ids)
	if len(successful) > 0 {
		go func() {
			_ = s.activityLogSvc.LogBulkActivities(models.DeactivateAction, models.JobDescription{}.TableName(), successful, nil, currentUserID, nil)
		}()
	}
	return err
}
