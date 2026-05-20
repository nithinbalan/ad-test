package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskUpdatePlanERMDeptController struct {
	service services.RiskUpdatePlanERMDeptServiceInterface
	db      *gorm.DB
}

func NewRiskUpdatePlanERMDeptController(service services.RiskUpdatePlanERMDeptServiceInterface) *RiskUpdatePlanERMDeptController {
	return &RiskUpdatePlanERMDeptController{
		service: service,
	}
}

func (rruedc *RiskUpdatePlanERMDeptController) GetERMDepartmentEvaluationRisks(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Parse query parameters
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"reference_number":          "rr.reference_number",
		"title":                     "rr.title",
		"description":               "rr.description",
		"risk_category_en":          "rc.title_en",
		"risk_category_ar":          "rc.title_ar",
		"organization_unit_en":      "ru.title_en",
		"organization_unit_ar":      "ru.title_ar",
		"responsible_user_en":       "u.first_name_en",
		"responsible_user_ar":       "u.first_name_ar",
		"risk_source_en":            "rs.title_en",
		"risk_source_ar":            "rs.title_ar",
		"risk_status_en":            "rst.title_en",
		"risk_status_ar":            "rst.title_ar",
		"inherent_likelihood_en":    "ill.title_en",
		"inherent_likelihood_ar":    "ill.title_ar",
		"inherent_impact_en":        "iil.title_en",
		"inherent_impact_ar":        "iil.title_ar",
		"inherent_risk_rating_en":   "irr.title_en",
		"inherent_risk_rating_ar":   "irr.title_ar",
		"risk_response_strategy_en": "rrs.title_en",
		"risk_response_strategy_ar": "rrs.title_ar",
		"residual_likelihood_en":    "rll.title_en",
		"residual_likelihood_ar":    "rll.title_ar",
		"residual_impact_en":        "ril.title_en",
		"residual_impact_ar":        "ril.title_ar",
		"residual_risk_rating_en":   "rrr.title_en",
		"residual_risk_rating_ar":   "rrr.title_ar",
		"reported_at":               "rr.reported_at",
		"review_frequency_en":       "rf.title_en",
		"review_frequency_ar":       "rf.title_ar",
		"next_review_date":          "rr.next_review_date",
		"division_en":               "div.title_en",
		"division_ar":               "div.title_ar",
		"section_en":                "sec.title_en",
		"section_ar":                "sec.title_ar",
		"sub_section_en":            "ss.title_en",
		"sub_section_ar":            "ss.title_ar",
		"risk_code":                 "rr.risk_code",
		"registration_date":         "rr.registration_date",
		"is_strategic":              "rr.is_strategic",
		"is_operational":            "rr.is_operational",
		"expected_cause":            "rr.expected_cause",
		"consequence":               "rr.consequence",
		"is_internal":               "rr.is_internal",
		"is_external":               "rr.is_external",
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rruedc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rruedc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := riskDTOs.CommonListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	mappedRisks, err := rruedc.service.GetERMDeptEvaluationRisks(uint(id), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rruedc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, mappedRisks)
}

func (rruedc *RiskUpdatePlanERMDeptController) GetERMDepartmentEvaluationProcesses(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Parse query parameters
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"reference_number":     "p.reference_number",
		"title_en":             "p.title_en",
		"title_ar":             "p.title_ar",
		"description":          "p.description",
		"process_group_en":     "pg.title_en",
		"process_group_ar":     "pg.title_ar",
		"organization_unit_en": "ou.title_en",
		"organization_unit_ar": "ou.title_ar",
		"owner_en":             "u.first_name_en",
		"owner_ar":             "u.first_name_ar",
		"status_en":            "ls.title_en",
		"status_ar":            "ls.title_ar",
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rruedc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rruedc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := riskDTOs.CommonListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	mappedProcesses, err := rruedc.service.GetERMDeptEvaluationProcesses(uint(id), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rruedc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, mappedProcesses)
}

