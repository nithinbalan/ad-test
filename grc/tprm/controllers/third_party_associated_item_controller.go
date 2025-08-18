package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/tprm/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThirdPartyAssociatedItemController struct {
	service services.ThirdPartyAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewThirdPartyAssociatedItemController(service services.ThirdPartyAssociatedItemServiceInterface) *ThirdPartyAssociatedItemController {
	return &ThirdPartyAssociatedItemController{
		service: service,
	}
}

func (paic *ThirdPartyAssociatedItemController) GetAllThirdPartyProcesses(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := paic.service.GetAllThirdPartyProcesses(uint(ThirdPartyID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}

func (paic *ThirdPartyAssociatedItemController) CreateThirdPartyProcess(c *gin.Context) {
	var input commonDTOs.AssociatedProcessMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	ThirdPartyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	err = paic.service.CreateThirdPartyProcess(uint(ThirdPartyID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (paic *ThirdPartyAssociatedItemController) DeleteThirdPartyProcess(c *gin.Context) {
	ThirdPartyID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	processID, _ := strconv.ParseUint(c.Param("processId"), 10, 64)
	err := paic.service.DeleteThirdPartyProcess(uint(ThirdPartyID), uint(processID))
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Process unmapped from  ThirdParty successfully"})
}

func (paic *ThirdPartyAssociatedItemController) GetAllThirdPartyStrategicObjectives(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	objectives, err := paic.service.GetAllThirdPartyStrategicObjectives(uint(ThirdPartyID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, objectives)
}

func (paic *ThirdPartyAssociatedItemController) CreateThirdPartyStrategicObjective(c *gin.Context) {
	var input commonDTOs.AssociatedStrategicObjectiveMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	ThirdPartyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	err = paic.service.CreateThirdPartyStrategicObjective(uint(ThirdPartyID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (paic *ThirdPartyAssociatedItemController) DeleteThirdPartyStrategicObjective(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")
	strategicObjectiveIDStr := c.Param("strategicObjectiveId")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategicObjectiveID, err := strconv.ParseUint(strategicObjectiveIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = paic.service.DeleteThirdPartyStrategicObjective(uint(ThirdPartyID), uint(strategicObjectiveID))
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Strategic Objective unmapped from third party successfully"})
}
