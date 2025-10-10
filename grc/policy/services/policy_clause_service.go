package services

import (
	"errors"
	"fmt"
	documentServices "jethur/grc/documents/services"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/grc/policy/repositories"
	"jethur/utility"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PolicyClauseServiceInterface interface {
	GetAll(params *policyDTOs.PolicyClauseListQueryParamsDTO, documentId uint) ([]*policyDTOs.PolicyClauseListDTO, error)
	GetFullDetailsByID(documentId uint, clauseId uint) (*policyDTOs.PolicyClauseDetailsDTO, error)
	Create(documentId uint, currentUserId uint, policyClauseFormDTO *policyDTOs.PolicyClauseFormDTO) (uint, error)
	Update(documentId uint, clauseId uint, currentUserId uint, policyClauseFormDTO *policyDTOs.PolicyClauseFormDTO) error
	Delete(clauseIDs []uint, currentUserID uint) ([]uint, error)
	GetByDocumentVersionId(documentId uint, documentVersionId uint, tx ...utility.DatabaseTransaction) ([]*policyDTOs.PolicyClauseListDTO, error)

	GetDocumentClauseHistories(documentId uint, clauseId uint) ([]*policyDTOs.ClauseHistoryDTO, error)
}

type PolicyClauseService struct {
	repo              repositories.PolicyClauseRepositoryInterface
	documentSvc       documentServices.DocumentServiceInterface
	documentStatusSvc documentServices.DocumentStatusServiceInterface
	activityLogSvc    services.ActivityLogServiceInterface
}

func NewPolicyClauseService(
	repo repositories.PolicyClauseRepositoryInterface,
	documentSvc documentServices.DocumentServiceInterface,
	documentStatusSvc documentServices.DocumentStatusServiceInterface,
	activityLogSvc services.ActivityLogServiceInterface) PolicyClauseServiceInterface {
	return &PolicyClauseService{
		repo:              repo,
		documentSvc:       documentSvc,
		documentStatusSvc: documentStatusSvc,
		activityLogSvc:    activityLogSvc,
	}
}

func (dcs *PolicyClauseService) GetAll(params *policyDTOs.PolicyClauseListQueryParamsDTO, documentId uint) ([]*policyDTOs.PolicyClauseListDTO, error) {
	// Get full document details (to know latest version ID)
	document, err := dcs.documentSvc.GetFullDetailsByID(documentId, nil)
	if err != nil {
		return nil, err
	}

	// Fetch all clauses (flattened DTOs with user info)
	clauses, err := dcs.repo.GetAll(params, document.LatestDocumentVersion.ID)
	if err != nil {
		return nil, err
	}

	// Map ID -> DTO
	dtoMap := make(map[uint]*policyDTOs.PolicyClauseListDTO)
	var roots []*policyDTOs.PolicyClauseListDTO

	for _, clause := range clauses {
		// Since repo.GetAll already returns PolicyClauseListDTO with aliased fields,
		// we can just copy it directly
		dto := &policyDTOs.PolicyClauseListDTO{
			ID:                          clause.ID,
			ClauseNumber:                clause.ClauseNumber,
			ClauseTitleEn:               clause.ClauseTitleEn,
			ClauseTitleAr:               clause.ClauseTitleAr,
			Description:                 clause.Description,
			CreatedUserID:               clause.CreatedUserID,
			CreatedUserFirstNameEn:      clause.CreatedUserFirstNameEn,
			CreatedUserFirstNameAr:      clause.CreatedUserFirstNameAr,
			CreatedUserLastNameEn:       clause.CreatedUserLastNameEn,
			CreatedUserLastNameAr:       clause.CreatedUserLastNameAr,
			CreatedUserInitial:          clause.CreatedUserInitial,
			CreatedUserColor:            clause.CreatedUserColor,
			CreatedUserEmail:            clause.CreatedUserEmail,
			CreatedUserProfileImagePath: clause.CreatedUserProfileImagePath,
			UpdatedUserID:               clause.UpdatedUserID,
			UpdatedUserFirstNameEn:      clause.UpdatedUserFirstNameEn,
			UpdatedUserFirstNameAr:      clause.UpdatedUserFirstNameAr,
			UpdatedUserLastNameEn:       clause.UpdatedUserLastNameEn,
			UpdatedUserLastNameAr:       clause.UpdatedUserLastNameAr,
			UpdatedUserInitial:          clause.UpdatedUserInitial,
			UpdatedUserColor:            clause.UpdatedUserColor,
			UpdatedUserEmail:            clause.UpdatedUserEmail,
			UpdatedUserProfileImagePath: clause.UpdatedUserProfileImagePath,
			CreatedAt:                   clause.CreatedAt,
			UpdatedAt:                   clause.UpdatedAt,
			DocumentVersionClauseID:     clause.DocumentVersionClauseID, // important for tree building
			CommentsCount:               clause.CommentsCount,
		}

		dtoMap[clause.ID] = dto
	}

	// Build parent-child hierarchy
	for _, clause := range clauses {
		dto := dtoMap[clause.ID]
		if clause.DocumentVersionClauseID != nil {
			// has parent → attach to parent
			if parent, ok := dtoMap[*clause.DocumentVersionClauseID]; ok {
				parent.Children = append(parent.Children, dto)
			}
		} else {
			// no parent → root
			roots = append(roots, dto)
		}
	}

	return roots, nil
}

func (dcs *PolicyClauseService) GetFullDetailsByID(documentId uint, clauseId uint) (*policyDTOs.PolicyClauseDetailsDTO, error) {
	var policyClauseDetailsDTO policyDTOs.PolicyClauseDetailsDTO

	clause, err := dcs.repo.GetFullDetailsByID(clauseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, clauseId)
		}
		return nil, err
	}

	err = mapstructure.Decode(clause, &policyClauseDetailsDTO)
	if err != nil {
		return nil, err
	}

	return &policyClauseDetailsDTO, nil
}

