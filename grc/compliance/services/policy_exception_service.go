package services

import (
	"errors"
	"fmt"
	"jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"

	documentSvc "jethur/grc/documents/services"
	logServices "jethur/grc/logs/services"
	masterSvc "jethur/grc/masters/services"
	organizationSvc "jethur/grc/organization/services"
	settingsSvc "jethur/grc/settings/services"
)

type PolicyExceptionServiceInterface interface {
	GetAll(params *dtos.PolicyExceptionListQueryParamsDTO) (*dtos.PolicyExceptionListResponseDTO, error)
	GetByID(id uint) (*models.PolicyException, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.PolicyExceptionDetailDTO, error)
	Create(currentUserID uint, dto *dtos.PolicyExceptionFormDTO) (uint, error)
	Update(id uint, dto *dtos.PolicyExceptionFormDTO, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	MapRisks(policyExceptionID uint, currentUserID uint, input *dtos.MapPolicyExceptionRisksDTO) error
	DeletePolicyExceptionRisk(policyExceptionID uint, riskID uint) error
}

type PolicyExceptionService struct {
	repo                        repositories.PolicyExceptionRepositoryInterface
	documentSvc                 documentSvc.DocumentServiceInterface
	workflowStatusSvc           masterSvc.WorkflowStatusServiceInterface
	organizationUnitSvc         settingsSvc.OrganizationUnitServiceInterface
	userSvc                     organizationSvc.UserServiceInterface
	documentClauseSvc           documentSvc.DocumentClauseServiceInterface
	documentRequirementGroupSvc documentSvc.DocumentRequirementGroupServiceInterface
	documentRequirementSvc      documentSvc.DocumentRequirementServiceInterface
	activityLogSvc              logServices.ActivityLogServiceInterface
}

func NewPolicyExceptionService(repo repositories.PolicyExceptionRepositoryInterface,
	documentSvc documentSvc.DocumentServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	organizationUnitSvc settingsSvc.OrganizationUnitServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	documentClauseSvc documentSvc.DocumentClauseServiceInterface,
	documentRequirementGroupSvc documentSvc.DocumentRequirementGroupServiceInterface,
	documentRequirementSvc documentSvc.DocumentRequirementServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) *PolicyExceptionService {
	return &PolicyExceptionService{
		repo:                        repo,
		documentSvc:                 documentSvc,
		workflowStatusSvc:           workflowStatusSvc,
		organizationUnitSvc:         organizationUnitSvc,
		userSvc:                     userSvc,
		documentClauseSvc:           documentClauseSvc,
		documentRequirementGroupSvc: documentRequirementGroupSvc,
		documentRequirementSvc:      documentRequirementSvc,
		activityLogSvc:              activityLogSvc,
	}
}
func (pes *PolicyExceptionService) GetAll(params *dtos.PolicyExceptionListQueryParamsDTO) (*dtos.PolicyExceptionListResponseDTO, error) {
	return pes.repo.GetAll(params)
}

func (pes PolicyExceptionService) GetByID(id uint) (*models.PolicyException, error) {
	policyException, err := pes.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return policyException, nil
}

func (pes PolicyExceptionService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.PolicyExceptionDetailDTO, error) {
	var detailDTO dtos.PolicyExceptionDetailDTO

	policyException, err := pes.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(policyException, &detailDTO)
	if err != nil {
		return nil, err
	}

	// Format date fields if present
	start := policyException.StartDate.Format("2006-01-02")
	detailDTO.StartDateString = start
	end := policyException.EndDate.Format("2006-01-02")
	detailDTO.EndDateString = end
	detailDTO.CreatedAt = policyException.CreatedAt

	if policyException.UpdatedAt != nil {
		detailDTO.UpdatedAt = policyException.UpdatedAt
	}

	// Add additional fields if needed in the detail screen
	// detailDTO.SomeField = policyException.SomeField

	return &detailDTO, nil
}
func (pes *PolicyExceptionService) Create(currentUserID uint, dto *dtos.PolicyExceptionFormDTO) (uint, error) {
	tx := pes.repo.BeginTransaction()

	if err := pes.validatePolicyExceptionFormDTO(dto); err != nil {
		tx.Rollback()
		return 0, err
	}
	dates, err := pes.parsePolicyExceptionDates(dto)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	document, err := pes.documentSvc.GetFullDetailsByID(dto.DocumentID, &currentUserID, tx)
	if err != nil || document.LatestDocumentVersion.ID == 0 {
		tx.Rollback()
		return 0, fmt.Errorf("failed to get latest version for document ID %d", dto.DocumentID)
	}

	workflowStatus, err := pes.workflowStatusSvc.GetByType("pending")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	referenceNumber, err := pes.repo.GenerateReferenceNumber()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	organizationUnitID, err := pes.userSvc.GetUserPrimaryUnitID(dto.RequestedByID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for user"), "not_exist", "organization unit")
	}
	model := &models.PolicyException{
		DocumentID:                 dto.DocumentID,
		DocumentVersionID:          document.LatestDocumentVersion.ID,
		DocumentVersionClauseID:    dto.DocumentVersionClauseID,
		DocumentRequirementGroupID: dto.DocumentRequirementGroupID,
		DocumentRequirementID:      dto.DocumentRequirementID,
		ReferenceNumber:            referenceNumber,
		Requirements:               dto.Requirements,
		Justification:              dto.Justification,
		StartDate:                  dates.StartDate,
		EndDate:                    &dates.EndDate,
		RequestedByID:              dto.RequestedByID,
		OrganizationUnitID:         organizationUnitID,
		WorkflowStatusID:           &workflowStatus.ID,
		CreatedByID:                currentUserID,
	}

	id, err := pes.repo.Create(model, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {
		fullDetails, err := pes.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch policy exception details", map[string]interface{}{"error": err.Error()})
			return
		}

		err = pes.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PolicyException{}.TableName(),
			id,
			fullDetails,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.PolicyException{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (pes *PolicyExceptionService) Update(id uint, dto *dtos.PolicyExceptionFormDTO, currentUserID uint) error {
	existing, err := pes.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	// Validate inputs
	if err := pes.validatePolicyExceptionFormDTO(dto); err != nil {
		return err
	}
	dates, err := pes.parsePolicyExceptionDates(dto)
	if err != nil {
		return err
	}
	tx := pes.repo.BeginTransaction()

	// Get document details for version ID
	document, err := pes.documentSvc.GetFullDetailsByID(dto.DocumentID, &currentUserID, tx)
	if err != nil {
		return err
	}
	if existing.RequestedByID != dto.RequestedByID {
		organizationUnitID, err := pes.userSvc.GetUserPrimaryUnitID(dto.RequestedByID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for requested user "), "not_exist", "organization unit")
		}
		existing.OrganizationUnitID = organizationUnitID
	}
	// Update fields
	existing.DocumentID = dto.DocumentID
	existing.DocumentVersionID = document.LatestDocumentVersion.ID
	existing.DocumentVersionClauseID = dto.DocumentVersionClauseID
	existing.DocumentRequirementGroupID = dto.DocumentRequirementGroupID
	existing.DocumentRequirementID = dto.DocumentRequirementID
	existing.Requirements = dto.Requirements
	existing.Justification = dto.Justification
	existing.StartDate = dates.StartDate
	existing.EndDate = &dates.EndDate
	existing.RequestedByID = dto.RequestedByID
	existing.UpdatedByID = &currentUserID

	if err := pes.repo.Update(existing, tx); err != nil {
		return err
	}
	tx.Commit()

	// Async Activity Logging
	go func() {
		fullDetails, err := pes.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch policy exception for logging", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		if err := pes.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PolicyException{}.TableName(),
			id,
			fullDetails,
			currentUserID,
			nil,
		); err != nil {
			utility.LogError("Failed to log activity for policy exception update", map[string]interface{}{
				"error":    err.Error(),
				"recordID": id,
				"userID":   currentUserID,
			})
		}
	}()

	return nil
}

func (pes *PolicyExceptionService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := pes.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedPolicyExceptions []map[string]interface{}

	for _, id := range ids {
		exception, err := pes.GetFullDetailsByID(id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for _, risk := range exception.Risks {
			err := pes.repo.DeleteSinglePolicyExceptionRisk(id, risk.RiskRegisterID, tx)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		err = pes.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		successfulIds = append(successfulIds, id)
		deletedPolicyExceptions = append(deletedPolicyExceptions, map[string]interface{}{
			"id":          id,
			"document_id": exception.Document.ID,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := pes.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.PolicyException{}.TableName(),
				successfulIds,
				deletedPolicyExceptions,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}()
	}

	return failedIds, nil
}

func (pes *PolicyExceptionService) MapRisks(policyExceptionID uint, currentUserID uint, input *dtos.MapPolicyExceptionRisksDTO) error {
	// Validate PolicyException existence via service
	if _, err := pes.repo.GetByID(policyExceptionID); err != nil {
		return err
	}

	// Validate risk IDs (optional - create a function if needed)
	if err := pes.validateRiskIDs(input.RiskIDs); err != nil {
		return err
	}

	var mappings []models.PolicyExceptionRisk
	for _, riskID := range input.RiskIDs {
		mappings = append(mappings, models.PolicyExceptionRisk{
			PolicyExceptionID: policyExceptionID,
			RiskRegisterID:    riskID,
			CreatedByID:       currentUserID,
		})
	}

	tx := pes.repo.BeginTransaction()

	// Insert new mappings
	if len(mappings) > 0 {
		if err := pes.repo.CreatePolicyExceptionRisks(mappings, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	// Log activities asynchronously
	go func() {
		ids := make([]uint, len(mappings))
		for i, m := range mappings {
			ids[i] = m.ID
		}
		err := pes.activityLogSvc.LogBulkActivities(models.CreateAction, models.PolicyExceptionRisk{}.TableName(), ids, mappings, currentUserID, nil)
		if err != nil {
			utility.LogError("Failed to log PolicyException Risk map activity", map[string]interface{}{"error": err})
		}
	}()

	return nil
}
func (pes *PolicyExceptionService) DeletePolicyExceptionRisk(policyExceptionID uint, riskID uint) error {
	if _, err := pes.repo.GetByID(policyExceptionID); err != nil {
		return err
	}
	if _, err := pes.repo.GetByRiskID(riskID); err != nil {
		return err
	}
	deletedMapping := map[string]interface{}{
		"policy_exception_id": policyExceptionID,
		"risk_id":             riskID,
	}
	tx := pes.repo.BeginTransaction()

	if err := pes.repo.DeleteSinglePolicyExceptionRisk(policyExceptionID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		err := pes.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.PolicyExceptionRisk{}.TableName(),
			0, // optional: use composite key if required
			deletedMapping,
			0, // optional: currentUserID, if available via context
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log PolicyException risk unmapping activity", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
	return nil
}

func (pes *PolicyExceptionService) validatePolicyExceptionFormDTO(dto *dtos.PolicyExceptionFormDTO) error {
	// Check if Document exists
	_, err := pes.documentSvc.GetByID(dto.DocumentID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("document does not exist"), "not_exist", "document")
	}

	// Check if RequestedBy user exists
	_, err = pes.userSvc.GetByID(dto.RequestedByID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("requested by user does not exist"), "not_exist", "requested_by")
	}

	// // Check if OrganizationUnit exists
	// _, err = pes.organizationUnitSvc.GetByID(dto.OrganizationUnitID)
	// if err != nil {
	// 	return utility.NewValidationError(fmt.Errorf("organization unit does not exist"), "not_exist", "organization_unit")
	// }

	// Check optional DocumentVersionClauseID if provided
	if dto.DocumentVersionClauseID != nil {
		_, err = pes.documentClauseSvc.GetFullDetailsByID(dto.DocumentID, *dto.DocumentVersionClauseID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("document version clause does not exist"), "not_exist", "document_version_clause")
		}
	}

	// Check optional DocumentRequirementGroupID if provided
	if dto.DocumentRequirementGroupID != nil {
		_, err = pes.documentRequirementGroupSvc.GetByID(dto.DocumentID, *dto.DocumentRequirementGroupID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("document requirement group does not exist"), "not_exist", "document_requirement_group")
		}
	}

	if dto.DocumentRequirementID != nil {
		_, err = pes.documentRequirementSvc.GetByID(dto.DocumentID, *dto.DocumentRequirementGroupID, *dto.DocumentRequirementID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("document requirement does not exist"), "not_exist", "document_requirement")
		}
	}

	return nil
}

func (pes *PolicyExceptionService) parsePolicyExceptionDates(dto *dtos.PolicyExceptionFormDTO) (*dtos.PolicyExceptionParsedDates, error) {
	parsed := &dtos.PolicyExceptionParsedDates{}
	var err error

	parsed.StartDate, err = utility.ParseMandatoryDate(dto.StartDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
	}

	parsed.EndDate, err = utility.ParseMandatoryDate(dto.EndDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid end date"), "invalid_date", "end date")
	}

	return parsed, nil
}

func (pes *PolicyExceptionService) validateRiskIDs(riskIDs []uint) error {
	if len(riskIDs) == 0 {
		return utility.NewValidationError(errors.New("no risks provided"), "required", "risk_ids")
	}

	for _, riskID := range riskIDs {
		_, err := pes.repo.GetByRiskID(riskID)
		if err != nil {
			return utility.NewValidationError(
				fmt.Errorf("risk with ID %d does not exist", riskID),
				"not_found",
				fmt.Sprintf("risk_id_%d", riskID),
			)
		}
	}
	return nil
}
