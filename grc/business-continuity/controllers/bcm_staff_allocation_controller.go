package controllers

import (
	bcDTOs "jethur/grc/business-continuity/dtos"
	bcServices "jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BCMStaffAllocationController struct {
	service bcServices.BCMStaffAllocationServiceInterface
	db      *gorm.DB
}

func NewBCMStaffAllocationController(service bcServices.BCMStaffAllocationServiceInterface) *BCMStaffAllocationController {
	return &BCMStaffAllocationController{service: service}
}

func (c *BCMStaffAllocationController) GetFullDetailsByProcessID(ctx *gin.Context) {
	idStr := ctx.Param("processId")

	processId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	result, err := c.service.GetFullDetailsByProcessID(uint(processId))
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *BCMStaffAllocationController) UpdateAllocation(ctx *gin.Context) {

	idStr := ctx.Param("processId")

	processId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	var dto bcDTOs.BCMStaffAllocationUpdateFormDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}

	if err := c.service.UpdateAllocation(uint(processId), currentUserID, &dto); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
		default:
			utility.SendInternalErrResp(c.db, ctx, e.Error())
		}
		return
	}

	ctx.Status(http.StatusOK)
}
