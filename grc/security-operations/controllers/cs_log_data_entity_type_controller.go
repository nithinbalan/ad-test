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

type CSLogDataEntityTypeController struct {
	service services.CSLogDataEntityTypeServiceInterface
	db      *gorm.DB
}

func NewCSLogDataEntityTypeController(service services.CSLogDataEntityTypeServiceInterface) *CSLogDataEntityTypeController {
	return &CSLogDataEntityTypeController{
		service: service,
	}
}

func (c *CSLogDataEntityTypeController) GetAll(ctx *gin.Context) {
	sort := ctx.Query("sort_order")
	sortBy := ctx.Query("sort_by")
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")

	entityTypes, err := c.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, entityTypes)
}

func (c *CSLogDataEntityTypeController) GetDetails(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	entityType, err := c.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	ctx.JSON(http.StatusOK, entityType)
}

func (c *CSLogDataEntityTypeController) Create(ctx *gin.Context) {
	var entityType *models.CSLogDataEntityType

	if err := ctx.ShouldBindJSON(&entityType); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	created, err := c.service.Create(currentUserID, entityType)
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

func (c *CSLogDataEntityTypeController) Update(ctx *gin.Context) {
	var updated *models.CSLogDataEntityType

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

func (c *CSLogDataEntityTypeController) Delete(ctx *gin.Context) {
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

func (c *CSLogDataEntityTypeController) Activate(ctx *gin.Context) {
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

func (c *CSLogDataEntityTypeController) Deactivate(ctx *gin.Context) {
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
