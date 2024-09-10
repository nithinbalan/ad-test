package controllers

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditTimelineController struct {
	service services.AuditTimelineServiceInterface
	db      *gorm.DB
}

func NewAuditTimelineController(service services.AuditTimelineServiceInterface) *AuditTimelineController {
	return &AuditTimelineController{
		service: service,
	}
}

func (atc *AuditTimelineController) GetAllAuditTimelines(c *gin.Context) {
	auditIDStr := c.Param("auditId")

	auditID, err := strconv.ParseUint(auditIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	queryParams := auditDTOs.AuditTimelineListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}

	auditTimelines, err := atc.service.GetAll(uint(auditID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(atc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditTimelines)
}
