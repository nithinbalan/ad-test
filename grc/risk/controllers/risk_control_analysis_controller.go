package controllers

import (
	riskDTOs "jethur/grc/risk/dtos"

	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskControlAnalysisController struct {
	service services.RiskControlAnalysisServiceInterface
	db      *gorm.DB
}

func NewRiskControlAnalysisController(service services.RiskControlAnalysisServiceInterface) *RiskControlAnalysisController {
	return &RiskControlAnalysisController{
		service: service,
	}
}

func (rcac *RiskControlAnalysisController) GetAllRiskControlAnalyses(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	queryParams := riskDTOs.RiskControlAnalysisListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
	}
	riskControlAnalyses, err := rcac.service.GetAll(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rcac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskControlAnalyses)
}

func (rcac *RiskControlAnalysisController) CreateRiskControlAnalysis(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcac.db, c, err.Error())
		return
	}

	err = rcac.service.Create(uint(riskID), currentUserID)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rcac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rcac *RiskControlAnalysisController) UpdateRiskControlAnalysis(c *gin.Context) {
	var controlTestFormDTO organizationDTOs.ControlTestFormDTO

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&controlTestFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcac.db, c, err.Error())
		return
	}

	err = rcac.service.Update(uint(riskID), uint(id), currentUserID, &controlTestFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rcac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rcac *RiskControlAnalysisController) DeleteRiskControlAnalyses(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcac.db, c, err.Error())
		return
	}

	failedIDs, err := rcac.service.Delete(uint(riskID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rcac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
