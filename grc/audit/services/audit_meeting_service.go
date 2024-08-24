package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditMeetingServiceInterface interface {
	GetAll(params *auditDTOs.AuditMeetingListQueryParamsDTO) (*auditDTOs.AuditMeetingListResponseDTO, error)
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

func (ams *AuditMeetingService) GetAll(params *auditDTOs.AuditMeetingListQueryParamsDTO) (*auditDTOs.AuditMeetingListResponseDTO, error) {
	return ams.repo.GetAll(params)
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

	err = mapstructure.Decode(auditMeeting, &auditMeetingDTO)
	if err != nil {
		return nil, err
	}

	auditMeetingDTO.Date = auditMeeting.Date.Format("2006-01-02")
	auditMeetingDTO.FromTime = auditMeeting.FromTime.Format("15:04")
	auditMeetingDTO.ToTime = auditMeeting.ToTime.Format("15:04")
	auditMeetingDTO.CreatedAt = auditMeeting.CreatedAt
	auditMeetingDTO.UpdatedAt = auditMeeting.UpdatedAt

	return &auditMeetingDTO, nil
}

func (ams *AuditMeetingService) Create(currentUserID uint, auditMeetingFormDTO *auditDTOs.AuditMeetingFormDTO) (uint, error) {
	tx := ams.repo.BeginTransaction()

	dates, err := ams.parseAuditMeetingDates(auditMeetingFormDTO)
	if err != nil {
		return 0, err
	}

	// if err := ams.validateAuditMeetingFormDTO(auditMeetingFormDTO); err != nil {
	// 	tx.Rollback()
	// 	return 0, err
	// }

	auditMeeting := &models.AuditMeeting{
		AuditsID:          auditMeetingFormDTO.AuditID,
		Title:             auditMeetingFormDTO.Title,
		Date:              dates.Date,
		FromTime:          dates.FromTime,
		ToTime:            dates.ToTime,
		OtherParticipants: auditMeetingFormDTO.OtherParticipants,
		Mom:               auditMeetingFormDTO.Mom,
		CreatedBYID:       currentUserID,
	}

	id, err := ams.repo.Create(auditMeeting, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if len(auditMeetingFormDTO.ParticipantIDs) > 0 {
		if err := ams.repo.AddParticipants(id, auditMeetingFormDTO.ParticipantIDs, tx); err != nil {
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

	existingAuditMeeting, err := ams.repo.GetByID(id)
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
		return err
	}

	existingAuditMeeting.AuditsID = auditMeetingFormDTO.AuditID
	existingAuditMeeting.Title = auditMeetingFormDTO.Title
	existingAuditMeeting.Date = dates.Date
	existingAuditMeeting.FromTime = dates.FromTime
	existingAuditMeeting.ToTime = dates.ToTime
	existingAuditMeeting.OtherParticipants = auditMeetingFormDTO.OtherParticipants
	existingAuditMeeting.Mom = auditMeetingFormDTO.Mom
	existingAuditMeeting.UpdatedBYID = &currentUserID

	err = ams.repo.Update(existingAuditMeeting, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := ams.repo.AddParticipants(id, auditMeetingFormDTO.ParticipantIDs, tx); err != nil {
		tx.Rollback()
		return err
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
		// Get existing audit meeting details
		_, err := ams.GetFullDetailsByID(id, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		// Delete the audit meeting itself
		if err := ams.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := ams.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditMeeting{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditMeeting{}.TableName(),
						"recordID":  id,
						"userID":    currentUserID,
						"error":     err.Error(),
					}

					utility.LogError("Failed to log activity", fields)
				}
			}
		}()
	} else {
		tx.Rollback()
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
