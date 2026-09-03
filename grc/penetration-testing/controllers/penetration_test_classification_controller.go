package controllers

import (
	"jethur/grc/penetration-testing/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenetrationTestClassificationController struct {
	service services.PenetrationTestClassificationServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestClassificationController(service services.PenetrationTestClassificationServiceInterface) *PenetrationTestClassificationController {
	return &PenetrationTestClassificationController{service: service}
}

func (ptcc *PenetrationTestClassificationController) GetAllPenetrationTestClassifications(c *gin.Context) {
	classifications, err := ptcc.service.GetAll(c.Query("sort_order"), c.Query("sort_by"), c.Query("keyword"), c.Query("status"))
	if err != nil {
		utility.SendInternalErrResp(ptcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, classifications)
}
