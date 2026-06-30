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

type ProcessDependencyTierController struct {
	service services.ProcessDependencyTierServiceInterface
	db      *gorm.DB
}

func NewProcessDependencyTierController(service services.ProcessDependencyTierServiceInterface) *ProcessDependencyTierController {
	return &ProcessDependencyTierController{
		service: service,
	}
}

func (c *ProcessDependencyTierController) GetAllProcessDependencyTiers(ctx *gin.Context) {
	sort := ctx.Query("sort_order")
	sortBy := ctx.Query("sort_by")
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")

	queryParams := dtos.ProcessDependencyTierListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
	}

	tiers, err := c.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tiers)
}

func (c *ProcessDependencyTierController) GetProcessDependencyTierDetails(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	tier, err := c.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	ctx.JSON(http.StatusOK, tier)
}

func (c *ProcessDependencyTierController) CreateProcessDependencyTier(ctx *gin.Context) {
	var formDTO dtos.ProcessDependencyTierFormDTO

	if err := ctx.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	id, err := c.service.Create(currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
		default:
			utility.SendInternalErrResp(c.db, ctx, e.Error())
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

func (c *ProcessDependencyTierController) UpdateProcessDependencyTier(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	var formDTO dtos.ProcessDependencyTierFormDTO
	if err := ctx.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	err = c.service.Update(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
		default:
			utility.SendInternalErrResp(c.db, ctx, e.Error())
		}
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *ProcessDependencyTierController) DeleteProcessDependencyTiers(ctx *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	failedIds, err := c.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	if len(failedIds) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message":    "Some process dependency tiers could not be deleted",
			"failed_ids": failedIds,
		})
		return
	}

	ctx.Status(http.StatusOK)
}
