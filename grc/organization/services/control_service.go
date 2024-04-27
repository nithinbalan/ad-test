package services

import (
	"errors"
	"fmt"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationModels "jethur/grc/organization/models"
	organizationRepositories "jethur/grc/organization/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ControlServiceInterface interface {
	GetAll(params *organizationDTOs.ControlListQueryParamsDTO) (*organizationDTOs.ControlListResponseDTO, error)
	GetFullDetailsByID(id uint) (*organizationDTOs.ControlDetailDTO, error)
	Create(currentUserID uint, controlFormDTO *organizationDTOs.ControlFormDTO) (uint, error)
	Update(id uint, currentUserID uint, controlFormDTO *organizationDTOs.ControlFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	Archive(ids []uint, currentUserID uint) ([]uint, error)
	Restore(ids []uint, currentUserID uint) ([]uint, error)
}

type ControlService struct {
	repo                           organizationRepositories.ControlRepositoryInterface
	activityLogSvc                 logServices.ActivityLogServiceInterface
	controlTypeSvc                 ControlTypeServiceInterface
	controlCategorySvc             ControlCategoryServiceInterface
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface
	controlEfficiencyStatusSvc     ControlEfficiencyStatusServiceInterface
	userSvc                        UserServiceInterface
	reviewFrequencySvc             masterServices.ReviewFrequencyServiceInterface
	lifecycleStatusSvc             masterServices.LifecycleStatusServiceInterface
}

func NewControlService(
	repo organizationRepositories.ControlRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	controlTypeSvc ControlTypeServiceInterface,
	controlCategorySvc ControlCategoryServiceInterface,
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface,
	controlEfficiencyStatusSvc ControlEfficiencyStatusServiceInterface,
	userSvc UserServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface) ControlServiceInterface {
	return &ControlService{
		repo:                           repo,
		activityLogSvc:                 activityLogSvc,
		controlTypeSvc:                 controlTypeSvc,
		controlCategorySvc:             controlCategorySvc,
		controlImplementationStatusSvc: controlImplementationStatusSvc,
		controlEfficiencyStatusSvc:     controlEfficiencyStatusSvc,
		userSvc:                        userSvc,
		reviewFrequencySvc:             reviewFrequencySvc,
		lifecycleStatusSvc:             lifecycleStatusSvc,
	}
}

func (cs *ControlService) GetAll(params *organizationDTOs.ControlListQueryParamsDTO) (*organizationDTOs.ControlListResponseDTO, error) {
	return cs.repo.GetAll(params)
}

func (cs *ControlService) GetFullDetailsByID(id uint) (*organizationDTOs.ControlDetailDTO, error) {

	// Initialize your DTO
	var controlDetailDTO organizationDTOs.ControlDetailDTO

	control, err := cs.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(control, &controlDetailDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &controlDetailDTO, nil
}

func (cs *ControlService) Create(currentUserID uint, controlFormDTO *organizationDTOs.ControlFormDTO) (uint, error) {

	if err := cs.validateControlFormDTO(controlFormDTO); err != nil {
		return 0, err
	}

	dates, err := cs.parseControlDates(controlFormDTO)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := cs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	lifecycleStatus, err := cs.lifecycleStatusSvc.GetByType("active")
	if err != nil {
		return 0, err
	}

	control := &organizationModels.Control{
		ControlTypeID:                 controlFormDTO.ControlTypeID,
		ControlCategoryID:             controlFormDTO.ControlCategoryID,
		OwnerUnitID:                   controlFormDTO.OwnerUnitID,
		OwnerID:                       controlFormDTO.OwnerID,
		ControlImplementationStatusID: controlFormDTO.ControlImplementationStatusID,
		ControlEfficiencyStatusID:     controlFormDTO.ControlEfficiencyStatusID,
		ReferenceNumber:               referenceNumber,
		Code:                          controlFormDTO.Code,
		TitleEn:                       controlFormDTO.TitleEn,
		TitleAr:                       controlFormDTO.TitleAr,
		Description:                   controlFormDTO.Description,
		Objective:                     controlFormDTO.Objective,
		ReviewFrequencyID:             controlFormDTO.ReviewFrequencyID,
		LastReviewDate:                dates.LastReviewDate,
		TestResponsibleUserID:         controlFormDTO.TestResponsibleUserID,
		NextReviewDate:                dates.NextReviewDate,
		LastReviewNote:                controlFormDTO.LastReviewNote,
		TestFrequencyID:               controlFormDTO.TestFrequencyID,
		ImplementedAt:                 dates.ImplementedAt,
		LastTestDate:                  dates.LastTestDate,
		NextTestDate:                  dates.NextTestDate,
		LifecycleStatusID:             lifecycleStatus.ID,
		CreatedBy:                     currentUserID,
	}

	id, err := cs.repo.Create(control)
	if err != nil {
		return 0, err
	}

	go func() {
		err := cs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			organizationModels.Control{}.TableName(),
			id,
			controlFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": organizationModels.Control{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (cs *ControlService) Update(id uint, currentUserID uint, controlFormDTO *organizationDTOs.ControlFormDTO) error {
	existingControl, err := cs.repo.GetByID(id)

	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := cs.validateControlFormDTO(controlFormDTO); err != nil {
		return err
	}

	dates, err := cs.parseControlDates(controlFormDTO)
	if err != nil {
		return err
	}

	existingControl.ControlTypeID = controlFormDTO.ControlTypeID
	existingControl.ControlCategoryID = controlFormDTO.ControlCategoryID
	existingControl.OwnerUnitID = controlFormDTO.OwnerUnitID
	existingControl.OwnerID = controlFormDTO.OwnerID
	existingControl.ControlImplementationStatusID = controlFormDTO.ControlImplementationStatusID
	existingControl.ControlEfficiencyStatusID = controlFormDTO.ControlEfficiencyStatusID
	existingControl.Code = controlFormDTO.Code
	existingControl.TitleEn = controlFormDTO.TitleEn
	existingControl.TitleAr = controlFormDTO.TitleAr
	existingControl.Description = controlFormDTO.Description
	existingControl.Objective = controlFormDTO.Objective
	existingControl.ReviewFrequencyID = controlFormDTO.ReviewFrequencyID
	existingControl.LastReviewDate = dates.LastReviewDate
	existingControl.TestResponsibleUserID = controlFormDTO.TestResponsibleUserID
	existingControl.NextReviewDate = dates.NextReviewDate
	existingControl.LastReviewNote = controlFormDTO.LastReviewNote
	existingControl.TestFrequencyID = controlFormDTO.TestFrequencyID
	existingControl.ImplementedAt = dates.ImplementedAt
	existingControl.LastTestDate = dates.LastTestDate
	existingControl.NextTestDate = dates.NextTestDate
	existingControl.UpdatedBy = &currentUserID

	err = cs.repo.Update(existingControl)

	if err != nil {
		return err
	}

	return nil
}

func (cs *ControlService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := cs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = cs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := cs.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
				organizationModels.Control{}.TableName(),
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

func (cs *ControlService) Archive(ids []uint, currentUserID uint) ([]uint, error) {
	lifecycleStatus, err := cs.lifecycleStatusSvc.GetByType("archived")
	if err != nil {
		return nil, err
	}
	idsToUpdate, err := cs.repo.GetIDsByLifecycleStatus(ids, "active")
	if err != nil {
		return nil, err
	}

	successfulIds, err := cs.repo.UpdateLifecycleStatusForIDs(idsToUpdate, lifecycleStatus.ID)

	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := cs.activityLogSvc.LogBulkActivities(
				logModels.ArchiveAction,
				organizationModels.Control{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk archive activity: %v\n", err)
			}
		}()
	}

	return successfulIds, nil
}

func (cs *ControlService) Restore(ids []uint, currentUserID uint) ([]uint, error) {
	lifecycleStatus, err := cs.lifecycleStatusSvc.GetByType("active")
	if err != nil {
		return nil, err
	}
	idsToUpdate, err := cs.repo.GetIDsByLifecycleStatus(ids, "archived")
	if err != nil {
		return nil, err
	}

	successfulIds, err := cs.repo.UpdateLifecycleStatusForIDs(idsToUpdate, lifecycleStatus.ID)

	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := cs.activityLogSvc.LogBulkActivities(
				logModels.RestoreAction,
				organizationModels.Control{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk archive activity: %v\n", err)
			}
		}()
	}

	return successfulIds, nil
}

func (cs *ControlService) validateControlFormDTO(controlFormDTO *organizationDTOs.ControlFormDTO) error {
	if _, err := cs.controlTypeSvc.GetByID(uint(controlFormDTO.ControlTypeID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("control type does not exist"), "not_exist", "control type")
	}

	if _, err := cs.controlCategorySvc.GetByID(uint(controlFormDTO.ControlCategoryID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("control category does not exist"), "not_exist", "control category")
	}

	if controlFormDTO.OwnerID != nil {
		if _, err := cs.userSvc.GetByID(uint(*controlFormDTO.OwnerID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("control owner does not exist"), "not_exist", "control owner")
		}
	}

	if controlFormDTO.ControlImplementationStatusID != nil {
		if _, err := cs.controlImplementationStatusSvc.GetByID(uint(*controlFormDTO.ControlImplementationStatusID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("control implementation status does not exist"), "not_exist", "control implementation status")
		}
	}

	if controlFormDTO.ControlEfficiencyStatusID != nil {
		if _, err := cs.controlImplementationStatusSvc.GetByID(uint(*controlFormDTO.ControlEfficiencyStatusID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("control efficiency status does not exist"), "not_exist", "control efficiency status")
		}
	}

	if controlFormDTO.ReviewFrequencyID != nil {
		if _, err := cs.reviewFrequencySvc.GetByID(uint(*controlFormDTO.ReviewFrequencyID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("review frequiency does not exist"), "not_exist", "review frequiency")
		}
	}

	if _, err := cs.userSvc.GetByID(uint(*&controlFormDTO.TestResponsibleUserID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("control test responsible user does not exist"), "not_exist", "test responsible user")
	}

	if controlFormDTO.TestFrequencyID != nil {
		if _, err := cs.reviewFrequencySvc.GetByID(uint(*controlFormDTO.TestFrequencyID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("review test frequiency does not exist"), "not_exist", "test frequiency")
		}
	}

	return nil
}

func (cs *ControlService) parseControlDates(dto *organizationDTOs.ControlFormDTO) (*organizationDTOs.ControlParsedDates, error) {
	parsed := &organizationDTOs.ControlParsedDates{}

	var err error
	if dto.LastReviewDate != nil {
		parsed.LastReviewDate, err = utility.ParseDate(*dto.LastReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid last review date"), "invalid_date", "last review date")
		}
	}

	if dto.NextReviewDate != nil {
		parsed.NextReviewDate, err = utility.ParseDate(*dto.NextReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next review date"), "invalid_date", "next review date")
		}
	}

	if dto.LastTestDate != nil {
		parsed.LastTestDate, err = utility.ParseDate(*dto.LastTestDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid last test date"), "invalid_date", "last test date")
		}
	}

	if dto.NextTestDate != nil {
		parsed.NextTestDate, err = utility.ParseDate(*dto.NextTestDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next test date"), "invalid_date", "next test date")
		}
	}

	if dto.ImplementedAt != nil {
		parsed.ImplementedAt, err = utility.ParseDate(*dto.ImplementedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid implemented at date"), "invalid_date", "implemented at date")
		}
	}

	return parsed, nil
}
