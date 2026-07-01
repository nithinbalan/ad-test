package controllers

import (
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProcessEquipmentToolController struct {
	service services.ProcessEquipmentToolServiceInterface
	db      *gorm.DB
}

func NewProcessEquipmentToolController(service services.ProcessEquipmentToolServiceInterface) *ProcessEquipmentToolController {
	return &ProcessEquipmentToolController{service: service}
}

func (petc *ProcessEquipmentToolController) GetAllProcessEquipmentTools(ctx *gin.Context) {
	sort := ctx.Query("sort_order")
	sortBy := ctx.Query("sort_by")
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")

	queryParams := dtos.ProcessEquipmentToolListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
	}

	tools, err := petc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(petc.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tools)
}

func (petc *ProcessEquipmentToolController) GetProcessEquipmentToolDetails(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	tool, err := petc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	ctx.JSON(http.StatusOK, tool)
}

func (petc *ProcessEquipmentToolController) CreateProcessEquipmentTool(ctx *gin.Context) {
	var formDTO dtos.ProcessEquipmentToolFormDTO

	if err := ctx.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(petc.db, ctx, err.Error())
		return
	}

	id, err := petc.service.Create(currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
		default:
			utility.SendInternalErrResp(petc.db, ctx, e.Error())
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

func (petc *ProcessEquipmentToolController) UpdateProcessEquipmentTool(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	var formDTO dtos.ProcessEquipmentToolFormDTO
	if err := ctx.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(petc.db, ctx, err.Error())
		return
	}

	err = petc.service.Update(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
		default:
			utility.SendInternalErrResp(petc.db, ctx, e.Error())
		}
		return
	}

	ctx.Status(http.StatusOK)
}

func (petc *ProcessEquipmentToolController) DeleteProcessEquipmentTools(ctx *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(petc.db, ctx, err.Error())
		return
	}

	failedIds, err := petc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(petc.db, ctx, err.Error())
		return
	}

	if len(failedIds) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message":    "Some process equipment tools could not be deleted",
			"failed_ids": failedIds,
		})
		return
	}

	ctx.Status(http.StatusOK)
}
