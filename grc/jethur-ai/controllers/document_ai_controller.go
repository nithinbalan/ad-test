package controllers

import (
	"jethur/grc/jethur-ai/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentAIController struct {
	service services.DocumentAIServiceInterface
	db      *gorm.DB
}

func NewDocumentAIController(service services.DocumentAIServiceInterface) *DocumentAIController {
	return &DocumentAIController{
		service: service,
	}
}

func (ac *DocumentAIController) GetDocumentSummary(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	summary, err := ac.service.GetAISummary(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, summary)
}
