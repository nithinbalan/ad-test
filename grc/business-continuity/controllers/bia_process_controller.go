package controllers

import (
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BiaProcessController struct {
	service services.BiaProcessServiceInterface
	db      *gorm.DB
}

func NewBiaProcessController(service services.BiaProcessServiceInterface) *BiaProcessController {
	return &BiaProcessController{
		service: service,
	}
}

func (bpc *BiaProcessController) GetBiaProcessDetails(c *gin.Context) {
	processIDStr := c.Param("processId")
	processID, err := strconv.ParseUint(processIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	biaProcess, err := bpc.service.GetFullDetailsByProcessID(uint(processID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(bpc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, biaProcess)
}
