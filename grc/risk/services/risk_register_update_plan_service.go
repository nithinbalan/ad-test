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
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
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

	// Leadership Assessment Operations
	GetLeadershipAssessmentUsers(planID uint, params *riskDTOs.UserListQueryParamsDTO) (*riskDTOs.LeaderShipAssessmentUserListResponseDTO, error)
	CreateLeadershipAssessment(planID uint, currentUserID uint) error
	MarkLeadershipAssessmentAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error
	GetLeadershipAssessmentDocuments(planID uint) ([]*riskDTOs.LeadershipAssessmentDocumentDTO, error)
	CreateLeadershipAssessmentDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error
	GetLeadershipAssessmentDocumentPath(planID uint, token string) (string, error)

	// Risk Team Assessment Operations
	GetRiskTeamAssessmentUser(planID uint) (*riskDTOs.RiskTeamAssessmentUserListDTO, error)
	CreateRiskTeamAssessment(planID uint, currentUserID uint) error
	MarkRiskTeamAssessmentAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error
	GetRiskTeamAssessmentDocuments(planID uint) ([]*riskDTOs.RiskTeamAssessmentDocumentDTO, error)
	CreateRiskTeamAssessmentDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error
	GetRiskTeamAssessmentDocumentPath(planID uint, token string) (string, error)

	//RiskIdentification
	GetAllRiskIdentifications(planID uint, params *riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO) (*riskDTOs.RiskIdentificationListResponseDTO, error)
	CreateRiskIdentification(planID uint, currentUserID uint, riskRegisterFormDTO *riskDTOs.UpdatePlanRiskRegisterFormDTO) (uint, error)
	UpdateRiskIdentification(planID uint, id uint, currentUserID uint, risk *riskDTOs.UpdatePlanRiskRegisterFormDTO) error
	DeleteRiskIdentification(planID uint, ids []uint, currentUserID uint) ([]uint, error)
	MarkRiskIdentificationAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error

	//Control Register Update
	GetControlCountsByImplementationStatus(planID uint) (*riskDTOs.ControlRegisterUpdatedDTO, error)
	CreateControlTest(planID uint, currentUserID uint, controlTestFormDTO *riskDTOs.UpdatePlanControlTestFormDTO) (uint, error)
	MarkControlRegisterUpdateAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error

	// Risk Assessment Workshop Document Operations
	GetRiskAssessmentWorkshopDocuments(planID uint) ([]*riskDTOs.RiskAssessmentWorkshopDocumentDTO, error)
	CreateRiskAssessmentWorkshopDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error
	GetRiskAssessmentWorkshopDocumentPath(planID uint, token string) (string, error)
	// Risk Assessment Workshop - Operations
	MarkRiskAssessmentWorkshopAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error

	// Risk Workshop Assessment Document Operations
	GetRiskWorkshopAssessmentDocuments(planID uint) ([]*riskDTOs.RiskWorkshopAssessmentDocumentDTO, error)
	CreateRiskWorkshopAssessmentDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error
	GetRiskWorkshopAssessmentDocumentPath(planID uint, token string) (string, error)
	// Risk Workshop Assessment - Operations
	MarkRiskWorkshopAssessmentAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error

	// Timeline Operations
	RecalculatePercentage(planID uint) error
}

type RiskRegisterUpdatePlanService struct {
	repo                 riskRepositories.RiskRegisterUpdatePlanRepositoryInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	userSvc              organizationServices.UserServiceInterface
	organizationUnitSvc  settingsServices.OrganizationUnitServiceInterface
	timelineTypeSvc      RiskUpdatePlanTimelineTypeServiceInterface
	riskRegisterSvc      RiskRegisterServiceInterface
	riskStatusSvc        RiskStatusServiceInterface
	riskMatrixSvc        RiskMatrixServiceInterface
	organizationUnitsSvc organizationServices.OrganizationUnitsServiceInterface
	notificationSvc      commonServices.NotificationServiceInterface
	riskSettingsSvc      RiskSettingsServiceInterface
	riskCategorySvc      RiskCategoryServiceInterface
	riskSourceSvc        RiskSourceServiceInterface
	controlTestSvc       organizationServices.ControlTestServiceInterface
}

