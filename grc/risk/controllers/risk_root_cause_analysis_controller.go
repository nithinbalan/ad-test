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

type RiskRootCauseAnalysisController struct {
	service services.RiskRootCauseAnalysisServiceInterface
	db      *gorm.DB
}

func NewRiskRootCauseAnalysisController(service services.RiskRootCauseAnalysisServiceInterface) *RiskRootCauseAnalysisController {
	return &RiskRootCauseAnalysisController{
		service: service,
	}
}

func (rrcac *RiskRootCauseAnalysisController) GetAllRiskRootCauseAnalyses(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	queryParams := riskDTOs.RiskRootCauseAnalysisListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
	}
	riskRootCauseAnalyses, err := rrcac.service.GetAll(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrcac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskRootCauseAnalyses)
}

func (rrcac *RiskRootCauseAnalysisController) CreateRiskRootCauseAnalysis(c *gin.Context) {
	var RiskRootCauseAnalysisFormDTO riskDTOs.RiskRootCauseAnalysisFormDTO

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&RiskRootCauseAnalysisFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrcac.db, c, err.Error())
		return
	}

	id, err := rrcac.service.Create(uint(riskID), currentUserID, &RiskRootCauseAnalysisFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrcac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (rrcac *RiskRootCauseAnalysisController) UpdateRiskRootCauseAnalysis(c *gin.Context) {
	var RiskRootCauseAnalysisFormDTO riskDTOs.RiskRootCauseAnalysisFormDTO

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

	if err := c.ShouldBindJSON(&RiskRootCauseAnalysisFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrcac.db, c, err.Error())
		return
	}

	err = rrcac.service.Update(uint(riskID), uint(id), currentUserID, &RiskRootCauseAnalysisFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrcac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrcac *RiskRootCauseAnalysisController) DeleteRiskRootCauseAnalysis(c *gin.Context) {
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

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrcac.db, c, err.Error())
		return
	}

	failedID, err := rrcac.service.Delete(uint(riskID), uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrcac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
