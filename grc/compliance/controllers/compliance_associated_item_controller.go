package controllers

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceAssociatedItemController struct {
	service services.ComplianceAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewComplianceAssociatedItemController(service services.ComplianceAssociatedItemServiceInterface) *ComplianceAssociatedItemController {
	return &ComplianceAssociatedItemController{
		service: service,
	}
}

func (caic *ComplianceAssociatedItemController) GetAllComplianceControls(c *gin.Context) {

	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
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

	queryParams := complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	compliances, err := caic.service.GetAllControls(uint(complianceID), &queryParams)

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

func (caic *ComplianceAssociatedItemController) GetAllComplianceRisks(c *gin.Context) {

	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	valueStr := c.Query("value")

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

	queryParams := complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Value:       value,
		QuickFilter: quickFilter,
	}
	compliances, err := caic.service.GetAllRisks(uint(complianceID), &queryParams)

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
	c.JSON(http.StatusOK, compliances)
}

func (caic *ComplianceAssociatedItemController) GetAllComplianceIncidents(c *gin.Context) {

	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
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

	queryParams := complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	compliances, err := caic.service.GetAllIncidents(uint(complianceID), &queryParams)
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
	c.JSON(http.StatusOK, compliances)
}

func (caic *ComplianceAssociatedItemController) GetAllComplianceStrategicObjectives(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
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

	queryParams := complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	compliances, err := caic.service.GetAllStrategicObjectives(uint(complianceID), &queryParams)
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
	c.JSON(http.StatusOK, compliances)
}

func (caic *ComplianceAssociatedItemController) GetAllAssociatedCounts(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	associatedCounts, err := caic.service.GetAllAssociatedCounts(uint(complianceID))
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

func (caic *ComplianceAssociatedItemController) GetComplianceRequirementCountByStatus(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	requirementCounts, err := caic.service.GetComplianceRequirementCountByStatus(uint(complianceID))
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
	c.JSON(http.StatusOK, requirementCounts)
}

func (caic *ComplianceAssociatedItemController) CreateControlsToCompliance(c *gin.Context) {
	var input complianceDTOs.ComplianceRegisterControlMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	complianceID, err := strconv.ParseUint(c.Param("complianceId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	err = caic.service.CreateControlsToCompliance(uint(complianceID), currentUserID, &input)
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

func (caic *ComplianceAssociatedItemController) DeleteControlFromCompliance(c *gin.Context) {
	complianceID, _ := strconv.ParseUint(c.Param("complianceId"), 10, 64)
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)

	err := caic.service.DeleteControlFromCompliance(uint(complianceID), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Control unmapped from compliance successfully"})
}
