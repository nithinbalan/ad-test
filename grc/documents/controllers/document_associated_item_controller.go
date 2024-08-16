package controllers

import (
	documentDTOs "jethur/grc/documents/dtos"
	documentService "jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentAssociatedItemController struct {
	service documentService.DocumentAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewDocumentAssociatedItemController(service documentService.DocumentAssociatedItemServiceInterface) *DocumentAssociatedItemController {
	return &DocumentAssociatedItemController{
		service: service,
	}
}

func (daic *DocumentAssociatedItemController) GetAllDocumentControls(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := documentDTOs.DocumentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	controls, err := daic.service.GetAllControls(uint(documentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}

func (daic *DocumentAssociatedItemController) GetAllDocumentComplianceRequirements(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := documentDTOs.DocumentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	complianceRequirements, err := daic.service.GetAllComplianceRequirements(uint(documentID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, complianceRequirements)
}

func (daic *DocumentAssociatedItemController) GetAllDocumentProcesses(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := documentDTOs.DocumentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	processes, err := daic.service.GetAllProcesses(uint(documentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, processes)
}

func (daic *DocumentAssociatedItemController) GetAllDocumentSLAs(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := documentDTOs.DocumentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	slas, err := daic.service.GetAllSLAs(uint(documentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, slas)
}

func (daic *DocumentAssociatedItemController) GetAllDocumentIncidents(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := documentDTOs.DocumentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	incidents, err := daic.service.GetAllIncidents(uint(documentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidents)
}

func (daic *DocumentAssociatedItemController) GetAllDocumentRisks(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := documentDTOs.DocumentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	risks, err := daic.service.GetAllRisks(uint(documentID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, risks)
}
