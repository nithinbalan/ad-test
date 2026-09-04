package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectMilestoneController struct {
	service services.ProjectMilestoneServiceInterface
	db      *gorm.DB
}

func NewProjectMilestoneController(service services.ProjectMilestoneServiceInterface) *ProjectMilestoneController {
	return &ProjectMilestoneController{
		service: service,
	}
}

func (pmc *ProjectMilestoneController) GetAllProjectMilestones(c *gin.Context) {
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
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "planned_date")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"project_ids",
		"project_milestone_status_ids",
		"project_phase_ids",
		"created_by_ids",
		"updated_by_ids",
	}

	validSortFields := map[string]string{
		"id":                     "project_milestones.id",
		"title_en":               "project_milestones.title_en",
		"title_ar":               "project_milestones.title_ar",
		"planned_date":           "project_milestones.planned_date",
		"actual_date":            "project_milestones.actual_date",
		"status":                 "project_milestones.project_milestone_status_id",
		"project_phase_title_en": "phase.title_en",
		"project_phase_title_ar": "phase.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(pmc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(pmc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := projectDTOs.ProjectMilestoneListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}

	projectMilestones, err := pmc.service.GetAll(uint(projectID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(pmc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, projectMilestones)
}

func (pmc *ProjectMilestoneController) GetProjectMilestoneDetails(c *gin.Context) {
	projectIDStr := c.Param("projectId")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("milestoneId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	projectMilestone, err := pmc.service.GetFullDetailsByID(uint(projectID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(pmc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, projectMilestone)
}

func (pmc *ProjectMilestoneController) CreateProjectMilestone(c *gin.Context) {
	var milestoneFormDTO projectDTOs.ProjectMilestoneFormDTO
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	if err := c.ShouldBindJSON(&milestoneFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pmc.db, c, err.Error())
		return
	}
	id, err := pmc.service.Create(uint(projectID), currentUserID, &milestoneFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pmc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (pmc *ProjectMilestoneController) UpdateProjectMilestone(c *gin.Context) {
	var milestoneFormDTO projectDTOs.ProjectMilestoneFormDTO
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	milestoneID, err := strconv.Atoi(c.Param("milestoneId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	if err := c.ShouldBindJSON(&milestoneFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pmc.db, c, err.Error())
		return
	}
	err = pmc.service.Update(uint(projectID), uint(milestoneID), currentUserID, &milestoneFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pmc.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (pmc *ProjectMilestoneController) UpdateProjectMilestoneStatus(c *gin.Context) {
	var form projectDTOs.ProjectMilestoneStatusUpdateDTO
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	milestoneID, err := strconv.ParseUint(c.Param("milestoneId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pmc.db, c, err.Error())
		return
	}
	err = pmc.service.UpdateStatus(uint(projectID), uint(milestoneID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(pmc.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (pmc *ProjectMilestoneController) DeleteProjectMilestones(c *gin.Context) {
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
		utility.SendInternalErrResp(pmc.db, c, err.Error())
		return
	}
	failedIDs, err := pmc.service.Delete(uint(projectID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pmc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
