package controllers

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BusinessImpactAssessmentController struct {
	service services.BusinessImpactAssessmentServiceInterface
	db      *gorm.DB
}

func NewBusinessImpactAssessmentController(service services.BusinessImpactAssessmentServiceInterface) *BusinessImpactAssessmentController {
	return &BusinessImpactAssessmentController{
		service: service,
	}
}

func (biac *BusinessImpactAssessmentController) GetAllBusinessImpactAssessments(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	keyword := c.Query("keyword")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"process_ids",
		"risk_level_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO{
		Sort:    sort,
		Keyword: keyword,
		Filters: filters,
	}
	businessImpactAssessments, err := biac.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, businessImpactAssessments)
}

func (biac *BusinessImpactAssessmentController) GetAllBusinessImpactAssessmentRecoveryObjectives(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	keyword := c.Query("keyword")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"process_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListQueryParamsDTO{
		Sort:    sort,
		Keyword: keyword,
		Filters: filters,
	}
	businessImpactAssessments, err := biac.service.GetAllBusinessImpactAssessmentRecoveryObjectives(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, businessImpactAssessments)
}

func (biac *BusinessImpactAssessmentController) CreateBusinessImpactAssessment(c *gin.Context) {
	var businessImpactAssessmentFormDTO businessContinuityDTOs.BusinessImpactFormDTO

	if err := c.ShouldBindJSON(&businessImpactAssessmentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	id, err := biac.service.Create(currentUserID, &businessImpactAssessmentFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(biac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (biac *BusinessImpactAssessmentController) UpdateImpactAssessment(c *gin.Context) {
	var BusinessImpactAssessmentFormDTO businessContinuityDTOs.BusinessImpactAssessmentFormDTO

	if err := c.ShouldBindJSON(&BusinessImpactAssessmentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	err = biac.service.UpdateImpactAssessment(BusinessImpactAssessmentFormDTO, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(biac.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (biac *BusinessImpactAssessmentController) UpdateRecoveryObjective(c *gin.Context) {
	var businessImpactAssessmentRecoveryObjectiveFormDTO businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveFormDTO

	if err := c.ShouldBindJSON(&businessImpactAssessmentRecoveryObjectiveFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	err = biac.service.UpdateRecoveryObjectives(businessImpactAssessmentRecoveryObjectiveFormDTO, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(biac.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}
