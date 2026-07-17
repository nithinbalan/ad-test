package services

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"jethur/grc/logs/services"
	"jethur/grc/models"
	settingsSvc "jethur/grc/settings/services"
	"jethur/grc/strategy/dtos"
	"jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/html"
	"gorm.io/gorm"

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
	ExportStrategicInitiatives(data []*dtos.StrategicIntiativeListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type StrategicIntiativeService struct {
	repo                                repositories.StrategicIntiativeRepositoryInterface
	strategicInitiativeStatusUpdateRepo repositories.StrategyInitiativeStatusUpdateRepositoryInterface
	strategicInitiativeStatusRepo       repositories.StrategicInitiativeStatusRepositoryInterface
	activityLogSvc                      services.ActivityLogServiceInterface
	organizationUnitSvc                 settingsSvc.OrganizationUnitServiceInterface
	userSvc                             organizationSvc.UserServiceInterface
	userSettingSvc                      organizationSvc.UserSettingServiceInterface
}

func NewStrategicIntiativeService(repo repositories.StrategicIntiativeRepositoryInterface,
	strategicInitiativeStatusUpdateRepo repositories.StrategyInitiativeStatusUpdateRepositoryInterface,
	strategicInitiativeStatusRepo repositories.StrategicInitiativeStatusRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	organizationUnitSvc settingsSvc.OrganizationUnitServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	userSettingSvc organizationSvc.UserSettingServiceInterface,
) StrategicIntiativeServiceInterface {
	return &StrategicIntiativeService{
		repo:                                repo,
		strategicInitiativeStatusUpdateRepo: strategicInitiativeStatusUpdateRepo,
		strategicInitiativeStatusRepo:       strategicInitiativeStatusRepo,
		activityLogSvc:                      activityLogSvc,
		organizationUnitSvc:                 organizationUnitSvc,
		userSvc:                             userSvc,
		userSettingSvc:                      userSettingSvc,
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

	statusUpdate := &models.StrategicInitiativeStatusUpdate{
		StrategyInitiativeID:        id,
		StrategicInitiativeStatusID: plannedStatus.ID,
		Comment:                     "Initial strategic initiative status update",
		CreatedByID:                 currentUserID,
		Percentage:                  0,
	}

	if _, err := sis.strategicInitiativeStatusUpdateRepo.Create(statusUpdate, tx); err != nil {
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
			"organization", "strategy_intiative", "Created Strategic Initiative", nil,
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
			"organization", "strategy_intiative", "Updated Strategic Initiative", nil,
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

func (sis *StrategicIntiativeService) ExportStrategicInitiatives(data []*dtos.StrategicIntiativeListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := sis.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"Title",
		"Description",
		"Objectives",
		"Recommendations",
		"Owner Department",
		"KPIs",
		"Expected Benefits",
		"Budget",
		"Resource Requirements",
		"Dependencies",
		"Responsible Unit",
		"Status",
		"Responsible User",
		"Start Date",
		"End Date",
	}

	rows := make([][]interface{}, len(data))
	for i, initiative := range data {
		budget := ""
		if initiative.Budget != nil {
			budget = fmt.Sprintf("%.2f", *initiative.Budget)
		}

		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				initiative.TitleAr,
				stripHTML(initiative.DescriptionAr),
				stripHTML(initiative.ObjectivesAr),
				stripHTML(initiative.RecommendationsAr),
				utility.GetStringValue(initiative.OwnerDepartmentTitleAr),
				stripHTML(initiative.KPIs),
				stripHTML(initiative.ExpectedBenefitsAr),
				budget,
				stripHTML(initiative.ResourceRequirementsAr),
				stripHTML(initiative.DependenciesAr),
				utility.GetStringValue(initiative.ResponsibleUnitTitleAr),
				utility.GetStringValue(initiative.StrategicInitiativeStatusTitleAr),
				fmt.Sprintf("%s %s", utility.GetStringValue(initiative.ResponsibleUserFirstNameAr), utility.GetStringValue(initiative.ResponsibleUserLastNameAr)),
				formatDateOnly(initiative.StartDate),
				formatDateOnly(initiative.EndDate),
			}
		} else {
			rows[i] = []interface{}{
				initiative.TitleEn,
				stripHTML(initiative.DescriptionEn),
				stripHTML(initiative.ObjectivesEn),
				stripHTML(initiative.RecommendationsEn),
				utility.GetStringValue(initiative.OwnerDepartmentTitleEn),
				stripHTML(initiative.KPIs),
				stripHTML(initiative.ExpectedBenefitsEn),
				budget,
				stripHTML(initiative.ResourceRequirementsEn),
				stripHTML(initiative.DependenciesEn),
				utility.GetStringValue(initiative.ResponsibleUnitTitleEn),
				utility.GetStringValue(initiative.StrategicInitiativeStatusTitleEn),
				fmt.Sprintf("%s %s", utility.GetStringValue(initiative.ResponsibleUserFirstNameEn), utility.GetStringValue(initiative.ResponsibleUserLastNameEn)),
				formatDateOnly(initiative.StartDate),
				formatDateOnly(initiative.EndDate),
			}
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func stripHTML(content *string) string {
	if content == nil {
		return ""
	}

	document, err := html.Parse(strings.NewReader(*content))
	if err != nil {
		return *content
	}

	var text strings.Builder
	var collectText func(*html.Node)
	collectText = func(node *html.Node) {
		if node.Type == html.TextNode {
			text.WriteString(node.Data)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			collectText(child)
		}
	}

	collectText(document)
	return strings.TrimSpace(text.String())
}

func formatDateOnly(date *string) string {
	if date == nil || strings.TrimSpace(*date) == "" {
		return ""
	}

	parsedDate, err := time.Parse("2006-01-02T15:04:05Z07:00", *date)
	if err != nil {
		parsedDate, err = time.Parse("2006-01-02 15:04:05", *date)
		if err != nil {
			return *date
		}
	}

	return parsedDate.Format("2006-01-02")
}
