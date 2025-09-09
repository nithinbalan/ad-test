package controllers

import (
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditFindingAssociatedItemController struct {
	service services.AuditFindingAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewAuditFindingAssociatedItemController(service services.AuditFindingAssociatedItemServiceInterface) *AuditFindingAssociatedItemController {
	return &AuditFindingAssociatedItemController{
		service: service,
	}
}

func (afaic *AuditFindingAssociatedItemController) GetAllAuditFindingDocuments(c *gin.Context) {

	findingIdStr := c.Param("findingId")

	findingId, err := strconv.ParseUint(findingIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := afaic.service.GetAllAuditFindingDocuments(uint(findingId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}
func (afaic *AuditFindingAssociatedItemController) CreateAuditFindingDocument(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	findingId, err := strconv.ParseUint(c.Param("findingId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	err = afaic.service.CreateAuditFindingDocument(uint(findingId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (afaic *AuditFindingAssociatedItemController) DeleteAuditFindingDocument(c *gin.Context) {
	findingId, _ := strconv.ParseUint(c.Param("findingId"), 10, 64)
	documentID, _ := strconv.ParseUint(c.Param("documentId"), 10, 64)
	err := afaic.service.DeleteAuditFindingDocument(uint(findingId), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from AuditFinding successfully"})
}

func (afaic *AuditFindingAssociatedItemController) GetAllAuditFindingControls(c *gin.Context) {

	findingIdStr := c.Param("findingId")

	findingId, err := strconv.ParseUint(findingIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := afaic.service.GetAllAuditFindingControls(uint(findingId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}
func (afaic *AuditFindingAssociatedItemController) CreateAuditFindingControl(c *gin.Context) {
	var input commonDTOs.AssociatedControlMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	findingId, err := strconv.ParseUint(c.Param("findingId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	err = afaic.service.CreateAuditFindingControl(uint(findingId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (afaic *AuditFindingAssociatedItemController) DeleteAuditFindingControl(c *gin.Context) {
	findingId, _ := strconv.ParseUint(c.Param("findingId"), 10, 64)
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)
	err := afaic.service.DeleteAuditFindingControl(uint(findingId), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Control unmapped from AuditFinding successfully"})
}

func (afaic *AuditFindingAssociatedItemController) GetAllAuditFindingProcesses(c *gin.Context) {

	findingIdStr := c.Param("findingId")

	findingId, err := strconv.ParseUint(findingIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := afaic.service.GetAllAuditFindingProcesses(uint(findingId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}
func (afaic *AuditFindingAssociatedItemController) CreateAuditFindingProcess(c *gin.Context) {
	var input commonDTOs.AssociatedProcessMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	findingId, err := strconv.ParseUint(c.Param("findingId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	err = afaic.service.CreateAuditFindingProcess(uint(findingId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (afaic *AuditFindingAssociatedItemController) DeleteAuditFindingProcess(c *gin.Context) {
	findingId, _ := strconv.ParseUint(c.Param("findingId"), 10, 64)
	processID, _ := strconv.ParseUint(c.Param("processId"), 10, 64)
	err := afaic.service.DeleteAuditFindingProcess(uint(findingId), uint(processID))
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Process unmapped from AuditFinding successfully"})
}

func (afaic *AuditFindingAssociatedItemController) GetAllAuditFindingStrategicObjectives(c *gin.Context) {

	findingIdStr := c.Param("findingId")

	findingId, err := strconv.ParseUint(findingIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	objectives, err := afaic.service.GetAllAuditFindingStrategicObjectives(uint(findingId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, objectives)
}
func (afaic *AuditFindingAssociatedItemController) CreateAuditFindingStrategicObjective(c *gin.Context) {
	var input commonDTOs.AssociatedStrategicObjectiveMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	findingId, err := strconv.ParseUint(c.Param("findingId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	err = afaic.service.CreateAuditFindingStrategicObjective(uint(findingId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (afaic *AuditFindingAssociatedItemController) DeleteAuditFindingStrategicObjective(c *gin.Context) {
	findingId, _ := strconv.ParseUint(c.Param("findingId"), 10, 64)
	strategicObjectiveID, _ := strconv.ParseUint(c.Param("strategicObjectiveId"), 10, 64)
	err := afaic.service.DeleteAuditFindingStrategicObjective(uint(findingId), uint(strategicObjectiveID))
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Strategic Objective unmapped from AuditFinding successfully"})
}

func (afaic *AuditFindingAssociatedItemController) GetAllAuditFindingAssets(c *gin.Context) {

	findingIdStr := c.Param("findingId")

	findingId, err := strconv.ParseUint(findingIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	assets, err := afaic.service.GetAllAuditFindingAssets(uint(findingId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, assets)
}
func (afaic *AuditFindingAssociatedItemController) CreateAuditFindingAsset(c *gin.Context) {
	var input commonDTOs.AssociatedAssetMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	findingId, err := strconv.ParseUint(c.Param("findingId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	err = afaic.service.CreateAuditFindingAsset(uint(findingId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (afaic *AuditFindingAssociatedItemController) DeleteAuditFindingAsset(c *gin.Context) {
	assetIDStr := c.Param("assetId")
	findingIdIDStr := c.Param("findingId")
	findingID, err := strconv.ParseUint(findingIdIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = afaic.service.DeleteAuditFindingAsset(uint(findingID), uint(assetID))
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset unmapped from AuditFinding successfully"})
}

func (afaic *AuditFindingAssociatedItemController) GetAllAuditFindingIncidents(c *gin.Context) {

	findingIdStr := c.Param("findingId")

	findingId, err := strconv.ParseUint(findingIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	incidents, err := afaic.service.GetAllAuditFindingIncidents(uint(findingId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidents)
}
func (afaic *AuditFindingAssociatedItemController) CreateAuditFindingIncident(c *gin.Context) {
	var input commonDTOs.AssociatedIncidentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	findingId, err := strconv.ParseUint(c.Param("findingId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	err = afaic.service.CreateAuditFindingIncident(uint(findingId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (afaic *AuditFindingAssociatedItemController) DeleteAuditFindingIncident(c *gin.Context) {
	incidentIDStr := c.Param("incidentId")
	findingIdIDStr := c.Param("findingId")
	findingID, err := strconv.ParseUint(findingIdIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = afaic.service.DeleteAuditFindingIncident(uint(findingID), uint(incidentID))
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident unmapped from AuditFinding successfully"})
}

func (afaic *AuditFindingAssociatedItemController) GetAllAuditFindingRisks(c *gin.Context) {

	findingIdStr := c.Param("findingId")

	findingId, err := strconv.ParseUint(findingIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := afaic.service.GetAllAuditFindingRisks(uint(findingId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}
func (afaic *AuditFindingAssociatedItemController) CreateAuditFindingRisk(c *gin.Context) {
	var input commonDTOs.AssociatedRiskMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	findingId, err := strconv.ParseUint(c.Param("findingId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	err = afaic.service.CreateAuditFindingRisk(uint(findingId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (afaic *AuditFindingAssociatedItemController) DeleteAuditFindingRisk(c *gin.Context) {
	riskIDStr := c.Param("riskId")
	findingIdIDStr := c.Param("findingId")
	findingID, err := strconv.ParseUint(findingIdIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = afaic.service.DeleteAuditFindingRisk(uint(findingID), uint(riskID))
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Risk unmapped from AuditFinding successfully"})
}

func (afaic *AuditFindingAssociatedItemController) GetAllAuditFindingKPI(c *gin.Context) {

	findingIdStr := c.Param("findingId")

	findingId, err := strconv.ParseUint(findingIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	kpis, err := afaic.service.GetAllAuditFindingKPI(uint(findingId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, kpis)
}
func (afaic *AuditFindingAssociatedItemController) CreateAuditFindingKPI(c *gin.Context) {
	var input commonDTOs.AssociatedKPIMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	findingId, err := strconv.ParseUint(c.Param("findingId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	err = afaic.service.CreateAuditFindingKPI(uint(findingId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (afaic *AuditFindingAssociatedItemController) DeleteAuditFindingKPI(c *gin.Context) {
	kpiIDStr := c.Param("kpiRegisterId")
	findingIdIDStr := c.Param("findingId")
	findingID, err := strconv.ParseUint(findingIdIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kpiID, err := strconv.ParseUint(kpiIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = afaic.service.DeleteAuditFindingKPI(uint(findingID), uint(kpiID))
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KPI unmapped from AuditFinding successfully"})
}

func (afaic *AuditFindingAssociatedItemController) GetAllClauses(c *gin.Context) {
	findingIDStr := c.Param("findingId")

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	clauses, err := afaic.service.GetAllClauses(uint(findingID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, clauses)
}
func (afaic *AuditFindingAssociatedItemController) UpdateClauses(c *gin.Context) {
	var clauseFormDTO commonDTOs.AssociatedClauseFormDTO

	findingID, err := strconv.Atoi(c.Param("findingId"))
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
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	err = afaic.service.UpdateClauses(uint(findingID), currentUserID, &clauseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Clauses against the audit finding updated successfully"})
}
