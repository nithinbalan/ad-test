package controllers

import (
	slaDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SLARenewController struct {
	service services.SLARenewServiceInterface
	db      *gorm.DB
}

func NewSLARenewController(service services.SLARenewServiceInterface) *SLARenewController {
	return &SLARenewController{
		service: service,
	}
}

func (src *SLARenewController) GetAllSLARenews(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	// quickFilter := c.QueryArray("quick_filter")

	// Parse  issue_date_from, issue_date_to, expiry_date_from, expiry_date_to parameters
	issueDateFrom := c.Query("issue_date_from")
	issueDateTo := c.Query("issue_date_to")
	expiryDateFrom := c.Query("expiry_date_from")
	expiryDateTo := c.Query("expiry_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := slaDTOs.SLARenewListQueryParamsDTO{
		Sort:           sort,
		SortBy:         sortBy,
		Keyword:        keyword,
		Page:           page,
		PageSize:       pageSize,
		IssueDateFrom:  issueDateFrom,
		IssueDateTo:    issueDateTo,
		ExpiryDateFrom: expiryDateFrom,
		ExpiryDateTo:   expiryDateTo,
	}
	slaRenewHistories, err := src.service.GetAll(uint(id), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, slaRenewHistories)

}

func (src *SLARenewController) Renew(c *gin.Context) {
	docID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var slaRenewFormDTO slaDTOs.SLARenewFormDTO
	if err := c.ShouldBindJSON(&slaRenewFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}
	err = src.service.Renew(uint(docID), currentUserID, &slaRenewFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(src.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (src *SLARenewController) Revert(c *gin.Context) {
	docID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	err = src.service.Revert(uint(docID), currentUserID)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
