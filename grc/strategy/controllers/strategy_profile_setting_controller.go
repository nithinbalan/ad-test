package controllers

import (
	"jethur/grc/strategy/dtos"
	"jethur/grc/strategy/services"
	"strconv"

	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategyProfileSettingController struct {
	service services.StrategyProfileSettingServiceInterface
	db      *gorm.DB
}

func NewStrategyProfileSettingController(service services.StrategyProfileSettingServiceInterface) *StrategyProfileSettingController {
	return &StrategyProfileSettingController{
		service: service,
	}
}

func (mc *StrategyProfileSettingController) GetAllStrategyProfileSettings(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	queryParams := dtos.StrategyProfileSettingListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		QuickFilter: quickFilter,
	}
	modules, err := mc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(mc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, modules)
}

func (ic *StrategyProfileSettingController) GetStrategyProfileSettingDetail(c *gin.Context) {

	settingIDStr := c.Param("settingId")

	settingID, err := strconv.ParseUint(settingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	module, err := ic.service.GetFullDetailsByID(uint(settingID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, module)

}

func (mc *StrategyProfileSettingController) EnableStrategyProfileSetting(c *gin.Context) {
	settingId, err := strconv.Atoi(c.Param("settingId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	_, err = mc.service.EnableStrategyProfileSetting(uint(settingId))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mc.db, c, e.Error())
			return
		}
	}
}

func (mc *StrategyProfileSettingController) DisableStrategyProfileSetting(c *gin.Context) {
	settingId, err := strconv.Atoi(c.Param("settingId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	_, err = mc.service.DisableStrategyProfileSetting(uint(settingId))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mc.db, c, e.Error())
			return
		}
	}
}
