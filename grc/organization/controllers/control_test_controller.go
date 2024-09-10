package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControlTestController struct {
	service services.ControlTestServiceInterface
	db      *gorm.DB
}

func NewControlTestController(service services.ControlTestServiceInterface) *ControlTestController {
	return &ControlTestController{
		service: service,
	}
}

func (ctc *ControlTestController) GetAllControlTests(c *gin.Context) {
	controlID, err := strconv.Atoi(c.Param("controlId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"control_implementation_status_ids",
		"control_efficiency_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	//Parse tested_at_from, tested_at_to parameters
	testedAtFrom := c.Query("tested_at_from")
	testedAtTo := c.Query("tested_at_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(ctc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(ctc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := organizationDTOs.ControlTestListQueryParamsDTO{
		Sort:         sort,
		SortBy:       sortBy,
		Keyword:      keyword,
		Status:       status,
		Page:         page,
		PageSize:     pageSize,
		Filters:      filters,
		QuickFilter:  quickFilter,
		TestedAtFrom: testedAtFrom,
		TestedAtTo:   testedAtTo,
	}
	controlTests, err := ctc.service.GetAll(uint(controlID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(ctc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, controlTests)
}

func (crc *ControlTestController) GetControlTestDetails(c *gin.Context) {
	idStr := c.Param("controlTestId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlTest, err := crc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controlTest)
}

func (ctc *ControlTestController) CreateControlTest(c *gin.Context) {
	var controlTestFormDTO organizationDTOs.ControlTestFormDTO

	if err := c.ShouldBindJSON(&controlTestFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	controlID, err := strconv.Atoi(c.Param("controlId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ctc.db, c, err.Error())
		return
	}

	id, err := ctc.service.Create(uint(controlID), currentUserID, nil, &controlTestFormDTO, "")

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ctc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (ctc *ControlTestController) UpdateControlTest(c *gin.Context) {

	controlID, err := strconv.Atoi(c.Param("controlId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlTestID, err := strconv.Atoi(c.Param("controlTestId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var controlTestFormDTO organizationDTOs.ControlTestFormDTO
	if err := c.ShouldBindJSON(&controlTestFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ctc.db, c, err.Error())
		return
	}

	err = ctc.service.Update(uint(controlID), uint(controlTestID), currentUserID, &controlTestFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ctc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)

}

func (ctc *ControlTestController) DeleteControlTest(c *gin.Context) {
	controlID, err := strconv.Atoi(c.Param("controlId"))
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
		utility.SendInternalErrResp(ctc.db, c, err.Error())
		return
	}

	failedIDs, err := ctc.service.Delete(uint(controlID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ctc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ctc *ControlTestController) GetControlTestTrendAnalysis(c *gin.Context) {
	idStr := c.Param("controlId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlTrend, err := ctc.service.GetTrendAnalysis(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ctc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controlTrend)
}
