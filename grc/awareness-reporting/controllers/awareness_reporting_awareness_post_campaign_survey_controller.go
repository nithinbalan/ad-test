package controllers

import (
	"errors"
	awarenessDTOs "jethur/grc/awareness-reporting/dtos"
	"jethur/grc/awareness-reporting/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AwarenessPostCampaignSurveyController struct {
	service services.AwarenessPostCampaignSurveyServiceInterface
	db      *gorm.DB
}

func NewAwarenessPostCampaignSurveyController(service services.AwarenessPostCampaignSurveyServiceInterface) *AwarenessPostCampaignSurveyController {
	return &AwarenessPostCampaignSurveyController{
		service: service,
	}
}

func (apcsc *AwarenessPostCampaignSurveyController) GetAllAwarenessPostCampaignSurveys(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "reference_number")
	keyword := c.Query("keyword")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(apcsc.db, c, "Invalid export parameter")
		return
	}

	allowedParams := []string{
		"linked_campaign_ids",
		"maturity_level_ids",
	}

	validSortFields := map[string]string{
		"reference_number":        "awareness_post_campaign_surveys.reference_number",
		"send_date":               "awareness_post_campaign_surveys.send_date",
		"target_audience":         "awareness_post_campaign_surveys.target_audience",
		"recipients":              "awareness_post_campaign_surveys.recipients",
		"response_rate":           "awareness_post_campaign_surveys.response_rate",
		"question_count":          "awareness_post_campaign_surveys.question_count",
		"average_rating":          "awareness_post_campaign_surveys.average_rating",
		"executing_party":         "awareness_post_campaign_surveys.executing_party",
		"linked_campaign_title":   "awareness_campaigns.campaign_title",
		"maturity_level_title_en": "awareness_maturity_levels.title_en",
		"maturity_level_title_ar": "awareness_maturity_levels.title_ar",
		"recommendations":         "awareness_post_campaign_surveys.recommendations",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(apcsc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(apcsc.db, c, "Invalid page_size parameter")
		return
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := awarenessDTOs.AwarenessPostCampaignSurveyListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
		Export:   export,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	surveys, err := apcsc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(apcsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, surveys)
}

func (apcsc *AwarenessPostCampaignSurveyController) GetAwarenessPostCampaignSurveyDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	survey, err := apcsc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(apcsc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, survey)
}

func (apcsc *AwarenessPostCampaignSurveyController) CreateAwarenessPostCampaignSurvey(c *gin.Context) {
	var formDTO awarenessDTOs.AwarenessPostCampaignSurveyFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apcsc.db, c, err.Error())
		return
	}

	id, err := apcsc.service.Create(currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apcsc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (apcsc *AwarenessPostCampaignSurveyController) UpdateAwarenessPostCampaignSurvey(c *gin.Context) {
	var formDTO awarenessDTOs.AwarenessPostCampaignSurveyFormDTO

	surveyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apcsc.db, c, err.Error())
		return
	}

	err = apcsc.service.Update(uint(surveyID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apcsc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (apcsc *AwarenessPostCampaignSurveyController) DeleteAwarenessPostCampaignSurveys(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apcsc.db, c, err.Error())
		return
	}

	failedIDs, err := apcsc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(apcsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
