package controllers

import (
	"net/http"
	"strconv"

	assetDTOs "jethur/grc/asset/dtos"
	assetServices "jethur/grc/asset/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SystemRegisterAssociatedItemController struct {
	db      *gorm.DB
	service assetServices.SystemRegisterAssociatedItemServiceInterface
}

func NewSystemRegisterAssociatedItemController(service assetServices.SystemRegisterAssociatedItemServiceInterface) *SystemRegisterAssociatedItemController {
	return &SystemRegisterAssociatedItemController{
		service: service,
	}
}

func (sraic *SystemRegisterAssociatedItemController) GetAllSystemRegisterSLAs(c *gin.Context) {

	systemRegisterIDStr := c.Param("id")

	systemRegisterID, err := strconv.ParseUint(systemRegisterIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Validate sort_order to prevent SQL injection
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(sraic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(sraic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := assetDTOs.SystemRegisterAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}

	slas, err := sraic.service.GetAllSLAs(uint(systemRegisterID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sraic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, slas)
}

func (sraic *SystemRegisterAssociatedItemController) GetAllSystemRegisterProcesses(c *gin.Context) {

	systemRegisterIDStr := c.Param("id")

	systemRegisterID, err := strconv.ParseUint(systemRegisterIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	// Validate sort_order to prevent SQL injection
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(sraic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(sraic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	processes, err := sraic.service.GetAllProcesses(uint(systemRegisterID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sraic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, processes)
}

func (sraic *SystemRegisterAssociatedItemController) CreateSystemRegisterSLAs(c *gin.Context) {

	var formDTO assetDTOs.SystemRegisterSLAFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	systemRegisterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sraic.db, c, err.Error())
		return
	}

	err = sraic.service.CreateSLAs(uint(systemRegisterID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sraic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (sraic *SystemRegisterAssociatedItemController) CreateSystemRegisterProcesses(c *gin.Context) {

	var formDTO assetDTOs.SystemRegisterProcessFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	systemRegisterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sraic.db, c, err.Error())
		return
	}

	err = sraic.service.CreateProcesses(uint(systemRegisterID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sraic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (sraic *SystemRegisterAssociatedItemController) DeleteSystemRegisterSLA(c *gin.Context) {

	systemRegisterID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	slaID, _ := strconv.ParseUint(c.Param("slaId"), 10, 64)

	err := sraic.service.DeleteSLA(uint(systemRegisterID), uint(slaID))
	if err != nil {
		utility.SendInternalErrResp(sraic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SLA unmapped from system register successfully"})
}

func (sraic *SystemRegisterAssociatedItemController) DeleteSystemRegisterProcess(c *gin.Context) {

	systemRegisterID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	processID, _ := strconv.ParseUint(c.Param("processId"), 10, 64)

	err := sraic.service.DeleteProcess(uint(systemRegisterID), uint(processID))
	if err != nil {
		utility.SendInternalErrResp(sraic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "process unmapped from system register successfully"})
}
