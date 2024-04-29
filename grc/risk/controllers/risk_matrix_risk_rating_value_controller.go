package controllers

import (
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskMatrixRiskRatingValueController struct {
	service services.RiskMatrixRiskRatingValueServiceInterface
	db      *gorm.DB
}

func NewRiskMatrixRiskRatingValueController(service services.RiskMatrixRiskRatingValueServiceInterface) *RiskMatrixRiskRatingValueController {
	return &RiskMatrixRiskRatingValueController{
		service: service,
	}
}

func (rmrrvc *RiskMatrixRiskRatingValueController) CreateRiskMatrixRiskRatingValue(c *gin.Context) {
	var riskRatingValue *riskDTOs.RiskMatrixRiskRatingValueFormDTO

	if err := c.ShouldBindJSON(&riskRatingValue); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err := rmrrvc.service.Create(riskRatingValue)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rmrrvc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
