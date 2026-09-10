package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	"jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
)

type ProjectSecurityGateChecklistServiceInterface interface {
	SetSecurityGateService(svc ProjectSecurityGateServiceInterface)
	GetAll(securityGateID uint, params *dtos.ProjectSecurityGateChecklistListQueryParamsDTO) (*dtos.ProjectSecurityGateChecklistListResponseDTO, error)
	GetChecklistSummary(securityGateID uint) (*dtos.ProjectSecurityGateChecklistSummaryDTO, error)
	GetFullDetailsByID(securityGateID uint, checklistID uint) (*dtos.ProjectSecurityGateChecklistDetailDTO, error)
	GetByID(securityGateID uint, checklistID uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklist, error)
	CreateChecklists(checklists []*models.ProjectSecurityGateChecklist, tx ...utility.DatabaseTransaction) error
	DeleteChecklistsBySecurityGateID(securityGateID uint, tx ...utility.DatabaseTransaction) error
	CreateChecklistResponse(securityGateID uint, checklistID uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistResponseFormDTO) (uint, error)
	UpdateChecklistResponse(securityGateID uint, checklistID uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistResponseFormDTO) error
	GetResponseDocumentPath(securityGateID uint, checklistID uint, token string) (string, error)
	UpdateChecklistControls(securityGateID uint, checklistID uint, currentUserID uint, form *commonDTOs.AssociatedControlMapDTO) error
}

type ProjectSecurityGateChecklistService struct {
	repo              repositories.ProjectSecurityGateChecklistRepositoryInterface
	securityGateSvc   ProjectSecurityGateServiceInterface
	responseStatusSvc ProjectSecurityGateChecklistResponseStatusServiceInterface
	controlSvc        organizationServices.ControlServiceInterface
	activityLog       logServices.ActivityLogServiceInterface
}

func NewProjectSecurityGateChecklistService(
	repo repositories.ProjectSecurityGateChecklistRepositoryInterface,
	responseStatusSvc ProjectSecurityGateChecklistResponseStatusServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
	activityLog logServices.ActivityLogServiceInterface,
) ProjectSecurityGateChecklistServiceInterface {
	return &ProjectSecurityGateChecklistService{
		repo:              repo,
		responseStatusSvc: responseStatusSvc,
		controlSvc:        controlSvc,
		activityLog:       activityLog,
	}
}

func (psgcs *ProjectSecurityGateChecklistService) SetSecurityGateService(svc ProjectSecurityGateServiceInterface) {
	psgcs.securityGateSvc = svc
}

func (psgcs *ProjectSecurityGateChecklistService) GetAll(securityGateID uint, params *dtos.ProjectSecurityGateChecklistListQueryParamsDTO) (*dtos.ProjectSecurityGateChecklistListResponseDTO, error) {
	if _, err := psgcs.securityGateSvc.GetByID(securityGateID); err != nil {
		return nil, err
	}

	return psgcs.repo.GetAll(securityGateID, params)
}

func (psgcs *ProjectSecurityGateChecklistService) GetChecklistSummary(securityGateID uint) (*dtos.ProjectSecurityGateChecklistSummaryDTO, error) {
	if _, err := psgcs.securityGateSvc.GetByID(securityGateID); err != nil {
		return nil, err
	}

	return psgcs.repo.GetSummaryBySecurityGateID(securityGateID)
}

func (psgcs *ProjectSecurityGateChecklistService) GetFullDetailsByID(securityGateID uint, checklistID uint) (*dtos.ProjectSecurityGateChecklistDetailDTO, error) {
	checklist, err := psgcs.repo.GetFullDetailsByID(checklistID)
	if err != nil {
		return nil, err
	}

	if checklist.SecurityGateID != securityGateID {
		return nil, &utility.ResourceNotFoundError{}
	}

	var detail dtos.ProjectSecurityGateChecklistDetailDTO
	if err := mapstructure.Decode(checklist, &detail); err != nil {
		return nil, err
	}

	detail.CreatedAt = checklist.CreatedAt
	detail.UpdatedAt = checklist.UpdatedAt

	for i := range detail.Comments {
		if i < len(checklist.Comments) {
			detail.Comments[i].CreatedAt = checklist.Comments[i].CreatedAt
		}
	}

	detail.Controls = make([]*organizationDTOs.ControlMinimalListDTO, 0, len(checklist.Controls))
	for _, mapping := range checklist.Controls {
		if mapping.Control == nil {
			continue
		}

		var controlDTO organizationDTOs.ControlMinimalListDTO
		if err := mapstructure.Decode(mapping.Control, &controlDTO); err != nil {
			return nil, err
		}
		detail.Controls = append(detail.Controls, &controlDTO)
	}

	return &detail, nil
}

func (psgcs *ProjectSecurityGateChecklistService) CreateChecklists(checklists []*models.ProjectSecurityGateChecklist, tx ...utility.DatabaseTransaction) error {
	return psgcs.repo.CreateChecklists(checklists, tx...)
}

