package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"

	"jethur/grc/project/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectSecurityGateChecklistExceptionWorkflowController struct {
	service services.ProjectSecurityGateChecklistExceptionWorkflowServiceInterface
	db      *gorm.DB
}

func NewProjectSecurityGateChecklistExceptionWorkflowController(service services.ProjectSecurityGateChecklistExceptionWorkflowServiceInterface) *ProjectSecurityGateChecklistExceptionWorkflowController {
	return &ProjectSecurityGateChecklistExceptionWorkflowController{
		service: service,
	}
}

func (psgcewc *ProjectSecurityGateChecklistExceptionWorkflowController) GetProjectSecurityGateChecklistExceptionWorkflowList(c *gin.Context) {
	securityGateID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	checklistID, err := strconv.ParseUint(c.Param("checklistId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := psgcewc.service.GetWorkflowList(uint(securityGateID), uint(checklistID))
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
			utility.SendInternalErrResp(psgcewc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}

func (psgcewc *ProjectSecurityGateChecklistExceptionWorkflowController) documentPath(c *gin.Context) (string, bool) {
	securityGateID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return "", false
	}

	checklistID, err := strconv.ParseUint(c.Param("checklistId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return "", false
	}

	historyID, err := strconv.ParseUint(c.Param("historyId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return "", false
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return "", false
	}

	documentPath, err := psgcewc.service.GetDocumentPath(uint(securityGateID), uint(checklistID), uint(historyID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(psgcewc.db, c, err.Error())
		}
		return "", false
	}

	return documentPath, true
}

func (psgcewc *ProjectSecurityGateChecklistExceptionWorkflowController) GetProjectSecurityGateChecklistExceptionDocumentPreview(c *gin.Context) {
	documentPath, ok := psgcewc.documentPath(c)
	if !ok {
		return
	}

	c.File(documentPath)
}

func (psgcewc *ProjectSecurityGateChecklistExceptionWorkflowController) DownloadProjectSecurityGateChecklistExceptionDocument(c *gin.Context) {
	documentPath, ok := psgcewc.documentPath(c)
	if !ok {
		return
	}

	filename := filepath.Base(documentPath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(documentPath)
}
