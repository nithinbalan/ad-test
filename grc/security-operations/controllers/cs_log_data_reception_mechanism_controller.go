package controllers

import (
	"jethur/grc/models"
	"jethur/grc/security-operations/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CSLogDataReceptionMechanismController struct {
	service services.CSLogDataReceptionMechanismServiceInterface
	db      *gorm.DB
}

func NewCSLogDataReceptionMechanismController(service services.CSLogDataReceptionMechanismServiceInterface) *CSLogDataReceptionMechanismController {
	return &CSLogDataReceptionMechanismController{
		service: service,
	}
}

func (c *CSLogDataReceptionMechanismController) GetAll(ctx *gin.Context) {
	sort := ctx.Query("sort_order")
	sortBy := ctx.Query("sort_by")
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")

	mechanisms, err := c.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, mechanisms)
}

func (c *CSLogDataReceptionMechanismController) GetDetails(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	mechanism, err := c.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	ctx.JSON(http.StatusOK, mechanism)
}

func (c *CSLogDataReceptionMechanismController) Create(ctx *gin.Context) {
	var mechanism *models.CSLogDataReceptionMechanism

	if err := ctx.ShouldBindJSON(&mechanism); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	created, err := c.service.Create(currentUserID, mechanism)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
			return
		default:
			utility.SendInternalErrResp(c.db, ctx, e.Error())
			return
		}
	}

	ctx.JSON(http.StatusCreated, created)
}

func (c *CSLogDataReceptionMechanismController) Update(ctx *gin.Context) {
	var updated *models.CSLogDataReceptionMechanism

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	if err := ctx.ShouldBindJSON(&updated); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	updated, err = c.service.Update(uint(id), currentUserID, updated)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
			return
		default:
			utility.SendInternalErrResp(c.db, ctx, e.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, updated)
}

func (c *CSLogDataReceptionMechanismController) Delete(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	failedIDs, err := c.service.Delete(req.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (c *CSLogDataReceptionMechanismController) Activate(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	if err := c.service.Activate(req.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *CSLogDataReceptionMechanismController) Deactivate(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	if err := c.service.Deactivate(req.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}
