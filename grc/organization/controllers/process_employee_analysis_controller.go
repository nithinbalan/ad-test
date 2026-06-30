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

type ProcessEmployeeAnalysisController struct {
	service services.ProcessEmployeeAnalysisServiceInterface
	db      *gorm.DB
}

func NewProcessEmployeeAnalysisController(service services.ProcessEmployeeAnalysisServiceInterface) *ProcessEmployeeAnalysisController {
	return &ProcessEmployeeAnalysisController{
		service: service,
	}
}

func (peac *ProcessEmployeeAnalysisController) GetAllProcessEmployeeAnalyses(ctx *gin.Context) {
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
		utility.SendInternalErrResp(peac.db, ctx, "Invalid export parameter")
		return
	}

	page, err := utility.ParseQueryParamToInt(ctx, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(peac.db, ctx, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(ctx, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(peac.db, ctx, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"created_at":           "process_employee_analysis.created_at",
		"designation_title_en": "designations.title_en",
		"designation_title_ar": "designations.title_ar",
		"normal_count":         "process_employee_analysis.normal_count",
		"critical_count":       "process_employee_analysis.critical_count",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	queryParams := orgDTOs.ProcessEmployeeAnalysisListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Export:   export,
	}

	result, err := peac.service.GetAll(uint(processID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(peac.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (peac *ProcessEmployeeAnalysisController) GetProcessEmployeeAnalysisDetails(ctx *gin.Context) {
	processIDStr := ctx.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	item, err := peac.service.GetFullDetailsByID(uint(processID), uint(id))
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
		default:
			utility.SendInternalErrResp(peac.db, ctx, err.Error())
		}
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (peac *ProcessEmployeeAnalysisController) CreateProcessEmployeeAnalysis(ctx *gin.Context) {
	processIDStr := ctx.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	var dto *orgDTOs.ProcessEmployeeAnalysisFormDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(peac.db, ctx, err.Error())
		return
	}

	id, err := peac.service.Create(uint(processID), currentUserID, dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
		default:
			utility.SendInternalErrResp(peac.db, ctx, e.Error())
		}
		return
	}

	ctx.JSON(http.StatusCreated, id)
}

func (peac *ProcessEmployeeAnalysisController) UpdateProcessEmployeeAnalysis(ctx *gin.Context) {
	processIDStr := ctx.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	var dto *orgDTOs.ProcessEmployeeAnalysisFormDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(peac.db, ctx, err.Error())
		return
	}

	err = peac.service.Update(uint(processID), uint(id), currentUserID, dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
		default:
			utility.SendInternalErrResp(peac.db, ctx, e.Error())
		}
		return
	}

	ctx.Status(http.StatusOK)
}

func (peac *ProcessEmployeeAnalysisController) DeleteProcessEmployeeAnalyses(ctx *gin.Context) {
	processIDStr := ctx.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	var dto struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(peac.db, ctx, err.Error())
		return
	}

	failedIds, err := peac.service.Delete(uint(processID), dto.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(peac.db, ctx, err.Error())
		return
	}

	if len(failedIds) > 0 {
		ctx.JSON(http.StatusPartialContent, gin.H{
			"failed_ids": failedIds,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"failed_ids": failedIds})
}
