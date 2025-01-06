package controllers

import (
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentChangeRequestWorkflowController struct {
	service services.DocumentChangeRequestWorkflowServiceInterface
	db      *gorm.DB
}

func NewDocumentChangeRequestWorkflowController(service services.DocumentChangeRequestWorkflowServiceInterface) *DocumentChangeRequestWorkflowController {
	return &DocumentChangeRequestWorkflowController{
		service: service,
	}
}

func (dcrwc *DocumentChangeRequestWorkflowController) GetWorkflowList(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := dcrwc.service.GetDocumentCRWorkflowList(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dcrwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}
