package controllers

import (
	"net/http"
	"strconv"

	complianceDTOs "jethur/grc/compliance/dtos"
	complianceServices "jethur/grc/compliance/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceAssessmentAcknowledgeDepartmentController struct {
	db      *gorm.DB
	service complianceServices.ComplianceAssessmentAcknowledgeDepartmentServiceInterface
}

func NewComplianceAssessmentAcknowledgeDepartmentController(
	service complianceServices.ComplianceAssessmentAcknowledgeDepartmentServiceInterface,
) *ComplianceAssessmentAcknowledgeDepartmentController {
	return &ComplianceAssessmentAcknowledgeDepartmentController{
		service: service,
	}
}

func (caadc *ComplianceAssessmentAcknowledgeDepartmentController) GetAllAcknowledgements(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort", "ASC")
	sortBy := c.DefaultQuery("sort_by", "users.first_name_en")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(caadc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(caadc.db, c, "Invalid page_size parameter")
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(caadc.db, c, err.Error())
		return
	}

	queryParams := complianceDTOs.ComplianceAssessmentAcknowledgeListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	acknowledgements, err := caadc.service.GetAll(uint(assessmentID), currentUserID, &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(caadc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, acknowledgements)
}

func (caadc *ComplianceAssessmentAcknowledgeDepartmentController) CreateAcknowledgement(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(caadc.db, c, err.Error())
		return
	}

	var acknowledgeFormDTO complianceDTOs.ComplianceAssessmentAcknowledgeFormDTO
	if err := c.ShouldBindJSON(&acknowledgeFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	acknowledgeID, err := caadc.service.Create(uint(assessmentID), currentUserID, &acknowledgeFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(caadc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      acknowledgeID,
		"message": "Assessment acknowledged successfully",
	})
}

func (caadc *ComplianceAssessmentAcknowledgeDepartmentController) DeleteAcknowledgement(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(caadc.db, c, err.Error())
		return
	}

	userID, err := caadc.service.Delete(uint(assessmentID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(caadc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"message": "Assessment acknowledgement removed successfully",
	})
}

func (caadc *ComplianceAssessmentAcknowledgeDepartmentController) GetAcknowledgementSummary(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := caadc.service.GetAssessmentAcknowledgmentSummary(uint(assessmentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(caadc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, summary)
}

func (caadc *ComplianceAssessmentAcknowledgeDepartmentController) GetAcknowledgementSummaryByOrgUnit(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summaries, err := caadc.service.GetAssessmentAcknowledgmentSummaryByOrgUnit(uint(assessmentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(caadc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, summaries)
}
