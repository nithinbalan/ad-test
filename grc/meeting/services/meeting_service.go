package services

import (
	"errors"
	"fmt"
	complianceServices "jethur/grc/compliance/services"
	incidentServices "jethur/grc/incident/services"
	logServices "jethur/grc/logs/services"
	meetingDTOs "jethur/grc/meeting/dtos"
	meetingRepositories "jethur/grc/meeting/repositories"
	organizationServices "jethur/grc/organization/services"
	riskServices "jethur/grc/risk/services"
	"time"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type MeetingServiceInterface interface {
	GetAll(currentUserID uint, params *meetingDTOs.MeetingListQueryParamsDTO) (*meetingDTOs.MeetingListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*meetingDTOs.MeetingDetailDTO, error)
	GetByID(id uint) (*models.Meeting, error)
	Create(currentUserID uint, meetingFormDTO *meetingDTOs.MeetingFormDTO) (uint, error)
	Update(id uint, currentUserID uint, meetingFormDTO *meetingDTOs.MeetingFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type MeetingService struct {
	repo                     meetingRepositories.MeetingRepositoryInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	meetingTypeSvc           MeetingTypeServiceInterface
	meetingFrequencySvc      MeetingFrequencyServiceInterface
	userSvc                  organizationServices.UserServiceInterface
	meetingStatusSvc         MeetingStatusServiceInterface
	controlSvc               organizationServices.ControlServiceInterface
	complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface
	riskRegisterSvc          riskServices.RiskRegisterServiceInterface
	processSvc               organizationServices.ProcessServiceInterface
	incidentRegisterSvc      incidentServices.IncidentRegisterServiceInterface
}

func NewMeetingService(
	repo meetingRepositories.MeetingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	meetingTypeSvc MeetingTypeServiceInterface,
	meetingFrequencySvc MeetingFrequencyServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	meetingStatusSvc MeetingStatusServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
	complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface,
	riskRegisterSvc riskServices.RiskRegisterServiceInterface,
	processSvc organizationServices.ProcessServiceInterface,
	incidentRegisterSvc incidentServices.IncidentRegisterServiceInterface,
) MeetingServiceInterface {
	return &MeetingService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		meetingTypeSvc:           meetingTypeSvc,
		meetingFrequencySvc:      meetingFrequencySvc,
		userSvc:                  userSvc,
		meetingStatusSvc:         meetingStatusSvc,
		controlSvc:               controlSvc,
		complianceRequirementSvc: complianceRequirementSvc,
		riskRegisterSvc:          riskRegisterSvc,
		processSvc:               processSvc,
		incidentRegisterSvc:      incidentRegisterSvc,
	}
}

func (ms *MeetingService) GetAll(currentUserID uint, params *meetingDTOs.MeetingListQueryParamsDTO) (*meetingDTOs.MeetingListResponseDTO, error) {
	return ms.repo.GetAll(currentUserID, params)
}

func (ms MeetingService) GetByID(id uint) (*models.Meeting, error) {
	meeting, err := ms.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return meeting, nil
}

func (ms MeetingService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*meetingDTOs.MeetingDetailDTO, error) {
	var meetingDTO meetingDTOs.MeetingDetailDTO

	meeting, err := ms.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(meeting, &meetingDTO)
	if err != nil {
		return nil, err
	}

	// Add extra fields required in details screen if any

	return &meetingDTO, nil
}

func (ms *MeetingService) Create(currentUserID uint, meetingFormDTO *meetingDTOs.MeetingFormDTO) (uint, error) {
	tx := ms.repo.BeginTransaction()

	if err := ms.validateMeetingFormDTO(meetingFormDTO); err != nil {
		return 0, err
	}

	dates, err := ms.parseMeetingDates(meetingFormDTO)
	if err != nil {
		return 0, err
	}

	// Extract time components from parsed StartTime and EndTime
	startTime := time.Date(1, 1, 1, dates.StartTime.Hour(), dates.StartTime.Minute(), dates.StartTime.Second(), 0, time.UTC)
	endTime := time.Date(1, 1, 1, dates.EndTime.Hour(), dates.EndTime.Minute(), dates.EndTime.Second(), 0, time.UTC)

	meeting := &models.Meeting{
		Title:              meetingFormDTO.Title,
		Location:           meetingFormDTO.Location,
		Date:               dates.Date,
		StartTime:          startTime,
		EndTime:            endTime,
		OrganizerID:        meetingFormDTO.OrganizerID,
		IsPublic:           meetingFormDTO.IsPublic,
		CreatedBy:          currentUserID,
		MeetingTypeID:      meetingFormDTO.MeetingTypeID,
		MeetingFrequencyID: meetingFormDTO.MeetingFrequencyID,
		MeetingID:          meetingFormDTO.MeetingID,
		MeetingStatusID:    meetingFormDTO.MeetingStatusID,
	}

	id, err := ms.repo.Create(meeting, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Save ComplianceRequirements
	var complianceRequirements []*models.MeetingComplianceRequirement
	for _, req := range meetingFormDTO.ComplianceRequirementIds {
		requirement := &models.MeetingComplianceRequirement{
			MeetingID:               id,
			ComplianceRequirementID: *req,
			CreatedBy:               currentUserID,
		}
		complianceRequirements = append(complianceRequirements, requirement)
	}

	if len(complianceRequirements) > 0 {
		if err := ms.repo.CreateMeetingComplianceRequirements(complianceRequirements, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Save Controls
	var controls []*models.MeetingControl
	for _, ctrl := range meetingFormDTO.ControlIds {
		control := &models.MeetingControl{
			MeetingID: id,
			ControlID: *ctrl,
			CreatedBy: currentUserID,
		}
		controls = append(controls, control)
	}

	if len(controls) > 0 {
		if err := ms.repo.CreateMeetingControls(controls, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Save Process
	var processes []*models.MeetingProcess
	for _, prcs := range meetingFormDTO.ProcessIds {
		control := &models.MeetingProcess{
			MeetingID: id,
			ProcessID: *prcs,
			CreatedBy: currentUserID,
		}
		processes = append(processes, control)
	}

	if len(processes) > 0 {
		if err := ms.repo.CreateMeetingProcesses(processes, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Save Risks
	var risks []*models.MeetingRisk
	for _, rsk := range meetingFormDTO.RiskIds {
		risk := &models.MeetingRisk{
			MeetingID:      id,
			RiskRegisterID: *rsk,
			CreatedBy:      currentUserID,
		}
		risks = append(risks, risk)
	}

	if len(risks) > 0 {
		if err := ms.repo.CreateMeetingRisks(risks, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Save Agendas
	var agendas []*models.MeetingAgenda
	for _, agenda := range meetingFormDTO.Agendas {
		meetingAgenda := &models.MeetingAgenda{
			MeetingID: id,
			Title:     agenda.Title,
			TimeFrom:  agenda.TimeFrom,
			TimeTo:    agenda.TimeTo,
			CreatedBy: currentUserID,
		}
		agendas = append(agendas, meetingAgenda)
	}

	if len(agendas) > 0 {
		if err := ms.repo.CreateMeetingAgendas(agendas, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Save Incident
	var incidents []*models.MeetingIncident
	for _, incd := range meetingFormDTO.IncidentIds {
		incident := &models.MeetingIncident{
			MeetingID:          id,
			IncidentRegisterID: *incd,
			CreatedBy:          currentUserID,
		}
		incidents = append(incidents, incident)
	}

	if len(incidents) > 0 {
		if err := ms.repo.CreateMeetingIncidents(incidents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	var participants []*models.MeetingParticipant
	for _, participantID := range meetingFormDTO.UserIDs {
		participant := &models.MeetingParticipant{
			MeetingID: id,
			UserID:    *participantID,
			CreatedBy: currentUserID,
		}
		participants = append(participants, participant)
	}

	if len(participants) > 0 {
		if err := ms.repo.CreateParticipants(participants, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if meetingFormDTO.MeetingFrequencyID != nil {
		err := ms.createRecurringMeetings(meeting, dates, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := ms.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Meeting{}.TableName(),
			id,
			meetingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.Meeting{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ms *MeetingService) Update(id uint, currentUserID uint, meetingFormDTO *meetingDTOs.MeetingFormDTO) error {
	tx := ms.repo.BeginTransaction()

	if err := ms.validateMeetingFormDTO(meetingFormDTO); err != nil {
		return err
	}

	existingMeeting, err := ms.GetByID(id)
	if err != nil {
		return err
	}

	dates, err := ms.parseMeetingDates(meetingFormDTO)
	if err != nil {
		return err
	}

	existingMeeting.MeetingID = meetingFormDTO.MeetingID
	existingMeeting.Title = meetingFormDTO.Title
	existingMeeting.Location = meetingFormDTO.Location
	existingMeeting.Date = dates.Date
	existingMeeting.StartTime = dates.StartTime
	existingMeeting.EndTime = dates.EndTime
	existingMeeting.OrganizerID = meetingFormDTO.OrganizerID
	existingMeeting.IsPublic = meetingFormDTO.IsPublic
	existingMeeting.UpdatedBy = &currentUserID
	existingMeeting.MeetingTypeID = meetingFormDTO.MeetingTypeID
	existingMeeting.MeetingFrequencyID = meetingFormDTO.MeetingFrequencyID
	existingMeeting.MeetingStatusID = meetingFormDTO.MeetingStatusID

	if err := ms.repo.Update(existingMeeting, tx); err != nil {
		tx.Rollback()
		return err
	}

	meeting, err := ms.repo.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingParticipantIDs := make(map[uint]bool)
	for _, participant := range meeting.Participants {
		existingParticipantIDs[participant.UserID] = true
	}

	newParticipantIDs := make(map[uint]bool)
	for _, userID := range meetingFormDTO.UserIDs {
		newParticipantIDs[*userID] = true
	}

	var participantsToAdd []*models.MeetingParticipant
	for _, userID := range meetingFormDTO.UserIDs {
		if _, exists := existingParticipantIDs[*userID]; !exists {
			participant := &models.MeetingParticipant{
				MeetingID: id,
				UserID:    *userID,
				CreatedBy: currentUserID,
			}
			participantsToAdd = append(participantsToAdd, participant)
		}
	}

	if len(participantsToAdd) > 0 {
		if err := ms.repo.CreateParticipants(participantsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	var participantIDsToDelete []uint
	for _, existingParticipant := range meeting.Participants {
		if _, exists := newParticipantIDs[existingParticipant.UserID]; !exists {
			participantIDsToDelete = append(participantIDsToDelete, existingParticipant.ID)
		}
	}

	// Delete participants not in the new list
	if len(participantIDsToDelete) > 0 {
		if err := ms.repo.DeleteParticipants(participantIDsToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		err := ms.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.Meeting{}.TableName(),
			id,
			meetingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.Meeting{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (ms *MeetingService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ms.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		// Fetch meeting details including participants
		meeting, err := ms.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Delete participants associated with the meeting
		if len(meeting.Participants) > 0 {
			var deleteParticipantIDs []uint
			for _, participant := range meeting.Participants {
				deleteParticipantIDs = append(deleteParticipantIDs, participant.ID)
			}
			if err := ms.repo.DeleteParticipants(deleteParticipantIDs, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		// Delete the meeting
		if err := ms.repo.Delete(id, tx); err != nil {
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
			err := ms.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.Meeting{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk archive activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (ms *MeetingService) validateMeetingFormDTO(meetingFormDTO *meetingDTOs.MeetingFormDTO) error {
	if _, err := ms.userSvc.GetByID(meetingFormDTO.OrganizerID); err != nil {
		return utility.NewValidationError(fmt.Errorf("organizer does not exist"), "not_exist", "organizer")
	}

	if _, err := ms.meetingTypeSvc.GetByID(meetingFormDTO.MeetingTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting type does not exist"), "not_exist", "meeting type")
	}

	if _, err := ms.meetingStatusSvc.GetByID(meetingFormDTO.MeetingStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting status does not exist"), "not_exist", "meeting status")
	}

	if meetingFormDTO.MeetingFrequencyID != nil {
		if _, err := ms.meetingFrequencySvc.GetByID(*meetingFormDTO.MeetingFrequencyID); err != nil {
			return utility.NewValidationError(fmt.Errorf("meeting frequency does not exist"), "not_exist", "meeting frequency")
		}
	}

	if meetingFormDTO.MeetingID != nil {
		if _, err := ms.GetByID(*meetingFormDTO.MeetingID); err != nil {
			return utility.NewValidationError(fmt.Errorf("meeting does not exist"), "not_exist", "meeting")
		}
	}

	//Controls exist
	controls := meetingFormDTO.ControlIds

	// Create a map to store control IDs
	controlIDs := make(map[uint]bool)

	// Check if each control exists
	for _, control := range controls {
		// Check if the control ID already exists in the map
		if controlIDs[*control] {
			// If the control ID already exists, return a validation error
			return utility.NewValidationError(fmt.Errorf("duplicate control ID: %d", control), "duplicate", "control")
		}

		// Check if the control exists in the database
		if _, err := ms.controlSvc.GetByID(*control); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", control), "not_exist", "control")
		}

		// Add the control ID to the map
		controlIDs[*control] = true
	}
	
	//Compliance Requirements Exist
	requirements := meetingFormDTO.ComplianceRequirementIds

	// Create a map to store requirement IDs
	requirementIDs := make(map[uint]bool)

	// Check if each requirement exists
	for _, requirement := range requirements {
		// Check if the requirement ID already exists in the map
		if requirementIDs[*requirement] {
			// If the requirement ID already exists, return a validation error
			return utility.NewValidationError(fmt.Errorf("duplicate compliance requirement ID: %d", requirement), "duplicate", "compliance requirement")
		}

		// Check if the requirement exists in the database
		if _, err := ms.complianceRequirementSvc.GetByComplianceRequirementID(*requirement); err != nil {
			return utility.NewValidationError(fmt.Errorf("compliance requirement does not exist for ID: %d", requirement), "not_exist", "compliance requirement")
		}

		// Add the requirement ID to the map
		requirementIDs[*requirement] = true
	}

	//Risks exist
	risks := meetingFormDTO.RiskIds
	// Create a map to store risk IDs
	riskIDs := make(map[uint]bool)

	// Check if each risk exists
	for _, risk := range risks {
		// Check if the risk ID already exists in the map
		if riskIDs[*risk] {
			// If the risk ID already exists, return a validation error
			return utility.NewValidationError(fmt.Errorf("duplicate risk ID: %d", risk), "duplicate", "risk")
		}

		// Check if the risk exists in the database
		if _, err := ms.riskRegisterSvc.GetByID(*risk); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk does not exist for ID: %d", risk), "not_exist", "risk")
		}

		// Add the risk ID to the map
		riskIDs[*risk] = true
	}

	//Process exist
	processes := meetingFormDTO.ProcessIds
	// Create a map to store process IDs
	processIds := make(map[uint]bool)

	// Check if each process exists
	for _, process := range processes {
		// Check if the process ID already exists in the map
		if processIds[*process] {
			// If the process ID already exists, return a validation error
			return utility.NewValidationError(fmt.Errorf("duplicate process ID: %d", process), "duplicate", "process")
		}

		// Check if the process exists in the database
		if _, err := ms.processSvc.GetByID(*process); err != nil {
			return utility.NewValidationError(fmt.Errorf("process does not exist for ID: %d", process), "not_exist", "process")
		}

		// Add the process ID to the map
		processIds[*process] = true
	}

	//Incidents exist
	incidents := meetingFormDTO.IncidentIds
	// Create a map to store incident IDs
	incidentIDs := make(map[uint]bool)

	// Check if each incident exists
	for _, incident := range incidents {
		// Check if the incident ID already exists in the map
		if incidentIDs[*incident] {
			// If the incident ID already exists, return a validation error
			return utility.NewValidationError(fmt.Errorf("duplicate incident ID: %d", incident), "duplicate", "incident")
		}

		// Check if the incident exists in the database
		if _, err := ms.incidentRegisterSvc.GetByID(*incident); err != nil {
			return utility.NewValidationError(fmt.Errorf("incident does not exist for ID: %d", incident), "not_exist", "incident")
		}

		// Add the incident ID to the map
		incidentIDs[*incident] = true
	}

	users := meetingFormDTO.UserIDs

	if len(users) == 0 {
		return nil
	}

	userIDs := make(map[uint]bool)

	for _, user := range users {
		if userIDs[*user] {
			return utility.NewValidationError(fmt.Errorf("duplicate user ID: %d", user), "duplicate", "user")
		}

		if _, err := ms.userSvc.GetByID(*user); err != nil {
			return utility.NewValidationError(fmt.Errorf("user does not exist for ID: %d", user), "not_exist", "user")
		}
		userIDs[*user] = true
	}
	return nil
}

func (ms *MeetingService) parseMeetingDates(dto *meetingDTOs.MeetingFormDTO) (*meetingDTOs.MeetingParsedDates, error) {
	parsed := &meetingDTOs.MeetingParsedDates{}

	var err error

	parsed.Date, err = utility.ParseMandatoryDate(dto.Date)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid date"), "invalid_date", "date")
	}

	parsed.StartTime, err = utility.ParseMandatoryTime(dto.StartTime)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid start time"), "invalid_time", "start time")
	}

	parsed.EndTime, err = utility.ParseMandatoryTime(dto.EndTime)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid end time"), "invalid_time", "end time")
	}

	return parsed, nil
}

func (ms *MeetingService) createRecurringMeetings(meeting *models.Meeting, dates *meetingDTOs.MeetingParsedDates, tx utility.DatabaseTransaction) error {
	frequency, err := ms.meetingFrequencySvc.GetByID(*meeting.MeetingFrequencyID)
	if err != nil {
		return err
	}

	fmt.Println(meeting)

	startDate := dates.Date

	endDate := startDate.AddDate(1, 0, 0) // add one year

	currentDate := startDate

	for currentDate.Before(endDate) {
		switch frequency.Type {
		case "weekly":
			currentDate = currentDate.AddDate(0, 0, 7)
		case "monthly":
			currentDate = currentDate.AddDate(0, 1, 0)
		case "quarterly":
			currentDate = currentDate.AddDate(0, 3, 0)
		case "half-yearly":
			currentDate = currentDate.AddDate(0, 6, 0)
		case "annually":
			currentDate = currentDate.AddDate(1, 0, 0)
		}

		if currentDate.After(endDate) {
			break
		}

		recurringMeeting := *meeting
		recurringMeeting.Date = currentDate
		recurringMeeting.ID = 0

		_, err := ms.repo.Create(&recurringMeeting, tx)
		if err != nil {
			return err
		}

	}

	return nil
}
