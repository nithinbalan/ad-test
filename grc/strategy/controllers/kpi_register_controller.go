package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KPIRegisterController struct {
	service services.KPIRegisterServiceInterface
	db      *gorm.DB
}

func NewKPIRegisterController(service services.KPIRegisterServiceInterface) *KPIRegisterController {
	return &KPIRegisterController{
		service: service,
	}
}

func (spc *KPIRegisterController) GetAllKPIRegister(c *gin.Context) {

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"title",
	}
	validSortFields := map[string]string{
		"title": "kpi_register.title",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := strategyDTOs.KPIRegisterListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}
	kpiRegister, err := spc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, kpiRegister)
}

func (spc *KPIRegisterController) GetKPIRegisterDetails(c *gin.Context) {
	idStr := c.Param("kpiRegisterId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kpiRegister, err := spc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(spc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, kpiRegister)
}

func (spc *KPIRegisterController) CreateKPIRegister(c *gin.Context) {

	var kpiRegisterFormDTO strategyDTOs.KPIRegisterFormDTO

	if err := c.ShouldBindJSON(&kpiRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	id, err := spc.service.Create(currentUserID, &kpiRegisterFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(spc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (spc *KPIRegisterController) UpdateKPIRegister(c *gin.Context) {
	var kpiRegisterFormDTO strategyDTOs.KPIRegisterFormDTO

	kpiRegisterID, err := strconv.Atoi(c.Param("kpiRegisterId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&kpiRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	err = spc.service.Update(uint(kpiRegisterID), currentUserID, &kpiRegisterFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(spc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (spc *KPIRegisterController) DeleteKPIRegister(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	failedIDs, err := spc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
