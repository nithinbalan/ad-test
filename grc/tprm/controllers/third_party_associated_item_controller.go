package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/tprm/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThirdPartyAssociatedItemController struct {
	service services.ThirdPartyAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewThirdPartyAssociatedItemController(service services.ThirdPartyAssociatedItemServiceInterface) *ThirdPartyAssociatedItemController {
	return &ThirdPartyAssociatedItemController{
		service: service,
	}
}

func (tpaic *ThirdPartyAssociatedItemController) GetAllThirdPartyProcesses(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := tpaic.service.GetAllThirdPartyProcesses(uint(ThirdPartyID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}

func (tpaic *ThirdPartyAssociatedItemController) CreateThirdPartyProcess(c *gin.Context) {
	var input commonDTOs.AssociatedProcessMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	ThirdPartyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	err = tpaic.service.CreateThirdPartyProcess(uint(ThirdPartyID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (tpaic *ThirdPartyAssociatedItemController) DeleteThirdPartyProcess(c *gin.Context) {
	ThirdPartyID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	processID, _ := strconv.ParseUint(c.Param("processId"), 10, 64)
	err := tpaic.service.DeleteThirdPartyProcess(uint(ThirdPartyID), uint(processID))
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Process unmapped from  ThirdParty successfully"})
}

func (tpaic *ThirdPartyAssociatedItemController) GetAllThirdPartyStrategicObjectives(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	objectives, err := tpaic.service.GetAllThirdPartyStrategicObjectives(uint(ThirdPartyID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, objectives)
}

func (tpaic *ThirdPartyAssociatedItemController) CreateThirdPartyStrategicObjective(c *gin.Context) {
	var input commonDTOs.AssociatedStrategicObjectiveMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	ThirdPartyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	err = tpaic.service.CreateThirdPartyStrategicObjective(uint(ThirdPartyID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (tpaic *ThirdPartyAssociatedItemController) DeleteThirdPartyStrategicObjective(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")
	strategicObjectiveIDStr := c.Param("strategicObjectiveId")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategicObjectiveID, err := strconv.ParseUint(strategicObjectiveIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = tpaic.service.DeleteThirdPartyStrategicObjective(uint(ThirdPartyID), uint(strategicObjectiveID))
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Strategic Objective unmapped from third party successfully"})
}

func (tpaic *ThirdPartyAssociatedItemController) GetAllThirdPartyAssets(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	assets, err := tpaic.service.GetAllThirdPartyAssets(uint(ThirdPartyID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, assets)
}

func (tpaic *ThirdPartyAssociatedItemController) CreateThirdPartyAsset(c *gin.Context) {
	var input commonDTOs.AssociatedAssetMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	ThirdPartyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	err = tpaic.service.CreateThirdPartyAsset(uint(ThirdPartyID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (tpaic *ThirdPartyAssociatedItemController) DeleteThirdPartyAsset(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")
	assetIDStr := c.Param("assetId")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = tpaic.service.DeleteThirdPartyAsset(uint(ThirdPartyID), uint(assetID))
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset unmapped from third party successfully"})
}

func (tpaic *ThirdPartyAssociatedItemController) GetAllThirdPartyIncidents(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	incidents, err := tpaic.service.GetAllThirdPartyIncidents(uint(ThirdPartyID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidents)
}

func (tpaic *ThirdPartyAssociatedItemController) CreateThirdPartyIncident(c *gin.Context) {
	var input commonDTOs.AssociatedIncidentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	ThirdPartyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	err = tpaic.service.CreateThirdPartyIncident(uint(ThirdPartyID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (tpaic *ThirdPartyAssociatedItemController) DeleteThirdPartyIncident(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")
	incidentIDStr := c.Param("incidentId")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = tpaic.service.DeleteThirdPartyIncident(uint(ThirdPartyID), uint(incidentID))
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident unmapped from third party successfully"})
}
func (tpaic *ThirdPartyAssociatedItemController) GetAllThirdPartyDocuments(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	documents, err := tpaic.service.GetAllThirdPartyDocuments(uint(ThirdPartyID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, documents)
}

func (tpaic *ThirdPartyAssociatedItemController) CreateThirdPartyDocument(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	ThirdPartyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	err = tpaic.service.CreateThirdPartyDocument(uint(ThirdPartyID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (tpaic *ThirdPartyAssociatedItemController) DeleteThirdPartyDocument(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")
	documentIDStr := c.Param("documentId")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = tpaic.service.DeleteThirdPartyDocument(uint(ThirdPartyID), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from third party successfully"})
}

func (tpaic *ThirdPartyAssociatedItemController) GetAllThirdPartyControls(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := tpaic.service.GetAllThirdPartyControls(uint(ThirdPartyID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}

func (tpaic *ThirdPartyAssociatedItemController) CreateThirdPartyControl(c *gin.Context) {
	var input commonDTOs.AssociatedControlMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	ThirdPartyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	err = tpaic.service.CreateThirdPartyControl(uint(ThirdPartyID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (tpaic *ThirdPartyAssociatedItemController) DeleteThirdPartyControl(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")
	controlIDStr := c.Param("controlId")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlID, err := strconv.ParseUint(controlIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = tpaic.service.DeleteThirdPartyControl(uint(ThirdPartyID), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Control unmapped from third party successfully"})
}

func (tpaic *ThirdPartyAssociatedItemController) GetAllThirdPartyRisks(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := tpaic.service.GetAllThirdPartyRisks(uint(ThirdPartyID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (tpaic *ThirdPartyAssociatedItemController) CreateThirdPartyRisk(c *gin.Context) {
	var input commonDTOs.AssociatedRiskMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	ThirdPartyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	err = tpaic.service.CreateThirdPartyRisk(uint(ThirdPartyID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (tpaic *ThirdPartyAssociatedItemController) DeleteThirdPartyRisk(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")
	riskIDStr := c.Param("riskId")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = tpaic.service.DeleteThirdPartyRisk(uint(ThirdPartyID), uint(riskID))
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Risk unmapped from third party successfully"})
}

func (tpaic *ThirdPartyAssociatedItemController) GetAllThirdPartyKPI(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	kpis, err := tpaic.service.GetAllThirdPartyKPI(uint(ThirdPartyID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, kpis)
}

func (tpaic *ThirdPartyAssociatedItemController) CreateThirdPartyKPI(c *gin.Context) {
	var input commonDTOs.AssociatedKPIMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	ThirdPartyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	err = tpaic.service.CreateThirdPartyKPI(uint(ThirdPartyID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (tpaic *ThirdPartyAssociatedItemController) DeleteThirdPartyKPI(c *gin.Context) {

	ThirdPartyIDStr := c.Param("id")
	kpiIDStr := c.Param("kpiRegisterId")

	ThirdPartyID, err := strconv.ParseUint(ThirdPartyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kpiID, err := strconv.ParseUint(kpiIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = tpaic.service.DeleteThirdPartyKPI(uint(ThirdPartyID), uint(kpiID))
	if err != nil {
		utility.SendInternalErrResp(tpaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KPI unmapped from third party successfully"})
}
