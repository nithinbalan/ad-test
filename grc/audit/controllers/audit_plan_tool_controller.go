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

type AuditPlanToolController struct {
	service services.AuditPlanToolServiceInterface
	db      *gorm.DB
}

func NewAuditPlanToolController(service services.AuditPlanToolServiceInterface) *AuditPlanToolController {
	return &AuditPlanToolController{
		service: service,
	}
}

func (aptc *AuditPlanToolController) GetAllAuditPlanTools(c *gin.Context) {
	auditPlanID, err := strconv.Atoi(c.Param("planId"))

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"audit_plan_ids",
		"tool_type_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := dtos.AuditPlanToolListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	auditPlanTools, err := aptc.service.GetAll(uint(auditPlanID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(aptc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditPlanTools)
}

func (aptc *AuditPlanToolController) GetAuditPlanToolDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditPlanTool, err := aptc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aptc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditPlanTool)
}

func (aptc *AuditPlanToolController) CreateAuditPlanTool(c *gin.Context) {
	auditPlanID, err := strconv.Atoi(c.Param("planId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditPlanToolFormDTO dtos.AuditPlanToolFormDTO

	if err := c.ShouldBindJSON(&auditPlanToolFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aptc.db, c, err.Error())
		return
	}

	id, err := aptc.service.Create(uint(auditPlanID), currentUserID, &auditPlanToolFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aptc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (aptc *AuditPlanToolController) UpdateAuditPlanTool(c *gin.Context) {
	auditPlanID, err := strconv.Atoi(c.Param("planId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditPlanToolFormDTO dtos.AuditPlanToolFormDTO

	auditPlanToolID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditPlanToolFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aptc.db, c, err.Error())
		return
	}

	err = aptc.service.Update(uint(auditPlanToolID), uint(auditPlanID), currentUserID, &auditPlanToolFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aptc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (aptc *AuditPlanToolController) DeleteAuditPlanTools(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aptc.db, c, err.Error())
		return
	}

	failedIDs, err := aptc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aptc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
