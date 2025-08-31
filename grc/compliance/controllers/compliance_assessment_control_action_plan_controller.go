package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceAssessmentControlActionPlanController struct {
	service services.ComplianceAssessmentControlActionPlanServiceInterface
	db      *gorm.DB
}

func NewComplianceAssessmentControlActionPlanController(service services.ComplianceAssessmentControlActionPlanServiceInterface) *ComplianceAssessmentControlActionPlanController {
	return &ComplianceAssessmentControlActionPlanController{
		service: service,
	}
}

func (cacac *ComplianceAssessmentControlActionPlanController) GetAllComplianceAssessmentControlActionPlans(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assessmentControlIDStr := c.Param("controlId")
	assessmentControlID, err := strconv.ParseUint(assessmentControlIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(cacac.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"compliance_requirement_ids",
		"compliance_assessment_control_ids",
		"severity_level_ids",
		"responsible_user_ids",
		"action_plan_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse target_date_from, target_date_to
	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(cacac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(cacac.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := complianceDTOs.ComplianceActionPlanListQueryParamsDTO{
		Sort:           sort,
		SortBy:         sortBy,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		Filters:        filters,
		QuickFilter:    quickFilter,
		TargetDateFrom: targetDateFrom,
		TargetDateTo:   targetDateTo,
		DateFrom:       dateFrom,
		DateTo:         dateTo,
		Export:         export,
	}

	complianceActionPlans, err := cacac.service.GetAll(uint(assessmentID), uint(assessmentControlID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cacac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, complianceActionPlans)
}

func (cacac *ComplianceAssessmentControlActionPlanController) GetComplianceAssessmentControlActionPlanDetails(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assessmentControlIDStr := c.Param("controlId")
	assessmentControlID, err := strconv.ParseUint(assessmentControlIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceActionPlan, err := cacac.service.GetFullDetailsByID(uint(assessmentID), uint(assessmentControlID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cacac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, complianceActionPlan)
}

func (cacac *ComplianceAssessmentControlActionPlanController) CreateComplianceAssessmentControlActionPlan(c *gin.Context) {
	var ComplianceActionPlanFormDTO complianceDTOs.ComplianceActionPlanFormDTO

	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assessmentControlIDStr := c.Param("controlId")
	assessmentControlID, err := strconv.ParseUint(assessmentControlIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&ComplianceActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cacac.db, c, err.Error())
		return
	}

	id, err := cacac.service.Create(uint(assessmentID), uint(assessmentControlID), currentUserID, &ComplianceActionPlanFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cacac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (cacac *ComplianceAssessmentControlActionPlanController) UpdateComplianceAssessmentControlActionPlan(c *gin.Context) {
	var ComplianceActionPlanFormDTO complianceDTOs.ComplianceActionPlanFormDTO

	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assessmentControlIDStr := c.Param("controlId")
	assessmentControlID, err := strconv.ParseUint(assessmentControlIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&ComplianceActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cacac.db, c, err.Error())
		return
	}

	err = cacac.service.Update(uint(assessmentID), uint(assessmentControlID), uint(id), currentUserID, &ComplianceActionPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cacac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (cacac *ComplianceAssessmentControlActionPlanController) DeleteComplianceAssessmentControlActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assessmentControlIDStr := c.Param("controlId")
	assessmentControlID, err := strconv.ParseUint(assessmentControlIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cacac.db, c, err.Error())
		return
	}

	failedIDs, err := cacac.service.Delete(uint(assessmentID), uint(assessmentControlID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cacac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
