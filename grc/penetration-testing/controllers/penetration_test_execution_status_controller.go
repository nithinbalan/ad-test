package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"jethur/grc/penetration-testing/services"
	"jethur/utility"
	"net/http"
)

type PenetrationTestExecutionStatusController struct {
	service services.PenetrationTestExecutionStatusServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestExecutionStatusController(service services.PenetrationTestExecutionStatusServiceInterface) *PenetrationTestExecutionStatusController {
	return &PenetrationTestExecutionStatusController{service: service}
}
func (ptesc *PenetrationTestExecutionStatusController) GetAllPenetrationTestExecutionStatuses(c *gin.Context) {
	values, err := ptesc.service.GetAll(c.Query("sort_order"), c.Query("sort_by"), c.Query("keyword"), c.Query("status"))
	if err != nil {
		utility.SendInternalErrResp(ptesc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, values)
}
