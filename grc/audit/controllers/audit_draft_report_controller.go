package controllers

import (
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditDraftReportController struct {
	service services.AuditDraftReportServiceInterface
	db      *gorm.DB
}

func NewAuditDraftReportController(service services.AuditDraftReportServiceInterface) *AuditDraftReportController {
	return &AuditDraftReportController{
		service: service,
	}
}

func (adrc *AuditDraftReportController) GetAllAuditDraftReports(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	allowedParams := []string{
		"audit_id",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(adrc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(adrc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := dtos.AuditDraftReportListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	auditDraftReports, err := adrc.service.GetAll(uint(auditID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(adrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditDraftReports)
}

func (adrc *AuditDraftReportController) GetAuditDraftReportDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditDraftReport, err := adrc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(adrc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditDraftReport)
}

func (adrc *AuditDraftReportController) CreateAuditDraftReport(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditDraftReportFormDTO dtos.AuditDraftReportFormDTO

	if err := c.ShouldBindJSON(&auditDraftReportFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrc.db, c, err.Error())
		return
	}

	id, err := adrc.service.Create(uint(auditID), currentUserID, &auditDraftReportFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (adrc *AuditDraftReportController) UpdateAuditDraftReport(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditDraftReportFormDTO dtos.AuditDraftReportFormDTO

	auditDraftReportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditDraftReportFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrc.db, c, err.Error())
		return
	}

	err = adrc.service.Update(uint(auditDraftReportID), uint(auditID), currentUserID, &auditDraftReportFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (adrc *AuditDraftReportController) DeleteAuditDraftReports(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrc.db, c, err.Error())
		return
	}

	failedIDs, err := adrc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(adrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
