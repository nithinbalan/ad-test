package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/grc/penetration-testing/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenetrationTestRetestController struct {
	service services.PenetrationTestRetestServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestRetestController(service services.PenetrationTestRetestServiceInterface) *PenetrationTestRetestController {
	return &PenetrationTestRetestController{
		service: service,
	}
}

func (ptrec *PenetrationTestRetestController) GetAllPenetrationTestRetests(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "start_date")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	allowedParams := []string{
		"penetration_test_executor_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	validSortFields := map[string]string{
		"reference_number":  "penetration_test_retests.reference_number",
		"start_date":        "penetration_test_retests.start_date",
		"end_date":          "penetration_test_retests.end_date",
		"status":            "penetration_test_retests.status",
		"executor_title_en": "penetration_test_executors.title_en",
		"executor_title_ar": "penetration_test_executors.title_ar",
		"created_at":        "penetration_test_retests.created_at",
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page parameter"))
		return
	}
	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page_size parameter"))
		return
	}
	if page < 1 || pageSize < 1 {
		utility.SendInvalidErrResp(c, errors.New("page and page_size must be greater than zero"))
		return
	}
	if sort != "asc" && sort != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	penetrationTestID, err := utility.ParseQueryParamToInt(c, "penetration_test_id", 0)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid penetration_test_id parameter"))
		return
	}
	if penetrationTestID < 0 {
		utility.SendInvalidErrResp(c, errors.New("penetration_test_id must be greater than zero"))
		return
	}

	queryParams := penetrationTestDTOs.PenetrationTestRetestListQueryParamsDTO{
		Sort:              sort,
		SortBy:            sortField,
		Keyword:           keyword,
		Status:            status,
		Page:              page,
		PageSize:          pageSize,
		PenetrationTestID: uint(penetrationTestID),
		Filters:           filters,
		QuickFilter:       quickFilter,
		DateFrom:          dateFrom,
		DateTo:            dateTo,
	}
	retests, err := ptrec.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(ptrec.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, retests)
}

func (ptrec *PenetrationTestRetestController) GetPenetrationTestRetestDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	retest, err := ptrec.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ptrec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, retest)
}

func (ptrec *PenetrationTestRetestController) CreatePenetrationTestRetest(c *gin.Context) {
	var formDTO penetrationTestDTOs.PenetrationTestRetestFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptrec.db, c, err.Error())
		return
	}
	id, err := ptrec.service.Create(currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptrec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (ptrec *PenetrationTestRetestController) UpdatePenetrationTestRetest(c *gin.Context) {
	var formDTO penetrationTestDTOs.PenetrationTestRetestFormDTO
	id, err := strconv.Atoi(c.Param("id"))
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
		utility.SendInternalErrResp(ptrec.db, c, err.Error())
		return
	}
	if err := ptrec.service.Update(uint(id), currentUserID, &formDTO); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptrec.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ptrec *PenetrationTestRetestController) VerifyPenetrationTestRetestFinding(c *gin.Context) {
	var formDTO penetrationTestDTOs.PenetrationTestRetestFindingFormDTO
	id, err := strconv.Atoi(c.Param("id"))
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
		utility.SendInternalErrResp(ptrec.db, c, err.Error())
		return
	}
	if err := ptrec.service.VerifyFinding(uint(id), currentUserID, &formDTO); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptrec.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ptrec *PenetrationTestRetestController) CompletePenetrationTestRetest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptrec.db, c, err.Error())
		return
	}
	if err := ptrec.service.Complete(uint(id), currentUserID); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptrec.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ptrec *PenetrationTestRetestController) DeletePenetrationTestRetests(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptrec.db, c, err.Error())
		return
	}
	failedIDs, err := ptrec.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ptrec.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
