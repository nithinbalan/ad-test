package controllers

import (
	"errors"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProcessExternalDependencyEntityTypeController struct {
	service services.ProcessExternalDependencyEntityTypeServiceInterface
	db      *gorm.DB
}

func NewProcessExternalDependencyEntityTypeController(service services.ProcessExternalDependencyEntityTypeServiceInterface) *ProcessExternalDependencyEntityTypeController {
	return &ProcessExternalDependencyEntityTypeController{
		service: service,
	}
}

func (c *ProcessExternalDependencyEntityTypeController) GetAllProcessExternalDependencyEntityTypes(ctx *gin.Context) {
	sort := ctx.DefaultQuery("sort_order", "desc")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")

	validSortFields := map[string]string{
		"title_en":   "title_en",
		"title_ar":   "title_ar",
		"status":     "status",
		"created_at": "created_at",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	queryParams := dtos.ProcessExternalDependencyEntityTypeListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortField,
		Keyword: keyword,
		Status:  status,
	}

	entityTypes, err := c.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, entityTypes)
}

func (c *ProcessExternalDependencyEntityTypeController) GetProcessExternalDependencyEntityTypeDetails(ctx *gin.Context) {
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

func (c *ProcessExternalDependencyEntityTypeController) CreateProcessExternalDependencyEntityType(ctx *gin.Context) {
	var formDTO dtos.ProcessExternalDependencyEntityTypeFormDTO

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

func (c *ProcessExternalDependencyEntityTypeController) UpdateProcessExternalDependencyEntityType(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	var formDTO dtos.ProcessExternalDependencyEntityTypeFormDTO
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

func (c *ProcessExternalDependencyEntityTypeController) DeleteProcessExternalDependencyEntityTypes(ctx *gin.Context) {
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
			"message":    "Some process external dependency entity types could not be deleted",
			"failed_ids": failedIds,
		})
		return
	}

	ctx.Status(http.StatusOK)
}
