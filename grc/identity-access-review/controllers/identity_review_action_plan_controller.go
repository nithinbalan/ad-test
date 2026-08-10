package controllers

import (
	"errors"
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IdentityReviewActionPlanController struct {
	service services.IdentityReviewActionPlanServiceInterface
	db      *gorm.DB
}

func NewIdentityReviewActionPlanController(service services.IdentityReviewActionPlanServiceInterface) *IdentityReviewActionPlanController {
	return &IdentityReviewActionPlanController{
		service: service,
	}
}

func (irapc *IdentityReviewActionPlanController) GetAllIdentityReviewActionPlans(c *gin.Context) {
	idStr := c.Param("id")

	var userID uint64
	var err error

	if idStr != "" {
		userID, err = strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		userID = 0
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	exportFromContext, exists := c.Get("export")
	var export bool

	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(irapc.db, c, "Invalid export parameter")
			return
		}
	}

	allowedParams := []string{
		"severity_level_ids",
		"action_plan_status_ids",
		"responsible_user_ids",
		"review_user_ids",
		"identity_review_user_ids",
	}

	validSortFields := map[string]string{
		"id":                             "identity_review_action_plans.id",
		"title":                          "identity_review_action_plans.title",
		"target_date":                    "identity_review_action_plans.target_date",
		"created_at":                     "identity_review_action_plans.created_at",
		"updated_at":                     "identity_review_action_plans.updated_at",
		"severity_level_title_en":        "severity_levels.title_en",
		"severity_level_title_ar":        "severity_levels.title_ar",
		"responsible_user_first_name_en": "responsible_user.first_name_en",
		"responsible_user_first_name_ar": "responsible_user.first_name_ar",
		"action_plan_status_title_en":    "action_plan_statuses.title_en",
		"action_plan_status_title_ar":    "action_plan_statuses.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(irapc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(irapc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := dtos.IdentityReviewActionPlanListQueryParamsDTO{
		Sort:           sort,
		SortBy:         dbField,
		Keyword:        keyword,
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

	identityReviewActionPlans, err := irapc.service.GetAll(uint(userID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapc.db, c, e.Error())
			return
		}
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(irapc.db, c, err.Error())
			return
		}

		excelFile, err := irapc.service.Export(identityReviewActionPlans.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(irapc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Identity review action plans.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, identityReviewActionPlans)
}

func (irapc *IdentityReviewActionPlanController) ExportIdentityReviewActionPlans(c *gin.Context) {
	c.Set("export", true)

	irapc.GetAllIdentityReviewActionPlans(c)
}

func (irapc *IdentityReviewActionPlanController) GetIdentityReviewActionPlanSummary(c *gin.Context) {
	idStr := c.Param("id")

	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	result, err := irapc.service.GetSummary(uint(userID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (irapc *IdentityReviewActionPlanController) GetIdentityReviewActionPlanDetails(c *gin.Context) {
	idStr := c.Param("id")

	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	identityReviewActionPlan, err := irapc.service.GetFullDetailsByID(uint(userID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, identityReviewActionPlan)
}

func (irapc *IdentityReviewActionPlanController) CreateIdentityReviewActionPlan(c *gin.Context) {
	var form dtos.IdentityReviewActionPlanFormDTO

	idStr := c.Param("id")

	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irapc.db, c, err.Error())
		return
	}

	id, err := irapc.service.Create(uint(userID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (irapc *IdentityReviewActionPlanController) UpdateIdentityReviewActionPlan(c *gin.Context) {
	var form dtos.IdentityReviewActionPlanFormDTO

	idStr := c.Param("id")

	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irapc.db, c, err.Error())
		return
	}

	err = irapc.service.Update(uint(userID), uint(id), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (irapc *IdentityReviewActionPlanController) DeleteIdentityReviewActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	idStr := c.Param("id")

	userID, err := strconv.ParseUint(idStr, 10, 64)
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
		utility.SendInternalErrResp(irapc.db, c, err.Error())
		return
	}

	failedIDs, err := irapc.service.Delete(uint(userID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (irapc *IdentityReviewActionPlanController) ForceDeleteIdentityReviewActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irapc.db, c, err.Error())
		return
	}

	failedIDs, err := irapc.service.ForceDelete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(irapc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
