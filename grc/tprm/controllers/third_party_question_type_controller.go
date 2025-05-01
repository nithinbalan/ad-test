package controllers

import (
	"jethur/grc/tprm/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThirdPartyQuestionTypeController struct {
	service services.ThirdPartyQuestionTypeServiceInterface
	db      *gorm.DB
}

func NewThirdPartyQuestionTypeController(service services.ThirdPartyQuestionTypeServiceInterface) *ThirdPartyQuestionTypeController {
	return &ThirdPartyQuestionTypeController{
		service: service,
	}
}

func (tpqtc *ThirdPartyQuestionTypeController) GetAllThirdPartyQuestionTypes(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	types, err := tpqtc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(tpqtc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, types)
}

func (tpqtc *ThirdPartyQuestionTypeController) GetThirdPartyQuestionTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendInternalErrResp(tpqtc.db, c, "Invalid id")
		return
	}

	questionType, err := tpqtc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, questionType)
}
