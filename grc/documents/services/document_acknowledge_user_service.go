package services

import (
	"fmt"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"
	settingServices "jethur/grc/settings/services"

	"jethur/grc/models"
	"jethur/utility"
)

type DocumentAcknowledgeUserServiceInterface interface {
	GetAll(documentID uint, params *documentDTOs.DocumentAcknowledgeListQueryParamsDTO) (*documentDTOs.DocumentAcknowledgementListResponseDTO, error)
	Create(docuemntID uint, currentUserID uint) (uint, error)
	Delete(docuemntID uint, currentUserID uint) (uint, error)
}

type DocumentAcknowledgeUserService struct {
	repo                 repositories.DocumentAcknowledgeUserRepositoryInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	documentSvc          DocumentServiceInterface
	userSvc              organizationServices.UserServiceInterface
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface
}

func NewDocumentAcknowledgeUserService(repo repositories.DocumentAcknowledgeUserRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	documentSvc DocumentServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface) DocumentAcknowledgeUserServiceInterface {
	return &DocumentAcknowledgeUserService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		documentSvc:          documentSvc,
		userSvc:              userSvc,
		organizationLevelSvc: organizationLevelSvc,
	}
}

func (daus *DocumentAcknowledgeUserService) GetAll(documentID uint, params *documentDTOs.DocumentAcknowledgeListQueryParamsDTO) (*documentDTOs.DocumentAcknowledgementListResponseDTO, error) {
	primaryLevel, err := daus.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}

	var primaryOrgUnitField string

	switch primaryLevel.Level {
	case 1:
		primaryOrgUnitField = "organization_unit_1_id"
	case 2:
		primaryOrgUnitField = "organization_unit_2_id"
	case 3:
		primaryOrgUnitField = "organization_unit_3_id"
	case 4:
		primaryOrgUnitField = "organization_unit_4_id"
	case 5:
		primaryOrgUnitField = "organization_unit_5_id"
	default:
		return nil, fmt.Errorf("unsupported primary level")
	}

	return daus.repo.GetAll(documentID, primaryOrgUnitField, params)
}

func (daus *DocumentAcknowledgeUserService) Create(docuemntID uint, currentUserID uint) (uint, error) {
	tx := daus.repo.BeginTransaction()

	document, err := daus.documentSvc.GetFullDetailsByID(docuemntID, nil, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Check if the entry already exists
	existing, err := daus.repo.Exists(docuemntID, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if existing {
		tx.Rollback()
		return 0, utility.NewForbiddenError("document acknowledgment by user already exists")
	}

	organizationUnitID, err := daus.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for owner"), "not_exist", "process owner unit")
	}

	documentAcknowledgeUser := &models.DocumentAcknowledgeUser{
		DocumentID:         docuemntID,
		UserID:             currentUserID,
		OrganizationUnitID: organizationUnitID,
		CreatedBYID:        currentUserID,
	}

	id, err := daus.repo.Create(documentAcknowledgeUser, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	//Add Audit logs
	if err := daus.documentSvc.CreateDocumentAuditLog(uint(docuemntID), document.LatestDocumentVersion.ID, currentUserID, "acknowledged", tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := daus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.DocumentAcknowledgeUser{}.TableName(),
			documentAcknowledgeUser.ID,
			nil,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.DocumentAcknowledgeUser{}.TableName(),
				"recordID":  id.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id.ID, nil
}

func (daus *DocumentAcknowledgeUserService) Delete(docuemntID uint, currentUserID uint) (uint, error) {
	tx := daus.repo.BeginTransaction()
	var failedId uint

	_, err := daus.documentSvc.GetByID(docuemntID)
	if err != nil {
		return 0, err
	}
	// Check if the entry already exists
	exists, err := daus.repo.Exists(docuemntID, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if !exists {
		tx.Rollback()
		return 0, utility.NewForbiddenError("document acknowledgment by user does not exist")
	}

	if err = daus.repo.Delete(uint(docuemntID), currentUserID, tx); err != nil {
		failedId = uint(docuemntID)
		return failedId, err
	}
	tx.Commit()
	// Log successful delete activity
	go func() {
		err := daus.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.DocumentAcknowledgeUser{}.TableName(),
			uint(docuemntID),
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			// Handle error if logging fails
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return failedId, nil
}
