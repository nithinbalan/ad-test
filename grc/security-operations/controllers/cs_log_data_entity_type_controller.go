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

func (cldetc *CSLogDataEntityTypeController) GetAllCSLogDataEntityTypes(ctx *gin.Context) {
	sort := ctx.Query("sort_order")
	sortBy := ctx.Query("sort_by")
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")

	entityTypes, err := cldetc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(cldetc.db, ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, entityTypes)
}

func (cldetc *CSLogDataEntityTypeController) GetCSLogDataEntityTypeDetails(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	entityType, err := cldetc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	ctx.JSON(http.StatusOK, entityType)
}

func (cldetc *CSLogDataEntityTypeController) CreateCSLogDataEntityType(ctx *gin.Context) {
	var entityType *models.CSLogDataEntityType

	if err := ctx.ShouldBindJSON(&entityType); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(cldetc.db, ctx, err.Error())
		return
	}

	created, err := cldetc.service.Create(currentUserID, entityType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
			return
		default:
			utility.SendInternalErrResp(cldetc.db, ctx, e.Error())
			return
		}
	}

	ctx.JSON(http.StatusCreated, created)
}

func (cldetc *CSLogDataEntityTypeController) UpdateCSLogDataEntityType(ctx *gin.Context) {
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
		utility.SendInternalErrResp(cldetc.db, ctx, err.Error())
		return
	}

	updated, err = cldetc.service.Update(uint(id), currentUserID, updated)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
			return
		default:
			utility.SendInternalErrResp(cldetc.db, ctx, e.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, updated)
}

func (cldetc *CSLogDataEntityTypeController) DeleteCSLogDataEntityTypes(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(cldetc.db, ctx, err.Error())
		return
	}

	failedIDs, err := cldetc.service.Delete(req.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cldetc.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (cldetc *CSLogDataEntityTypeController) ActivateCSLogDataEntityTypes(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(cldetc.db, ctx, err.Error())
		return
	}

	if err := cldetc.service.Activate(req.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(cldetc.db, ctx, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (cldetc *CSLogDataEntityTypeController) DeactivateCSLogDataEntityTypes(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(cldetc.db, ctx, err.Error())
		return
	}

	if err := cldetc.service.Deactivate(req.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(cldetc.db, ctx, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}