func (dcs *PolicyClauseService) Create(documentId uint, currentUserId uint, form *policyDTOs.PolicyClauseFormDTO) (uint, error) {
	tx := dcs.repo.BeginTransaction()

	var clauseNo string

	// Get latest document version
	document, err := dcs.documentSvc.GetFullDetailsByID(documentId, nil, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if form.PreviousClauseID != nil {
		afterClause, err := dcs.repo.GetByID(*form.PreviousClauseID, tx)
		if err != nil || afterClause == nil {
			tx.Rollback()
			return 0, fmt.Errorf("invalid after_clause_id")
		}

		if form.ParentClauseID != nil {
			parentClause, err := dcs.repo.GetByID(*form.ParentClauseID, tx)
			if err != nil || parentClause == nil {
				tx.Rollback()
				return 0, fmt.Errorf("invalid parent_clause_id")
			}

			if afterClause.DocumentVersionClauseID == nil ||
				*afterClause.DocumentVersionClauseID != *form.ParentClauseID {
				tx.Rollback()
				return 0, fmt.Errorf("after_clause_id must belong to parent_clause_id")
			}
		}

		if err := dcs.repo.ShiftClausesAfter(
			document.LatestDocumentVersion.ID,
			*afterClause.ClauseNumber,
			tx,
		); err != nil {
			tx.Rollback()
			return 0, err
		}

		parts := strings.Split(*afterClause.ClauseNumber, ".")
		last, _ := strconv.Atoi(parts[len(parts)-1])
		parts[len(parts)-1] = fmt.Sprintf("%d", last+1)
		clauseNo = strings.Join(parts, ".")

	} else if form.ParentClauseID != nil {
		parentClause, err := dcs.repo.GetByID(*form.ParentClauseID, tx)
		if err != nil || parentClause == nil {
			tx.Rollback()
			return 0, fmt.Errorf("invalid parent_clause_id")
		}

		nextSeq, err := dcs.repo.GetNextClauseSequence(*form.ParentClauseID, documentId, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		clauseNo = fmt.Sprintf("%s.%d", *parentClause.ClauseNumber, nextSeq)

	} else {
		nextSeq, err := dcs.repo.GetNextTopLevelClauseSequence(documentId, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		clauseNo = fmt.Sprintf("%d", nextSeq)
	}

	newClause := &models.DocumentVersionClause{
		DocumentVersionID:       document.LatestDocumentVersion.ID,
		ClauseNumber:            &clauseNo,
		ClauseTitleEn:           form.ClauseTitleEn,
		ClauseTitleAr:           form.ClauseTitleAr,
		Description:             form.Description,
		CreatedBYID:             currentUserId,
		DocumentVersionClauseID: form.ParentClauseID,
	}

	id, err := dcs.repo.Create(newClause, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return id, nil
}

func (dcs *PolicyClauseService) Update(documentId uint, clauseId uint, currentUserId uint, policyClauseFormDTO *policyDTOs.PolicyClauseFormDTO) error {
	tx := dcs.repo.BeginTransaction()
	existingClause, err := dcs.repo.GetByID(clauseId, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(fmt.Errorf("document clause not found"), clauseId)
		}
		return err
	}

	oldPolicyClause, err := dcs.GetFullDetailsByID(documentId, clauseId)
	if err != nil {
		return err
	}

	document, err := dcs.documentSvc.GetFullDetailsByID(documentId, nil, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingClause.DocumentVersionID = document.LatestDocumentVersion.ID
	existingClause.ClauseTitleEn = policyClauseFormDTO.ClauseTitleEn
	existingClause.ClauseTitleAr = policyClauseFormDTO.ClauseTitleAr
	existingClause.Description = policyClauseFormDTO.Description
	existingClause.UpdatedBYID = &currentUserId

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

	// Prepare log data
	changedData := make(map[string]interface{})

	//log changes
	oldData := map[string]interface{}{
		"clause_number": oldPolicyClause.ClauseNumber,
		"title_en":      oldPolicyClause.ClauseTitleEn,
		"description":   oldPolicyClause.Description,
	}

	newData := map[string]interface{}{
		"clause_number": existingClauseFullDetails.ClauseNumber,
		"title_en":      existingClauseFullDetails.ClauseTitleEn,
		"description":   existingClauseFullDetails.Description,
	}
	// Compare other fields and add them to changedData
	fieldChanges := utility.CompareFieldChanges(oldData, newData)
	for k, v := range fieldChanges {
		changedData[k] = v
	}
	inReviewStatus, err := dcs.documentStatusSvc.GetByType("in-review")
	if err != nil {
		return err
	}
	if document.DocumentStatus.ID == int(inReviewStatus.ID) {
		var clauseEdit []*models.DocumentVersionClauseHistory
		for key, raw := range fieldChanges {
			data := raw.(map[string]interface{})
			oldVal := data["old"].(string)
			newVal := data["new"].(string)
			newClauseHistory := &models.DocumentVersionClauseHistory{
				DocumentVersionClauseID: clauseId,
				OldValue:                &oldVal,
				NewValue:                &newVal,
				ChangedField:            &key,
				CreatedBYID:             currentUserId,
			}
			clauseEdit = append(clauseEdit, newClauseHistory)
		}
		if len(clauseEdit) > 0 {
			err = dcs.repo.CreateClauseHistory(clauseEdit, tx)
			if err != nil {
				tx.Rollback()
				return err

			}
		}
	}

	tx.Commit()
	go func() {

		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := dcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.DocumentVersionClause{}.TableName(),
				clauseId,
				changedData,
				currentUserId,
				"", "", "", nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.DocumentVersionClause{}.TableName(),
					"recordID":  clauseId,
					"userID":    currentUserId,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()
	return nil
}

func (dcs *PolicyClauseService) Delete(clauseIDs []uint, currentUserID uint) ([]uint, error) {
	tx := dcs.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedClause []map[string]interface{}

	for _, id := range clauseIDs {
		clause, err := dcs.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// 1. Check subclauses
		hasSubs, err := dcs.repo.HasSubClauses(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if hasSubs {
			// Cannot delete if subclauses exist
			failedIds = append(failedIds, id)
			continue
		}

		// 2. Delete clause
		if err := dcs.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}

		// 3. Shift sibling clause numbers
		if err := dcs.repo.ShiftClauseNumberAfterDeleting(clause.DocumentVersionID, clause.DocumentVersionClauseID, *clause.ClauseNumber, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedClause = append(deletedClause, map[string]interface{}{
			"document_id":         clause.DocumentVersion.DocumentID,
			"document_version_id": clause.DocumentVersionID,
			"clause_id":           clause.ID,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	// Async activity log
	if len(successfulIds) > 0 {
		go func() {
			err := dcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.DocumentVersionClause{}.TableName(),
				successfulIds,
				deletedClause,
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

func (dcs *PolicyClauseService) GetByDocumentVersionId(documentId uint, documentVersionId uint, tx ...utility.DatabaseTransaction) ([]*policyDTOs.PolicyClauseListDTO, error) {
	// Get full document details (to know latest version ID)
	_, err := dcs.documentSvc.GetFullDetailsByID(documentId, nil, tx...)
	if err != nil {
		return nil, err
	}

	// Fetch all clauses (flattened DTOs with user info)
	clauses, err := dcs.repo.GetByDocumentVersionId(documentId, documentVersionId)
	if err != nil {
		return nil, err
	}

	// Map ID -> DTO
	dtoMap := make(map[uint]*policyDTOs.PolicyClauseListDTO)
	var roots []*policyDTOs.PolicyClauseListDTO

	for _, clause := range clauses {
		// Since repo.GetAll already returns PolicyClauseListDTO with aliased fields,
		// we can just copy it directly

		if clause.Description == nil {
			clause.Description = nil
		}

		dto := &policyDTOs.PolicyClauseListDTO{
			ID:                          clause.ID,
			ClauseNumber:                clause.ClauseNumber,
			ClauseTitleEn:               clause.ClauseTitleEn,
			ClauseTitleAr:               clause.ClauseTitleAr,
			Description:                 clause.Description,
			CreatedUserID:               clause.CreatedUserID,
			CreatedUserFirstNameEn:      clause.CreatedUserFirstNameEn,
			CreatedUserFirstNameAr:      clause.CreatedUserFirstNameAr,
			CreatedUserLastNameEn:       clause.CreatedUserLastNameEn,
			CreatedUserLastNameAr:       clause.CreatedUserLastNameAr,
			CreatedUserInitial:          clause.CreatedUserInitial,
			CreatedUserColor:            clause.CreatedUserColor,
			CreatedUserEmail:            clause.CreatedUserEmail,
			CreatedUserProfileImagePath: clause.CreatedUserProfileImagePath,
			UpdatedUserID:               clause.UpdatedUserID,
			UpdatedUserFirstNameEn:      clause.UpdatedUserFirstNameEn,
			UpdatedUserFirstNameAr:      clause.UpdatedUserFirstNameAr,
			UpdatedUserLastNameEn:       clause.UpdatedUserLastNameEn,
			UpdatedUserLastNameAr:       clause.UpdatedUserLastNameAr,
			UpdatedUserInitial:          clause.UpdatedUserInitial,
			UpdatedUserColor:            clause.UpdatedUserColor,
			UpdatedUserEmail:            clause.UpdatedUserEmail,
			UpdatedUserProfileImagePath: clause.UpdatedUserProfileImagePath,
			CreatedAt:                   clause.CreatedAt,
			UpdatedAt:                   clause.UpdatedAt,
			DocumentVersionClauseID:     clause.DocumentVersionClauseID, // important for tree building
		}

		dtoMap[clause.ID] = dto
	}

	// Build parent-child hierarchy
	for _, clause := range clauses {
		dto := dtoMap[clause.ID]
		if clause.DocumentVersionClauseID != nil {
			// has parent → attach to parent
			if parent, ok := dtoMap[*clause.DocumentVersionClauseID]; ok {
				parent.Children = append(parent.Children, dto)
			}
		} else {
			// no parent → root
			roots = append(roots, dto)
		}
	}

	return roots, nil
}

func (pcs *PolicyClauseService) GetDocumentClauseHistories(documentId uint, clauseId uint) ([]*policyDTOs.ClauseHistoryDTO, error) {
	return pcs.repo.GetClauseHistories(clauseId)
}
