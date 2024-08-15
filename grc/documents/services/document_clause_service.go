package services

import (
	"errors"
	"fmt"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"
	"strings"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentClauseServiceInterface interface {
	GetAll(params *documentDTOs.DocumentClauseListQueryParamsDTO, documentId uint) ([]*documentDTOs.DocumentClauseListDTO, error)
	GetFullDetailsByID(documentId uint, clauseId uint) (*documentDTOs.DocumentClauseDetailsDTO, error)
	Create(documentId uint, currentUserId uint, documentClauseFormDTO *documentDTOs.DocumentClauseFormDTO) (uint, error)
	Update(documentId uint, clauseId uint, currentUserId uint, documentClauseFormDTO *documentDTOs.DocumentClauseFormDTO) error
	Delete(clauseIDs []uint, currentUserID uint) ([]uint, error)
}

type DocumentClauseService struct {
	repo           repositories.DocumentClauseRepositoryInterface
	documentSvc    DocumentServiceInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewDocumentClauseService(
	repo repositories.DocumentClauseRepositoryInterface,
	documentSvc DocumentServiceInterface,
	activityLogSvc services.ActivityLogServiceInterface) DocumentClauseServiceInterface {
	return &DocumentClauseService{
		repo:           repo,
		documentSvc:    documentSvc,
		activityLogSvc: activityLogSvc,
	}
}

func (dcs *DocumentClauseService) GetAll(params *documentDTOs.DocumentClauseListQueryParamsDTO, documentId uint) ([]*documentDTOs.DocumentClauseListDTO, error) {
	document, err := dcs.documentSvc.GetFullDetailsByID(documentId, nil)
	if err != nil {
		return nil, err
	}

	clauses, err := dcs.repo.GetAll(params, document.LatestDocumentVersion.ID)
	if err != nil {
		return nil, err
	}

	dtos := []*documentDTOs.DocumentClauseListDTO{}

	for _, clause := range clauses {
		dto := &documentDTOs.DocumentClauseListDTO{
			ID:                          clause.ID,
			ClauseNumber:                clause.ClauseNumber,
			ClauseTitleEn:               clause.ClauseTitleEn,
			ClauseTitleAr:               clause.ClauseTitleAr,
			Description:                 clause.Description,
			CreatedUserFirstNameEn:      clause.CreatedByUser.FirstNameEn,
			CreatedUserFirstNameAr:      clause.CreatedByUser.FirstNameAr,
			CreatedUserLastNameEn:       clause.CreatedByUser.LastNameEn,
			CreatedUserLastNameAr:       clause.CreatedByUser.LastNameAr,
			CreatedUserInitial:          clause.CreatedByUser.Initial,
			CreatedUserColor:            clause.CreatedByUser.Color,
			CreatedUserEmail:            clause.CreatedByUser.Email,
			CreatedUserProfileImagePath: clause.CreatedByUser.ProfileImagePath,
			CreatedAt:                   clause.CreatedAt,
			UpdatedAt:                   clause.UpdatedAt,
			Controls:                    []*organizationDTOs.ControlMinimalListDTO{}, // Initialize the slice to prevent nil slice in JSON
		}

		if clause.UpdatedByUser != nil {
			dto.UpdatedUserFirstNameEn = clause.UpdatedByUser.FirstNameEn
			dto.UpdatedUserFirstNameAr = clause.UpdatedByUser.FirstNameAr
			dto.UpdatedUserLastNameEn = clause.UpdatedByUser.LastNameEn
			dto.UpdatedUserLastNameAr = clause.UpdatedByUser.LastNameAr
			dto.UpdatedUserInitial = clause.UpdatedByUser.Initial
			dto.UpdatedUserColor = clause.UpdatedByUser.Color
			dto.UpdatedUserEmail = clause.UpdatedByUser.Email
			dto.UpdatedUserProfileImagePath = clause.UpdatedByUser.ProfileImagePath
		}

		// Map each DocumentControl to ControlMinimalListDTO
		for _, documentControl := range clause.DocumentControls {
			controlDto := &organizationDTOs.ControlMinimalListDTO{
				ReferenceNumber: documentControl.Control.ReferenceNumber, // Assuming the ReferenceNumber is the string ID
				Code:            documentControl.Control.Code,            // Ensure the Code field exists and is of type *string in DocumentControl model
				TitleEn:         documentControl.Control.TitleEn,
				TitleAr:         documentControl.Control.TitleAr,
			}
			dto.Controls = append(dto.Controls, controlDto)
		}

		dtos = append(dtos, dto)
	}
	return dtos, nil
}

func (dcs *DocumentClauseService) GetFullDetailsByID(documentId uint, clauseId uint) (*documentDTOs.DocumentClauseDetailsDTO, error) {
	var documentClauseDetailsDTO documentDTOs.DocumentClauseDetailsDTO

	clause, err := dcs.repo.GetFullDetailsByID(clauseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, clauseId)
		}
		return nil, err
	}

	err = mapstructure.Decode(clause, &documentClauseDetailsDTO)
	if err != nil {
		return nil, err
	}

	return &documentClauseDetailsDTO, nil
}

