package controllers

import (
	"errors"
	"jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectController struct {
	service services.ProjectServiceInterface
	db      *gorm.DB
}

func NewProjectController(service services.ProjectServiceInterface) *ProjectController {
	return &ProjectController{
		service: service,
	}
}

func (pc *ProjectController) GetAllProjects(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	status := c.Query("status")
	keyword := c.Query("keyword")

	// Support for quick filter
	quickFilter := c.QueryArray("quick_filter")

	// Support for multi-field filters
	allowedParams := []string{
		"project_type_ids",
		"project_status_ids",
		"responsible_user_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Handle pagination
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, "Invalid page_size parameter")
		return
	}

	// Define allowed sortable fields mapping
	validSortFields := map[string]string{
		"title_en":          "projects.title_en",
		"title_ar":          "projects.title_ar",
		"description_en":    "projects.description_en",
		"description_ar":    "projects.description_ar",
		"start_date":        "projects.start_date",
		"end_date":          "projects.end_date",
		"created_at":        "projects.created_at",
		"updated_at":        "projects.updated_at",
		"project_type_en":   "project_types.title_en",
		"project_type_ar":   "project_types.title_ar",
		"project_status_en": "project_statuses.title_en",
		"project_status_ar": "project_statuses.title_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	params := &dtos.ProjectListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Status:      status,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	projects, err := pc.service.GetAll(params)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (pc *ProjectController) GetProjectDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	project, err := pc.service.GetByID(uint(id))
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(pc.db, c, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, project)
}

func (pc *ProjectController) CreateProject(c *gin.Context) {
	var form dtos.ProjectFormDTO

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	id, err := pc.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(pc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (pc *ProjectController) UpdateProject(c *gin.Context) {
	var form dtos.ProjectFormDTO
	projectID, err := strconv.Atoi(c.Param("id"))
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
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	err = pc.service.Update(uint(projectID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(pc.db, c, e.Error())
		}
		return
	}
	c.Status(http.StatusOK)
}

func (pc *ProjectController) DeleteProjects(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	failedIDs, err := pc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
