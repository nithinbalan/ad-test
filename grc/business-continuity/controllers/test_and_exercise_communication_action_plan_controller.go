package controllers

import (
	"errors"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestAndExerciseCommunicationActionPlanController struct {
	service services.TestAndExerciseCommunicationActionPlanServiceInterface
	db      *gorm.DB
}

func NewTestAndExerciseCommunicationActionPlanController(service services.TestAndExerciseCommunicationActionPlanServiceInterface) *TestAndExerciseCommunicationActionPlanController {
	return &TestAndExerciseCommunicationActionPlanController{
		service: service,
	}
}

func (tecapc *TestAndExerciseCommunicationActionPlanController) GetAllTestAndExerciseCommunicationActionPlans(c *gin.Context) {

	idStr := c.Param("testAndExerciseId")

	testAndExerciseId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"responsible_user_ids",
		"bcp_phase_ids",
	}
	validSortFields := map[string]string{
		"title":                          "test_and_exercise_communication_action_plans.title",
		"responsible_user_first_name_en": "responsible_user.first_name_en",
		"responsible_user_first_name_ar": "responsible_user.first_name_ar",
		"responsible_user_last_name_en":  "responsible_user.last_name_en",
		"responsible_user_last_name_ar":  "responsible_user.last_name_ar",
		"responsible_user_email":         "responsible_user.email",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(tecapc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(tecapc.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := businessContinuityDTOs.TestAndExerciseCommunicationActionPlanListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}
	communicationActionPlans, err := tecapc.service.GetAll(uint(testAndExerciseId), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(tecapc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, communicationActionPlans)
}

func (tecapc *TestAndExerciseCommunicationActionPlanController) GetTestAndExerciseCommunicationActionPlanDetails(c *gin.Context) {
	testAndExerciseIdStr := c.Param("testAndExerciseId")

	testAndExerciseId, err := strconv.ParseUint(testAndExerciseIdStr, 10, 32)
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

	bcpCommunicationPlan, err := tecapc.service.GetFullDetailsByID(uint(testAndExerciseId), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tecapc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, bcpCommunicationPlan)
}

func (tecapc *TestAndExerciseCommunicationActionPlanController) CreateTestAndExerciseCommunicationActionPlan(c *gin.Context) {

	idStr := c.Param("testAndExerciseId")

	testAndExerciseId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var bcpCommunicationPlanFormDTO businessContinuityDTOs.TestAndExerciseCommunicationActionPlanFormDTO

	if err := c.ShouldBindJSON(&bcpCommunicationPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tecapc.db, c, err.Error())
		return
	}

	id, err := tecapc.service.Create(uint(testAndExerciseId), currentUserID, &bcpCommunicationPlanFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tecapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (tecapc *TestAndExerciseCommunicationActionPlanController) UpdateTestAndExerciseCommunicationActionPlan(c *gin.Context) {
	testAndExerciseIdStr := c.Param("testAndExerciseId")

	testAndExerciseId, err := strconv.ParseUint(testAndExerciseIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var bcpCommunicationPlanFormDTO businessContinuityDTOs.TestAndExerciseCommunicationActionPlanFormDTO

	bcpCommunicationPlanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&bcpCommunicationPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tecapc.db, c, err.Error())
		return
	}

	err = tecapc.service.Update(uint(testAndExerciseId), uint(bcpCommunicationPlanID), currentUserID, &bcpCommunicationPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tecapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tecapc *TestAndExerciseCommunicationActionPlanController) DeleteTestAndExerciseCommunicationActionPlans(c *gin.Context) {
	testAndExerciseIdStr := c.Param("testAndExerciseId")

	testAndExerciseId, err := strconv.ParseUint(testAndExerciseIdStr, 10, 32)
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
		utility.SendInternalErrResp(tecapc.db, c, err.Error())
		return
	}

	failedIDs, err := tecapc.service.Delete(uint(testAndExerciseId), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tecapc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
