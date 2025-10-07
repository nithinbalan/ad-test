package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceServices "jethur/grc/compliance/services"
	incidentServices "jethur/grc/incident/services"
	logServices "jethur/grc/logs/services"
	meetingDTOs "jethur/grc/meeting/dtos"
	meetingRepositories "jethur/grc/meeting/repositories"
	organizationServices "jethur/grc/organization/services"
	riskServices "jethur/grc/risk/services"

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
	GetAllCounts(meetingID uint) (*meetingDTOs.MeetingsCountResponseDTO, error)
	GetMeetingEffectiveness(meetingID uint) (float64, error)

	ExportMeeting(data []*meetingDTOs.MeetingListDTO, currentUserID uint) (*os.File, error)
	GetDocumentPath(id uint, token string, tx ...utility.DatabaseTransaction) (string, error)
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
	userSettingSvc           organizationServices.UserSettingServiceInterface
	notificationSvc          commonServices.NotificationServiceInterface
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
	userSettingSvc organizationServices.UserSettingServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
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
		userSettingSvc:           userSettingSvc,
		notificationSvc:          notificationSvc,
	}
}

func (ms *MeetingService) GetAll(currentUserID uint, params *meetingDTOs.MeetingListQueryParamsDTO) (*meetingDTOs.MeetingListResponseDTO, error) {
	// Retrieve user details to check for super admin role
	userDetails, err := ms.userSvc.GetCurrentUserDetails(currentUserID)
	if err != nil {
		return nil, err
	}

	// Determine if the user is a super admin
	isSuperAdmin := false
	if userDetails.Roles != nil {
		for _, role := range *userDetails.Roles {
			if role.IsSuperAdmin {
				isSuperAdmin = true
				break
			}
		}
	}

	// Pass the isSuperAdmin flag to the repository method
	return ms.repo.GetAll(currentUserID, params, isSuperAdmin)
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
	date := meeting.Date.Format("2006-01-02")
	meetingDTO.DateString = &date

	if meeting.Meeting != nil {
		meetingDate := meeting.Meeting.Date.Format("2006-01-02")
		meetingDTO.Meeting.DateString = &meetingDate
	}

	meetingDTO.CreatedAt = meeting.CreatedAt
	meetingDTO.UpdatedAt = meeting.UpdatedAt

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

	referenceNumber, err := ms.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}
	meeting := &models.Meeting{
		ReferenceNumber:    referenceNumber,
		Title:              meetingFormDTO.Title,
		Description:        meetingFormDTO.Description,
		Location:           meetingFormDTO.Location,
		Date:               dates.Date,
		StartTime:          meetingFormDTO.StartTime,
		EndTime:            meetingFormDTO.EndTime,
		OrganizerID:        meetingFormDTO.OrganizerID,
		IsPublic:           meetingFormDTO.IsPublic,
		CreatedByID:        currentUserID,
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

	// Save Documents
	var documents []*models.MeetingDocument
	for _, doc := range meetingFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("meeting", "meetings", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
		document := &models.MeetingDocument{
			MeetingID:    id,
			URL:          url,
			Name:         doc.FileName,
			Size:         doc.Size,
			Token:        encryptedFileName,
			OriginalName: &doc.OriginalName,
			Type:         doc.Type,
			CreatedAt:    time.Now(),
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := ms.repo.CreateMeetingDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Save Agendas
	// var agendas []*models.MeetingAgenda
	// for _, agenda := range meetingFormDTO.Agendas {
	// 	meetingAgenda := &models.MeetingAgenda{
	// 		MeetingID:   id,
	// 		Title:       agenda.Title,
	// 		TimeFrom:    agenda.TimeFrom,
	// 		TimeTo:      agenda.TimeTo,
	// 		CreatedBYID: currentUserID,
	// 	}
	// 	agendas = append(agendas, meetingAgenda)

	// }

	// if len(agendas) > 0 {
	// 	if err := ms.repo.CreateMeetingAgendas(agendas, tx); err != nil {
	// 		tx.Rollback()
	// 		return 0, err
	// 	}
	// }

	// Save Meeting Action Plans
	var actionPlans []*models.MeetingActionPlan

	for _, actionPlan := range meetingFormDTO.ActionPlans {
		meetingActionPlan := &models.MeetingActionPlan{
			MeetingID:         *meetingFormDTO.MeetingID,
			Title:             actionPlan.Title,
			Description:       actionPlan.Description,
			SeverityLevelID:   actionPlan.SeverityLevelID,
			ResponsibleUserID: actionPlan.ResponsibleUserID,
			// TargetDate:         actionPlan.TargetDate,
			Percentage:         actionPlan.Percentage,
			ActionPlanStatusID: actionPlan.ActionPlanStatusID,
			CreatedByID:        currentUserID,
		}
		actionPlans = append(actionPlans, meetingActionPlan)
	}
	if len(actionPlans) > 0 {
		if err := ms.repo.CreateMeetingActionPlans(actionPlans, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	var participants []*models.MeetingParticipant
	for _, participantID := range meetingFormDTO.UserIDs {
		participant := &models.MeetingParticipant{
			MeetingID:   id,
			UserID:      *participantID,
			CreatedByID: currentUserID,
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
		err := ms.createRecurringMeetings(meetingFormDTO, dates, currentUserID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	for _, doc := range documents {
		storagePath := filepath.Join("meeting", "meetings", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			return 0, err
		}
	}

	go func() {
		existingMeeting, err := ms.GetFullDetailsByID(meeting.ID)
		if err != nil {
			utility.LogError("Failed to fetch meeting", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = ms.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Meeting{}.TableName(),
			id,
			existingMeeting,
			currentUserID,
			"management_review", "meeting", "meeting", nil,
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

		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "meeting"
		action := "participant-assign"

		// Send notification to all participants
		for _, participantID := range meetingFormDTO.UserIDs {
			participant, err := ms.userSvc.GetByID(*participantID)
			if err != nil {
				utility.LogError("Failed to fetch meeting participant", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*participantID},
				ModuleIdentifier:     "meeting",
				NotificationType:     string(commonDTOs.AssignmentType),
				ReferenceNumber:      referenceNumber,
				HeadingEn:            "Meeting Scheduled",
				HeadingAr:            "جدولة الاجتماع",
				TitleEn:              meeting.Title,
				TitleAr:              meeting.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			createdBy, err := ms.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			// Fetch Meeting Status Title
			var meetingStatusTitle string
			meetingStatus, err := ms.meetingStatusSvc.GetByID(meetingFormDTO.MeetingStatusID)
			if err != nil {
				utility.LogError("Failed to fetch meeting status", map[string]interface{}{
					"error":    err.Error(),
					"statusID": meetingFormDTO.MeetingStatusID,
				})
			} else {
				if meetingStatus != nil {
					meetingStatusTitle = meetingStatus.TitleEn
				}
			}

			// Fetch Meeting Type Title
			var meetingTypeTitle string
			meetingType, err := ms.meetingTypeSvc.GetByID(meetingFormDTO.MeetingTypeID)
			if err != nil {
				utility.LogError("Failed to fetch meeting type", map[string]interface{}{
					"error":  err.Error(),
					"typeID": meetingFormDTO.MeetingTypeID,
				})
			} else {
				if meetingType != nil {
					meetingTypeTitle = meetingType.TitleEn
				}
			}

			// Send notification
			emailTo := []string{participant.Email}
			subject := fmt.Sprintf("Meeting/%s: You have been assigned", meeting.Title)
			appURL := utility.GetAppURL()

			emailData := map[string]interface{}{
				"FirstNameEn":        participant.FirstNameEn,
				"LastNameEn":         participant.LastNameEn,
				"ReferenceNumber":    meeting.ReferenceNumber,
				"Title":              meeting.Title,
				"Status":             meetingStatusTitle,
				"Type":               meetingTypeTitle,
				"Date":               meeting.Date.Format("2006-01-02"),
				"StartTime":          meeting.StartTime,
				"EndTime":            meeting.EndTime,
				"Description":        meeting.Description,
				"CreatedByFirstName": createdBy.FirstNameEn,
				"CreatedByLastName":  createdBy.LastNameEn,
				"AssignedOn":         time.Now().Format("2006-01-02"),
				"RedirectURL":        fmt.Sprintf("%s/review-meeting/meetings/%d/overview", appURL, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "meeting/templates/email/meeting_create_user_mail.html",
				EmailData:            emailData,
			}

			err = ms.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "meeting")
			if err != nil {
				utility.LogError("Failed to send meeting notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (ms *MeetingService) Update(id uint, currentUserID uint, meetingFormDTO *meetingDTOs.MeetingFormDTO) error {
	tx := ms.repo.BeginTransaction()

	oldMeeting, err := ms.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

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
	existingMeeting.Description = meetingFormDTO.Description
	existingMeeting.Location = meetingFormDTO.Location
	existingMeeting.Date = dates.Date
	existingMeeting.StartTime = meetingFormDTO.StartTime
	existingMeeting.EndTime = meetingFormDTO.EndTime
	existingMeeting.OrganizerID = meetingFormDTO.OrganizerID
	existingMeeting.IsPublic = meetingFormDTO.IsPublic
	existingMeeting.UpdatedByID = &currentUserID
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

	// Handle documents
	var documents []*models.MeetingDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.MeetingDocument

	for _, doc := range meetingFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Handle new documents
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("meeting", "meetings", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
			document := &models.MeetingDocument{
				MeetingID:    id,
				URL:          url,
				Name:         doc.FileName,
				Size:         doc.Size,
				Token:        encryptedFileName,
				OriginalName: &doc.OriginalName,
				Type:         doc.Type,
				CreatedAt:    time.Now(),
			}
			documents = append(documents, document)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			// Handle documents to be deleted
			if doc.ID != nil {
				existingDocument, err := ms.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, id)
					}
					return err
				}
				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			}
		}
	}

	// Create new documents
	if len(documents) > 0 {
		if err := ms.repo.CreateMeetingDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Delete marked documents
	if len(deleteDocumentIDs) > 0 {
		if err := ms.repo.DeleteMeetingDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update Participants
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
				MeetingID:   id,
				UserID:      *userID,
				CreatedByID: currentUserID,
			}
			participantsToAdd = append(participantsToAdd, participant)
		}
	}

	var removedParticipants []uint
	for oldUserID := range existingParticipantIDs {
		if _, exists := newParticipantIDs[oldUserID]; !exists {
			removedParticipants = append(removedParticipants, oldUserID)
		}
	}

	// Step 4: Identify New Participants (users in new list but not in old list)
	var newParticipants []uint
	for newUserID := range newParticipantIDs {
		if _, exists := existingParticipantIDs[newUserID]; !exists {
			newParticipants = append(newParticipants, newUserID)
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

	for _, doc := range documents {
		storagePath := filepath.Join("meeting", "meetings", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			fmt.Printf("Failed to store file %s: %v\n", doc.Name, err)
		}
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("meeting", "meetings", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			fmt.Printf("Failed to delete file %s: %v\n", doc.Name, err)
		}
	}

	go func() {
		existingReviewMeeting, err := ms.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		// Determine control and document changes

		participantsAdded := []uint{}
		participantsRemoved := []uint{}
		for id := range newParticipantIDs {
			if !existingParticipantIDs[id] {
				participantsAdded = append(participantsAdded, id)
			}
		}
		for id := range existingParticipantIDs {
			if !newParticipantIDs[id] {
				participantsRemoved = append(participantsRemoved, id)
			}
		}

		// Prepare log data
		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})
		if len(participantsAdded) > 0 || len(participantsRemoved) > 0 {
			associatedItems["participant_ids"] = map[string]interface{}{
				"added":   participantsAdded,
				"removed": participantsRemoved,
			}
		}

		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}
		// log changes
		oldData := map[string]interface{}{
			"title":             oldMeeting.Title,
			"description":       oldMeeting.Description,
			"location":          oldMeeting.Location,
			"date":              oldMeeting.DateString,
			"start_time":        oldMeeting.StartTime,
			"end_time":          oldMeeting.EndTime,
			"organizer":         fmt.Sprintf("%s %s", oldMeeting.Organizer.FirstNameEN, oldMeeting.Organizer.LastNameEN),
			"is_public":         oldMeeting.IsPublic,
			"meeting_type":      oldMeeting.MeetingType.TitleEn,
			"meeting_frequency": utility.GetFieldValue(oldMeeting.MeetingFrequency).(string),
			"meeting_status":    oldMeeting.MeetingStatus.TitleEn,
		}

		newData := map[string]interface{}{
			"title":             existingReviewMeeting.Title,
			"description":       existingReviewMeeting.Description,
			"location":          existingReviewMeeting.Location,
			"date":              existingReviewMeeting.DateString,
			"start_time":        existingReviewMeeting.StartTime,
			"end_time":          existingReviewMeeting.EndTime,
			"organizer":         fmt.Sprintf("%s %s", existingReviewMeeting.Organizer.FirstNameEN, existingReviewMeeting.Organizer.LastNameEN),
			"is_public":         existingReviewMeeting.IsPublic,
			"meeting_type":      existingReviewMeeting.MeetingType.TitleEn,
			"meeting_frequency": utility.GetFieldValue(existingReviewMeeting.MeetingFrequency).(string),
			"meeting_status":    existingReviewMeeting.MeetingStatus.TitleEn,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		// Log activity if there are any changes
		if len(changedData) > 0 {

			err := ms.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.Meeting{}.TableName(),
				id,
				changedData,
				currentUserID,
				"management_review", "meeting", "meeting", nil,
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
		}

		createdBy, err := ms.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Fetch Meeting Status Title
		var meetingStatusTitle string
		meetingStatus, err := ms.meetingStatusSvc.GetByID(meetingFormDTO.MeetingStatusID)
		if err != nil {
			utility.LogError("Failed to fetch meeting status", map[string]interface{}{
				"error":    err.Error(),
				"statusID": meetingFormDTO.MeetingStatusID,
			})
		} else {
			if meetingStatus != nil {
				meetingStatusTitle = meetingStatus.TitleEn
			}
		}

		// Fetch Meeting Type Title
		var meetingTypeTitle string
		meetingType, err := ms.meetingTypeSvc.GetByID(meetingFormDTO.MeetingTypeID)
		if err != nil {
			utility.LogError("Failed to fetch meeting type", map[string]interface{}{
				"error":  err.Error(),
				"typeID": meetingFormDTO.MeetingTypeID,
			})
		} else {
			if meetingType != nil {
				meetingTypeTitle = meetingType.TitleEn
			}
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "meeting"
		removeAction := "participant-remove"
		assignAction := "participant-assign"
		appURL := utility.GetAppURL()

		for _, removedUserID := range removedParticipants {

			removedUser, err := ms.userSvc.GetByID(removedUserID)
			if err != nil {
				utility.LogError("Failed to fetch removed participant details", map[string]interface{}{
					"userID": removedUserID,
					"error":  err.Error(),
				})
				continue
			}

			removedUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{removedUserID},
				ModuleIdentifier:     "meeting",
				NotificationType:     string(commonDTOs.InformationType),
				ReferenceNumber:      existingMeeting.ReferenceNumber,
				HeadingEn:            "Meeting Participation Removed",
				HeadingAr:            "تمت إزالة المشاركة في الاجتماع",
				TitleEn:              existingMeeting.Title,
				TitleAr:              existingMeeting.Title,
				Entity:               &entity,
				Action:               &removeAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			removedUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{removedUser.Email},
				Subject:              fmt.Sprintf("Meeting/%s: You have been removed", existingMeeting.Title),
				TemplateRelativePath: "meeting/templates/email/meeting_remove_user_mail.html",
				EmailData: map[string]interface{}{
					"FirstNameEn":        removedUser.FirstNameEn,
					"LastNameEn":         removedUser.LastNameEn,
					"ReferenceNumber":    existingMeeting.ReferenceNumber,
					"TitleEn":            existingMeeting.Title,
					"Description":        existingMeeting.Description,
					"Status":             meetingStatusTitle,
					"Type":               meetingTypeTitle,
					"Date":               meeting.Date.Format("2006-01-02"),
					"StartTime":          meeting.StartTime,
					"EndTime":            meeting.EndTime,
					"CreatedByFirstName": createdBy.FirstNameEn,
					"CreatedByLastName":  createdBy.LastNameEn,
					"RedirectURL":        fmt.Sprintf("%s/review-meeting/meetings/%d/overview", appURL, id),
				},
			}

			err = ms.notificationSvc.SendNotification(removedUserNotificationDTO, removedUserEmailNotificationDTO, "meeting")
			if err != nil {
				utility.LogError("Failed to send removal notification", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}

		for _, newUserID := range newParticipants {

			newUser, err := ms.userSvc.GetByID(newUserID)
			if err != nil {
				utility.LogError("Failed to fetch new participant details", map[string]interface{}{
					"userID": newUserID,
					"error":  err.Error(),
				})
				continue
			}

			newUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{newUserID},
				ModuleIdentifier:     "meeting",
				NotificationType:     string(commonDTOs.AssignmentType),
				ReferenceNumber:      existingMeeting.ReferenceNumber,
				HeadingEn:            "Meeting Participation Assigned",
				HeadingAr:            "تم تعيين المشاركة في الاجتماع",
				TitleEn:              existingMeeting.Title,
				TitleAr:              existingMeeting.Title,
				Entity:               &entity,
				Action:               &assignAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			newUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{newUser.Email},
				Subject:              fmt.Sprintf("Meeting/%s: You have been added", existingMeeting.Title),
				TemplateRelativePath: "meeting/templates/email/meeting_change_user_mail.html",
				EmailData: map[string]interface{}{
					"FirstNameEn":        newUser.FirstNameEn,
					"LastNameEn":         newUser.LastNameEn,
					"ReferenceNumber":    existingMeeting.ReferenceNumber,
					"Title":              existingMeeting.Title,
					"Description":        existingMeeting.Description,
					"Status":             meetingStatusTitle,
					"Type":               meetingTypeTitle,
					"Date":               meeting.Date.Format("2006-01-02"),
					"StartTime":          meeting.StartTime,
					"EndTime":            meeting.EndTime,
					"CreatedByFirstName": createdBy.FirstNameEn,
					"CreatedByLastName":  createdBy.LastNameEn,
					"AssignedOn":         time.Now().Format("2006-01-02"),
					"RedirectURL":        fmt.Sprintf("%s/review-meeting/meetings/%d/overview", appURL, id),
				},
			}

			err = ms.notificationSvc.SendNotification(newUserNotificationDTO, newUserEmailNotificationDTO, "meeting")
			if err != nil {
				utility.LogError("Failed to send assignment notification", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (ms *MeetingService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ms.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var deletedMeeting []map[string]interface{}

	for _, id := range ids {
		// Fetch meeting details including participants
		meeting, err := ms.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Notify participants about the meeting deletion
		if len(meeting.Participants) > 0 {
			entity := "meeting"
			action := "delete"
			urlParam := fmt.Sprintf(`{"id":%d}`, id)

			// Loop through participants to send notifications and emails
			for _, participant := range meeting.Participants {
				// Create system notification
				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{participant.UserID},
					ModuleIdentifier:     "meeting",
					NotificationType:     string(commonDTOs.InformationType),
					ReferenceNumber:      meeting.ReferenceNumber,
					HeadingEn:            "Meeting Deleted",
					HeadingAr:            "تم حذف الاجتماع",
					TitleEn:              meeting.Title,
					TitleAr:              meeting.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				// Fetch participant details
				user, err := ms.userSvc.GetByID(participant.UserID)
				if err != nil {
					utility.LogError("Failed to fetch participant details for deletion notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				// Prepare email notification
				emailTo := []string{user.Email}
				subject := fmt.Sprintf("Meeting/%s: Deleted", meeting.Title)
				emailData := map[string]interface{}{
					"FirstNameEn":     user.FirstNameEn,
					"LastNameEn":      user.LastNameEn,
					"Title":           meeting.Title,
					"Description":     meeting.Description,
					"ReferenceNumber": meeting.ReferenceNumber,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "meeting/templates/email/meeting_delete_user_mail.html",
					EmailData:            emailData,
				}

				// Send system and email notifications
				err = ms.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "meeting")
				if err != nil {
					utility.LogError("Failed to send meeting deletion notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
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
		deletedMeeting = append(deletedMeeting, map[string]interface{}{
			"id":    id,
			"title": meeting.Title,
		})
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
				deletedMeeting,
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

	// parsed.StartTime, err = utility.ParseMandatoryTime(*dto.StartTime)
	// if err != nil {
	// 	return nil, utility.NewValidationError(fmt.Errorf("invalid start time"), "invalid_time", "start time")
	// }

	// parsed.EndTime, err = utility.ParseMandatoryTime(*dto.EndTime)
	// if err != nil {
	// 	return nil, utility.NewValidationError(fmt.Errorf("invalid end time"), "invalid_time", "end time")
	// }

	return parsed, nil
}

func (ms *MeetingService) createRecurringMeetings(meetingFormDTO *meetingDTOs.MeetingFormDTO, dates *meetingDTOs.MeetingParsedDates, currentUserID uint, tx utility.DatabaseTransaction) error {
	frequency, err := ms.meetingFrequencySvc.GetByID(*meetingFormDTO.MeetingFrequencyID)
	if err != nil {
		return err
	}

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
		default:
			return fmt.Errorf("unsupported frequency type: %s", frequency.Type)
		}

		if currentDate.After(endDate) {
			break
		}

		// Create a new meeting instance for the recurring date
		// startTime := time.Date(1, 1, 1, dates.StartTime.Hour(), dates.StartTime.Minute(), dates.StartTime.Second(), 0, time.UTC)
		// endTime := time.Date(1, 1, 1, dates.EndTime.Hour(), dates.EndTime.Minute(), dates.EndTime.Second(), 0, time.UTC)

		referenceNumber, err := ms.repo.GenerateReferenceNumber()
		if err != nil {
			return err
		}
		recurringMeeting := &models.Meeting{
			ReferenceNumber:    referenceNumber,
			Title:              meetingFormDTO.Title,
			Description:        meetingFormDTO.Description,
			Location:           meetingFormDTO.Location,
			Date:               currentDate,
			StartTime:          meetingFormDTO.StartTime,
			EndTime:            meetingFormDTO.EndTime,
			OrganizerID:        meetingFormDTO.OrganizerID,
			IsPublic:           meetingFormDTO.IsPublic,
			CreatedByID:        currentUserID,
			MeetingTypeID:      meetingFormDTO.MeetingTypeID,
			MeetingFrequencyID: meetingFormDTO.MeetingFrequencyID,
			MeetingID:          meetingFormDTO.MeetingID,
			MeetingStatusID:    meetingFormDTO.MeetingStatusID,
		}

		newMeetingID, err := ms.repo.Create(recurringMeeting, tx)
		if err != nil {
			return err
		}

		// Add participants to the new recurring meeting
		var participants []*models.MeetingParticipant
		for _, participantID := range meetingFormDTO.UserIDs {
			participant := &models.MeetingParticipant{
				MeetingID:   newMeetingID,
				UserID:      *participantID,
				CreatedByID: currentUserID,
			}
			participants = append(participants, participant)
		}

		if len(participants) > 0 {
			if err := ms.repo.CreateParticipants(participants, tx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ms *MeetingService) GetAllCounts(meetingID uint) (*meetingDTOs.MeetingsCountResponseDTO, error) {
	if _, err := ms.GetByID(meetingID); err != nil {
		return nil, err
	}
	response, err := ms.repo.GetAllCounts(meetingID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (ms *MeetingService) GetMeetingEffectiveness(meetingID uint) (float64, error) {
	if _, err := ms.GetByID(meetingID); err != nil {
		return 0, err
	}
	response, err := ms.repo.GetMeetingEffectiveness(meetingID)
	if err != nil {
		return 0, err
	}

	return response, nil
}

func (ms *MeetingService) ExportMeeting(data []*meetingDTOs.MeetingListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := ms.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Ref No",
			"Title",
			"Meeting Category",
			"Date",
			"Time",
			"Location",
			"Organizer",
			"Status",
		}
	} else {
		headers = []string{
			"SL No.",
			"Ref No",
			"Title",
			"Meeting Category",
			"Date",
			"Time",
			"Location",
			"Organizer",
			"Status",
		}
	}
	// Prepare the data rows
	rows := make([][]interface{}, len(data))
	for i, meeting := range data {
		var meetingCategory, location, organizer, status string

		if userSetting.Language == "ar" {
			meetingCategory = utility.GetStringValue(meeting.MeetingTypeTitleAr)
			location = utility.GetStringValue(meeting.Location)
			organizer = utility.GetStringValue(meeting.OrganizerUserFirstNameAr)
			status = utility.GetStringValue(meeting.MeetingStatusTitleAr)
		} else {
			meetingCategory = utility.GetStringValue(meeting.MeetingTypeTitleAr)
			location = utility.GetStringValue(meeting.Location)
			organizer = utility.GetStringValue(meeting.OrganizerUserFirstNameAr)
			status = utility.GetStringValue(meeting.MeetingStatusTitleAr)
		}

		rows[i] = []interface{}{
			i + 1,
			meeting.ReferenceNumber,
			meeting.Title,
			meetingCategory,
			meeting.Date,
			meeting.StartTime,
			location,
			organizer,
			status,
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}

func (ms *MeetingService) GetDocumentPath(id uint, token string, tx ...utility.DatabaseTransaction) (string, error) {
	existingMeeting, err := ms.GetFullDetailsByID(id, tx...)
	if err != nil {
		return "", err
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("meeting", "meetings", strconv.FormatUint(uint64(existingMeeting.ID), 10), "documents")

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
