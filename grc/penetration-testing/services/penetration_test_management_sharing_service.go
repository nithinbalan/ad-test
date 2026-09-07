package services

import (
	"errors"
	"fmt"
	"time"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	penetrationTestRepositories "jethur/grc/penetration-testing/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PenetrationTestManagementSharingServiceInterface interface {
	GetByID(penetrationTestID uint) (*models.PenetrationTestManagementSharing, error)
	GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*penetrationTestDTOs.PenetrationTestManagementSharingDetailDTO, error)
	Create(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestManagementSharingFormDTO) (uint, error)
	Update(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestManagementSharingFormDTO) error
	Delete(penetrationTestID uint, currentUserID uint) (uint, error)
}

// Service and implementations
type PenetrationTestManagementSharingService struct {
	repo               penetrationTestRepositories.PenetrationTestManagementSharingRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	penetrationTestSvc PenetrationTestServiceInterface
}

func NewPenetrationTestManagementSharingService(
	repo penetrationTestRepositories.PenetrationTestManagementSharingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	penetrationTestSvc PenetrationTestServiceInterface,
) PenetrationTestManagementSharingServiceInterface {
	return &PenetrationTestManagementSharingService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		penetrationTestSvc: penetrationTestSvc,
	}
}

func (ptms PenetrationTestManagementSharingService) GetByID(penetrationTestID uint) (*models.PenetrationTestManagementSharing, error) {
	if _, err := ptms.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		return nil, err
	}

	managementSharing, err := ptms.repo.GetByID(penetrationTestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, penetrationTestID)
		}
		return nil, err
	}
	return managementSharing, nil
}

func (ptms PenetrationTestManagementSharingService) GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*penetrationTestDTOs.PenetrationTestManagementSharingDetailDTO, error) {
	var managementSharingDTO penetrationTestDTOs.PenetrationTestManagementSharingDetailDTO

	if _, err := ptms.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		return nil, err
	}

	managementSharing, err := ptms.repo.GetFullDetailsByID(penetrationTestID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, penetrationTestID)
		}
		return nil, err
	}

	err = mapstructure.Decode(managementSharing, &managementSharingDTO)
	if err != nil {
		return nil, err
	}

	if managementSharing.ManagementNotifyDate != nil {
		managementNotifyDate := managementSharing.ManagementNotifyDate.Format("2006-01-02")
		managementSharingDTO.ManagementNotifyDateString = &managementNotifyDate
	}

	if managementSharing.ManagementResponseDate != nil {
		managementResponseDate := managementSharing.ManagementResponseDate.Format("2006-01-02")
		managementSharingDTO.ManagementResponseDateString = &managementResponseDate
	}

	managementSharingDTO.CreatedAt = managementSharing.CreatedAt
	managementSharingDTO.UpdatedAt = managementSharing.UpdatedAt

	return &managementSharingDTO, nil
}

func (ptms *PenetrationTestManagementSharingService) Create(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestManagementSharingFormDTO) (uint, error) {
	tx := ptms.repo.BeginTransaction()

	if _, err := ptms.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		tx.Rollback()
		return 0, err
	}

	var managementNotifyDate *time.Time
	if form.ManagementNotifyDate != nil {
		parsedDate, err := utility.ParseDate(*form.ManagementNotifyDate)
		if err != nil {
			tx.Rollback()
			return 0, utility.NewValidationError(fmt.Errorf("invalid management notify date"), "invalid_date", "management_notify_date")
		}
		managementNotifyDate = parsedDate
	}

	var managementResponseDate *time.Time
	if form.ManagementResponseDate != nil {
		parsedDate, err := utility.ParseDate(*form.ManagementResponseDate)
		if err != nil {
			tx.Rollback()
			return 0, utility.NewValidationError(fmt.Errorf("invalid management response date"), "invalid_date", "management_response_date")
		}
		managementResponseDate = parsedDate
	}

	managementSharing := &models.PenetrationTestManagementSharing{
		PenetrationTestID:           penetrationTestID,
		ResultsSharedWithManagement: form.ResultsSharedWithManagement,
		ManagementNotifyDate:        managementNotifyDate,
		ManagementResponseDate:      managementResponseDate,
		RemediationNotes:            form.RemediationNotes,
		CreatedBYID:                 currentUserID,
	}

	id, err := ptms.repo.Create(managementSharing, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingManagementSharing, err := ptms.GetFullDetailsByID(penetrationTestID)
		if err != nil {
			utility.LogError("Failed to fetch penetration test management sharing", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = ptms.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PenetrationTestManagementSharing{}.TableName(),
			id,
			existingManagementSharing,
			currentUserID,
			"penetration_testing", "penetration_testing_pentest_log", "management sharing", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.PenetrationTestManagementSharing{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ptms *PenetrationTestManagementSharingService) Update(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestManagementSharingFormDTO) error {
	tx := ptms.repo.BeginTransaction()

	if _, err := ptms.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		tx.Rollback()
		return err
	}

	existingManagementSharing, err := ptms.GetByID(penetrationTestID)
	if err != nil {
		tx.Rollback()
		return err
	}

	var managementNotifyDate *time.Time
	if form.ManagementNotifyDate != nil {
		parsedDate, err := utility.ParseDate(*form.ManagementNotifyDate)
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("invalid management notify date"), "invalid_date", "management_notify_date")
		}
		managementNotifyDate = parsedDate
	}

	var managementResponseDate *time.Time
	if form.ManagementResponseDate != nil {
		parsedDate, err := utility.ParseDate(*form.ManagementResponseDate)
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("invalid management response date"), "invalid_date", "management_response_date")
		}
		managementResponseDate = parsedDate
	}

	existingManagementSharing.ResultsSharedWithManagement = form.ResultsSharedWithManagement
	existingManagementSharing.ManagementNotifyDate = managementNotifyDate
	existingManagementSharing.ManagementResponseDate = managementResponseDate
	existingManagementSharing.RemediationNotes = form.RemediationNotes
	existingManagementSharing.UpdatedBYID = &currentUserID

	err = ptms.repo.Update(existingManagementSharing, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingEntry, err := ptms.GetFullDetailsByID(penetrationTestID)
		if err != nil {
			utility.LogError("Failed to fetch penetration test management sharing", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = ptms.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PenetrationTestManagementSharing{}.TableName(),
			existingManagementSharing.ID,
			existingEntry,
			currentUserID,
			"penetration_testing", "penetration_testing_pentest_log", "management sharing", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.PenetrationTestManagementSharing{}.TableName(),
				"recordID":  existingManagementSharing.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (ptms *PenetrationTestManagementSharingService) Delete(penetrationTestID uint, currentUserID uint) (uint, error) {
	tx := ptms.repo.BeginTransaction()
	var failedId uint
	var deletedManagementSharing []map[string]interface{}

	if _, err := ptms.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		tx.Rollback()
		return 0, err
	}

	managementSharing, err := ptms.GetByID(penetrationTestID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err = ptms.repo.Delete(managementSharing.ID, tx); err != nil {
		tx.Rollback()
		failedId = managementSharing.ID
		return failedId, err
	}

	deletedManagementSharing = append(deletedManagementSharing, map[string]interface{}{
		"id":                  managementSharing.ID,
		"penetration_test_id": penetrationTestID,
	})

	tx.Commit()

	go func() {
		err := ptms.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.PenetrationTestManagementSharing{}.TableName(),
			managementSharing.ID,
			deletedManagementSharing,
			currentUserID,
			"penetration_testing", "penetration_testing_pentest_log", "management sharing", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return failedId, nil
}
