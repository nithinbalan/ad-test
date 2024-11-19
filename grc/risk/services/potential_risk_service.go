package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
	"os"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PotentialRiskServiceInterface interface {
	GetAll(params *riskDTOs.PotentialRiskListQueryParamsDTO) (*riskDTOs.PotentialRiskListResponseDTO, error)
	GetFullDetailsByID(id uint) (*riskDTOs.PotentialRiskDetailsDTO, error)
	GetByID(id uint) (*models.PotentialRisk, error)
	Create(currentUserID uint, potentialRiskFormDTO *riskDTOs.PotentialRiskFormDTO) (uint, error)
	Update(id uint, currentUserID uint, potentialRiskFormDTO *riskDTOs.PotentialRiskFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(id uint, currentUserID uint, potentialRiskStatusID uint) error

	ExportPotentialRisk(data []*riskDTOs.PotentialRiskListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type PotentialRiskService struct {
	repo                   riskRepositories.PotentialRiskRepositoryInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	severityLevelSvc       masterServices.SeverityLevelServiceInterface
	riskCategorySvc        RiskCategoryServiceInterface
	potentialRiskStatusSvc PotentialRiskStatusServiceInterface
	userSettingSvc         organizationServices.UserSettingServiceInterface
}

func NewPotentialRiskService(
	repo riskRepositories.PotentialRiskRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	riskCategorySvc RiskCategoryServiceInterface,
	potentialRiskStatusSvc PotentialRiskStatusServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface) PotentialRiskServiceInterface {
	return &PotentialRiskService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		severityLevelSvc:       severityLevelSvc,
		riskCategorySvc:        riskCategorySvc,
		potentialRiskStatusSvc: potentialRiskStatusSvc,
		userSettingSvc:         userSettingSvc,
	}
}

func (prs *PotentialRiskService) GetAll(params *riskDTOs.PotentialRiskListQueryParamsDTO) (*riskDTOs.PotentialRiskListResponseDTO, error) {
	return prs.repo.GetAll(params)
}

func (prs PotentialRiskService) GetByID(id uint) (*models.PotentialRisk, error) {
	potentialRisk, err := prs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return potentialRisk, nil
}

func (prs PotentialRiskService) GetFullDetailsByID(id uint) (*riskDTOs.PotentialRiskDetailsDTO, error) {

	var potentialRiskDTO riskDTOs.PotentialRiskDetailsDTO

	potentialRisk, err := prs.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	err = mapstructure.Decode(potentialRisk, &potentialRiskDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &potentialRiskDTO, nil
}

func (prs *PotentialRiskService) Create(currentUserID uint, potentialRiskFormDTO *riskDTOs.PotentialRiskFormDTO) (uint, error) {

	if err := prs.validatePotentialRiskFormDTO(potentialRiskFormDTO, nil); err != nil {
		return 0, err
	}

	potentialRiskStatus, err := prs.potentialRiskStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	potentialRisk := &models.PotentialRisk{
		RiskCategoryID:        potentialRiskFormDTO.RiskCategoryID,
		SeverityLevelID:       potentialRiskFormDTO.SeverityLevelID,
		PotentialRiskStatusID: potentialRiskStatus.ID,
		Title:                 potentialRiskFormDTO.Title,
		Description:           potentialRiskFormDTO.Description,
		RiskSource:            "manual",
		CreatedBYID:           currentUserID,
		UpdatedAt:             nil,
	}

	id, err := prs.repo.Create(potentialRisk)
	if err != nil {
		return 0, err
	}

	go func() {

		existingPotentialRisk, err := prs.GetFullDetailsByID(potentialRisk.ID)
		if err != nil {
			utility.LogError("Failed to fetch potential risk", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = prs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PotentialRisk{}.TableName(),
			id,
			existingPotentialRisk,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.PotentialRisk{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (prs *PotentialRiskService) Update(id uint, currentUserID uint, potentialRiskFormDTO *riskDTOs.PotentialRiskFormDTO) error {

	existingPotentialRisk, err := prs.GetByID(id)

	if err != nil {
		return err
	}

	oldPotentialRisk, err := prs.GetFullDetailsByID(id)
	if err != nil {
		return err
	}

	if err := prs.validatePotentialRiskFormDTO(potentialRiskFormDTO, &id); err != nil {
		return err
	}

	existingPotentialRisk.RiskCategoryID = potentialRiskFormDTO.RiskCategoryID
	existingPotentialRisk.SeverityLevelID = potentialRiskFormDTO.SeverityLevelID
	existingPotentialRisk.Title = potentialRiskFormDTO.Title
	existingPotentialRisk.Description = potentialRiskFormDTO.Description
	existingPotentialRisk.UpdatedBYID = &currentUserID

	err = prs.repo.Update(existingPotentialRisk)

	if err != nil {
		return err
	}

	go func() {
		existingRisk, err := prs.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"title":          oldPotentialRisk.Title,
			"description":    oldPotentialRisk.Description,
			"risk_cateogry":  oldPotentialRisk.RiskCategory.TitleEn,
			"severity_level": oldPotentialRisk.SeverityLevel.TitleEn,
		}

		newData := map[string]interface{}{
			"title":          existingRisk.Title,
			"description":    existingRisk.Description,
			"risk_cateogry":  existingRisk.RiskCategory.TitleEn,
			"severity_level": existingRisk.SeverityLevel.TitleEn,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err = prs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.PotentialRisk{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.PotentialRisk{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}

func (prs *PotentialRiskService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint

	var deletedPotentialRisk []map[string]interface{}

	for _, id := range ids {

		_, err := prs.repo.GetByID(id)

		if err != nil {
			continue
		}

		potentialRisk, err := prs.repo.GetFullDetailsByID(id)
		if err != nil {
			return nil, err
		}

		err = prs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedPotentialRisk = append(deletedPotentialRisk, map[string]interface{}{
			"id":    id,
			"title": potentialRisk.Title,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := prs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.PotentialRisk{}.TableName(),
				successfulIds,
				deletedPotentialRisk,
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

func (prs *PotentialRiskService) StatusUpdate(id uint, currentUserID uint, potentialRiskStatusID uint) error {
	existingPotentialRisk, err := prs.GetByID(id)

	if err != nil {
		return err
	}

	existingPotentialRisk.PotentialRiskStatusID = potentialRiskStatusID
	existingPotentialRisk.UpdatedBYID = &currentUserID

	err = prs.repo.Update(existingPotentialRisk)

	if err != nil {
		return err
	}

	return nil
}

func (prs *PotentialRiskService) validatePotentialRiskFormDTO(potentialRiskFormDTO *riskDTOs.PotentialRiskFormDTO, id *uint) error {
	if _, err := prs.riskCategorySvc.GetByID(potentialRiskFormDTO.RiskCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk category does not exist"), "not_exist", "risk category")
	}

	if _, err := prs.severityLevelSvc.GetByID(potentialRiskFormDTO.SeverityLevelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}

	if id == nil {
		if exists, err := prs.repo.ExistsByTitle(potentialRiskFormDTO.Title, potentialRiskFormDTO.RiskCategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	} else {
		if exists, err := prs.repo.ExistsByTitleExcludingID(*id, potentialRiskFormDTO.Title, potentialRiskFormDTO.RiskCategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}

func (prs *PotentialRiskService) ExportPotentialRisk(data []*riskDTOs.PotentialRiskListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := prs.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Title",
			"Description",
			"Risk Category",
			"Severity Level",
			"Status",
		}
	} else {
		headers = []string{
			"SL No.",
			"Title",
			"Description",
			"Risk Category",
			"Severity Level",
			"Status",
		}
	}
	// Prepare the data rows
	rows := make([][]interface{}, len(data))
	for i, risk := range data {
		var riskCategory, severityLevel, status string

		if userSetting.Language == "ar" {
			riskCategory = risk.RiskCategoryTitleAr
			severityLevel = risk.SeverityLevelTitleAr
			status = risk.PotentialRiskStatusTitleAr
		} else {
			riskCategory = risk.RiskCategoryTitleEn
			severityLevel = risk.SeverityLevelTitleEn
			status = risk.PotentialRiskStatusTitleEn
		}

		rows[i] = []interface{}{
			i + 1,
			risk.Title,
			utility.GetStringValue(risk.Description),
			riskCategory,
			severityLevel,
			status,
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}
