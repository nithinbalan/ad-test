package services

import (
	"fmt"
	"log"

	"jethur/grc/incident/dtos"
	"jethur/grc/incident/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
)

type IncidentSLAEscalationMatrixServiceInterface interface {
	GetAllEscalationMatrices(sort, sortBy, keyword, status string) ([]*dtos.IncidentEscalationMatrixListDTO, error)
	GetAllResolutionMatrices(sort, sortBy, keyword, status string) ([]*dtos.IncidentResolutionMatrixListDTO, error)
	UpdateEscalationMatrix(currentUserID uint, req *dtos.IncidentEscalationMatrixUpdateDTO) ([]*models.IncidentSLAEscalationMatrix, error)
	UpdateResolutionMatrix(currentUserID uint, req *dtos.IncidentResolutionMatrixUpdateDTO) ([]*models.IncidentSLAEscalationMatrix, error)
}

type IncidentSLAEscalationMatrixService struct {
	repo           repositories.IncidentSLAEscalationMatrixRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewIncidentSLAEscalationMatrixService(repo repositories.IncidentSLAEscalationMatrixRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) IncidentSLAEscalationMatrixServiceInterface {
	return &IncidentSLAEscalationMatrixService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (isems *IncidentSLAEscalationMatrixService) GetAllEscalationMatrices(sort, sortBy, keyword, status string) ([]*dtos.IncidentEscalationMatrixListDTO, error) {
	return isems.repo.GetAllEscalationMatrices(sort, sortBy, keyword, status)
}

func (isems *IncidentSLAEscalationMatrixService) GetAllResolutionMatrices(sort, sortBy, keyword, status string) ([]*dtos.IncidentResolutionMatrixListDTO, error) {
	return isems.repo.GetAllResolutionMatrices(sort, sortBy, keyword, status)
}

func (isems *IncidentSLAEscalationMatrixService) UpdateEscalationMatrix(currentUserID uint, req *dtos.IncidentEscalationMatrixUpdateDTO) ([]*models.IncidentSLAEscalationMatrix, error) {
	fieldsToCompare := []string{
		"NotifyWithin", "NotifyWithinUnitID", "NotifyToIncidentOwner", "NotifyToResponseTeam", "NotifyRoleID",
		"EscalateAfter", "EscalateAfterUnitID", "EscalateAfterIncidentOwner", "EscalateAfterResponseTeam", "EscalateRoleID",
	}

	var updatedMatrices []*models.IncidentSLAEscalationMatrix
	var updatedIDs []uint
	var allChanges []utility.ChangeDetail

	for _, form := range req.Settings {
		existingMatrix, err := isems.repo.GetBySeverityID(form.SeverityID, "escalation_matrix")
		if err != nil {
			return nil, err
		}

		changes, err := utility.CompareFields(existingMatrix, form, fieldsToCompare)
		if err != nil {
			log.Println("error :", err.Error())
		}
		allChanges = append(allChanges, changes...)

		existingMatrix.NotifyWithin = form.NotifyWithin
		existingMatrix.NotifyWithinUnitID = form.NotifyWithinUnitID
		existingMatrix.NotifyToIncidentOwner = form.NotifyToIncidentOwner
		existingMatrix.NotifyToResponseTeam = form.NotifyToResponseTeam
		existingMatrix.NotifyRoleID = form.NotifyRoleID
		existingMatrix.EscalateAfterIncidentOwner = form.EscalateAfterIncidentOwner
		existingMatrix.EscalateAfterResponseTeam = form.EscalateAfterResponseTeam
		existingMatrix.EscalateRoleID = form.EscalateRoleID

		if err := isems.repo.UpdateEscalationMatrix(existingMatrix); err != nil {
			return nil, err
		}

		updatedMatrices = append(updatedMatrices, existingMatrix)
		updatedIDs = append(updatedIDs, existingMatrix.ID)
	}

	if len(updatedIDs) > 0 {
		go func() {
			err := isems.activityLogSvc.LogBulkActivities(
				models.UpdateAction,
				models.IncidentSLAEscalationMatrix{}.TableName(),
				updatedIDs,
				allChanges,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk update activity: %v\n", err)
			}
		}()
	}

	return updatedMatrices, nil
}

func (isems *IncidentSLAEscalationMatrixService) UpdateResolutionMatrix(currentUserID uint, req *dtos.IncidentResolutionMatrixUpdateDTO) ([]*models.IncidentSLAEscalationMatrix, error) {
	fieldsToCompare := []string{
		"CloseWithin", "CloseWithinUnitID", "CloseWithinIncidentOwner", "CloseWithinResponseTeam", "CloseWithinRoleID",
	}

	var updatedMatrices []*models.IncidentSLAEscalationMatrix
	var updatedIDs []uint
	var allChanges []utility.ChangeDetail

	for _, form := range req.Settings {
		existingMatrix, err := isems.repo.GetBySeverityID(form.SeverityID, "resolution_matrix")
		if err != nil {
			return nil, err
		}

		changes, err := utility.CompareFields(existingMatrix, form, fieldsToCompare)
		if err != nil {
			log.Println("error :", err.Error())
		}
		allChanges = append(allChanges, changes...)

		existingMatrix.CloseWithin = form.CloseWithin
		existingMatrix.CloseWithinUnitID = form.CloseWithinUnitID
		existingMatrix.CloseWithinIncidentOwner = form.CloseWithinIncidentOwner
		existingMatrix.CloseWithinResponseTeam = form.CloseWithinResponseTeam
		existingMatrix.CloseWithinRoleID = form.CloseWithinRoleID

		if err := isems.repo.UpdateResolutionMatrix(existingMatrix); err != nil {
			return nil, err
		}

		updatedMatrices = append(updatedMatrices, existingMatrix)
		updatedIDs = append(updatedIDs, existingMatrix.ID)
	}

	if len(updatedIDs) > 0 {
		go func() {
			err := isems.activityLogSvc.LogBulkActivities(
				models.UpdateAction,
				models.IncidentSLAEscalationMatrix{}.TableName(),
				updatedIDs,
				allChanges,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk update activity: %v\n", err)
			}
		}()
	}

	return updatedMatrices, nil
}
