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

type RiskUpdatePlanRiskMitigationController struct {
	service services.RiskUpdatePlanRiskMitigationServiceInterface
	db      *gorm.DB
}

func NewRiskUpdatePlanRiskMitigationController(service services.RiskUpdatePlanRiskMitigationServiceInterface) *RiskUpdatePlanRiskMitigationController {
	return &RiskUpdatePlanRiskMitigationController{
		service: service,
	}
}

func (ruprmc *RiskUpdatePlanRiskMitigationController) UpdateRiskMitigation(c *gin.Context) {
	var formDTO riskDTOs.RiskMitigationUpdatePlanFormDTO

	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ruprmc.db, c, err.Error())
		return
	}

	err = ruprmc.service.UpdateRiskMitigation(uint(planID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ruprmc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ruprmc *RiskUpdatePlanRiskMitigationController) DeleteRiskMitigation(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskIDStr := c.Param("riskId")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ruprmc.db, c, err.Error())
		return
	}

	err = ruprmc.service.DeleteRiskMitigation(uint(planID), currentUserID, uint(riskID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ruprmc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ruprmc *RiskUpdatePlanRiskMitigationController) MarkRiskMitigationAsNotApplicable(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ruprmc.db, c, err.Error())
		return
	}

	var formDTO riskDTOs.MarkTimelineNotApplicableDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = ruprmc.service.MarkRiskMitigationAsNotApplicable(uint(planID), currentUserID, &formDTO)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ruprmc.db, c, err.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ruprmc *RiskUpdatePlanRiskMitigationController) GetActionPlanDetails(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskIDStr := c.Param("riskId")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")
	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanDetails, err := ruprmc.service.GetActionPlanDetails(uint(planID), uint(riskID), uint(actionPlanID))
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ruprmc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, actionPlanDetails)
}

func (ruprmc *RiskUpdatePlanRiskMitigationController) CreateActionPlan(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskIDStr := c.Param("riskId")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ruprmc.db, c, err.Error())
		return
	}

	var formDTO riskDTOs.RiskActionPlanFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	actionPlanID, err := ruprmc.service.CreateActionPlan(uint(planID), uint(riskID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ruprmc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": actionPlanID})
}

func (ruprmc *RiskUpdatePlanRiskMitigationController) UpdateActionPlan(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskIDStr := c.Param("riskId")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")
	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ruprmc.db, c, err.Error())
		return
	}

	var formDTO riskDTOs.RiskActionPlanFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = ruprmc.service.UpdateActionPlan(uint(planID), uint(riskID), uint(actionPlanID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ruprmc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ruprmc *RiskUpdatePlanRiskMitigationController) DeleteActionPlan(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskIDStr := c.Param("riskId")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ruprmc.db, c, err.Error())
		return
	}

	var requestDTO struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	failedIDs, err := ruprmc.service.DeleteActionPlan(uint(planID), uint(riskID), requestDTO.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ruprmc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
