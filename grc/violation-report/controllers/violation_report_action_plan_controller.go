package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	"jethur/grc/violation-report/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ViolationReportActionPlanController struct {
	service services.ViolationReportActionPlanServiceInterface
	db      *gorm.DB
}

func NewViolationReportActionPlanController(service services.ViolationReportActionPlanServiceInterface) *ViolationReportActionPlanController {
	return &ViolationReportActionPlanController{
		service: service,
	}
}

func (vapc *ViolationReportActionPlanController) GetAllViolationReportActionPlans(c *gin.Context) {
	violationReportIDStr := c.Param("id")

	var violationReportID uint64
	var err error

	if violationReportIDStr != "" {
		violationReportID, err = strconv.ParseUint(violationReportIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		violationReportID = 0
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	exportFromContext, exists := c.Get("export")
	var export bool

	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(vapc.db, c, "Invalid export parameter")
			return
		}
	}

	allowedParams := []string{
		"severity_level_ids",
		"action_plan_status_ids",
		"responsible_unit_ids",
		"responsible_user_ids",
		"violation_report_ids",
	}

	validSortFields := map[string]string{
		"title":               "violation_report_action_plans.title",
		"description":         "violation_report_action_plans.description",
		"violation_title":     "violation_reports.title",
		"severity_level_en":   "severity_levels.title_en",
		"severity_level_ar":   "severity_levels.title_ar",
		"responsible_user_en": "responsible_user.first_name_en",
		"responsible_user_ar": "responsible_user.first_name_ar",
		"target_date":         "violation_report_action_plans.target_date",
		"percentage":          "violation_report_action_plans.percentage",
		"status_en":           "action_plan_statuses.title_en",
		"status_ar":           "action_plan_statuses.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(vapc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(vapc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := violationReportDTOs.ViolationReportActionPlanListQueryParamsDTO{
		Sort:           sort,
		SortBy:         dbField,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		Filters:        filters,
		QuickFilter:    quickFilter,
		TargetDateFrom: targetDateFrom,
		TargetDateTo:   targetDateTo,
		Export:         export,
		DateFrom:       dateFrom,
		DateTo:         dateTo,
	}

	violationActionPlans, err := vapc.service.GetAll(uint(violationReportID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(vapc.db, c, e.Error())
			return
		}
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(vapc.db, c, err.Error())
			return
		}

		excelFile, err := vapc.service.ExportViolationReportActionPlan(violationActionPlans.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(vapc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Violation report action plan.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())
		return
	}

	c.JSON(http.StatusOK, violationActionPlans)
}

func (vapc *ViolationReportActionPlanController) GetViolationReportActionPlanDetails(c *gin.Context) {
	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	violationActionPlan, err := vapc.service.GetFullDetailsByID(uint(violationReportID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(vapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, violationActionPlan)
}

// func (vapc *ViolationReportActionPlanController) GetUsedActionPlanMasterIDs(c *gin.Context) {
// 	violationReportIDStr := c.Param("id")

// 	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
// 	if err != nil {
// 		utility.SendNotFoundErrResp(c)
// 		return
// 	}

// 	ids, err := vapc.service.GetUsedActionPlanMasterIDs(uint(violationReportID))
// 	if err != nil {
// 		switch e := err.(type) {
// 		case *utility.ResourceNotFoundError:
// 			utility.SendNotFoundErrResp(c)
// 			return
// 		default:
// 			utility.SendInternalErrResp(vapc.db, c, e.Error())
// 			return
// 		}
// 	}

// 	c.JSON(http.StatusOK, ids)
// }

func (vapc *ViolationReportActionPlanController) CreateViolationReportActionPlan(c *gin.Context) {
	var formDTO violationReportDTOs.ViolationReportActionPlanFormDTO

	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vapc.db, c, err.Error())
		return
	}

	id, err := vapc.service.Create(uint(violationReportID), currentUserID, &formDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(vapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (vapc *ViolationReportActionPlanController) UpdateViolationReportActionPlan(c *gin.Context) {
	var formDTO violationReportDTOs.ViolationReportActionPlanFormDTO

	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vapc.db, c, err.Error())
		return
	}

	err = vapc.service.Update(uint(violationReportID), uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(vapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (vapc *ViolationReportActionPlanController) DeleteViolationReportActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vapc.db, c, err.Error())
		return
	}

	failedIDs, err := vapc.service.Delete(uint(violationReportID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(vapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (vapc *ViolationReportActionPlanController) GetViolationReportActionPlanExport(c *gin.Context) {
	c.Set("export", true)
	vapc.GetAllViolationReportActionPlans(c)
}

func (vapc *ViolationReportActionPlanController) DeleteForceViolationReportActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vapc.db, c, err.Error())
		return
	}

	failedIDs, err := vapc.service.ForceDelete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(vapc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
