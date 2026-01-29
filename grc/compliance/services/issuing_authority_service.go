package services

import (
	"fmt"
	"log"

	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
)

type IssuingAuthorityServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceIssuingAuthority, error)
	GetByID(id uint) (*models.ComplianceIssuingAuthority, error)
	Create(currentUserID uint, authority *models.ComplianceIssuingAuthority) (*models.ComplianceIssuingAuthority, error)
	Update(id uint, currentUserID uint, authority *models.ComplianceIssuingAuthority) (*models.ComplianceIssuingAuthority, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementation
type IssuingAuthorityService struct {
	repo           repositories.IssuingAuthorityRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewIssuingAuthorityService(
	repo repositories.IssuingAuthorityRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) IssuingAuthorityServiceInterface {
	return &IssuingAuthorityService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ias *IssuingAuthorityService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceIssuingAuthority, error) {
	return ias.repo.GetAll(sort, sortBy, keyword, status)
}

func (ias *IssuingAuthorityService) GetByID(id uint) (*models.ComplianceIssuingAuthority, error) {
	return ias.repo.GetByID(id)
}

func (ias *IssuingAuthorityService) Create(currentUserID uint, authority *models.ComplianceIssuingAuthority) (*models.ComplianceIssuingAuthority, error) {
	if exists, err := ias.repo.ExistsByTitle(authority.TitleEn, authority.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(
			fmt.Errorf("the title must be unique"),
			"unique",
			"title",
		)
	}

	createdAuthority, err := ias.repo.Create(authority)
	if err != nil {
		return nil, err
	}

	go func() {
		err := ias.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceIssuingAuthority{}.TableName(),
			createdAuthority.ID,
			createdAuthority,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdAuthority, nil
}

func (ias *IssuingAuthorityService) Update(
	id uint,
	currentUserID uint,
	updatedAuthority *models.ComplianceIssuingAuthority,
) (*models.ComplianceIssuingAuthority, error) {

	existingAuthority, err := ias.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedAuthority.TitleEn != existingAuthority.TitleEn ||
		updatedAuthority.TitleAr != existingAuthority.TitleAr {

		if exists, err := ias.repo.ExistsByTitleExcludingID(
			id,
			updatedAuthority.TitleEn,
			updatedAuthority.TitleAr,
		); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(
				fmt.Errorf("the title must be unique"),
				"unique",
				"title",
			)
		}
	}

	result, err := ias.repo.Update(id, updatedAuthority)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingAuthority, updatedAuthority, fieldsToCompare)
	if err != nil {
		log.Println("compare error:", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("json error:", err.Error())
	}

	go func() {
		err := ias.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ComplianceIssuingAuthority{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()

	return result, nil
}

func (ias *IssuingAuthorityService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deletedAuthorities []map[string]interface{}

	for _, id := range ids {
		authority, err := ias.repo.GetByID(id)
		if err != nil {
			continue
		}

		if err := ias.repo.Delete(id); err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
		}

		deletedAuthorities = append(deletedAuthorities, map[string]interface{}{
			"id":       id,
			"title_en": authority.TitleEn,
			"title_ar": authority.TitleAr,
		})
	}

	if len(successfulIDs) > 0 {
		go func() {
			err := ias.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ComplianceIssuingAuthority{}.TableName(),
				successfulIDs,
				deletedAuthorities,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return failedIDs, nil
}

func (ias *IssuingAuthorityService) Activate(ids []uint, currentUserID uint) error {
	successfulIDs, err := ias.repo.Activate(ids)

	if len(successfulIDs) > 0 {
		go func() {
			err := ias.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ComplianceIssuingAuthority{}.TableName(),
				successfulIDs,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (ias *IssuingAuthorityService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIDs, err := ias.repo.Deactivate(ids)

	if len(successfulIDs) > 0 {
		go func() {
			err := ias.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ComplianceIssuingAuthority{}.TableName(),
				successfulIDs,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}
