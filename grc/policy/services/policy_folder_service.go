package services

import (
	"fmt"
	documentServices "jethur/grc/documents/services"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	"jethur/grc/policy/dtos"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/grc/policy/repositories"
	settingSvc "jethur/grc/settings/services"
	"jethur/utility"
	"log"
)

type PolicyFolderServiceInterface interface {
	GetAll(params *policyDTOs.PolicyRegisterListQueryParamsDTO, currentUserID uint) (*policyDTOs.PolicyRegisterListResponseDTO, error)
	GetByID(id uint) (*models.Document, error)
	Create(currentUserID uint, formDTO *dtos.PolicyFolderFormDTO) (*models.Document, error)
	Update(id uint, currentUserID uint, formDTO *dtos.PolicyFolderFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type PolicyFolderService struct {
	repo                  repositories.PolicyFolderRepositoryInterface
	activityLogSvc        services.ActivityLogServiceInterface
	documentAccessTypeSvc documentServices.DocumentAccessTypeServiceInterface
	orgLevelSvc           settingSvc.OrganizationLevelServiceInterface
	userSvc               organizationSvc.UserServiceInterface
}

func NewPolicyFolderService(repo repositories.PolicyFolderRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface, documentAccessTypeSvc documentServices.DocumentAccessTypeServiceInterface, orgLevelSvc settingSvc.OrganizationLevelServiceInterface, userSvc organizationSvc.UserServiceInterface) PolicyFolderServiceInterface {
	return &PolicyFolderService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		documentAccessTypeSvc: documentAccessTypeSvc,
		orgLevelSvc:           orgLevelSvc,
		userSvc:               userSvc,
	}
}

func (cfs *PolicyFolderService) GetAll(params *policyDTOs.PolicyRegisterListQueryParamsDTO, currentUserID uint) (*policyDTOs.PolicyRegisterListResponseDTO, error) {
	// Get user details to check for super admin role
	user, err := cfs.userSvc.GetFullDetailsByID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Check if user has super admin role
	isUserSuperAdmin := false
	if user.Roles != nil {
		for _, role := range *user.Roles {
			if role.IsSuperAdmin {
				isUserSuperAdmin = true
				break
			}
		}
	}

	// Get primary organization unit field
	primaryLevel, err := cfs.orgLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}

	// Determine the field name based on the primary level
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
	return cfs.repo.GetAll(params, currentUserID, isUserSuperAdmin, primaryOrgUnitField)
}

func (cfs PolicyFolderService) GetByID(id uint) (*models.Document, error) {
	return cfs.repo.GetByID(id)
}

func (cfs *PolicyFolderService) Create(currentUserID uint, formDTO *dtos.PolicyFolderFormDTO) (*models.Document, error) {
	if exists, err := cfs.repo.ExistsByTitle(formDTO.TitleEn, formDTO.TitleAr, formDTO.ParentDocumentID); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	accessType, err := cfs.documentAccessTypeSvc.GetByType("public")
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("document access type does not exist"), "not_exist", "document access type")
	}

	folder := &models.Document{
		TitleEn:              formDTO.TitleEn,
		TitleAr:              formDTO.TitleEn,
		IsFolder:             true,
		ParentDocumentID:     formDTO.ParentDocumentID,
		DocumentAccessTypeID: uint(accessType.ID),
		CreatedBYID:          currentUserID,
	}

	createdFolder, err := cfs.repo.Create(folder)
	if err != nil {
		return nil, err
	}

	go func() {
		err := cfs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Document{}.TableName(),
			folder.ID,
			createdFolder,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdFolder, nil
}

func (cfs *PolicyFolderService) Update(id uint, currentUserID uint, formDTO *dtos.PolicyFolderFormDTO) error {
	existingFolder, err := cfs.repo.GetByID(id)

	if err != nil {
		return err
	}

	if formDTO.TitleEn != existingFolder.TitleEn || formDTO.TitleAr != existingFolder.TitleAr {
		if exists, err := cfs.repo.ExistsByTitleExcludingID(id, formDTO.TitleEn, formDTO.TitleAr, formDTO.ParentDocumentID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	existingFolder.UpdatedBYID = &currentUserID

	result, err := cfs.repo.Update(id, existingFolder)

	if err != nil {
		return err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingFolder, existingFolder, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := cfs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.Document{}.TableName(),
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
	return err
}

func (cfs *PolicyFolderService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedFolders []map[string]interface{}

	for _, id := range ids {
		folder, err := cfs.repo.GetByID(id)
		if err != nil {
			continue
		}

		err = cfs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedFolders = append(deletedFolders, map[string]interface{}{
			"id":       id,
			"title_en": folder.TitleEn,
			"title_ar": folder.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := cfs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.Document{}.TableName(),
				successfulIds,
				deletedFolders,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}
