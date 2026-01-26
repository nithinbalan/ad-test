package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/violation-report/repositories"
	"jethur/utility"
	"log"
)

// Interface
type ViolationReportChannelServiceInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*models.ViolationReportChannel, error)
	GetByID(id uint) (*models.ViolationReportChannel, error)
	Create(currentUserID uint, vrc *models.ViolationReportChannel) (*models.ViolationReportChannel, error)
	Update(id uint, currentUserID uint, vrc *models.ViolationReportChannel) (*models.ViolationReportChannel, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Implementation
type ViolationReportChannelService struct {
	repo           repositories.ViolationReportChannelRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewViolationReportChannelService(repo repositories.ViolationReportChannelRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ViolationReportChannelServiceInterface {
	return &ViolationReportChannelService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (vrcs *ViolationReportChannelService) GetAll(sort, sortBy, keyword, status string) ([]*models.ViolationReportChannel, error) {
	return vrcs.repo.GetAll(sort, sortBy, keyword, status)
}

func (vrcs *ViolationReportChannelService) GetByID(id uint) (*models.ViolationReportChannel, error) {
	return vrcs.repo.GetByID(id)
}

func (vrcs *ViolationReportChannelService) Create(currentUserID uint, vrc *models.ViolationReportChannel) (*models.ViolationReportChannel, error) {
	if exists, err := vrcs.repo.ExistsByTitle(vrc.TitleEn, vrc.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	created, err := vrcs.repo.Create(vrc)
	if err != nil {
		return nil, err
	}

	go func() {
		err := vrcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ViolationReportChannel{}.TableName(),
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

func (vrcs *ViolationReportChannelService) Update(id uint, currentUserID uint, updated *models.ViolationReportChannel) (*models.ViolationReportChannel, error) {
	existing, err := vrcs.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updated.TitleEn != existing.TitleEn || updated.TitleAr != existing.TitleAr {
		if exists, err := vrcs.repo.ExistsByTitleExcludingID(id, updated.TitleEn, updated.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := vrcs.repo.Update(id, updated)
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
		err := vrcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ViolationReportChannel{}.TableName(),
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

func (vrcs *ViolationReportChannelService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deletedViolationReportChannels []map[string]interface{}

	for _, id := range ids {
		vrc, err := vrcs.repo.GetByID(id)
		if err != nil {
			continue
		}

		err = vrcs.repo.Delete(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
			deletedViolationReportChannels = append(deletedViolationReportChannels, map[string]interface{}{
				"id":       id,
				"title_en": vrc.TitleEn,
				"title_ar": vrc.TitleAr,
			})
		}
	}

	if len(successfulIDs) > 0 {
		go func() {
			err := vrcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ViolationReportChannel{}.TableName(),
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
