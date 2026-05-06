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

	// ERM Department Evaluation - Auto Fetch All
	CreateERMDepartmentEvaluation(planID uint, currentUserID uint) error

	// ERM Department Evaluation - Risks
	MapRisksToTimeline(planID uint, currentUserID uint, formDTO *riskDTOs.MapRisksToUpdatePlanDTO) error
	GetERMDeptEvaluationRisks(planID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationListResponseDTO, error)

	// ERM Department Evaluation - Processes
	MapProcessesToTimeline(planID uint, currentUserID uint, formDTO *riskDTOs.MapProcessesToUpdatePlanDTO) error
	GetERMDeptEvaluationProcesses(planID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationProcessListResponseDTO, error)

	// ERM Department Evaluation - Controls
	MapControlsToTimeline(planID uint, currentUserID uint, formDTO *riskDTOs.MapControlsToUpdatePlanDTO) error
	GetERMDeptEvaluationControls(planID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationControlListResponseDTO, error)
	// ERM Department Evaluation - Operations
	MarkERMDeptEvaluationAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error

	// ERM Department Evaluation - Documents
	GetERMDeptEvaluationDocuments(planID uint) ([]*riskDTOs.ERMDeptEvaluationDocumentDTO, error)
	CreateERMDeptEvaluationDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error
	GetERMDeptEvaluationDocumentPath(planID uint, token string) (string, error)

	// Introductory Meeting Document Operations
	GetIntroductoryMeetingDocuments(planID uint) ([]*riskDTOs.IntroductoryMeetingDocumentDTO, error)
	CreateIntroductoryMeetingDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error
	GetIntroductoryMeetingDocumentPath(planID uint, token string) (string, error)
	// Introductory Meeting - Operations
	MarkIntroductoryMeetingAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error

	// Self-Awareness Training Document Operations
	GetSelfAwarenessTrainingDocuments(planID uint) ([]*riskDTOs.SelfAwarenessTrainingDocumentDTO, error)
	CreateSelfAwarenessTrainingDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error
	GetSelfAwarenessTrainingDocumentPath(planID uint, token string) (string, error)
	// Self-Awareness Training - Operations
	MarkSelfAwarenessTrainingAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error

	// Business Context Operations
	CreateBusinessContext(planID uint, currentUserID uint) error
	GetBusinessContextByPlanID(planID uint) ([]*riskDTOs.BusinessContextDetailDTO, error)
	UpdateBusinessContext(id uint, currentUserID uint, formDTO *riskDTOs.UpdateBusinessContextFormDTO) error
	UpdateBusinessContextStatus(id uint, currentUserID uint, formDTO *riskDTOs.UpdateBusinessContextStatusFormDTO) error
	MarkBusinessContextAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error

	// Timeline Operations
	RecalculatePercentage(planID uint) error
}

type RiskRegisterUpdatePlanService struct {
	repo                riskRepositories.RiskRegisterUpdatePlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	userSvc             organizationServices.UserServiceInterface
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface
	timelineTypeSvc     RiskUpdatePlanTimelineTypeServiceInterface
}

