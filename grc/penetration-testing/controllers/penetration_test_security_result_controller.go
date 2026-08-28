package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"jethur/grc/penetration-testing/services"
	"jethur/utility"
	"net/http"
)

type PenetrationTestSecurityResultController struct {
	service services.PenetrationTestSecurityResultServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestSecurityResultController(service services.PenetrationTestSecurityResultServiceInterface) *PenetrationTestSecurityResultController {
	return &PenetrationTestSecurityResultController{service: service}
}
func (ptsrc *PenetrationTestSecurityResultController) GetAllPenetrationTestSecurityResults(c *gin.Context) {
	values, err := ptsrc.service.GetAll(c.Query("sort_order"), c.Query("sort_by"), c.Query("keyword"), c.Query("status"))
	if err != nil {
		utility.SendInternalErrResp(ptsrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, values)
}
