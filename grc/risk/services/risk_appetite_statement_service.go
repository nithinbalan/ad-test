package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskAppetiteStatementServiceInterface interface {
	GetAll(params *riskDTOs.RiskAppetiteStatementListQueryParamsDTO) (*riskDTOs.RiskAppetiteStatementListResponseDTO, error)
	GetAllByCategories(params *riskDTOs.RiskAppetiteStatementListQueryParamsDTO) (*riskDTOs.RiskAppetiteStatementListGroupedResponseDTO, error)
	GetByID(id uint) (*models.RiskAppetiteStatement, error)
	GetFullDetailsByID(id uint) (*riskDTOs.RiskAppetiteStatementDetailsDTO, error)
	Create(currentUserID uint, RiskAppetiteStatementFormDTO *riskDTOs.RiskAppetiteStatementFormDTO) (uint, error)
	Update(id uint, currentUserID uint, RiskAppetiteStatementFormDTO *riskDTOs.RiskAppetiteStatementFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type RiskAppetiteStatementService struct {
	repo            repositories.RiskAppetiteStatementRepositoryInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	riskCategorySvc RiskCategoryServiceInterface
}

func NewRiskAppetiteStatementService(
	repo repositories.RiskAppetiteStatementRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskCategorySvc RiskCategoryServiceInterface) RiskAppetiteStatementServiceInterface {
	return &RiskAppetiteStatementService{
		repo:            repo,
		activityLogSvc:  activityLogSvc,
		riskCategorySvc: riskCategorySvc,
	}
}

func (rass *RiskAppetiteStatementService) GetAll(params *riskDTOs.RiskAppetiteStatementListQueryParamsDTO) (*riskDTOs.RiskAppetiteStatementListResponseDTO, error) {
	return rass.repo.GetAll(params)
}

func (rass *RiskAppetiteStatementService) GetAllByCategories(params *riskDTOs.RiskAppetiteStatementListQueryParamsDTO) (*riskDTOs.RiskAppetiteStatementListGroupedResponseDTO, error) {
	return rass.repo.GetAllByCategories(params)
}

func (rass RiskAppetiteStatementService) GetByID(id uint) (*models.RiskAppetiteStatement, error) {
	riskAppetiteStatement, err := rass.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskAppetiteStatement, nil
}

func (rass RiskAppetiteStatementService) GetFullDetailsByID(id uint) (*riskDTOs.RiskAppetiteStatementDetailsDTO, error) {

	var riskAppetiteDetailsDTO riskDTOs.RiskAppetiteStatementDetailsDTO

	riskAppetiteDetails, err := rass.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(riskAppetiteDetails, &riskAppetiteDetailsDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &riskAppetiteDetailsDTO, nil
}

func (rass *RiskAppetiteStatementService) Create(currentUserID uint, riskAppetiteStatementFormDTO *riskDTOs.RiskAppetiteStatementFormDTO) (uint, error) {

	if err := rass.validateRiskAppetiteStatementFormDTO(riskAppetiteStatementFormDTO, nil); err != nil {
		return 0, err
	}

	referenceNumber, err := rass.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	riskAppetiteStatement := &models.RiskAppetiteStatement{
		ReferenceNumber: referenceNumber,
		RiskCategoryID:  riskAppetiteStatementFormDTO.RiskCategoryID,
		StatementEn:     riskAppetiteStatementFormDTO.StatementEn,
		StatementAr:     riskAppetiteStatementFormDTO.StatementAr,
		Status:          riskAppetiteStatementFormDTO.Status,
		CreatedByID:     currentUserID,
		UpdatedAt:       nil,
	}

	id, err := rass.repo.Create(riskAppetiteStatement)
	if err != nil {
		return 0, err
	}

	go func() {
		err := rass.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskAppetiteStatement{}.TableName(),
			id,
			riskAppetiteStatementFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskAppetiteStatement{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (rass *RiskAppetiteStatementService) Update(id uint, currentUserID uint, riskAppetiteStatementFormDTO *riskDTOs.RiskAppetiteStatementFormDTO) error {

	existingStatement, err := rass.GetByID(id)

	if err != nil {
		return err
	}

	if err := rass.validateRiskAppetiteStatementFormDTO(riskAppetiteStatementFormDTO, &id); err != nil {
		return err
	}

	existingStatement.RiskCategoryID = riskAppetiteStatementFormDTO.RiskCategoryID
	existingStatement.StatementEn = riskAppetiteStatementFormDTO.StatementEn
	existingStatement.StatementAr = riskAppetiteStatementFormDTO.StatementAr
	existingStatement.Status = riskAppetiteStatementFormDTO.Status
	existingStatement.UpdatedByID = &currentUserID

	err = rass.repo.Update(existingStatement)

	if err != nil {
		return err
	}

	return nil
}

func (rass *RiskAppetiteStatementService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := rass.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = rass.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := rass.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.RiskAppetiteStatement{}.TableName(),
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

func (rass *RiskAppetiteStatementService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := rass.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rass.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.RiskAppetiteStatement{}.TableName(),
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

func (rass *RiskAppetiteStatementService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := rass.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rass.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.RiskAppetiteStatement{}.TableName(),
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

func (rass *RiskAppetiteStatementService) validateRiskAppetiteStatementFormDTO(riskAppetiteStatementFormDTO *riskDTOs.RiskAppetiteStatementFormDTO, id *uint) error {
	if _, err := rass.riskCategorySvc.GetByID(riskAppetiteStatementFormDTO.RiskCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk category does not exist"), "not_exist", "risk category")
	}

	if id == nil {
		if exists, err := rass.repo.ExistsByStatement(riskAppetiteStatementFormDTO.StatementEn, riskAppetiteStatementFormDTO.StatementAr, riskAppetiteStatementFormDTO.RiskCategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the statment must be unique"), "unique", "statment")
		}

	} else {
		if exists, err := rass.repo.ExistsByStatementExcludingID(*id, riskAppetiteStatementFormDTO.StatementEn, riskAppetiteStatementFormDTO.StatementAr, riskAppetiteStatementFormDTO.RiskCategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the statment must be unique"), "unique", "statment")
		}
	}

	return nil
}
