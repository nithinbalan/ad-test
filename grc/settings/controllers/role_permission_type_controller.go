package controllers

import (
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RolePermissionTypeController struct {
	service services.RolePermissionTypeServiceInterface
	db      *gorm.DB
}

func NewRolePermissionTypeController(service services.RolePermissionTypeServiceInterface) *RolePermissionTypeController {
	return &RolePermissionTypeController{
		service: service,
	}
}

func (rptc *RolePermissionTypeController) GetAllRolePermissionTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	rolePermissionTypes, err := rptc.service.GetAll(sort, sortBy, keyword)
	if err != nil {
		utility.SendInternalErrResp(rptc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, rolePermissionTypes)
}

func (rptc *RolePermissionTypeController) GetRolePermissionTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	rolePermissionType, err := rptc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, rolePermissionType)
}
