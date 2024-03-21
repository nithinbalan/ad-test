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

// func (olc *OrganizationLevelController) CreateOrganizationLevel(c *gin.Context) {
// 	fmt.Println("hiii")
// 	var organizationLevel []*models.OrganizationLevel

// 	if err := c.ShouldBindJSON(&organizationLevel); err != nil {
// 		utility.SendInvalidErrResp(c, err)
// 		return
// 	}

// 	currentUserID, err := utility.GetCurrentUserID(c)
// 	if err != nil {
// 		utility.SendInternalErrResp(olc.db, c, err.Error())
// 		return
// 	}

// 	createdOrganizationLevel, err := olc.service.Create(currentUserID, organizationLevel)
// 	if err != nil {
// 		switch e := err.(type) {
// 		case *utility.ValidationError:
// 			utility.SendInvalidErrResp(c, e)
// 			return
// 		default:
// 			utility.SendInternalErrResp(olc.db, c, e.Error())
// 			return
// 		}
// 	}

// 	c.JSON(http.StatusCreated, createdOrganizationLevel)
// }

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

// func (olc *OrganizationLevelController) DeleteOrganizationLevel(c *gin.Context) {
// 	// Extract the ID from the URL parameter
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		utility.SendInvalidErrResp(c, err)
// 		return
// 	}

// 	currentUserID, err := utility.GetCurrentUserID(c)
// 	// Call the service layer to delete the organization level
// 	err = olc.service.Delete(uint(id), currentUserID)
// 	if err != nil {
// 		// Handle potential errors, such as "not found" or "internal error"
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// If no error, return a success response
// 	c.JSON(http.StatusOK, gin.H{"message": "Organization level deleted successfully"})
// }
