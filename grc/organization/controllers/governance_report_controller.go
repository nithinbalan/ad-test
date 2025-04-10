package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GovernanceReportController struct {
	service services.GovernanceReportServiceInterface
	db      *gorm.DB
}

func NewGovernanceReportController(service services.GovernanceReportServiceInterface) *GovernanceReportController {
	return &GovernanceReportController{
		service: service,
	}
}

func (orc *GovernanceReportController) GetAllGovernanceReports(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"title":      "governance_reports.title",
		"start_date": "governance_reports.start_date",
		"end_date":   "governance_reports.end_date",
		"created_at": "governance_reports.created_at",
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

	queryParams := dtos.GovernanceReportListQueryParamsDTO{
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

func (orc *GovernanceReportController) GetGovernanceReportDetails(c *gin.Context) {
	governanceReportID, err := strconv.Atoi(c.Param("governanceReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := orc.service.GetGovernanceReportDetails(uint(governanceReportID))
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

func (orc *GovernanceReportController) CreateGovernanceReport(c *gin.Context) {
	var request dtos.GovernanceReportFormDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	governanceReportID, err := orc.service.Create(currentUserID, &request)
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
		"message":              "Governance report generated successfully",
		"governance_report_id": governanceReportID,
	})
}

func (orc *GovernanceReportController) GetGovernanceReportSection(c *gin.Context) {
	governanceReportID, err := strconv.Atoi(c.Param("governanceReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := orc.service.GetSectionByID(uint(governanceReportID), uint(sectionID))
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

func (orc *GovernanceReportController) UpdateGovernanceReportSection(c *gin.Context) {
	governanceReportID, err := strconv.Atoi(c.Param("governanceReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.GovernanceReportSectionFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	err = orc.service.UpdateSection(uint(governanceReportID), uint(sectionID), currentUserID, &formDTO)
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

func (orc *GovernanceReportController) DeleteGovernanceReportSection(c *gin.Context) {
	governanceReportID, err := strconv.Atoi(c.Param("governanceReportId"))
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

	err = orc.service.DeleteGovernanceReportSection(uint(governanceReportID), uint(sectionID), currentUserID)
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

func (orc *GovernanceReportController) DeleteGovernanceReport(c *gin.Context) {
	governanceReportID, err := strconv.Atoi(c.Param("governanceReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	err = orc.service.DeleteGovernanceReport(uint(governanceReportID), currentUserID)
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
