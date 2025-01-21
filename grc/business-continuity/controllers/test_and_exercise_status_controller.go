package controllers

import (
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestAndExerciseStatusController struct {
	service services.TestAndExerciseStatusServiceInterface
	db      *gorm.DB
}

func NewTestAndExerciseStatusController(service services.TestAndExerciseStatusServiceInterface) *TestAndExerciseStatusController {
	return &TestAndExerciseStatusController{
		service: service,
	}
}

func (tesc *TestAndExerciseStatusController) GetAllTestAndExerciseStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	testAndExerciseStatuses, err := tesc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(tesc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, testAndExerciseStatuses)
}

func (tesc *TestAndExerciseStatusController) GetTestAndExerciseStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	testAndExerciseStatus, err := tesc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, testAndExerciseStatus)
}
