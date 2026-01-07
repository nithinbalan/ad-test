package controllers

import (
	"errors"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		"user_status_ids",
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

func (occ *OrganizationChartController) ChangeReportingUser(c *gin.Context) {

	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid user id"))
		return
	}

	var changeReportingUserFormDTO organizationDTOs.ChangeReportingUserFormDTO
	if err := c.ShouldBindJSON(&changeReportingUserFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = occ.service.ChangeReportingUser(uint(userID), &changeReportingUserFormDTO)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(occ.db, c, err.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
