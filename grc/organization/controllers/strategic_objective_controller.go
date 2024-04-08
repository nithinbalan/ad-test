package controllers

import (
	"fmt"
	"jethur/grc/organization/models"
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

func (ccc *StrategicOBjectiveController) GetAllStrategicObjectives(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	strategicObjectives, err := ccc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, strategicObjectives)
}

func (ccc *StrategicOBjectiveController) GetStrategicObjectiveDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategicObjective, err := ccc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, strategicObjective)
}

func (ccc *StrategicOBjectiveController) CreateStrategicObjective(c *gin.Context) {

	var strategicObjective *models.StrategicObjective

	if err := c.ShouldBindJSON(&strategicObjective); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	fmt.Println("hiiiiiiiiiiiiiiiiii")
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	createdStrategicObjective, err := ccc.service.Create(currentUserID, strategicObjective)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ccc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdStrategicObjective)
}

func (ccc *StrategicOBjectiveController) UpdateStrategicObjective(c *gin.Context) {

	strategicObjectiveID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedStrategicObjective *models.StrategicObjective
	if err := c.ShouldBindJSON(&updatedStrategicObjective); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	updatedStrategicObjective, err = ccc.service.Update(uint(strategicObjectiveID), currentUserID, updatedStrategicObjective)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ccc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedStrategicObjective)
}

func (ccc *StrategicOBjectiveController) DeleteStrategicObjectives(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	failedIDs, err := ccc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ccc *StrategicOBjectiveController) ActivateStrategicObjectives(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	err = ccc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (ccc *StrategicOBjectiveController) DeactivateStrategicObjectives(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	err = ccc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
