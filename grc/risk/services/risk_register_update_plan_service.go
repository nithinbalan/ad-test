package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskRegisterUpdatePlanServiceInterface interface {
	GetAll(params *riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO) (*riskDTOs.RiskRegisterUpdatePlanListResponseDTO, error)
	GetByID(id uint) (*models.RiskRegisterUpdatePlan, error)
	GetFullDetailsByID(id uint) (*riskDTOs.RiskRegisterUpdatePlanDetailDTO, error)
	Create(currentUserID uint, formDTO *riskDTOs.RiskRegisterUpdatePlanFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *riskDTOs.RiskRegisterUpdatePlanFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	// ERM Department Evaluation Operations
	CreateERMDepartmentEvaluation(planID uint, currentUserID uint) (*riskDTOs.AutoFetchRisksResponseDTO, error)
	UpdateERMDepartmentEvaluation(planID uint, currentUserID uint, formDTO *riskDTOs.MapRisksToUpdatePlanDTO) (*riskDTOs.MapRisksToUpdatePlanResponseDTO, error)
	GetERMDepartmentEvaluation(planID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationListResponseDTO, error)
	DeleteERMDepartmentEvaluation(id uint, currentUserID uint) error

	// Introductory Meeting Document Operations
	GetIntroductoryMeetingDocuments(planID uint) ([]*riskDTOs.IntroductoryMeetingDocumentDTO, error)
	CreateIntroductoryMeetingDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error
	GetIntroductoryMeetingDocumentPath(planID uint, token string) (string, error)

	// Business Context Operations
	CreateBusinessContext(planID uint, currentUserID uint) error
	GetBusinessContextByPlanID(planID uint) ([]*riskDTOs.BusinessContextDetailDTO, error)
	UpdateBusinessContextStatus(id uint, currentUserID uint, formDTO *riskDTOs.UpdateBusinessContextStatusFormDTO) error
}

type RiskRegisterUpdatePlanService struct {
	repo                riskRepositories.RiskRegisterUpdatePlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	userSvc             organizationServices.UserServiceInterface
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface
}

func NewRiskRegisterUpdatePlanService(
	repo riskRepositories.RiskRegisterUpdatePlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
) RiskRegisterUpdatePlanServiceInterface {
	return &RiskRegisterUpdatePlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		userSvc:             userSvc,
		organizationUnitSvc: organizationUnitSvc,
	}
}

func (s *RiskRegisterUpdatePlanService) GetAll(params *riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO) (*riskDTOs.RiskRegisterUpdatePlanListResponseDTO, error) {
	return s.repo.GetAll(params)
}

func (s *RiskRegisterUpdatePlanService) GetByID(id uint) (*models.RiskRegisterUpdatePlan, error) {
	plan, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return plan, nil
}

func (s *RiskRegisterUpdatePlanService) GetFullDetailsByID(id uint) (*riskDTOs.RiskRegisterUpdatePlanDetailDTO, error) {
	var planDTO riskDTOs.RiskRegisterUpdatePlanDetailDTO

	plan, err := s.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(plan, &planDTO)
	if err != nil {
		return nil, err
	}

	planDTO.StartDateString = plan.StartDate.Format("2006-01-02")
	planDTO.EndDateString = plan.EndDate.Format("2006-01-02")
	planDTO.CreatedAt = plan.CreatedAt
	planDTO.UpdatedAt = plan.UpdatedAt

	return &planDTO, nil
}

