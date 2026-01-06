package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"
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

	allowedParams := []string{
		"reporting_user_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	params := organizationDTOs.OrganizationChartUserListQueryParamsDTO{
		Filters: filters,
	}

	users, err := occ.service.GetAll(&params)
	if err != nil {
		utility.SendInternalErrResp(occ.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

func (occ *OrganizationChartController) GetUserFullDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	user, err := occ.service.GetUserFullDetailsByID(uint(id))
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
