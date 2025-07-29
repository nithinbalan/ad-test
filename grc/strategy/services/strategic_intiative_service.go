package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	settingsSvc "jethur/grc/settings/services"
	"jethur/grc/strategy/dtos"
	"jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"

	organizationSvc "jethur/grc/organization/services"
)

type StrategicIntiativeServiceInterface interface {
	GetAll(params *dtos.StrategicIntiativeListQueryParamsDTO) (*dtos.StrategicIntiativeListResponseDTO, error)
	GetByID(id uint) (*models.StrategicIntiative, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.StrategicIntiativeDetailDTO, error)
	Create(currentUserID uint, strategicIntiativeFormDTO *dtos.StrategicIntiativeFormDTO) (uint, error)
	Update(id uint, currentUserID uint, dto *dtos.StrategicIntiativeFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type StrategicIntiativeService struct {
	repo                repositories.StrategicIntiativeRepositoryInterface
	activityLogSvc      services.ActivityLogServiceInterface
	organizationUnitSvc settingsSvc.OrganizationUnitServiceInterface
	userSvc             organizationSvc.UserServiceInterface
}

func NewStrategicIntiativeService(repo repositories.StrategicIntiativeRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	organizationUnitSvc settingsSvc.OrganizationUnitServiceInterface,
	userSvc organizationSvc.UserServiceInterface,

) StrategicIntiativeServiceInterface {
	return &StrategicIntiativeService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		organizationUnitSvc: organizationUnitSvc,
		userSvc:             userSvc,
	}
}

func (sis *StrategicIntiativeService) GetAll(params *dtos.StrategicIntiativeListQueryParamsDTO) (*dtos.StrategicIntiativeListResponseDTO, error) {
	return sis.repo.GetAll(params)
}

func (sis StrategicIntiativeService) GetByID(id uint) (*models.StrategicIntiative, error) {
	strategicIntiative, err := sis.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return strategicIntiative, nil

}

func (sis *StrategicIntiativeService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.StrategicIntiativeDetailDTO, error) {
	// Initialize DTO
	var strategicIntiativeDTO dtos.StrategicIntiativeDetailDTO

	// Get from repo
	strategicIntiative, err := sis.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map to DTO
	err = mapstructure.Decode(strategicIntiative, &strategicIntiativeDTO)
	if err != nil {
		return nil, err
	}

	// Format date fields safely
	if strategicIntiative.StartDate != nil {
		formattedDate := strategicIntiative.StartDate.Format("2006-01-02")
		strategicIntiativeDTO.StartDateString = &formattedDate
	}
	if strategicIntiative.EndDate != nil {
		formattedDate := strategicIntiative.EndDate.Format("2006-01-02")
		strategicIntiativeDTO.EndDateString = &formattedDate
	}

	strategicIntiativeDTO.CreatedAt = strategicIntiative.CreatedAt
	strategicIntiativeDTO.UpdatedAt = strategicIntiative.UpdatedAt

	return &strategicIntiativeDTO, nil
}

func (sis *StrategicIntiativeService) Create(currentUserID uint, dto *dtos.StrategicIntiativeFormDTO) (uint, error) {
	tx := sis.repo.BeginTransaction()

	// Validate responsible user
	responsibleUser, err := sis.userSvc.GetByID(dto.ResponsibleUserID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	organizationUnitID, err := sis.userSvc.GetUserPrimaryUnitID(dto.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for owner"), "not_exist", "process owner unit")
	}

	// Parse dates
	dates, err := sis.parseStrategicIntiativeDates(dto)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Build model
	model := &models.StrategicIntiative{
		TitleEn:           dto.TitleEn,
		TitleAr:           dto.TitleAr,
		DescriptionEn:     dto.DescriptionEn,
		DescriptionAr:     dto.DescriptionAr,
		StartDate:         dates.StartDate,
		EndDate:           dates.EndDate,
		ResponsibleUserID: responsibleUser.ID,
		ResponsibleUnitID: organizationUnitID,
		CreatedBYID:       currentUserID,
	}

	// Create
	id, err := sis.repo.Create(model, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	// Async activity log
	go func() {
		fullDetails, err := sis.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch strategic objective details", map[string]interface{}{"error": err.Error()})
			return
		}

		err = sis.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StrategicIntiative{}.TableName(),
			id,
			fullDetails,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StrategicIntiative{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (sis *StrategicIntiativeService) Update(id uint, currentUserID uint, dto *dtos.StrategicIntiativeFormDTO) error {
	existing, err := sis.GetByID(id)
	if err != nil {
		return err
	}

	// Validate responsible user
	responsibleUser, err := sis.userSvc.GetByID(dto.ResponsibleUserID)
	if err != nil {
		return err
	}

	//Parse dates
	dates, err := sis.parseStrategicIntiativeDates(dto)
	if err != nil {
		return err
	}

	// Update fields
	existing.TitleEn = dto.TitleEn
	existing.TitleAr = dto.TitleAr
	existing.DescriptionEn = dto.DescriptionEn
	existing.DescriptionAr = dto.DescriptionAr
	existing.StartDate = dates.StartDate
	existing.EndDate = dates.EndDate
	existing.ResponsibleUserID = responsibleUser.ID
	existing.UpdatedBYID = &currentUserID

	// Resolve Responsible Unit dynamically
	if dto.ResponsibleUserID != 0 {
		unitID, err := sis.userSvc.GetUserPrimaryUnitID(dto.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "organization_unit")
		}
		existing.ResponsibleUnitID = unitID
	} else {
		existing.ResponsibleUnitID = 0
	}

	// Save
	err = sis.repo.Update(existing)
	if err != nil {
		return err
	}

	// Log activity asynchronously
	go func() {
		fullDetails, err := sis.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch strategic objective details for activity log", map[string]interface{}{"error": err.Error()})
			return
		}

		err = sis.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.StrategicIntiative{}.TableName(),
			id,
			fullDetails,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log update activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.StrategicIntiative{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}

func (sis *StrategicIntiativeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := sis.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = sis.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := sis.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.StrategicIntiative{}.TableName(),
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

func (sis *StrategicIntiativeService) parseStrategicIntiativeDates(dto *dtos.StrategicIntiativeFormDTO) (*dtos.StrategicIntiativeParsedDates, error) {
	parsed := &dtos.StrategicIntiativeParsedDates{}

	if dto.StartDate != nil {
		sd, err := utility.ParseDate(*dto.StartDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start_date")
		}
		parsed.StartDate = sd
	}
	if dto.EndDate != nil {
		ed, err := utility.ParseDate(*dto.EndDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid end date"), "invalid_date", "end_date")
		}
		parsed.EndDate = ed
	}

	return parsed, nil

}
