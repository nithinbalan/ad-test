package controllers

import (
	"net/http"

	"jethur/grc/penetration-testing/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenetrationTestRetestResultController struct {
	service services.PenetrationTestRetestResultServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestRetestResultController(service services.PenetrationTestRetestResultServiceInterface) *PenetrationTestRetestResultController {
	return &PenetrationTestRetestResultController{service: service}
}

func (ptrrc *PenetrationTestRetestResultController) GetAllPenetrationTestRetestResults(c *gin.Context) {
	values, err := ptrrc.service.GetAll(c.Query("sort_order"), c.Query("sort_by"), c.Query("keyword"), c.Query("status"))
	if err != nil {
		utility.SendInternalErrResp(ptrrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, values)
}
