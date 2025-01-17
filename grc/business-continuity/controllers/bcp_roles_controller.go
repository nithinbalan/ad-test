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

func (brc *BcpRolesController) GetAllBcpRoles(ctx *gin.Context) {
	bcpIdStr := ctx.Param("bcpId")
	bcpId, err := strconv.ParseUint(bcpIdStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(ctx, errors.New("invalid business continuity plan ID"))
		return
	}

	sort := ctx.DefaultQuery("sort_order", "asc")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	keyword := ctx.Query("keyword")

	allowedParams := []string{"responsible_user_id", "responsible_unit_id"}
	validSortFields := map[string]string{
		"created_at":     "bcp_roles.created_at",
		"responsibility": "bcp_roles.responsibility",
	}

	filters := utility.ParseQueryFilters(allowedParams, ctx.Request.URL.Query())

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(ctx, errors.New("invalid 'sort_by' parameter"))
		return
	}

	queryParams := businessContinuityDTOs.BcpRolesListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortField,
		Keyword: keyword,
		Filters: filters,
	}

	roles, err := brc.service.GetAll(uint(bcpId), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(brc.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, roles)
}

func (brc *BcpRolesController) GetBcpRoleByID(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	role, err := brc.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utility.SendNotFoundErrResp(ctx)
			return
		}
		utility.SendInternalErrResp(brc.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, role)
}

func (brc *BcpRolesController) CreateBcpRole(ctx *gin.Context) {
	idStr := ctx.Param("bcpId")
	bcpId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	var form businessContinuityDTOs.BcpRolesFormDTO
	if err := ctx.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(brc.db, ctx, err.Error())
		return
	}

	id, err := brc.service.Create(uint(bcpId), currentUserID, &form)
	if err != nil {
		utility.SendInternalErrResp(brc.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (brc *BcpRolesController) UpdateBcpRole(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	var form businessContinuityDTOs.BcpRolesFormDTO
	if err := ctx.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(brc.db, ctx, err.Error())
		return
	}

	err = brc.service.Update(uint(id), currentUserID, &form)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utility.SendNotFoundErrResp(ctx)
			return
		}
		utility.SendInternalErrResp(brc.db, ctx, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (brc *BcpRolesController) DeleteBcpRole(ctx *gin.Context) {
	bcpIdParam := ctx.Param("bcpId")
	bcpId, err := strconv.Atoi(bcpIdParam)
	if err != nil || bcpId <= 0 {
		utility.SendInvalidErrResp(ctx, errors.New("invalid Business Continuity Plan ID"))
		return
	}

	var request struct {
		RoleID uint `json:"role_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	err = brc.service.Delete(request.RoleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utility.SendNotFoundErrResp(ctx)
			return
		}
		utility.SendInternalErrResp(brc.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
