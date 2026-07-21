package services

import (
	"errors"
	"fmt"
	"os"

	awarenessDTOs "jethur/grc/awareness-reporting/dtos"
	awarenessRepositories "jethur/grc/awareness-reporting/repositories"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AwarenessPostCampaignSurveyServiceInterface interface {
	GetAll(params *awarenessDTOs.AwarenessPostCampaignSurveyListQueryParamsDTO) (*awarenessDTOs.AwarenessPostCampaignSurveyListResponseDTO, error)
	GetByID(id uint) (*models.AwarenessPostCampaignSurvey, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*awarenessDTOs.AwarenessPostCampaignSurveyDetailResponseDTO, error)
	Create(currentUserID uint, formDTO *awarenessDTOs.AwarenessPostCampaignSurveyFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *awarenessDTOs.AwarenessPostCampaignSurveyFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	ExportAwarenessPostCampaignSurvey(data []*awarenessDTOs.AwarenessPostCampaignSurveyListDTO, currentUserID uint) (*os.File, error)
}

type AwarenessPostCampaignSurveyService struct {
	repo                        awarenessRepositories.AwarenessPostCampaignSurveyRepositoryInterface
	activityLogSvc              logServices.ActivityLogServiceInterface
	awarenessMaturityLevelSvc   AwarenessMaturityLevelServiceInterface
	awarenessCampaignService    AwarenessCampaignServiceInterface
	awarenessAssessmentThemeSvc AwarenessAssessmentThemeServiceInterface
	userSettingSvc              organizationServices.UserSettingServiceInterface
}

func NewAwarenessPostCampaignSurveyService(
	repo awarenessRepositories.AwarenessPostCampaignSurveyRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	awarenessMaturityLevelSvc AwarenessMaturityLevelServiceInterface,
	awarenessCampaignService AwarenessCampaignServiceInterface,
	awarenessAssessmentThemeSvc AwarenessAssessmentThemeServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) AwarenessPostCampaignSurveyServiceInterface {
	return &AwarenessPostCampaignSurveyService{
		repo:                        repo,
		activityLogSvc:              activityLogSvc,
		awarenessMaturityLevelSvc:   awarenessMaturityLevelSvc,
		awarenessCampaignService:    awarenessCampaignService,
		awarenessAssessmentThemeSvc: awarenessAssessmentThemeSvc,
		userSettingSvc:              userSettingSvc,
	}
}

func (apcss *AwarenessPostCampaignSurveyService) GetAll(params *awarenessDTOs.AwarenessPostCampaignSurveyListQueryParamsDTO) (*awarenessDTOs.AwarenessPostCampaignSurveyListResponseDTO, error) {
	return apcss.repo.GetAll(params)
}

func (apcss *AwarenessPostCampaignSurveyService) GetByID(id uint) (*models.AwarenessPostCampaignSurvey, error) {
	survey, err := apcss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return survey, nil
}

func (apcss *AwarenessPostCampaignSurveyService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*awarenessDTOs.AwarenessPostCampaignSurveyDetailResponseDTO, error) {
	var surveyDTO awarenessDTOs.AwarenessPostCampaignSurveyDetailResponseDTO

	survey, err := apcss.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if err := mapstructure.Decode(survey, &surveyDTO); err != nil {
		return nil, err
	}

	surveyDTO.CreatedAt = survey.CreatedAt
	surveyDTO.UpdatedAt = survey.UpdatedAt

	surveyDTO.SurveyThemes = nil
	for _, theme := range survey.SurveyThemes {
		if theme.AwarenessAssessmentTheme != nil {
			surveyDTO.SurveyThemes = append(surveyDTO.SurveyThemes, &awarenessDTOs.AwarenessAssessmentThemeDTO{
				ID:      theme.AwarenessAssessmentTheme.ID,
				TitleEn: theme.AwarenessAssessmentTheme.TitleEn,
				TitleAr: theme.AwarenessAssessmentTheme.TitleAr,
			})
		}
	}

	return &surveyDTO, nil
}

func (apcss *AwarenessPostCampaignSurveyService) Create(currentUserID uint, formDTO *awarenessDTOs.AwarenessPostCampaignSurveyFormDTO) (uint, error) {
	tx := apcss.repo.BeginTransaction()

	if err := apcss.validateAwarenessPostCampaignSurveyFormDTO(formDTO); err != nil {
		return 0, err
	}

	referenceNumber, err := apcss.repo.GenerateReferenceNumber()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	survey := &models.AwarenessPostCampaignSurvey{
		ReferenceNumber:  referenceNumber,
		LinkedCampaignID: formDTO.LinkedCampaignID,
		SendDate:         formDTO.SendDate,
		TargetAudience:   formDTO.TargetAudience,
		Recipients:       formDTO.Recipients,
		ResponseRate:     formDTO.ResponseRate,
		QuestionCount:    formDTO.QuestionCount,
		AverageRating:    formDTO.AverageRating,
		MaturityLevelID:  formDTO.MaturityLevelID,
		Recommendations:  formDTO.Recommendations,
		ExecutingParty:   formDTO.ExecutingParty,
		ResultsLink:      formDTO.ResultsLink,
		CreatedBYID:      currentUserID,
	}

	id, err := apcss.repo.Create(survey, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if len(formDTO.AwarenessAssessmentThemeIDs) > 0 {
		var themes []*models.AwarenessPostCampaignSurveyTheme
		for _, themeID := range formDTO.AwarenessAssessmentThemeIDs {
			themes = append(themes, &models.AwarenessPostCampaignSurveyTheme{
				AwarenessPostCampaignSurveyID: id,
				AwarenessAssessmentThemeID:    themeID,
				CreatedByID:                   currentUserID,
			})
		}
		if err := apcss.repo.CreateSurveyThemes(themes, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		existingSurvey, err := apcss.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch post campaign survey", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = apcss.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AwarenessPostCampaignSurvey{}.TableName(),
			id,
			existingSurvey,
			currentUserID,
			"awareness", "awareness_post_campaign_survey", "survey", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AwarenessPostCampaignSurvey{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (apcss *AwarenessPostCampaignSurveyService) Update(id uint, currentUserID uint, formDTO *awarenessDTOs.AwarenessPostCampaignSurveyFormDTO) error {
	tx := apcss.repo.BeginTransaction()

	existingSurvey, err := apcss.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldSurvey, err := apcss.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := apcss.validateAwarenessPostCampaignSurveyFormDTO(formDTO); err != nil {
		return err
	}

	existingSurvey.LinkedCampaignID = formDTO.LinkedCampaignID
	existingSurvey.SendDate = formDTO.SendDate
	existingSurvey.TargetAudience = formDTO.TargetAudience
	existingSurvey.Recipients = formDTO.Recipients
	existingSurvey.ResponseRate = formDTO.ResponseRate
	existingSurvey.QuestionCount = formDTO.QuestionCount
	existingSurvey.AverageRating = formDTO.AverageRating
	existingSurvey.MaturityLevelID = formDTO.MaturityLevelID
	existingSurvey.Recommendations = formDTO.Recommendations
	existingSurvey.ExecutingParty = formDTO.ExecutingParty
	existingSurvey.ResultsLink = formDTO.ResultsLink
	existingSurvey.UpdatedBYID = &currentUserID

	if err := apcss.repo.Update(existingSurvey, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := apcss.repo.DeleteSurveyThemes(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	if len(formDTO.AwarenessAssessmentThemeIDs) > 0 {
		var themes []*models.AwarenessPostCampaignSurveyTheme
		for _, themeID := range formDTO.AwarenessAssessmentThemeIDs {
			themes = append(themes, &models.AwarenessPostCampaignSurveyTheme{
				AwarenessPostCampaignSurveyID: id,
				AwarenessAssessmentThemeID:    themeID,
				CreatedByID:                   currentUserID,
			})
		}
		if err := apcss.repo.CreateSurveyThemes(themes, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		existingFullSurvey, err := apcss.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		changedData := make(map[string]interface{})

		oldMaturityLevel := ""
		if oldSurvey.MaturityLevel != nil {
			oldMaturityLevel = oldSurvey.MaturityLevel.TitleEn
		}
		newMaturityLevel := ""
		if existingFullSurvey.MaturityLevel != nil {
			newMaturityLevel = existingFullSurvey.MaturityLevel.TitleEn
		}

		oldData := map[string]interface{}{
			"maturity_level":  oldMaturityLevel,
			"target_audience": utility.GetStringValue(oldSurvey.TargetAudience),
			"recipients":      oldSurvey.Recipients,
			"response_rate":   oldSurvey.ResponseRate,
			"question_count":  oldSurvey.QuestionCount,
			"average_rating":  oldSurvey.AverageRating,
			"recommendations": utility.GetStringValue(oldSurvey.Recommendations),
			"executing_party": utility.GetStringValue(oldSurvey.ExecutingParty),
			"results_link":    utility.GetStringValue(oldSurvey.ResultsLink),
			"send_date":       oldSurvey.SendDate,
		}

		newData := map[string]interface{}{
			"maturity_level":  newMaturityLevel,
			"target_audience": utility.GetStringValue(existingFullSurvey.TargetAudience),
			"recipients":      existingFullSurvey.Recipients,
			"response_rate":   existingFullSurvey.ResponseRate,
			"question_count":  existingFullSurvey.QuestionCount,
			"average_rating":  existingFullSurvey.AverageRating,
			"recommendations": utility.GetStringValue(existingFullSurvey.Recommendations),
			"executing_party": utility.GetStringValue(existingFullSurvey.ExecutingParty),
			"results_link":    utility.GetStringValue(existingFullSurvey.ResultsLink),
			"send_date":       existingFullSurvey.SendDate,
		}

		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		if len(changedData) > 0 {
			err := apcss.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AwarenessPostCampaignSurvey{}.TableName(),
				id,
				changedData,
				currentUserID,
				"awareness", "awareness_post_campaign_survey", "survey", nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AwarenessPostCampaignSurvey{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}
	}()

	return nil
}

func (apcss *AwarenessPostCampaignSurveyService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := apcss.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedSurveys []map[string]interface{}

	for _, id := range ids {
		survey, err := apcss.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedSurveys = append(deletedSurveys, map[string]interface{}{
			"id":               id,
			"reference_number": survey.ReferenceNumber,
		})

		if err := apcss.repo.DeleteSurveyThemes(id, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := apcss.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := apcss.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AwarenessPostCampaignSurvey{}.TableName(),
				successfulIds,
				deletedSurveys,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}()
	}

	return failedIds, nil
}

func (apcss *AwarenessPostCampaignSurveyService) validateAwarenessPostCampaignSurveyFormDTO(formDTO *awarenessDTOs.AwarenessPostCampaignSurveyFormDTO) error {
	if formDTO.LinkedCampaignID != nil {
		if _, err := apcss.awarenessCampaignService.GetByID(*formDTO.LinkedCampaignID); err != nil {
			return utility.NewValidationError(fmt.Errorf("linked campaign does not exist"), "not_exist", "linked_campaign_id")
		}
	}
	if formDTO.MaturityLevelID != nil {
		if _, err := apcss.awarenessMaturityLevelSvc.GetByID(*formDTO.MaturityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("maturity level does not exist"), "not_exist", "maturity level")
		}
	}

	for _, themeID := range formDTO.AwarenessAssessmentThemeIDs {
		if _, err := apcss.awarenessAssessmentThemeSvc.GetByID(themeID); err != nil {
			return utility.NewValidationError(fmt.Errorf("awareness assessment theme with id %d does not exist", themeID), "not_exist", "awareness_assessment_theme_ids")
		}
	}

	return nil
}

func (apcss *AwarenessPostCampaignSurveyService) ExportAwarenessPostCampaignSurvey(data []*awarenessDTOs.AwarenessPostCampaignSurveyListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := apcss.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"الرقم التسلسلي",
			"الرقم المرجعي",
			"الحملة المرتبطة",
			"تاريخ الإرسال",
			"الجمهور المستهدف",
			"عدد المستلمين",
			"معدل الاستجابة",
			"عدد الأسئلة",
			"متوسط التقييم",
			"مستوى النضج",
			"الجهة المنفذة",
			"التوصيات",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Linked Campaign",
			"Send Date",
			"Target Audience",
			"Recipients",
			"Response Rate",
			"Question Count",
			"Average Rating",
			"Maturity Level",
			"Executing Party",
			"Recommendations",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, survey := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				survey.ReferenceNumber,
				utility.GetStringValue(survey.LinkedCampaignTitle),
				survey.SendDate,
				utility.GetStringValue(survey.TargetAudience),
				survey.Recipients,
				survey.ResponseRate,
				survey.QuestionCount,
				survey.AverageRating,
				utility.GetStringValue(survey.MaturityLevelTitleAr),
				utility.GetStringValue(survey.ExecutingParty),
				utility.GetStringValue(survey.Recommendations),
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				survey.ReferenceNumber,
				utility.GetStringValue(survey.LinkedCampaignTitle),
				survey.SendDate,
				utility.GetStringValue(survey.TargetAudience),
				survey.Recipients,
				survey.ResponseRate,
				survey.QuestionCount,
				survey.AverageRating,
				utility.GetStringValue(survey.MaturityLevelTitleEn),
				utility.GetStringValue(survey.ExecutingParty),
				utility.GetStringValue(survey.Recommendations),
			}
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}
