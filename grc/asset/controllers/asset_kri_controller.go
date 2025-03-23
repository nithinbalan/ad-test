package controllers

import (
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetKriController struct {
	service services.AssetKriServiceInterface
	db      *gorm.DB
}

func NewAssetKriController(service services.AssetKriServiceInterface) *AssetKriController {
	return &AssetKriController{
		service: service,
	}
}

func (akc *AssetKriController) GetAllAssetKris(c *gin.Context) {
	assetIDStr := c.Param("id")

	var assetID uint64
	var err error

	if assetIDStr != "" {
		assetID, err = strconv.ParseUint(assetIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		assetID = 0
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"kri_register_ids",
		"asset_register_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(akc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(akc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := assetDTOs.AssetKriListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	assetKris, err := akc.service.GetAll(uint(assetID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(akc.db, c, e.Error())
		}
		return
	}
	c.JSON(http.StatusOK, assetKris)
}

func (akc *AssetKriController) GetAssetKriDetails(c *gin.Context) {
	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("kriId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetKri, err := akc.service.GetFullDetailsByID(uint(assetID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(akc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusOK, assetKri)
}

func (akc *AssetKriController) CreateAssetKri(c *gin.Context) {
	var form assetDTOs.AssetKriFormDTO

	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(akc.db, c, err.Error())
		return
	}

	ids, err := akc.service.Create(uint(assetID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(akc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, commonDTOs.SuccessfulIDsResponse{SuccessfulIds: ids})
}

func (akc *AssetKriController) UpdateAssetKri(c *gin.Context) {
	var form assetDTOs.AssetKriFormDTO

	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(akc.db, c, err.Error())
		return
	}

	err = akc.service.Update(uint(assetID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(akc.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (akc *AssetKriController) DeleteAssetKris(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(akc.db, c, err.Error())
		return
	}

	failedIDs, err := akc.service.Delete(uint(assetID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(akc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
