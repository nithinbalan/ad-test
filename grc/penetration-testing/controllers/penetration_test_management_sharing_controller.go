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

type PenetrationTestManagementSharingController struct {
	service services.PenetrationTestManagementSharingServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestManagementSharingController(service services.PenetrationTestManagementSharingServiceInterface) *PenetrationTestManagementSharingController {
	return &PenetrationTestManagementSharingController{
		service: service,
	}
}

func (ptmc *PenetrationTestManagementSharingController) GetPenetrationTestManagementSharingDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	managementSharing, err := ptmc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ptmc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, managementSharing)
}

func (ptmc *PenetrationTestManagementSharingController) CreatePenetrationTestManagementSharing(c *gin.Context) {
	var formDTO penetrationTestDTOs.PenetrationTestManagementSharingFormDTO

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
		utility.SendInternalErrResp(ptmc.db, c, err.Error())
		return
	}

	createdID, err := ptmc.service.Create(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptmc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (ptmc *PenetrationTestManagementSharingController) UpdatePenetrationTestManagementSharing(c *gin.Context) {
	var formDTO penetrationTestDTOs.PenetrationTestManagementSharingFormDTO

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
		utility.SendInternalErrResp(ptmc.db, c, err.Error())
		return
	}

	err = ptmc.service.Update(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptmc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ptmc *PenetrationTestManagementSharingController) DeletePenetrationTestManagementSharings(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptmc.db, c, err.Error())
		return
	}

	failedID, err := ptmc.service.Delete(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptmc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
