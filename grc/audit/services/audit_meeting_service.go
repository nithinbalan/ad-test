package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditMeetingServiceInterface interface {
	GetAll(auditID uint, params *auditDTOs.AuditMeetingListQueryParamsDTO) (*auditDTOs.AuditMeetingListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditMeetingDetailDTO, error)
	GetByID(id uint) (*models.AuditMeeting, error)
	Create(currentUserID uint, auditMeetingFormDTO *auditDTOs.AuditMeetingFormDTO) (uint, error)
	Update(id uint, currentUserID uint, auditMeetingFormDTO *auditDTOs.AuditMeetingFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type AuditMeetingService struct {
	repo           repositories.AuditMeetingRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	auditSvc       AuditServiceInterface
}

func NewAuditMeetingService(
	repo repositories.AuditMeetingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
) AuditMeetingServiceInterface {
	return &AuditMeetingService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		auditSvc:       auditSvc,
	}
}

func (ams *AuditMeetingService) GetAll(auditID uint, params *auditDTOs.AuditMeetingListQueryParamsDTO) (*auditDTOs.AuditMeetingListResponseDTO, error) {
	// Fetch the audit meetings and pagination info from the repository
	auditMeetings, pagination, err := ams.repo.GetAll(auditID, params)
	if err != nil {
		return nil, err
	}

	// Map the audit meetings to the response DTO
	auditMeetingDTOs := make([]*auditDTOs.AuditMeetingListDTO, len(auditMeetings))
	for i, meeting := range auditMeetings {
		var updatedAt *string
		if meeting.UpdatedAt != nil {
			formattedUpdatedAt := meeting.UpdatedAt.Format("2006-01-02 15:04:05")
			updatedAt = &formattedUpdatedAt
		}

		// Initialize the UpdatedUser fields to empty strings
		var updatedUserFirstNameEn, updatedUserFirstNameAr, updatedUserLastNameEn, updatedUserLastNameAr, updatedUserInitial, updatedUserColor, updatedUserEmail string
		var updatedUserProfileImagePath *string

		// Check if UpdatedByUser is not nil before accessing its fields
		if meeting.UpdatedByUser != nil {
			updatedUserFirstNameEn = meeting.UpdatedByUser.FirstNameEn
			updatedUserFirstNameAr = meeting.UpdatedByUser.FirstNameAr
			updatedUserLastNameEn = meeting.UpdatedByUser.LastNameEn
			updatedUserLastNameAr = meeting.UpdatedByUser.LastNameAr
			updatedUserInitial = meeting.UpdatedByUser.Initial
			updatedUserColor = meeting.UpdatedByUser.Color
			updatedUserEmail = meeting.UpdatedByUser.Email
			updatedUserProfileImagePath = meeting.UpdatedByUser.ProfileImagePath
		}

		auditMeetingDTOs[i] = &auditDTOs.AuditMeetingListDTO{
			ID:                          meeting.ID,
			AuditID:                     meeting.AuditID,
			Title:                       meeting.Title,
			Date:                        meeting.Date.Format("2006-01-02"),
			FromTime:                    meeting.FromTime,
			ToTime:                      meeting.ToTime,
			OtherParticipants:           meeting.OtherParticipants,
			Mom:                         meeting.Mom,
			CreatedAt:                   meeting.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:                   updatedAt,
			CreatedUserID:               meeting.CreatedBYID,
			CreatedUserFirstNameEn:      meeting.CreatedByUser.FirstNameEn,
			CreatedUserFirstNameAr:      meeting.CreatedByUser.FirstNameAr,
			CreatedUserLastNameEn:       meeting.CreatedByUser.LastNameEn,
			CreatedUserLastNameAr:       meeting.CreatedByUser.LastNameAr,
			CreatedUserInitial:          meeting.CreatedByUser.Initial,
			CreatedUserColor:            meeting.CreatedByUser.Color,
			CreatedUserEmail:            meeting.CreatedByUser.Email,
			CreatedUserProfileImagePath: meeting.CreatedByUser.ProfileImagePath,
			UpdatedUserID:               meeting.UpdatedBYID,
			UpdatedUserFirstNameEn:      updatedUserFirstNameEn,
			UpdatedUserFirstNameAr:      updatedUserFirstNameAr,
			UpdatedUserLastNameEn:       updatedUserLastNameEn,
			UpdatedUserLastNameAr:       updatedUserLastNameAr,
			UpdatedUserInitial:          updatedUserInitial,
			UpdatedUserColor:            updatedUserColor,
			UpdatedUserEmail:            updatedUserEmail,
			UpdatedUserProfileImagePath: updatedUserProfileImagePath,
			Participants:                mapParticipants(meeting.Participants),
		}
	}

	// Prepare and return the response with pagination info
	response := &auditDTOs.AuditMeetingListResponseDTO{
		Items:      auditMeetingDTOs,
		Pagination: *pagination,
	}
	return response, nil
}

func (ams *AuditMeetingService) GetByID(id uint) (*models.AuditMeeting, error) {
	auditMeeting, err := ams.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditMeeting, nil
}

func (ams *AuditMeetingService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditMeetingDetailDTO, error) {
	var auditMeetingDTO auditDTOs.AuditMeetingDetailDTO

	auditMeeting, err := ams.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Manually map fields from AuditMeeting to AuditMeetingDetailDTO
	auditMeetingDTO.ID = auditMeeting.ID
	auditMeetingDTO.AuditID = auditMeeting.AuditID
	auditMeetingDTO.Title = auditMeeting.Title
	auditMeetingDTO.Date = auditMeeting.Date.Format("2006-01-02")
	auditMeetingDTO.FromTime = auditMeeting.FromTime
	auditMeetingDTO.ToTime = auditMeeting.ToTime
	auditMeetingDTO.OtherParticipants = auditMeeting.OtherParticipants
	auditMeetingDTO.Mom = auditMeeting.Mom
	auditMeetingDTO.CreatedAt = auditMeeting.CreatedAt
	auditMeetingDTO.UpdatedAt = auditMeeting.UpdatedAt

	// Map CreatedByUser and UpdatedByUser to UserMinimalDTO
	auditMeetingDTO.CreatedByUser = organizationDTOs.UserMinimalDTO{
		ID:               int(auditMeeting.CreatedByUser.ID),
		FirstNameEN:      auditMeeting.CreatedByUser.FirstNameEn,
		FirstNameAR:      auditMeeting.CreatedByUser.FirstNameAr,
		LastNameEN:       auditMeeting.CreatedByUser.LastNameEn,
		LastNameAR:       auditMeeting.CreatedByUser.LastNameAr,
		Email:            auditMeeting.CreatedByUser.Email,
		Initial:          auditMeeting.CreatedByUser.Initial,
		Color:            auditMeeting.CreatedByUser.Color,
		ProfileImagePath: auditMeeting.CreatedByUser.ProfileImagePath,
	}

	if auditMeeting.UpdatedByUser != nil {
		auditMeetingDTO.UpdatedByUser = &organizationDTOs.UserMinimalDTO{
			ID:               int(auditMeeting.UpdatedByUser.ID),
			FirstNameEN:      auditMeeting.UpdatedByUser.FirstNameEn,
			FirstNameAR:      auditMeeting.UpdatedByUser.FirstNameAr,
			LastNameEN:       auditMeeting.UpdatedByUser.LastNameEn,
			LastNameAR:       auditMeeting.UpdatedByUser.LastNameAr,
			Email:            auditMeeting.UpdatedByUser.Email,
			Initial:          auditMeeting.UpdatedByUser.Initial,
			Color:            auditMeeting.UpdatedByUser.Color,
			ProfileImagePath: auditMeeting.UpdatedByUser.ProfileImagePath,
		}
	}

	// Map Participants to a slice of UserMinimalDTO
	participants := make([]*organizationDTOs.UserMinimalDTO, 0, len(auditMeeting.Participants))
	for _, participant := range auditMeeting.Participants {
		if participant.User != nil {
			participants = append(participants, &organizationDTOs.UserMinimalDTO{
				ID:               int(participant.User.ID),
				FirstNameEN:      participant.User.FirstNameEn,
				FirstNameAR:      participant.User.FirstNameAr,
				LastNameEN:       participant.User.LastNameEn,
				LastNameAR:       participant.User.LastNameAr,
				Email:            participant.User.Email,
				Initial:          participant.User.Initial,
				Color:            participant.User.Color,
				ProfileImagePath: participant.User.ProfileImagePath,
			})
		}
	}
	auditMeetingDTO.Participants = participants

	return &auditMeetingDTO, nil
}

func (ams *AuditMeetingService) Create(currentUserID uint, auditMeetingFormDTO *auditDTOs.AuditMeetingFormDTO) (uint, error) {
	tx := ams.repo.BeginTransaction()

	dates, err := ams.parseAuditMeetingDates(auditMeetingFormDTO)
	if err != nil {
		return 0, err
	}

	if err := ams.validateAuditMeetingFormDTO(auditMeetingFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	auditMeeting := &models.AuditMeeting{
		AuditID:           auditMeetingFormDTO.AuditID,
		Title:             auditMeetingFormDTO.Title,
		Date:              dates.Date,
		FromTime:          &auditMeetingFormDTO.FromTime,
		ToTime:            &auditMeetingFormDTO.ToTime,
		OtherParticipants: auditMeetingFormDTO.OtherParticipants,
		Mom:               auditMeetingFormDTO.Mom,
		CreatedBYID:       currentUserID,
	}

	id, err := ams.repo.Create(auditMeeting, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var participants []*models.AuditMeetingParticipant
	for _, participantID := range auditMeetingFormDTO.ParticipantIDs {
		participant := &models.AuditMeetingParticipant{
			AuditMeetingID: id,
			UserID:         participantID,
		}
		participants = append(participants, participant)
	}

	if len(participants) > 0 {
		if err := ams.repo.CreateParticipants(participants, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := ams.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditMeeting{}.TableName(),
			id,
			auditMeetingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditMeeting{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ams *AuditMeetingService) Update(id uint, currentUserID uint, auditMeetingFormDTO *auditDTOs.AuditMeetingFormDTO) error {
	tx := ams.repo.BeginTransaction()

	existingAuditMeeting, err := ams.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := ams.validateAuditMeetingFormDTO(auditMeetingFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	dates, err := ams.parseAuditMeetingDates(auditMeetingFormDTO)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingAuditMeeting.AuditID = auditMeetingFormDTO.AuditID
	existingAuditMeeting.Title = auditMeetingFormDTO.Title
	existingAuditMeeting.Date = dates.Date
	existingAuditMeeting.FromTime = &auditMeetingFormDTO.FromTime
	existingAuditMeeting.ToTime = &auditMeetingFormDTO.ToTime
	existingAuditMeeting.OtherParticipants = auditMeetingFormDTO.OtherParticipants
	existingAuditMeeting.Mom = auditMeetingFormDTO.Mom
	existingAuditMeeting.UpdatedBYID = &currentUserID

	err = ams.repo.Update(existingAuditMeeting, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingAuditMeetingDetails, err := ams.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingParticipantIDs := make(map[uint]bool)
	for _, participant := range existingAuditMeetingDetails.Participants {
		existingParticipantIDs[uint(participant.ID)] = true
	}

	newParticipantIDs := make(map[uint]bool)
	for _, participantID := range auditMeetingFormDTO.ParticipantIDs {
		newParticipantIDs[participantID] = true
	}

	var participantsToAdd []*models.AuditMeetingParticipant
	for _, newParticipantID := range auditMeetingFormDTO.ParticipantIDs {
		if _, exists := existingParticipantIDs[newParticipantID]; !exists {
			participant := &models.AuditMeetingParticipant{
				AuditMeetingID: id,
				UserID:         newParticipantID,
			}
			participantsToAdd = append(participantsToAdd, participant)
		}
	}

	if len(participantsToAdd) > 0 {
		if err := ams.repo.CreateParticipants(participantsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	var deleteParticipantIDs []*uint
	for _, existingParticipant := range existingAuditMeetingDetails.Participants {
		if _, exists := newParticipantIDs[uint(existingParticipant.ID)]; !exists {
			convertedID := uint(existingParticipant.ID)
			deleteParticipantIDs = append(deleteParticipantIDs, &convertedID)
		}
	}

	if len(deleteParticipantIDs) > 0 {
		if err := ams.repo.DeleteParticipant(existingAuditMeeting.ID, deleteParticipantIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		err := ams.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditMeeting{}.TableName(),
			id,
			auditMeetingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditMeeting{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (ams *AuditMeetingService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ams.repo.BeginTransaction()

	var failedIDs []uint
	var successfulIDs []uint

	for _, id := range ids {
		auditMeeting, err := ams.GetFullDetailsByID(id, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		// Delete associated participants if any
		if len(auditMeeting.Participants) > 0 {
			var deleteParticipantIDs []*uint
			for _, participant := range auditMeeting.Participants {
				convertedID := uint(participant.ID)
				deleteParticipantIDs = append(deleteParticipantIDs, &convertedID)
			}
			if err := ams.repo.DeleteParticipant(id, deleteParticipantIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if err := ams.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		} else {
			successfulIDs = append(successfulIDs, id)
		}
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()

	if len(successfulIDs) > 0 {
		go func() {
			err := ams.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AuditMeeting{}.TableName(),
				successfulIDs,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIDs, nil
}

func (ams *AuditMeetingService) validateAuditMeetingFormDTO(auditMeetingDTO *auditDTOs.AuditMeetingFormDTO) error {
	if _, err := ams.auditSvc.GetByID(auditMeetingDTO.AuditID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit does not exist"), "not_exist", "audit")
	}
	return nil
}

func (ams *AuditMeetingService) parseAuditMeetingDates(dto *auditDTOs.AuditMeetingFormDTO) (*auditDTOs.AuditMeetingParsedDates, error) {
	parsed := &auditDTOs.AuditMeetingParsedDates{}

	var err error

	parsed.Date, err = utility.ParseMandatoryDate(dto.Date)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid meeting date"), "invalid_date", "date")
	}

	return parsed, nil
}

func mapParticipants(participants []*models.AuditMeetingParticipant) []*organizationDTOs.UserMinimalDTO {
	var participantDTOs []*organizationDTOs.UserMinimalDTO
	for _, participant := range participants {
		if participant.User != nil {
			participantDTO := &organizationDTOs.UserMinimalDTO{
				ID:               int(participant.UserID),
				FirstNameEN:      participant.User.FirstNameEn,
				LastNameEN:       participant.User.LastNameEn,
				FirstNameAR:      participant.User.FirstNameAr,
				LastNameAR:       participant.User.LastNameAr,
				Email:            participant.User.Email,
				Initial:          participant.User.Initial,
				Color:            participant.User.Color,
				ProfileImagePath: participant.User.ProfileImagePath,
			}
			participantDTOs = append(participantDTOs, participantDTO)
		}
	}
	return participantDTOs
}
