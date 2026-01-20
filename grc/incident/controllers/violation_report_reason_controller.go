package controllers

import (
	"jethur/grc/incident/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ViolationReportReasonController struct {
	service services.ViolationReportReasonServiceInterface
	db      *gorm.DB
}

func NewViolationReportReasonController(service services.ViolationReportReasonServiceInterface) *ViolationReportReasonController {
	return &ViolationReportReasonController{
		service: service,
	}
}

func (vrrc *ViolationReportReasonController) GetAllViolationReportReasons(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	violationReportReasons, err := vrrc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(vrrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, violationReportReasons)
}

func (vrrc *ViolationReportReasonController) GetViolationReportReasonDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	violationReportReason, err := vrrc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, violationReportReason)
}

func (vrrc *ViolationReportReasonController) CreateViolationReportReason(c *gin.Context) {
	var vrr *models.ViolationReportReason

	if err := c.ShouldBindJSON(&vrr); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrrc.db, c, err.Error())
		return
	}

	createdViolationReportReason, err := vrrc.service.Create(currentUserID, vrr)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(vrrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdViolationReportReason)
}

func (vrrc *ViolationReportReasonController) UpdateViolationReportReason(c *gin.Context) {
	var updatedVRR *models.ViolationReportReason

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedVRR); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrrc.db, c, err.Error())
		return
	}

	updatedVRR, err = vrrc.service.Update(uint(id), currentUserID, updatedVRR)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(vrrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedVRR)
}

func (vrrc *ViolationReportReasonController) DeleteViolationReportReasons(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrrc.db, c, err.Error())
		return
	}

	failedIDs, err := vrrc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(vrrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
