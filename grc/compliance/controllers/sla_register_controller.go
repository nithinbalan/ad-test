package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	slaDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SLARegisterController struct {
	service services.SLARegisterServiceInterface
	db      *gorm.DB
}

func NewSLARegisterController(service services.SLARegisterServiceInterface) *SLARegisterController {
	return &SLARegisterController{
		service: service,
	}
}

func (src *SLARegisterController) GetAllSLARegisters(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"sla_category_ids",
		"severity_level_ids",
		"responsible_unit_ids",
		"lifecycle_status_ids",
		"compliance_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse  issue_date_from, issue_date_to, expiry_date_from, expiry_date_to parameters
	issueDateFrom := c.Query("issue_date_from")
	issueDateTo := c.Query("issue_date_to")
	expiryDateFrom := c.Query("expiry_date_from")
	expiryDateTo := c.Query("expiry_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := slaDTOs.SLARegisterListQueryParamsDTO{
		Sort:           sort,
		SortBy:         sortBy,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		Filters:        filters,
		QuickFilter:    quickFilter,
		IssueDateFrom:  issueDateFrom,
		IssueDateTo:    issueDateTo,
		ExpiryDateFrom: expiryDateFrom,
		ExpiryDateTo:   expiryDateTo,
	}
	slaRegisters, err := src.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, slaRegisters)

}

func (src *SLARegisterController) GetSLARegisterDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sla, err := src.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, sla)
}

func (src *SLARegisterController) CreateSLARegister(c *gin.Context) {

	var slaRegFormDTO slaDTOs.SLARegisterFormDTO

	if err := c.ShouldBindJSON(&slaRegFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	id, err := src.service.Create(currentUserID, &slaRegFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(src.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (src *SLARegisterController) UpdateSLARegister(c *gin.Context) {
	slaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var slaRegFormDTO slaDTOs.SLARegisterFormDTO
	if err := c.ShouldBindJSON(&slaRegFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}
	err = src.service.Update(uint(slaID), currentUserID, &slaRegFormDTO)
	// @Review
	// We have to manage file upload failure error message,
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(src.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (src *SLARegisterController) DeleteSLARegisters(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	failedIDs, err := src.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (src *SLARegisterController) ArchiveSLARegisters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}
	successfulIds, err := src.service.Archive(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"successfulIds": successfulIds})
}

func (src *SLARegisterController) RestoreSLARegisters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	successfulIds, err := src.service.Restore(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"successfulIds": successfulIds})
}

