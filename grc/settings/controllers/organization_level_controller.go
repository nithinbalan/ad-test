package controllers

import (
	"jethur/grc/models"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

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

func (olc *OrganizationLevelController) GetAllOrganizationLevels(c *gin.Context) {
	organizationLevels, err := olc.service.GetAll()

	if err != nil {
		utility.SendInternalErrResp(olc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, organizationLevels)
}

func (olc *OrganizationLevelController) CreateOrganizationLevels(c *gin.Context) {
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

	createdOrganizationLevel, err := olc.service.CreateOrganizationLevels(currentUserID, updatedOrganizationLevel)
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

func (olc *OrganizationLevelController) CreateOrganizationLevel(c *gin.Context) {

	var orgLevel *models.OrganizationLevel

	if err := c.ShouldBindJSON(&orgLevel); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(olc.db, c, err.Error())
		return
	}

	createdOrgLevel, err := olc.service.Create(currentUserID, orgLevel)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(olc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdOrgLevel)
}

func (olc *OrganizationLevelController) UpdateOrganizationLevel(c *gin.Context) {

	orgLevelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedOrgLevel *models.OrganizationLevel
	if err := c.ShouldBindJSON(&updatedOrgLevel); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(olc.db, c, err.Error())
		return
	}

	updatedOrgLevel, err = olc.service.Update(uint(orgLevelID), currentUserID, updatedOrgLevel)
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

	c.JSON(http.StatusOK, updatedOrgLevel)
}

func (olc *OrganizationLevelController) DeleteOrganizationLevel(c *gin.Context) {

	orgLevelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(olc.db, c, err.Error())
		return
	}

	if err := olc.service.Delete(uint(orgLevelID), currentUserID); err != nil {
		switch e := err.(type) {
		case *utility.ResourceConflictError:
			utility.SendConflictErrResp(c, e.Error())
			return
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

	c.JSON(http.StatusOK, gin.H{"message": "Organization level deleted successfully"})
}
