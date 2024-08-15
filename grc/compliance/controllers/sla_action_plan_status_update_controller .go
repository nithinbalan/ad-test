package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SLAActionPlanStatusUpdateController struct {
	service services.SLAActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewSLAActionPlanStatusUpdateController(service services.SLAActionPlanStatusUpdateServiceInterface) *SLAActionPlanStatusUpdateController {
	return &SLAActionPlanStatusUpdateController{
		service: service,
	}
}

func (sasuc *SLAActionPlanStatusUpdateController) GetAllSLAActionPlanStatusUpdates(c *gin.Context) {
	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	slaStatuses, err := sasuc.service.GetAll(uint(slaID), uint(actionPlanID), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sasuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, slaStatuses)
}

func (sasuc *SLAActionPlanStatusUpdateController) CreateSLAActionPlanStatusUpdate(c *gin.Context) {
	var slaStatusUpdateFormDTO complianceDTOs.SLAActionPlanStatusUpdateFormDTO

	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&slaStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sasuc.db, c, err.Error())
		return
	}

	createdID, err := sasuc.service.Create(uint(slaID), uint(actionPlanID), currentUserID, &slaStatusUpdateFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sasuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (sasuc *SLAActionPlanStatusUpdateController) UpdateSLAActionPlanStatusUpdate(c *gin.Context) {
	var slaStatusUpdateFormDTO complianceDTOs.SLAActionPlanStatusUpdateFormDTO

	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&slaStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sasuc.db, c, err.Error())
		return
	}

	err = sasuc.service.Update(uint(slaID), uint(actionPlanID), uint(id), currentUserID, &slaStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sasuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (sasuc *SLAActionPlanStatusUpdateController) GetLatestSLAActionPlanStatusUpdate(c *gin.Context) {
	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	slaActionPlanStatusUpdate, err := sasuc.service.GetLatestSLAActionPlanStatusUpdate(uint(slaID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sasuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, slaActionPlanStatusUpdate)
}
