package controllers

import (
	"jethur/grc/strategy/dtos"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategicObjectiveKPIController struct {
	service services.StrategicObjectiveKPIServiceInterface
	db      *gorm.DB
}

func NewStrategicObjectiveKPIController(service services.StrategicObjectiveKPIServiceInterface) *StrategicObjectiveKPIController {
	return &StrategicObjectiveKPIController{
		service: service,
	}
}

func (sokc *StrategicObjectiveKPIController) CreateStrategicObjectiveKPI(c *gin.Context) {
	strategicObjectiveID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var strategicObjectiveKPIFormDTO *dtos.StrategicObjectiveKPIFormDTO

	if err := c.ShouldBindJSON(&strategicObjectiveKPIFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sokc.db, c, err.Error())
		return
	}
	err = sokc.service.Create(uint(strategicObjectiveID), currentUserID, strategicObjectiveKPIFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sokc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (sokc *StrategicObjectiveKPIController) DeleteStrategicObjectiveKPI(c *gin.Context) {
	// Parse StrategicObjective ID from URL
	strategicObjectiveID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Parse KPIRegister ID from URL
	kpiID, err := strconv.Atoi(c.Param("kpiRegisterId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	// Call service
	err = sokc.service.DeleteStrategicObjectiveKPI(uint(strategicObjectiveID), uint(kpiID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(sokc.db, c, e.Error())
		}
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{"message": "Strategic Objective KPI unlinked successfully"})
}
