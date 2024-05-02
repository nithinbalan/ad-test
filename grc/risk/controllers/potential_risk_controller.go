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

type PotentialRiskController struct {
	service services.PotentialRiskServiceInterface
	db      *gorm.DB
}

func NewPotentialRiskController(service services.PotentialRiskServiceInterface) *PotentialRiskController {
	return &PotentialRiskController{
		service: service,
	}
}

func (prc *PotentialRiskController) GetAllPotentialRisks(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"risk_category_ids",
		"severity_level_ids",
		"potential_risk_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := riskDTOs.PotentialRiskListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}
	potentialRisks, err := prc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, potentialRisks)
}

func (prc *PotentialRiskController) GetPotentialRiskDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	potentialRisk, err := prc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(prc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, potentialRisk)
}

func (prc *PotentialRiskController) CreatePotentialRisk(c *gin.Context) {

	var potentialRiskFormDTO riskDTOs.PotentialRiskFormDTO

	if err := c.ShouldBindJSON(&potentialRiskFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	id, err := prc.service.Create(currentUserID, &potentialRiskFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(prc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (prc *PotentialRiskController) UpdatePotentialRisk(c *gin.Context) {

	potentialRiskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var potentialRiskFormDTO riskDTOs.PotentialRiskFormDTO
	if err := c.ShouldBindJSON(&potentialRiskFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	err = prc.service.Update(uint(potentialRiskID), currentUserID, &potentialRiskFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(prc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (prc *PotentialRiskController) DeletePotentialRisks(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	failedIDs, err := prc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
