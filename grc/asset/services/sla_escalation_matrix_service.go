package services

import (
	"errors"
	"fmt"
	"log"

	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type SLAEscalationMatrixServiceInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*dtos.SLAEscalationMatrixListDTO, error)
	GetByID(id uint) (*models.SLAEscalationMatrix, error)
	Create(currentUserID uint, matrix *models.SLAEscalationMatrix) (*models.SLAEscalationMatrix, error)
	Update(id, currentUserID uint, matrix *models.SLAEscalationMatrix) (*models.SLAEscalationMatrix, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type SLAEscalationMatrixService struct {
	repo           repositories.SLAEscalationMatrixRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewSLAEscalationMatrixService(repo repositories.SLAEscalationMatrixRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) SLAEscalationMatrixServiceInterface {
	return &SLAEscalationMatrixService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (sems *SLAEscalationMatrixService)GetAll(sort, sortBy, keyword, status string) ([]*dtos.SLAEscalationMatrixListDTO, error) {
	return sems.repo.GetAll(sort, sortBy, keyword, status)
}

func (sems *SLAEscalationMatrixService) GetByID(id uint) (*models.SLAEscalationMatrix, error) {
	matrix, err := sems.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return matrix, nil
}

func (sems *SLAEscalationMatrixService) Create(currentUserID uint, matrix *models.SLAEscalationMatrix) (*models.SLAEscalationMatrix, error) {
	if exists, err := sems.repo.ExistsBySeverityID(matrix.SeverityID); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("this severity level already has an SLA entry"), "unique", "severity_id")
	}

	createdMatrix, err := sems.repo.Create(matrix)
	if err != nil {
		return nil, err
	}

	go func() {
		err := sems.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SLAEscalationMatrix{}.TableName(),
			createdMatrix.ID,
			createdMatrix,
			currentUserID,
			"governance",
			"sla_escalation_matrix",
			"sla_escalation_matrix",
			nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdMatrix, nil
}

func (sems *SLAEscalationMatrixService) Update(id, currentUserID uint, updatedMatrix *models.SLAEscalationMatrix) (*models.SLAEscalationMatrix, error) {
	existingMatrix, err := sems.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedMatrix.SeverityID != existingMatrix.SeverityID {
		if exists, err := sems.repo.ExistsBySeverityIDExcludingID(id, updatedMatrix.SeverityID); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("this severity level already has an SLA entry"), "unique", "severity_id")
		}
	}

	fieldsToCompare := []string{"SeverityID", "RemediationDays", "IsEmail", "IsAlert", "RoleID", "Status"}
	changes, err := utility.CompareFields(existingMatrix, updatedMatrix, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	existingMatrix.SeverityID = updatedMatrix.SeverityID
	existingMatrix.RemediationDays = updatedMatrix.RemediationDays
	existingMatrix.IsEmail = updatedMatrix.IsEmail
	existingMatrix.IsAlert = updatedMatrix.IsAlert
	existingMatrix.RoleID = updatedMatrix.RoleID
	existingMatrix.Status = updatedMatrix.Status

	if err := sems.repo.Update(existingMatrix); err != nil {
		return nil, err
	}

	go func() {
		err := sems.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.SLAEscalationMatrix{}.TableName(),
			existingMatrix.ID,
			changedDataJSON,
			currentUserID,
			"governance",
			"sla_escalation_matrix",
			"sla_escalation_matrix",
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log update activity: %v\n", err)
		}
	}()

	return existingMatrix, nil
}

func (sems *SLAEscalationMatrixService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedMatrices []map[string]interface{}
	for _, id := range ids {
		matrix, err := sems.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = sems.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedMatrices = append(deletedMatrices, map[string]interface{}{
			"id":          id,
			"severity_id": matrix.SeverityID,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := sems.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.SLAEscalationMatrix{}.TableName(),
				successfulIds,
				deletedMatrices,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (sems *SLAEscalationMatrixService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := sems.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := sems.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.SLAEscalationMatrix{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk activate activity: %v\n", err)
			}
		}()
	}

	return err
}

func (sems *SLAEscalationMatrixService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := sems.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := sems.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.SLAEscalationMatrix{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk deactivate activity: %v\n", err)
			}
		}()
	}

	return err
}