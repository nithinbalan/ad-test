package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	riskServices "jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetRiskAssociatedItemController struct {
	service riskServices.RiskAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewAssetRiskAssociatedItemController(service riskServices.RiskAssociatedItemServiceInterface) *AssetRiskAssociatedItemController {
	return &AssetRiskAssociatedItemController{
		service: service,
	}
}

func (araic *AssetRiskAssociatedItemController) GetAllAssetRiskControls(c *gin.Context) {

	riskIDStr := c.Param("id")

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
		utility.SendInternalErrResp(araic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := araic.service.GetAllControls(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (araic *AssetRiskAssociatedItemController) GetAllAssetRiskProcesses(c *gin.Context) {

	riskIDStr := c.Param("id")

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
		utility.SendInternalErrResp(araic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := araic.service.GetAllProcesses(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (araic *AssetRiskAssociatedItemController) GetAllAssetRiskStrategicObjectives(c *gin.Context) {

	riskIDStr := c.Param("id")

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
		utility.SendInternalErrResp(araic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := araic.service.GetAllStrategicObjectives(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (araic *AssetRiskAssociatedItemController) GetAllAssetRiskIncidents(c *gin.Context) {

	riskIDStr := c.Param("id")

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
		utility.SendInternalErrResp(araic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, "Invalid page_size parameter")
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
	risks, err := araic.service.GetAllIncidents(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (araic *AssetRiskAssociatedItemController) GetAllAssetRiskDocuments(c *gin.Context) {
	riskIDStr := c.Param("id")

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
		utility.SendInternalErrResp(araic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, "Invalid page_size parameter")
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
	documents, err := araic.service.GetAllDocuments(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, documents)
}

func (araic *AssetRiskAssociatedItemController) GetAllAssetRiskAssets(c *gin.Context) {

	riskIDStr := c.Param("id")

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
		utility.SendInternalErrResp(araic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := araic.service.GetAllAssets(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (araic *AssetRiskAssociatedItemController) GetAllAssetRiskKPI(c *gin.Context) {

	riskIDStr := c.Param("id")

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
		utility.SendInternalErrResp(araic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := araic.service.GetAllRiskKPI(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (araic *AssetRiskAssociatedItemController) GetAllAssetRiskClauses(c *gin.Context) {
	riskIDStr := c.Param("id")

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
		utility.SendInternalErrResp(araic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	clauses, err := araic.service.GetAllClauses(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, clauses)
}

func (araic *AssetRiskAssociatedItemController) UpdateAssetRiskClauses(c *gin.Context) {
	var clauseFormDTO commonDTOs.AssociatedClauseFormDTO

	riskID, err := strconv.Atoi(c.Param("id"))
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
		utility.SendInternalErrResp(araic.db, c, err.Error())
		return
	}

	err = araic.service.UpdateClauses(uint(riskID), currentUserID, &clauseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Clauses against the risk updated successfully"})
}

func (araic *AssetRiskAssociatedItemController) CreateAssetRiskStrategicObjective(c *gin.Context) {
	var input commonDTOs.AssociatedStrategicObjectiveMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, err.Error())
		return
	}

	err = araic.service.CreatRiskStrategicObjective(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (araic *AssetRiskAssociatedItemController) DeleteAssetRiskStrategicObjective(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	strategicObjectiveID, _ := strconv.ParseUint(c.Param("strategicObjectiveId"), 10, 64)
	err := araic.service.DeletetRiskStrategicObjective(uint(riskID), uint(strategicObjectiveID))
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Strategic objective unmapped from risk successfully"})
}

func (araic *AssetRiskAssociatedItemController) CreateAssetRiskKPI(c *gin.Context) {
	var input commonDTOs.AssociatedKPIMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, err.Error())
		return
	}

	err = araic.service.CreatRiskKPI(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (araic *AssetRiskAssociatedItemController) DeleteAssetRiskKPI(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	kpiID, _ := strconv.ParseUint(c.Param("kpiId"), 10, 64)
	err := araic.service.DeletetRiskKPI(uint(riskID), uint(kpiID))
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KPI unmapped from risk successfully"})
}

func (araic *AssetRiskAssociatedItemController) CreateAssetRiskProcess(c *gin.Context) {
	var input commonDTOs.AssociatedProcessMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, err.Error())
		return
	}

	err = araic.service.CreatRiskProcess(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (araic *AssetRiskAssociatedItemController) DeleteAssetRiskProcess(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	processID, _ := strconv.ParseUint(c.Param("processId"), 10, 64)
	err := araic.service.DeletetRiskProcess(uint(riskID), uint(processID))
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "process unmapped from risk successfully"})
}

func (araic *AssetRiskAssociatedItemController) CreateAssetRiskDocument(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, err.Error())
		return
	}

	err = araic.service.CreatRiskDocument(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (araic *AssetRiskAssociatedItemController) DeleteAssetRiskDocument(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	documentID, _ := strconv.ParseUint(c.Param("documentId"), 10, 64)
	err := araic.service.DeletetRiskDocument(uint(riskID), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from risk successfully"})
}

func (araic *AssetRiskAssociatedItemController) CreateAssetRiskIncident(c *gin.Context) {
	var input commonDTOs.AssociatedIncidentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, err.Error())
		return
	}

	err = araic.service.CreatRiskIncident(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (araic *AssetRiskAssociatedItemController) DeleteAssetRiskIncident(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	incidentID, _ := strconv.ParseUint(c.Param("incidentId"), 10, 64)
	err := araic.service.DeletetRiskIncident(uint(riskID), uint(incidentID))
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident unmapped from risk successfully"})
}

func (araic *AssetRiskAssociatedItemController) CreateAssetRiskControl(c *gin.Context) {
	var input commonDTOs.AssociatedControlMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, err.Error())
		return
	}

	err = araic.service.CreatRiskControl(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (araic *AssetRiskAssociatedItemController) DeleteAssetRiskControl(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)
	err := araic.service.DeletetRiskControl(uint(riskID), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Control unmapped from risk successfully"})
}

func (araic *AssetRiskAssociatedItemController) CreateAssetRiskAsset(c *gin.Context) {
	var input commonDTOs.AssociatedAssetMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	riskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, err.Error())
		return
	}

	err = araic.service.CreatRiskAsset(uint(riskID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(araic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (araic *AssetRiskAssociatedItemController) DeleteAssetRiskAsset(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	assetID, _ := strconv.ParseUint(c.Param("assetId"), 10, 64)
	err := araic.service.DeletetRiskAsset(uint(riskID), uint(assetID))
	if err != nil {
		utility.SendInternalErrResp(araic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset unmapped from risk successfully"})
}
