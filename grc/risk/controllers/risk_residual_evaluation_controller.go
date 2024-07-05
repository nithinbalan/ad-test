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

type RiskResidualEvaluationController struct {
	service services.RiskResidualEvaluationServiceInterface
	db      *gorm.DB
}

func NewRiskResidualEvaluationController(service services.RiskResidualEvaluationServiceInterface) *RiskResidualEvaluationController {
	return &RiskResidualEvaluationController{
		service: service,
	}
}

func (rrec *RiskResidualEvaluationController) GetRiskResidualEvaluation(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskResidualEvaluations, err := rrec.service.GetFullDetailsByID(uint(riskID))

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrec.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskResidualEvaluations)
}

func (rrec *RiskResidualEvaluationController) CreateRiskResidualEvaluation(c *gin.Context) {
	var RiskResidualEvaluationFormDTO riskDTOs.RiskResidualEvaluationFormDTO

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&RiskResidualEvaluationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrec.db, c, err.Error())
		return
	}

	id, err := rrec.service.Create(uint(riskID), currentUserID, &RiskResidualEvaluationFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (rrec *RiskResidualEvaluationController) UpdateRiskResidualEvaluation(c *gin.Context) {
	var RiskResidualEvaluationFormDTO riskDTOs.RiskResidualEvaluationFormDTO

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

	if err := c.ShouldBindJSON(&RiskResidualEvaluationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrec.db, c, err.Error())
		return
	}

	err = rrec.service.Update(uint(riskID), uint(id), currentUserID, &RiskResidualEvaluationFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrec.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrec *RiskResidualEvaluationController) DeleteRiskResidualEvaluation(c *gin.Context) {
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
		utility.SendInternalErrResp(rrec.db, c, err.Error())
		return
	}

	failedID, err := rrec.service.Delete(uint(riskID), uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
