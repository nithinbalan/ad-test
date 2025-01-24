package controllers

import (
	"errors"
	"net/http"
	"strconv"

	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestAndExerciseRecoveryActionPlanController struct {
	service services.TestAndExerciseRecoveryActionPlanServiceInterface
	db      *gorm.DB
}

func NewTestAndExerciseRecoveryActionPlanController(service services.TestAndExerciseRecoveryActionPlanServiceInterface) *TestAndExerciseRecoveryActionPlanController {
	return &TestAndExerciseRecoveryActionPlanController{
		service: service,
	}
}

func (terapc *TestAndExerciseRecoveryActionPlanController) GetAllTestAndExerciseRecoveryActionPlans(c *gin.Context) {
	idStr := c.Param("testAndExerciseId")

	testAndExerciseID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")

	allowedParams := []string{
		"bcp_phase_ids",
		"process_ids",
		"responsible_user_ids",
		"severity_level_ids",
		"action_plan_status_id",
	}
	validSortFields := map[string]string{
		"title":                          "test_and_exercise_recovery_action_plans.title",
		"responsible_user_first_name_en": "responsible_user.first_name_en",
		"responsible_user_first_name_ar": "responsible_user.first_name_ar",
		"responsible_user_last_name_en":  "responsible_user.last_name_en",
		"responsible_user_last_name_ar":  "responsible_user.last_name_ar",
		"responsible_user_email":         "responsible_user.email",
		"responsible_unit_title_en":      "organization_units.title_en",
		"responsible_unit_title_ar":      "organization_units.title_ar",
		"severity_level_title_en":        "severity_levels.title_en",
		"severity_level_title_ar":        "severity_levels.title_ar",
		"action_plan_status_title_en":    "action_plan_statuses.title_en",
		"action_plan_status_title_ar":    "action_plan_statuses.title_ar",
		"created_user_first_name_en":     "created_user.first_name_en",
		"created_user_last_name_en":      "created_user.last_name_en",
		"updated_user_first_name_en":     "updated_user.first_name_en",
		"updated_user_last_name_en":      "updated_user.last_name_en",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := businessContinuityDTOs.TestAndExerciseRecoveryActionPlanQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortField,
		Keyword: keyword,
		Filters: filters,
	}

	actionPlans, err := terapc.service.GetAll(uint(testAndExerciseID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(terapc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, actionPlans)
}

func (terapc *TestAndExerciseRecoveryActionPlanController) GetByID(c *gin.Context) {
	testAndExerciseIDStr := c.Param("testAndExerciseId")
	testAndExerciseID, err := strconv.ParseUint(testAndExerciseIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlan, err := terapc.service.GetByID(uint(testAndExerciseID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(terapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, actionPlan)
}

func (terapc *TestAndExerciseRecoveryActionPlanController) GetTestAndExerciseRecoveryActionPlanDetails(c *gin.Context) {
	testAndExerciseIDStr := c.Param("testAndExerciseId")
	testAndExerciseID, err := strconv.ParseUint(testAndExerciseIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlan, err := terapc.service.GetFullDetailsByID(uint(testAndExerciseID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(terapc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, actionPlan)
}

func (terapc *TestAndExerciseRecoveryActionPlanController) CreateTestAndExerciseRecoveryActionPlan(c *gin.Context) {
	idStr := c.Param("testAndExerciseId")
	testAndExerciseID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO businessContinuityDTOs.TestAndExerciseRecoveryActionPlanFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(terapc.db, c, err.Error())
		return
	}

	id, err := terapc.service.Create(uint(testAndExerciseID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(terapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (terapc *TestAndExerciseRecoveryActionPlanController) UpdateTestAndExerciseRecoveryActionPlan(c *gin.Context) {
	idStr := c.Param("testAndExerciseId")
	testAndExerciseID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO businessContinuityDTOs.TestAndExerciseRecoveryActionPlanFormDTO

	actionPlanID, err := strconv.Atoi(c.Param("id"))
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
		utility.SendInternalErrResp(terapc.db, c, err.Error())
		return
	}

	err = terapc.service.Update(uint(testAndExerciseID), uint(actionPlanID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(terapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (terapc *TestAndExerciseRecoveryActionPlanController) DeleteTestAndExerciseRecoveryActionPlans(c *gin.Context) {
	idStr := c.Param("testAndExerciseId")
	testAndExerciseID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(terapc.db, c, err.Error())
		return
	}

	failedIDs, err := terapc.service.Delete(uint(testAndExerciseID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(terapc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