func (psgcs *ProjectSecurityGateChecklistService) DeleteChecklistsBySecurityGateID(securityGateID uint, tx ...utility.DatabaseTransaction) error {
	return psgcs.repo.DeleteBySecurityGateID(securityGateID, tx...)
}

func (psgcs *ProjectSecurityGateChecklistService) GetByID(securityGateID uint, checklistID uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklist, error) {
	checklist, err := psgcs.repo.GetByID(checklistID, tx...)
	if err != nil {
		return nil, err
	}

	if checklist.SecurityGateID != securityGateID {
		return nil, &utility.ResourceNotFoundError{}
	}

	return checklist, nil
}

func (psgcs *ProjectSecurityGateChecklistService) GetResponseDocumentPath(securityGateID uint, checklistID uint, token string) (string, error) {
	checklist, err := psgcs.GetByID(securityGateID, checklistID)
	if err != nil {
		return "", err
	}

	if len(checklist.Documents) == 0 {
		return "", &utility.ResourceNotFoundError{}
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", &utility.ResourceNotFoundError{}
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", &utility.ResourceNotFoundError{}
	}

	storagePath := filepath.Join(
		"project",
		"security-gates", strconv.FormatUint(uint64(securityGateID), 10),
		"checklists", strconv.FormatUint(uint64(checklistID), 10),
		"documents",
	)

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", &utility.ResourceNotFoundError{}
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (psgcs *ProjectSecurityGateChecklistService) CreateChecklistResponse(securityGateID uint, checklistID uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistResponseFormDTO) (id uint, err error) {
	tx := psgcs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	checklist, err := psgcs.GetByID(securityGateID, checklistID, tx)
	if err != nil {
		return 0, err
	}

	if checklist.CreatedBYID != nil {
		err = utility.NewValidationError(
			fmt.Errorf("a response already exists for this checklist item"),
			"unique",
			"response",
		)
		return 0, err
	}

	var responseStatusID *uint
	var exceptionJustification *string
	if form.IsException {
		exceptionJustification = form.ExceptionJustification
	} else {
		if form.ResponseStatusID != nil {
			if _, err = psgcs.responseStatusSvc.GetByID(*form.ResponseStatusID, tx); err != nil {
				return 0, err
			}
		}
		responseStatusID = form.ResponseStatusID
	}

	now := time.Now()
	checklist.ResponseStatusID = responseStatusID
	checklist.IsException = form.IsException
	checklist.ExceptionJustification = exceptionJustification
	checklist.Comment = form.Comment
	checklist.CreatedBYID = &currentUserID
	checklist.CreatedAt = &now

	if err = psgcs.repo.UpdateResponse(checklist, tx); err != nil {
		return 0, err
	}

	id = checklist.ID

	var documents []*models.ProjectSecurityGateChecklistResponseDocument

	for _, doc := range form.Documents {
		encryptedFileName, encErr := utility.EncryptData([]byte(doc.FileName))
		if encErr != nil {
			err = encErr
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"project",
			"security-gates", strconv.FormatUint(uint64(securityGateID), 10),
			"checklists", strconv.FormatUint(uint64(checklistID), 10),
			"documents",
		)

		documents = append(documents, &models.ProjectSecurityGateChecklistResponseDocument{
			ChecklistID:  checklistID,
			URL:          url,
			Name:         doc.FileName,
			Size:         doc.Size,
			Token:        encryptedFileName,
			OriginalName: doc.OriginalName,
			Type:         doc.Type,
			CreatedBYID:  currentUserID,
		})
	}

	if len(documents) > 0 {
		if err = psgcs.repo.CreateResponseDocuments(documents, tx); err != nil {
			return 0, err
		}

		storagePath := filepath.Join(
			"project",
			"security-gates", strconv.FormatUint(uint64(securityGateID), 10),
			"checklists", strconv.FormatUint(uint64(checklistID), 10),
			"documents",
		)

		for _, document := range documents {
			if err = utility.StoreFile(document.Name, storagePath); err != nil {
				return 0, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		activityErr := psgcs.activityLog.LogActivity(
			models.CreateAction,
			models.ProjectSecurityGateChecklist{}.TableName(),
			id,
			checklist,
			currentUserID,
			"project", "project_security_gate", "project_security_gate_checklist_response_create", nil,
		)
		if activityErr != nil {
			log.Printf("Failed to log activity: %v", activityErr)
		}
	}()

	return id, err
}

func (psgcs *ProjectSecurityGateChecklistService) UpdateChecklistResponse(securityGateID uint, checklistID uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistResponseFormDTO) (err error) {
	tx := psgcs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	checklist, err := psgcs.GetByID(securityGateID, checklistID, tx)
	if err != nil {
		return err
	}

	if checklist.CreatedBYID == nil {
		err = &utility.ResourceNotFoundError{}
		return err
	}

	var responseStatusID *uint
	var exceptionJustification *string
	if form.IsException {
		exceptionJustification = form.ExceptionJustification
	} else {
		if form.ResponseStatusID != nil {
			if _, err = psgcs.responseStatusSvc.GetByID(*form.ResponseStatusID, tx); err != nil {
				return err
			}
		}
		responseStatusID = form.ResponseStatusID
	}

	checklist.ResponseStatusID = responseStatusID
	checklist.IsException = form.IsException
	checklist.ExceptionJustification = exceptionJustification
	checklist.Comment = form.Comment
	checklist.UpdatedBYID = &currentUserID
	now := time.Now()
	checklist.UpdatedAt = &now

	if err = psgcs.repo.UpdateResponse(checklist, tx); err != nil {
		return err
	}

	var newDocuments []*models.ProjectSecurityGateChecklistResponseDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.ProjectSecurityGateChecklistResponseDocument

	for _, doc := range form.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			encryptedFileName, encErr := utility.EncryptData([]byte(doc.FileName))
			if encErr != nil {
				err = encErr
				return err
			}

			url := utility.GenerateDocumentURL(
				"project",
				"security-gates", strconv.FormatUint(uint64(securityGateID), 10),
				"checklists", strconv.FormatUint(uint64(checklistID), 10),
				"documents",
			)

			newDocuments = append(newDocuments, &models.ProjectSecurityGateChecklistResponseDocument{
				ChecklistID:  checklist.ID,
				URL:          url,
				Name:         doc.FileName,
				Size:         doc.Size,
				Token:        encryptedFileName,
				OriginalName: doc.OriginalName,
				Type:         doc.Type,
				CreatedBYID:  currentUserID,
			})
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, docErr := psgcs.repo.GetResponseDocumentByID(*doc.ID)
				if docErr != nil {
					err = docErr
					return err
				}
				if existingDocument.ChecklistID != checklist.ID {
					err = &utility.ResourceNotFoundError{}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			}
		}
	}

	if len(newDocuments) > 0 {
		if err = psgcs.repo.CreateResponseDocuments(newDocuments, tx); err != nil {
			return err
		}
	}

	if len(deleteDocumentIDs) > 0 {
		if err = psgcs.repo.DeleteResponseDocuments(deleteDocumentIDs, tx); err != nil {
			return err
		}
	}

	storagePath := filepath.Join(
		"project",
		"security-gates", strconv.FormatUint(uint64(securityGateID), 10),
		"checklists", strconv.FormatUint(uint64(checklistID), 10),
		"documents",
	)

	for _, document := range newDocuments {
		if err = utility.StoreFile(document.Name, storagePath); err != nil {
			return err
		}
	}

	for _, document := range documentsToDelete {
		if err = utility.RemoveFile(document.Name, storagePath); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	go func() {
		activityErr := psgcs.activityLog.LogActivity(
			models.UpdateAction,
			models.ProjectSecurityGateChecklist{}.TableName(),
			checklist.ID,
			checklist,
			currentUserID,
			"project", "project_security_gate", "project_security_gate_checklist_response_update", nil,
		)
		if activityErr != nil {
			log.Printf("Failed to log update activity: %v", activityErr)
		}
	}()

	return err
}

func (psgcs *ProjectSecurityGateChecklistService) UpdateChecklistControls(securityGateID uint, checklistID uint, currentUserID uint, form *commonDTOs.AssociatedControlMapDTO) (err error) {
	tx := psgcs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	if _, err = psgcs.GetByID(securityGateID, checklistID, tx); err != nil {
		return err
	}

	now := time.Now()
	seen := make(map[uint]bool, len(form.ControlIDs))
	controls := make([]*models.ProjectSecurityGateChecklistControl, 0, len(form.ControlIDs))

	for _, controlID := range form.ControlIDs {
		if seen[controlID] {
			err = utility.NewValidationError(
				fmt.Errorf("duplicate control ID: %d", controlID),
				"duplicate",
				"control_id",
			)
			return err
		}
		seen[controlID] = true

		if _, err = psgcs.controlSvc.GetByID(controlID, tx); err != nil {
			return err
		}

		controls = append(controls, &models.ProjectSecurityGateChecklistControl{
			ChecklistID: checklistID,
			ControlID:   controlID,
			CreatedBYID: currentUserID,
			CreatedAt:   now,
		})
	}

	if err = psgcs.repo.DeleteChecklistControls(checklistID, tx); err != nil {
		return err
	}

	if err = psgcs.repo.CreateChecklistControls(controls, tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	go func() {
		activityErr := psgcs.activityLog.LogActivity(
			models.UpdateAction,
			// Log against the parent checklist
			models.ProjectSecurityGateChecklist{}.TableName(),
			checklistID,
			map[string]interface{}{
				"checklist_id": checklistID,
				"control_ids":  form.ControlIDs,
			},
			currentUserID,
			"project", "project_security_gate", "project_security_gate_checklist_control_map", nil,
		)
		if activityErr != nil {
			log.Printf("Failed to log activity: %v", activityErr)
		}
	}()

	return err
}