func (rruedc *RiskUpdatePlanERMDeptController) GetERMDepartmentEvaluationControls(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Parse query parameters
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"reference_number":         "c.reference_number",
		"code":                     "c.code",
		"title_en":                 "c.title_en",
		"title_ar":                 "c.title_ar",
		"description":              "c.description",
		"objective":                "c.objective",
		"control_category_en":      "cc.title_en",
		"control_category_ar":      "cc.title_ar",
		"control_type_en":          "ct.title_en",
		"control_type_ar":          "ct.title_ar",
		"owner_en":                 "owner.first_name_en",
		"owner_ar":                 "owner.first_name_ar",
		"review_frequency_en":      "rf.title_en",
		"review_frequency_ar":      "rf.title_ar",
		"implementation_status_en": "cis.title_en",
		"implementation_status_ar": "cis.title_ar",
		"implementation_date":      "c.implemented_at",
		"efficiency_status_en":     "ces.title_en",
		"efficiency_status_ar":     "ces.title_ar",
		"last_test_date":           "c.last_test_date",
		"next_test_date":           "c.next_test_date",
		"test_frequency_en":        "tf.title_en",
		"test_frequency_ar":        "tf.title_ar",
		"last_review_note":         "c.last_review_note",
		"next_review_date":         "c.next_review_date",
		"status_en":                "ls.title_en",
		"status_ar":                "ls.title_ar",
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rruedc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rruedc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := riskDTOs.CommonListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	mappedControls, err := rruedc.service.GetERMDeptEvaluationControls(uint(id), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rruedc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, mappedControls)
}
func (rruedc *RiskUpdatePlanERMDeptController) CreateERMDepartmentEvaluation(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rruedc.db, c, err.Error())
		return
	}

	err = rruedc.service.CreateERMDepartmentEvaluation(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rruedc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rruedc *RiskUpdatePlanERMDeptController) UpdateERMDepartmentEvaluationRisks(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO riskDTOs.MapRisksToUpdatePlanDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rruedc.db, c, err.Error())
		return
	}

	err = rruedc.service.MapRisksToTimeline(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rruedc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rruedc *RiskUpdatePlanERMDeptController) UpdateERMDepartmentEvaluationProcesses(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO riskDTOs.MapProcessesToUpdatePlanDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rruedc.db, c, err.Error())
		return
	}

	err = rruedc.service.MapProcessesToTimeline(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rruedc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rruedc *RiskUpdatePlanERMDeptController) UpdateERMDepartmentEvaluationControls(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO riskDTOs.MapControlsToUpdatePlanDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rruedc.db, c, err.Error())
		return
	}

	err = rruedc.service.MapControlsToTimeline(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rruedc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rruedc *RiskUpdatePlanERMDeptController) MarkERMDeptEvaluationAsNotApplicable(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
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
		utility.SendInternalErrResp(rruedc.db, c, err.Error())
		return
	}

	err = rruedc.service.MarkERMDeptEvaluationAsNotApplicable(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rruedc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rruedc *RiskUpdatePlanERMDeptController) GetERMDeptEvaluationDocuments(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documents, err := rruedc.service.GetERMDeptEvaluationDocuments(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rruedc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, documents)
}

func (rruedc *RiskUpdatePlanERMDeptController) CreateERMDeptEvaluationDocuments(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO riskDTOs.CommonDocumentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rruedc.db, c, err.Error())
		return
	}

	err = rruedc.service.CreateERMDeptEvaluationDocuments(uint(id), currentUserID, formDTO.Documents)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rruedc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rruedc *RiskUpdatePlanERMDeptController) GetERMDeptEvaluationDocumentPreview(c *gin.Context) {
	idStr := c.Param("id")
	token := c.Query("token")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if token == "" {
		utility.SendInvalidErrResp(c, errors.New("token is required"))
		return
	}

	filePath, err := rruedc.service.GetERMDeptEvaluationDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rruedc.db, c, err.Error())
			return
		}
	}

	c.File(filePath)
}

func (rruedc *RiskUpdatePlanERMDeptController) DownloadERMDeptEvaluationDocument(c *gin.Context) {
	idStr := c.Param("id")
	token := c.Query("token")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if token == "" {
		utility.SendInvalidErrResp(c, errors.New("token is required"))
		return
	}

	filePath, err := rruedc.service.GetERMDeptEvaluationDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rruedc.db, c, err.Error())
			return
		}
	}

	c.FileAttachment(filePath, filepath.Base(filePath))
}
