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

type AuditFindingServiceInterface interface {
	GetAll(params *auditDTOs.AuditFindingListQueryParamsDTO) (*auditDTOs.AuditFindingListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingDetailDTO, error)
	GetByID(id uint) (*models.AuditFinding, error)
	Create(currentUserID uint, auditFindingFormDTO *auditDTOs.AuditFindingFormDTO) (uint, error)
	Update(id uint, currentUserID uint, auditFindingFormDTO *auditDTOs.AuditFindingFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type AuditFindingService struct {
	repo             repositories.AuditFindingRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	auditPlanSvc     AuditPlanServiceInterface
	auditTimelineSvc AuditTimelineServiceInterface
}

func NewAuditFindingService(
	repo repositories.AuditFindingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditPlanSvc AuditPlanServiceInterface,
	auditTimelineSvc AuditTimelineServiceInterface,
) AuditFindingServiceInterface {
	return &AuditFindingService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		auditPlanSvc:     auditPlanSvc,
		auditTimelineSvc: auditTimelineSvc,
	}
}

func (afs *AuditFindingService) GetAll(params *auditDTOs.AuditFindingListQueryParamsDTO) (*auditDTOs.AuditFindingListResponseDTO, error) {
	return afs.repo.GetAll(params)
}

func (afs *AuditFindingService) GetByID(id uint) (*models.AuditFinding, error) {
	auditFinding, err := afs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditFinding, nil
}

func (afs *AuditFindingService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingDetailDTO, error) {
	var auditFindingDTO auditDTOs.AuditFindingDetailDTO

	auditFinding, err := afs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditFinding, &auditFindingDTO)
	if err != nil {
		return nil, err
	}

	auditFindingDTO.StartDateString = auditFinding.StartDate.Format("2006-01-02")
	auditFindingDTO.EndDateString = auditFinding.EndDate.Format("2006-01-02")

	auditFindingDTO.CreatedAt = auditFinding.CreatedAt
	auditFindingDTO.UpdatedAt = auditFinding.UpdatedAt

	return &auditFindingDTO, nil
}

func (afs *AuditFindingService) Create(currentUserID uint, auditFindingFormDTO *auditDTOs.AuditFindingFormDTO) (uint, error) {
	tx := afs.repo.BeginTransaction()

	dates, err := afs.parseAuditPlanDates(auditFindingFormDTO)
	if err != nil {
		return 0, err
	}

	if err := afs.validateAuditFindingFormDTO(auditFindingFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	auditFinding := &models.AuditFinding{
		AuditPlanID: auditFindingFormDTO.AuditPlanID,
		StartDate:   dates.StartDate,
		EndDate:     dates.EndDate,
		CreatedBYID: currentUserID,
	}

	id, err := afs.repo.Create(auditFinding, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := afs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditFinding{}.TableName(),
			id,
			auditFindingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditFinding{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (afs *AuditFindingService) Update(id uint, currentUserID uint, auditFindingFormDTO *auditDTOs.AuditFindingFormDTO) error {
	tx := afs.repo.BeginTransaction()

	existingAuditFinding, err := afs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := afs.validateAuditFindingFormDTO(auditFindingFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	dates, err := afs.parseAuditPlanDates(auditFindingFormDTO)
	if err != nil {
		return err
	}

	existingAuditFinding.AuditPlanID = auditFindingFormDTO.AuditPlanID
	existingAuditFinding.StartDate = dates.StartDate
	existingAuditFinding.EndDate = dates.EndDate
	existingAuditFinding.UpdatedBYID = &currentUserID

	err = afs.repo.Update(existingAuditFinding, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := afs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditFinding{}.TableName(),
			id,
			auditFindingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditFinding{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (afs *AuditFindingService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := afs.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint

	for _, id := range ids {
		// Get existing audit finding details
		_, err := afs.GetFullDetailsByID(id, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		if err := afs.auditTimelineSvc.Delete(id, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

		// Delete the audit finding itself
		if err := afs.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := afs.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditFinding{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditFinding{}.TableName(),
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

func (afs *AuditFindingService) validateAuditFindingFormDTO(dto *auditDTOs.AuditFindingFormDTO) error {
	if _, err := afs.auditPlanSvc.GetByID(dto.AuditPlanID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit plan does not exist"), "not_exist", "audit plan")
	}

	return nil
}

func (afs *AuditFindingService) parseAuditPlanDates(dto *auditDTOs.AuditFindingFormDTO) (*auditDTOs.AuditFindingParsedDates, error) {
	parsed := &auditDTOs.AuditFindingParsedDates{}

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
