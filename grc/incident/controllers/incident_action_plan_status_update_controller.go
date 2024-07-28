package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentActionPlanStatusUpdateController struct {
	service services.IncidentActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewIncidentActionPlanStatusUpdateController(service services.IncidentActionPlanStatusUpdateServiceInterface) *IncidentActionPlanStatusUpdateController {
	return &IncidentActionPlanStatusUpdateController{
		service: service,
	}
}

func (iasuc *IncidentActionPlanStatusUpdateController) GetAllIncidentActionPlanStatusUpdates(c *gin.Context) {
	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 32)
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

	incidentStatuses, err := iasuc.service.GetAll(uint(incidentID), uint(actionPlanID), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(iasuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidentStatuses)
}

func (iasuc *IncidentActionPlanStatusUpdateController) CreateIncidentActionPlanStatusUpdate(c *gin.Context) {
	var incidentStatusUpdateFormDTO incidentDTOs.IncidentActionPlanStatusUpdateFormDTO

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 32)
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

	if err := c.ShouldBindJSON(&incidentStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iasuc.db, c, err.Error())
		return
	}

	createdID, err := iasuc.service.Create(uint(incidentID), uint(actionPlanID), currentUserID, &incidentStatusUpdateFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iasuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (iasuc *IncidentActionPlanStatusUpdateController) UpdateIncidentActionPlanStatusUpdate(c *gin.Context) {
	var incidentStatusUpdateFormDTO incidentDTOs.IncidentActionPlanStatusUpdateFormDTO

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
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

	if err := c.ShouldBindJSON(&incidentStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iasuc.db, c, err.Error())
		return
	}

	err = iasuc.service.Update(uint(incidentID), uint(actionPlanID), uint(id), currentUserID, &incidentStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iasuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}


func (iasuc *IncidentActionPlanStatusUpdateController) GetLatestIncidentActionPlanStatusUpdate(c *gin.Context) {
	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
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

	incidentActionPlanStatusUpdate, err := iasuc.service.GetLatestIncidentActionPlanStatusUpdate(uint(incidentID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(iasuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, incidentActionPlanStatusUpdate)
}