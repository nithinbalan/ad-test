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

type ViolationReportChannelController struct {
	service services.ViolationReportChannelServiceInterface
	db      *gorm.DB
}

func NewViolationReportChannelController(service services.ViolationReportChannelServiceInterface) *ViolationReportChannelController {
	return &ViolationReportChannelController{
		service: service,
	}
}

func (vrc *ViolationReportChannelController) GetAllViolationReportChannels(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	violationReportChannels, err := vrc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, violationReportChannels)
}

func (vrc *ViolationReportChannelController) GetViolationReportChannelDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	violationReportChannel, err := vrc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, violationReportChannel)
}

func (vrc *ViolationReportChannelController) CreateViolationReportChannel(c *gin.Context) {
	var vrcModel *models.ViolationReportChannel

	if err := c.ShouldBindJSON(&vrcModel); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	createdViolationReportChannel, err := vrc.service.Create(currentUserID, vrcModel)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(vrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdViolationReportChannel)
}

func (vrc *ViolationReportChannelController) UpdateViolationReportChannel(c *gin.Context) {
	var updatedVRC *models.ViolationReportChannel

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedVRC); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	updatedVRC, err = vrc.service.Update(uint(id), currentUserID, updatedVRC)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(vrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedVRC)
}

func (vrc *ViolationReportChannelController) DeleteViolationReportChannels(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	failedIDs, err := vrc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
