package controllers

import (
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentAssociatedItemController struct {
	service services.IncidentAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewIncidentAssociatedItemController(service services.IncidentAssociatedItemServiceInterface) *IncidentAssociatedItemController {
	return &IncidentAssociatedItemController{
		service: service,
	}
}

func (iaic *IncidentAssociatedItemController) GetAllIncidentControls(c *gin.Context) {

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	queryParams := incidentDTOs.IncidentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}
	incidents, err := iaic.service.GetAllControls(uint(incidentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidents)
}

func (iaic *IncidentAssociatedItemController) GetAllIncidentComplianceRequirements(c *gin.Context) {

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	queryParams := incidentDTOs.IncidentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}
	incidents, err := iaic.service.GetAllComplianceRequirements(uint(incidentID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func (iaic *IncidentAssociatedItemController) GetAllIncidentRisks(c *gin.Context) {

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	queryParams := incidentDTOs.IncidentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}
	incidents, err := iaic.service.GetAllRisks(uint(incidentID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func (iaic *IncidentAssociatedItemController) GetAllIncidentProcesses(c *gin.Context) {

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	queryParams := incidentDTOs.IncidentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}
	incidents, err := iaic.service.GetAllProcesses(uint(incidentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidents)
}
