package controllers

import (
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditWorkpaperController struct {
	service services.AuditWorkpaperServiceInterface
	db      *gorm.DB
}

func NewAuditWorkpaperController(service services.AuditWorkpaperServiceInterface) *AuditWorkpaperController {
	return &AuditWorkpaperController{
		service: service,
	}
}

func (awc *AuditWorkpaperController) GetAllAuditWorkpapers(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"audit_ids",
		"workpaper_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(awc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(awc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := dtos.AuditWorkpaperListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	auditWorkpapers, err := awc.service.GetAll(uint(auditID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(awc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditWorkpapers)
}

func (awc *AuditWorkpaperController) GetAuditWorkpaperDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditWorkpaper, err := awc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(awc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditWorkpaper)
}

func (awc *AuditWorkpaperController) CreateAuditWorkpaper(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var auditWorkpaperFormDTO dtos.AuditWorkpaperFormDTO

	if err := c.ShouldBindJSON(&auditWorkpaperFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(awc.db, c, err.Error())
		return
	}

	ids, err := awc.service.Create(uint(auditID), currentUserID, &auditWorkpaperFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(awc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"ids": ids})
}

func (awc *AuditWorkpaperController) UpdateAuditWorkpaper(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditWorkpaperID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var auditWorkpaperUpdateFormDTO dtos.AuditWorkpaperUpdateFormDTO

	if err := c.ShouldBindJSON(&auditWorkpaperUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(awc.db, c, err.Error())
		return
	}

	err = awc.service.Update(uint(auditWorkpaperID), uint(auditID), currentUserID, &auditWorkpaperUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(awc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (awc *AuditWorkpaperController) DeleteAuditWorkpapers(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(awc.db, c, err.Error())
		return
	}

	failedIDs, err := awc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(awc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (awc *AuditWorkpaperController) DownloadWorkpaperDocument(c *gin.Context) {
	auditIdStr := c.Param("auditId")
	auditId, err := strconv.ParseUint(auditIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	workpaperDocumentPath, err := awc.service.GetDocumentPath(uint(auditId), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(awc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(workpaperDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(workpaperDocumentPath)
}
