package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ThemeServiceInterface interface {
	GetAll(Keyword string,SortBy string,Sort string) ([]*dtos.ThemeListDTO, error)
	GetFullDetailsByID(themeID uint) (*dtos.ThemeDetailDTO, error)
	Create(currentUserID uint,form *dtos.ThemeFormDTO) (uint, error)
	Update(id uint,currentUserID uint,form *dtos.ThemeFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type ThemeService struct {
	repo            repositories.ThemeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewThemeService(repo repositories.ThemeRepositoryInterface,activityLogSvc services.ActivityLogServiceInterface,) *ThemeService {
	return &ThemeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}


func (ts *ThemeService) GetAll(Keyword string,SortBy string,Sort string) ([]*dtos.ThemeListDTO, error) {
	return ts.repo.GetAll(Keyword ,SortBy ,Sort)
}


func (ts *ThemeService) GetFullDetailsByID(themeID uint) (*dtos.ThemeDetailDTO, error) {
	var themeDTO dtos.ThemeDetailDTO

	theme, err := ts.repo.GetFullDetailsByID(themeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, themeID)
		}
		return nil, err
	}

	err = mapstructure.Decode(theme, &themeDTO)
	if err != nil {
		return nil, err
	}

	return &themeDTO, nil
}


func (ts *ThemeService) Create(currentUserID uint,themeForm *dtos.ThemeFormDTO) (uint, error) {
	if err := ts.validateThemeForm(themeForm); err != nil {
		return 0, err
	}
	
	theme := &models.Theme{
		TitleEn:        themeForm.TitleEn,
		TitleAr:        themeForm.TitleAr,
		DescriptionEn:  themeForm.DescriptionEn,
		DescriptionAr:  themeForm.DescriptionAr,
		PrimaryColor:   themeForm.PrimaryColor,
		SecondaryColor: themeForm.SecondaryColor,
		LinkColor:      themeForm.LinkColor,
		SidebarBgColor: themeForm.SidebarBgColor,
		IsSystem:       themeForm.IsSystem,
		CreatedBYID:    currentUserID,
	}

	createdTheme, err := ts.repo.Create(theme)
	if err != nil {
		return 0, err
	}

	 
	go func() {
		err := ts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Theme{}.TableName(),
			createdTheme.ID,
			createdTheme,
			currentUserID,
			"setting", "general-setting", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.Theme{}.TableName(),
				"recordID":  createdTheme.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return createdTheme.ID, nil
}

func (ts *ThemeService) Update(id uint,currentUserID uint,themeForm *dtos.ThemeFormDTO) error {

	existingTheme, err := ts.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := ts.validateThemeForm(themeForm); err != nil {
		return err
	}

	updatedTheme := *existingTheme

	updatedTheme.TitleEn = themeForm.TitleEn
	updatedTheme.TitleAr = themeForm.TitleAr
	updatedTheme.DescriptionEn = themeForm.DescriptionEn
	updatedTheme.DescriptionAr = themeForm.DescriptionAr
	updatedTheme.PrimaryColor = themeForm.PrimaryColor
	updatedTheme.SecondaryColor = themeForm.SecondaryColor
	updatedTheme.LinkColor = themeForm.LinkColor
	updatedTheme.SidebarBgColor = themeForm.SidebarBgColor
	updatedTheme.IsSystem = themeForm.IsSystem

	updatedTheme.UpdatedBYID = &currentUserID
	now := time.Now()
	updatedTheme.UpdatedAt = &now

	if err := ts.repo.Update(&updatedTheme); err != nil {
		return err
	}

	fieldsToCompare := []string{
		"TitleEn",
		"TitleAr",
		"DescriptionEn",
		"DescriptionAr",
		"PrimaryColor",
		"SecondaryColor",
		"LinkColor",
		"SidebarBgColor",
		"IsSystem",
	}

	changes, err := utility.CompareFields(existingTheme, &updatedTheme, fieldsToCompare)
	if err != nil {
		log.Println("compare error:", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("changes json error:", err.Error())
	}

	go func() {
		if err := ts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.Theme{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			"setting", "general-setting", "", nil,
		); err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.Theme{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}


func (ts *ThemeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {

	var failedIDs []uint
	var successfulIDs []uint
	var deletedThemes []map[string]interface{}

	for _, id := range ids {
		theme, err := ts.repo.GetByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if err := ts.repo.Delete(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
		deletedThemes = append(deletedThemes, map[string]interface{}{
			"id":       id,
			"title_en": theme.TitleEn,
		})
	}

	if len(successfulIDs) > 0 {
		go func() {
			err := ts.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.DocumentVersionClauseComment{}.TableName(),
				successfulIDs,
				deletedThemes,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIDs, nil
}

func (ts *ThemeService) validateThemeForm(form *dtos.ThemeFormDTO) error {

	HexColorRegex := regexp.MustCompile(`^#([A-Fa-f0-9]{6})$`)

	if strings.TrimSpace(form.TitleEn) == "" {
		return utility.NewValidationError(nil, "title_en", "Title (English) is required")
	}

	if strings.TrimSpace(form.TitleAr) == "" {
		return utility.NewValidationError(nil, "title_ar", "Title (Arabic) is required")
	}

	if strings.TrimSpace(form.PrimaryColor) == "" {
		return utility.NewValidationError(nil, "primary_color", "Primary color is required")
	}

	if !HexColorRegex.MatchString(form.PrimaryColor) {
		return utility.NewValidationError(
			nil,
			"primary_color",
			"Primary color must be a valid hex value (e.g. #165D31)",
		)
	}

	if form.SecondaryColor != nil && !HexColorRegex.MatchString(*form.SecondaryColor) {
		return utility.NewValidationError(nil, "secondary_color", "Secondary color must be a valid hex value")
	}

	if form.LinkColor != nil && !HexColorRegex.MatchString(*form.LinkColor) {
		return utility.NewValidationError(nil, "link_color", "Link color must be a valid hex value")
	}

	if form.SidebarBgColor != nil && !HexColorRegex.MatchString(*form.SidebarBgColor) {
		return utility.NewValidationError(nil, "sidebar_bg_color", "Sidebar background color must be a valid hex value")
	}

	return nil
}






