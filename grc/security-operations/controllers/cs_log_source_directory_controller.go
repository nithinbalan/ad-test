package controllers

import (
	"errors"
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
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

	page, err := utility.ParseQueryParamToInt(ctx, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(ctx, errors.New("invalid page parameter"))
		return
	}
	pageSize, err := utility.ParseQueryParamToInt(ctx, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(ctx, errors.New("invalid page_size parameter"))
		return
	}
	if page < 1 || pageSize < 1 {
		utility.SendInvalidErrResp(ctx, errors.New("page and page_size must be greater than zero"))
		return
	}

	validSortFields := map[string]string{
		"id":                          "sd.id",
		"reference_number":            "sd.reference_number",
		"entity_name":                 "sd.entity_name",
		"created_at":                  "sd.created_at",
		"updated_at":                  "sd.updated_at",
		"data_entity_type_title_en":   "data_entity_types.title_en",
		"data_entity_type_title_ar":   "data_entity_types.title_ar",
		"responsible_entity_title_en": "responsible_entity.title_en",
		"responsible_entity_title_ar": "responsible_entity.title_ar",
		"beneficiary_entity_title_en": "beneficiary_entity.title_en",
		"beneficiary_entity_title_ar": "beneficiary_entity.title_ar",
	}
	if sortBy == "" {
		sortBy = "id"
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(ctx, errors.New("invalid sort_by parameter"))
		return
	}
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

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

	entries, err := csdc.service.GetAll(&socDtos.CSLogSourceDirectoryListQueryParamsDTO{
		Sort:                sort,
		SortBy:              sortField,
		Keyword:             keyword,
		Page:                page,
		PageSize:            pageSize,
		DataEntityTypeID:    dataEntityTypeID,
		ResponsibleEntityID: responsibleEntityID,
		BeneficiaryEntityID: beneficiaryEntityID,
	})
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

	entry, err := csdc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		var notFoundErr *utility.ResourceNotFoundError
		if errors.As(err, &notFoundErr) {
			utility.SendNotFoundErrResp(ctx)
		} else {
			utility.SendInternalErrResp(csdc.db, ctx, err.Error())
		}
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

	createdID, err := csdc.service.Create(currentUserID, form)
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

	ctx.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (csdc *CSLogSourceDirectoryController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
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

	err = csdc.service.Update(uint(id), currentUserID, form)
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

	ctx.Status(http.StatusOK)
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
