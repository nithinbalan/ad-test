package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	settingServices "jethur/grc/settings/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	kpiRegisterRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type KPIRegisterServiceInterface interface {
	GetAll(params *strategyDTOs.KPIRegisterListQueryParamsDTO) ([]*strategyDTOs.KPIRegisterListDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.KPIRegisterDetailDTO, error)
	GetByID(id uint) (*models.KPIRegister, error)
	Create(currentUserID uint, strategyDTOs *strategyDTOs.KPIRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, kpiRegister *strategyDTOs.KPIRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type KPIRegisterService struct {
	repo                kpiRegisterRepositories.KPIRegisterRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	userSvc             organizationServices.UserServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
	organizationUnitSvc settingServices.OrganizationUnitServiceInterface
	kpiStatusSvc        KPIStatusServiceInterface
}

func NewKPIRegisterService(
	repo kpiRegisterRepositories.KPIRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	organizationUnitSvc settingServices.OrganizationUnitServiceInterface,
	kpiStatusSvc KPIStatusServiceInterface,
) KPIRegisterServiceInterface {
	return &KPIRegisterService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		userSvc:             userSvc,
		userSettingSvc:      userSettingSvc,
		organizationUnitSvc: organizationUnitSvc,
		kpiStatusSvc:        kpiStatusSvc,
	}
}

func (krs *KPIRegisterService) GetAll(params *strategyDTOs.KPIRegisterListQueryParamsDTO) ([]*strategyDTOs.KPIRegisterListDTO, error) {
	return krs.repo.GetAll(params)
}

func (krs KPIRegisterService) GetByID(id uint) (*models.KPIRegister, error) {
	kpiRegister, err := krs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return kpiRegister, nil
}

func (krs KPIRegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.KPIRegisterDetailDTO, error) {
	// Initialize your DTO
	var kpiRegisterDTO strategyDTOs.KPIRegisterDetailDTO

	kpiRegister, err := krs.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(kpiRegister, &kpiRegisterDTO)
	if err != nil {
		return nil, err
	}

	kpiRegisterDTO.CreatedAt = kpiRegister.CreatedAt
	kpiRegisterDTO.UpdatedAt = kpiRegister.UpdatedAt

	return &kpiRegisterDTO, nil
}

