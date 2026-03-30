package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategicInitiativeAssociatedItemController struct {
	service services.StrategicInitiativeAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewStrategicInitiativeAssociatedItemController(service services.StrategicInitiativeAssociatedItemServiceInterface) *StrategicInitiativeAssociatedItemController {
	return &StrategicInitiativeAssociatedItemController{
		service: service,
	}
}

func (siai *StrategicInitiativeAssociatedItemController) GetAllStrategicIntiativeRisks(c *gin.Context) {
	intiativeIDStr := c.Param("id")

	strategicIntiativeID, err := strconv.ParseUint(intiativeIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	valueStr := c.Query("value")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(siai.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(siai.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(siai.db, c, "Invalid value parameter")
			return
		}
	}
	allowedParams := []string{
		"current_likelihood_ids",
		"current_impact_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Value:    value,
		Filters:  filters,
	}
	risks, err := siai.service.GetAllStrategicIntiativeRisks(uint(strategicIntiativeID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(siai.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (siai *StrategicInitiativeAssociatedItemController) CreateStrategicIntiativeRisk(c *gin.Context) {
	var input commonDTOs.AssociatedRiskMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	strategicIntiativeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(siai.db, c, err.Error())
		return
	}

	err = siai.service.CreateStrategicIntiativeRisk(uint(strategicIntiativeID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(siai.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (siai *StrategicInitiativeAssociatedItemController) DeleteStrategicIntiativeRisk(c *gin.Context) {
	strategicIntiativeID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	err := siai.service.DeleteStrategicIntiativeRisk(uint(strategicIntiativeID), uint(riskID))
	if err != nil {
		utility.SendInternalErrResp(siai.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Risk unmapped from strategic intiative successfully"})
}
