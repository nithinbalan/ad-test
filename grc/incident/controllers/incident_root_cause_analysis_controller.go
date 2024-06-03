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

type IncidentRootCauseAnalysisController struct {
	service services.IncidentRootCauseAnalysisServiceInterface
	db      *gorm.DB
}

func NewIncidentRootCauseAnalysisController(service services.IncidentRootCauseAnalysisServiceInterface) *IncidentRootCauseAnalysisController {
	return &IncidentRootCauseAnalysisController{
		service: service,
	}
}

func (ircac *IncidentRootCauseAnalysisController) GetAllIncidentRootCauseAnalyses(c *gin.Context) {
	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	queryParams := incidentDTOs.IncidentRootCauseAnalysisListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
	}
	incidentRootCauseAnalyses, err := ircac.service.GetAll(uint(incidentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ircac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidentRootCauseAnalyses)
}

func (ircac *IncidentRootCauseAnalysisController) CreateIncidentRootCauseAnalysis(c *gin.Context) {
	var IncidentRootCauseAnalysisFormDTO incidentDTOs.IncidentRootCauseAnalysisFormDTO

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&IncidentRootCauseAnalysisFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ircac.db, c, err.Error())
		return
	}

	id, err := ircac.service.Create(uint(incidentID), currentUserID, &IncidentRootCauseAnalysisFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ircac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (ircac *IncidentRootCauseAnalysisController) UpdateIncidentRootCauseAnalysis(c *gin.Context) {
	var IncidentRootCauseAnalysisFormDTO incidentDTOs.IncidentRootCauseAnalysisFormDTO

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&IncidentRootCauseAnalysisFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ircac.db, c, err.Error())
		return
	}

	err = ircac.service.Update(uint(incidentID), currentUserID, &IncidentRootCauseAnalysisFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ircac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ircac *IncidentRootCauseAnalysisController) DeleteIncidentRootCauseAnalyses(c *gin.Context) {
	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ircac.db, c, err.Error())
		return
	}

	failedID, err := ircac.service.Delete(uint(incidentID), uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ircac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
