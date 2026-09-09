package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectSecurityGateChecklistCommentController struct {
	service services.ProjectSecurityGateChecklistCommentServiceInterface
	db      *gorm.DB
}

func NewProjectSecurityGateChecklistCommentController(service services.ProjectSecurityGateChecklistCommentServiceInterface) *ProjectSecurityGateChecklistCommentController {
	return &ProjectSecurityGateChecklistCommentController{
		service: service,
	}
}

func (psgccc *ProjectSecurityGateChecklistCommentController) CreateProjectSecurityGateChecklistComment(c *gin.Context) {
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

	var form dtos.ProjectSecurityGateChecklistCommentFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgccc.db, c, err.Error())
		return
	}

	id, err := psgccc.service.Create(uint(securityGateID), uint(checklistID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psgccc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (psgccc *ProjectSecurityGateChecklistCommentController) UpdateProjectSecurityGateChecklistComment(c *gin.Context) {
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

	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.ProjectSecurityGateChecklistCommentFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgccc.db, c, err.Error())
		return
	}

	err = psgccc.service.Update(uint(securityGateID), uint(checklistID), uint(commentID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psgccc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (psgccc *ProjectSecurityGateChecklistCommentController) DeleteProjectSecurityGateChecklistComments(c *gin.Context) {
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

	var request struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgccc.db, c, err.Error())
		return
	}

	failedIDs, err := psgccc.service.Delete(uint(securityGateID), uint(checklistID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(psgccc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (psgccc *ProjectSecurityGateChecklistCommentController) GetProjectSecurityGateChecklistCommentDocumentPreview(c *gin.Context) {
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

	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := psgccc.service.GetDocumentPath(uint(securityGateID), uint(checklistID), uint(commentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(psgccc.db, c, err.Error())
		}
		return
	}

	c.File(documentPath)
}

func (psgccc *ProjectSecurityGateChecklistCommentController) DownloadProjectSecurityGateChecklistCommentDocument(c *gin.Context) {
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

	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := psgccc.service.GetDocumentPath(uint(securityGateID), uint(checklistID), uint(commentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(psgccc.db, c, err.Error())
		}
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