func NewRiskRegisterUpdatePlanService(
	repo riskRepositories.RiskRegisterUpdatePlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
	timelineTypeSvc RiskUpdatePlanTimelineTypeServiceInterface,
) RiskRegisterUpdatePlanServiceInterface {
	return &RiskRegisterUpdatePlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		userSvc:             userSvc,
		organizationUnitSvc: organizationUnitSvc,
		timelineTypeSvc:     timelineTypeSvc,
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

	// Automatically create 14 timelines (one for each timeline type)
	timelineTypes, err := s.repo.GetAllTimelineTypes()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	now := time.Now()
	for _, timelineType := range timelineTypes {
		// Mark "risk_register_update_plan" type as completed since the plan is being created
		isCompleted := timelineType.Type == "risk_register_update_plan"

		timeline := &models.RiskUpdatePlanTimeline{
			RiskUpdatePlanID:             planID,
			RiskUpdatePlanTimelineTypeID: timelineType.ID,
			IsBusinessContextUpdated:     false,
			IsNotApplicable:              false,
			IsCompleted:                  isCompleted,
			CreatedBYID:                  currentUserID,
			CreatedAt:                    now,
		}

		if err := s.repo.CreateTimeline(timeline, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Calculate and update initial percentage (should be ~7.14% since 1 out of 14 is completed)
	if err := s.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return 0, err
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

		// Delete timelines first
		if err := s.repo.DeleteTimelinesByPlanID(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Delete team members
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

// ============================================
// ERM Department Evaluation - Auto Fetch All
// ============================================

func (s *RiskRegisterUpdatePlanService) CreateERMDepartmentEvaluation(planID uint, currentUserID uint) error {
	plan, err := s.GetByID(planID)
	if err != nil {
		return err
	}

	// Get timeline type for ERM Department Evaluation
	timelineType, err := s.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Fetch risks, processes, and controls for the business unit
	riskIDs, err := s.repo.GetRisksByBusinessUnitID(plan.BusinessUnitID)
	if err != nil {
		tx.Rollback()
		return err
	}

	processIDs, err := s.repo.GetProcessesByBusinessUnitID(plan.BusinessUnitID)
	if err != nil {
		tx.Rollback()
		return err
	}

	controlIDs, err := s.repo.GetControlsByBusinessUnitID(plan.BusinessUnitID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create risk mappings
	riskMappings := make([]*models.RiskUpdatePlanRisk, 0, len(riskIDs))
	for _, riskID := range riskIDs {
		riskMappings = append(riskMappings, &models.RiskUpdatePlanRisk{
			RiskUpdatePlanTimelineID: timeline.ID,
			RiskRegisterID:           riskID,
			CreatedBYID:              &currentUserID,
			CreatedAt:                time.Now(),
		})
	}

	if err := s.repo.CreateRiskMappings(riskMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Create process mappings
	processMappings := make([]*models.RiskUpdatePlanProcess, 0, len(processIDs))
	for _, processID := range processIDs {
		processMappings = append(processMappings, &models.RiskUpdatePlanProcess{
			RiskUpdatePlanTimelineID: timeline.ID,
			ProcessID:                processID,
			CreatedBYID:              &currentUserID,
			CreatedAt:                time.Now(),
		})
	}

	if err := s.repo.CreateProcessMappings(processMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Create control mappings
	controlMappings := make([]*models.RiskUpdatePlanControl, 0, len(controlIDs))
	for _, controlID := range controlIDs {
		controlMappings = append(controlMappings, &models.RiskUpdatePlanControl{
			RiskUpdatePlanTimelineID: timeline.ID,
			ControlID:                controlID,
			CreatedBYID:              &currentUserID,
			CreatedAt:                time.Now(),
		})
	}

	if err := s.repo.CreateControlMappings(controlMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Mark timeline as completed
	if err := s.repo.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := s.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// ============================================
// ERM Department Evaluation - Risks
// ============================================

func (s *RiskRegisterUpdatePlanService) MapRisksToTimeline(planID uint, currentUserID uint, formDTO *riskDTOs.MapRisksToUpdatePlanDTO) error {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := s.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete all existing risk mappings for this timeline
	if err := s.repo.DeleteAllRiskMappingsByTimelineID(timeline.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Create new risk mappings
	riskMappings := make([]*models.RiskUpdatePlanRisk, 0, len(formDTO.RiskRegisterIDs))
	for _, riskID := range formDTO.RiskRegisterIDs {
		riskMappings = append(riskMappings, &models.RiskUpdatePlanRisk{
			RiskUpdatePlanTimelineID: timeline.ID,
			RiskRegisterID:           riskID,
			CreatedBYID:              &currentUserID,
			CreatedAt:                time.Now(),
		})
	}

	if err := s.repo.CreateRiskMappings(riskMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *RiskRegisterUpdatePlanService) GetERMDeptEvaluationRisks(planID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationListResponseDTO, error) {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := s.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetERMDeptEvaluationRisks(timeline.ID, params)
}

// ============================================
// ERM Department Evaluation - Processes
// ============================================

func (s *RiskRegisterUpdatePlanService) MapProcessesToTimeline(planID uint, currentUserID uint, formDTO *riskDTOs.MapProcessesToUpdatePlanDTO) error {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := s.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete all existing process mappings for this timeline
	if err := s.repo.DeleteAllProcessMappingsByTimelineID(timeline.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Create new process mappings
	processMappings := make([]*models.RiskUpdatePlanProcess, 0, len(formDTO.ProcessIDs))
	for _, processID := range formDTO.ProcessIDs {
		processMappings = append(processMappings, &models.RiskUpdatePlanProcess{
			RiskUpdatePlanTimelineID: timeline.ID,
			ProcessID:                processID,
			CreatedBYID:              &currentUserID,
			CreatedAt:                time.Now(),
		})
	}

	if err := s.repo.CreateProcessMappings(processMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *RiskRegisterUpdatePlanService) GetERMDeptEvaluationProcesses(planID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationProcessListResponseDTO, error) {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := s.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetERMDeptEvaluationProcesses(timeline.ID, params)
}

// ============================================
// ERM Department Evaluation - Controls
// ============================================

func (s *RiskRegisterUpdatePlanService) MapControlsToTimeline(planID uint, currentUserID uint, formDTO *riskDTOs.MapControlsToUpdatePlanDTO) error {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := s.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete all existing control mappings for this timeline
	if err := s.repo.DeleteAllControlMappingsByTimelineID(timeline.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Create new control mappings
	controlMappings := make([]*models.RiskUpdatePlanControl, 0, len(formDTO.ControlIDs))
	for _, controlID := range formDTO.ControlIDs {
		controlMappings = append(controlMappings, &models.RiskUpdatePlanControl{
			RiskUpdatePlanTimelineID: timeline.ID,
			ControlID:                controlID,
			CreatedBYID:              &currentUserID,
			CreatedAt:                time.Now(),
		})
	}

	if err := s.repo.CreateControlMappings(controlMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *RiskRegisterUpdatePlanService) GetERMDeptEvaluationControls(planID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationControlListResponseDTO, error) {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := s.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetERMDeptEvaluationControls(timeline.ID, params)
}

// ============================================
// ERM Department Evaluation Operations
// ============================================

func (s *RiskRegisterUpdatePlanService) MarkERMDeptEvaluationAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := s.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	// Update timeline - mark as not applicable and completed
	now := time.Now()
	timeline.IsNotApplicable = true
	timeline.IsCompleted = true
	timeline.Justification = &formDTO.Justification
	timeline.UpdatedBYID = &currentUserID
	timeline.UpdatedAt = &now

	if err := tx.DB().Save(&timeline).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := s.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// ERM Department Evaluation Document Operations

func (s *RiskRegisterUpdatePlanService) GetERMDeptEvaluationDocuments(planID uint) ([]*riskDTOs.ERMDeptEvaluationDocumentDTO, error) {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := s.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	// Get documents
	documents, err := s.repo.GetERMDeptDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	var documentDTOs []*riskDTOs.ERMDeptEvaluationDocumentDTO
	for _, doc := range documents {
		dto := &riskDTOs.ERMDeptEvaluationDocumentDTO{
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

func (s *RiskRegisterUpdatePlanService) CreateERMDeptEvaluationDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error {

	// Verify plan exists
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for ERM Department Evaluation
	timelineType, err := s.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	// Handle documents
	var newDocuments []*models.RiskUpdatePlanErmDeptDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskUpdatePlanErmDeptDocument

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
				"erm-dept-evaluation-documents",
			)

			document := &models.RiskUpdatePlanErmDeptDocument{
				RiskUpdatePlanTimelineID: timeline.ID,
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
			existingDocument, err := s.repo.GetERMDeptDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Verify document belongs to this timeline
			if existingDocument.RiskUpdatePlanTimelineID != timeline.ID {
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
			_, err := s.repo.CreateERMDeptDocument(doc, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Delete documents
	if len(deleteDocumentIDs) > 0 {
		for _, docID := range deleteDocumentIDs {
			if err := s.repo.DeleteERMDeptDocument(*docID, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// File system operations - store new files
	for _, doc := range newDocuments {
		storagePath := filepath.Join(
			"risk",
			"risk-register-update-plan",
			strconv.FormatUint(uint64(planID), 10),
			"erm-dept-evaluation-documents",
		)

		if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	// File system operations - delete files
	for _, doc := range documentsToDelete {
		storagePath := filepath.Join(
			"risk",
			"risk-register-update-plan",
			strconv.FormatUint(uint64(planID), 10),
			"erm-dept-evaluation-documents",
		)

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
			models.RiskUpdatePlanErmDeptDocument{}.TableName(),
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

func (s *RiskRegisterUpdatePlanService) GetERMDeptEvaluationDocumentPath(planID uint, token string) (string, error) {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := s.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return "", err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return "", err
	}

	// Get all documents for this timeline
	documents, err := s.repo.GetERMDeptDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return "", err
	}

	// Find document with matching token
	for _, doc := range documents {
		if doc.Token != nil && *doc.Token == token {
			// Decrypt filename
			decryptedFilename, err := utility.DecryptData(*doc.Token)
			if err != nil {
				return "", err
			}

			// Construct file path
			storagePath := filepath.Join(
				"risk",
				"risk-register-update-plan",
				strconv.FormatUint(uint64(planID), 10),
				"erm-dept-evaluation-documents",
			)
			baseStoragePath := config.GetEnv("STORAGE_PATH", "")
			fullPath := filepath.Join(baseStoragePath, storagePath, string(decryptedFilename))

			// Check if the file exists
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				return "", utility.NewResourceNotFoundError(err, 0)
			}

			return fullPath, nil
		}
	}

	return "", utility.NewResourceNotFoundError(errors.New("document not found"), 0)
}

// Introductory Meeting Operations

func (s *RiskRegisterUpdatePlanService) MarkIntroductoryMeetingAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	// Get timeline type for Introductory Meeting
	timelineType, err := s.timelineTypeSvc.GetByType("introductory_meeting")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	// Update timeline - mark as not applicable and completed
	now := time.Now()
	timeline.IsNotApplicable = true
	timeline.IsCompleted = true
	timeline.Justification = &formDTO.Justification
	timeline.UpdatedBYID = &currentUserID
	timeline.UpdatedAt = &now

	if err := tx.DB().Save(&timeline).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := s.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Introductory Meeting Document Operations

func (s *RiskRegisterUpdatePlanService) GetIntroductoryMeetingDocumentPath(planID uint, token string) (string, error) {
	// Get timeline type for Introductory Meeting
	timelineType, err := s.timelineTypeSvc.GetByType("introductory_meeting")
	if err != nil {
		return "", err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return "", err
	}

	// Get all documents for this timeline
	documents, err := s.repo.GetIntroductoryMeetingDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return "", err
	}

	// Find document with matching token
	for _, doc := range documents {
		if doc.Token != nil && *doc.Token == token {
			// Decrypt filename
			decryptedFilename, err := utility.DecryptData(*doc.Token)
			if err != nil {
				return "", err
			}

			// Construct file path
			storagePath := filepath.Join("risk", "risk-register-update-plan", strconv.FormatUint(uint64(planID), 10), "introductory-meeting-documents")
			baseStoragePath := config.GetEnv("STORAGE_PATH", "")
			fullPath := filepath.Join(baseStoragePath, storagePath, string(decryptedFilename))

			// Check if the file exists
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				return "", utility.NewResourceNotFoundError(err, 0)
			}

			return fullPath, nil
		}
	}

	return "", utility.NewResourceNotFoundError(errors.New("document not found"), 0)
}

// Introductory Meeting Document Operations

func (s *RiskRegisterUpdatePlanService) GetIntroductoryMeetingDocuments(planID uint) ([]*riskDTOs.IntroductoryMeetingDocumentDTO, error) {
	// Get timeline type for Introductory Meeting
	timelineType, err := s.timelineTypeSvc.GetByType("introductory_meeting")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	// Get documents by timeline ID
	documents, err := s.repo.GetIntroductoryMeetingDocumentsByTimelineID(timeline.ID)
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

	// Get timeline type for Introductory Meeting
	timelineType, err := s.timelineTypeSvc.GetByType("introductory_meeting")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
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
				RiskUpdatePlanTimelineID: timeline.ID,
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

			// Verify document belongs to this timeline
			if existingDocument.RiskUpdatePlanTimelineID != timeline.ID {
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

	// Mark timeline as completed
	if err := s.repo.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := s.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
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

// Self-Awareness Training Document Operations

func (s *RiskRegisterUpdatePlanService) GetSelfAwarenessTrainingDocuments(planID uint) ([]*riskDTOs.SelfAwarenessTrainingDocumentDTO, error) {
	// Verify plan exists
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, planID)
		}
		return nil, err
	}

	// Get timeline type for Self-Awareness Training
	timelineType, err := s.timelineTypeSvc.GetByType("self_awareness_training")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	// Get documents by timeline ID
	documents, err := s.repo.GetSelfAwarenessTrainingDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return nil, err
	}

	var dtos []*riskDTOs.SelfAwarenessTrainingDocumentDTO
	for _, doc := range documents {
		dto := &riskDTOs.SelfAwarenessTrainingDocumentDTO{
			ID:           &doc.ID,
			URL:          doc.URL,
			Name:         doc.Name,
			Token:        doc.Token,
			OriginalName: doc.OriginalName,
			Type:         doc.Type,
			Size:         doc.Size,
		}
		dtos = append(dtos, dto)
	}

	return dtos, nil
}

func (s *RiskRegisterUpdatePlanService) CreateSelfAwarenessTrainingDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error {

	// Verify plan exists
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for Self-Awareness Training
	timelineType, err := s.timelineTypeSvc.GetByType("self_awareness_training")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	// Handle documents
	var newDocuments []*models.RiskUpdatePlanSelfAwarenessTraining
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskUpdatePlanSelfAwarenessTraining

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
				"self-awareness-training-documents",
			)

			document := &models.RiskUpdatePlanSelfAwarenessTraining{
				RiskUpdatePlanTimelineID: timeline.ID,
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
			existingDocument, err := s.repo.GetSelfAwarenessTrainingDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Verify document belongs to this timeline
			if existingDocument.RiskUpdatePlanTimelineID != timeline.ID {
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
			_, err := s.repo.CreateSelfAwarenessTrainingDocument(doc, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Delete documents
	if len(deleteDocumentIDs) > 0 {
		for _, docID := range deleteDocumentIDs {
			if err := s.repo.DeleteSelfAwarenessTrainingDocument(*docID, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// File system operations - store new files
	for _, doc := range newDocuments {
		storagePath := filepath.Join(
			"risk",
			"risk-register-update-plan",
			strconv.FormatUint(uint64(planID), 10),
			"self-awareness-training-documents",
		)

		if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	// File system operations - delete files
	for _, doc := range documentsToDelete {
		storagePath := filepath.Join(
			"risk",
			"risk-register-update-plan",
			strconv.FormatUint(uint64(planID), 10),
			"self-awareness-training-documents",
		)

		if doc.Name != nil {
			if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Mark timeline as completed
	if err := s.repo.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := s.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity
	go func() {
		err := s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskUpdatePlanSelfAwarenessTraining{}.TableName(),
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

func (s *RiskRegisterUpdatePlanService) GetSelfAwarenessTrainingDocumentPath(planID uint, token string) (string, error) {

	// Get timeline type for Self-Awareness Training
	timelineType, err := s.timelineTypeSvc.GetByType("self_awareness_training")
	if err != nil {
		return "", err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return "", err
	}

	// Get all documents for this timeline
	documents, err := s.repo.GetSelfAwarenessTrainingDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return "", err
	}

	// Find document with matching token
	for _, doc := range documents {
		if doc.Token != nil && *doc.Token == token {

			// Decrypt filename
			decryptedFilename, err := utility.DecryptData(*doc.Token)
			if err != nil {
				return "", err
			}

			// Construct file path
			storagePath := filepath.Join(
				"risk",
				"risk-register-update-plan",
				strconv.FormatUint(uint64(planID), 10),
				"self-awareness-training-documents",
			)

			baseStoragePath := config.GetEnv("STORAGE_PATH", "")
			fullPath := filepath.Join(baseStoragePath, storagePath, string(decryptedFilename))

			// Check if file exists
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				return "", utility.NewResourceNotFoundError(err, 0)
			}

			return fullPath, nil
		}
	}

	return "", utility.NewResourceNotFoundError(errors.New("document not found"), 0)
}

func (s *RiskRegisterUpdatePlanService) MarkSelfAwarenessTrainingAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	// Get timeline type for Self-Awareness Training
	timelineType, err := s.timelineTypeSvc.GetByType("self_awareness_training")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	// Update timeline - mark as not applicable and completed
	now := time.Now()
	timeline.IsNotApplicable = true
	timeline.IsCompleted = true
	timeline.Justification = &formDTO.Justification
	timeline.UpdatedBYID = &currentUserID
	timeline.UpdatedAt = &now

	if err := tx.DB().Save(&timeline).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := s.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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

	// Get timeline type for Business Context
	timelineType, err := s.timelineTypeSvc.GetByType("business_context")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	// Get business unit ID from plan
	businessUnitID, err := s.repo.GetBusinessUnitIDByPlanID(planID)
	if err != nil {
		return err
	}

	// Get OrganizationUnitBusinessContext data
	orgBusinessContext, err := s.repo.GetOrganizationUnitBusinessContextByBusinessUnit(businessUnitID)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	// Create a single business context record with data from OrganizationUnitBusinessContext
	businessContext := &models.RiskUpdatePlanBusinessContext{
		RiskUpdatePlanTimelineID:   timeline.ID,
		InternalAuditNote:          orgBusinessContext.InternalAuditNote,
		ExternalContext:            orgBusinessContext.ExternalContext,
		AnnualObjectiveDescription: orgBusinessContext.AnnualObjectiveDescription,
		CreatedBYID:                &currentUserID,
		CreatedAt:                  time.Now(),
	}

	if err := s.repo.CreateBusinessContext(businessContext, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Mark timeline as completed
	if err := s.repo.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := s.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
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
				"timeline_id":      timeline.ID,
			},
			currentUserID,
			"risk", "risk_update_plan_business_context", "create", nil,
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

func (s *RiskRegisterUpdatePlanService) MarkBusinessContextAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	// Get timeline type for Business Context
	timelineType, err := s.timelineTypeSvc.GetByType("business_context")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	// Update timeline - mark as not applicable and completed
	now := time.Now()
	timeline.IsNotApplicable = true
	timeline.IsCompleted = true
	timeline.Justification = &formDTO.Justification
	timeline.UpdatedBYID = &currentUserID
	timeline.UpdatedAt = &now

	if err := tx.DB().Save(&timeline).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := s.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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

	// Get timeline type for Business Context
	timelineType, err := s.timelineTypeSvc.GetByType("business_context")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	// Get business contexts by timeline ID
	businessContexts, err := s.repo.GetBusinessContextByTimelineID(timeline.ID)
	if err != nil {
		return nil, err
	}

	var dtos []*riskDTOs.BusinessContextDetailDTO

	err = mapstructure.Decode(businessContexts, &dtos)
	if err != nil {
		return nil, err
	}

	return dtos, nil
}

func (s *RiskRegisterUpdatePlanService) UpdateBusinessContext(id uint, currentUserID uint, formDTO *riskDTOs.UpdateBusinessContextFormDTO) error {
	// Get existing business context
	businessContext, err := s.repo.GetBusinessContextByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	// Update fields from DTO
	if formDTO.InternalAuditNote != nil {
		businessContext.InternalAuditNote = formDTO.InternalAuditNote
	}
	if formDTO.ExternalContext != nil {
		businessContext.ExternalContext = formDTO.ExternalContext
	}
	if formDTO.AnnualObjectiveDescription != nil {
		businessContext.AnnualObjectiveDescription = formDTO.AnnualObjectiveDescription
	}

	tx := s.repo.BeginTransaction()

	// Update business context
	if err := s.repo.UpdateBusinessContext(businessContext, tx); err != nil {
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
				"internal_audit_note":          formDTO.InternalAuditNote,
				"external_context":             formDTO.ExternalContext,
				"annual_objective_description": formDTO.AnnualObjectiveDescription,
			},
			currentUserID,
			"risk", "risk_update_plan_business_context", "update", nil,
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

func (s *RiskRegisterUpdatePlanService) UpdateBusinessContextStatus(id uint, currentUserID uint, formDTO *riskDTOs.UpdateBusinessContextStatusFormDTO) error {
	// Verify plan exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	// Get timeline type for Business Context
	timelineType, err := s.timelineTypeSvc.GetByType("business_context")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(id, timelineType.ID)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	// Update status
	if err := s.repo.UpdateBusinessContextStatus(timeline.ID, *formDTO.IsBusinessContextUpdated, tx); err != nil {
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
				"is_business_context_updated": *formDTO.IsBusinessContextUpdated,
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

// Timeline Operations

func (s *RiskRegisterUpdatePlanService) RecalculatePercentage(planID uint) error {
	// Verify plan exists
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	tx := s.repo.BeginTransaction()

	if err := s.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
