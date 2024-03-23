package controllers

import (
	"jethur/grc/settings/models"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrganizationLevelController struct {
	service services.OrganizationLevelServiceInterface
	db      *gorm.DB
}

func NewOrganizationLevelController(service services.OrganizationLevelServiceInterface) *OrganizationLevelController {
	return &OrganizationLevelController{
		service: service,
	}
}

func (olc *OrganizationLevelController) GetAllOrganizationLevel(c *gin.Context) {

	organizationLevels, err := olc.service.GetAll()

	if err != nil {
		utility.SendInternalErrResp(olc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, organizationLevels)
}

func (olc *OrganizationLevelController) CreateOrganizationLevel(c *gin.Context) {

	var updatedOrganizationLevel []*models.OrganizationLevel
	if err := c.ShouldBindJSON(&updatedOrganizationLevel); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(olc.db, c, err.Error())
		return
	}

	createdOrganizationLevel, err := olc.service.Create(currentUserID, updatedOrganizationLevel)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(olc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, createdOrganizationLevel)
}
