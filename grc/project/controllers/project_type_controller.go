package controllers

import (
	"jethur/grc/models"
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectTypeController struct {
	service services.ProjectTypeServiceInterface
	db      *gorm.DB
}

func NewProjectTypeController(service services.ProjectTypeServiceInterface) *ProjectTypeController {
	return &ProjectTypeController{
		service: service,
	}
}

func (ptc *ProjectTypeController) GetAllProjectTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	projectTypes, err := ptc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(ptc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, projectTypes)
}

func (ptc *ProjectTypeController) GetProjectTypeDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	projectType, err := ptc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, projectType)
}

func (ptc *ProjectTypeController) CreateProjectType(c *gin.Context) {
	var pt *models.ProjectType

	if err := c.ShouldBindJSON(&pt); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptc.db, c, err.Error())
		return
	}

	createdProjectType, err := ptc.service.Create(currentUserID, pt)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdProjectType)
}

func (ptc *ProjectTypeController) UpdateProjectType(c *gin.Context) {
	var updatedPT *models.ProjectType

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedPT); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptc.db, c, err.Error())
		return
	}

	updatedPT, err = ptc.service.Update(uint(id), currentUserID, updatedPT)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedPT)
}

func (ptc *ProjectTypeController) DeleteProjectTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptc.db, c, err.Error())
		return
	}

	failedIDs, err := ptc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ptc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ptc *ProjectTypeController) ActivateProjectTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptc.db, c, err.Error())
		return
	}

	err = ptc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ptc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (ptc *ProjectTypeController) DeactivateProjectTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptc.db, c, err.Error())
		return
	}

	err = ptc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ptc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
