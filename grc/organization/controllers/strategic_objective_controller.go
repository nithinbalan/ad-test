package controllers

import (
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategicOBjectiveController struct {
	service services.StrategicObjectiveServiceInterface
	db      *gorm.DB
}

func NewStrategicOBjectiveController(service services.StrategicObjectiveServiceInterface) *StrategicOBjectiveController {
	return &StrategicOBjectiveController{
		service: service,
	}
}

func (soc *StrategicOBjectiveController) GetAllStrategicObjectives(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	strategicObjectives, err := soc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, strategicObjectives)
}

func (soc *StrategicOBjectiveController) GetStrategicObjectiveDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategicObjective, err := soc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, strategicObjective)
}

func (soc *StrategicOBjectiveController) CreateStrategicObjective(c *gin.Context) {

	var strategicObjectiveFormDTO *dtos.StrategicObjectiveFormDTO

	if err := c.ShouldBindJSON(&strategicObjectiveFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	createdStrategicObjectiveID, err := soc.service.Create(currentUserID, strategicObjectiveFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdStrategicObjectiveID)
}

func (soc *StrategicOBjectiveController) UpdateStrategicObjective(c *gin.Context) {

	strategicObjectiveID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var strategicObjectiveFormDTO *dtos.StrategicObjectiveFormDTO
	if err := c.ShouldBindJSON(&strategicObjectiveFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	err = soc.service.Update(uint(strategicObjectiveID), currentUserID, strategicObjectiveFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (soc *StrategicOBjectiveController) DeleteStrategicObjectives(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	failedIDs, err := soc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (soc *StrategicOBjectiveController) ActivateStrategicObjectives(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	err = soc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (soc *StrategicOBjectiveController) DeactivateStrategicObjectives(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	err = soc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
