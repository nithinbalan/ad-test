package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/grc/penetration-testing/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenetrationTestEnvironmentPreparationController struct {
	service services.PenetrationTestEnvironmentPreparationServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestEnvironmentPreparationController(service services.PenetrationTestEnvironmentPreparationServiceInterface) *PenetrationTestEnvironmentPreparationController {
	return &PenetrationTestEnvironmentPreparationController{
		service: service,
	}
}

func (ptpc *PenetrationTestEnvironmentPreparationController) GetPenetrationTestEnvironmentPreparationDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	environmentPreparation, err := ptpc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ptpc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, environmentPreparation)
}

func (ptpc *PenetrationTestEnvironmentPreparationController) CreatePenetrationTestEnvironmentPreparation(c *gin.Context) {
	var formDTO penetrationTestDTOs.PenetrationTestEnvironmentPreparationFormDTO

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptpc.db, c, err.Error())
		return
	}

	createdID, err := ptpc.service.Create(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptpc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (ptpc *PenetrationTestEnvironmentPreparationController) UpdatePenetrationTestEnvironmentPreparation(c *gin.Context) {
	var formDTO penetrationTestDTOs.PenetrationTestEnvironmentPreparationFormDTO

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptpc.db, c, err.Error())
		return
	}

	err = ptpc.service.Update(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptpc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ptpc *PenetrationTestEnvironmentPreparationController) DeletePenetrationTestEnvironmentPreparations(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptpc.db, c, err.Error())
		return
	}

	failedID, err := ptpc.service.Delete(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptpc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
