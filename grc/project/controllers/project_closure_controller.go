package controllers

import (
	"net/http"
	"strconv"

	"jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectClosureController struct {
	service services.ProjectClosureServiceInterface
	db      *gorm.DB
}

func NewProjectClosureController(service services.ProjectClosureServiceInterface, db *gorm.DB) *ProjectClosureController {
	return &ProjectClosureController{
		service: service,
		db:      db,
	}
}

func (pcc *ProjectClosureController) GetAllProjectClosures(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	closures, err := pcc.service.GetAll(uint(projectID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(pcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, closures)
}

func (pcc *ProjectClosureController) GetProjectClosureDetails(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	closureIDStr := c.Param("closureId")
	closureID, err := strconv.ParseUint(closureIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	closure, err := pcc.service.GetFullDetailsByID(uint(projectID), uint(closureID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(pcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, closure)
}

func (pcc *ProjectClosureController) CreateProjectClosure(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.ProjectClosureFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcc.db, c, err.Error())
		return
	}

	id, err := pcc.service.Create(uint(projectID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (pcc *ProjectClosureController) UpdateProjectClosure(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	closureIDStr := c.Param("closureId")
	closureID, err := strconv.ParseUint(closureIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.ProjectClosureFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcc.db, c, err.Error())
		return
	}

	err = pcc.service.Update(uint(projectID), uint(closureID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusNoContent)
}

func (pcc *ProjectClosureController) DeleteProjectClosure(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	closureIDStr := c.Param("closureId")
	closureID, err := strconv.ParseUint(closureIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcc.db, c, err.Error())
		return
	}

	err = pcc.service.Delete(uint(projectID), uint(closureID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(pcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusNoContent)
}

func (pcc *ProjectClosureController) ApproveProjectClosure(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	closureIDStr := c.Param("closureId")
	closureID, err := strconv.ParseUint(closureIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcc.db, c, err.Error())
		return
	}

	err = pcc.service.Approve(uint(projectID), uint(closureID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			c.JSON(http.StatusForbidden, gin.H{"error": e.Error()})
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project closure approved successfully and project phase updated to go live"})
}
