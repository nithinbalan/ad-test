package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskInherentEvaluationController struct {
	service services.RiskInherentEvaluationServiceInterface
	db      *gorm.DB
}

func NewRiskInherentEvaluationController(service services.RiskInherentEvaluationServiceInterface) *RiskInherentEvaluationController {
	return &RiskInherentEvaluationController{
		service: service,
	}
}

func (riec *RiskInherentEvaluationController) GetRiskInherentEvaluation(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskInherentEvaluations, err := riec.service.GetAll(uint(riskID))

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(riec.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskInherentEvaluations)
}

func (riec *RiskInherentEvaluationController) CreateRiskInherentEvaluation(c *gin.Context) {
	var RiskInherentEvaluationFormDTO riskDTOs.RiskInherentEvaluationFormDTO

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&RiskInherentEvaluationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(riec.db, c, err.Error())
		return
	}

	id, err := riec.service.Create(uint(riskID), currentUserID, &RiskInherentEvaluationFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(riec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (riec *RiskInherentEvaluationController) UpdateRiskInherentEvaluation(c *gin.Context) {
	var RiskInherentEvaluationFormDTO riskDTOs.RiskInherentEvaluationFormDTO

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
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

	if err := c.ShouldBindJSON(&RiskInherentEvaluationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(riec.db, c, err.Error())
		return
	}

	err = riec.service.Update(uint(riskID), uint(id), currentUserID, &RiskInherentEvaluationFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(riec.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (riec *RiskInherentEvaluationController) DeleteRiskInherentEvaluation(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
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

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(riec.db, c, err.Error())
		return
	}

	failedID, err := riec.service.Delete(uint(riskID), uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(riec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
