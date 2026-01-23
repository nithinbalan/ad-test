package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	"jethur/grc/violation-report/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ViolationReportRootCauseAnalysisController struct {
	service services.ViolationReportRootCauseAnalysisServiceInterface
	db      *gorm.DB
}

func NewViolationReportRootCauseAnalysisController(service services.ViolationReportRootCauseAnalysisServiceInterface) *ViolationReportRootCauseAnalysisController {
	return &ViolationReportRootCauseAnalysisController{
		service: service,
	}
}

func (ircac *ViolationReportRootCauseAnalysisController) GetAllViolationReportRootCauseAnalyses(c *gin.Context) {
	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	queryParams := violationReportDTOs.ViolationReportRootCauseAnalysisListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
	}
	violationReportRootCauseAnalyses, err := ircac.service.GetAll(uint(violationReportID), &queryParams)

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
	c.JSON(http.StatusOK, violationReportRootCauseAnalyses)
}

func (ircac *ViolationReportRootCauseAnalysisController) CreateViolationReportRootCauseAnalysis(c *gin.Context) {
	var violationReportRootCauseAnalysisFormDTO violationReportDTOs.ViolationReportRootCauseAnalysisFormDTO

	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&violationReportRootCauseAnalysisFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ircac.db, c, err.Error())
		return
	}

	id, err := ircac.service.Create(uint(violationReportID), currentUserID, &violationReportRootCauseAnalysisFormDTO)

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

func (ircac *ViolationReportRootCauseAnalysisController) UpdateViolationReportRootCauseAnalysis(c *gin.Context) {
	var violationReportRootCauseAnalysisFormDTO violationReportDTOs.ViolationReportRootCauseAnalysisFormDTO

	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("rootCauseAnalysisID")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&violationReportRootCauseAnalysisFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ircac.db, c, err.Error())
		return
	}

	err = ircac.service.Update(uint(violationReportID), uint(id), currentUserID, &violationReportRootCauseAnalysisFormDTO)
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

func (ircac *ViolationReportRootCauseAnalysisController) DeleteViolationReportRootCauseAnalyses(c *gin.Context) {
	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("rootCauseAnalysisID")

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

	failedID, err := ircac.service.Delete(uint(violationReportID), uint(id), currentUserID)
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
