package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KRIAssociatedItemController struct {
	service services.KRIAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewKRIAssociatedItemController(service services.KRIAssociatedItemServiceInterface) *KRIAssociatedItemController {
	return &KRIAssociatedItemController{
		service: service,
	}
}

func (kaic *KRIAssociatedItemController) GetAllKRIDocuments(c *gin.Context) {

	kriIDStr := c.Param("id")

	kriID, err := strconv.ParseUint(kriIDStr, 10, 64)
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
	controls, err := kaic.service.GetAllKRIDocuments(uint(kriID), &queryParams)

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
func (kaic *KRIAssociatedItemController) CreateKRIDocument(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kriID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	err = kaic.service.CreateKRIDocument(uint(kriID), currentUserID, &input)
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

func (kaic *KRIAssociatedItemController) DeleteKRIDocument(c *gin.Context) {
	kriID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	documentID, _ := strconv.ParseUint(c.Param("documentId"), 10, 64)
	err := kaic.service.DeleteKRIDocument(uint(kriID), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from KRI successfully"})
}

func (kaic *KRIAssociatedItemController) GetAllKRIControls(c *gin.Context) {

	kriIDStr := c.Param("id")

	kriID, err := strconv.ParseUint(kriIDStr, 10, 64)
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
	controls, err := kaic.service.GetAllKRIControls(uint(kriID), &queryParams)

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
func (kaic *KRIAssociatedItemController) CreateKRIControl(c *gin.Context) {
	var input commonDTOs.AssociatedControlMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kriID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	err = kaic.service.CreateKRIControl(uint(kriID), currentUserID, &input)
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

func (kaic *KRIAssociatedItemController) DeleteKRIControl(c *gin.Context) {
	kriID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)
	err := kaic.service.DeleteKRIControl(uint(kriID), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Control unmapped from KRI successfully"})
}

func (kaic *KRIAssociatedItemController) GetAllKRIAssets(c *gin.Context) {

	kriIDStr := c.Param("id")

	kriID, err := strconv.ParseUint(kriIDStr, 10, 64)
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
	controls, err := kaic.service.GetAllKRIAssets(uint(kriID), &queryParams)

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
func (kaic *KRIAssociatedItemController) CreateKRIAsset(c *gin.Context) {
	var input commonDTOs.AssociatedAssetMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kriID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	err = kaic.service.CreateKRIAsset(uint(kriID), currentUserID, &input)
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

func (kaic *KRIAssociatedItemController) DeleteKRIAsset(c *gin.Context) {
	kriIDStr := c.Param("id")
	assetIDStr := c.Param("assetId")
	kriID, err := strconv.ParseUint(kriIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = kaic.service.DeleteKRIAsset(uint(kriID), uint(assetID))
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset unmapped from KRI successfully"})
}

func (kaic *KRIAssociatedItemController) GetAllKRIIncidents(c *gin.Context) {

	kriIDStr := c.Param("id")

	kriID, err := strconv.ParseUint(kriIDStr, 10, 64)
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
	incidents, err := kaic.service.GetAllKRIIncidents(uint(kriID), &queryParams)

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
	c.JSON(http.StatusOK, incidents)
}
func (kaic *KRIAssociatedItemController) CreateKRIIncident(c *gin.Context) {
	var input commonDTOs.AssociatedIncidentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kriID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	err = kaic.service.CreateKRIIncident(uint(kriID), currentUserID, &input)
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

func (kaic *KRIAssociatedItemController) DeleteKRIIncident(c *gin.Context) {
	kriIDStr := c.Param("id")
	incidentIDStr := c.Param("incidentId")
	kriID, err := strconv.ParseUint(kriIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = kaic.service.DeleteKRIIncident(uint(kriID), uint(incidentID))
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident unmapped from KRI successfully"})
}

func (kaic *KRIAssociatedItemController) GetAllKRIRisks(c *gin.Context) {

	kriIDStr := c.Param("id")

	kriID, err := strconv.ParseUint(kriIDStr, 10, 64)
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
	risks, err := kaic.service.GetAllKRIRisks(uint(kriID), &queryParams)

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
func (kaic *KRIAssociatedItemController) CreateKRIRisk(c *gin.Context) {
	var input commonDTOs.AssociatedRiskMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kriID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	err = kaic.service.CreateKRIRisk(uint(kriID), currentUserID, &input)
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

func (kaic *KRIAssociatedItemController) DeleteKRIRisk(c *gin.Context) {
	kriIDStr := c.Param("id")
	riskIDStr := c.Param("riskId")
	kriID, err := strconv.ParseUint(kriIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = kaic.service.DeleteKRIRisk(uint(kriID), uint(riskID))
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Risk unmapped from KRI successfully"})
}
func (kaic *KRIAssociatedItemController) GetAllKRIKPI(c *gin.Context) {

	kriIDStr := c.Param("id")

	kriID, err := strconv.ParseUint(kriIDStr, 10, 64)
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
	kpis, err := kaic.service.GetAllKRIKPI(uint(kriID), &queryParams)

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
	c.JSON(http.StatusOK, kpis)
}
func (kaic *KRIAssociatedItemController) CreateKRIKPI(c *gin.Context) {
	var input commonDTOs.AssociatedKPIMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kriID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	err = kaic.service.CreateKRIKPI(uint(kriID), currentUserID, &input)
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

func (kaic *KRIAssociatedItemController) DeleteKRIKPI(c *gin.Context) {
	kriIDStr := c.Param("id")
	kpiIDStr := c.Param("kpiRegisterId")
	kriID, err := strconv.ParseUint(kriIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kpiID, err := strconv.ParseUint(kpiIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = kaic.service.DeleteKRIKPI(uint(kriID), uint(kpiID))
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KPI unmapped from KRI successfully"})
}

func (kaic *KRIAssociatedItemController) GetAllKRIProcesses(c *gin.Context) {

	kriIDStr := c.Param("id")

	kriID, err := strconv.ParseUint(kriIDStr, 10, 64)
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
	proceeses, err := kaic.service.GetAllKRIProcesses(uint(kriID), &queryParams)

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
	c.JSON(http.StatusOK, proceeses)
}
func (kaic *KRIAssociatedItemController) CreateKRIProcess(c *gin.Context) {
	var input commonDTOs.AssociatedProcessMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kriID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	err = kaic.service.CreateKRIProcess(uint(kriID), currentUserID, &input)
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

func (kaic *KRIAssociatedItemController) DeleteKRIProcess(c *gin.Context) {
	kriIDStr := c.Param("id")
	processIDStr := c.Param("processId")
	kriID, err := strconv.ParseUint(kriIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	processID, err := strconv.ParseUint(processIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = kaic.service.DeleteKRIProcess(uint(kriID), uint(processID))
	if err != nil {
		utility.SendInternalErrResp(kaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Process unmapped from KRI successfully"})
}

func (kaic *KRIAssociatedItemController) GetAllClauses(c *gin.Context) {
	projectIDStr := c.Param("id")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
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
	clauses, err := kaic.service.GetAllClauses(uint(projectID), &queryParams)

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
func (kaic *KRIAssociatedItemController) UpdateClauses(c *gin.Context) {
	var clauseFormDTO commonDTOs.AssociatedClauseFormDTO

	projectID, err := strconv.Atoi(c.Param("id"))
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

	err = kaic.service.UpdateClauses(uint(projectID), currentUserID, &clauseFormDTO)
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

	c.JSON(http.StatusCreated, gin.H{"message": "Clauses against the kri updated successfully"})
}
