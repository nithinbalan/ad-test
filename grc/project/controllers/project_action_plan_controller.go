package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectActionPlanController struct {
	service services.ProjectActionPlanServiceInterface
	db      *gorm.DB
}

func NewProjectActionPlanController(service services.ProjectActionPlanServiceInterface) *ProjectActionPlanController {
	return &ProjectActionPlanController{
		service: service,
	}
}

func (papc *ProjectActionPlanController) GetAllProjectActionPlans(c *gin.Context) {
	projectIDStr := c.Param("projectId")

	var projectID uint64
	var err error

	// Check if projectID is provided in the URL
	if projectIDStr != "" {
		projectID, err = strconv.ParseUint(projectIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		projectID = 0
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"action_plan_status_ids",
		"responsible_user_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(papc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(papc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := projectDTOs.ProjectActionPlanListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	projectActionPlans, err := papc.service.GetAll(uint(projectID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(papc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, projectActionPlans)
}

func (papc *ProjectActionPlanController) GetProjectActionPlanDetails(c *gin.Context) {
	projectIDStr := c.Param("projectId")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	projectActionPlan, err := papc.service.GetFullDetailsByID(uint(projectID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(papc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, projectActionPlan)
}

func (papc *ProjectActionPlanController) CreateProjectActionPlan(c *gin.Context) {
	var ProjectActionPlanFormDTO projectDTOs.ProjectActionPlanFormDTO

	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&ProjectActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(papc.db, c, err.Error())
		return
	}

	id, err := papc.service.Create(uint(projectID), currentUserID, &ProjectActionPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(papc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (papc *ProjectActionPlanController) UpdateProjectActionPlan(c *gin.Context) {
	var ProjectActionPlanFormDTO projectDTOs.ProjectActionPlanFormDTO

	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	id, err := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&ProjectActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(papc.db, c, err.Error())
		return
	}

	err = papc.service.Update(uint(projectID), uint(id), currentUserID, &ProjectActionPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(papc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (papc *ProjectActionPlanController) DeleteProjectActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(papc.db, c, err.Error())
		return
	}

	failedIDs, err := papc.service.Delete(uint(projectID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(papc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
