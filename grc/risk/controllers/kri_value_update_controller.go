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

type KRIValueUpdateController struct {
	service services.KRIValueUpdateServiceInterface
	db      *gorm.DB
}

func NewKRIValueUpdateController(service services.KRIValueUpdateServiceInterface) *KRIValueUpdateController {
	return &KRIValueUpdateController{
		service: service,
	}
}

func (kvuc *KRIValueUpdateController) GetAllKRIValueUpdates(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	reviewDateFrom := c.Query("review_date_from")
	reviewDateTO := c.Query("review_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := riskDTOs.KRIValueUpdateListQueryParamsDTO{
		Sort:           sort,
		SortBy:         sortBy,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		ReviewDateFrom: reviewDateFrom,
		ReviewDateTo:   reviewDateTO,
	}
	kriRegisterUpdates, err := kvuc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, kriRegisterUpdates)
}

func (kvuc *KRIValueUpdateController) CreateKRIValueUpdate(c *gin.Context) {

	var kriValueUpdateFormDTO riskDTOs.KRIValueUpdateFormDTO
	if err := c.ShouldBindJSON(&kriValueUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, err.Error())
		return
	}
	id, err := kvuc.service.Create(currentUserID, &kriValueUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kvuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (kvuc *KRIValueUpdateController) UpdateKRIValueUpdate(c *gin.Context) {

	kriValueUpdateID, err := strconv.Atoi(c.Param("updateId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var kriValueUpdateFormDTO riskDTOs.KRIValueUpdateFormDTO
	if err := c.ShouldBindJSON(&kriValueUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, err.Error())
		return
	}

	err = kvuc.service.Update(uint(kriValueUpdateID), currentUserID, &kriValueUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kvuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (kvuc *KRIValueUpdateController) DeleteKRIValueUpdates(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, err.Error())
		return
	}

	failedIDs, err := kvuc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
