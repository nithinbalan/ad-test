package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IntegrationController struct {
	service services.IntegrationServiceInterface
	db      *gorm.DB
}

func NewIntegrationController(service services.IntegrationServiceInterface) *IntegrationController {
	return &IntegrationController{
		service: service,
	}
}

func (icc *IntegrationController) GetAllIntegrations(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"integration_category_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := settingsDTOs.IntegrationListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}

	integrations, err := icc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, integrations)
}

func (ic *IntegrationController) GetIntegrationDetails(c *gin.Context) {

	integrationIDStr := c.Param("integrationId")

	integrationID, err := strconv.ParseUint(integrationIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	integrationCredential, err := ic.service.GetFullDetailsByID(uint(integrationID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, integrationCredential)

}

func (ic *IntegrationController) EnableIntegration(c *gin.Context) {
	integrationID, err := strconv.Atoi(c.Param("integrationId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var integrationCredentialFormDTO settingsDTOs.IntegrationCredentialFormDTO

	if err := c.ShouldBindJSON(&integrationCredentialFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ic.db, c, err.Error())
		return
	}

	id, err := ic.service.EnableIntegration(uint(integrationID), currentUserID, &integrationCredentialFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (ic *IntegrationController) UpdateCredential(c *gin.Context) {

	var integrationCredentialFormDTO settingsDTOs.IntegrationCredentialFormDTO

	integrationID, err := strconv.Atoi(c.Param("integrationId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&integrationCredentialFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ic.db, c, err.Error())
		return
	}

	id, err := ic.service.UpdateCredential(uint(integrationID), currentUserID, &integrationCredentialFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, commonDTOs.IDResponse{ID: id})
}

func (ic *IntegrationController) DisableIntegration(c *gin.Context) {

	integrationID, err := strconv.Atoi(c.Param("integrationId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ic.db, c, err.Error())
		return
	}

	id, err := ic.service.DisableIntegration(uint(integrationID), currentUserID)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, commonDTOs.IDResponse{ID: id})
}
