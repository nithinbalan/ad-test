package controllers

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControlAssociatedItemController struct {
	service services.ControlAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewControlAssociatedItemController(service services.ControlAssociatedItemServiceInterface) *ControlAssociatedItemController {
	return &ControlAssociatedItemController{
		service: service,
	}
}

func (caic *ControlAssociatedItemController) GetAllControlProcesses(c *gin.Context) {
	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := organizationDTOs.ControlAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	processes, err := caic.service.GetAllProcesses(uint(controlID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(caic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, processes)
}

func (caic *ControlAssociatedItemController) GetAllControlComplianceRequirements(c *gin.Context) {
	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := organizationDTOs.ControlAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	compliances, err := caic.service.GetAllComplianceRequirements(uint(controlID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(caic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, compliances)
}

func (caic *ControlAssociatedItemController) GetAllControlRisks(c *gin.Context) {
	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	valueStr := c.Query("value")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(caic.db, c, "Invalid value parameter")
			return
		}
	}

	queryParams := organizationDTOs.ControlAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Value:       value,
		QuickFilter: quickFilter,
	}
	risks, err := caic.service.GetAllRisks(uint(controlID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(caic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (caic *ControlAssociatedItemController) GetAllControlIncidents(c *gin.Context) {
	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := organizationDTOs.ControlAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	incidents, err := caic.service.GetAllIncidents(uint(controlID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(caic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidents)
}

func (caic *ControlAssociatedItemController) GetAllAssociatedCounts(c *gin.Context) {
	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	associatedCounts, err := caic.service.GetAllAssociatedCounts(uint(controlID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(caic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, associatedCounts)
}

func (caic *ControlAssociatedItemController) GetAllClauses(c *gin.Context) {
	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := organizationDTOs.ControlAssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	associatedClauses, err := caic.service.GetAllClauses(uint(controlID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(caic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, associatedClauses)
}

func (caic *ControlAssociatedItemController) UpdateClauses(c *gin.Context) {
	var ControlClauseFormDTO organizationDTOs.ControlClauseFormDTO

	controlID, err := strconv.Atoi(c.Param("controlId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&ControlClauseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	err = caic.service.UpdateClauses(uint(controlID), currentUserID, &ControlClauseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(caic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Clauses against control updated successfully"})
}

func (caic *ControlAssociatedItemController) CreateControlRisk(c *gin.Context) {
	var input organizationDTOs.AssociatedRiskMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	controlID, err := strconv.ParseUint(c.Param("controlId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	err = caic.service.CreateControlRisk(uint(controlID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(caic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (caic *ControlAssociatedItemController) DeleteControlRisk(c *gin.Context) {
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	err := caic.service.DeleteControlRisk(uint(controlID), uint(riskID))
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Risk unmapped from control successfully"})
}
