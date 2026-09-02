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

type PenetrationTestController struct {
	service services.PenetrationTestServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestController(service services.PenetrationTestServiceInterface) *PenetrationTestController {
	return &PenetrationTestController{
		service: service,
	}
}

func (ptc *PenetrationTestController) GetAllPenetrationTests(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")
	status := c.Query("status")

	validSortFields := map[string]string{
		"reference_number":        "reference_number",
		"system_application_name": "system_application_name",
		"year":                    "year",
		"period":                  "period",
		"created_at":              "created_at",
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

	queryParams := penetrationTestDTOs.PenetrationTestListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
	}
	penetrationTests, err := ptc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(ptc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, penetrationTests)
}

func (ptc *PenetrationTestController) GetPenetrationTestDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	penetrationTest, err := ptc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ptc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, penetrationTest)
}

func (ptc *PenetrationTestController) CreatePenetrationTest(c *gin.Context) {
	var formDTO penetrationTestDTOs.PenetrationTestFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptc.db, c, err.Error())
		return
	}

	id, err := ptc.service.Create(currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (ptc *PenetrationTestController) UpdatePenetrationTest(c *gin.Context) {
	var formDTO penetrationTestDTOs.PenetrationTestFormDTO

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
		utility.SendInternalErrResp(ptc.db, c, err.Error())
		return
	}

	if err := ptc.service.Update(uint(id), currentUserID, &formDTO); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ptc *PenetrationTestController) DeletePenetrationTests(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptc.db, c, err.Error())
		return
	}

	failedIDs, err := ptc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ptc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
