package controllers

import (
	"errors"
	"net/http"
	"strconv"

	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BcpRolesController struct {
	service services.BcpRolesServiceInterface
	db      *gorm.DB
}

func NewBcpRolesController(service services.BcpRolesServiceInterface) *BcpRolesController {
	return &BcpRolesController{
		service: service,
	}
}

func (brc *BcpRolesController) GetAllBcpRoles(c *gin.Context) {
	bcpIdStr := c.Param("bcpId")
	bcpId, err := strconv.ParseUint(bcpIdStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid business continuity plan ID"))
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "bcp_roles.responsibility")
	keyword := c.Query("keyword")

	allowedParams := []string{"responsible_user_id", "backup_user_id", "responsible_unit_id"}
	validSortFields := map[string]string{
		"responsibility":                 "bcp_roles.responsibility",
		"responsible_user_first_name":    "responsible_user.first_name_en",
		"responsible_user_last_name":     "responsible_user.last_name_en",
		"responsible_user_first_name_ar": "responsible_user.first_name_ar",
		"responsible_user_last_name_ar":  "responsible_user.last_name_ar",
		"responsible_unit_title_en":      "organization_units.title_en",
		"responsible_unit_title_ar":      "organization_units.title_ar",
		"backup_user_first_name":         "backup_user.first_name_en",
		"backup_user_last_name":          "backup_user.last_name_en",
		"backup_user_first_name_ar":      "backup_user.first_name_ar",
		"backup_user_last_name_ar":       "backup_user.last_name_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	queryParams := businessContinuityDTOs.BcpRolesListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortField,
		Keyword: keyword,
		Filters: filters,
	}

	roles, err := brc.service.GetAll(&queryParams, uint(bcpId))
	if err != nil {
		utility.SendInternalErrResp(brc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, roles)
}

func (brc *BcpRolesController) GetBcpRoleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	role, err := brc.service.GetByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(brc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, role)
}

func (brc *BcpRolesController) CreateBcpRole(c *gin.Context) {
	idStr := c.Param("bcpId")
	bcpId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form businessContinuityDTOs.BcpRolesFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(brc.db, c, err.Error())
		return
	}

	id, err := brc.service.Create(uint(bcpId), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(brc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (brc *BcpRolesController) UpdateBcpRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form businessContinuityDTOs.BcpRolesFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(brc.db, c, err.Error())
		return
	}

	err = brc.service.Update(uint(id), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(brc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (brc *BcpRolesController) DeleteBcpRoles(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(brc.db, c, err.Error())
		return
	}

	failedIDs, err := brc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(brc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