func (s *RiskRegisterUpdatePlanService) Create(currentUserID uint, formDTO *riskDTOs.RiskRegisterUpdatePlanFormDTO) (uint, error) {

	if err := s.validateRiskRegisterUpdatePlanFormDTO(formDTO); err != nil {
		return 0, err
	}

	dates, err := s.parseRiskRegisterUpdatePlanDates(formDTO, currentUserID)
	if err != nil {
		return 0, err
	}

	tx := s.repo.BeginTransaction()

	plan := &models.RiskRegisterUpdatePlan{
		DivisionID:                     formDTO.DivisionID,
		BusinessUnitID:                 formDTO.BusinessUnitID,
		RiskRegisterUpdatePlanStatusID: formDTO.RiskRegisterUpdatePlanStatusID,
		NumberOfOperatingEntities:      formDTO.NumberOfOperatingEntities,
		Workdays:                       formDTO.Workdays,
		RiskChampionID:                 formDTO.RiskChampionID,
		StartDate:                      dates.StartDate,
		EndDate:                        dates.EndDate,
		Notes:                          formDTO.Notes,
		CreatedBYID:                    currentUserID,
		CreatedAt:                      time.Now(),
	}

	planID, err := s.repo.Create(plan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if len(formDTO.TeamMemberIds) > 0 {
		teamMembers := make([]*models.RiskRegisterUpdatePlanTeamMember, 0)
		for _, userID := range formDTO.TeamMemberIds {
			if userID != nil {
				member := &models.RiskRegisterUpdatePlanTeamMember{
					RiskRegisterUpdatePlanID: planID,
					UserID:                   *userID,
					CreatedBYID:              &currentUserID,
					CreatedAt:                time.Now(),
				}
				teamMembers = append(teamMembers, member)
			}
		}

		if err := s.repo.CreateTeamMembers(teamMembers, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := s.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskRegisterUpdatePlan{}.TableName(),
			planID,
			nil,
			currentUserID,
			"risk", "risk_register_update_plan", "create", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskRegisterUpdatePlan{}.TableName(),
				"recordID":  planID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return planID, nil
}

func (s *RiskRegisterUpdatePlanService) Update(id uint, currentUserID uint, formDTO *riskDTOs.RiskRegisterUpdatePlanFormDTO) error {

	if err := s.validateRiskRegisterUpdatePlanFormDTO(formDTO); err != nil {
		return err
	}

	dates, err := s.parseRiskRegisterUpdatePlanDates(formDTO, currentUserID)
	if err != nil {
		return err
	}
	existingPlan, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	tx := s.repo.BeginTransaction()

	oldData, err := s.repo.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	now := time.Now()
	existingPlan.DivisionID = formDTO.DivisionID
	existingPlan.BusinessUnitID = formDTO.BusinessUnitID
	existingPlan.RiskRegisterUpdatePlanStatusID = formDTO.RiskRegisterUpdatePlanStatusID
	existingPlan.NumberOfOperatingEntities = formDTO.NumberOfOperatingEntities
	existingPlan.Workdays = formDTO.Workdays
	existingPlan.RiskChampionID = formDTO.RiskChampionID
	existingPlan.StartDate = dates.StartDate
	existingPlan.EndDate = dates.EndDate
	existingPlan.Notes = formDTO.Notes
	existingPlan.UpdatedBYID = &currentUserID
	existingPlan.UpdatedAt = &now

	if err := s.repo.Update(existingPlan, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteTeamMembersByPlanID(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	if len(formDTO.TeamMemberIds) > 0 {
		teamMembers := make([]*models.RiskRegisterUpdatePlanTeamMember, 0)
		for _, userID := range formDTO.TeamMemberIds {
			if userID != nil {
				member := &models.RiskRegisterUpdatePlanTeamMember{
					RiskRegisterUpdatePlanID: id,
					UserID:                   *userID,
					CreatedBYID:              &currentUserID,
					CreatedAt:                time.Now(),
				}
				teamMembers = append(teamMembers, member)
			}
		}

		if err := s.repo.CreateTeamMembers(teamMembers, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {

		oldData := map[string]interface{}{
			"division_id":      oldData.DivisionID,
			"business_unit_id": oldData.BusinessUnitID,
			"status_id":        oldData.RiskRegisterUpdatePlanStatusID,
			"risk_champion_id": oldData.RiskChampionID,
			"start_date":       oldData.StartDate.Format("2006-01-02"),
			"end_date":         oldData.EndDate.Format("2006-01-02"),
		}

		newData := map[string]interface{}{
			"division_id":      formDTO.DivisionID,
			"business_unit_id": formDTO.BusinessUnitID,
			"status_id":        formDTO.RiskRegisterUpdatePlanStatusID,
			"risk_champion_id": formDTO.RiskChampionID,
			"start_date":       formDTO.StartDate,
			"end_date":         formDTO.EndDate,
		}

		changedData := utility.CompareFieldChanges(oldData, newData)

		if len(changedData) > 0 {
			err := s.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskRegisterUpdatePlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				"risk", "risk_register_update_plan", "update", nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.RiskRegisterUpdatePlan{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}
	}()

	return nil
}

func (s *RiskRegisterUpdatePlanService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := s.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedPlans []map[string]interface{}

	for _, id := range ids {
		plan, err := s.GetFullDetailsByID(id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedPlans = append(deletedPlans, map[string]interface{}{
			"id":               id,
			"division_id":      plan.DivisionID,
			"business_unit_id": plan.BusinessUnitID,
			"risk_champion_id": plan.RiskChampionID,
			"start_date":       plan.StartDateString,
			"end_date":         plan.EndDateString,
		})

		// Delete team members first
		if err := s.repo.DeleteTeamMembersByPlanID(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Delete the plan
		if err := s.repo.Delete(id, tx); err != nil {
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

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := s.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskRegisterUpdatePlan{}.TableName(),
				successfulIds,
				deletedPlans,
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

func (s *RiskRegisterUpdatePlanService) validateRiskRegisterUpdatePlanFormDTO(formDTO *riskDTOs.RiskRegisterUpdatePlanFormDTO) error {
	// Validate division
	if _, err := s.organizationUnitSvc.GetByID(formDTO.DivisionID); err != nil {
		return utility.NewValidationError(fmt.Errorf("division does not exist"), "not_exist", "division")
	}

	// Validate business unit
	if _, err := s.organizationUnitSvc.GetByID(formDTO.BusinessUnitID); err != nil {
		return utility.NewValidationError(fmt.Errorf("business unit does not exist"), "not_exist", "business_unit")
	}

	// Validate risk champion
	if _, err := s.userSvc.GetByID(formDTO.RiskChampionID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk champion does not exist"), "not_exist", "risk_champion")
	}

	// Validate team members if provided
	if len(formDTO.TeamMemberIds) > 0 {
		for _, userID := range formDTO.TeamMemberIds {
			if userID != nil {
				if _, err := s.userSvc.GetByID(*userID); err != nil {
					return utility.NewValidationError(fmt.Errorf("team member user does not exist"), "not_exist", "team_member")
				}
			}
		}
	}

	return nil
}

func (s *RiskRegisterUpdatePlanService) parseRiskRegisterUpdatePlanDates(dto *riskDTOs.RiskRegisterUpdatePlanFormDTO, currentUserID uint) (*riskDTOs.RiskRegisterUpdatePlanParsedDates, error) {
	parsed := &riskDTOs.RiskRegisterUpdatePlanParsedDates{}

	// Parse start date
	if strings.TrimSpace(dto.StartDate) != "" {
		startDate, err := utility.ParseDate(strings.TrimSpace(dto.StartDate))
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid start date: %w", err), "invalid_date", "start_date")
		}
		if startDate != nil {
			parsed.StartDate = *startDate
		}
	} else {
		currentDate, _ := utility.ParseDate(time.Now().Format("2006-01-02"))
		if currentDate != nil {
			parsed.StartDate = *currentDate
		}
	}

	if strings.TrimSpace(dto.EndDate) != "" {
		endDate, err := utility.ParseDate(strings.TrimSpace(dto.EndDate))
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid end date: %w", err), "invalid_date", "end_date")
		}
		if endDate != nil {
			parsed.EndDate = *endDate
		}
	} else {
		currentDate, _ := utility.ParseDate(time.Now().Format("2006-01-02"))
		if currentDate != nil {
			parsed.EndDate = *currentDate
		}
	}

	// Validate that end date is not before start date
	if parsed.EndDate.Before(parsed.StartDate) {
		return nil, utility.NewValidationError(fmt.Errorf("end date cannot be before start date"), "invalid", "end_date")
	}

	return parsed, nil
}

// ERM Department Evaluation Operations

func (s *RiskRegisterUpdatePlanService) CreateERMDepartmentEvaluation(planID uint, currentUserID uint) (*riskDTOs.AutoFetchRisksResponseDTO, error) {
	// Get the plan to retrieve business unit ID
	plan, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, planID)
		}
		return nil, err
	}

	// Get all risks for the business unit
	riskIDs, err := s.repo.GetRisksByBusinessUnitID(plan.BusinessUnitID)
	if err != nil {
		return nil, err
	}

	response := &riskDTOs.AutoFetchRisksResponseDTO{
		TotalRisks:  len(riskIDs),
		MappedCount: 0,
		FailedCount: 0,
		FailedIDs:   []uint{},
	}

	if len(riskIDs) == 0 {
		return response, nil
	}

	tx := s.repo.BeginTransaction()

	// Map each risk
	for _, riskID := range riskIDs {
		// Check if already mapped
		alreadyMapped, err := s.repo.CheckRiskAlreadyMapped(planID, riskID)
		if err != nil {
			response.FailedCount++
			response.FailedIDs = append(response.FailedIDs, riskID)
			continue
		}

		if alreadyMapped {
			continue // Skip already mapped risks
		}

		evaluation := &models.RiskUpdatePlanErmDeptEvaluation{
			RiskRegisterUpdatePlanID: planID,
			RiskRegisterID:           riskID,
			CreatedBYID:              currentUserID,
			CreatedAt:                time.Now(),
		}

		if err := s.repo.CreateERMDepartmentEvaluation(evaluation, tx); err != nil {
			response.FailedCount++
			response.FailedIDs = append(response.FailedIDs, riskID)
			continue
		}

		response.MappedCount++
	}

	if response.FailedCount > 0 && response.MappedCount == 0 {
		tx.Rollback()
		return nil, fmt.Errorf("failed to map any risks")
	}

	tx.Commit()

	// Log activity
	go func() {
		err := s.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskUpdatePlanErmDeptEvaluation{}.TableName(),
			planID,
			map[string]interface{}{
				"total_risks":  response.TotalRisks,
				"mapped_count": response.MappedCount,
			},
			currentUserID,
			"risk", "risk_update_plan_evaluation", "auto_fetch", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action": models.CreateAction,
				"planID": planID,
				"userID": currentUserID,
				"error":  err.Error(),
			})
		}
	}()

	return response, nil
}

func (s *RiskRegisterUpdatePlanService) UpdateERMDepartmentEvaluation(planID uint, currentUserID uint, formDTO *riskDTOs.MapRisksToUpdatePlanDTO) (*riskDTOs.MapRisksToUpdatePlanResponseDTO, error) {
	tx := s.repo.BeginTransaction()

	// Verify plan exists
	if _, err := s.repo.GetByID(planID); err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, planID)
		}
		return nil, err
	}

	// Get existing mapped risk IDs
	existingRiskIDs, err := s.repo.GetExistingMappedRiskIDs(planID, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create maps for comparison
	existingRisks := make(map[uint]bool)
	for _, riskID := range existingRiskIDs {
		existingRisks[riskID] = true
	}

	newRisks := make(map[uint]bool)
	for _, riskID := range formDTO.RiskRegisterIDs {
		newRisks[riskID] = true
	}

	// Find risks to add (in new but not in existing)
	var risksToAdd []*models.RiskUpdatePlanErmDeptEvaluation
	now := time.Now()
	for riskID := range newRisks {
		if !existingRisks[riskID] {
			risksToAdd = append(risksToAdd, &models.RiskUpdatePlanErmDeptEvaluation{
				RiskRegisterUpdatePlanID: planID,
				RiskRegisterID:           riskID,
				CreatedBYID:              currentUserID,
				CreatedAt:                now,
			})
		}
	}

	// Add new risks
	if len(risksToAdd) > 0 {
		for _, evaluation := range risksToAdd {
			if err := s.repo.CreateERMDepartmentEvaluation(evaluation, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Find risks to delete (in existing but not in new)
	var deleteRiskIDs []uint
	for existingID := range existingRisks {
		if !newRisks[existingID] {
			deleteRiskIDs = append(deleteRiskIDs, existingID)
		}
	}

	// Delete removed risks
	if len(deleteRiskIDs) > 0 {
		if err := s.repo.DeleteERMDepartmentEvaluationsByRiskIDs(planID, deleteRiskIDs, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	// Prepare response
	response := &riskDTOs.MapRisksToUpdatePlanResponseDTO{
		SuccessCount: len(risksToAdd),
		FailedCount:  0,
		FailedIDs:    []uint{},
	}

	// Log activity
	go func() {
		err := s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskUpdatePlanErmDeptEvaluation{}.TableName(),
			planID,
			map[string]interface{}{
				"added_count":   len(risksToAdd),
				"deleted_count": len(deleteRiskIDs),
			},
			currentUserID,
			"risk", "risk_update_plan_evaluation", "update", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action": models.UpdateAction,
				"planID": planID,
				"userID": currentUserID,
				"error":  err.Error(),
			})
		}
	}()

	return response, nil
}

func (s *RiskRegisterUpdatePlanService) GetERMDepartmentEvaluation(planID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationListResponseDTO, error) {
	// Verify plan exists
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, planID)
		}
		return nil, err
	}

	return s.repo.GetERMDepartmentEvaluation(planID, params)
}

func (s *RiskRegisterUpdatePlanService) DeleteERMDepartmentEvaluation(id uint, currentUserID uint) error {
	tx := s.repo.BeginTransaction()

	if err := s.repo.DeleteERMDepartmentEvaluation(id, tx); err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	tx.Commit()

	// Log activity
	go func() {
		err := s.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.RiskUpdatePlanErmDeptEvaluation{}.TableName(),
			id,
			nil,
			currentUserID,
			"risk", "risk_update_plan_evaluation", "delete", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action": models.DeleteAction,
				"id":     id,
				"userID": currentUserID,
				"error":  err.Error(),
			})
		}
	}()

	return nil
}

