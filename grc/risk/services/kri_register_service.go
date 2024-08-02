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

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type KRIRegisterServiceInterface interface {
	GetAll(params *riskDTOs.KRIRegisterListQueryParamsDTO) (*riskDTOs.KRIRegisterListResponseDTO, error)
	GetFullDetailsByID(id uint) (*riskDTOs.KRIRegisterDetailsDTO, error)
	GetByID(id uint) (*models.KRIRegister, error)
	Create(currentUserID uint, kriRegisterFormDTO *riskDTOs.KRIRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, kriRegisterFormDTO *riskDTOs.KRIRegisterFormDTO) error
	UpdateValue(kriRegister *models.KRIRegister, tx utility.DatabaseTransaction) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type KRIRegisterService struct {
	repo               riskRepositories.KRIRegisterRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	kriCategorySvc     KRICategoryServiceInterface
	userSvc            organizationServices.UserServiceInterface
	riskLevelSvc       RiskLevelServiceInterface
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface
	unitSvc            masterServices.UnitServiceInterface
}

func NewKRIRegisterService(
	repo riskRepositories.KRIRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	kriCategorySvc KRICategoryServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	riskLevelSvc RiskLevelServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	unitSvc masterServices.UnitServiceInterface,
) KRIRegisterServiceInterface {
	return &KRIRegisterService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		kriCategorySvc:     kriCategorySvc,
		userSvc:            userSvc,
		riskLevelSvc:       riskLevelSvc,
		reviewFrequencySvc: reviewFrequencySvc,
		unitSvc:            unitSvc,
	}
}

func (krs *KRIRegisterService) GetAll(params *riskDTOs.KRIRegisterListQueryParamsDTO) (*riskDTOs.KRIRegisterListResponseDTO, error) {
	return krs.repo.GetAll(params)
}

func (krs KRIRegisterService) GetByID(id uint) (*models.KRIRegister, error) {
	kriRegister, err := krs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return kriRegister, nil
}

func (krs KRIRegisterService) GetFullDetailsByID(id uint) (*riskDTOs.KRIRegisterDetailsDTO, error) {
	var kriRegisterDTO riskDTOs.KRIRegisterDetailsDTO

	kriRegister, err := krs.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(kriRegister, &kriRegisterDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &kriRegisterDTO, nil
}

func (krs *KRIRegisterService) Create(currentUserID uint, kriRegisterFormDTO *riskDTOs.KRIRegisterFormDTO) (uint, error) {
	if err := krs.validateKRIRegisterFormDTO(kriRegisterFormDTO, nil); err != nil {
		return 0, err
	}
	dates, err := krs.parseKRIRegisterDates(kriRegisterFormDTO)
	if err != nil {
		return 0, err
	}
	referenceNumber, err := krs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	var organizationUnitID uint
	unitID, err := krs.userSvc.GetUserPrimaryUnitID(kriRegisterFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("responsible unit not set forthe responsible user"), "not_exist", "responsible unit")
	}
	organizationUnitID = unitID
	kriRegister := &models.KRIRegister{
		ReferenceNumber:   referenceNumber,
		KRICategoryID:     kriRegisterFormDTO.KRICategoryID,
		ResponsibleUnitID: organizationUnitID,
		ResponsibleUserID: kriRegisterFormDTO.ResponsibleUserID,
		ReviewFrequencyID: kriRegisterFormDTO.ReviewFrequencyID,
		RiskLevelID:       kriRegisterFormDTO.RiskLevelID,
		UnitID:            kriRegisterFormDTO.UnitID,
		Title:             kriRegisterFormDTO.Title,
		Description:       kriRegisterFormDTO.Description,
		CalculationType:   kriRegisterFormDTO.CalculationType,
		LowThreshold:      kriRegisterFormDTO.LowThreshold,
		MediumThreshold:   kriRegisterFormDTO.MediumThreshold,
		HighThreshold:     kriRegisterFormDTO.HighThreshold,
		CurrentValue:      kriRegisterFormDTO.CurrentValue,
		NextReviewDate:    *dates.NextReviewDate,
		LastReviewNote:    kriRegisterFormDTO.LastReviewNote,
		CreatedBYID:       currentUserID,
		UpdatedAt:         nil,
	}
	if kriRegisterFormDTO.LastReviewedDate != "" {
		kriRegister.LastReviewedDate = dates.LastReviewedDate
	}

	id, err := krs.repo.Create(kriRegister)
	if err != nil {
		return 0, err
	}
	go func() {
		err := krs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.KRIRegister{}.TableName(),
			id,
			kriRegisterFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.KRIRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (krs *KRIRegisterService) Update(id uint, currentUserID uint, kriRegisterFormDTO *riskDTOs.KRIRegisterFormDTO) error {
	existingKRIRegister, err := krs.GetByID(id)

	if err != nil {
		return err
	}

	if err := krs.validateKRIRegisterFormDTO(kriRegisterFormDTO, &id); err != nil {
		return err
	}

	dates, err := krs.parseKRIRegisterDates(kriRegisterFormDTO)
	if err != nil {
		return err
	}
	if existingKRIRegister.ResponsibleUserID != kriRegisterFormDTO.ResponsibleUserID {
		organizationUnitID, err := krs.userSvc.GetUserPrimaryUnitID(uint(kriRegisterFormDTO.ResponsibleUserID))
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("responsible unit not set for the responsible user"), "not_exist", "responsible unit")
		}
		existingKRIRegister.ResponsibleUnitID = organizationUnitID
	}
	existingKRIRegister.KRICategoryID = kriRegisterFormDTO.KRICategoryID
	existingKRIRegister.ResponsibleUserID = kriRegisterFormDTO.ResponsibleUserID
	existingKRIRegister.ReviewFrequencyID = kriRegisterFormDTO.ReviewFrequencyID
	existingKRIRegister.RiskLevelID = kriRegisterFormDTO.RiskLevelID
	existingKRIRegister.UnitID = kriRegisterFormDTO.UnitID
	existingKRIRegister.Title = kriRegisterFormDTO.Title
	existingKRIRegister.Description = kriRegisterFormDTO.Description
	existingKRIRegister.CalculationType = kriRegisterFormDTO.CalculationType
	existingKRIRegister.LowThreshold = kriRegisterFormDTO.LowThreshold
	existingKRIRegister.MediumThreshold = kriRegisterFormDTO.MediumThreshold
	existingKRIRegister.HighThreshold = kriRegisterFormDTO.HighThreshold
	existingKRIRegister.CurrentValue = kriRegisterFormDTO.CurrentValue

	if kriRegisterFormDTO.LastReviewedDate != "" {
		existingKRIRegister.LastReviewedDate = dates.LastReviewedDate
	}
	existingKRIRegister.NextReviewDate = *dates.NextReviewDate
	existingKRIRegister.LastReviewNote = kriRegisterFormDTO.LastReviewNote
	existingKRIRegister.UpdatedBYID = &currentUserID

	err = krs.repo.Update(existingKRIRegister)
	if err != nil {
		return err
	}

	return nil
}

