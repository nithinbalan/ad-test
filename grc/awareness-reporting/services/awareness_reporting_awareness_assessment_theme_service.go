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

type AwarenessAssessmentThemeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessAssessmentTheme, error)
	GetByID(id uint) (*models.AwarenessAssessmentTheme, error)
	Create(currentUserID uint, theme *models.AwarenessAssessmentTheme) (*models.AwarenessAssessmentTheme, error)
	Update(id uint, currentUserID uint, theme *models.AwarenessAssessmentTheme) (*models.AwarenessAssessmentTheme, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type AwarenessAssessmentThemeService struct {
	repo           repositories.AwarenessAssessmentThemeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewAwarenessAssessmentThemeService(repo repositories.AwarenessAssessmentThemeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) AwarenessAssessmentThemeServiceInterface {
	return &AwarenessAssessmentThemeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (aats *AwarenessAssessmentThemeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessAssessmentTheme, error) {
	return aats.repo.GetAll(sort, sortBy, keyword, status)
}

func (aats *AwarenessAssessmentThemeService) GetByID(id uint) (*models.AwarenessAssessmentTheme, error) {
	theme, err := aats.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return theme, nil
}

func (aats *AwarenessAssessmentThemeService) Create(currentUserID uint, theme *models.AwarenessAssessmentTheme) (*models.AwarenessAssessmentTheme, error) {
	if exists, err := aats.repo.ExistsByTitle(theme.TitleEn, theme.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdTheme, err := aats.repo.Create(theme)
	if err != nil {
		return nil, err
	}

	go func() {
		err := aats.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AwarenessAssessmentTheme{}.TableName(),
			createdTheme.ID,
			createdTheme,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdTheme, nil
}

func (aats *AwarenessAssessmentThemeService) Update(id uint, currentUserID uint, updatedTheme *models.AwarenessAssessmentTheme) (*models.AwarenessAssessmentTheme, error) {
	existingTheme, err := aats.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedTheme.TitleEn != existingTheme.TitleEn || updatedTheme.TitleAr != existingTheme.TitleAr {
		if exists, err := aats.repo.ExistsByTitleExcludingID(id, updatedTheme.TitleEn, updatedTheme.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := aats.repo.Update(id, updatedTheme)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingTheme, updatedTheme, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := aats.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AwarenessAssessmentTheme{}.TableName(),
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

func (aats *AwarenessAssessmentThemeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedThemes []map[string]interface{}

	for _, id := range ids {
		theme, err := aats.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = aats.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
		deletedThemes = append(deletedThemes, map[string]interface{}{
			"id":       id,
			"title_en": theme.TitleEn,
			"title_ar": theme.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := aats.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AwarenessAssessmentTheme{}.TableName(),
				successfulIds,
				deletedThemes,
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

func (aats *AwarenessAssessmentThemeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := aats.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := aats.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.AwarenessAssessmentTheme{}.TableName(),
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

func (aats *AwarenessAssessmentThemeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := aats.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := aats.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.AwarenessAssessmentTheme{}.TableName(),
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