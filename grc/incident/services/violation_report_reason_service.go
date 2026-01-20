package services

import (
	"fmt"
	"jethur/grc/incident/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"
)

// Interface
type ViolationReportReasonServiceInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*models.ViolationReportReason, error)
	GetByID(id uint) (*models.ViolationReportReason, error)
	Create(currentUserID uint, vrr *models.ViolationReportReason) (*models.ViolationReportReason, error)
	Update(id uint, currentUserID uint, vrr *models.ViolationReportReason) (*models.ViolationReportReason, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Implementation
type ViolationReportReasonService struct {
	repo           repositories.ViolationReportReasonRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewViolationReportReasonService(repo repositories.ViolationReportReasonRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ViolationReportReasonServiceInterface {
	return &ViolationReportReasonService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (vrrs *ViolationReportReasonService) GetAll(sort, sortBy, keyword, status string) ([]*models.ViolationReportReason, error) {
	return vrrs.repo.GetAll(sort, sortBy, keyword, status)
}

func (vrrs *ViolationReportReasonService) GetByID(id uint) (*models.ViolationReportReason, error) {
	return vrrs.repo.GetByID(id)
}

func (vrrs *ViolationReportReasonService) Create(currentUserID uint, vrr *models.ViolationReportReason) (*models.ViolationReportReason, error) {
	if exists, err := vrrs.repo.ExistsByTitle(vrr.TitleEn, vrr.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	created, err := vrrs.repo.Create(vrr)
	if err != nil {
		return nil, err
	}

	go func() {
		err := vrrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ViolationReportReason{}.TableName(),
			created.ID,
			created,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return created, nil
}

func (vrrs *ViolationReportReasonService) Update(id uint, currentUserID uint, updated *models.ViolationReportReason) (*models.ViolationReportReason, error) {
	existing, err := vrrs.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updated.TitleEn != existing.TitleEn || updated.TitleAr != existing.TitleAr {
		if exists, err := vrrs.repo.ExistsByTitleExcludingID(id, updated.TitleEn, updated.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := vrrs.repo.Update(id, updated)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existing, updated, fieldsToCompare)
	if err != nil {
		log.Println("compare error:", err)
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("json error:", err)
	}

	go func() {
		err := vrrs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ViolationReportReason{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log update activity: %v", err)
		}
	}()

	return result, nil
}

func (vrrs *ViolationReportReasonService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deletedViolationReportChannels []map[string]interface{}

	for _, id := range ids {
		vrr, err := vrrs.repo.GetByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		err = vrrs.repo.Delete(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
			deletedViolationReportChannels = append(deletedViolationReportChannels, map[string]interface{}{
				"id":       id,
				"title_en": vrr.TitleEn,
				"title_ar": vrr.TitleAr,
			})
		}
	}

	if len(successfulIDs) > 0 {
		go func() {
			err := vrrs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ViolationReportReason{}.TableName(),
				successfulIDs,
				deletedViolationReportChannels,
				currentUserID,
				nil,
			)
			if err != nil {
				log.Printf("Failed to log delete activity: %v", err)
			}
		}()
	}

	return failedIDs, nil
}
