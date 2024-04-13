package controllers

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategicProfileController struct {
	service services.StrategicProfileServiceInterface
	db      *gorm.DB
}

func NewStrategicProfileController(service services.StrategicProfileServiceInterface) *StrategicProfileController {
	return &StrategicProfileController{
		service: service,
	}
}

func (spc *StrategicProfileController) GetStrategicProfile(c *gin.Context) {

	strategicProfile, err := spc.service.GetStrategicProfile()
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(spc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, strategicProfile)
}

func (spc *StrategicProfileController) UpdateStrategicProfile(c *gin.Context) {

	var strategicProfileFormDTO organizationDTOs.StrategicProfileFormDTO
	if err := c.ShouldBindJSON(&strategicProfileFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	err = spc.service.Update(currentUserID, &strategicProfileFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(spc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (spc *StrategicProfileController) UpdateVision(c *gin.Context) {

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	var strategicProfileVisionFormDTO organizationDTOs.StrategicProfileVisionFormDTO
	if err := c.ShouldBindJSON(&strategicProfileVisionFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = spc.service.UpdateVision(currentUserID, &strategicProfileVisionFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(spc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)

}

func (spc *StrategicProfileController) UpdateMission(c *gin.Context) {

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	var strategicProfileMissionFormDTO organizationDTOs.StrategicProfileMisionFormDTO
	if err := c.ShouldBindJSON(&strategicProfileMissionFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = spc.service.UpdateMission(currentUserID, &strategicProfileMissionFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(spc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)

}
