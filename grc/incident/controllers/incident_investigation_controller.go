package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentInvestigationController struct {
	service services.IncidentInvestigationServiceInterface
	db      *gorm.DB
}

func NewIncidentInvestigationController(service services.IncidentInvestigationServiceInterface) *IncidentInvestigationController {
	return &IncidentInvestigationController{
		service: service,
	}
}

func (iic *IncidentInvestigationController) GetAllIncidentInvestigations(c *gin.Context) {
	incidentIDStr := c.Param("incidentId")

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

	// Define an array of allowed parameter names
	allowedParams := []string{
		"user_ids",
		"team_lead_ids",
		"investigation_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	initiatedAt := c.Query("initiated_at")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := incidentDTOs.IncidentInvestigationListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		InitiatedAt: initiatedAt,
	}
	incidentInvestigations, err := iic.service.GetAll(uint(incidentID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidentInvestigations)
}

func (iic *IncidentInvestigationController) GetIncidentInvestigationDetails(c *gin.Context) {
	incidentIDStr := c.Param("incidentId")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentInvestigation, err := iic.service.GetFullDetailsByID(uint(incidentID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(iic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidentInvestigation)
}

func (iic *IncidentInvestigationController) CreateIncidentInvestigation(c *gin.Context) {
	var incidentInvestigationFormDTO incidentDTOs.IncidentInvestigationFormDTO

	incidentIDStr := c.Param("incidentId")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&incidentInvestigationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}

	id, err := iic.service.Create(uint(incidentID), currentUserID, &incidentInvestigationFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (iic *IncidentInvestigationController) UpdateIncidentInvestigation(c *gin.Context) {
	var incidentInvestigationFormDTO incidentDTOs.IncidentInvestigationFormDTO

	incidentIDStr := c.Param("incidentId")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentInvestigationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&incidentInvestigationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}

	err = iic.service.Update(uint(incidentID), uint(incidentInvestigationID), currentUserID, &incidentInvestigationFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (iic *IncidentInvestigationController) DeleteIncidentInvestigations(c *gin.Context) {
	incidentIDStr := c.Param("incidentId")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}

	failedIDs, err := iic.service.Delete(uint(incidentID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
