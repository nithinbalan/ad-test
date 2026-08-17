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

type CSLogCommunicationChannelController struct {
	service services.CSLogCommunicationChannelServiceInterface
	db      *gorm.DB
}

func NewCSLogCommunicationChannelController(service services.CSLogCommunicationChannelServiceInterface) *CSLogCommunicationChannelController {
	return &CSLogCommunicationChannelController{
		service: service,
	}
}

func (clccc *CSLogCommunicationChannelController) GetAllCSLogCommunicationChannels(ctx *gin.Context) {
	sort := ctx.Query("sort_order")
	sortBy := ctx.Query("sort_by")
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")

	channels, err := clccc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(clccc.db, ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, channels)
}

func (clccc *CSLogCommunicationChannelController) GetCSLogCommunicationChannelDetails(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	channel, err := clccc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	ctx.JSON(http.StatusOK, channel)
}

func (clccc *CSLogCommunicationChannelController) CreateCSLogCommunicationChannel(ctx *gin.Context) {
	var channel *models.CSLogCommunicationChannel

	if err := ctx.ShouldBindJSON(&channel); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(clccc.db, ctx, err.Error())
		return
	}

	created, err := clccc.service.Create(currentUserID, channel)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
			return
		default:
			utility.SendInternalErrResp(clccc.db, ctx, e.Error())
			return
		}
	}

	ctx.JSON(http.StatusCreated, created)
}

func (clccc *CSLogCommunicationChannelController) UpdateCSLogCommunicationChannel(ctx *gin.Context) {
	var updated *models.CSLogCommunicationChannel

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
		utility.SendInternalErrResp(clccc.db, ctx, err.Error())
		return
	}

	updated, err = clccc.service.Update(uint(id), currentUserID, updated)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
			return
		default:
			utility.SendInternalErrResp(clccc.db, ctx, e.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, updated)
}

func (clccc *CSLogCommunicationChannelController) DeleteCSLogCommunicationChannels(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(clccc.db, ctx, err.Error())
		return
	}

	failedIDs, err := clccc.service.Delete(req.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(clccc.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (clccc *CSLogCommunicationChannelController) ActivateCSLogCommunicationChannels(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(clccc.db, ctx, err.Error())
		return
	}

	if err := clccc.service.Activate(req.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(clccc.db, ctx, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (clccc *CSLogCommunicationChannelController) DeactivateCSLogCommunicationChannels(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(clccc.db, ctx, err.Error())
		return
	}

	if err := clccc.service.Deactivate(req.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(clccc.db, ctx, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}
