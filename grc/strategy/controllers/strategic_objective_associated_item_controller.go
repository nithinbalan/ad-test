package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategicObjectiveAssociatedItemController struct {
	service services.StrategicObjectiveAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewStrategicObjectiveAssociatedItemController(service services.StrategicObjectiveAssociatedItemServiceInterface) *StrategicObjectiveAssociatedItemController {
	return &StrategicObjectiveAssociatedItemController{
		service: service,
	}
}

func (soai *StrategicObjectiveAssociatedItemController) GetAllClauses(c *gin.Context) {
	strategyIDStr := c.Param("id")

	strategicObjectiveID, err := strconv.ParseUint(strategyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	documents, err := soai.service.GetAllClauses(uint(strategicObjectiveID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soai.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, documents)
}
func (soai *StrategicObjectiveAssociatedItemController) UpdateClauses(c *gin.Context) {
	var ControlClauseFormDTO commonDTOs.AssociatedClauseFormDTO

	strategicObjectiveID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&ControlClauseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, err.Error())
		return
	}

	err = soai.service.UpdateClauses(uint(strategicObjectiveID), currentUserID, &ControlClauseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soai.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Clauses against the stategic objective updated successfully"})
}

func (soai *StrategicObjectiveAssociatedItemController) GetAllStrategicObjectiveRisks(c *gin.Context) {
	strategyIDStr := c.Param("id")

	strategicObjectiveID, err := strconv.ParseUint(strategyIDStr, 10, 64)
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
		utility.SendInternalErrResp(soai.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(soai.db, c, "Invalid value parameter")
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
	risks, err := soai.service.GetAllRisks(uint(strategicObjectiveID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soai.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (soai *StrategicObjectiveAssociatedItemController) CreateStrategicObjectiveRisk(c *gin.Context) {
	var input commonDTOs.AssociatedRiskMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	strategicObjectiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, err.Error())
		return
	}

	err = soai.service.CreateStrategicObjectiveRisk(uint(strategicObjectiveID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soai.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (soai *StrategicObjectiveAssociatedItemController) DeleteStrategicObjectiveRisk(c *gin.Context) {
	strategicObjectiveID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	err := soai.service.DeleteStrategicObjectiveRisk(uint(strategicObjectiveID), uint(riskID))
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Risk unmapped from strategic objective successfully"})
}

func (soai *StrategicObjectiveAssociatedItemController) GetAllStrategicObjectiveDocuments(c *gin.Context) {
	strategyIDStr := c.Param("id")

	strategicObjectiveID, err := strconv.ParseUint(strategyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := soai.service.GetAllStrategicObjectiveDocuments(uint(strategicObjectiveID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soai.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (soai *StrategicObjectiveAssociatedItemController) CreateStrategicObjectiveDocuments(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	strategicObjectiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, err.Error())
		return
	}

	err = soai.service.CreateStrategicObjectiveDocuments(uint(strategicObjectiveID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soai.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (soai *StrategicObjectiveAssociatedItemController) DeleteStrategicObjectiveDocument(c *gin.Context) {
	strategicObjectiveID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	documentID, _ := strconv.ParseUint(c.Param("documentId"), 10, 64)
	err := soai.service.DeleteStrategicObjectiveDocument(uint(strategicObjectiveID), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from strategic objective successfully"})
}

func (soai *StrategicObjectiveAssociatedItemController) GetAllStrategicObjectivControls(c *gin.Context) {
	strategyIDStr := c.Param("id")

	strategicObjectiveID, err := strconv.ParseUint(strategyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := soai.service.GetAllStrategicObjectivControls(uint(strategicObjectiveID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soai.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (soai *StrategicObjectiveAssociatedItemController) CreateStrategicObjectiveControl(c *gin.Context) {
	var input commonDTOs.AssociatedControlMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	strategicObjectiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, err.Error())
		return
	}

	err = soai.service.CreateStrategicObjectiveControl(uint(strategicObjectiveID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soai.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (soai *StrategicObjectiveAssociatedItemController) DeleteStrategicObjectiveControl(c *gin.Context) {
	strategicObjectiveID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)
	err := soai.service.DeleteStrategicObjectiveControl(uint(strategicObjectiveID), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "control unmapped from strategic objective successfully"})
}

func (soai *StrategicObjectiveAssociatedItemController) GetAllStrategicObjectiveProcesses(c *gin.Context) {
	strategyIDStr := c.Param("id")

	strategicObjectiveID, err := strconv.ParseUint(strategyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := soai.service.GetAllStrategicObjectiveProcesses(uint(strategicObjectiveID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soai.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (soai *StrategicObjectiveAssociatedItemController) CreateStrategicObjectiveProcess(c *gin.Context) {
	var input commonDTOs.AssociatedProcessMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	strategicObjectiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, err.Error())
		return
	}

	err = soai.service.CreateStrategicObjectiveProcess(uint(strategicObjectiveID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soai.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (soai *StrategicObjectiveAssociatedItemController) DeleteStrategicObjectiveProcess(c *gin.Context) {
	strategicObjectiveID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	processID, _ := strconv.ParseUint(c.Param("processId"), 10, 64)
	err := soai.service.DeleteStrategicObjectiveProcess(uint(strategicObjectiveID), uint(processID))
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "process unmapped from strategic objective successfully"})
}

func (soai *StrategicObjectiveAssociatedItemController) GetAllStrategicObjectiveAssets(c *gin.Context) {
	strategyIDStr := c.Param("id")

	strategicObjectiveID, err := strconv.ParseUint(strategyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := soai.service.GetAllStrategicObjectiveAssets(uint(strategicObjectiveID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soai.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (soai *StrategicObjectiveAssociatedItemController) CreateStrategicObjectiveAsset(c *gin.Context) {
	var input commonDTOs.AssociatedAssetMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	strategicObjectiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, err.Error())
		return
	}

	err = soai.service.CreateStrategicObjectiveAsset(uint(strategicObjectiveID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soai.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (soai *StrategicObjectiveAssociatedItemController) DeleteStrategicObjectiveAsset(c *gin.Context) {
	strategicObjectiveIDStr := c.Param("id")
	assetIDStr := c.Param("assetId")
	strategicObjectiveID, err := strconv.ParseUint(strategicObjectiveIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = soai.service.DeleteStrategicObjectiveAsset(uint(strategicObjectiveID), uint(assetID))
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "asset unmapped from strategic objective successfully"})
}
