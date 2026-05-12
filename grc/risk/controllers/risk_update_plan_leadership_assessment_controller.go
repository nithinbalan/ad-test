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

type RiskUpdatePlanLeadershipAssessmentController struct {
	service services.RiskUpdatePlanLeadershipAssessmentServiceInterface
	db      *gorm.DB
}

func NewRiskUpdatePlanLeadershipAssessmentController(service services.RiskUpdatePlanLeadershipAssessmentServiceInterface) *RiskUpdatePlanLeadershipAssessmentController {
	return &RiskUpdatePlanLeadershipAssessmentController{
		service: service,
	}
}

func (rrupc *RiskUpdatePlanLeadershipAssessmentController) GetLeadershipAssessmentUsers(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "first_name_en")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"first_name_en":  "users.first_name_en",
		"first_name_ar":  "users.first_name_ar",
		"email":          "users.email",
		"mobile_number":  "users.mobile_number",
		"designation_en": "designations.title_en",
		"designation_ar": "designations.title_ar",
		"user_status_en": "user_statuses.title_en",
		"user_status_ar": "user_statuses.title_ar",
		"name_en":        "CONCAT(users.first_name_en, ' ', users.last_name_en)",
		"name_ar":        "CONCAT(users.first_name_ar, ' ', users.last_name_ar)",
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

	queryParams := riskDTOs.CommonListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	users, err := rrupc.service.GetLeadershipAssessmentUsers(uint(planID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

func (rrupc *RiskUpdatePlanLeadershipAssessmentController) CreateLeadershipAssessment(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	err = rrupc.service.CreateLeadershipAssessment(uint(planID), currentUserID)
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

func (rrupc *RiskUpdatePlanLeadershipAssessmentController) MarkLeadershipAssessmentAsNotApplicable(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	var formDTO riskDTOs.MarkTimelineNotApplicableDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = rrupc.service.MarkLeadershipAssessmentAsNotApplicable(uint(planID), currentUserID, &formDTO)
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

func (rrupc *RiskUpdatePlanLeadershipAssessmentController) GetLeadershipAssessmentDocuments(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documents, err := rrupc.service.GetLeadershipAssessmentDocuments(uint(planID))
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, documents)
}

func (rrupc *RiskUpdatePlanLeadershipAssessmentController) CreateLeadershipAssessmentDocuments(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	var documentsWrapper riskDTOs.CommonDocumentFormDTO
	if err := c.ShouldBindJSON(&documentsWrapper); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = rrupc.service.CreateLeadershipAssessmentDocuments(uint(planID), currentUserID, documentsWrapper.Documents)
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

func (rrupc *RiskUpdatePlanLeadershipAssessmentController) GetLeadershipAssessmentDocumentPreview(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendInvalidErrResp(c, errors.New("token is required"))
		return
	}

	filePath, err := rrupc.service.GetLeadershipAssessmentDocumentPath(uint(planID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, err.Error())
			return
		}
	}

	c.File(filePath)
}

func (rrupc *RiskUpdatePlanLeadershipAssessmentController) DownloadLeadershipAssessmentDocument(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendInvalidErrResp(c, errors.New("token is required"))
		return
	}

	filePath, err := rrupc.service.GetLeadershipAssessmentDocumentPath(uint(planID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, err.Error())
			return
		}
	}

	filename := filepath.Base(filePath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(filePath)
}
