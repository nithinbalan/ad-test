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

type TestAndExerciseController struct {
	service services.TestAndExerciseServiceInterface
	db      *gorm.DB
}

func NewTestAndExerciseController(service services.TestAndExerciseServiceInterface) *TestAndExerciseController {
	return &TestAndExerciseController{
		service: service,
	}
}

func (tec *TestAndExerciseController) GetAllTestAndExercises(c *gin.Context) {

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"responsible_user_ids",
	}
	validSortFields := map[string]string{
		"title":                          "bcp_test_and_exercises.title",
		"responsible_user_first_name_en": "responsible_user.first_name_en",
		"responsible_user_first_name_ar": "responsible_user.first_name_ar",
		"responsible_user_last_name_en":  "responsible_user.last_name_en",
		"responsible_user_last_name_ar":  "responsible_user.last_name_ar",
		"responsible_user_email":         "responsible_user.email",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(tec.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(tec.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := businessContinuityDTOs.TestAndExerciseListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}
	testAndExercises, err := tec.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(tec.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, testAndExercises)
}

func (tec *TestAndExerciseController) GetTestAndExerciseDetails(c *gin.Context) {
	idStr := c.Param("testAndExerciseId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	testAndExercise, err := tec.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tec.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, testAndExercise)
}

func (tec *TestAndExerciseController) CreateTestAndExercise(c *gin.Context) {

	var testAndExerciseFormDTO businessContinuityDTOs.TestAndExerciseFormDTO

	if err := c.ShouldBindJSON(&testAndExerciseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tec.db, c, err.Error())
		return
	}

	id, err := tec.service.Create(currentUserID, &testAndExerciseFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (tec *TestAndExerciseController) UpdateTestAndExercise(c *gin.Context) {
	var testAndExerciseFormDTO businessContinuityDTOs.TestAndExerciseFormDTO

	testAndExerciseID, err := strconv.Atoi(c.Param("testAndExerciseId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&testAndExerciseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tec.db, c, err.Error())
		return
	}

	err = tec.service.Update(uint(testAndExerciseID), currentUserID, &testAndExerciseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tec.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tec *TestAndExerciseController) DeleteTestAndExercises(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tec.db, c, err.Error())
		return
	}

	failedIDs, err := tec.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tec.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
