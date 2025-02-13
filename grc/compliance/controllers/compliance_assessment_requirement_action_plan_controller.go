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

type ComplianceAssessmentRequirementActionPlanController struct {
	service services.ComplianceAssessmentRequirementActionPlanServiceInterface
	db      *gorm.DB
}

func NewComplianceAssessmentRequirementActionPlanController(service services.ComplianceAssessmentRequirementActionPlanServiceInterface) *ComplianceAssessmentRequirementActionPlanController {
	return &ComplianceAssessmentRequirementActionPlanController{
		service: service,
	}
}

func (carapc *ComplianceAssessmentRequirementActionPlanController) GetAllComplianceAssessmentRequirementActionPlans(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assessmentRequirementIDStr := c.Param("requirementId")
	assessmentRequirementID, err := strconv.ParseUint(assessmentRequirementIDStr, 10, 32)
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
		utility.SendInternalErrResp(carapc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"compliance_requirement_ids",
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
		utility.SendInternalErrResp(carapc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(carapc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := complianceDTOs.ComplianceAssessmentRequirementActionPlanListQueryParamsDTO{
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

	complianceActionPlans, err := carapc.service.GetAll(uint(assessmentID), uint(assessmentRequirementID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(carapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, complianceActionPlans)
}

func (carapc *ComplianceAssessmentRequirementActionPlanController) GetComplianceAssessmentRequirementActionPlanDetails(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assessmentRequirementIDStr := c.Param("requirementId")
	assessmentRequirementID, err := strconv.ParseUint(assessmentRequirementIDStr, 10, 32)
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

	complianceActionPlan, err := carapc.service.GetFullDetailsByID(uint(assessmentID), uint(assessmentRequirementID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(carapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, complianceActionPlan)
}

func (carapc *ComplianceAssessmentRequirementActionPlanController) CreateComplianceAssessmentRequirementActionPlan(c *gin.Context) {
	var ComplianceAssessmentRequirementActionPlanFormDTO complianceDTOs.ComplianceAssessmentRequirementActionPlanFormDTO

	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assessmentRequirementIDStr := c.Param("requirementId")
	assessmentRequirementID, err := strconv.ParseUint(assessmentRequirementIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&ComplianceAssessmentRequirementActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(carapc.db, c, err.Error())
		return
	}

	id, err := carapc.service.Create(uint(assessmentID), uint(assessmentRequirementID), currentUserID, &ComplianceAssessmentRequirementActionPlanFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(carapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (carapc *ComplianceAssessmentRequirementActionPlanController) UpdateComplianceAssessmentRequirementActionPlan(c *gin.Context) {
	var ComplianceAssessmentRequirementActionPlanFormDTO complianceDTOs.ComplianceAssessmentRequirementActionPlanFormDTO

	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assessmentRequirementIDStr := c.Param("requirementId")
	assessmentRequirementID, err := strconv.ParseUint(assessmentRequirementIDStr, 10, 32)
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

	if err := c.ShouldBindJSON(&ComplianceAssessmentRequirementActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(carapc.db, c, err.Error())
		return
	}

	err = carapc.service.Update(uint(assessmentID), uint(assessmentRequirementID), uint(id), currentUserID, &ComplianceAssessmentRequirementActionPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(carapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (carapc *ComplianceAssessmentRequirementActionPlanController) DeleteComplianceAssessmentRequirementActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assessmentRequirementIDStr := c.Param("requirementId")
	assessmentRequirementID, err := strconv.ParseUint(assessmentRequirementIDStr, 10, 32)
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
		utility.SendInternalErrResp(carapc.db, c, err.Error())
		return
	}

	failedIDs, err := carapc.service.Delete(uint(assessmentID), uint(assessmentRequirementID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(carapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
