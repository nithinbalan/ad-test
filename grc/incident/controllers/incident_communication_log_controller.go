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

type IncidentCommunicationLogController struct {
	service services.IncidentCommunicationLogServiceInterface
	db      *gorm.DB
}

func NewIncidentCommunicationLogController(service services.IncidentCommunicationLogServiceInterface) *IncidentCommunicationLogController {
	return &IncidentCommunicationLogController{
		service: service,
	}
}

func (iclc *IncidentCommunicationLogController) GetAllIncidentCommunicationLogs(c *gin.Context) {
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

	// Parse target_date_from, target_date_to
	communicatedAt := c.Query("communicated_at")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(iclc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(iclc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := incidentDTOs.IncidentCommunicationLogListQueryParamsDTO{
		Sort:           sort,
		SortBy:         sortBy,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		QuickFilter:    quickFilter,
		CommunicatedAt: communicatedAt,
	}
	incidentCommunicationLogs, err := iclc.service.GetAll(uint(incidentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(iclc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidentCommunicationLogs)
}

func (iclc *IncidentCommunicationLogController) GetIncidentCommunicationLogDetails(c *gin.Context) {
	incidentIDStr := c.Param("incidentID")
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

	incidentCommunicationLog, err := iclc.service.GetFullDetailsByID(uint(incidentID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(iclc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, incidentCommunicationLog)
}

func (iclc *IncidentCommunicationLogController) CreateIncidentCommunicationLog(c *gin.Context) {
	var IncidentCommunicationLogFormDTO incidentDTOs.IncidentCommunicationLogFormDTO

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&IncidentCommunicationLogFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iclc.db, c, err.Error())
		return
	}

	id, err := iclc.service.Create(uint(incidentID), currentUserID, &IncidentCommunicationLogFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iclc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (iclc *IncidentCommunicationLogController) UpdateIncidentCommunicationLog(c *gin.Context) {
	var IncidentCommunicationLogFormDTO incidentDTOs.IncidentCommunicationLogFormDTO

	incidentIDStr := c.Param("incidentID")

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

	if err := c.ShouldBindJSON(&IncidentCommunicationLogFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iclc.db, c, err.Error())
		return
	}

	err = iclc.service.Update(uint(incidentID), uint(id), currentUserID, &IncidentCommunicationLogFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iclc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (iclc *IncidentCommunicationLogController) DeleteIncidentCommunicationLogs(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iclc.db, c, err.Error())
		return
	}

	failedIDs, err := iclc.service.Delete(uint(incidentID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iclc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
