package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/strategy/dtos"
	"jethur/grc/strategy/repositories"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"

	masterService "jethur/grc/masters/services"
	mastersSvc "jethur/grc/masters/services"
	organizationSvc "jethur/grc/organization/services"
	settingsSvc "jethur/grc/settings/services"
)

type StrategicObjectiveServiceInterface interface {
	GetAll(params *dtos.StrategicObjectiveListQueryParamsDTO) (*dtos.StrategicObjectiveListResponseDTO, error)
	GetByID(id uint) (*models.StrategicObjective, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.StrategicObjectiveDetailDTO, error)
	Create(currentUserID uint, strategicObjectiveFormDTO *dtos.StrategicObjectiveFormDTO) (uint, error)
	Update(id uint, currentUserID uint, dto *dtos.StrategicObjectiveFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
	UpdateStrategicObjective(strategicObjective *models.StrategicObjective, tx ...utility.DatabaseTransaction) error
}

// Service and implementations
type StrategicObjectiveService struct {
	repo                          repositories.StrategicObjectiveRepositoryInterface
	activityLogSvc                services.ActivityLogServiceInterface
	strategicObjectiveCategorySvc StrategicObjectiveCategoryServiceInterface
	organizationUnitSvc           settingsSvc.OrganizationUnitServiceInterface
	userSvc                       organizationSvc.UserServiceInterface
	reviewFrequencySvc            mastersSvc.ReviewFrequencyServiceInterface
	reviewStatusSvc               masterService.ReviewStatusServiceInterface
}

func NewStrategicObjectiveService(repo repositories.StrategicObjectiveRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	strategicObjectiveCategorySvc StrategicObjectiveCategoryServiceInterface,
	organizationUnitSvc settingsSvc.OrganizationUnitServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	reviewFrequencySvc mastersSvc.ReviewFrequencyServiceInterface,
	reviewStatusSvc masterService.ReviewStatusServiceInterface,
) StrategicObjectiveServiceInterface {
	return &StrategicObjectiveService{
		repo:                          repo,
		activityLogSvc:                activityLogSvc,
		strategicObjectiveCategorySvc: strategicObjectiveCategorySvc,
		organizationUnitSvc:           organizationUnitSvc,
		userSvc:                       userSvc,
		reviewFrequencySvc:            reviewFrequencySvc,
		reviewStatusSvc:               reviewStatusSvc,
	}
}

func (sos *StrategicObjectiveService) GetAll(params *dtos.StrategicObjectiveListQueryParamsDTO) (*dtos.StrategicObjectiveListResponseDTO, error) {
	return sos.repo.GetAll(params)
}

func (sos StrategicObjectiveService) GetByID(id uint) (*models.StrategicObjective, error) {
	strategicObjective, err := sos.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return strategicObjective, nil

}

func (sos *StrategicObjectiveService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.StrategicObjectiveDetailDTO, error) {
	// Initialize DTO
	var strategicObjectiveDTO dtos.StrategicObjectiveDetailDTO

	// Get from repo
	strategicObjective, err := sos.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map to DTO
	err = mapstructure.Decode(strategicObjective, &strategicObjectiveDTO)
	if err != nil {
		return nil, err
	}

	// Format date fields safely
	if strategicObjective.StartDate != nil {
		strategicObjectiveDTO.StartDateString = strategicObjective.StartDate.Format("2006-01-02")
	}
	if strategicObjective.EndDate != nil {
		strategicObjectiveDTO.EndDateString = strategicObjective.EndDate.Format("2006-01-02")
	}

	strategicObjectiveDTO.CreatedAt = strategicObjective.CreatedAt
	strategicObjectiveDTO.UpdatedAt = strategicObjective.UpdatedAt

	return &strategicObjectiveDTO, nil
}

func (sos *StrategicObjectiveService) Create(currentUserID uint, dto *dtos.StrategicObjectiveFormDTO) (uint, error) {
	tx := sos.repo.BeginTransaction()

	// Check for unique title
	if exists, err := sos.repo.ExistsByTitle(dto.TitleEn, dto.TitleAr); err != nil {
		tx.Rollback()
		return 0, err
	} else if exists {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	// Validate related entities
	if err := sos.validateStrategicObjectiveFormDTO(dto); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Parse dates
	dates, err := sos.parseStrategicObjectiveDates(dto)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Generate reference number
	referenceNumber, err := sos.repo.GenerateReferenceNumber()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var nextReviewDate time.Time
	if dto.ReviewFrequencyID != nil {
		reviewFrequency, err := sos.reviewFrequencySvc.GetByID(*dto.ReviewFrequencyID)
		if err != nil {
			return 0, err
		}
		nextReviewDate = utility.CalculateNextReviewDate(reviewFrequency.Type)
	}
	reviewStatus, err := sos.reviewStatusSvc.GetByType("review-pending")
	if err != nil {
		return 0, err
	}
	// Build model
	model := &models.StrategicObjective{
		TitleEn:                      dto.TitleEn,
		TitleAr:                      dto.TitleAr,
		ReferenceNumber:              referenceNumber,
		StrategicObjectiveCategoryID: dto.StrategicObjectiveCategoryID,
		DescriptionEn:                dto.DescriptionEn,
		DescriptionAr:                dto.DescriptionAr,
		StartDate:                    dates.StartDate,
		EndDate:                      dates.EndDate,
		ResponsibleUserID:            dto.ResponsibleUserID,
		ResponsibleUnitID:            dto.ResponsibleUnitID,
		ReviewUserID:                 dto.ReviewUserID,
		ReviewFrequencyID:            dto.ReviewFrequencyID,
		Status:                       dto.Status,
		NextReviewDate:               &nextReviewDate,
		ReviewStatusID:               &reviewStatus.ID,
		CreatedBYID:                  currentUserID,
	}

	// Create
	id, err := sos.repo.Create(model, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	// Async activity log
	go func() {
		fullDetails, err := sos.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch strategic objective details", map[string]interface{}{"error": err.Error()})
			return
		}

		err = sos.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StrategicObjective{}.TableName(),
			id,
			fullDetails,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StrategicObjective{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (sos *StrategicObjectiveService) Update(id uint, currentUserID uint, dto *dtos.StrategicObjectiveFormDTO) error {
	existing, err := sos.GetByID(id)
	if err != nil {
		return err
	}

	// Check title uniqueness if changed
	if dto.TitleEn != existing.TitleEn || dto.TitleAr != existing.TitleAr {
		if exists, err := sos.repo.ExistsByTitleExcludingID(id, dto.TitleEn, dto.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	// Validate related entities
	if err := sos.validateStrategicObjectiveFormDTO(dto); err != nil {
		return err
	}

	// Parse dates
	dates, err := sos.parseStrategicObjectiveDates(dto)
	if err != nil {
		return err
	}
	var nextReviewDate time.Time
	if dto.ReviewFrequencyID != nil {
		reviewFrequency, err := sos.reviewFrequencySvc.GetByID(*dto.ReviewFrequencyID)
		if err != nil {
			return err
		}
		nextReviewDate = utility.CalculateNextReviewDate(reviewFrequency.Type)
	}
	// Update fields
	existing.TitleEn = dto.TitleEn
	existing.TitleAr = dto.TitleAr
	existing.StrategicObjectiveCategoryID = dto.StrategicObjectiveCategoryID
	existing.DescriptionEn = dto.DescriptionEn
	existing.DescriptionAr = dto.DescriptionAr
	existing.StartDate = dates.StartDate
	existing.EndDate = dates.EndDate
	existing.ResponsibleUserID = dto.ResponsibleUserID
	existing.ReviewUserID = dto.ReviewUserID
	existing.ReviewFrequencyID = dto.ReviewFrequencyID
	existing.NextReviewDate = &nextReviewDate
	existing.Status = dto.Status
	existing.UpdatedBYID = &currentUserID

	// Resolve Responsible Unit dynamically
	if dto.ResponsibleUserID != nil {
		unitID, err := sos.userSvc.GetUserPrimaryUnitID(*dto.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "organization_unit")
		}
		existing.ResponsibleUnitID = &unitID
	} else {
		existing.ResponsibleUnitID = nil
	}

	// Save
	err = sos.repo.Update(existing)
	if err != nil {
		return err
	}

	// Log activity asynchronously
	go func() {
		fullDetails, err := sos.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch strategic objective details for activity log", map[string]interface{}{"error": err.Error()})
			return
		}

		err = sos.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.StrategicObjective{}.TableName(),
			id,
			fullDetails,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log update activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.StrategicObjective{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}

func (sos *StrategicObjectiveService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := sos.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = sos.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := sos.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.StrategicObjective{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (sos *StrategicObjectiveService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := sos.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := sos.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.StrategicObjective{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (sos *StrategicObjectiveService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := sos.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := sos.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.StrategicObjective{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (sos *StrategicObjectiveService) UpdateStrategicObjective(strategicObjective *models.StrategicObjective, tx ...utility.DatabaseTransaction) error {
	err := sos.repo.Update(strategicObjective, tx...)
	if err != nil {
		return err
	}

	return nil
}

func (sos *StrategicObjectiveService) validateStrategicObjectiveFormDTO(dto *dtos.StrategicObjectiveFormDTO) error {
	// Example: validate required category
	_, err := sos.strategicObjectiveCategorySvc.GetByID(dto.StrategicObjectiveCategoryID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("invalid strategic objective category"), "not_exist", "strategic_objective_category")
	}

	// Validate responsible user exists if given
	if dto.ResponsibleUserID != nil {
		_, err := sos.userSvc.GetByID(*dto.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible_user")
		}
	}

	// Validate review user exists if given
	if dto.ReviewUserID != nil {
		_, err := sos.userSvc.GetByID(*dto.ReviewUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("review user does not exist"), "not_exist", "review_user")
		}
	}

	// Validate review frequency exists if given
	if dto.ReviewFrequencyID != nil {
		_, err := sos.reviewFrequencySvc.GetByID(*dto.ReviewFrequencyID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review_frequency")
		}
	}
	_, err = sos.organizationUnitSvc.GetByID(*dto.ResponsibleUnitID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("organization unit does not exist"), "not_exist", "organization_unit")
	}

	return nil
}

func (sos *StrategicObjectiveService) parseStrategicObjectiveDates(dto *dtos.StrategicObjectiveFormDTO) (*dtos.StrategicObjectiveParsedDates, error) {
	parsed := &dtos.StrategicObjectiveParsedDates{}

	startDate, err := utility.ParseMandatoryDate(dto.StartDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start_date")
	}

	endDate, err := utility.ParseMandatoryDate(dto.EndDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid end date"), "invalid_date", "end_date")
	}

	parsed.StartDate = &startDate
	parsed.EndDate = &endDate

	return parsed, nil
}
