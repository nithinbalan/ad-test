package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KPIAssociatedItemController struct {
	service services.KPIAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewKPIAssociatedItemController(service services.KPIAssociatedItemServiceInterface) *KPIAssociatedItemController {
	return &KPIAssociatedItemController{
		service: service,
	}
}

func (kaic *KPIAssociatedItemController) GetAllKPIControls(c *gin.Context) {

	kpiIdStr := c.Param("kpiRegisterId")

	kpiId, err := strconv.ParseUint(kpiIdStr, 10, 64)
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
		utility.SendInternalErrResp(kaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := strategyDTOs.KPIAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	controls, err := kaic.service.GetAllControls(uint(kpiId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}

func (kaic *KPIAssociatedItemController) CreateKPIControls(c *gin.Context) {
	var input strategyDTOs.KPIControlMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kpiId, err := strconv.ParseUint(c.Param("kpiRegisterId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	err = kaic.service.CreateKPIControls(uint(kpiId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (kaic *KPIAssociatedItemController) DeleteKPIControl(c *gin.Context) {
	kpiId, _ := strconv.ParseUint(c.Param("kpiRegisterId"), 10, 64)
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)

	err := kaic.service.DeleteKPIControl(uint(kpiId), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Control unmapped from KPI successfully"})
}

func (kaic *KPIAssociatedItemController) GetAllKPIProcesses(c *gin.Context) {

	kpiIdStr := c.Param("kpiRegisterId")

	kpiId, err := strconv.ParseUint(kpiIdStr, 10, 64)
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
		utility.SendInternalErrResp(kaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := strategyDTOs.KPIAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	processes, err := kaic.service.GetAllProcesses(uint(kpiId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, processes)
}

func (kaic *KPIAssociatedItemController) CreateKPIProcesses(c *gin.Context) {
	var input strategyDTOs.KPIProcessMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kpiId, err := strconv.ParseUint(c.Param("kpiRegisterId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	err = kaic.service.CreateKPIProcesses(uint(kpiId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (kaic *KPIAssociatedItemController) DeleteKPIProcess(c *gin.Context) {
	kpiId, _ := strconv.ParseUint(c.Param("kpiRegisterId"), 10, 64)
	processID, _ := strconv.ParseUint(c.Param("processId"), 10, 64)

	err := kaic.service.DeleteKPIProcess(uint(kpiId), uint(processID))
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Process unmapped from KPI successfully"})
}

func (kaic *KPIAssociatedItemController) GetAllKPICompliances(c *gin.Context) {

	kpiIdStr := c.Param("kpiRegisterId")

	kpiId, err := strconv.ParseUint(kpiIdStr, 10, 64)
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
		utility.SendInternalErrResp(kaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := strategyDTOs.KPIAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	compliances, err := kaic.service.GetAllCompliances(uint(kpiId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, compliances)
}

func (kaic *KPIAssociatedItemController) CreateKPICompliances(c *gin.Context) {
	var input strategyDTOs.KPIComplianceMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kpiId, err := strconv.ParseUint(c.Param("kpiRegisterId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	err = kaic.service.CreateKPICompliances(uint(kpiId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (kaic *KPIAssociatedItemController) DeleteKPICompliance(c *gin.Context) {
	kpiId, _ := strconv.ParseUint(c.Param("kpiRegisterId"), 10, 64)
	complianceID, _ := strconv.ParseUint(c.Param("complianceId"), 10, 64)

	err := kaic.service.DeleteKPICompliance(uint(kpiId), uint(complianceID))
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compliance unmapped from KPI successfully"})
}

func (kaic *KPIAssociatedItemController) GetAllKPIAssets(c *gin.Context) {

	kpiIdStr := c.Param("kpiRegisterId")

	kpiId, err := strconv.ParseUint(kpiIdStr, 10, 64)
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
		utility.SendInternalErrResp(kaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := strategyDTOs.KPIAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	assets, err := kaic.service.GetAllAssets(uint(kpiId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, assets)
}

func (kaic *KPIAssociatedItemController) CreateKPIAssets(c *gin.Context) {
	var input strategyDTOs.KPIAssetMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kpiId, err := strconv.ParseUint(c.Param("kpiRegisterId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	err = kaic.service.CreateKPIAssets(uint(kpiId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (kaic *KPIAssociatedItemController) DeleteKPIAsset(c *gin.Context) {
	kpiId, _ := strconv.ParseUint(c.Param("kpiRegisterId"), 10, 64)
	assetID, _ := strconv.ParseUint(c.Param("assetId"), 10, 64)

	err := kaic.service.DeleteKPIAsset(uint(kpiId), uint(assetID))
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset unmapped from KPI successfully"})
}

func (kaic *KPIAssociatedItemController) GetAllKPIRisks(c *gin.Context) {

	kpiIdStr := c.Param("kpiRegisterId")

	kpiId, err := strconv.ParseUint(kpiIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, "Invalid page_size parameter")
		return
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
		Filters:  filters,
	}
	risks, err := kaic.service.GetAllRisks(uint(kpiId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (kaic *KPIAssociatedItemController) CreateKPIRisks(c *gin.Context) {
	var input strategyDTOs.KPIRiskMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kpiId, err := strconv.ParseUint(c.Param("kpiRegisterId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	err = kaic.service.CreateKPIRisks(uint(kpiId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (kaic *KPIAssociatedItemController) DeleteKPIRisk(c *gin.Context) {
	kpiId, _ := strconv.ParseUint(c.Param("kpiRegisterId"), 10, 64)
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)

	err := kaic.service.DeleteKPIRisk(uint(kpiId), uint(riskID))
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Risk unmapped from KPI successfully"})
}

func (kaic *KPIAssociatedItemController) GetAllKPIThirdParties(c *gin.Context) {

	kpiIdStr := c.Param("kpiRegisterId")

	kpiId, err := strconv.ParseUint(kpiIdStr, 10, 64)
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
		utility.SendInternalErrResp(kaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := strategyDTOs.KPIAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	thirdParties, err := kaic.service.GetAllThirdParties(uint(kpiId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, thirdParties)
}

func (kaic *KPIAssociatedItemController) CreateKPIThirdParties(c *gin.Context) {
	var input strategyDTOs.KPIThirdPartyMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kpiId, err := strconv.ParseUint(c.Param("kpiRegisterId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	err = kaic.service.CreateKPIThirdParties(uint(kpiId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (kaic *KPIAssociatedItemController) DeleteKPIThirdParty(c *gin.Context) {
	kpiId, _ := strconv.ParseUint(c.Param("kpiRegisterId"), 10, 64)
	thirdPartyID, _ := strconv.ParseUint(c.Param("thirdPartyRegisterId"), 10, 64)

	err := kaic.service.DeleteKPIThirdParty(uint(kpiId), uint(thirdPartyID))
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Third Party Register unmapped from KPI successfully"})
}

func (kaic *KPIAssociatedItemController) GetAllKPIPolicies(c *gin.Context) {

	kpiIdStr := c.Param("kpiRegisterId")

	kpiId, err := strconv.ParseUint(kpiIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	documents, err := kaic.service.GetAllPolicies(uint(kpiId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, documents)
}

func (kaic *KPIAssociatedItemController) CreateKPIPolicies(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kpiId, err := strconv.ParseUint(c.Param("kpiRegisterId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	err = kaic.service.CreateKPIPolicies(uint(kpiId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (kaic *KPIAssociatedItemController) DeleteKPIPolicy(c *gin.Context) {
	kpiId, _ := strconv.ParseUint(c.Param("kpiRegisterId"), 10, 64)
	documentID, _ := strconv.ParseUint(c.Param("documentId"), 10, 64)

	err := kaic.service.DeleteKPIPolicy(uint(kpiId), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from KPI successfully"})
}

func (kaic *KPIAssociatedItemController) GetAllClauses(c *gin.Context) {
	kpiIDStr := c.Param("kpiRegisterId")

	riskID, err := strconv.ParseUint(kpiIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	clauses, err := kaic.service.GetAllClauses(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, clauses)
}
func (kaic *KPIAssociatedItemController) UpdateClauses(c *gin.Context) {
	var clauseFormDTO commonDTOs.AssociatedClauseFormDTO

	kpiID, err := strconv.Atoi(c.Param("kpiRegisterId"))
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
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	err = kaic.service.UpdateClauses(uint(kpiID), currentUserID, &clauseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kaic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Clauses against the KPI updated successfully"})
}
