package controllers

import (
	"net/http"

	"jethur/grc/penetration-testing/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenetrationTestRetestStatusController struct {
	service services.PenetrationTestRetestStatusServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestRetestStatusController(service services.PenetrationTestRetestStatusServiceInterface) *PenetrationTestRetestStatusController {
	return &PenetrationTestRetestStatusController{service: service}
}

func (ptrsc *PenetrationTestRetestStatusController) GetAllPenetrationTestRetestStatuses(c *gin.Context) {
	values, err := ptrsc.service.GetAll(c.Query("sort_order"), c.Query("sort_by"), c.Query("keyword"), c.Query("status"))
	if err != nil {
		utility.SendInternalErrResp(ptrsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, values)
}
