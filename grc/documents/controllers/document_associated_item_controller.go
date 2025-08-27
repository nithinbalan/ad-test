package controllers

import (
	documentDTOs "jethur/grc/documents/dtos"
	documentService "jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentAssociatedItemController struct {
	service documentService.DocumentAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewDocumentAssociatedItemController(service documentService.DocumentAssociatedItemServiceInterface) *DocumentAssociatedItemController {
	return &DocumentAssociatedItemController{
		service: service,
	}
}

func (daic *DocumentAssociatedItemController) GetAllDocumentControls(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := daic.service.GetAllControls(uint(documentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}

func (daic *DocumentAssociatedItemController) GetAllDocumentComplianceRequirements(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
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
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := documentDTOs.DocumentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	complianceRequirements, err := daic.service.GetAllComplianceRequirements(uint(documentID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, complianceRequirements)
}

func (daic *DocumentAssociatedItemController) GetAllDocumentProcesses(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	processes, err := daic.service.GetAllProcesses(uint(documentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, processes)
}

func (daic *DocumentAssociatedItemController) GetAllDocumentSLAs(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
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
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := documentDTOs.DocumentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	slas, err := daic.service.GetAllSLAs(uint(documentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, slas)
}

func (daic *DocumentAssociatedItemController) GetAllDocumentIncidents(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
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
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := documentDTOs.DocumentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	incidents, err := daic.service.GetAllIncidents(uint(documentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidents)
}

func (daic *DocumentAssociatedItemController) GetAllDocumentRisks(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
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
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}
	allowedParams := []string{
		"current_likelihood_ids",
		"current_impact_ids",
	}
	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(daic.db, c, "Invalid value parameter")
			return
		}
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
	risks, err := daic.service.GetAllRisks(uint(documentID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, risks)
}
func (daic *DocumentAssociatedItemController) GetAllDocumentClauses(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	clauses, err := daic.service.GetAllClauses(uint(documentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, clauses)
}
func (daic *DocumentAssociatedItemController) UpdateClauses(c *gin.Context) {
	var clauseFormDTO commonDTOs.AssociatedClauseFormDTO

	documentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&clauseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}

	err = daic.service.UpdateClauses(uint(documentID), currentUserID, &clauseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Clauses against document updated successfully"})
}

func (daic *DocumentAssociatedItemController) CreateDocumentRisk(c *gin.Context) {
	var input commonDTOs.AssociatedRiskMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	documentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}

	err = daic.service.CreateDocumentRisk(uint(documentID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (daic *DocumentAssociatedItemController) DeleteDocumentRisk(c *gin.Context) {
	documentID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	err := daic.service.DeleteDocumentRisk(uint(documentID), uint(riskID))
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Risk unmapped from document successfully"})
}
func (daic *DocumentAssociatedItemController) CreateDocumentControl(c *gin.Context) {
	var input commonDTOs.AssociatedControlMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	documentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}

	err = daic.service.CreateDocumentControl(uint(documentID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (daic *DocumentAssociatedItemController) DeleteDocumentControl(c *gin.Context) {
	documentID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)
	err := daic.service.DeleteDocumentControl(uint(documentID), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Control unmapped from document successfully"})
}

func (daic *DocumentAssociatedItemController) CreateDocumentProcess(c *gin.Context) {
	var input commonDTOs.AssociatedProcessMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	documentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}

	err = daic.service.CreateDocumentProcess(uint(documentID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (daic *DocumentAssociatedItemController) DeleteDocumentProcess(c *gin.Context) {
	documentID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	processID, _ := strconv.ParseUint(c.Param("processId"), 10, 64)
	err := daic.service.DeleteDocumentProcess(uint(documentID), uint(processID))
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Process unmapped from document successfully"})
}

func (daic *DocumentAssociatedItemController) GetAllDocumentStrategicObjectives(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	objectives, err := daic.service.GetAllStrategicObjectives(uint(documentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, objectives)
}
func (daic *DocumentAssociatedItemController) CreateDocumentStrategicObjective(c *gin.Context) {
	var input commonDTOs.AssociatedStrategicObjectiveMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	documentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}

	err = daic.service.CreateDocumentStrategicObjective(uint(documentID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (daic *DocumentAssociatedItemController) DeleteDocumentStrategicObjective(c *gin.Context) {
	documentID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	strategicObjectiveID, _ := strconv.ParseUint(c.Param("strategicObjectiveId"), 10, 64)
	err := daic.service.DeleteDocumentStrategicObjective(uint(documentID), uint(strategicObjectiveID))
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Strategic Objective unmapped from document successfully"})
}

func (daic *DocumentAssociatedItemController) GetAllDocumentKPI(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	kpis, err := daic.service.GetAllDocumentKPI(uint(documentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, kpis)
}
func (daic *DocumentAssociatedItemController) CreateDocumentKPI(c *gin.Context) {
	var input commonDTOs.AssociatedKPIMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	documentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}

	err = daic.service.CreateDocumentKPI(uint(documentID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (daic *DocumentAssociatedItemController) DeleteDocumentKPI(c *gin.Context) {
	documentIDStr := c.Param("id")
	kpiIDStr := c.Param("kpiRegisterId")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kpiID, err := strconv.ParseUint(kpiIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = daic.service.DeleteDocumentKPI(uint(documentID), uint(kpiID))
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KPI unmapped from document successfully"})
}

func (daic *DocumentAssociatedItemController) GetAllAssociatedDocument(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	documents, err := daic.service.GetAllAssociatedDocument(uint(documentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, documents)
}
func (daic *DocumentAssociatedItemController) CreateAssociatedDocument(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	documentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}

	err = daic.service.CreateAssociatedDocument(uint(documentID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (daic *DocumentAssociatedItemController) DeleteAssociatedDocument(c *gin.Context) {
	documentIDStr := c.Param("id")
	associatedDocumentIDStr := c.Param("associatedDocumentId")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	associatedDocumentID, err := strconv.ParseUint(associatedDocumentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = daic.service.DeleteAssociatedDocument(uint(documentID), uint(associatedDocumentID))
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Associated document unmapped from document successfully"})
}

func (daic *DocumentAssociatedItemController) GetAllDocumentAsset(c *gin.Context) {

	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	assets, err := daic.service.GetAllDocumentAsset(uint(documentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, assets)
}
func (daic *DocumentAssociatedItemController) CreateDocumentAsset(c *gin.Context) {
	var input commonDTOs.AssociatedAssetMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	documentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}

	err = daic.service.CreateDocumentAsset(uint(documentID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(daic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (daic *DocumentAssociatedItemController) DeleteDocumentAsset(c *gin.Context) {
	documentIDStr := c.Param("id")
	assetIDStr := c.Param("assetId")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = daic.service.DeleteDocumentAsset(uint(documentID), uint(assetID))
	if err != nil {
		utility.SendInternalErrResp(daic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset unmapped from document successfully"})
}
