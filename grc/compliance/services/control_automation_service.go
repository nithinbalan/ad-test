package services

import (
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
)

type ControlAutomationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlAutomation, error)
	GetByID(id uint) (*models.ControlAutomation, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type ControlAutomationService struct {
	repo           repositories.ControlAutomationRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewControlAutomationService(repo repositories.ControlAutomationRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) ControlAutomationServiceInterface {
	return &ControlAutomationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (cas *ControlAutomationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlAutomation, error) {
	return cas.repo.GetAll(sort, sortBy, keyword, status)
}

func (cas *ControlAutomationService) GetByID(id uint) (*models.ControlAutomation, error) {
	controlAutomation, err := cas.repo.GetByID(id)
	if err != nil {
		if _, ok := err.(*utility.ResourceNotFoundError); ok {
			return nil, err
		}
		return nil, err
	}
	return controlAutomation, nil
}

func (cas *ControlAutomationService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := cas.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := cas.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ControlAutomation{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
			}
		}()
	}

	return err
}

func (cas *ControlAutomationService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := cas.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := cas.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ControlAutomation{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
			}
		}()
	}

	return err
}


