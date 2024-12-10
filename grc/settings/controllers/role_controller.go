package controllers

import (
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleController struct {
	service services.RoleServiceInterface
	db      *gorm.DB
}

func NewRoleController(service services.RoleServiceInterface) *RoleController {
	return &RoleController{
		service: service,
	}
}

func (rc *RoleController) GetAllRoles(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	roles, err := rc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (rc *RoleController) GetRoleDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	role, err := rc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, role)
}

func (rc *RoleController) CreateRole(c *gin.Context) {

	var role *models.Role

	if err := c.ShouldBindJSON(&role); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	createdRole, err := rc.service.Create(currentUserID, role)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdRole)
}

func (rc *RoleController) UpdateRole(c *gin.Context) {

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedRole *models.Role
	if err := c.ShouldBindJSON(&updatedRole); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	updatedRole, err = rc.service.Update(uint(roleID), currentUserID, updatedRole)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedRole)
}

func (rc *RoleController) DeleteRoles(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	failedIDs, err := rc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (rc *RoleController) ActivateRoles(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	err = rc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (rc *RoleController) DeactivateRoles(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	err = rc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (rc *RoleController) GetRolePermissions(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	quickFilter := []string{"active"}

	queryParams := dtos.ModuleListQueryParamsDTO{
		Sort:              sort,
		SortBy:            sortBy,
		Keyword:           keyword,
		QuickFilter:       quickFilter,
		IsActiveSubModule: true,
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	role, err := rc.service.GetAllWithPermissions(&queryParams, uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, role)
}

func (rc *RoleController) UpdateRolePermissions(c *gin.Context) {
	var rolePermissionsFormDTOs []*dtos.RolePermissionFormDTO

	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&rolePermissionsFormDTOs); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	err = rc.service.UpdateRolePermissions(uint(roleID), currentUserID, rolePermissionsFormDTOs)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
