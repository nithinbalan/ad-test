package controllers

import (
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QualysController struct {
	service services.QualysServiceInterface
}

func NewQualysController(service services.QualysServiceInterface) *QualysController {
	return &QualysController{service: service}
}

func (qkc *QualysController) GetKBDetailsByQID(c *gin.Context) {
	qid := c.Param("qid")

	data, err := qkc.service.GetKBDetailsByQID(qid)
	if err != nil {
		switch e := err.(type) {
		case *utility.QualysConnectionError:
			utility.SendCustomErrResp(c, http.StatusUnprocessableEntity, err.Error())
			return
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(nil, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, data)
}
