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

type ProcessExternalDependencyController struct {
	service services.ProcessExternalDependencyServiceInterface
	db      *gorm.DB
}

func NewProcessExternalDependencyController(service services.ProcessExternalDependencyServiceInterface) *ProcessExternalDependencyController {
	return &ProcessExternalDependencyController{
		service: service,
	}
}

func (pedc *ProcessExternalDependencyController) GetAllProcessExternalDependencies(ctx *gin.Context) {
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
		utility.SendInternalErrResp(pedc.db, ctx, "Invalid export parameter")
		return
	}

	page, err := utility.ParseQueryParamToInt(ctx, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(pedc.db, ctx, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(ctx, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(pedc.db, ctx, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"created_at":                 "process_external_dependencies.created_at",
		"entity_name":                "process_external_dependencies.entity_name",
		"entity_type_title_en":       "entity_type.title_en",
		"entity_type_title_ar":       "entity_type.title_ar",
		"dependency_tier_title_en":   "process_dependency_tiers.title_en",
		"dependency_tier_title_ar":   "process_dependency_tiers.title_ar",
		"needed_support_description": "process_external_dependencies.needed_support_description",
		"vendor_rto":                 "process_external_dependencies.vendor_rto",
		"contact_number":             "process_external_dependencies.contact_number",
		"contact_email":              "process_external_dependencies.contact_email",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	queryParams := orgDTOs.ProcessExternalDependencyListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Export:   export,
	}

	result, err := pedc.service.GetAll(uint(processID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(pedc.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (pedc *ProcessExternalDependencyController) GetProcessExternalDependencyDetails(ctx *gin.Context) {
	processIDStr := ctx.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	idStr := ctx.Param("externalDependencyId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	item, err := pedc.service.GetFullDetailsByID(uint(processID), uint(id))
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
		default:
			utility.SendInternalErrResp(pedc.db, ctx, err.Error())
		}
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (pedc *ProcessExternalDependencyController) CreateProcessExternalDependency(ctx *gin.Context) {
	processIDStr := ctx.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	var dto orgDTOs.ProcessExternalDependencyFormDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(pedc.db, ctx, err.Error())
		return
	}

	id, err := pedc.service.Create(uint(processID), currentUserID, &dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
		default:
			utility.SendInternalErrResp(pedc.db, ctx, e.Error())
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

func (pedc *ProcessExternalDependencyController) UpdateProcessExternalDependency(ctx *gin.Context) {
	processIDStr := ctx.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	idStr := ctx.Param("externalDependencyId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	var dto orgDTOs.ProcessExternalDependencyFormDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(pedc.db, ctx, err.Error())
		return
	}

	err = pedc.service.Update(uint(processID), uint(id), currentUserID, &dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
		default:
			utility.SendInternalErrResp(pedc.db, ctx, e.Error())
		}
		return
	}

	ctx.Status(http.StatusOK)
}

func (pedc *ProcessExternalDependencyController) DeleteProcessExternalDependencies(ctx *gin.Context) {
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
		utility.SendInternalErrResp(pedc.db, ctx, err.Error())
		return
	}

	failedIDs, err := pedc.service.Delete(uint(processID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pedc.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
