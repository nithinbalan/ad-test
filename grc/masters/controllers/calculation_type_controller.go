package controllers

import (
	"fmt"
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CalculationTypeController struct {
	service services.CalculationTypeServiceInterface
	db      *gorm.DB
}

func NewCalculationTypeController(service services.CalculationTypeServiceInterface) *CalculationTypeController {
	return &CalculationTypeController{
		service: service,
	}
}

func (ctc *CalculationTypeController) GetAllCalculationTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	fmt.Println("1")
	calculationTypes, err := ctc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(ctc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, calculationTypes)
}

func (ctc *CalculationTypeController) GetCalculationTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	calculationType, err := ctc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, calculationType)
}
