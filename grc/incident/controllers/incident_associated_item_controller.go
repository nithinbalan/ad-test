package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentAssociatedItemController struct {
	service services.IncidentAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewIncidentAssociatedItemController(service services.IncidentAssociatedItemServiceInterface) *IncidentAssociatedItemController {
	return &IncidentAssociatedItemController{
		service: service,
	}
}

func (iaic *IncidentAssociatedItemController) GetAllIncidentControls(c *gin.Context) {

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
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
		utility.SendInternalErrResp(iaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := incidentDTOs.IncidentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}

	incidents, err := iaic.service.GetAllControls(uint(incidentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidents)
}

func (iaic *IncidentAssociatedItemController) GetAllIncidentComplianceRequirements(c *gin.Context) {

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
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
		utility.SendInternalErrResp(iaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := incidentDTOs.IncidentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	incidents, err := iaic.service.GetAllComplianceRequirements(uint(incidentID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func (iaic *IncidentAssociatedItemController) GetAllIncidentRisks(c *gin.Context) {

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	incidents, err := iaic.service.GetAllRisks(uint(incidentID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func (iaic *IncidentAssociatedItemController) GetAllIncidentProcesses(c *gin.Context) {

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
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
		utility.SendInternalErrResp(iaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := incidentDTOs.IncidentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	incidents, err := iaic.service.GetAllProcesses(uint(incidentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidents)
}

func (iaic *IncidentAssociatedItemController) GetAllIncidentSLAs(c *gin.Context) {

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
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
		utility.SendInternalErrResp(iaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := incidentDTOs.IncidentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	incidents, err := iaic.service.GetAllSLAs(uint(incidentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidents)
}

func (ic *IncidentAssociatedItemController) GetAllIncidentDocuments(c *gin.Context) {
	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(ic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(ic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	documents, err := ic.service.GetAllDocuments(uint(incidentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, documents)
}

func (iaic *IncidentAssociatedItemController) GetAllIncidentAssets(c *gin.Context) {

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
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
		utility.SendInternalErrResp(iaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := incidentDTOs.IncidentAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	incidents, err := iaic.service.GetAllAssets(uint(incidentID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func (iaic *IncidentAssociatedItemController) GetAllIncidentClauses(c *gin.Context) {

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	clauses, err := iaic.service.GetAllClauses(uint(incidentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, clauses)
}

func (iaic *IncidentAssociatedItemController) UpdateClauses(c *gin.Context) {
	var clauseFormDTO commonDTOs.AssociatedClauseFormDTO

	incidentID, err := strconv.Atoi(c.Param("incidentID"))
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
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}

	err = iaic.service.UpdateClauses(uint(incidentID), currentUserID, &clauseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iaic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Clauses against incident updated successfully"})
}

func (iaic *IncidentAssociatedItemController) CreateIncidentRisk(c *gin.Context) {
	var input commonDTOs.AssociatedRiskMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	incidentID, err := strconv.ParseUint(c.Param("incidentID"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}

	err = iaic.service.CreateIncidentRisk(uint(incidentID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (iaic *IncidentAssociatedItemController) DeleteIncidentRisk(c *gin.Context) {
	incidentID, _ := strconv.ParseUint(c.Param("incidentID"), 10, 64)
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	err := iaic.service.DeleteIncidentRisk(uint(incidentID), uint(riskID))
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Risk unmapped from incident successfully"})
}

func (iaic *IncidentAssociatedItemController) CreateIncidentDocument(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	incidentID, err := strconv.ParseUint(c.Param("incidentID"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}

	err = iaic.service.CreateIncidentDocument(uint(incidentID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (iaic *IncidentAssociatedItemController) DeleteIncidentDocument(c *gin.Context) {
	incidentID, _ := strconv.ParseUint(c.Param("incidentID"), 10, 64)
	documentID, _ := strconv.ParseUint(c.Param("documentId"), 10, 64)
	err := iaic.service.DeleteIncidentDocument(uint(incidentID), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from incident successfully"})
}
