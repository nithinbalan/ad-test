package controllers

import (
	"errors"
	"net/http"
	"strconv"

	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskUpdatePlanRiskAnalysisController struct {
	service services.RiskUpdatePlanRiskAnalysisServiceInterface
	db      *gorm.DB
}

func NewRiskUpdatePlanRiskAnalysisController(service services.RiskUpdatePlanRiskAnalysisServiceInterface) *RiskUpdatePlanRiskAnalysisController {
	return &RiskUpdatePlanRiskAnalysisController{
		service: service,
	}
}

func (rrupc *RiskUpdatePlanRiskAnalysisController) GetRiskAnalysisList(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"reference_number":           "combined_risks.reference_number",
		"title":                      "combined_risks.title",
		"description":                "combined_risks.description",
		"inherent_likelihood_en":     "combined_risks.inherent_likelihood_title_en",
		"inherent_likelihood_ar":     "combined_risks.inherent_likelihood_title_ar",
		"inherent_impact_en":         "combined_risks.inherent_impact_title_en",
		"inherent_impact_ar":         "combined_risks.inherent_impact_title_ar",
		"inherent_risk_rating_en":    "combined_risks.inherent_risk_rating_title_en",
		"inherent_risk_rating_ar":    "combined_risks.inherent_risk_rating_title_ar",
		"inherent_risk_rating_value": "combined_risks.inherent_risk_rating_value",
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	risks, err := rrupc.service.GetRiskAnalysisList(uint(planID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, risks)
}

func (rrupc *RiskUpdatePlanRiskAnalysisController) CreateRiskAnalysis(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var formDTO riskDTOs.RiskAnalysisFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	err = rrupc.service.CreateRiskAnalysis(uint(planID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, err)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrupc *RiskUpdatePlanRiskAnalysisController) MarkRiskAnalysisAsNotApplicable(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO riskDTOs.MarkTimelineNotApplicableDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	err = rrupc.service.MarkRiskAnalysisAsNotApplicable(uint(planID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
