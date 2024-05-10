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

type IncidentRegisterController struct {
	service services.IncidentRegisterServiceInterface
	db      *gorm.DB
}

func NewIncidentRegisterController(service services.IncidentRegisterServiceInterface) *IncidentRegisterController {
	return &IncidentRegisterController{
		service: service,
	}
}

func (irc *IncidentRegisterController) GetAllIncidentRegisters(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"incident_category_ids",
		"incident_sub_category_ids",
		"severity_level_ids",
		"incident_status_ids",
		"reported_user_ids",
		"incident_type_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	incidentAt := c.Query("incident_at")
	reportedAt := c.Query("reported")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := incidentDTOs.IncidentRegisterListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		IncidentAt:  incidentAt,
		ReportedAt:  reportedAt,
	}
	incidents, err := irc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func (irc *IncidentRegisterController) GetIncidentRegisterDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incident, err := irc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(irc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incident)
}

func (irc *IncidentRegisterController) CreateIncidentRegister(c *gin.Context) {
	var incidentRegisterFormDTO incidentDTOs.IncidentRegisterFormDTO

	if err := c.ShouldBindJSON(&incidentRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	id, err := irc.service.Create(currentUserID, &incidentRegisterFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(irc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (irc *IncidentRegisterController) UpdateIncidentRegister(c *gin.Context) {
	var incidentRegisterFormDTO incidentDTOs.IncidentRegisterFormDTO

	incidentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&incidentRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	err = irc.service.Update(uint(incidentID), currentUserID, &incidentRegisterFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(irc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (irc *IncidentRegisterController) DeleteIncidentRegisters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	failedIDs, err := irc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
