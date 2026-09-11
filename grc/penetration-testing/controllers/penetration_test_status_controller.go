package controllers

import (
	"net/http"

	"jethur/grc/penetration-testing/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenetrationTestStatusController struct {
	service services.PenetrationTestStatusServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestStatusController(service services.PenetrationTestStatusServiceInterface) *PenetrationTestStatusController {
	return &PenetrationTestStatusController{service: service}
}

func (ptsc *PenetrationTestStatusController) GetAllPenetrationTestStatuses(c *gin.Context) {
	values, err := ptsc.service.GetAll(c.Query("sort_order"), c.Query("sort_by"), c.Query("keyword"), c.Query("status"))
	if err != nil {
		utility.SendInternalErrResp(ptsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, values)
}
