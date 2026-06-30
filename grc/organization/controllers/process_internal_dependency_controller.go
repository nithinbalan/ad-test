package controllers

import (
	"errors"
	orgDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProcessInternalDependencyController struct {
	service services.ProcessInternalDependencyServiceInterface
	db      *gorm.DB
}

func NewProcessInternalDependencyController(service services.ProcessInternalDependencyServiceInterface) *ProcessInternalDependencyController {
	return &ProcessInternalDependencyController{
		service: service,
	}
}

func (pidc *ProcessInternalDependencyController) GetAllProcessInternalDependencies(ctx *gin.Context) {
	processIDStr := ctx.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	sort := ctx.DefaultQuery("sort_order", "desc")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	keyword := ctx.Query("keyword")

	export, err := strconv.ParseBool(ctx.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(pidc.db, ctx, "Invalid export parameter")
		return
	}

	page, err := utility.ParseQueryParamToInt(ctx, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(pidc.db, ctx, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(ctx, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(pidc.db, ctx, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"created_at":                 "process_internal_dependencies.created_at",
		"supporting_unit_title_en":   "supporting_unit.title_en",
		"supporting_unit_title_ar":   "supporting_unit.title_ar",
		"dependency_tier_title_en":   "process_dependency_tiers.title_en",
		"dependency_tier_title_ar":   "process_dependency_tiers.title_ar",
		"dependency_percentage_from": "process_internal_dependencies.dependency_percentage_from",
		"dependency_percentage_to":   "process_internal_dependencies.dependency_percentage_to",
		"required_support":           "process_internal_dependencies.required_support",
		"notes":                      "process_internal_dependencies.notes",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	queryParams := orgDTOs.ProcessInternalDependencyListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Export:   export,
	}

	result, err := pidc.service.GetAll(uint(processID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(pidc.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (pidc *ProcessInternalDependencyController) GetProcessInternalDependencyDetails(ctx *gin.Context) {
	processIDStr := ctx.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	idStr := ctx.Param("internalDependencyId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	item, err := pidc.service.GetFullDetailsByID(uint(processID), uint(id))
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
		default:
			utility.SendInternalErrResp(pidc.db, ctx, err.Error())
		}
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (pidc *ProcessInternalDependencyController) CreateProcessInternalDependency(ctx *gin.Context) {
	processIDStr := ctx.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	var dto orgDTOs.ProcessInternalDependencyFormDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(pidc.db, ctx, err.Error())
		return
	}

	id, err := pidc.service.Create(uint(processID), currentUserID, &dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
		default:
			utility.SendInternalErrResp(pidc.db, ctx, e.Error())
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

func (pidc *ProcessInternalDependencyController) UpdateProcessInternalDependency(ctx *gin.Context) {
	processIDStr := ctx.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	idStr := ctx.Param("internalDependencyId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	var dto orgDTOs.ProcessInternalDependencyFormDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(pidc.db, ctx, err.Error())
		return
	}

	err = pidc.service.Update(uint(processID), uint(id), currentUserID, &dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
		default:
			utility.SendInternalErrResp(pidc.db, ctx, e.Error())
		}
		return
	}

	ctx.Status(http.StatusOK)
}

func (pidc *ProcessInternalDependencyController) DeleteProcessInternalDependencies(ctx *gin.Context) {
	processIDStr := ctx.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(pidc.db, ctx, err.Error())
		return
	}

	failedIDs, err := pidc.service.Delete(uint(processID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pidc.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
