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

type ComplianceRequirementAssociatedItemController struct {
	service services.ComplianceRequirementAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewComplianceRequirementAssociatedItemController(service services.ComplianceRequirementAssociatedItemServiceInterface) *ComplianceRequirementAssociatedItemController {
	return &ComplianceRequirementAssociatedItemController{
		service: service,
	}
}

func (craic *ComplianceRequirementAssociatedItemController) GetAllComplianceRequirementControls(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	requirementIDStr := c.Param("requirementId")

	complianceRequirementID, err := strconv.ParseUint(requirementIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	queryParams := complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}
	compliances, err := craic.service.GetAllControls(uint(complianceID), uint(complianceRequirementID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(craic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, compliances)
}

func (craic *ComplianceRequirementAssociatedItemController) GetAllComplianceRequirementStrategicObjectives(c *gin.Context) {

	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	requirementIDStr := c.Param("requirementId")

	complianceRequirementID, err := strconv.ParseUint(requirementIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	queryParams := complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}
	compliances, err := craic.service.GetAllStrategicObjectives(uint(complianceID), uint(complianceRequirementID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(craic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, compliances)
}

func (craic *ComplianceRequirementAssociatedItemController) GetAllComplianceRequirementRisks(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	requirementIDStr := c.Param("requirementId")

	complianceRequirementID, err := strconv.ParseUint(requirementIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	queryParams := complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}
	compliances, err := craic.service.GetAllRisks(uint(complianceID), uint(complianceRequirementID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(craic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, compliances)
}

func (craic *ComplianceRequirementAssociatedItemController) GetAllComplianceRequirementIncidents(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	requirementIDStr := c.Param("requirementId")

	complianceRequirementID, err := strconv.ParseUint(requirementIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	queryParams := complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}
	compliances, err := craic.service.GetAllIncidents(uint(complianceID), uint(complianceRequirementID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(craic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, compliances)
}