func NewRiskRegisterUpdatePlanService(
	repo riskRepositories.RiskRegisterUpdatePlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
	timelineTypeSvc RiskUpdatePlanTimelineTypeServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	riskStatusSvc RiskStatusServiceInterface,
	riskMatrixSvc RiskMatrixServiceInterface,
	organizationUnitsSvc organizationServices.OrganizationUnitsServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	riskSettingsSvc RiskSettingsServiceInterface,
	riskCategorySvc RiskCategoryServiceInterface,
	riskSourceSvc RiskSourceServiceInterface,
	controlTestSvc organizationServices.ControlTestServiceInterface,
) RiskRegisterUpdatePlanServiceInterface {
	return &RiskRegisterUpdatePlanService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		userSvc:              userSvc,
		organizationUnitSvc:  organizationUnitSvc,
		timelineTypeSvc:      timelineTypeSvc,
		riskRegisterSvc:      riskRegisterSvc,
		riskStatusSvc:        riskStatusSvc,
		riskMatrixSvc:        riskMatrixSvc,
		organizationUnitsSvc: organizationUnitsSvc,
		notificationSvc:      notificationSvc,
		riskSettingsSvc:      riskSettingsSvc,
		riskCategorySvc:      riskCategorySvc,
		riskSourceSvc:        riskSourceSvc,
		controlTestSvc:       controlTestSvc,
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

// Leadership Assessment Operations

func (s *RiskRegisterUpdatePlanService) GetLeadershipAssessmentUsers(planID uint, params *riskDTOs.UserListQueryParamsDTO) (*riskDTOs.LeaderShipAssessmentUserListResponseDTO, error) {
	// Get timeline type for Leadership Assessment
	timelineType, err := s.timelineTypeSvc.GetByType("leadership_assessment")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetLeadershipAssessmentUsers(timeline.ID, params)
}

func (s *RiskRegisterUpdatePlanService) CreateLeadershipAssessment(planID uint, currentUserID uint) error {
	plan, err := s.GetByID(planID)
	if err != nil {
		return err
	}

	// Get timeline type for Leadership Assessment
	timelineType, err := s.timelineTypeSvc.GetByType("leadership_assessment")
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

	// Fetch users for the business unit
	users, err := s.userSvc.GetUsersByOrganizationUnit(plan.BusinessUnitID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create leadership assessment mappings
	assessments := make([]*models.RiskUpdatePlanLeadershipAssessment, 0, len(users))
	for _, user := range users {
		assessments = append(assessments, &models.RiskUpdatePlanLeadershipAssessment{
			RiskUpdatePlanTimelineID: timeline.ID,
			UserID:                   user.ID,
			CreatedBYID:              &currentUserID,
			CreatedAt:                time.Now(),
		})
	}

	if err := s.repo.CreateLeadershipAssessments(assessments, tx); err != nil {
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

func (s *RiskRegisterUpdatePlanService) MarkLeadershipAssessmentAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := s.timelineTypeSvc.GetByType("leadership_assessment")
	if err != nil {
		return err
	}

	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	if err := s.repo.DeleteLeadershipAssessmentsByTimelineID(timeline.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

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

	if err := s.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *RiskRegisterUpdatePlanService) GetLeadershipAssessmentDocuments(planID uint) ([]*riskDTOs.LeadershipAssessmentDocumentDTO, error) {
	timelineType, err := s.timelineTypeSvc.GetByType("leadership_assessment")
	if err != nil {
		return nil, err
	}

	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	documents, err := s.repo.GetLeadershipDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return nil, err
	}

	var documentDTOs []*riskDTOs.LeadershipAssessmentDocumentDTO
	for _, doc := range documents {
		dto := &riskDTOs.LeadershipAssessmentDocumentDTO{
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

func (s *RiskRegisterUpdatePlanService) CreateLeadershipAssessmentDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error {
	// Verify plan exists
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for Leadership Assessment
	timelineType, err := s.timelineTypeSvc.GetByType("leadership_assessment")
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
	var newDocuments []*models.RiskUpdatePlanLeadershipDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskUpdatePlanLeadershipDocument

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
				"leadership-assessment-documents",
			)

			document := &models.RiskUpdatePlanLeadershipDocument{
				RiskUpdatePlanLeadershipAssessmentID: timeline.ID,
				URL:                                  &url,
				Name:                                 &doc.FileName,
				Size:                                 &doc.Size,
				Token:                                &encryptedFilename,
				OriginalName:                         &doc.OriginalName,
				Type:                                 &doc.Type,
				CreatedBYID:                          &currentUserID,
				CreatedAt:                            now,
			}

			newDocuments = append(newDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {

			// DELETE FILE
			existingDocument, err := s.repo.GetLeadershipDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Verify document belongs to this timeline
			if existingDocument.RiskUpdatePlanLeadershipAssessmentID != timeline.ID {
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
			_, err := s.repo.CreateLeadershipDocument(doc, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Delete documents
	if len(deleteDocumentIDs) > 0 {
		for _, docID := range deleteDocumentIDs {
			if err := s.repo.DeleteLeadershipDocument(*docID, tx); err != nil {
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
			"leadership-assessment-documents",
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
			"leadership-assessment-documents",
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
			models.RiskUpdatePlanLeadershipDocument{}.TableName(),
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

func (s *RiskRegisterUpdatePlanService) GetLeadershipAssessmentDocumentPath(planID uint, token string) (string, error) {
	timelineType, err := s.timelineTypeSvc.GetByType("leadership_assessment")
	if err != nil {
		return "", err
	}

	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return "", err
	}

	documents, err := s.repo.GetLeadershipDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return "", err
	}

	for _, doc := range documents {
		if doc.Token != nil && *doc.Token == token {
			decryptedFilename, err := utility.DecryptData(*doc.Token)
			if err != nil {
				return "", err
			}

			storagePath := filepath.Join(
				"risk",
				"risk-register-update-plan",
				strconv.FormatUint(uint64(planID), 10),
				"leadership-assessment-documents",
			)
			baseStoragePath := config.GetEnv("STORAGE_PATH", "")
			fullPath := filepath.Join(baseStoragePath, storagePath, string(decryptedFilename))

			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				return "", utility.NewResourceNotFoundError(err, 0)
			}

			return fullPath, nil
		}
	}

	return "", utility.NewResourceNotFoundError(fmt.Errorf("document not found"), 0)
}

// Risk Team Assessment Operations

func (s *RiskRegisterUpdatePlanService) GetRiskTeamAssessmentUser(planID uint) (*riskDTOs.RiskTeamAssessmentUserListDTO, error) {

	timelineType, err := s.timelineTypeSvc.GetByType("risk_team_assessment")
	if err != nil {
		return nil, err
	}

	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	assessments, err := s.repo.GetRiskTeamAssessmentUsers(timeline.ID)
	if err != nil {
		return nil, err
	}

	if len(assessments) == 0 {
		return nil, utility.NewResourceNotFoundError(fmt.Errorf("no risk team assessment found"), planID)
	}

	var dto riskDTOs.RiskTeamAssessmentUserListDTO

	err = mapstructure.Decode(assessments[0], &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}

func (s *RiskRegisterUpdatePlanService) CreateRiskTeamAssessment(planID uint, currentUserID uint) error {
	plan, err := s.GetByID(planID)
	if err != nil {
		return err
	}

	timelineType, err := s.timelineTypeSvc.GetByType("risk_team_assessment")
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	assessments := make([]*models.RiskUpdatePlanRiskTeamAssessment, 0, 1)
	assessments = append(assessments, &models.RiskUpdatePlanRiskTeamAssessment{
		RiskUpdatePlanTimelineID: timeline.ID,
		RiskChampionID:           plan.RiskChampionID,
		CreatedBYID:              &currentUserID,
		CreatedAt:                time.Now(),
	})

	if err := s.repo.CreateRiskTeamAssessments(assessments, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *RiskRegisterUpdatePlanService) MarkRiskTeamAssessmentAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := s.timelineTypeSvc.GetByType("risk_team_assessment")
	if err != nil {
		return err
	}

	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	if err := s.repo.DeleteRiskTeamAssessmentsByTimelineID(timeline.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

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

	if err := s.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *RiskRegisterUpdatePlanService) GetRiskTeamAssessmentDocuments(planID uint) ([]*riskDTOs.RiskTeamAssessmentDocumentDTO, error) {
	timelineType, err := s.timelineTypeSvc.GetByType("risk_team_assessment")
	if err != nil {
		return nil, err
	}

	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	documents, err := s.repo.GetRiskTeamDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return nil, err
	}

	var documentDTOs []*riskDTOs.RiskTeamAssessmentDocumentDTO
	for _, doc := range documents {
		dto := &riskDTOs.RiskTeamAssessmentDocumentDTO{
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

func (s *RiskRegisterUpdatePlanService) CreateRiskTeamAssessmentDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error {
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := s.timelineTypeSvc.GetByType("risk_team_assessment")
	if err != nil {
		return err
	}

	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	var newDocuments []*models.RiskUpdatePlanRiskTeamDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskUpdatePlanRiskTeamDocument

	now := time.Now()

	for _, doc := range documents {

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
				"risk-team-assessment-documents",
			)

			document := &models.RiskUpdatePlanRiskTeamDocument{
				RiskUpdatePlanRiskTeamAssessmentID: timeline.ID,
				URL:                                &url,
				Name:                               &doc.FileName,
				Size:                               &doc.Size,
				Token:                              &encryptedFilename,
				OriginalName:                       &doc.OriginalName,
				Type:                               &doc.Type,
				CreatedBYID:                        &currentUserID,
				CreatedAt:                          now,
			}

			newDocuments = append(newDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {

			existingDocument, err := s.repo.GetRiskTeamDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				return err
			}

			if existingDocument.RiskUpdatePlanRiskTeamAssessmentID != timeline.ID {
				tx.Rollback()
				return utility.NewResourceNotFoundError(fmt.Errorf("document not found"), *doc.ID)
			}

			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			documentsToDelete = append(documentsToDelete, existingDocument)
		}
	}

	if len(newDocuments) > 0 {
		for _, doc := range newDocuments {
			_, err := s.repo.CreateRiskTeamDocument(doc, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if len(deleteDocumentIDs) > 0 {
		for _, docID := range deleteDocumentIDs {
			if err := s.repo.DeleteRiskTeamDocument(*docID, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join(
			"risk",
			"risk-register-update-plan",
			strconv.FormatUint(uint64(planID), 10),
			"risk-team-assessment-documents",
		)

		if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join(
			"risk",
			"risk-register-update-plan",
			strconv.FormatUint(uint64(planID), 10),
			"risk-team-assessment-documents",
		)

		if doc.Name != nil {
			if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	return nil
}

func (s *RiskRegisterUpdatePlanService) GetRiskTeamAssessmentDocumentPath(planID uint, token string) (string, error) {
	timelineType, err := s.timelineTypeSvc.GetByType("risk_team_assessment")
	if err != nil {
		return "", err
	}

	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return "", err
	}

	documents, err := s.repo.GetRiskTeamDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return "", err
	}

	for _, doc := range documents {
		if doc.Token != nil && *doc.Token == token {
			decryptedFilename, err := utility.DecryptData(*doc.Token)
			if err != nil {
				return "", err
			}

			storagePath := filepath.Join(
				"risk",
				"risk-register-update-plan",
				strconv.FormatUint(uint64(planID), 10),
				"risk-team-assessment-documents",
			)
			baseStoragePath := config.GetEnv("STORAGE_PATH", "")
			fullPath := filepath.Join(baseStoragePath, storagePath, string(decryptedFilename))

			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				return "", utility.NewResourceNotFoundError(err, 0)
			}

			return fullPath, nil
		}
	}

	return "", utility.NewResourceNotFoundError(fmt.Errorf("document not found"), 0)
}

// RiskIdentification

func (rrups *RiskRegisterUpdatePlanService) GetAllRiskIdentifications(planID uint, params *riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO) (*riskDTOs.RiskIdentificationListResponseDTO, error) {
	return rrups.repo.GetAllRiskIdentifications(planID, params)
}
func (rrups *RiskRegisterUpdatePlanService) CreateRiskIdentification(planID uint, currentUserID uint, riskRegisterFormDTO *riskDTOs.UpdatePlanRiskRegisterFormDTO) (uint, error) {
	// Verify plan exists
	_, err := rrups.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, utility.NewResourceNotFoundError(err, planID)
		}
		return 0, err
	}

	// Get timeline type for risk_identification
	timelineType, err := rrups.timelineTypeSvc.GetByType("risk_identification")
	if err != nil {
		return 0, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rrups.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return 0, err
	}

	tx := rrups.repo.BeginTransaction()

	if err := rrups.validateRiskRegisterFormDTO(riskRegisterFormDTO); err != nil {
		return 0, err
	}

	referenceNumber, err := rrups.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	riskMatrix, err := rrups.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return 0, err
	}

	riskStatus, err := rrups.riskStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}

	registrationDate := time.Now()

	riskRegister := &models.RiskRegister{
		ReferenceNumber:          referenceNumber,
		RiskCategoryID:           riskRegisterFormDTO.RiskCategoryID,
		RiskSourceID:             riskRegisterFormDTO.RiskSourceID,
		RiskMatrixID:             riskMatrix.ID,
		ResponsibleUnitID:        riskRegisterFormDTO.ResponsibleUnitID,
		ResponsibleUserID:        riskRegisterFormDTO.ResponsibleUserID,
		RiskStatusID:             riskStatus.ID,
		Title:                    riskRegisterFormDTO.Title,
		Description:              riskRegisterFormDTO.Description,
		ReviewFrequencyID:        riskRegisterFormDTO.ReviewFrequencyID,
		LastReviewNote:           riskRegisterFormDTO.LastReviewNote,
		RiskCode:                 riskRegisterFormDTO.RiskCode,
		RegistrationDate:         &registrationDate,
		DivisionID:               riskRegisterFormDTO.DivisionID,
		SectionID:                riskRegisterFormDTO.SectionID,
		IsStrategic:              riskRegisterFormDTO.IsStrategic,
		IsOperational:            riskRegisterFormDTO.IsOperational,
		ExpectedCause:            riskRegisterFormDTO.ExpectedCause,
		Consequence:              riskRegisterFormDTO.Consequence,
		IsInternal:               riskRegisterFormDTO.IsInternal,
		IsExternal:               riskRegisterFormDTO.IsExternal,
		RiskRegisterUpdatePlanID: &planID,
		CreatedBYID:              currentUserID,
	}

	id, err := rrups.repo.CreateRiskIdentification(riskRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Check if this is the first risk for this plan
	riskCount, err := rrups.repo.CountRisksByPlanID(planID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Update timeline completion and percentage only if this is the first risk
	if riskCount == 1 {
		if err := rrups.repo.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		if err := rrups.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		existingRiskRegister, err := rrups.riskRegisterSvc.GetFullDetailsByID(riskRegister.ID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// log
		err = rrups.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskRegister{}.TableName(),
			id,
			existingRiskRegister,
			currentUserID,
			"risk", "risk_register", "risk_register_update_plan", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// Prepare common variables for notifications
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "risk-register"

		createdBy, err := rrups.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch created by user", map[string]interface{}{
				"error": err.Error(),
			})
		}

		// Send notification to responsible user and/or responsible unit head
		action := "responsible-user-assign"
		var notificationUserIDs []uint
		var emailTo []string

		// If responsible user is assigned, add them to notification
		if riskRegisterFormDTO.ResponsibleUserID != nil {
			notificationUserIDs = append(notificationUserIDs, *riskRegisterFormDTO.ResponsibleUserID)
			responsibleUser, err := rrups.userSvc.GetByID(*riskRegisterFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
			} else {
				emailTo = append(emailTo, responsibleUser.Email)
			}
		}

		// Also send notification to responsible unit head
		if riskRegisterFormDTO.ResponsibleUnitID != nil {
			headUsers, err := rrups.organizationUnitsSvc.GetHeadUsersByOrganizationUnitID(*riskRegisterFormDTO.ResponsibleUnitID)
			if err != nil {
				utility.LogError("Failed to get responsible unit head users for notification", map[string]interface{}{
					"unit_id": *riskRegisterFormDTO.ResponsibleUnitID,
					"error":   err.Error(),
				})
			} else {
				for _, headUser := range headUsers {
					// Avoid duplicate notifications if the head user is also the responsible user
					alreadyNotified := false
					for _, notifiedUserID := range notificationUserIDs {
						if notifiedUserID == headUser.ID {
							alreadyNotified = true
							break
						}
					}
					if !alreadyNotified {
						notificationUserIDs = append(notificationUserIDs, headUser.ID)
						emailTo = append(emailTo, headUser.Email)
					}
				}
			}
		}

		// Send notification if there are recipients
		if len(notificationUserIDs) > 0 {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              notificationUserIDs,
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Risk Register Responsibility Assigned",
				HeadingAr:            "مسؤولية تسجيل المخاطر المعينة",
				TitleEn:              riskRegisterFormDTO.Title,
				ReferenceNumber:      referenceNumber,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Use responsible user details for email if available, otherwise use unit head
			var firstNameEn, lastNameEn string
			if riskRegisterFormDTO.ResponsibleUserID != nil {
				responsibleUser, _ := rrups.userSvc.GetByID(*riskRegisterFormDTO.ResponsibleUserID)
				if responsibleUser != nil {
					firstNameEn = responsibleUser.FirstNameEn
					lastNameEn = responsibleUser.LastNameEn
				}
			}
			if firstNameEn == "" && len(emailTo) > 0 {
				// Fallback to first recipient's name
				firstUser, _ := rrups.userSvc.GetByID(notificationUserIDs[0])
				if firstUser != nil {
					firstNameEn = firstUser.FirstNameEn
					lastNameEn = firstUser.LastNameEn
				}
			}

			subject := fmt.Sprintf("Risk Register/%s: Responsibility Assigned", riskRegister.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":          firstNameEn,
				"LastNameEn":           lastNameEn,
				"RiskReferenceNumber":  existingRiskRegister.ReferenceNumber,
				"RiskTitle":            existingRiskRegister.Title,
				"RiskDescription":      utility.GetStringValue(existingRiskRegister.Description),
				"CreatedByFirstNameEn": createdBy.FirstNameEn,
				"CreatedByLastNameEn":  createdBy.LastNameEn,
				"AssignedOn":           time.Now().Format("2006-01-02"),
				"RedirectURL":          fmt.Sprintf("%s/risk/identified-risk/%d/workflow", appURL, id),
			}
			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "risk/templates/email/risk_register_create_responsible_user_mail.html",
				EmailData:            emailData,
			}
			err = rrups.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}

		// Fetch risk settings
		riskSetting, err := rrups.riskSettingsSvc.Get()
		if err != nil {
			utility.LogError("Failed to fetch risk setting", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Get Review User IDs
		reviewUserIDs, err := rrups.userSvc.GetUserIdsByRole(riskSetting.ReviewRoleID)
		if err != nil {
			utility.LogError("Failed to fetch review user IDs", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Ensure we have valid reviewers
		if len(reviewUserIDs) == 0 {
			fmt.Println("No review users found, skipping notifications")
			return
		}

		reviewEmailTo := make([]string, len(reviewUserIDs))
		reviewAction := "review-user-assign"

		for i, userID := range reviewUserIDs {
			user, err := rrups.userSvc.GetByID(userID)
			if err != nil {
				utility.LogError("Failed to fetch reviewer details", map[string]interface{}{
					"error":  err.Error(),
					"userID": userID,
				})
				continue
			}
			reviewEmailTo[i] = user.Email
		}

		dbNotificationReview := &commonDTOs.DBNotificationDTO{
			UserIDs:              reviewUserIDs,
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.InformationType),
			HeadingEn:            "Risk Register Review Assigned",
			HeadingAr:            "تم تعيين مراجعة سجل المخاطر",
			TitleEn:              riskRegister.Title,
			TitleAr:              riskRegister.Title,
			ReferenceNumber:      existingRiskRegister.ReferenceNumber,
			Entity:               &entity,
			Action:               &reviewAction,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		err = rrups.notificationSvc.SendNotification(dbNotificationReview, nil, "risk")
		if err != nil {
			utility.LogError("Failed to send DB notification", map[string]interface{}{
				"error": err.Error(),
			})
		}

		emailNotificationReview := &commonDTOs.EmailNotificationDTO{
			EmailTo:              reviewEmailTo,
			Subject:              fmt.Sprintf("Risk Register/%s: Review Assigned", riskRegister.Title),
			TemplateRelativePath: "risk/templates/email/risk_register_review_user_mail.html",
			EmailData: map[string]interface{}{
				"FirstNameEn":          riskSetting.ReviewRole.TitleEn,
				"RiskReferenceNumber":  existingRiskRegister.ReferenceNumber,
				"RiskTitle":            existingRiskRegister.Title,
				"RiskDescription":      utility.GetStringValue(existingRiskRegister.Description),
				"CreatedByFirstNameEn": createdBy.FirstNameEn,
				"CreatedByLastNameEn":  createdBy.LastNameEn,
				"AssignedOn":           time.Now().Format("2006-01-02"),
				"RedirectURL":          fmt.Sprintf("%s/risk/identified-risk/%d/workflow", utility.GetAppURL(), riskRegister.ID),
			},
		}

		err = rrups.notificationSvc.SendNotification(nil, emailNotificationReview, "risk")
		if err != nil {
			utility.LogError("Failed to send review user email", map[string]interface{}{
				"error":       err.Error(),
				"user_emails": reviewEmailTo,
			})
		}
	}()

	return id, nil
}

func (rrups *RiskRegisterUpdatePlanService) UpdateRiskIdentification(planID uint, id uint, currentUserID uint, riskRegisterFormDTO *riskDTOs.UpdatePlanRiskRegisterFormDTO) error {
	// Verify plan exists
	_, err := rrups.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	tx := rrups.repo.BeginTransaction()

	oldRiskRegister, err := rrups.riskRegisterSvc.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := rrups.validateRiskRegisterFormDTO(riskRegisterFormDTO); err != nil {
		return err
	}

	existingRiskRegister, err := rrups.riskRegisterSvc.GetByID(id)
	if err != nil {
		return err
	}

	// Verify risk belongs to this plan
	if existingRiskRegister.RiskRegisterUpdatePlanID == nil || *existingRiskRegister.RiskRegisterUpdatePlanID != planID {
		tx.Rollback()
		return utility.NewResourceNotFoundError(fmt.Errorf("risk not found in this plan"), id)
	}

	oldResponsibleUserID := existingRiskRegister.ResponsibleUserID

	existingRiskRegister.RiskCategoryID = riskRegisterFormDTO.RiskCategoryID
	existingRiskRegister.RiskSourceID = riskRegisterFormDTO.RiskSourceID
	existingRiskRegister.ResponsibleUserID = riskRegisterFormDTO.ResponsibleUserID
	existingRiskRegister.ResponsibleUnitID = riskRegisterFormDTO.ResponsibleUnitID
	existingRiskRegister.Title = riskRegisterFormDTO.Title
	existingRiskRegister.Description = riskRegisterFormDTO.Description
	existingRiskRegister.ReviewFrequencyID = riskRegisterFormDTO.ReviewFrequencyID
	existingRiskRegister.LastReviewNote = riskRegisterFormDTO.LastReviewNote
	existingRiskRegister.RiskCode = riskRegisterFormDTO.RiskCode
	existingRiskRegister.DivisionID = riskRegisterFormDTO.DivisionID
	existingRiskRegister.SectionID = riskRegisterFormDTO.SectionID
	existingRiskRegister.IsStrategic = riskRegisterFormDTO.IsStrategic
	existingRiskRegister.IsOperational = riskRegisterFormDTO.IsOperational
	existingRiskRegister.ExpectedCause = riskRegisterFormDTO.ExpectedCause
	existingRiskRegister.Consequence = riskRegisterFormDTO.Consequence
	existingRiskRegister.IsInternal = riskRegisterFormDTO.IsInternal
	existingRiskRegister.IsExternal = riskRegisterFormDTO.IsExternal
	existingRiskRegister.UpdatedBYID = &currentUserID

	if err := rrups.repo.UpdateRiskIdentification(existingRiskRegister, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingRisk, err := rrups.riskRegisterSvc.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {

		// Prepare log data
		changedData := make(map[string]interface{})

		// log changes
		oldData := map[string]interface{}{
			"title":            oldRiskRegister.Title,
			"description":      oldRiskRegister.Description,
			"risk_categrory":   oldRiskRegister.RiskCategory.TitleEn,
			"risk_source":      utility.GetFieldValue(oldRiskRegister.RiskSource).(string),
			"reported_at":      oldRiskRegister.ReportedAtString,
			"next_review_date": oldRiskRegister.NextReviewDateString,
		}

		if oldRiskRegister.ResponsibleUser != nil {
			oldData["responsible_user"] = fmt.Sprintf("%s %s", oldRiskRegister.ResponsibleUser.FirstNameEN, oldRiskRegister.ResponsibleUser.LastNameEN)
		}

		newData := map[string]interface{}{
			"title":            existingRisk.Title,
			"description":      existingRisk.Description,
			"risk_categrory":   existingRisk.RiskCategory.TitleEn,
			"risk_source":      utility.GetFieldValue(existingRisk.RiskSource).(string),
			"reported_at":      existingRisk.ReportedAtString,
			"next_review_date": existingRisk.NextReviewDateString,
		}

		if existingRisk.ResponsibleUser != nil {
			newData["responsible_user"] = fmt.Sprintf("%s %s", existingRisk.ResponsibleUser.FirstNameEN, existingRisk.ResponsibleUser.LastNameEN)
		}

		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := rrups.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskRegister{}.TableName(),
				id,
				changedData,
				currentUserID,
				"risk", "risk_register", "", nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"error":    err.Error(),
					"recordID": id,
				})
			}
		}

		// Send notifications to old and new responsible users if the responsible user has changed
		if oldResponsibleUserID != existingRiskRegister.ResponsibleUserID {
			// URL param for notifications
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "risk-register"

			// Fetch details of the user making the assignment
			updatedByUser, err := rrups.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			// Send notification to old responsible user if exists
			if oldResponsibleUserID != nil {
				oldResponsibleUserAction := "responsible-user-remove"
				oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*oldResponsibleUserID},
					ModuleIdentifier:     "risk",
					NotificationType:     string(commonDTOs.InformationType),
					HeadingEn:            "Risk Register Responsibility Removed",
					HeadingAr:            "تمت إزالة مسؤولية تسجيل المخاطر",
					TitleEn:              existingRiskRegister.Title,
					ReferenceNumber:      existingRiskRegister.ReferenceNumber,
					Entity:               &entity,
					Action:               &oldResponsibleUserAction,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				oldResponsibleUser, err := rrups.userSvc.GetByID(*oldResponsibleUserID)
				if err != nil {
					utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
						"error": err.Error(),
					})
				} else {
					oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
					oldResponsibleUserSubject := fmt.Sprintf("Risk Register/%s: Responsibility Removed", existingRiskRegister.Title)
					oldResponsibleUserEmailData := map[string]interface{}{
						"FirstNameEn":          oldResponsibleUser.FirstNameEn,
						"LastNameEn":           oldResponsibleUser.LastNameEn,
						"RiskReferenceNumber":  existingRiskRegister.ReferenceNumber,
						"RiskTitle":            existingRiskRegister.Title,
						"RiskDescription":      utility.GetStringValue(existingRiskRegister.Description),
						"UpdatedByFirstNameEn": updatedByUser.FirstNameEn,
						"UpdatedByLastNameEn":  updatedByUser.LastNameEn,
						"RedirectURL":          fmt.Sprintf("%s/risk/identified-risk/%d/workflow", appURL, id),
					}

					oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
						EmailTo:              oldResponsibleUserEmailTo,
						Subject:              oldResponsibleUserSubject,
						TemplateRelativePath: "risk/templates/email/risk_register_remove_responsible_user_mail.html",
						EmailData:            oldResponsibleUserEmailData,
					}

					err = rrups.notificationSvc.SendNotification(oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "risk")
					if err != nil {
						utility.LogError("Failed to send notifications to old responsible user", map[string]interface{}{
							"error": err.Error(),
						})
					}
				}
			}

			// Send notification to new responsible user and/or responsible unit head
			var notificationUserIDs []uint
			var emailTo []string
			newResponsibleUserAction := "responsible-user-assign"

			// If responsible user is assigned, add them to notification
			if riskRegisterFormDTO.ResponsibleUserID != nil {
				notificationUserIDs = append(notificationUserIDs, *riskRegisterFormDTO.ResponsibleUserID)
				newResponsibleUser, err := rrups.userSvc.GetByID(*riskRegisterFormDTO.ResponsibleUserID)
				if err != nil {
					utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
						"error": err.Error(),
					})
				} else {
					emailTo = append(emailTo, newResponsibleUser.Email)
				}
			}

			// Also send notification to responsible unit head
			if riskRegisterFormDTO.ResponsibleUnitID != nil {
				headUsers, err := rrups.organizationUnitsSvc.GetHeadUsersByOrganizationUnitID(*riskRegisterFormDTO.ResponsibleUnitID)
				if err != nil {
					utility.LogError("Failed to get responsible unit head users for notification", map[string]interface{}{
						"unit_id": *riskRegisterFormDTO.ResponsibleUnitID,
						"error":   err.Error(),
					})
				} else {
					for _, headUser := range headUsers {
						// Avoid duplicate notifications if the head user is also the responsible user
						alreadyNotified := false
						for _, notifiedUserID := range notificationUserIDs {
							if notifiedUserID == headUser.ID {
								alreadyNotified = true
								break
							}
						}
						if !alreadyNotified {
							notificationUserIDs = append(notificationUserIDs, headUser.ID)
							emailTo = append(emailTo, headUser.Email)
						}
					}
				}
			}

			// Send notification if there are recipients
			if len(notificationUserIDs) > 0 {
				newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              notificationUserIDs,
					ModuleIdentifier:     "risk",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Risk Register Responsibility Assigned",
					HeadingAr:            "مسؤولية تسجيل المخاطر المعينة",
					TitleEn:              existingRiskRegister.Title,
					ReferenceNumber:      existingRiskRegister.ReferenceNumber,
					Entity:               &entity,
					Action:               &newResponsibleUserAction,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				// Use responsible user details for email if available, otherwise use unit head
				var firstNameEn, lastNameEn string
				if riskRegisterFormDTO.ResponsibleUserID != nil {
					newResponsibleUser, _ := rrups.userSvc.GetByID(*riskRegisterFormDTO.ResponsibleUserID)
					if newResponsibleUser != nil {
						firstNameEn = newResponsibleUser.FirstNameEn
						lastNameEn = newResponsibleUser.LastNameEn
					}
				}
				if firstNameEn == "" && len(emailTo) > 0 {
					// Fallback to first recipient's name
					firstUser, _ := rrups.userSvc.GetByID(notificationUserIDs[0])
					if firstUser != nil {
						firstNameEn = firstUser.FirstNameEn
						lastNameEn = firstUser.LastNameEn
					}
				}

				newResponsibleUserSubject := fmt.Sprintf("Risk Register/%s: Responsibility Assigned", existingRiskRegister.Title)
				newResponsibleUserEmailData := map[string]interface{}{
					"FirstNameEn":           firstNameEn,
					"LastNameEn":            lastNameEn,
					"RiskReferenceNumber":   existingRiskRegister.ReferenceNumber,
					"RiskTitle":             existingRiskRegister.Title,
					"RiskDescription":       utility.GetStringValue(existingRiskRegister.Description),
					"AssignedByFirstNameEn": updatedByUser.FirstNameEn,
					"AssignedByLastNameEn":  updatedByUser.LastNameEn,
					"AssignedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":           fmt.Sprintf("%s/risk/identified-risk/%d/workflow", appURL, id),
				}

				newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              newResponsibleUserSubject,
					TemplateRelativePath: "risk/templates/email/risk_register_change_responsible_user_mail.html",
					EmailData:            newResponsibleUserEmailData,
				}

				err = rrups.notificationSvc.SendNotification(newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "risk")
				if err != nil {
					utility.LogError("Failed to send notifications to new responsible user", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}
	}()

	return nil
}

func (rrups *RiskRegisterUpdatePlanService) DeleteRiskIdentification(planID uint, ids []uint, currentUserID uint) ([]uint, error) {
	// Verify plan exists
	_, err := rrups.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, planID)
		}
		return nil, err
	}

	// Get timeline type for risk_identification
	timelineType, err := rrups.timelineTypeSvc.GetByType("risk_identification")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rrups.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	tx := rrups.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var notifications []struct {
		OwnerID           *uint
		ResponsibleUnitID *uint
		TitleEn           string
		TitleAr           string
		Description       *string
		ReferenceNumber   string
	}

	var deletedRiskRegister []map[string]interface{}

	for _, id := range ids {
		// Get risk model to verify plan ownership
		riskModel, err := rrups.riskRegisterSvc.GetByID(id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Verify risk belongs to this plan
		if riskModel.RiskRegisterUpdatePlanID == nil || *riskModel.RiskRegisterUpdatePlanID != planID {
			failedIds = append(failedIds, id)
			continue
		}

		risk, err := rrups.riskRegisterSvc.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		riskRegister, err := rrups.riskRegisterSvc.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		// Prepare notification data before deletion
		var responsibleUnitID *uint
		if riskRegister.ResponsibleUnit != nil {
			unitID := uint(riskRegister.ResponsibleUnit.ID)
			responsibleUnitID = &unitID
		}

		notifications = append(notifications, struct {
			OwnerID           *uint
			ResponsibleUnitID *uint
			TitleEn           string
			TitleAr           string
			Description       *string
			ReferenceNumber   string
		}{
			OwnerID:           &risk.ResponsibleUserID,
			ResponsibleUnitID: responsibleUnitID,
			TitleEn:           risk.Title,
			Description:       risk.Description,
			ReferenceNumber:   risk.ReferenceNumber,
		})

		deletedRiskRegister = append(deletedRiskRegister, map[string]interface{}{
			"id":    id,
			"title": risk.Title,
		})

		// Delete the risk
		if err := rrups.repo.DeleteRiskIdentification(id, tx); err != nil {
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

	// Check remaining risks after deletion
	remainingRiskCount, err := rrups.repo.CountRisksByPlanID(planID, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// If no risks remain, mark timeline as incomplete
	if remainingRiskCount == 0 {
		if err := rrups.repo.MarkTimelineAsIncomplete(timeline.ID, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Recalculate percentage
	if err := rrups.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := rrups.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.RiskRegister{}.TableName(),
				successfulIds,
				deletedRiskRegister,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			// Send DB and email notifications for each successfully deleted risk register
			for _, notification := range notifications {
				urlParam := ""
				entity := "risk-register"
				action := "delete"

				var notificationUserIDs []uint
				var emailTo []string

				// If responsible user is assigned, add them to notification
				if notification.OwnerID != nil {
					notificationUserIDs = append(notificationUserIDs, *notification.OwnerID)
					ownerUser, err := rrups.userSvc.GetByID(*notification.OwnerID)
					if err != nil {
						utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
							"error": err.Error(),
						})
					} else {
						emailTo = append(emailTo, ownerUser.Email)
					}
				}

				// Also send notification to responsible unit head
				if notification.ResponsibleUnitID != nil {
					headUsers, err := rrups.organizationUnitsSvc.GetHeadUsersByOrganizationUnitID(*notification.ResponsibleUnitID)
					if err != nil {
						utility.LogError("Failed to get responsible unit head users for notification", map[string]interface{}{
							"unit_id": *notification.ResponsibleUnitID,
							"error":   err.Error(),
						})
					} else {
						for _, headUser := range headUsers {
							// Avoid duplicate notifications if the head user is also the responsible user
							alreadyNotified := false
							for _, notifiedUserID := range notificationUserIDs {
								if notifiedUserID == headUser.ID {
									alreadyNotified = true
									break
								}
							}
							if !alreadyNotified {
								notificationUserIDs = append(notificationUserIDs, headUser.ID)
								emailTo = append(emailTo, headUser.Email)
							}
						}
					}
				}

				// Send notification if there are recipients
				if len(notificationUserIDs) > 0 {
					dbNotificationDTO := &commonDTOs.DBNotificationDTO{
						UserIDs:              notificationUserIDs,
						ModuleIdentifier:     "risk",
						NotificationType:     string(commonDTOs.DeletionType),
						HeadingEn:            "Risk Register Deleted",
						HeadingAr:            "تم حذف سجل المخاطر",
						TitleEn:              notification.TitleEn,
						ReferenceNumber:      notification.ReferenceNumber,
						Entity:               &entity,
						Action:               &action,
						URLParam:             urlParam,
						IsSystemNotification: false,
						CreatedByID:          &currentUserID,
					}

					// Use responsible user details for email if available, otherwise use unit head
					var firstNameEn, lastNameEn string
					if notification.OwnerID != nil {
						ownerUser, _ := rrups.userSvc.GetByID(*notification.OwnerID)
						if ownerUser != nil {
							firstNameEn = ownerUser.FirstNameEn
							lastNameEn = ownerUser.LastNameEn
						}
					}
					if firstNameEn == "" && len(emailTo) > 0 {
						// Fallback to first recipient's name
						firstUser, _ := rrups.userSvc.GetByID(notificationUserIDs[0])
						if firstUser != nil {
							firstNameEn = firstUser.FirstNameEn
							lastNameEn = firstUser.LastNameEn
						}
					}

					subject := fmt.Sprintf("Risk Register/%s: Deleted", notification.TitleEn)
					templateRelativePath := "risk/templates/email/risk_register_delete_mail.html"

					emailData := map[string]interface{}{
						"FirstNameEn":         firstNameEn,
						"LastNameEn":          lastNameEn,
						"RiskReferenceNumber": notification.ReferenceNumber,
						"RiskTitle":           notification.TitleEn,
						"RiskDescription":     notification.Description,
					}

					emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
						EmailTo:              emailTo,
						Subject:              subject,
						TemplateRelativePath: templateRelativePath,
						EmailData:            emailData,
					}

					// Send DB and email notifications
					err = rrups.notificationSvc.SendNotification(
						dbNotificationDTO, emailNotificationDTO, "risk")
					if err != nil {
						utility.LogError("Failed to send deletion notifications", map[string]interface{}{
							"error": err.Error(),
						})
					}
				}
			}
		}()
	}

	return failedIds, nil
}

func (rrups *RiskRegisterUpdatePlanService) MarkRiskIdentificationAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	_, err := rrups.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := rrups.timelineTypeSvc.GetByType("risk_identification")
	if err != nil {
		return err
	}

	timeline, err := rrups.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := rrups.repo.BeginTransaction()

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

	if err := rrups.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (rrups *RiskRegisterUpdatePlanService) validateRiskRegisterFormDTO(riskRegisterFormDTO *riskDTOs.UpdatePlanRiskRegisterFormDTO) error {

	if riskRegisterFormDTO.RiskCategoryID == nil {
		if _, err := rrups.riskCategorySvc.GetByID(*riskRegisterFormDTO.RiskCategoryID); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk category does not exist"), "not_exist", "risk category")
		}
	}

	if riskRegisterFormDTO.RiskSourceID != nil {
		if _, err := rrups.riskSourceSvc.GetByID(*riskRegisterFormDTO.RiskSourceID); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk source does not exist"), "not_exist", "risk source")
		}
	}

	if riskRegisterFormDTO.ResponsibleUserID != nil {
		if _, err := rrups.userSvc.GetByID(*riskRegisterFormDTO.ResponsibleUserID); err != nil {
			return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
		}
	}

	// Validate that division contains the responsible unit

	if !rrups.organizationUnitSvc.IsUnitUnderParent(*riskRegisterFormDTO.ResponsibleUnitID, riskRegisterFormDTO.DivisionID) {
		return utility.NewValidationError(fmt.Errorf("responsible unit must be under the selected division"), "invalid_hierarchy", "responsible_unit")
	}

	// Validate that section is under responsible unit
	if riskRegisterFormDTO.SectionID != nil && riskRegisterFormDTO.ResponsibleUnitID != nil {
		if !rrups.organizationUnitSvc.IsUnitUnderParent(*riskRegisterFormDTO.SectionID, *riskRegisterFormDTO.ResponsibleUnitID) {
			return utility.NewValidationError(fmt.Errorf("section must be under the responsible unit"), "invalid_hierarchy", "section")
		}
	}

	return nil
}

//Control Register Update

func (rrups *RiskRegisterUpdatePlanService) GetControlCountsByImplementationStatus(planID uint) (*riskDTOs.ControlRegisterUpdatedDTO, error) {
	return rrups.repo.GetControlCountsByImplementationStatus(planID)
}

func (rrups *RiskRegisterUpdatePlanService) CreateControlTest(planID uint, currentUserID uint, controlTestFormDTO *riskDTOs.UpdatePlanControlTestFormDTO) (uint, error) {
	// Verify plan exists
	_, err := rrups.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, utility.NewResourceNotFoundError(err, planID)
		}
		return 0, err
	}

	// Get timeline type for control_register_updated
	timelineType, err := rrups.timelineTypeSvc.GetByType("control_register_updated")
	if err != nil {
		return 0, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rrups.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return 0, err
	}

	tx := rrups.repo.BeginTransaction()
	defer tx.Rollback()

	// Convert UpdatePlanControlTestFormDTO to ControlTestFormDTO
	orgControlTestDTO := &organizationDTOs.ControlTestFormDTO{
		ControlImplementationStatusID: controlTestFormDTO.ControlImplementationStatusID,
		ControlEfficiencyStatusID:     controlTestFormDTO.ControlEfficiencyStatusID,
		Title:                         controlTestFormDTO.Title,
		Description:                   controlTestFormDTO.Description,
		Score:                         controlTestFormDTO.Score,
		Comments:                      controlTestFormDTO.Comments,
		TestedAt:                      controlTestFormDTO.TestedAt,
		TestedBy:                      controlTestFormDTO.TestedBy,
		Documents:                     controlTestFormDTO.Documents,
	}

	// Create control test using existing service with update_plan source
	controlTestID, err := rrups.controlTestSvc.Create(controlTestFormDTO.ControlID, currentUserID, nil, orgControlTestDTO, "update_plan", tx)
	if err != nil {
		return 0, err
	}

	// Update the control test with RiskRegisterUpdatePlanID
	if err := rrups.repo.UpdateControlTestPlanID(controlTestID, planID, tx); err != nil {
		return 0, err
	}

	// Update timeline completion
	if err := rrups.repo.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		return 0, err
	}

	// Recalculate percentage
	if err := rrups.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		return 0, err
	}

	tx.Commit()

	return controlTestID, nil
}

func (rrups *RiskRegisterUpdatePlanService) MarkControlRegisterUpdateAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	_, err := rrups.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := rrups.timelineTypeSvc.GetByType("control_register_updated")
	if err != nil {
		return err
	}

	timeline, err := rrups.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := rrups.repo.BeginTransaction()

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

	if err := rrups.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// Risk Assessment Workshop Document Operations

func (s *RiskRegisterUpdatePlanService) GetRiskAssessmentWorkshopDocuments(planID uint) ([]*riskDTOs.RiskAssessmentWorkshopDocumentDTO, error) {
	// Get timeline type for Risk Assessment Workshop
	timelineType, err := s.timelineTypeSvc.GetByType("risk_assessment_workshop")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	// Get documents by timeline ID
	documents, err := s.repo.GetRiskAssessmentWorkshopDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	var documentDTOs []*riskDTOs.RiskAssessmentWorkshopDocumentDTO
	for _, doc := range documents {
		dto := &riskDTOs.RiskAssessmentWorkshopDocumentDTO{
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

func (s *RiskRegisterUpdatePlanService) CreateRiskAssessmentWorkshopDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error {
	// Verify plan exists
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for Risk Assessment Workshop
	timelineType, err := s.timelineTypeSvc.GetByType("risk_assessment_workshop")
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
	var newDocuments []*models.RiskUpdatePlanRiskAssessmentWorkshopDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskUpdatePlanRiskAssessmentWorkshopDocument

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
				"risk-assessment-workshop-documents",
			)

			document := &models.RiskUpdatePlanRiskAssessmentWorkshopDocument{
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
			existingDocument, err := s.repo.GetRiskAssessmentWorkshopDocumentByID(*doc.ID)
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
			_, err := s.repo.CreateRiskAssessmentWorkshopDocument(doc, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Delete documents
	if len(deleteDocumentIDs) > 0 {
		for _, docID := range deleteDocumentIDs {
			if err := s.repo.DeleteRiskAssessmentWorkshopDocument(*docID, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// File system operations - store new files
	for _, doc := range newDocuments {
		storagePath := filepath.Join("risk", "risk-register-update-plan", strconv.FormatUint(uint64(planID), 10), "risk-assessment-workshop-documents")
		if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	// File system operations - delete files
	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("risk", "risk-register-update-plan", strconv.FormatUint(uint64(planID), 10), "risk-assessment-workshop-documents")
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
			models.RiskUpdatePlanRiskAssessmentWorkshopDocument{}.TableName(),
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
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (s *RiskRegisterUpdatePlanService) GetRiskAssessmentWorkshopDocumentPath(planID uint, token string) (string, error) {
	timelineType, err := s.timelineTypeSvc.GetByType("risk_assessment_workshop")
	if err != nil {
		return "", err
	}

	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return "", err
	}

	documents, err := s.repo.GetRiskAssessmentWorkshopDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return "", err
	}

	for _, doc := range documents {
		if doc.Token != nil && *doc.Token == token {
			decryptedFilename, err := utility.DecryptData(*doc.Token)
			if err != nil {
				return "", err
			}

			storagePath := filepath.Join(
				"risk",
				"risk-register-update-plan",
				strconv.FormatUint(uint64(planID), 10),
				"risk-assessment-workshop-documents",
			)

			baseStoragePath := config.GetEnv("STORAGE_PATH", "")
			fullPath := filepath.Join(baseStoragePath, storagePath, string(decryptedFilename))

			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				return "", utility.NewResourceNotFoundError(err, 0)
			}

			return fullPath, nil
		}
	}

	return "", utility.NewResourceNotFoundError(fmt.Errorf("document not found"), 0)
}

func (s *RiskRegisterUpdatePlanService) MarkRiskAssessmentWorkshopAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := s.timelineTypeSvc.GetByType("risk_assessment_workshop")
	if err != nil {
		return err
	}

	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

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

	if err := s.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// Risk Workshop Assessment Document Operations

func (s *RiskRegisterUpdatePlanService) GetRiskWorkshopAssessmentDocuments(planID uint) ([]*riskDTOs.RiskWorkshopAssessmentDocumentDTO, error) {
	// Get timeline type for Risk Workshop Assessment
	timelineType, err := s.timelineTypeSvc.GetByType("risk_workshop_assessment")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	// Get documents by timeline ID
	documents, err := s.repo.GetRiskWorkshopAssessmentDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	var documentDTOs []*riskDTOs.RiskWorkshopAssessmentDocumentDTO
	for _, doc := range documents {
		dto := &riskDTOs.RiskWorkshopAssessmentDocumentDTO{
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

func (s *RiskRegisterUpdatePlanService) CreateRiskWorkshopAssessmentDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error {
	// Verify plan exists
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for Risk Workshop Assessment
	timelineType, err := s.timelineTypeSvc.GetByType("risk_workshop_assessment")
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
	var newDocuments []*models.RiskUpdatePlanRiskWorkshopAssessmentDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskUpdatePlanRiskWorkshopAssessmentDocument

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
				"risk-workshop-assessment-documents",
			)

			document := &models.RiskUpdatePlanRiskWorkshopAssessmentDocument{
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
			existingDocument, err := s.repo.GetRiskWorkshopAssessmentDocumentByID(*doc.ID)
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
			_, err := s.repo.CreateRiskWorkshopAssessmentDocument(doc, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Delete documents
	if len(deleteDocumentIDs) > 0 {
		for _, docID := range deleteDocumentIDs {
			if err := s.repo.DeleteRiskWorkshopAssessmentDocument(*docID, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// File system operations - store new files
	for _, doc := range newDocuments {
		storagePath := filepath.Join("risk", "risk-register-update-plan", strconv.FormatUint(uint64(planID), 10), "risk-workshop-assessment-documents")
		if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	// File system operations - delete files
	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("risk", "risk-register-update-plan", strconv.FormatUint(uint64(planID), 10), "risk-workshop-assessment-documents")
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
			models.RiskUpdatePlanRiskWorkshopAssessmentDocument{}.TableName(),
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
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (s *RiskRegisterUpdatePlanService) GetRiskWorkshopAssessmentDocumentPath(planID uint, token string) (string, error) {
	timelineType, err := s.timelineTypeSvc.GetByType("risk_workshop_assessment")
	if err != nil {
		return "", err
	}

	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return "", err
	}

	documents, err := s.repo.GetRiskWorkshopAssessmentDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return "", err
	}

	for _, doc := range documents {
		if doc.Token != nil && *doc.Token == token {
			decryptedFilename, err := utility.DecryptData(*doc.Token)
			if err != nil {
				return "", err
			}

			storagePath := filepath.Join(
				"risk",
				"risk-register-update-plan",
				strconv.FormatUint(uint64(planID), 10),
				"risk-workshop-assessment-documents",
			)

			baseStoragePath := config.GetEnv("STORAGE_PATH", "")
			fullPath := filepath.Join(baseStoragePath, storagePath, string(decryptedFilename))

			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				return "", utility.NewResourceNotFoundError(err, 0)
			}

			return fullPath, nil
		}
	}

	return "", utility.NewResourceNotFoundError(fmt.Errorf("document not found"), 0)
}

func (s *RiskRegisterUpdatePlanService) MarkRiskWorkshopAssessmentAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	_, err := s.repo.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := s.timelineTypeSvc.GetByType("risk_workshop_assessment")
	if err != nil {
		return err
	}

	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

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

	if err := s.repo.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
