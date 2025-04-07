package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/grc/organizational-context/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StakeholderNeedsAndExpectationController struct {
	service services.StakeholderNeedsAndExpectationServiceInterface
	db      *gorm.DB
}

func NewStakeholderNeedsAndExpectationController(service services.StakeholderNeedsAndExpectationServiceInterface) *StakeholderNeedsAndExpectationController {
	return &StakeholderNeedsAndExpectationController{
		service: service,
	}
}

func (snec *StakeholderNeedsAndExpectationController) GetAllStakeholderNeedsAndExpectations(c *gin.Context) {
	stakeholderIDStr := c.Param("stakeholderId")
	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(snec.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"evaluation_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(snec.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(snec.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := organizationalContextDTOs.StakeholderNeedsAndExpectationListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      export,
	}

	needsAndExpectation, err := snec.service.GetAll(uint(stakeholderID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(snec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, needsAndExpectation)
}

func (snec *StakeholderNeedsAndExpectationController) GetStakeholderNeedsAndExpectationDetails(c *gin.Context) {
	stakeholderIDStr := c.Param("stakeholderId")
	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("needsAndExpectationId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	needsAndExpectation, err := snec.service.GetFullDetailsByID(uint(stakeholderID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(snec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, needsAndExpectation)
}

func (snec *StakeholderNeedsAndExpectationController) CreateStakeholderNeedsAndExpectation(c *gin.Context) {
	var NeedsAndExpectationFormDTO organizationalContextDTOs.StakeholderNeedsAndExpectationFormDTO

	stakeholderIDStr := c.Param("stakeholderId")
	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&NeedsAndExpectationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(snec.db, c, err.Error())
		return
	}

	id, err := snec.service.Create(uint(stakeholderID), currentUserID, &NeedsAndExpectationFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(snec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (snec *StakeholderNeedsAndExpectationController) UpdateStakeholderNeedsAndExpectation(c *gin.Context) {
	var NeedsAndExpectationFormDTO organizationalContextDTOs.StakeholderNeedsAndExpectationFormDTO

	stakeholderIDStr := c.Param("stakeholderId")
	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("needsAndExpectationId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&NeedsAndExpectationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(snec.db, c, err.Error())
		return
	}

	err = snec.service.Update(uint(stakeholderID), uint(id), currentUserID, &NeedsAndExpectationFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(snec.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (snec *StakeholderNeedsAndExpectationController) DeleteStakeholderNeedsAndExpectations(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	stakeholderIDStr := c.Param("stakeholderId")
	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 32)
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
		utility.SendInternalErrResp(snec.db, c, err.Error())
		return
	}

	failedIDs, err := snec.service.Delete(uint(stakeholderID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(snec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