// Introductory Meeting Document Operations

func (s *RiskRegisterUpdatePlanService) GetIntroductoryMeetingDocuments(planID uint) ([]*riskDTOs.IntroductoryMeetingDocumentDTO, error) {
	// Verify plan exists
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, planID)
		}
		return nil, err
	}

	// Get documents
	documents, err := s.repo.GetIntroductoryMeetingDocumentsByPlanID(planID)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	var documentDTOs []*riskDTOs.IntroductoryMeetingDocumentDTO
	for _, doc := range documents {
		dto := &riskDTOs.IntroductoryMeetingDocumentDTO{
			ID:           &doc.ID,
			URL:          doc.URL,
			Name:         doc.Name,
			Token:        doc.Token,
			OriginalName: doc.OriginalName,
			Type:         doc.Type,
			Size:         doc.Size,
		}
		documentDTOs = append(documentDTOs, dto)
	}

	return documentDTOs, nil
}

func (s *RiskRegisterUpdatePlanService) CreateIntroductoryMeetingDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error {
	// Verify plan exists
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	tx := s.repo.BeginTransaction()

	// Handle documents
	var newDocuments []*models.RiskUpdatePlanIntroductoryMeetingDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskUpdatePlanIntroductoryMeetingDocument

	now := time.Now()

	for _, doc := range documents {
		// NEW FILE
		if doc.IsNewFile != nil && *doc.IsNewFile {
			encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL(
				"risk",
				"risk-register-update-plan",
				strconv.FormatUint(uint64(planID), 10),
				"introductory-meeting-documents",
			)

			document := &models.RiskUpdatePlanIntroductoryMeetingDocument{
				RiskRegisterUpdatePlanID: planID,
				URL:                      &url,
				Name:                     &doc.FileName,
				Size:                     &doc.Size,
				Token:                    &encryptedFilename,
				OriginalName:             &doc.OriginalName,
				Type:                     &doc.Type,
				CreatedBYID:              &currentUserID,
				CreatedAt:                now,
			}

			newDocuments = append(newDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {
			// DELETE FILE
			existingDocument, err := s.repo.GetIntroductoryMeetingDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Verify document belongs to the plan
			if existingDocument.RiskRegisterUpdatePlanID != planID {
				tx.Rollback()
				return utility.NewResourceNotFoundError(fmt.Errorf("document not found"), *doc.ID)
			}

			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			documentsToDelete = append(documentsToDelete, existingDocument)
		}
	}

	// Create new documents
	if len(newDocuments) > 0 {
		for _, doc := range newDocuments {
			_, err := s.repo.CreateIntroductoryMeetingDocument(doc, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Delete documents
	if len(deleteDocumentIDs) > 0 {
		for _, docID := range deleteDocumentIDs {
			if err := s.repo.DeleteIntroductoryMeetingDocument(*docID, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// File system operations - store new files
	for _, doc := range newDocuments {
		storagePath := filepath.Join("risk", "risk-register-update-plan", strconv.FormatUint(uint64(planID), 10), "introductory-meeting-documents")
		if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	// File system operations - delete files
	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("risk", "risk-register-update-plan", strconv.FormatUint(uint64(planID), 10), "introductory-meeting-documents")
		if doc.Name != nil {
			if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	// Log activity
	go func() {
		err := s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskUpdatePlanIntroductoryMeetingDocument{}.TableName(),
			planID,
			map[string]interface{}{
				"new_documents_count":     len(newDocuments),
				"deleted_documents_count": len(documentsToDelete),
			},
			currentUserID,
			"risk", "risk_register_update_plan", "documents_update", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action": models.UpdateAction,
				"planID": planID,
				"userID": currentUserID,
				"error":  err.Error(),
			})
		}
	}()

	return nil
}

func (s *RiskRegisterUpdatePlanService) GetIntroductoryMeetingDocumentPath(planID uint, token string) (string, error) {
	// Verify plan exists
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", utility.NewResourceNotFoundError(err, planID)
		}
		return "", err
	}

	// Get documents for the plan
	documents, err := s.repo.GetIntroductoryMeetingDocumentsByPlanID(planID)
	if err != nil {
		return "", err
	}

	if len(documents) == 0 {
		return "", fmt.Errorf("no documents found for the plan")
	}

	// Decrypt the token
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	// Construct file path
	storagePath := filepath.Join("risk", "risk-register-update-plan", fmt.Sprintf("%d", planID), "introductory-meeting-documents")
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

// Business Context Operations

func (s *RiskRegisterUpdatePlanService) CreateBusinessContext(planID uint, currentUserID uint) error {
	// Verify plan exists
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get business unit ID from plan
	businessUnitID, err := s.repo.GetBusinessUnitIDByPlanID(planID)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	// Get all business context IDs for the business unit
	businessContextIDs, err := s.repo.GetBusinessContextIDsByBusinessUnit(businessUnitID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Get all process IDs for the business unit
	processIDs, err := s.repo.GetProcessIDsByBusinessUnit(businessUnitID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Get all control IDs for the business unit
	controlIDs, err := s.repo.GetControlIDsByBusinessUnit(businessUnitID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create business context records
	for _, bcID := range businessContextIDs {
		businessContext := &models.RiskUpdatePlanBusinessContext{
			RiskRegisterUpdatePlanID: planID,
			BusinessContextID:        &bcID,
			CreatedBYID:              &currentUserID,
		}
		if err := s.repo.CreateBusinessContext(businessContext, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Create process records
	for _, processID := range processIDs {
		businessContext := &models.RiskUpdatePlanBusinessContext{
			RiskRegisterUpdatePlanID: planID,
			ProcessID:                &processID,
			CreatedBYID:              &currentUserID,
		}
		if err := s.repo.CreateBusinessContext(businessContext, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Create control records
	for _, controlID := range controlIDs {
		businessContext := &models.RiskUpdatePlanBusinessContext{
			RiskRegisterUpdatePlanID: planID,
			ControlID:                &controlID,
			CreatedBYID:              &currentUserID,
		}
		if err := s.repo.CreateBusinessContext(businessContext, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	// Log activity
	go func() {
		err := s.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskUpdatePlanBusinessContext{}.TableName(),
			planID,
			map[string]interface{}{
				"plan_id":          planID,
				"business_unit_id": businessUnitID,
				"total_created":    len(businessContextIDs) + len(processIDs) + len(controlIDs),
			},
			currentUserID,
			"risk", "risk_update_plan_business_context", "create_all", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action": models.CreateAction,
				"planID": planID,
				"userID": currentUserID,
				"error":  err.Error(),
			})
		}
	}()

	return nil
}

func (s *RiskRegisterUpdatePlanService) GetBusinessContextByPlanID(planID uint) ([]*riskDTOs.BusinessContextDetailDTO, error) {

	// Verify plan exists
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, planID)
		}
		return nil, err
	}

	// Get business contexts
	businessContexts, err := s.repo.GetBusinessContextByPlanID(planID)
	if err != nil {
		return nil, err
	}

	var dtos []*riskDTOs.BusinessContextDetailDTO

	for _, businessContext := range businessContexts {
		var dto riskDTOs.BusinessContextDetailDTO

		err := mapstructure.Decode(businessContext, &dto)
		if err != nil {
			return nil, err
		}

		dto.CreatedAt = businessContext.CreatedAt

		dtos = append(dtos, &dto)
	}

	return dtos, nil
}

func (s *RiskRegisterUpdatePlanService) UpdateBusinessContextStatus(id uint, currentUserID uint, formDTO *riskDTOs.UpdateBusinessContextStatusFormDTO) error {
	// Verify business context exists
	_, err := s.repo.GetBusinessContextByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	tx := s.repo.BeginTransaction()

	// Update status
	if err := s.repo.UpdateBusinessContextStatus(id, formDTO.IsBusinessContextUpdated, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity
	go func() {
		err := s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskUpdatePlanBusinessContext{}.TableName(),
			id,
			map[string]interface{}{
				"is_business_context_updated": formDTO.IsBusinessContextUpdated,
			},
			currentUserID,
			"risk", "risk_update_plan_business_context", "update_status", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action": models.UpdateAction,
				"id":     id,
				"userID": currentUserID,
				"error":  err.Error(),
			})
		}
	}()

	return nil
}
