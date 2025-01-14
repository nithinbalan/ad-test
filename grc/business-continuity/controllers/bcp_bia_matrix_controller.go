package controllers

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BcpBiaMatrixController struct {
	service services.BcpBiaMatrixServiceInterface
	db      *gorm.DB
}

func NewBcpBiaMatrixController(service services.BcpBiaMatrixServiceInterface) *BcpBiaMatrixController {
	return &BcpBiaMatrixController{
		service: service,
	}
}

func (bbmc *BcpBiaMatrixController) GetBcpBiaMatrixDetails(c *gin.Context) {
	idStr := c.Param("bcpId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	businessContinuityPlan, err := bbmc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(bbmc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, businessContinuityPlan)
}

func (bbmc *BcpBiaMatrixController) UpdateBcpBiaMatrix(c *gin.Context) {
	var businessContinuityPlanFormDTO businessContinuityDTOs.BcpProcessFormDTO

	businessContinuityPlanID, err := strconv.Atoi(c.Param("bcpId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&businessContinuityPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bbmc.db, c, err.Error())
		return
	}

	err = bbmc.service.Update(uint(businessContinuityPlanID), currentUserID, &businessContinuityPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(bbmc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
