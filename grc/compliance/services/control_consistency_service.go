package services

import (
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
)

type ControlConsistencyServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlConsistency, error)
	GetByID(id uint) (*models.ControlConsistency, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type ControlConsistencyService struct {
	repo           repositories.ControlConsistencyRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewControlConsistencyService(repo repositories.ControlConsistencyRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) ControlConsistencyServiceInterface {
	return &ControlConsistencyService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ccs *ControlConsistencyService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlConsistency, error) {
	return ccs.repo.GetAll(sort, sortBy, keyword, status)
}

func (ccs *ControlConsistencyService) GetByID(id uint) (*models.ControlConsistency, error) {
	controlConsistency, err := ccs.repo.GetByID(id)
	if err != nil {
		if _, ok := err.(*utility.ResourceNotFoundError); ok {
			return nil, err
		}
		return nil, err
	}
	return controlConsistency, nil
}

func (ccs *ControlConsistencyService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ccs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ControlConsistency{}.TableName(),
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

func (ccs *ControlConsistencyService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ccs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ControlConsistency{}.TableName(),
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


