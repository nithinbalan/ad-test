package controllers

import (
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserJobDescriptionController struct {
	service services.UserJobDescriptionServiceInterface
	db      *gorm.DB
}

func NewUserJobDescriptionController(service services.UserJobDescriptionServiceInterface) *UserJobDescriptionController {
	return &UserJobDescriptionController{
		service: service,
	}
}

func (ujdc *UserJobDescriptionController) GetAllUserJobDescriptions(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")

	dtos, err := ujdc.service.GetAll(sort, sortBy, keyword, uint(userID))
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ujdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, dtos)
}

func (ujdc *UserJobDescriptionController) GetUserJobDescriptionDetailsByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	detail, err := ujdc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(ujdc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusOK, detail)
}

func (ujdc *UserJobDescriptionController) AssignJobDescriptionsToUser(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var input dtos.AssignJobDescriptionsDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ujdc.db, c, err.Error())
		return
	}

	_, err = ujdc.service.Create(uint(userID), &input, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ujdc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (ctrl *UserJobDescriptionController) DeleteJobDescriptionUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	jdID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ctrl.db, c, err.Error())
		return
	}

	if err := ctrl.service.Delete(uint(userID), uint(jdID), currentUserID); err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(ctrl.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (ctrl *UserJobDescriptionController) EvaluateCompetency(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	_, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	competencyID, err := strconv.ParseUint(c.Param("cId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var dto dtos.EvaluateUserJDCompetencyDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ctrl.db, c, err.Error())
		return
	}

	if err := ctrl.service.Update(uint(userID), uint(competencyID), dto, currentUserID); err != nil {
		utility.SendInternalErrResp(ctrl.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (ujdc *UserJobDescriptionController) DownloadUserJobDescriptionDocument(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	jdcIDStr := c.Param("id")
	jdcID, err := strconv.ParseUint(jdcIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	competencyID, err := strconv.ParseUint(c.Param("cId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	docId, err := strconv.ParseUint(c.Param("docId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := ujdc.service.GetUserJobDescriptionDocumentPath(uint(userID), uint(jdcID), uint(competencyID), uint(docId), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(ujdc.db, c, err.Error())
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

func (ujdc *UserJobDescriptionController) GetUserJobDescriptionDocumentPreview(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	jdcIDStr := c.Param("id")
	jdcID, err := strconv.ParseUint(jdcIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	competencyID, err := strconv.ParseUint(c.Param("cId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	docId, err := strconv.ParseUint(c.Param("docId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := ujdc.service.GetUserJobDescriptionDocumentPath(uint(userID), uint(jdcID), uint(competencyID), uint(docId), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(ujdc.db, c, err.Error())
		}
		return
	}

	c.File(documentPath)
}
