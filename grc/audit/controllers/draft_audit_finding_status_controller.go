package controllers

import (
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DraftAuditFindingStatusController struct {
	service services.DraftAuditFindingStatusServiceInterface
	db      *gorm.DB
}

func NewDraftAuditFindingStatusController(service services.DraftAuditFindingStatusServiceInterface) *DraftAuditFindingStatusController {
	return &DraftAuditFindingStatusController{
		service: service,
	}
}

func (dafsc *DraftAuditFindingStatusController) GetAllDraftAuditFindingStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	draftAuditFindingStatuses, err := dafsc.service.GetAll(sort, sortBy, keyword)

	if err != nil {
		utility.SendInternalErrResp(dafsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, draftAuditFindingStatuses)
}

func (dafsc *DraftAuditFindingStatusController) GetDraftAuditFindingStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	draftAuditFindingStatus, err := dafsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, draftAuditFindingStatus)
}
