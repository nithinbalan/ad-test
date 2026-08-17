package controllers

import (
	"net/http"
	"strconv"

	socDtos "jethur/grc/security-operations/dtos"
	"jethur/grc/security-operations/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CSLogSourceDirectoryController struct {
	service services.CSLogSourceDirectoryServiceInterface
	db      *gorm.DB
}

func NewCSLogSourceDirectoryController(service services.CSLogSourceDirectoryServiceInterface) *CSLogSourceDirectoryController {
	return &CSLogSourceDirectoryController{service: service}
}

func (csdc *CSLogSourceDirectoryController) GetAll(ctx *gin.Context) {
	sort := ctx.Query("sort_order")
	sortBy := ctx.Query("sort_by")
	keyword := ctx.Query("keyword")

	var dataEntityTypeID, responsibleEntityID, beneficiaryEntityID *uint
	if v := ctx.Query("data_entity_type_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			u := uint(id)
			dataEntityTypeID = &u
		}
	}
	if v := ctx.Query("responsible_entity_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			u := uint(id)
			responsibleEntityID = &u
		}
	}
	if v := ctx.Query("beneficiary_entity_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			u := uint(id)
			beneficiaryEntityID = &u
		}
	}

	entries, err := csdc.service.GetAll(sort, sortBy, keyword, dataEntityTypeID, responsibleEntityID, beneficiaryEntityID)
	if err != nil {
		utility.SendInternalErrResp(csdc.db, ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, entries)
}

func (csdc *CSLogSourceDirectoryController) GetDetails(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	entry, err := csdc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}
	ctx.JSON(http.StatusOK, entry)
}

func (csdc *CSLogSourceDirectoryController) Create(ctx *gin.Context) {
	var form *socDtos.CSLogSourceDirectoryFormDTO
	if err := ctx.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(csdc.db, ctx, err.Error())
		return
	}

	created, err := csdc.service.Create(currentUserID, form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
			return
		default:
			utility.SendInternalErrResp(csdc.db, ctx, e.Error())
			return
		}
	}

	ctx.JSON(http.StatusCreated, created)
}

func (csdc *CSLogSourceDirectoryController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	var form *socDtos.CSLogSourceDirectoryFormDTO
	if err := ctx.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(csdc.db, ctx, err.Error())
		return
	}

	updated, err := csdc.service.Update(uint(id), currentUserID, form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
			return
		default:
			utility.SendInternalErrResp(csdc.db, ctx, e.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, updated)
}

func (csdc *CSLogSourceDirectoryController) Delete(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(csdc.db, ctx, err.Error())
		return
	}

	failedIDs, err := csdc.service.Delete(req.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(csdc.db, ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
