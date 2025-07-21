package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectActionPlanStatusUpdateController struct {
	service services.ProjectActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewProjectActionPlanStatusUpdateController(service services.ProjectActionPlanStatusUpdateServiceInterface) *ProjectActionPlanStatusUpdateController {
	return &ProjectActionPlanStatusUpdateController{
		service: service,
	}
}

func (pasuc *ProjectActionPlanStatusUpdateController) GetAllProjectActionPlanStatusUpdates(c *gin.Context) {
	projectIDStr := c.Param("projectId")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	projectActionPlanStatuses, err := pasuc.service.GetAll(uint(projectID), uint(actionPlanID), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(pasuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, projectActionPlanStatuses)
}

func (pasuc *ProjectActionPlanStatusUpdateController) CreateProjectActionPlanStatusUpdate(c *gin.Context) {
	var projectActionPlanStatusUpdateFormDTO projectDTOs.ProjectActionPlanStatusUpdateFormDTO

	projectIDStr := c.Param("projectId")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&projectActionPlanStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pasuc.db, c, err.Error())
		return
	}

	createdID, err := pasuc.service.Create(uint(projectID), uint(actionPlanID), currentUserID, &projectActionPlanStatusUpdateFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pasuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (pasuc *ProjectActionPlanStatusUpdateController) UpdateProjectActionPlanStatusUpdate(c *gin.Context) {
	var projectActionPlanStatusUpdateFormDTO projectDTOs.ProjectActionPlanStatusUpdateFormDTO

	projectIDStr := c.Param("projectId")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&projectActionPlanStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pasuc.db, c, err.Error())
		return
	}

	err = pasuc.service.Update(uint(projectID), uint(actionPlanID), uint(id), currentUserID, &projectActionPlanStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pasuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (pasuc *ProjectActionPlanStatusUpdateController) GetLatestProjectActionPlanStatusUpdate(c *gin.Context) {
	projectIDStr := c.Param("projectId")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	projectActionPlanStatusUpdate, err := pasuc.service.GetLatestProjectActionPlanStatusUpdate(uint(projectID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(pasuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, projectActionPlanStatusUpdate)
}

func (pasuc *ProjectActionPlanStatusUpdateController) GetProjectActionPlanStatusUpdateDocumentPreview(c *gin.Context) {
	projectIDStr := c.Param("projectId")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	projectActionPlanStatusUpdateDocumentPath, err := pasuc.service.GetDocumentPath(uint(projectID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(pasuc.db, c, err.Error())
		}
		return
	}
	c.File(projectActionPlanStatusUpdateDocumentPath)
}

func (pasuc *ProjectActionPlanStatusUpdateController) DownloadProjectActionPlanStatusUpdateDocument(c *gin.Context) {
	projectIDStr := c.Param("projectId")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	projectActionPlanStatusUpdateDocumentPath, err := pasuc.service.GetDocumentPath(uint(projectID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(pasuc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(projectActionPlanStatusUpdateDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(projectActionPlanStatusUpdateDocumentPath)
}

func (pasuc *ProjectActionPlanStatusUpdateController) ApproveProjectActionPlan(c *gin.Context) {
	var projectActionPlanStatusReviewUserUpdateFormDTO projectDTOs.ProjectActionPlanReviewUserStatusUpdateFormDTO

	projectID, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	projectActionPlanID, err := strconv.Atoi(c.Param("actionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&projectActionPlanStatusReviewUserUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pasuc.db, c, err.Error())
		return
	}

	err = pasuc.service.Approve(uint(projectID), uint(projectActionPlanID), currentUserID, &projectActionPlanStatusReviewUserUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(pasuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (pasuc *ProjectActionPlanStatusUpdateController) RevertProjectActionPlan(c *gin.Context) {
	var projectActionPlanStatusReviewUserUpdateFormDTO projectDTOs.ProjectActionPlanReviewUserStatusUpdateFormDTO

	projectID, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	projectActionPlanID, err := strconv.Atoi(c.Param("actionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&projectActionPlanStatusReviewUserUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pasuc.db, c, err.Error())
		return
	}

	err = pasuc.service.Revert(uint(projectID), uint(projectActionPlanID), currentUserID, &projectActionPlanStatusReviewUserUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(pasuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
