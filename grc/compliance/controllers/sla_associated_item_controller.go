package controllers

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SLAAssociatedItemController struct {
	service services.SLAAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewSLAAssociatedItemController(service services.SLAAssociatedItemServiceInterface) *SLAAssociatedItemController {
	return &SLAAssociatedItemController{
		service: service,
	}
}

func (saic *SLAAssociatedItemController) GetAllSLAProcesses(c *gin.Context) {

	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
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
		utility.SendInternalErrResp(saic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(saic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := complianceDTOs.SLAAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	slas, err := saic.service.GetAllProcesses(uint(slaID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(saic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, slas)
}

func (saic *SLAAssociatedItemController) GetAllAssociatedCounts(c *gin.Context) {
	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	associatedCounts, err := saic.service.GetAllAssociatedCounts(uint(slaID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(saic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, associatedCounts)
}
