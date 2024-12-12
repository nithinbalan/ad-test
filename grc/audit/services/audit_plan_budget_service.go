package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditPlanBudgetServiceInterface interface {
	GetAll(auditPlanID uint, params *auditDTOs.AuditPlanBudgetListQueryParamsDTO) ([]*auditDTOs.AuditPlanBudgetListDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanBudgetDetailDTO, error)
	GetByID(id uint) (*models.AuditPlanBudget, error)
	Create(auditPlanID uint, currentUserID uint, auditPlanBudgetFormDTO *auditDTOs.AuditPlanBudgetFormDTO) (uint, error)
	Update(id uint, auditPlanID uint, currentUserID uint, auditPlanBudgetFormDTO *auditDTOs.AuditPlanBudgetFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type AuditPlanBudgetService struct {
	repo               repositories.AuditPlanBudgetRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	auditPlanSvc       AuditPlanServiceInterface
	expenseCategorySvc masterServices.ExpenseCategoryServiceInterface
}

func NewAuditPlanBudgetService(
	repo repositories.AuditPlanBudgetRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditPlanSvc AuditPlanServiceInterface,
	expenseCategorySvc masterServices.ExpenseCategoryServiceInterface,
) AuditPlanBudgetServiceInterface {
	return &AuditPlanBudgetService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		auditPlanSvc:       auditPlanSvc,
		expenseCategorySvc: expenseCategorySvc,
	}
}

func (apbs *AuditPlanBudgetService) GetAll(auditPlanID uint, params *auditDTOs.AuditPlanBudgetListQueryParamsDTO) ([]*auditDTOs.AuditPlanBudgetListDTO, error) {
	if _, err := apbs.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return nil, err
	}
	return apbs.repo.GetAll(auditPlanID, params)
}

func (apbs *AuditPlanBudgetService) GetByID(id uint) (*models.AuditPlanBudget, error) {
	auditPlanBudget, err := apbs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditPlanBudget, nil
}

func (apbs *AuditPlanBudgetService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanBudgetDetailDTO, error) {
	var auditPlanBudgetDTO auditDTOs.AuditPlanBudgetDetailDTO

	auditPlanBudget, err := apbs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditPlanBudget, &auditPlanBudgetDTO)
	if err != nil {
		return nil, err
	}

	auditPlanBudgetDTO.CreatedAt = auditPlanBudget.CreatedAt
	auditPlanBudgetDTO.UpdatedAt = auditPlanBudget.UpdatedAt

	return &auditPlanBudgetDTO, nil
}

func (apbs *AuditPlanBudgetService) Create(auditPlanID uint, currentUserID uint, auditPlanBudgetFormDTO *auditDTOs.AuditPlanBudgetFormDTO) (uint, error) {
	tx := apbs.repo.BeginTransaction()
	if _, err := apbs.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return 0, err
	}
	if err := apbs.validateAuditPlanBudgetFormDTO(auditPlanBudgetFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	auditPlanBudget := &models.AuditPlanBudget{
		AuditPlanID:       auditPlanID,
		ExpenseCategoryID: auditPlanBudgetFormDTO.ExpenseCategoryID,
		Title:             auditPlanBudgetFormDTO.Title,
		EstimatedCost:     auditPlanBudgetFormDTO.EstimatedCost,
		Comments:          auditPlanBudgetFormDTO.Comments,
		CreatedBYID:       currentUserID,
	}

	id, err := apbs.repo.Create(auditPlanBudget, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingAuditPlanBudget, err := apbs.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch audit plan tool", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = apbs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditPlanBudget{}.TableName(),
			id,
			existingAuditPlanBudget,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditPlanBudget{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (apbs *AuditPlanBudgetService) Update(id uint, auditPlanID uint, currentUserID uint, auditPlanBudgetFormDTO *auditDTOs.AuditPlanBudgetFormDTO) error {
	tx := apbs.repo.BeginTransaction()
	if _, err := apbs.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return err
	}
	if err := apbs.validateAuditPlanBudgetFormDTO(auditPlanBudgetFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	existingAuditPlanBudget, err := apbs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	oldAuditPlanBudget, err := apbs.GetFullDetailsByID(id)
	if err != nil {
		return err
	}

	existingAuditPlanBudget.ExpenseCategoryID = auditPlanBudgetFormDTO.ExpenseCategoryID
	existingAuditPlanBudget.Title = auditPlanBudgetFormDTO.Title
	existingAuditPlanBudget.EstimatedCost = auditPlanBudgetFormDTO.EstimatedCost
	existingAuditPlanBudget.Comments = auditPlanBudgetFormDTO.Comments
	existingAuditPlanBudget.UpdatedBYID = &currentUserID

	err = apbs.repo.Update(existingAuditPlanBudget, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingPlanBudget, err := apbs.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"expense_category": oldAuditPlanBudget.ExpenseCategory.TitleEn,
			"title":            oldAuditPlanBudget.Title,
			"estimated_cost":   oldAuditPlanBudget.EstimatedCost,
			"comments":         oldAuditPlanBudget.Comments,
		}

		newData := map[string]interface{}{
			"expense_category": existingPlanBudget.ExpenseCategory.TitleEn,
			"title":            existingPlanBudget.Title,
			"estimated_cost":   existingPlanBudget.EstimatedCost,
			"comments":         existingPlanBudget.Comments,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := apbs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AuditPlanBudget{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AuditPlanBudget{}.TableName(),
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

func (apbs *AuditPlanBudgetService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := apbs.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint
	var deletedAuditPlanBudget []map[string]interface{}
	for _, id := range ids {
		auditPlanBudget, err := apbs.repo.GetByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		if err := apbs.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)

		deletedAuditPlanBudget = append(deletedAuditPlanBudget, map[string]interface{}{
			"id":    id,
			"title": auditPlanBudget.Title,
		})
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := apbs.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditPlanBudget{}.TableName(),
					id,
					deletedAuditPlanBudget,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditPlanBudget{}.TableName(),
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

func (apbs *AuditPlanBudgetService) validateAuditPlanBudgetFormDTO(dto *auditDTOs.AuditPlanBudgetFormDTO) error {
	if _, err := apbs.expenseCategorySvc.GetByID(dto.ExpenseCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("expense category does not exist"), "not_exist", "expense category")
	}

	return nil
}
