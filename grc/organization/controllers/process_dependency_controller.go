package controllers

import (
	processDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProcessDependencyController struct {
	service services.ProcessDependencyServiceInterface
	db      *gorm.DB
}

func NewProcessDependencyController(service services.ProcessDependencyServiceInterface) *ProcessDependencyController {
	return &ProcessDependencyController{
		service: service,
	}
}

func (p *ProcessDependencyController) GetAllProcessDependencies(c *gin.Context) {
	processIDStr := c.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	allowedParams := []string{
		"process_dependency_sub_type_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := processDTOs.ProcessDependencyListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Filters: filters,
	}
	result, err := p.service.GetAll(uint(processID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(p.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (p *ProcessDependencyController) GetProcessDependencyDetails(c *gin.Context) {
	processIDStr := c.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("dependencyId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	item, err := p.service.GetByID(uint(processID), uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (p *ProcessDependencyController) CreateProcessDependency(c *gin.Context) {
	processIDStr := c.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var dto *processDTOs.ProcessDependencyFormDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(p.db, c, err.Error())
		return
	}
	id, err := p.service.Create(uint(processID), currentUserID, dto)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(p.db, c, e.Error())
		}
		return
	}
	c.JSON(http.StatusCreated, id)
}

func (p *ProcessDependencyController) UpdateProcessDependency(c *gin.Context) {

	processIDStr := c.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	id, err := strconv.Atoi(c.Param("dependencyId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var dto *processDTOs.ProcessDependencyFormDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(p.db, c, err.Error())
		return
	}

	err = p.service.Update(uint(processID), uint(id), currentUserID, dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(p.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (p *ProcessDependencyController) DeleteProcessDependencies(c *gin.Context) {

	processIDStr := c.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(p.db, c, err.Error())
		return
	}

	failedIDs, err := p.service.Delete(uint(processID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(p.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