func (krs *KPIRegisterService) Create(currentUserID uint, kpiRegisterFormDTO *strategyDTOs.KPIRegisterFormDTO) (uint, error) {

	tx := krs.repo.BeginTransaction()

	dates, err := krs.parseKPIRegisterDates(kpiRegisterFormDTO)
	if err != nil {
		return 0, err
	}

	kpiStatus, err := krs.kpiStatusSvc.GetByType("active")
	if err != nil {
		return 0, err
	}

	kpiRegister := &models.KPIRegister{
		Title:             kpiRegisterFormDTO.Title,
		Description:       kpiRegisterFormDTO.Description,
		KPICategoryID:     kpiRegisterFormDTO.KPICategoryID,
		ResponsibleUserID: kpiRegisterFormDTO.ResponsibleUserID,
		ReviewFrequencyID: kpiRegisterFormDTO.ReviewFrequencyID,
		CalculationTypeID: kpiRegisterFormDTO.CalculationTypeID,
		Formula:           kpiRegisterFormDTO.Formula,
		UnitID:            kpiRegisterFormDTO.UnitID,
		Source:            kpiRegisterFormDTO.Source,
		TargetValue:       kpiRegisterFormDTO.TargetValue,
		Weightage:         kpiRegisterFormDTO.Weightage,
		StartDate:         dates.StartDate,
		EndDate:           dates.EndDate,
		Remarks:           kpiRegisterFormDTO.Remarks,
		KPIStatusID:       kpiStatus.ID,
		CreatedBYID:       currentUserID,
	}
	id, err := krs.repo.Create(kpiRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Prepare input data records
	var inputs []*models.KPIRegisterDataInput
	for _, input := range kpiRegisterFormDTO.DataInputs {
		if input != nil {
			inputData := &models.KPIRegisterDataInput{
				KPIRegisterID: id,
				Input:         input.Input,
				Variable:      input.Variable,
				CreatedBYID:   currentUserID,
			}
			inputs = append(inputs, inputData)
		}
	}

	if len(inputs) > 0 {
		if err := krs.repo.CreateKPIRegisterInpuData(inputs, tx); err != nil {
			return 0, err
		}
	}

	tx.Commit()
	go func() {
		existingKPIRegister, err := krs.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch kpi register", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = krs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.KPIRegister{}.TableName(),
			id,
			existingKPIRegister,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.KPIRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (krs *KPIRegisterService) Update(id uint, currentUserID uint, kpiRegisterFormDTO *strategyDTOs.KPIRegisterFormDTO) error {
	tx := krs.repo.BeginTransaction()

	existingKPIRegister, err := krs.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldKPIRegister, err := krs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	dates, err := krs.parseKPIRegisterDates(kpiRegisterFormDTO)
	if err != nil {
		return err
	}

	existingKPIRegister.Title = kpiRegisterFormDTO.Title
	existingKPIRegister.Description = kpiRegisterFormDTO.Description
	existingKPIRegister.KPICategoryID = kpiRegisterFormDTO.KPICategoryID
	existingKPIRegister.ResponsibleUserID = kpiRegisterFormDTO.ResponsibleUserID
	existingKPIRegister.ReviewFrequencyID = kpiRegisterFormDTO.ReviewFrequencyID
	existingKPIRegister.CalculationTypeID = kpiRegisterFormDTO.CalculationTypeID
	existingKPIRegister.Formula = kpiRegisterFormDTO.Formula
	existingKPIRegister.UnitID = kpiRegisterFormDTO.UnitID
	existingKPIRegister.Source = kpiRegisterFormDTO.Source
	existingKPIRegister.TargetValue = kpiRegisterFormDTO.TargetValue
	existingKPIRegister.LowThreshold = kpiRegisterFormDTO.LowThreshold
	existingKPIRegister.MediumThreshold = kpiRegisterFormDTO.MediumThreshold
	existingKPIRegister.Weightage = kpiRegisterFormDTO.Weightage
	existingKPIRegister.StartDate = dates.StartDate
	existingKPIRegister.EndDate = dates.EndDate
	existingKPIRegister.Remarks = kpiRegisterFormDTO.Remarks
	existingKPIRegister.UpdatedBYID = &currentUserID

	err = krs.repo.Update(existingKPIRegister, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := krs.repo.DeleteKPIRegisterInputData(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Prepare input data records
	var inputs []*models.KPIRegisterDataInput
	for _, input := range kpiRegisterFormDTO.DataInputs {
		if input != nil {
			inputData := &models.KPIRegisterDataInput{
				KPIRegisterID: id,
				Input:         input.Input,
				Variable:      input.Variable,
				CreatedBYID:   currentUserID,
			}
			inputs = append(inputs, inputData)
		}
	}

	if len(inputs) > 0 {
		if err := krs.repo.CreateKPIRegisterInpuData(inputs, tx); err != nil {
			return err
		}
	}

	tx.Commit()
	go func() {
		existingKPIRegister, err := krs.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"title":            oldKPIRegister.Title,
			"description":      oldKPIRegister.Description,
			"kpi_category":     oldKPIRegister.KPICategory.TitleEn,
			"responsible_user": utility.GetFieldValue(oldKPIRegister.ResponsibleUser),
			"review_frequency": utility.GetFieldValue(oldKPIRegister.ReviewFrequency),
			"calculation_type": utility.GetFieldValue(oldKPIRegister.CalculationType),
			"formula":          oldKPIRegister.Formula,
			"unit":             utility.GetFieldValue(oldKPIRegister.Unit),
			"source":           oldKPIRegister.Source,
			"target_value":     oldKPIRegister.TargetValue,
			"low_threshold":    oldKPIRegister.LowThreshold,
			"medium_threshold": oldKPIRegister.MediumThreshold,
			"weightage":        oldKPIRegister.Weightage,
			"start_date":       oldKPIRegister.StartDateString,
			"end_date":         oldKPIRegister.EndDateString,
			"remarks":          oldKPIRegister.Remarks,
		}

		newData := map[string]interface{}{
			"title":            existingKPIRegister.Title,
			"description":      existingKPIRegister.Description,
			"kpi_category":     existingKPIRegister.KPICategory.TitleEn,
			"responsible_user": utility.GetFieldValue(existingKPIRegister.ResponsibleUser),
			"review_frequency": utility.GetFieldValue(existingKPIRegister.ReviewFrequency),
			"calculation_type": utility.GetFieldValue(existingKPIRegister.CalculationType),
			"formula":          existingKPIRegister.Formula,
			"unit":             utility.GetFieldValue(existingKPIRegister.Unit),
			"source":           existingKPIRegister.Source,
			"target_value":     existingKPIRegister.TargetValue,
			"low_threshold":    existingKPIRegister.LowThreshold,
			"medium_threshold": existingKPIRegister.MediumThreshold,
			"weightage":        existingKPIRegister.Weightage,
			"start_date":       existingKPIRegister.StartDateString,
			"end_date":         existingKPIRegister.EndDateString,
			"remarks":          existingKPIRegister.Remarks,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := krs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.KPIRegister{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.KPIRegister{}.TableName(),
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

func (krs *KPIRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := krs.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedKPIRegister []map[string]interface{}

	for _, id := range ids {
		kpiRegister, err := krs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedKPIRegister = append(deletedKPIRegister, map[string]interface{}{
			"id":    id,
			"title": kpiRegister.Title,
		})

		// Delete the business
		if err := krs.repo.Delete(id, tx); err != nil {
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

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := krs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.KPIRegister{}.TableName(),
				successfulIds,
				deletedKPIRegister,
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

func (krs *KPIRegisterService) parseKPIRegisterDates(dto *strategyDTOs.KPIRegisterFormDTO) (*strategyDTOs.KPIRegisterParsedDates, error) {
	parsed := &strategyDTOs.KPIRegisterParsedDates{}

	var err error

	if dto.StartDate != nil {
		parsed.StartDate, err = utility.ParseDate(*dto.StartDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
		}
	}

	if dto.EndDate != nil {
		parsed.EndDate, err = utility.ParseDate(*dto.EndDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid end date"), "invalid_date", "end date")
		}
	}

	return parsed, nil
}
