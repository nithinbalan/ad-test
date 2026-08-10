package services

import (
	"errors"
	"fmt"
	"log"

	"jethur/grc/identity-access-review/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type PrivilegedAccountServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PrivilegedAccount, error)
	GetByID(id uint) (*models.PrivilegedAccount, error)
	Create(currentUserID uint, privilegedAccount *models.PrivilegedAccount) (*models.PrivilegedAccount, error)
	Update(id uint, currentUserID uint, privilegedAccount *models.PrivilegedAccount) (*models.PrivilegedAccount, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type PrivilegedAccountService struct {
	repo           repositories.PrivilegedAccountRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewPrivilegedAccountService(
	repo repositories.PrivilegedAccountRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) PrivilegedAccountServiceInterface {
	return &PrivilegedAccountService{repo: repo, activityLogSvc: activityLogSvc}
}

func (pas *PrivilegedAccountService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PrivilegedAccount, error) {
	return pas.repo.GetAll(sort, sortBy, keyword, status)
}

func (pas *PrivilegedAccountService) GetByID(id uint) (*models.PrivilegedAccount, error) {
	privilegedAccount, err := pas.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return privilegedAccount, nil
}

func (pas *PrivilegedAccountService) Create(currentUserID uint, privilegedAccount *models.PrivilegedAccount) (*models.PrivilegedAccount, error) {
	if exists, err := pas.repo.ExistsByTitle(privilegedAccount.TitleEn, privilegedAccount.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdAccount, err := pas.repo.Create(privilegedAccount)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := pas.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PrivilegedAccount{}.TableName(),
			createdAccount.ID,
			createdAccount,
			currentUserID,
			"", "", "", nil,
		); err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdAccount, nil
}

func (pas *PrivilegedAccountService) Update(id uint, currentUserID uint, updatedAccount *models.PrivilegedAccount) (*models.PrivilegedAccount, error) {
	existingAccount, err := pas.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedAccount.TitleEn != existingAccount.TitleEn || updatedAccount.TitleAr != existingAccount.TitleAr {
		if exists, err := pas.repo.ExistsByTitleExcludingID(id, updatedAccount.TitleEn, updatedAccount.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := pas.repo.Update(id, updatedAccount)
	if err != nil {
		return nil, err
	}

	changes, compareErr := utility.CompareFields(existingAccount, updatedAccount, []string{"TitleEn", "TitleAr", "Status"})
	if compareErr != nil {
		log.Println("error :", compareErr.Error())
	}
	changedDataJSON, jsonErr := utility.ChangesToJSON(changes)
	if jsonErr != nil {
		log.Println("error :", jsonErr.Error())
	}

	go func() {
		if err := pas.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PrivilegedAccount{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		); err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return result, nil
}

func (pas *PrivilegedAccountService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deletedAccounts []map[string]interface{}

	for _, id := range ids {
		account, err := pas.repo.GetByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if err := pas.repo.Delete(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
		deletedAccounts = append(deletedAccounts, map[string]interface{}{
			"id": id, "title_en": account.TitleEn, "title_ar": account.TitleAr,
		})
	}

	if len(successfulIDs) > 0 {
		go func() {
			if err := pas.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.PrivilegedAccount{}.TableName(),
				successfulIDs,
				deletedAccounts,
				currentUserID,
				nil,
			); err != nil {
				log.Printf("Failed to log bulk delete activity: %v", err)
			}
		}()
	}

	return failedIDs, nil
}

func (pas *PrivilegedAccountService) Activate(ids []uint, currentUserID uint) error {
	return pas.changeStatus(ids, currentUserID, models.ActivateAction, pas.repo.Activate)
}

func (pas *PrivilegedAccountService) Deactivate(ids []uint, currentUserID uint) error {
	return pas.changeStatus(ids, currentUserID, models.DeactivateAction, pas.repo.Deactivate)
}

func (pas *PrivilegedAccountService) changeStatus(
	ids []uint,
	currentUserID uint,
	action models.Action,
	update func([]uint) ([]uint, error),
) error {
	successfulIDs, err := update(ids)
	if len(successfulIDs) > 0 {
		go func() {
			if logErr := pas.activityLogSvc.LogBulkActivities(
				action,
				models.PrivilegedAccount{}.TableName(),
				successfulIDs,
				nil,
				currentUserID,
				nil,
			); logErr != nil {
				log.Printf("Failed to log bulk status activity: %v", logErr)
			}
		}()
	}
	return err
}
