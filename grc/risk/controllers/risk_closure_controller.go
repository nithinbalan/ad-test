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

type RiskClosureController struct {
	service services.RiskClosureServiceInterface
	db      *gorm.DB
}

func NewRiskClosureController(service services.RiskClosureServiceInterface) *RiskClosureController {
	return &RiskClosureController{
		service: service,
	}
}

func (rcc *RiskClosureController) GetRiskClosure(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskClosures, err := rcc.service.GetFullDetailsByID(uint(riskID))

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskClosures)
}

func (rcc *RiskClosureController) CreateRiskClosure(c *gin.Context) {
	var RiskClosureFormDTO riskDTOs.RiskClosureFormDTO

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&RiskClosureFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	id, err := rcc.service.Create(uint(riskID), currentUserID, &RiskClosureFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (rcc *RiskClosureController) UpdateRiskClosure(c *gin.Context) {
	var RiskClosureFormDTO riskDTOs.RiskClosureFormDTO

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

	if err := c.ShouldBindJSON(&RiskClosureFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	err = rcc.service.Update(uint(riskID), uint(id), currentUserID, &RiskClosureFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rcc *RiskClosureController) DeleteRiskClosure(c *gin.Context) {
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
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	failedID, err := rcc.service.Delete(uint(riskID), uint(id), currentUserID)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
