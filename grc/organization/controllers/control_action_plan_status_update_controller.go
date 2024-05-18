package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControlActionPlanStatusUpdateController struct {
	service services.ControlActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewControlActionPlanStatusUpdateController(service services.ControlActionPlanStatusUpdateServiceInterface) *ControlActionPlanStatusUpdateController {
	return &ControlActionPlanStatusUpdateController{
		service: service,
	}
}

func (casuc *ControlActionPlanStatusUpdateController) GetAllControlActionPlanStatusUpdates(c *gin.Context) {
	controlID, err := strconv.Atoi(c.Param("controlId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlActionPlanID, err := strconv.Atoi(c.Param("controlActionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	controlActionPlanStatusUpdates, err := casuc.service.GetAll(uint(controlID), uint(controlActionPlanID), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(casuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controlActionPlanStatusUpdates)
}

func (casuc *ControlActionPlanStatusUpdateController) CreateControlActionPlanStatusUpdate(c *gin.Context) {
	var controlActionPlanStatusUpdateFormDTO organizationDTOs.ControlActionPlanStatusUpdateFormDTO

	controlID, err := strconv.Atoi(c.Param("controlId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlActionPlanID, err := strconv.Atoi(c.Param("controlActionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&controlActionPlanStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(casuc.db, c, err.Error())
		return
	}

	createdID, err := casuc.service.Create(uint(controlID), uint(controlActionPlanID), currentUserID, &controlActionPlanStatusUpdateFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(casuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (casuc *ControlActionPlanStatusUpdateController) UpdateControlActionPlanStatusUpdate(c *gin.Context) {
	var controlActionPlanStatusUpdateFormDTO organizationDTOs.ControlActionPlanStatusUpdateFormDTO

	controlID, err := strconv.Atoi(c.Param("controlId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlActionPlanID, err := strconv.Atoi(c.Param("controlActionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlActionPlanUpdateID, err := strconv.Atoi(c.Param("controlActionPlanUpdateId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&controlActionPlanStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(casuc.db, c, err.Error())
		return
	}

	err = casuc.service.Update(uint(controlID), uint(controlActionPlanID), uint(controlActionPlanUpdateID), currentUserID, &controlActionPlanStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(casuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (crc *ControlActionPlanStatusUpdateController) GetControlActionPlanStatusUpdateDetails(c *gin.Context) {
	controlID, err := strconv.Atoi(c.Param("controlId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlActionPlanID, err := strconv.Atoi(c.Param("controlActionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlActionPlanUpdateID, err := strconv.Atoi(c.Param("controlActionPlanUpdateId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlTest, err := crc.service.GetByID(uint(controlID), uint(controlActionPlanID), uint(controlActionPlanUpdateID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controlTest)
}
