package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GRCReportController struct {
	service services.GRCReportServiceInterface
	db      *gorm.DB
}

func NewGRCReportController(service services.GRCReportServiceInterface) *GRCReportController {
	return &GRCReportController{
		service: service,
	}
}

func (orc *GRCReportController) GetAllGRCReports(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"title":      "grc_reports.title",
		"start_date": "grc_reports.start_date",
		"end_date":   "grc_reports.end_date",
		"created_at": "grc_reports.created_at",
	}

	filters := utility.ParseQueryFilters([]string{}, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := dtos.GRCReportListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	reports, err := orc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (orc *GRCReportController) GetGRCReportDetails(c *gin.Context) {
	grcReportID, err := strconv.Atoi(c.Param("grcReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := orc.service.GetGRCReportDetails(uint(grcReportID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(orc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, details)
}

func (orc *GRCReportController) CreateGRCReport(c *gin.Context) {
	var request dtos.GRCReportFormDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	grcReportID, err := orc.service.Create(currentUserID, &request)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(orc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "GRC report generated successfully",
		"grc_report_id": grcReportID,
	})
}

func (orc *GRCReportController) GetGRCReportSection(c *gin.Context) {
	grcReportID, err := strconv.Atoi(c.Param("grcReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := orc.service.GetSectionByID(uint(grcReportID), uint(sectionID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(orc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, details)
}

func (orc *GRCReportController) UpdateGRCReportSection(c *gin.Context) {
	grcReportID, err := strconv.Atoi(c.Param("grcReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.GRCReportSectionFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	err = orc.service.UpdateSection(uint(grcReportID), uint(sectionID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(orc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (orc *GRCReportController) DeleteGRCReportSection(c *gin.Context) {
	grcReportID, err := strconv.Atoi(c.Param("grcReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	err = orc.service.DeleteGRCReportSection(uint(grcReportID), uint(sectionID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(orc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (orc *GRCReportController) DeleteGRCReport(c *gin.Context) {
	grcReportID, err := strconv.Atoi(c.Param("grcReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	err = orc.service.DeleteGRCReport(uint(grcReportID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(orc.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (orc *GRCReportController) UpdateHighlights(c *gin.Context) {
	grcReportID, err := strconv.Atoi(c.Param("grcReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.GRCReportHighlightChallengeFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	err = orc.service.UpdateHighlights(uint(grcReportID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(orc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (orc *GRCReportController) UpdateChallenges(c *gin.Context) {
	grcReportID, err := strconv.Atoi(c.Param("grcReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.GRCReportHighlightChallengeFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	err = orc.service.UpdateChallenges(uint(grcReportID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(orc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (orc *GRCReportController) UpdateRiskApprovalNote(c *gin.Context) {
	grcReportID, err := strconv.Atoi(c.Param("grcReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskID, err := strconv.Atoi(c.Param("riskId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.GRCReportHighlightChallengeFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	err = orc.service.UpdateRiskApprovalNote(uint(grcReportID), uint(riskID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(orc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
