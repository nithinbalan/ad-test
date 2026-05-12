package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

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

	// Timeline Operations
	GetTimelineByPlanIDAndTypeID(planID uint, timelineTypeID uint, tx ...utility.DatabaseTransaction) (*models.RiskUpdatePlanTimeline, error)
	UpdateTimelineCompletion(timelineID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
	CalculateAndUpdatePercentage(planID uint, tx ...utility.DatabaseTransaction) error
	MarkTimelineAsIncomplete(timelineID uint, tx ...utility.DatabaseTransaction) error

	// Get by plan ID
	GetByPlanID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlan, error)
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

// Timeline Operations

func (s *RiskRegisterUpdatePlanService) GetTimelineByPlanIDAndTypeID(planID uint, timelineTypeID uint, tx ...utility.DatabaseTransaction) (*models.RiskUpdatePlanTimeline, error) {
	timeline, err := s.repo.GetTimelineByPlanIDAndTypeID(planID, timelineTypeID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, planID)
		}
		return nil, err
	}

	return timeline, nil
}

func (s *RiskRegisterUpdatePlanService) UpdateTimelineCompletion(timelineID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	return s.repo.UpdateTimelineCompletion(timelineID, currentUserID, tx...)
}

func (s *RiskRegisterUpdatePlanService) CalculateAndUpdatePercentage(planID uint, tx ...utility.DatabaseTransaction) error {
	return s.repo.CalculateAndUpdatePercentage(planID, tx...)
}

func (s *RiskRegisterUpdatePlanService) MarkTimelineAsIncomplete(timelineID uint, tx ...utility.DatabaseTransaction) error {
	return s.repo.MarkTimelineAsIncomplete(timelineID, tx...)
}

func (s *RiskRegisterUpdatePlanService) GetByPlanID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlan, error) {
	plan, err := s.repo.GetByPlanID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return plan, nil
}
