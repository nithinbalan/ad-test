package controllers

import (
	"fmt"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EscalationController struct {
	service services.EscalationServiceInterface
	db      *gorm.DB
}

func NewEscalationController(service services.EscalationServiceInterface) *EscalationController {
	return &EscalationController{
		service: service,
	}
}

func (rc *EscalationController) GetAllEscalations(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	allowedParams := []string{
		"module_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := dtos.EscalationListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}

	escalations, err := rc.service.GetAll(queryParams)

	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, escalations)
}

func (rc *EscalationController) GetEscalationDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	fmt.Println("1")
	escalation, err := rc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, escalation)
}

func (rc *EscalationController) GetAllEnabledModules(c *gin.Context) {
	modules, err := rc.service.GetAllEnabledModules()
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, modules)
}

func (rc *EscalationController) GetEscalationsByModuleID(c *gin.Context) {
	moduleIDStr := c.Param("moduleId")

	moduleID, err := strconv.ParseUint(moduleIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	stakeholderStatuses, err := rc.service.GetAllEscalationsByModuleID(uint(moduleID))

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, stakeholderStatuses)
}
