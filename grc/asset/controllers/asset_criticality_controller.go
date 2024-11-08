package controllers

import (
	"net/http"
	"strconv"

	"jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetCriticalityController struct {
	service services.AssetCriticalityServiceInterface
	db      *gorm.DB
}

func NewAssetCriticalityController(service services.AssetCriticalityServiceInterface) *AssetCriticalityController {
	return &AssetCriticalityController{
		service: service,
	}
}

func (acc *AssetCriticalityController) GetAssetCriticality(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	queryParams := dtos.AssetCriticalityContributionListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
	}
	contributions, err := acc.service.GetAll(uint(id), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, contributions)
}

func (acc *AssetCriticalityController) CreateAssetCriticality(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var contributions []dtos.AssetCriticalityFormDTO
	if err := c.ShouldBindJSON(&contributions); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	err = acc.service.Create(uint(id), contributions, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (acc *AssetCriticalityController) UpdateAssetCriticality(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	criticalityIdStr := c.Param("criticalityId")
	criticalityId, err := strconv.ParseUint(criticalityIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var contribution dtos.AssetCriticalityFormDTO
	if err := c.ShouldBindJSON(&contribution); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	err = acc.service.Update(uint(id), uint(criticalityId), contribution, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (acc *AssetCriticalityController) DeleteAssetCriticality(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	if err := acc.service.Delete(uint(id), currentUserID); err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
