package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditFindingReassessmentController struct {
	service services.AuditFindingReassessmentServiceInterface
	db      *gorm.DB
}

func NewAuditFindingReassessmentController(service services.AuditFindingReassessmentServiceInterface) *AuditFindingReassessmentController {
	return &AuditFindingReassessmentController{
		service: service,
	}
}

func (afrc *AuditFindingReassessmentController) GetAllAuditFindingReassessments(c *gin.Context) {
	findingID, err := strconv.ParseUint(c.Param("findingId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := auditDTOs.AuditFindingReassessmentListQueryParamsDTO{
		Sort:     c.Query("sort_order"),
		SortBy:   c.Query("sort_by"),
		Keyword:  c.Query("keyword"),
		Page:     page,
		PageSize: pageSize,
	}

	reassessments, err := afrc.service.GetAll(uint(findingID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, reassessments)
}

func (afrc *AuditFindingReassessmentController) GetAllAuditFindingReassessmentActionPlans(c *gin.Context) {
	findingID, err := strconv.ParseUint(c.Param("findingId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reassessmentID, err := strconv.ParseUint(c.Param("reassessmentId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	sortOrder := c.DefaultQuery("sort_order", "asc")
	if sortOrder != "asc" && sortOrder != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}

	validSortFields := map[string]string{
		"id":                             "audit_finding_action_plans.id",
		"audit_finding_id":               "audit_finding_action_plans.audit_finding_id",
		"audit_finding_reference_number": "audit_findings.reference_number",
		"audit_finding_title":            "audit_findings.title",
		"audit_finding_description":      "audit_findings.description",
		"title":                          "audit_finding_action_plans.title",
		"description":                    "audit_finding_action_plans.description",
		"severity_level_title_en":        "severity_levels.title_en",
		"severity_level_title_ar":        "severity_levels.title_ar",
		"responsible_user_id":            "responsible_users.id",
		"responsible_user_first_name_en": "responsible_users.first_name_en",
		"responsible_user_first_name_ar": "responsible_users.first_name_ar",
		"responsible_user_last_name_en":  "responsible_users.last_name_en",
		"responsible_user_last_name_ar":  "responsible_users.last_name_ar",
		"responsible_user_email":         "responsible_users.email",
		"review_user_id":                 "review_users.id",
		"review_user_first_name_en":      "review_users.first_name_en",
		"review_user_first_name_ar":      "review_users.first_name_ar",
		"review_user_last_name_en":       "review_users.last_name_en",
		"review_user_last_name_ar":       "review_users.last_name_ar",
		"review_user_email":              "review_users.email",
		"target_date":                    "audit_finding_action_plans.target_date",
		"percentage":                     "audit_finding_action_plans.percentage",
		"comments":                       "audit_finding_action_plans.comments",
		"action_plan_status_title_en":    "action_plan_statuses.title_en",
		"action_plan_status_title_ar":    "action_plan_statuses.title_ar",
	}

	sortField, valid := validSortFields[c.DefaultQuery("sort_by", "title")]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	queryParams := auditDTOs.AuditFindingActionPlanListQueryParamsDTO{
		Sort:        sortOrder,
		SortBy:      sortField,
		Keyword:     c.Query("keyword"),
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: c.QueryArray("quick_filter"),
	}

	actionPlans, err := afrc.service.GetActionPlans(uint(findingID), uint(reassessmentID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, actionPlans)
}

func (afrc *AuditFindingReassessmentController) CreateAuditFindingReassessmentActionPlan(c *gin.Context) {
	findingID, err := strconv.ParseUint(c.Param("findingId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reassessmentID, err := strconv.ParseUint(c.Param("reassessmentId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO auditDTOs.AuditFindingActionPlanFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afrc.db, c, err.Error())
		return
	}

	id, err := afrc.service.CreateActionPlan(uint(findingID), uint(reassessmentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (afrc *AuditFindingReassessmentController) GetAuditFindingReassessmentDetails(c *gin.Context) {
	findingID, err := strconv.ParseUint(c.Param("findingId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reassessmentID, err := strconv.ParseUint(c.Param("reassessmentId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reassessment, err := afrc.service.GetFullDetailsByID(uint(findingID), uint(reassessmentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, reassessment)
}

func (afrc *AuditFindingReassessmentController) CreateAuditFindingReassessment(c *gin.Context) {
	findingID, err := strconv.ParseUint(c.Param("findingId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO auditDTOs.AuditFindingReassessmentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afrc.db, c, err.Error())
		return
	}

	id, err := afrc.service.Create(uint(findingID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (afrc *AuditFindingReassessmentController) UpdateAuditFindingReassessment(c *gin.Context) {
	findingID, err := strconv.ParseUint(c.Param("findingId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reassessmentID, err := strconv.ParseUint(c.Param("reassessmentId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO auditDTOs.AuditFindingReassessmentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afrc.db, c, err.Error())
		return
	}

	err = afrc.service.Update(uint(findingID), uint(reassessmentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (afrc *AuditFindingReassessmentController) DeleteAuditFindingReassessment(c *gin.Context) {
	findingID, err := strconv.ParseUint(c.Param("findingId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reassessmentID, err := strconv.ParseUint(c.Param("reassessmentId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afrc.db, c, err.Error())
		return
	}

	err = afrc.service.Delete(uint(findingID), uint(reassessmentID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (afrc *AuditFindingReassessmentController) GetAuditFindingReassessmentDocumentPreview(c *gin.Context) {
	findingID, err := strconv.ParseUint(c.Param("findingId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reassessmentID, err := strconv.ParseUint(c.Param("reassessmentId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := afrc.service.GetDocumentPath(uint(findingID), uint(reassessmentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(afrc.db, c, err.Error())
		}
		return
	}

	c.File(documentPath)
}

func (afrc *AuditFindingReassessmentController) DownloadAuditFindingReassessmentDocument(c *gin.Context) {
	findingID, err := strconv.ParseUint(c.Param("findingId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reassessmentID, err := strconv.ParseUint(c.Param("reassessmentId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := afrc.service.GetDocumentPath(uint(findingID), uint(reassessmentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(afrc.db, c, err.Error())
		}
		return
	}

	filename := filepath.Base(documentPath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(documentPath)
}
