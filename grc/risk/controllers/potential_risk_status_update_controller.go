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

type PotentialRiskStatusUpdateController struct {
	service services.PotentialRiskStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewPotentialRiskStatusUpdateController(service services.PotentialRiskStatusUpdateServiceInterface) *PotentialRiskStatusUpdateController {
	return &PotentialRiskStatusUpdateController{
		service: service,
	}
}

func (prsuc *PotentialRiskStatusUpdateController) GetAllPotentialRiskStatusUpdates(c *gin.Context) {

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	potentialRiskStatuses, err := prsuc.service.GetAll(uint(id), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(prsuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, potentialRiskStatuses)
}

func (prsuc *PotentialRiskStatusUpdateController) CreatePotentialRiskStatusUpdate(c *gin.Context) {

	var potentialRiskStatusUpdateFormDTO riskDTOs.PotentialRiskStatusUpdateFormDTO

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&potentialRiskStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prsuc.db, c, err.Error())
		return
	}

	createdID, err := prsuc.service.Create(uint(id), currentUserID, &potentialRiskStatusUpdateFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ResourceNotFoundError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(prsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}
