package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EmailStatusController struct {
	service services.EmailStatusServiceInterface
	db      *gorm.DB
}

func NewEmailStatusController(service services.EmailStatusServiceInterface) *EmailStatusController {
	return &EmailStatusController{
		service: service,
	}
}

func (c *EmailStatusController) GetAllEmailStatuses(ctx *gin.Context) {
	sort := ctx.Query("sort_order")
	sortBy := ctx.Query("sort_by")
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")

	emailStatuses, err := c.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, emailStatuses)
}

func (c *EmailStatusController) GetEmailStatusDetails(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, "Invalid id")
		return
	}

	status, err := c.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	ctx.JSON(http.StatusOK, status)
}