func (dcs *DocumentClauseService) Create(documentId uint, currentUserId uint, documentClauseFormDTO *documentDTOs.DocumentClauseFormDTO) (uint, error) {
	tx := dcs.repo.BeginTransaction()
	document, err := dcs.documentSvc.GetFullDetailsByID(documentId, nil, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if documentClauseFormDTO.ClauseNumber != nil {
		exists, err := dcs.repo.ExistsByClauseNumber(*documentClauseFormDTO.ClauseNumber, document.LatestDocumentVersion.ID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if exists {
			tx.Rollback()
			return 0, utility.NewValidationError(fmt.Errorf("the clause number already exists"), "unique", "clause number")
		}
	}

	documentVersionClause := &models.DocumentVersionClause{
		DocumentVersionID: document.LatestDocumentVersion.ID,
		ClauseNumber:      documentClauseFormDTO.ClauseNumber,
		ClauseTitleEn:     documentClauseFormDTO.ClauseTitleEn,
		ClauseTitleAr:     documentClauseFormDTO.ClauseTitleAr,
		Description:       documentClauseFormDTO.Description,
		CreatedBYID:       currentUserId,
	}

	if documentClauseFormDTO.ParantClauseID != nil {
		parantClause, err := dcs.repo.GetByID(*documentClauseFormDTO.ParantClauseID, tx)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, utility.NewValidationError(fmt.Errorf("the parant clause not found"), "not_exist", "clause")
			}
			return 0, err
		}
		if parantClause.ClauseNumber != nil {
			if !isValidChildClauseNumber(*parantClause.ClauseNumber, *documentClauseFormDTO.ClauseNumber) {
				return 0, utility.NewValidationError(fmt.Errorf("invalid clause number '%s' for the specified parent '%s'", *documentClauseFormDTO.ClauseNumber, *parantClause.ClauseNumber), "invalid_clause_number", "clause number")
			}
		}

		documentVersionClause.DocumentVersionClauseID = &parantClause.ID
	}

	id, err := dcs.repo.Create(documentVersionClause, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if documentClauseFormDTO.ControlIds != nil {

		var controls []*models.DocumentControl
		for _, controlID := range documentClauseFormDTO.ControlIds {
			control := &models.DocumentControl{
				DocumentID:              documentId,
				DocumentVersionID:       &document.LatestDocumentVersion.ID,
				DocumentVersionClauseID: &id,
				ControlID:               controlID,
				CreatedBYID:             currentUserId,
			}
			controls = append(controls, control)
		}

		err = dcs.repo.CreateDocumentControls(controls, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := dcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.DocumentVersionClause{}.TableName(),
			id,
			documentClauseFormDTO,
			currentUserId,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.DocumentVersionClause{}.TableName(),
				"recordID":  id,
				"userID":    currentUserId,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (dcs *DocumentClauseService) Update(documentId uint, clauseId uint, currentUserId uint, documentClauseFormDTO *documentDTOs.DocumentClauseFormDTO) error {
	tx := dcs.repo.BeginTransaction()
	existingClause, err := dcs.repo.GetByID(clauseId, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(fmt.Errorf("document clause not found"), clauseId)
		}
		return err
	}

	document, err := dcs.documentSvc.GetFullDetailsByID(documentId, nil, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingClause.ClauseNumber != nil {
		if *existingClause.ClauseNumber != *documentClauseFormDTO.ClauseNumber {
			exists, err := dcs.repo.ExistsByClauseNumber(*documentClauseFormDTO.ClauseNumber, document.LatestDocumentVersion.ID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			if exists {
				tx.Rollback()
				return utility.NewValidationError(fmt.Errorf("the clause number already exists"), "unique", "clause number")
			}
		}
	} else {
		exists, err := dcs.repo.ExistsByClauseNumber(*documentClauseFormDTO.ClauseNumber, document.LatestDocumentVersion.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		if exists {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("the clause number already exists"), "unique", "clause number")
		}
	}

	existingClause.DocumentVersionID = document.LatestDocumentVersion.ID
	existingClause.ClauseNumber = documentClauseFormDTO.ClauseNumber
	existingClause.ClauseTitleEn = documentClauseFormDTO.ClauseTitleEn
	existingClause.ClauseTitleAr = documentClauseFormDTO.ClauseTitleAr
	existingClause.Description = documentClauseFormDTO.Description
	existingClause.UpdatedBYID = &currentUserId

	if documentClauseFormDTO.ParantClauseID != nil {
		parantClause, err := dcs.repo.GetByID(*documentClauseFormDTO.ParantClauseID, tx)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utility.NewValidationError(fmt.Errorf("parent clause not found"), "not_exist", "clause number")
			}
			return err
		}
		if parantClause.ClauseNumber != nil {
			if !isValidChildClauseNumber(*parantClause.ClauseNumber, *documentClauseFormDTO.ClauseNumber) {
				return utility.NewValidationError(fmt.Errorf("invalid clause number '%s' for the specified parent '%s'", *documentClauseFormDTO.ClauseNumber, *parantClause.ClauseNumber), "invalid_clause_number", "clause number")
			}
		}
		existingClause.DocumentVersionClauseID = &parantClause.ID
	}

	err = dcs.repo.Update(existingClause, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingClauseFullDetails, err := dcs.repo.GetFullDetailsByID(clauseId, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Initialize control maps
	existingControlIDs := make(map[uint]bool)
	if existingClauseFullDetails.DocumentControls != nil {
		for _, ctrl := range existingClauseFullDetails.DocumentControls {
			existingControlIDs[*ctrl.ControlID] = true
		}
	}

	// Prepare to add new controls
	newControlIDs := make(map[uint]bool)
	for _, ctrlID := range documentClauseFormDTO.ControlIds {
		newControlIDs[*ctrlID] = true
	}

	// Add new controls
	var controlsToAdd []*models.DocumentControl
	for _, newControlID := range documentClauseFormDTO.ControlIds {
		if _, exists := existingControlIDs[*newControlID]; !exists {
			control := &models.DocumentControl{
				DocumentID:              documentId,
				DocumentVersionID:       &document.LatestDocumentVersion.ID,
				DocumentVersionClauseID: &clauseId,
				ControlID:               newControlID,
				CreatedBYID:             currentUserId,
			}
			controlsToAdd = append(controlsToAdd, control)
		}
	}

	// Save new controls
	if len(controlsToAdd) > 0 {
		if err := dcs.repo.CreateDocumentControls(controlsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove controls not in the new list
	var deleteControlIDs []*uint
	for _, existingControl := range existingClauseFullDetails.DocumentControls {
		if _, exists := newControlIDs[*existingControl.ControlID]; !exists {
			deleteControlIDs = append(deleteControlIDs, &existingControl.ID)
		}
	}

	// Delete controls not in the new list
	if len(deleteControlIDs) > 0 {
		if err := dcs.repo.DeleteDocumentControls(deleteControlIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (dcs *DocumentClauseService) Delete(clauseIDs []uint, currentUserID uint) ([]uint, error) {
	tx := dcs.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint

	for _, id := range clauseIDs {
		clause, err := dcs.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(clause.DocumentControls) > 0 {
			var controlIDs []*uint
			for _, control := range clause.DocumentControls {
				controlIDs = append(controlIDs, &control.ID)
			}
			if err := dcs.repo.DeleteDocumentControls(controlIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if err := dcs.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := dcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.DocumentVersionClause{}.TableName(),
				successfulIds,
				nil,
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

func isValidChildClauseNumber(parentNumber, childNumber string) bool {
	return strings.HasPrefix(childNumber, parentNumber+".")
}
