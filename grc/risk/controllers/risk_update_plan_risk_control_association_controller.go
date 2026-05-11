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

type RiskUpdatePlanRiskControlAssociationController struct {
	service services.RiskUpdatePlanRiskControlAssociationServiceInterface
	db      *gorm.DB
}

func NewRiskUpdatePlanRiskControlAssociationController(service services.RiskUpdatePlanRiskControlAssociationServiceInterface) *RiskUpdatePlanRiskControlAssociationController {
	return &RiskUpdatePlanRiskControlAssociationController{
		service: service,
	}
}

func (rrupc *RiskUpdatePlanRiskControlAssociationController) CreateRiskControlAssociation(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	var formDTO riskDTOs.RiskControlAssociationFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = rrupc.service.CreateRiskControlAssociation(uint(planID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, err)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrupc *RiskUpdatePlanRiskControlAssociationController) GetRiskControlAssociationList(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	isCritical := c.Query("is_critical") == "true"

	risks, err := rrupc.service.GetRiskControlAssociationList(uint(planID), isCritical)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, risks)
}

func (rrupc *RiskUpdatePlanRiskControlAssociationController) MarkRiskControlAssociationAsNotApplicable(c *gin.Context) {
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

	err = rrupc.service.MarkRiskControlAssociationAsNotApplicable(uint(planID), currentUserID, &formDTO)
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
