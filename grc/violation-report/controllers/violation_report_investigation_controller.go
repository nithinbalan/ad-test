package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	"jethur/grc/violation-report/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ViolationReportInvestigationController struct {
	service services.ViolationReportInvestigationServiceInterface
	db      *gorm.DB
}

func NewViolationReportInvestigationController(service services.ViolationReportInvestigationServiceInterface) *ViolationReportInvestigationController {
	return &ViolationReportInvestigationController{
		service: service,
	}
}

func (iic *ViolationReportInvestigationController) GetAllViolationReportInvestigations(c *gin.Context) {
	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
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
		"department_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	investigationDate := c.Query("investigation_date")

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

	queryParams := violationReportDTOs.ViolationReportInvestigationListQueryParamsDTO{
		Sort:              sort,
		SortBy:            sortBy,
		Keyword:           keyword,
		Status:            status,
		Page:              page,
		PageSize:          pageSize,
		Filters:           filters,
		QuickFilter:       quickFilter,
		InvestigationDate: investigationDate,
	}
	violationReportInvestigations, err := iic.service.GetAll(uint(violationReportID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, violationReportInvestigations)
}

func (iic *ViolationReportInvestigationController) GetViolationReportInvestigationDetails(c *gin.Context) {
	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	violationReportInvestigation, err := iic.service.GetFullDetailsByID(uint(violationReportID))
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
	c.JSON(http.StatusOK, violationReportInvestigation)
}

func (iic *ViolationReportInvestigationController) CreateViolationReportInvestigation(c *gin.Context) {
	var violationReportInvestigationFormDTO violationReportDTOs.ViolationReportInvestigationFormDTO

	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&violationReportInvestigationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}

	id, err := iic.service.Create(uint(violationReportID), currentUserID, &violationReportInvestigationFormDTO)

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

func (iic *ViolationReportInvestigationController) UpdateViolationReportInvestigation(c *gin.Context) {
	var violationReportInvestigationFormDTO violationReportDTOs.ViolationReportInvestigationFormDTO

	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&violationReportInvestigationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}

	err = iic.service.Update(uint(violationReportID), currentUserID, &violationReportInvestigationFormDTO)
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

func (iic *ViolationReportInvestigationController) DeleteViolationReportInvestigations(c *gin.Context) {
	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}

	failedID, err := iic.service.Delete(uint(violationReportID), currentUserID)
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

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}

func (iasuc *ViolationReportInvestigationController) GetViolationReportInvestigationDocumentPreview(c *gin.Context) {

	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	violationReportInvestigationDocumentPath, err := iasuc.service.GetDocumentPath(uint(violationReportID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(iasuc.db, c, err.Error())
		}
		return
	}
	c.File(violationReportInvestigationDocumentPath)
}

func (iasuc *ViolationReportInvestigationController) DownloadViolationReportInvestigationDocument(c *gin.Context) {
	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	violationReportInvestigationDocumentPath, err := iasuc.service.GetDocumentPath(uint(violationReportID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(iasuc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(violationReportInvestigationDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(violationReportInvestigationDocumentPath)
}
