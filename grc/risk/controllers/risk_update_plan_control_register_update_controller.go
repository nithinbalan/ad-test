package controllers

import (
	"net/http"
	"strconv"

	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskUpdatePlanControlRegisterUpdateController struct {
	service services.RiskUpdatePlanControlRegisterUpdateServiceInterface
	db      *gorm.DB
}

func NewRiskUpdatePlanControlRegisterUpdateController(service services.RiskUpdatePlanControlRegisterUpdateServiceInterface) *RiskUpdatePlanControlRegisterUpdateController {
	return &RiskUpdatePlanControlRegisterUpdateController{
		service: service,
	}
}

func (rrupc *RiskUpdatePlanControlRegisterUpdateController) GetControlCountsByImplementationStatus(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	counts, err := rrupc.service.GetControlCountsByImplementationStatus(uint(planID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, counts)
}

func (rrupc *RiskUpdatePlanControlRegisterUpdateController) CreateControlTest(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var controlTestFormDTO riskDTOs.UpdatePlanControlTestFormDTO
	if err := c.ShouldBindJSON(&controlTestFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	id, err := rrupc.service.CreateControlTest(uint(planID), currentUserID, &controlTestFormDTO)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, err)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (rrupc *RiskUpdatePlanControlRegisterUpdateController) MarkControlRegisterUpdateAsNotApplicable(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO riskDTOs.MarkTimelineNotApplicableDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	err = rrupc.service.MarkControlRegisterUpdateAsNotApplicable(uint(planID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
