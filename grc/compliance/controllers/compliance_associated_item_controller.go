package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
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


func (caic *ComplianceAssociatedItemController) GetAllComplianceAssets(c *gin.Context) {
	complianceID, err := strconv.Atoi(c.Param("complianceId"))
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

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := caic.service.GetAllComplianceAssets(uint(complianceID), &queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, risks)
}

func (caic *ComplianceAssociatedItemController) GetAllComplianceProcesses(c *gin.Context) {
	complianceID, err := strconv.Atoi(c.Param("complianceId"))
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

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := caic.service.GetAllComplianceProcesses(uint(complianceID), &queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, risks)
}

func (caic *ComplianceAssociatedItemController) GetAllComplianceKPI(c *gin.Context) {
	complianceID, err := strconv.Atoi(c.Param("complianceId"))
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

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := caic.service.GetAllComplianceKPI(uint(complianceID), &queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, risks)
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

func (caic *ComplianceAssociatedItemController) CreateComplianceControls(c *gin.Context) {
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

	err = caic.service.CreateComplianceControls(uint(complianceID), currentUserID, &input)
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

func (caic *ComplianceAssociatedItemController) DeleteComplianceControl(c *gin.Context) {
	complianceID, _ := strconv.ParseUint(c.Param("complianceId"), 10, 64)
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)

	err := caic.service.DeleteComplianceControl(uint(complianceID), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Control unmapped from compliance successfully"})
}

func (caic *ComplianceAssociatedItemController) CreateComplianceRisk(c *gin.Context) {
	var input commonDTOs.AssociatedRiskMapDTO
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

	err = caic.service.CreateComplianceRisk(uint(complianceID), currentUserID, &input)
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

func (caic *ComplianceAssociatedItemController) DeleteComplianceRisk(c *gin.Context) {
	complianceID, _ := strconv.ParseUint(c.Param("complianceId"), 10, 64)
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	err := caic.service.DeleteComplianceRisk(uint(complianceID), uint(riskID))
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Risk unmapped from compliance successfully"})
}

func (caic *ComplianceAssociatedItemController) GetAllComplianceDocuments(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
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

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	documents, err := caic.service.GetAllDocuments(uint(complianceID), &queryParams)

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
	c.JSON(http.StatusOK, documents)
}

func (caic *ComplianceAssociatedItemController) CreateComplianceDocument(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
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

	err = caic.service.CreateComplianceDocument(uint(complianceID), currentUserID, &input)
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

func (caic *ComplianceAssociatedItemController) DeleteComplianceDocument(c *gin.Context) {
	complianceID, _ := strconv.ParseUint(c.Param("complianceId"), 10, 64)
	documentID, _ := strconv.ParseUint(c.Param("documentId"), 10, 64)
	err := caic.service.DeleteComplianceDocument(uint(complianceID), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from compliance register successfully"})
}

func (caic *ComplianceAssociatedItemController) GetComplianceControlCountByStatus(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlCounts, err := caic.service.GetComplianceControlCountByStatus(uint(complianceID))
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

	c.JSON(http.StatusOK, controlCounts)
}
