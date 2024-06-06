package controllers

import (
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategicProfileValueController struct {
	service services.StrategicProfileValueServiceInterface
	db      *gorm.DB
}

func NewStrategicProfileValueController(service services.StrategicProfileValueServiceInterface) *StrategicProfileValueController {
	return &StrategicProfileValueController{
		service: service,
	}
}

func (spvc *StrategicProfileValueController) GetStrategicProfileValueDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategicProfileValue, err := spvc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(spvc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, strategicProfileValue)
}

func (spvc *StrategicProfileValueController) CreateStrategicProfileValue(c *gin.Context) {
	var strategicProfileValueFormDTO *dtos.StrategicProfileValueFormDTO

	if err := c.ShouldBindJSON(&strategicProfileValueFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(spvc.db, c, err.Error())
		return
	}

	createdStrategicProfileValueID, err := spvc.service.Create(currentUserID, strategicProfileValueFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(spvc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdStrategicProfileValueID)
}

func (spvc *StrategicProfileValueController) UpdateStrategicProfileValue(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var strategicProfileValueFormDTO *dtos.StrategicProfileValueFormDTO
	if err := c.ShouldBindJSON(&strategicProfileValueFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(spvc.db, c, err.Error())
		return
	}

	err = spvc.service.Update(uint(id), currentUserID, strategicProfileValueFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(spvc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (spvc *StrategicProfileValueController) DeleteStrategicProfileValues(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(spvc.db, c, err.Error())
		return
	}

	failedIDs, err := spvc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(spvc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
