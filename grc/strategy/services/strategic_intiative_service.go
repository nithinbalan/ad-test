package services

import (
	"errors"
	"fmt"
	"time"

	"jethur/grc/logs/services"
	"jethur/grc/models"
	settingsSvc "jethur/grc/settings/services"
	"jethur/grc/strategy/dtos"
	"jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"

	organizationDTOs "jethur/grc/organization/dtos"
	organizationSvc "jethur/grc/organization/services"
)

type StrategicIntiativeServiceInterface interface {
	GetAll(params *dtos.StrategicIntiativeListQueryParamsDTO) (*dtos.StrategicIntiativeListResponseDTO, error)
	GetByID(id uint) (*models.StrategicIntiative, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.StrategicIntiativeDetailDTO, error)
	Create(currentUserID uint, strategicIntiativeFormDTO *dtos.StrategicIntiativeFormDTO) (uint, error)
	Update(id uint, currentUserID uint, dto *dtos.StrategicIntiativeFormDTO) error
	UpdateStatusAndCompletion(id uint, currentUserID uint, dto *dtos.StrategicIntiativeStatusCompletionUpdateDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type StrategicIntiativeService struct {
	repo                          repositories.StrategicIntiativeRepositoryInterface
	strategicInitiativeStatusRepo repositories.StrategicInitiativeStatusRepositoryInterface
	activityLogSvc                services.ActivityLogServiceInterface
	organizationUnitSvc           settingsSvc.OrganizationUnitServiceInterface
	userSvc                       organizationSvc.UserServiceInterface
}

func NewStrategicIntiativeService(repo repositories.StrategicIntiativeRepositoryInterface,
	strategicInitiativeStatusRepo repositories.StrategicInitiativeStatusRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	organizationUnitSvc settingsSvc.OrganizationUnitServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
) StrategicIntiativeServiceInterface {
	return &StrategicIntiativeService{
		repo:                          repo,
		strategicInitiativeStatusRepo: strategicInitiativeStatusRepo,
		activityLogSvc:                activityLogSvc,
		organizationUnitSvc:           organizationUnitSvc,
		userSvc:                       userSvc,
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

	if strategicIntiative.OwnerDepartment.ID != 0 {
		strategicIntiativeDTO.OwnerDepartment = organizationDTOs.OrganizationUnit1DTO{
			ID:      int(strategicIntiative.OwnerDepartment.ID),
			TitleEn: strategicIntiative.OwnerDepartment.TitleEn,
			TitleAr: strategicIntiative.OwnerDepartment.TitleAr,
			Status:  strategicIntiative.OwnerDepartment.Status,
		}
	}
	if strategicIntiative.StrategicInitiativeStatus.ID != 0 {
		strategicIntiativeDTO.StrategicInitiativeStatus = dtos.StrategicInitiativeStatusDTO{
			ID:        strategicIntiative.StrategicInitiativeStatus.ID,
			TitleEn:   strategicIntiative.StrategicInitiativeStatus.TitleEn,
			TitleAr:   strategicIntiative.StrategicInitiativeStatus.TitleAr,
			ColorCode: strategicIntiative.StrategicInitiativeStatus.ColorCode,
			Type:      strategicIntiative.StrategicInitiativeStatus.Type,
			Status:    strategicIntiative.StrategicInitiativeStatus.Status,
		}
	}
	if strategicIntiative.ResponsibleUnit.ID != 0 {
		strategicIntiativeDTO.ResponsibleUnit = organizationDTOs.OrganizationUnit1DTO{
			ID:      int(strategicIntiative.ResponsibleUnit.ID),
			TitleEn: strategicIntiative.ResponsibleUnit.TitleEn,
			TitleAr: strategicIntiative.ResponsibleUnit.TitleAr,
			Status:  strategicIntiative.ResponsibleUnit.Status,
		}
	}

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

	if _, err := sis.organizationUnitSvc.GetByID(dto.ResponsibleUnitID); err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("responsible unit does not exist"), "not_exist", "responsible_unit")
	}

	ownerDepartmentID, err := sis.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("owner department is not set for the current user"), "not_exist", "owner_department")
	}

	// Parse dates
	dates, err := sis.parseStrategicIntiativeDates(dto)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	plannedStatus, err := sis.strategicInitiativeStatusRepo.GetByType("planned")
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
		ResponsibleUnitID: dto.ResponsibleUnitID,
		OwnerDepartmentID: ownerDepartmentID,

		ObjectivesEn:           dto.ObjectivesEn,
		ObjectivesAr:           dto.ObjectivesAr,
		RecommendationsEn:      dto.RecommendationsEn,
		RecommendationsAr:      dto.RecommendationsAr,
		KPIs:                   dto.KPIs,
		ExpectedBenefitsEn:     dto.ExpectedBenefitsEn,
		ExpectedBenefitsAr:     dto.ExpectedBenefitsAr,
		Budget:                 dto.Budget,
		ResourceRequirementsEn: dto.ResourceRequirementsEn,
		ResourceRequirementsAr: dto.ResourceRequirementsAr,
		DependenciesEn:         dto.DependenciesEn,
		DependenciesAr:         dto.DependenciesAr,

		StrategicInitiativeStatusID: plannedStatus.ID,
		CreatedBYID:                 currentUserID,
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
			"organization", "strategy_intiative", "", nil,
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

	if _, err := sis.organizationUnitSvc.GetByID(dto.ResponsibleUnitID); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible unit does not exist"), "not_exist", "responsible_unit")
	}

	ownerDepartmentID, err := sis.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("owner department is not set for the current user"), "not_exist", "owner_department")
	}

	// Parse dates
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
	existing.ResponsibleUnitID = dto.ResponsibleUnitID
	existing.OwnerDepartmentID = ownerDepartmentID

	existing.ObjectivesEn = dto.ObjectivesEn
	existing.ObjectivesAr = dto.ObjectivesAr
	existing.RecommendationsEn = dto.RecommendationsEn
	existing.RecommendationsAr = dto.RecommendationsAr
	existing.KPIs = dto.KPIs
	existing.ExpectedBenefitsEn = dto.ExpectedBenefitsEn
	existing.ExpectedBenefitsAr = dto.ExpectedBenefitsAr
	existing.Budget = dto.Budget
	existing.ResourceRequirementsEn = dto.ResourceRequirementsEn
	existing.ResourceRequirementsAr = dto.ResourceRequirementsAr
	existing.DependenciesEn = dto.DependenciesEn
	existing.DependenciesAr = dto.DependenciesAr

	existing.UpdatedBYID = &currentUserID

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
			"organization", "strategy_intiative", "", nil,
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

func (sis *StrategicIntiativeService) UpdateStatusAndCompletion(id uint, currentUserID uint, dto *dtos.StrategicIntiativeStatusCompletionUpdateDTO) error {
	if dto.CompletionPercentage < 0 || dto.CompletionPercentage > 100 {
		return utility.NewValidationError(fmt.Errorf("completion percentage must be between 0 and 100"), "invalid_range", "completion_percentage")
	}

	existing, err := sis.GetByID(id)
	if err != nil {
		return err
	}

	if _, err := sis.strategicInitiativeStatusRepo.GetByID(dto.StrategicInitiativeStatusID); err != nil {
		return err
	}

	cp := dto.CompletionPercentage
	existing.StrategicInitiativeStatusID = dto.StrategicInitiativeStatusID
	existing.CompletionPercentage = &cp
	existing.UpdatedBYID = &currentUserID
	now := time.Now()
	existing.UpdatedAt = &now

	err = sis.repo.Update(existing)
	if err != nil {
		return err
	}

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
			"organization", "strategy_intiative", "", nil,
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

	if parsed.EndDate.Before(*parsed.StartDate) {
		return nil, utility.NewValidationError(fmt.Errorf("end date must be after start date"), "invalid_date", "end_date")
	}

	return parsed, nil
}
