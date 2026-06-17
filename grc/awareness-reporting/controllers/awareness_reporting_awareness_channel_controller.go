package controllers

import (
	"jethur/grc/awareness-reporting/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AwarenessChannelController struct {
	service services.AwarenessChannelServiceInterface
	db      *gorm.DB
}

func NewAwarenessChannelController(service services.AwarenessChannelServiceInterface) *AwarenessChannelController {
	return &AwarenessChannelController{
		service: service,
	}
}

func (acc *AwarenessChannelController) GetAllAwarenessChannels(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	awarenessChannels, err := acc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, awarenessChannels)
}

func (acc *AwarenessChannelController) GetAwarenessChannelDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	awarenessChannel, err := acc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, awarenessChannel)
}

func (acc *AwarenessChannelController) CreateAwarenessChannel(c *gin.Context) {
	var awarenessChannel *models.AwarenessChannel

	if err := c.ShouldBindJSON(&awarenessChannel); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	createdChannel, err := acc.service.Create(currentUserID, awarenessChannel)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdChannel)
}

func (acc *AwarenessChannelController) UpdateAwarenessChannel(c *gin.Context) {
	var updatedChannel *models.AwarenessChannel

	awarenessChannelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedChannel); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	updatedChannel, err = acc.service.Update(uint(awarenessChannelID), currentUserID, updatedChannel)
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

	c.JSON(http.StatusOK, updatedChannel)
}

func (acc *AwarenessChannelController) DeleteAwarenessChannels(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	failedIDs, err := acc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (acc *AwarenessChannelController) ActivateAwarenessChannels(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	err = acc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (acc *AwarenessChannelController) DeactivateAwarenessChannels(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	err = acc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}