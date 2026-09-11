package controllers

import (
	tiDTOs "jethur/grc/threat-intelligence/dtos"
	"jethur/grc/threat-intelligence/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThreatIntelligenceRiskLevelValueController struct {
	service services.ThreatIntelligenceRiskLevelValueServiceInterface
	db      *gorm.DB
}

func NewThreatIntelligenceRiskLevelValueController(service services.ThreatIntelligenceRiskLevelValueServiceInterface) *ThreatIntelligenceRiskLevelValueController {
	return &ThreatIntelligenceRiskLevelValueController{
		service: service,
	}
}

func (tirlvc *ThreatIntelligenceRiskLevelValueController) GetAllThreatIntelligenceRiskLevelValues(c *gin.Context) {
	riskRatingValues, err := tirlvc.service.GetRiskAll()
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tirlvc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskRatingValues)
}

func (tirlvc *ThreatIntelligenceRiskLevelValueController) CreateThreatIntelligenceRiskLevelValue(c *gin.Context) {
	var riskRatingValue *tiDTOs.ThreatIntelligenceRiskLevelValueFormDTO

	if err := c.ShouldBindJSON(&riskRatingValue); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err := tirlvc.service.Create(riskRatingValue)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(tirlvc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tirlvc *ThreatIntelligenceRiskLevelValueController) GetThreatIntelligenceRiskLevelMatrix(c *gin.Context) {
	riskMatrix, err := tirlvc.service.GetRiskMatrix()

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tirlvc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskMatrix)
}
