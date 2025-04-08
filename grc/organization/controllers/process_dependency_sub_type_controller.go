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

type ProcessDependencySubTypeController struct {
	service services.ProcessDependencySubTypeServiceInterface
	db      *gorm.DB
}

func NewProcessDependencySubTypeController(service services.ProcessDependencySubTypeServiceInterface) *ProcessDependencySubTypeController {
	return &ProcessDependencySubTypeController{
		service: service,
	}
}

func (p *ProcessDependencySubTypeController) GetAllProcessDependencySubTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	allowedParams := []string{
		"process_dependency_type_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := processDTOs.ProcessDependencySubTypeListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}
	result, err := p.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(p.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (p *ProcessDependencySubTypeController) GetProcessDependencySubTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	item, err := p.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (p *ProcessDependencySubTypeController) CreateProcessDependencySubType(c *gin.Context) {
	var dto *processDTOs.ProcessDependencySubTypeFormDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(p.db, c, err.Error())
		return
	}
	id, err := p.service.Create(currentUserID, dto)

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

func (p *ProcessDependencySubTypeController) UpdateProcessDependencySubType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var dto *processDTOs.ProcessDependencySubTypeFormDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(p.db, c, err.Error())
		return
	}

	err = p.service.Update(uint(id), currentUserID, dto)
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

func (p *ProcessDependencySubTypeController) DeleteProcessDependencySubTypes(c *gin.Context) {
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

	failedIDs, err := p.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(p.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (p *ProcessDependencySubTypeController) ActivateProcessDependencySubTypes(c *gin.Context) {
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

	err = p.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(p.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (p *ProcessDependencySubTypeController) DeactivateProcessDependencySubTypes(c *gin.Context) {
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

	err = p.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(p.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
