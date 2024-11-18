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

type AuditTemplateSettingController struct {
	service services.AuditTemplateSettingServiceInterface
	db      *gorm.DB
}

func NewAuditTemplateSettingController(service services.AuditTemplateSettingServiceInterface) *AuditTemplateSettingController {
	return &AuditTemplateSettingController{
		service: service,
	}
}

func (atsc *AuditTemplateSettingController) GetAllAuditTemplateSettings(c *gin.Context) {
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
	queryParams := settingsDTOs.AuditTemplateSettingListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Filters: filters,
	}

	// Fetch data from the service layer
	auditTemplateSettings, err := atsc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(atsc.db, c, err.Error())
		return
	}

	// Send the response
	c.JSON(http.StatusOK, auditTemplateSettings)
}

func (atsc *AuditTemplateSettingController) UpdateAuditTemplateSetting(c *gin.Context) {
	var updatedTemplateSetting *settingsDTOs.AuditTemplateSettingFormDTO

	auditTemplateSettingID, err := strconv.Atoi(c.Param("id"))
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
		utility.SendInternalErrResp(atsc.db, c, err.Error())
		return
	}

	err = atsc.service.Update(uint(auditTemplateSettingID), currentUserID, updatedTemplateSetting)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(atsc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, nil)
}
