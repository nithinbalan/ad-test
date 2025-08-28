package controllers

import (
	"jethur/grc/business-continuity/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BcpAssociatedItemController struct {
	service services.BcpAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewBcpAssociatedItemController(service services.BcpAssociatedItemServiceInterface) *BcpAssociatedItemController {
	return &BcpAssociatedItemController{
		service: service,
	}
}

func (baic *BcpAssociatedItemController) GetAllBcpDocuments(c *gin.Context) {

	BcpIDStr := c.Param("id")

	BcpID, err := strconv.ParseUint(BcpIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := baic.service.GetAllBcpDocuments(uint(BcpID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(baic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}
func (baic *BcpAssociatedItemController) CreateBcpDocument(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	BcpID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, err.Error())
		return
	}

	err = baic.service.CreateBcpDocument(uint(BcpID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(baic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (baic *BcpAssociatedItemController) DeleteBcpDocument(c *gin.Context) {
	BcpID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	documentID, _ := strconv.ParseUint(c.Param("documentId"), 10, 64)
	err := baic.service.DeleteBcpDocument(uint(BcpID), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from Bcp successfully"})
}

func (baic *BcpAssociatedItemController) GetAllBcpControl(c *gin.Context) {

	BcpIDStr := c.Param("id")

	BcpID, err := strconv.ParseUint(BcpIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := baic.service.GetAllBcpControl(uint(BcpID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(baic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}
func (baic *BcpAssociatedItemController) CreateBcpControl(c *gin.Context) {
	var input commonDTOs.AssociatedControlMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	BcpID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, err.Error())
		return
	}

	err = baic.service.CreateBcpControl(uint(BcpID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(baic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (baic *BcpAssociatedItemController) DeleteBcpControl(c *gin.Context) {
	BcpID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)
	err := baic.service.DeleteBcpControl(uint(BcpID), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Control unmapped from Bcp successfully"})
}

func (baic *BcpAssociatedItemController) GetAllBcpKPI(c *gin.Context) {

	BcpIDStr := c.Param("id")

	BcpID, err := strconv.ParseUint(BcpIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	kpis, err := baic.service.GetAllBcpKPI(uint(BcpID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(baic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, kpis)
}

func (baic *BcpAssociatedItemController) CreateBcpKPI(c *gin.Context) {
	var input commonDTOs.AssociatedKPIMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	BcpID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, err.Error())
		return
	}

	err = baic.service.CreateBcpKPI(uint(BcpID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(baic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (baic *BcpAssociatedItemController) DeleteBcpKPI(c *gin.Context) {

	BcpIDStr := c.Param("id")
	kpiRegisterIDStr := c.Param("kpiRegisterId")
	BcpID, err := strconv.ParseUint(BcpIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kpiRegisterID, err := strconv.ParseUint(kpiRegisterIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = baic.service.DeleteBcpKPI(uint(BcpID), uint(kpiRegisterID))
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KPI unmapped from Bcp successfully"})
}

func (baic *BcpAssociatedItemController) GetAllBcpRisk(c *gin.Context) {

	BcpIDStr := c.Param("id")

	BcpID, err := strconv.ParseUint(BcpIDStr, 10, 64)
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
		utility.SendInternalErrResp(baic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, "Invalid page_size parameter")
		return
	}
	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(baic.db, c, "Invalid value parameter")
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
		Filters:  filters,
		Value:    value,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := baic.service.GetAllBcpRisk(uint(BcpID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(baic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (baic *BcpAssociatedItemController) CreateBcpRisk(c *gin.Context) {
	var input commonDTOs.AssociatedRiskMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	BcpID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, err.Error())
		return
	}

	err = baic.service.CreateBcpRisk(uint(BcpID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(baic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (baic *BcpAssociatedItemController) DeleteBcpRisk(c *gin.Context) {

	BcpIDStr := c.Param("id")
	riskIDStr := c.Param("riskId")
	BcpID, err := strconv.ParseUint(BcpIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = baic.service.DeleteBcpRisk(uint(BcpID), uint(riskID))
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Risk unmapped from Bcp successfully"})
}

func (baic *BcpAssociatedItemController) GetAllBcpAsset(c *gin.Context) {

	BcpIDStr := c.Param("id")

	BcpID, err := strconv.ParseUint(BcpIDStr, 10, 64)
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
		utility.SendInternalErrResp(baic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, "Invalid page_size parameter")
		return
	}
	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(baic.db, c, "Invalid value parameter")
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
		Filters:  filters,
		Value:    value,
		Page:     page,
		PageSize: pageSize,
	}
	assets, err := baic.service.GetAllBcpAsset(uint(BcpID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(baic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, assets)
}

func (baic *BcpAssociatedItemController) CreateBcpAsset(c *gin.Context) {
	var input commonDTOs.AssociatedAssetMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	BcpID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, err.Error())
		return
	}

	err = baic.service.CreateBcpAsset(uint(BcpID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(baic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (baic *BcpAssociatedItemController) DeleteBcpAsset(c *gin.Context) {

	BcpIDStr := c.Param("id")
	assetIDStr := c.Param("assetId")
	BcpID, err := strconv.ParseUint(BcpIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = baic.service.DeleteBcpAsset(uint(BcpID), uint(assetID))
	if err != nil {
		utility.SendInternalErrResp(baic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset unmapped from Bcp successfully"})
}
