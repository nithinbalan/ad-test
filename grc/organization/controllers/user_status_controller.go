package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserStatusController struct {
	service services.UserStatusServiceInterface
	db      *gorm.DB
}

func NewUserStatusController(service services.UserStatusServiceInterface) *UserStatusController {
	return &UserStatusController{
		service: service,
	}
}

func (usc *UserStatusController) GetAllUserStatuses(c *gin.Context) {
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	userStatuses, err := usc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(usc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, userStatuses)
}

func (usc *UserStatusController) GetUserStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	userStatus, err := usc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, userStatus)
}
