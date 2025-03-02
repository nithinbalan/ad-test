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

type BcpTestAndExerciseController struct {
	service services.BcpTestAndExerciseServiceInterface
	db      *gorm.DB
}

func NewBcpTestAndExerciseController(service services.BcpTestAndExerciseServiceInterface) *BcpTestAndExerciseController {
	return &BcpTestAndExerciseController{
		service: service,
	}
}

func (btec *BcpTestAndExerciseController) GetAllBcpTestAndExercises(c *gin.Context) {

	idStr := c.Param("bcpId")

	bcpId, err := strconv.ParseUint(idStr, 10, 32)
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
		"review_frequency_ids",
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
		utility.SendInternalErrResp(btec.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(btec.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := businessContinuityDTOs.BcpTestAndExerciseListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}
	bcpTestAndExercises, err := btec.service.GetAll(uint(bcpId), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(btec.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, bcpTestAndExercises)
}

func (btec *BcpTestAndExerciseController) GetBcpTestAndExerciseDetails(c *gin.Context) {
	bcpIDStr := c.Param("bcpId")

	bcpId, err := strconv.ParseUint(bcpIDStr, 10, 64)
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

	bcpTestAndExercise, err := btec.service.GetFullDetailsByID(uint(bcpId), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(btec.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, bcpTestAndExercise)
}

func (btec *BcpTestAndExerciseController) CreateBcpTestAndExercise(c *gin.Context) {

	idStr := c.Param("bcpId")

	bcpId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var bcpTestAndExerciseFormDTO businessContinuityDTOs.BcpTestAndExerciseFormDTO

	if err := c.ShouldBindJSON(&bcpTestAndExerciseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(btec.db, c, err.Error())
		return
	}

	id, err := btec.service.Create(uint(bcpId), currentUserID, &bcpTestAndExerciseFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(btec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (btec *BcpTestAndExerciseController) UpdateBcpTestAndExercise(c *gin.Context) {
	var bcpTestAndExerciseFormDTO businessContinuityDTOs.BcpTestAndExerciseFormDTO

	bcpIDStr := c.Param("bcpId")

	bcpId, err := strconv.ParseUint(bcpIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	bcpTestAndExerciseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&bcpTestAndExerciseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(btec.db, c, err.Error())
		return
	}

	err = btec.service.Update(uint(bcpId), uint(bcpTestAndExerciseID), currentUserID, &bcpTestAndExerciseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(btec.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (btec *BcpTestAndExerciseController) DeleteBcpTestAndExercises(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	bcpIDStr := c.Param("bcpId")

	bcpId, err := strconv.ParseUint(bcpIDStr, 10, 64)
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
		utility.SendInternalErrResp(btec.db, c, err.Error())
		return
	}

	failedIDs, err := btec.service.Delete(uint(bcpId), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(btec.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
