package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KRIRegisterController struct {
	service services.KRIRegisterServiceInterface
	db      *gorm.DB
}

func NewKRIRegisterController(service services.KRIRegisterServiceInterface) *KRIRegisterController {
	return &KRIRegisterController{
		service: service,
	}
}

func (krc *KRIRegisterController) GetAllKRIRegisters(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"kri_category_ids",
		"risk_level_ids",
		"review_frequency_ids",
		"owner_ids",
		"owner_unit_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := riskDTOs.KRIRegisterListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}
	kriRegisters, err := krc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, kriRegisters)
}

func (krc *KRIRegisterController) GetKRIRegisterDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kriRegister, err := krc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(krc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, kriRegister)
}

func (krc *KRIRegisterController) CreateKRIRegister(c *gin.Context) {

	var kriRegisterFormDTO riskDTOs.KRIRegisterFormDTO
	if err := c.ShouldBindJSON(&kriRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}
	id, err := krc.service.Create(currentUserID, &kriRegisterFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(krc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (krc *KRIRegisterController) UpdateKRIRegister(c *gin.Context) {

	kriRegisterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var kriRegisterFormDTO riskDTOs.KRIRegisterFormDTO
	if err := c.ShouldBindJSON(&kriRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	err = krc.service.Update(uint(kriRegisterID), currentUserID, &kriRegisterFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(krc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (krc *KRIRegisterController) DeleteKRIRegisters(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	failedIDs, err := krc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
