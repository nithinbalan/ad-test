package controllers

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RCMController struct {
	service services.RCMServiceInterface
	db      *gorm.DB
}

func NewRCMController(service services.RCMServiceInterface) *RCMController {
	return &RCMController{
		service: service,
	}
}

func (rc *RCMController) GetAllRiskControls(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, "Invalid page_size parameter")
		return
	}

	auditIdStr := c.Param("auditId")
	auditId, err := strconv.ParseUint(auditIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	queryParams := auditDTOs.RCMListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	riskControls, err := rc.service.GetAll(&queryParams, uint(auditId))
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, riskControls)
}
