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
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "planned_date")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"project_id",
		"project_milestone_status_id",
		"created_by_id",
		"updated_by_id",
	}

	validSortFields := map[string]string{
		"id":           "project_milestones.id",
		"title_en":     "project_milestones.title_en",
		"title_ar":     "project_milestones.title_ar",
		"planned_date": "project_milestones.planned_date",
		"actual_date":  "project_milestones.actual_date",
		"status":       "project_milestones.project_milestone_status_id",
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

	milestones, err := pmc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(pmc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, milestones)
}

func (pmc *ProjectMilestoneController) GetProjectMilestoneDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	milestone, err := pmc.service.GetFullDetailsByID(uint(id))
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
	c.JSON(http.StatusOK, milestone)
}

func (pmc *ProjectMilestoneController) CreateProjectMilestone(c *gin.Context) {
	var milestoneFormDTO projectDTOs.ProjectMilestoneFormDTO
	if err := c.ShouldBindJSON(&milestoneFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pmc.db, c, err.Error())
		return
	}
	id, err := pmc.service.Create(currentUserID, &milestoneFormDTO)
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
	milestoneID, err := strconv.Atoi(c.Param("id"))
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
	err = pmc.service.Update(uint(milestoneID), currentUserID, &milestoneFormDTO)
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

func (pmc *ProjectMilestoneController) DeleteProjectMilestones(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
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
	failedIDs, err := pmc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pmc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
