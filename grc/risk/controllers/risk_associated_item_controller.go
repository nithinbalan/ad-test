package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskAssociatedItemController struct {
	service services.RiskAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewRiskAssociatedItemController(service services.RiskAssociatedItemServiceInterface) *RiskAssociatedItemController {
	return &RiskAssociatedItemController{
		service: service,
	}
}

func (raic *RiskAssociatedItemController) GetAllRiskControls(c *gin.Context) {

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := raic.service.GetAllControls(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(raic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (raic *RiskAssociatedItemController) GetAllRiskComplianceRequirements(c *gin.Context) {

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
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
		utility.SendInternalErrResp(raic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := riskDTOs.RiskAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	risks, err := raic.service.GetAllComplianceRequirements(uint(riskID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(raic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, risks)
}

func (raic *RiskAssociatedItemController) GetAllRiskProcesses(c *gin.Context) {

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := raic.service.GetAllProcesses(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(raic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (raic *RiskAssociatedItemController) GetAllRiskStrategicObjectives(c *gin.Context) {

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := raic.service.GetAllStrategicObjectives(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(raic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (raic *RiskAssociatedItemController) GetAllRiskSLAs(c *gin.Context) {

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
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
		utility.SendInternalErrResp(raic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := riskDTOs.RiskAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	risks, err := raic.service.GetAllSLAs(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(raic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (raic *RiskAssociatedItemController) GetAllRiskIncidents(c *gin.Context) {

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	risks, err := raic.service.GetAllIncidents(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(raic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (rc *RiskAssociatedItemController) GetAllRiskDocuments(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	documents, err := rc.service.GetAllDocuments(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, documents)
}

func (raic *RiskAssociatedItemController) GetAllRiskAssets(c *gin.Context) {

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := raic.service.GetAllAssets(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(raic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (raic *RiskAssociatedItemController) GetAllRiskThreats(c *gin.Context) {

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
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
		utility.SendInternalErrResp(raic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := riskDTOs.RiskAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	threats, err := raic.service.GetAllThreats(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(raic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, threats)
}

func (raic *RiskAssociatedItemController) GetAllRiskKPI(c *gin.Context) {

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := raic.service.GetAllRiskKPI(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(raic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (rc *RiskAssociatedItemController) GetAllClauses(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	clauses, err := rc.service.GetAllClauses(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, clauses)
}
func (rc *RiskAssociatedItemController) UpdateClauses(c *gin.Context) {
	var clauseFormDTO commonDTOs.AssociatedClauseFormDTO

	riskID, err := strconv.Atoi(c.Param("riskId"))
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
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	err = rc.service.UpdateClauses(uint(riskID), currentUserID, &clauseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Clauses against the risk updated successfully"})
}

func (rc *RiskAssociatedItemController) CreatRiskStrategicObjective(c *gin.Context) {
	var input commonDTOs.AssociatedStrategicObjectiveMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("riskId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	err = rc.service.CreatRiskStrategicObjective(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (rc *RiskAssociatedItemController) DeletetRiskStrategicObjective(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	strategicObjectiveID, _ := strconv.ParseUint(c.Param("strategicObjectiveId"), 10, 64)
	err := rc.service.DeletetRiskStrategicObjective(uint(riskID), uint(strategicObjectiveID))
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Strategic objective unmapped from risk successfully"})
}

func (rc *RiskAssociatedItemController) CreatRiskKPI(c *gin.Context) {
	var input commonDTOs.AssociatedKPIMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("riskId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	err = rc.service.CreatRiskKPI(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (rc *RiskAssociatedItemController) DeletetRiskKPI(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	kpiID, _ := strconv.ParseUint(c.Param("kpiId"), 10, 64)
	err := rc.service.DeletetRiskKPI(uint(riskID), uint(kpiID))
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KPI unmapped from risk successfully"})
}
func (rc *RiskAssociatedItemController) CreatRiskProcess(c *gin.Context) {
	var input commonDTOs.AssociatedProcessMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("riskId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	err = rc.service.CreatRiskProcess(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (rc *RiskAssociatedItemController) DeletetRiskProcess(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	processID, _ := strconv.ParseUint(c.Param("processId"), 10, 64)
	err := rc.service.DeletetRiskProcess(uint(riskID), uint(processID))
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "process unmapped from risk successfully"})
}

func (rc *RiskAssociatedItemController) CreatRiskDocument(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("riskId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	err = rc.service.CreatRiskDocument(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (rc *RiskAssociatedItemController) DeletetRiskDocument(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	documentID, _ := strconv.ParseUint(c.Param("documentId"), 10, 64)
	err := rc.service.DeletetRiskDocument(uint(riskID), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from risk successfully"})
}

func (rc *RiskAssociatedItemController) CreatRiskIncident(c *gin.Context) {
	var input commonDTOs.AssociatedIncidentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("riskId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	err = rc.service.CreatRiskIncident(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (rc *RiskAssociatedItemController) DeletetRiskIncident(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	incidentID, _ := strconv.ParseUint(c.Param("incidentId"), 10, 64)
	err := rc.service.DeletetRiskIncident(uint(riskID), uint(incidentID))
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident unmapped from risk successfully"})
}

func (rc *RiskAssociatedItemController) CreatRiskControl(c *gin.Context) {
	var input commonDTOs.AssociatedControlMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("riskId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	err = rc.service.CreatRiskControl(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (rc *RiskAssociatedItemController) DeletetRiskControl(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)
	err := rc.service.DeletetRiskControl(uint(riskID), uint(controlID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Control unmapped from risk successfully"})
}

func (rc *RiskAssociatedItemController) CreatRiskAsset(c *gin.Context) {
	var input commonDTOs.AssociatedAssetMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("riskId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	err = rc.service.CreatRiskAsset(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (rc *RiskAssociatedItemController) DeletetRiskAsset(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	assetID, _ := strconv.ParseUint(c.Param("assetId"), 10, 64)
	err := rc.service.DeletetRiskAsset(uint(riskID), uint(assetID))
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset unmapped from risk successfully"})
}
func (raic *RiskAssociatedItemController) GetAllRiskBCP(c *gin.Context) {

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := raic.service.GetAllRiskBCP(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(raic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}
func (rc *RiskAssociatedItemController) CreatRiskBCP(c *gin.Context) {
	var input commonDTOs.AssociatedBCPMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("riskId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	err = rc.service.CreatRiskBCP(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (rc *RiskAssociatedItemController) DeletetRiskBCP(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	bcpID, _ := strconv.ParseUint(c.Param("bcpId"), 10, 64)
	err := rc.service.DeletetRiskBCP(uint(riskID), uint(bcpID))
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "BCP unmapped from risk successfully"})
}

func (raic *RiskAssociatedItemController) GetAllRiskProject(c *gin.Context) {

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(raic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := raic.service.GetAllRiskProject(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(raic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}
func (rc *RiskAssociatedItemController) CreatRiskProject(c *gin.Context) {
	var input commonDTOs.AssociatedProjectMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("riskId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	err = rc.service.CreatRiskProject(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (rc *RiskAssociatedItemController) DeletetRiskProject(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	projectID, _ := strconv.ParseUint(c.Param("projectId"), 10, 64)
	err := rc.service.DeletetRiskProject(uint(riskID), uint(projectID))
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project unmapped from risk successfully"})
}