func (krs *KRIRegisterService) UpdateValue(kriRegister *models.KRIRegister, tx utility.DatabaseTransaction) error {
	err := krs.repo.Update(kriRegister, tx)
	if err != nil {
		return err
	}

	return nil
}

func (krs *KRIRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := krs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = krs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := krs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.KRIRegister{}.TableName(),
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

func (krs *KRIRegisterService) validateKRIRegisterFormDTO(kriRegisterFormDTO *riskDTOs.KRIRegisterFormDTO, id *uint) error {
	if _, err := krs.kriCategorySvc.GetByID(kriRegisterFormDTO.KRICategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("kri category does not exist"), "not_exist", "kri category")
	}

	if _, err := krs.userSvc.GetByID(kriRegisterFormDTO.ResponsibleUserID); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}
	if kriRegisterFormDTO.RiskLevelID != nil {
		if _, err := krs.riskLevelSvc.GetByID(*kriRegisterFormDTO.RiskLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk level does not exist"), "not_exist", "risk level")
		}
	}
	if _, err := krs.riskLevelSvc.GetByID(kriRegisterFormDTO.UnitID); err != nil {
		return utility.NewValidationError(fmt.Errorf("unit does not exist"), "not_exist", "unit")
	}

	if _, err := krs.reviewFrequencySvc.GetByID(kriRegisterFormDTO.ReviewFrequencyID); err != nil {
		return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
	}

	if id == nil {
		if exists, err := krs.repo.ExistsByTitle(kriRegisterFormDTO.Title, kriRegisterFormDTO.KRICategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	} else {
		if exists, err := krs.repo.ExistsByTitleExcludingID(*id, kriRegisterFormDTO.Title, kriRegisterFormDTO.KRICategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}

func (krs *KRIRegisterService) parseKRIRegisterDates(dto *riskDTOs.KRIRegisterFormDTO) (*riskDTOs.KRIRegisterParsedDates, error) {
	parsed := &riskDTOs.KRIRegisterParsedDates{}

	var err error
	if dto.LastReviewedDate != "" {
		parsed.LastReviewedDate, err = utility.ParseDate(dto.LastReviewedDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid last review date"), "invalid_date", "last review date")
		}
	}

	if dto.NextReviewDate != "" {
		parsed.NextReviewDate, err = utility.ParseDate(dto.NextReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next review date"), "invalid_date", "next review date")
		}
	}

	return parsed, nil
}
