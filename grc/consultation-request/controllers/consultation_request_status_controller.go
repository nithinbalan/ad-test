package controllers

import (
	"jethur/grc/consultation-request/services"

	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ConsultationRequestStatusController handles consultation request statuses
type ConsultationRequestStatusController struct {
	service services.ConsultationRequestStatusServiceInterface
	db      *gorm.DB
}

// NewConsultationRequestStatusController initializes the controller
func NewConsultationRequestStatusController(service services.ConsultationRequestStatusServiceInterface) *ConsultationRequestStatusController {
	return &ConsultationRequestStatusController{
		service: service,
	}
}

// GetAllConsultationRequestStatuses fetches all statuses with optional filters
func (crsc *ConsultationRequestStatusController) GetAllConsultationRequestStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	statuses, err := crsc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(crsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, statuses)
}

// GetConsultationRequestStatusDetails fetches a single status by ID
func (crsc *ConsultationRequestStatusController) GetConsultationRequestStatusDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	status, err := crsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, status)
}

func (crc *ConsultationRequestController) GetAllConsultationRequestComment(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comments, err := crc.service.GetAllConsultationRequestComments(uint(id))
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, comments)
}
