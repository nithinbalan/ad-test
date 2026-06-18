package services

import (
	"errors"
	"fmt"
	"log"

	"jethur/grc/awareness-reporting/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AwarenessChannelServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessChannel, error)
	GetByID(id uint) (*models.AwarenessChannel, error)
	Create(currentUserID uint, awarenessChannelForm *models.AwarenessChannel) (*models.AwarenessChannel, error)
	Update(id uint, currentUserID uint, awarenessChannelForm *models.AwarenessChannel) (*models.AwarenessChannel, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type AwarenessChannelService struct {
	repo           repositories.AwarenessChannelRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewAwarenessChannelService(
	repo repositories.AwarenessChannelRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
) AwarenessChannelServiceInterface {
	return &AwarenessChannelService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (acs *AwarenessChannelService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessChannel, error) {
	return acs.repo.GetAll(sort, sortBy, keyword, status)
}

func (acs *AwarenessChannelService) GetByID(id uint) (*models.AwarenessChannel, error) {
	awarenessChannel, err := acs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return awarenessChannel, nil
}

func (acs *AwarenessChannelService) Create(currentUserID uint, awarenessChannel *models.AwarenessChannel) (*models.AwarenessChannel, error) {
	if exists, err := acs.repo.ExistsByTitle(awarenessChannel.TitleEn, awarenessChannel.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdChannel, err := acs.repo.Create(awarenessChannel)
	if err != nil {
		return nil, err
	}

	go func() {
		err := acs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AwarenessChannel{}.TableName(),
			createdChannel.ID,
			createdChannel,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdChannel, nil
}

func (acs *AwarenessChannelService) Update(id uint, currentUserID uint, updatedChannel *models.AwarenessChannel) (*models.AwarenessChannel, error) {
	existingChannel, err := acs.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedChannel.TitleEn != existingChannel.TitleEn || updatedChannel.TitleAr != existingChannel.TitleAr {
		if exists, err := acs.repo.ExistsByTitleExcludingID(id, updatedChannel.TitleEn, updatedChannel.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := acs.repo.Update(id, updatedChannel)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}

	changes, err := utility.CompareFields(existingChannel, updatedChannel, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := acs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AwarenessChannel{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

	return result, nil
}

func (acs *AwarenessChannelService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedChannels []map[string]interface{}

	for _, id := range ids {
		channel, err := acs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = acs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedChannels = append(deletedChannels, map[string]interface{}{
			"id":       id,
			"title_en": channel.TitleEn,
			"title_ar": channel.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := acs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AwarenessChannel{}.TableName(),
				successfulIds,
				deletedChannels,
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

func (acs *AwarenessChannelService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := acs.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := acs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.AwarenessChannel{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (acs *AwarenessChannelService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := acs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := acs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.AwarenessChannel{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}
