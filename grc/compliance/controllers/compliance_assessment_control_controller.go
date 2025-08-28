package controllers

import (
	"errors"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceAssessmentControlController struct {
	service services.ComplianceAssessmentControlServiceInterface
	db      *gorm.DB
}

func NewComplianceAssessmentControlController(service services.ComplianceAssessmentControlServiceInterface) *ComplianceAssessmentControlController {
	return &ComplianceAssessmentControlController{
		service: service,
	}
}

// ====================== CLAUSES ======================
func (cacc *ComplianceAssessmentControlController) GetComplianceAssessmentControlClauses(ctx *gin.Context) {
	assessmentIDStr := ctx.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(cacc.db, ctx, err.Error())
		return
	}

	clauses, err := cacc.service.GetAllClauses(uint(assessmentID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
		default:
			utility.SendInternalErrResp(cacc.db, ctx, e.Error())
		}
		return
	}

	ctx.JSON(http.StatusOK, clauses)
}

// ====================== LIST ======================
func (cacc *ComplianceAssessmentControlController) GetComplianceAssessmentControls(ctx *gin.Context) {
	assessmentIDStr := ctx.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	sort := ctx.DefaultQuery("sort_order", "asc")
	sortBy := ctx.DefaultQuery("sort_by", "compliance_clause_no")
	keyword := ctx.Query("keyword")
	quickFilter := ctx.QueryArray("quick_filter")

	validSortFields := map[string]string{
		"compliance_clause_no":       "compliance_clauses.clause_no",
		"compliance_clause_title_en": "compliance_clauses.title_en",
		"compliance_clause_title_ar": "compliance_clauses.title_ar",
		"control_title_en":           "compliance_controls.title_en",
		"control_title_ar":           "compliance_controls.title_ar",
		"assessment_completed":       "compliance_assessment_controls.is_completed",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(ctx, errors.New("invalid 'sort_by' parameter"))
		return
	}

	allowedParams := []string{
		"compliance_control_ids",
		"compliance_clause_ids",
		"status_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, ctx.Request.URL.Query())

	page, _ := utility.ParseQueryParamToInt(ctx, "page", 1)
	pageSize, _ := utility.ParseQueryParamToInt(ctx, "page_size", 10)

	clauseIDStr := ctx.Query("clause_id")
	var clauseID uint
	if clauseIDStr != "" {
		parsedClauseID, _ := strconv.ParseUint(clauseIDStr, 10, 32)
		clauseID = uint(parsedClauseID)
	}

	queryParams := complianceDTOs.ComplianceAssessmentControlListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Filters:     filters,
		QuickFilter: quickFilter,
		Page:        page,
		PageSize:    pageSize,
		ClauseID:    clauseID,
	}

	controls, err := cacc.service.GetAll(uint(assessmentID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(cacc.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, controls)
}

// ====================== DETAILS ======================
func (cacc *ComplianceAssessmentControlController) GetComplianceAssessmentControlDetails(ctx *gin.Context) {
	assessmentIDStr := ctx.Param("assessmentId")
	assessmentID, _ := strconv.ParseUint(assessmentIDStr, 10, 64)

	controlIDStr := ctx.Param("controlId")
	controlID, _ := strconv.ParseUint(controlIDStr, 10, 64)

	control, err := cacc.service.GetFullDetailsByID(uint(assessmentID), uint(controlID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
		default:
			utility.SendInternalErrResp(cacc.db, ctx, e.Error())
		}
		return
	}

	ctx.JSON(http.StatusOK, control)
}

// ====================== UPDATE ======================
func (cacc *ComplianceAssessmentControlController) UpdateComplianceAssessmentControl(ctx *gin.Context) {
	var updateForm complianceDTOs.ComplianceAssessmentControlUpdateFormDTO

	assessmentID, _ := strconv.Atoi(ctx.Param("assessmentId"))
	controlID, _ := strconv.Atoi(ctx.Param("controlId"))

	if err := ctx.ShouldBindJSON(&updateForm); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(cacc.db, ctx, err.Error())
		return
	}

	err = cacc.service.Update(uint(assessmentID), uint(controlID), currentUserID, &updateForm)
	if err != nil {
		utility.SendInternalErrResp(cacc.db, ctx, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

// ====================== COMPLETE / OPEN ======================
func (cacc *ComplianceAssessmentControlController) MarkAsCompleted(ctx *gin.Context) {
	assessmentID, _ := strconv.Atoi(ctx.Param("assessmentId"))
	controlID, _ := strconv.Atoi(ctx.Param("controlId"))

	currentUserID, _ := utility.GetCurrentUserID(ctx)

	if err := cacc.service.MarkAsCompleted(uint(assessmentID), uint(controlID), currentUserID); err != nil {
		utility.SendInternalErrResp(cacc.db, ctx, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (cacc *ComplianceAssessmentControlController) MarkAsOpen(ctx *gin.Context) {
	assessmentID, _ := strconv.Atoi(ctx.Param("assessmentId"))
	controlID, _ := strconv.Atoi(ctx.Param("controlId"))

	currentUserID, _ := utility.GetCurrentUserID(ctx)

	if err := cacc.service.MarkAsOpen(uint(assessmentID), uint(controlID), currentUserID); err != nil {
		utility.SendInternalErrResp(cacc.db, ctx, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

// ====================== APPROVE / REVERT / RESUBMIT ======================
func (cacc *ComplianceAssessmentControlController) ApproveComplianceAssessmentControl(ctx *gin.Context) {
	var form complianceDTOs.ComplianceAssessmentControlCommentFormDTO

	assessmentID, _ := strconv.Atoi(ctx.Param("assessmentId"))
	controlID, _ := strconv.Atoi(ctx.Param("controlId"))

	if err := ctx.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}
	currentUserID, _ := utility.GetCurrentUserID(ctx)

	if err := cacc.service.ApproveComplianceControl(uint(assessmentID), uint(controlID), currentUserID, &form); err != nil {
		utility.SendInternalErrResp(cacc.db, ctx, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (cacc *ComplianceAssessmentControlController) RevertComplianceAssessmentControl(ctx *gin.Context) {
	var form complianceDTOs.ComplianceAssessmentControlCommentFormDTO

	assessmentID, _ := strconv.Atoi(ctx.Param("assessmentId"))
	controlID, _ := strconv.Atoi(ctx.Param("controlId"))

	if err := ctx.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}
	currentUserID, _ := utility.GetCurrentUserID(ctx)

	if err := cacc.service.RevertComplianceControl(uint(assessmentID), uint(controlID), currentUserID, &form); err != nil {
		utility.SendInternalErrResp(cacc.db, ctx, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (cacc *ComplianceAssessmentControlController) ResubmitComplianceAssessmentControl(ctx *gin.Context) {
	var form complianceDTOs.ComplianceAssessmentControlCommentFormDTO

	assessmentID, _ := strconv.Atoi(ctx.Param("assessmentId"))
	controlID, _ := strconv.Atoi(ctx.Param("controlId"))

	if err := ctx.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}
	currentUserID, _ := utility.GetCurrentUserID(ctx)

	if err := cacc.service.ResubmitComplianceControl(uint(assessmentID), uint(controlID), currentUserID, &form); err != nil {
		utility.SendInternalErrResp(cacc.db, ctx, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}
