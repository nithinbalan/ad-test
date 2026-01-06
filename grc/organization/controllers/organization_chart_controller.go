package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
    "strconv"
	"jethur/grc/organization/services"
	"jethur/utility"
)

type OrganizationChartController struct {
	service services.OrganizationChartServiceInterface
	db      *gorm.DB
}

func NewOrganizationChartController(service services.OrganizationChartServiceInterface) *OrganizationChartController {
	return &OrganizationChartController{
		service: service,
	}
}

func (occ *OrganizationChartController) GetUsers(c *gin.Context) {
	users, err := occ.service.GetUsers()
	if err != nil {
		utility.SendInternalErrResp(occ.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

func (occ *OrganizationChartController) GetUserDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	user, err := occ.service.GetUserDetailsByID(uint(id))
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(nil, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, user)
}
