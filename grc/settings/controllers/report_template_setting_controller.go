package controllers

import (
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportTemplateSettingController struct {
	service services.ReportTemplateSettingServiceInterface
	db      *gorm.DB
}

func NewReportTemplateSettingController(service services.ReportTemplateSettingServiceInterface) *ReportTemplateSettingController {
	return &ReportTemplateSettingController{
		service: service,
	}
}

func (rtsc *ReportTemplateSettingController) GetAllReportTemplateSettings(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	// Define an array of allowed parameter names for filters
	allowedParams := []string{
		"created_by_ids",
		"updated_by_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Build the query parameters DTO
	queryParams := settingsDTOs.ReportTemplateSettingListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Filters: filters,
	}

	// Fetch data from the service layer
	reportTemplateSettings, err := rtsc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(rtsc.db, c, err.Error())
		return
	}

	// Send the response
	c.JSON(http.StatusOK, reportTemplateSettings)
}

func (rtsc *ReportTemplateSettingController) UpdateReportTemplateSetting(c *gin.Context) {
	var updatedTemplateSetting *settingsDTOs.ReportTemplateSettingFormDTO

	reportTemplateSettingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedTemplateSetting); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rtsc.db, c, err.Error())
		return
	}

	err = rtsc.service.Update(uint(reportTemplateSettingID), currentUserID, updatedTemplateSetting)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rtsc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, nil)
}

func (rtsc *ReportTemplateSettingController) GetReportTemplateSettingPreview(c *gin.Context) {
	templateType := c.Query("type")
	if templateType == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	reportTemplateSettingImagePath, err := rtsc.service.GetImagePath(templateType)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rtsc.db, c, err.Error())
		}
		return
	}
	c.File(reportTemplateSettingImagePath)
}
