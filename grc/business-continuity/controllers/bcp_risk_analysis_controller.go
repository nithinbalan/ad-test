package controllers

import (
	"jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BcpRiskAnalysisController struct {
	service services.BcpRiskAnalysisServiceInterface
	db      *gorm.DB
}

func NewBcpRiskAnalysisController(service services.BcpRiskAnalysisServiceInterface) *BcpRiskAnalysisController {
	return &BcpRiskAnalysisController{
		service: service,
	}
}

func (brac *BcpRiskAnalysisController) GetAllRisksByBcpId(c *gin.Context) {
	idStr := c.Param("bcpId")

	bcpId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	valueStr := c.Query("value")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(brac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(brac.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(brac.db, c, "Invalid value parameter")
			return
		}
	}

	queryParams := dtos.BcpRiskListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
		Value:       value,
	}

	risks, err := brac.service.GetAllRisksByBcpId(uint(bcpId), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(brac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, risks)
}

func (brac *BcpRiskAnalysisController) AddRiskToBcp(c *gin.Context) {

	idStr := c.Param("bcpId")

	bcpId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var bcpRiskForm dtos.BcpRiskForm
	if err := c.ShouldBindJSON(&bcpRiskForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(brac.db, c, err.Error())
		return
	}

	err = brac.service.AddRiskToBcp(uint(bcpId), bcpRiskForm, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(brac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Risk added to BCP successfully"})
}

func (brac *BcpRiskAnalysisController) DeleteRiskFromBcp(c *gin.Context) {
	idStr := c.Param("bcpId")

	bcpId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var bcpRiskForm dtos.BcpRiskForm
	if err := c.ShouldBindJSON(&bcpRiskForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = brac.service.DeleteRiskFromBcp(uint(bcpId), bcpRiskForm)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(brac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Risk deleted from BCP successfully"})
}

func (brac *BcpRiskAnalysisController) GetAllProcessRisks(c *gin.Context) {
	idStr := c.Param("bcpId")

	bcpId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	valueStr := c.Query("value")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(brac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(brac.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(brac.db, c, "Invalid value parameter")
			return
		}
	}

	queryParams := dtos.BcpRiskListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
		Value:       value,
	}

	risks, err := brac.service.GetAllProcessRisks(uint(bcpId), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(brac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, risks)
}
