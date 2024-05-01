package controllers

import (
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskMatrixRiskRatingController struct {
	service services.RiskMatrixRiskRatingServiceInterface
	db      *gorm.DB
}

func NewRiskMatrixRiskRatingController(service services.RiskMatrixRiskRatingServiceInterface) *RiskMatrixRiskRatingController {
	return &RiskMatrixRiskRatingController{
		service: service,
	}
}

func (rmrrc *RiskMatrixRiskRatingController) GetAllRiskMatrixRiskRatings(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"risk_matrix_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := riskDTOs.RiskMatrixRiskRatingListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}

	riskRatings, err := rmrrc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(rmrrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, riskRatings)
}

func (rmrrc *RiskMatrixRiskRatingController) GetRiskMatrixRiskRatingDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskRiskRating, err := rmrrc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, riskRiskRating)
}

func (rmrrc *RiskMatrixRiskRatingController) CreateRiskMatrixRiskRating(c *gin.Context) {
	var riskRating *riskDTOs.RiskMatrixRiskRatingFormDTO

	if err := c.ShouldBindJSON(&riskRating); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmrrc.db, c, err.Error())
		return
	}

	id, err := rmrrc.service.Create(currentUserID, riskRating)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rmrrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, id)
}

func (rmrrc *RiskMatrixRiskRatingController) UpdateRiskMatrixRiskRating(c *gin.Context) {
	riskRatingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedRiskRating *riskDTOs.RiskMatrixRiskRatingFormDTO
	if err := c.ShouldBindJSON(&updatedRiskRating); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmrrc.db, c, err.Error())
		return
	}

	err = rmrrc.service.Update(uint(riskRatingID), currentUserID, updatedRiskRating)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rmrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rmrrc *RiskMatrixRiskRatingController) DeleteRiskMatrixRiskRatings(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmrrc.db, c, err.Error())
		return
	}

	failedIDs, err := rmrrc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rmrrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (rmrrc *RiskMatrixRiskRatingController) ActivateRiskMatrixRiskRating(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmrrc.db, c, err.Error())
		return
	}

	err = rmrrc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rmrrc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (rmrrc *RiskMatrixRiskRatingController) DeactivateRiskMatrixRiskRating(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmrrc.db, c, err.Error())
		return
	}

	err = rmrrc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rmrrc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
