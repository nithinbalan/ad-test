package controllers

import (
	assetDTOs "jethur/grc/asset/dtos"
	assetService "jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetAssociatedItemController struct {
	service assetService.AssetAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewAssetAssociatedItemController(service assetService.AssetAssociatedItemServiceInterface) *AssetAssociatedItemController {
	return &AssetAssociatedItemController{
		service: service,
	}
}

func (aaic *AssetAssociatedItemController) GetAllAssetControls(c *gin.Context) {

	assetIDStr := c.Param("id")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
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
		utility.SendInternalErrResp(aaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := assetDTOs.AssetAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	controls, err := aaic.service.GetAllControls(uint(assetID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}
