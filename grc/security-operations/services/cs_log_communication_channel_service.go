package services

import (
	"errors"
	"fmt"
	"log"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/security-operations/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type CSLogCommunicationChannelServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CSLogCommunicationChannel, error)
	GetByID(id uint) (*models.CSLogCommunicationChannel, error)
	Create(currentUserID uint, channel *models.CSLogCommunicationChannel) (*models.CSLogCommunicationChannel, error)
	Update(id uint, currentUserID uint, channel *models.CSLogCommunicationChannel) (*models.CSLogCommunicationChannel, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type CSLogCommunicationChannelService struct {
	repo           repositories.CSLogCommunicationChannelRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewCSLogCommunicationChannelService(repo repositories.CSLogCommunicationChannelRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) CSLogCommunicationChannelServiceInterface {
	return &CSLogCommunicationChannelService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (clccs *CSLogCommunicationChannelService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CSLogCommunicationChannel, error) {
	return clccs.repo.GetAll(sort, sortBy, keyword, status)
}

func (clccs *CSLogCommunicationChannelService) GetByID(id uint) (*models.CSLogCommunicationChannel, error) {
	channel, err := clccs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return channel, nil
}

func (clccs *CSLogCommunicationChannelService) Create(currentUserID uint, channel *models.CSLogCommunicationChannel) (*models.CSLogCommunicationChannel, error) {
	if exists, err := clccs.repo.ExistsByTitle(channel.TitleEn, channel.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	created, err := clccs.repo.Create(channel)
	if err != nil {
		return nil, err
	}

	go func() {
		err := clccs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.CSLogCommunicationChannel{}.TableName(),
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

func (clccs *CSLogCommunicationChannelService) Update(id uint, currentUserID uint, updated *models.CSLogCommunicationChannel) (*models.CSLogCommunicationChannel, error) {
	existing, err := clccs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if updated.TitleEn != existing.TitleEn || updated.TitleAr != existing.TitleAr {
		if exists, err := clccs.repo.ExistsByTitleExcludingID(id, updated.TitleEn, updated.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := clccs.repo.Update(id, updated)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existing, updated, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := clccs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.CSLogCommunicationChannel{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log update activity: %v\n", err)
		}
	}()
	return result, err
}

func (clccs *CSLogCommunicationChannelService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deleted []map[string]any
	for _, id := range ids {
		channel, err := clccs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		if err := clccs.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deleted = append(deleted, map[string]any{
			"id":       id,
			"title_en": channel.TitleEn,
			"title_ar": channel.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := clccs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.CSLogCommunicationChannel{}.TableName(),
				successfulIds,
				deleted,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (clccs *CSLogCommunicationChannelService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := clccs.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := clccs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.CSLogCommunicationChannel{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk activate activity: %v\n", err)
			}
		}()
	}

	return err
}

func (clccs *CSLogCommunicationChannelService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := clccs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := clccs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.CSLogCommunicationChannel{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk deactivate activity: %v\n", err)
			}
		}()
	}

	return err
}
