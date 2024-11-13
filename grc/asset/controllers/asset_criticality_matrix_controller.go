package controllers

import (
	"errors"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetCriticalityMatrixController struct {
	service services.AssetCriticalityMatrixServiceInterface
	db      *gorm.DB
}

func NewAssetCriticalityMatrixController(service services.AssetCriticalityMatrixServiceInterface) *AssetCriticalityMatrixController {
	return &AssetCriticalityMatrixController{
		service: service,
	}
}

func (acmc *AssetCriticalityMatrixController) GetAllMatrices(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page parameter"))
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page_size parameter"))
		return
	}

	params := dtos.AssetCriticalityMatrixListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	matrices, err := acmc.service.GetAll(&params)
	if err != nil {
		utility.SendInternalErrResp(acmc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, matrices)
}

func (acmc *AssetCriticalityMatrixController) GetMatrixDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	matrix, err := acmc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, matrix)
}

func (acmc *AssetCriticalityMatrixController) CreateMatrix(c *gin.Context) {
	var matrixForm dtos.AssetCriticalityMatrixFormDTO
	if err := c.ShouldBindJSON(&matrixForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acmc.db, c, err.Error())
		return
	}
	matrixID, err := acmc.service.Create(&matrixForm, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acmc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": matrixID})
}

func (acmc *AssetCriticalityMatrixController) UpdateMatrix(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var matrixForm dtos.AssetCriticalityMatrixFormDTO
	if err := c.ShouldBindJSON(&matrixForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acmc.db, c, err.Error())
		return
	}

	err = acmc.service.Update(uint(id), &matrixForm, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acmc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (acmc *AssetCriticalityMatrixController) DeleteMatrices(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acmc.db, c, err.Error())
		return
	}

	failedIDs, err := acmc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(acmc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
