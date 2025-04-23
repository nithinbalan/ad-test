package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/grc/organizational-context/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StakeholderNeedsAndExpectationEvaluationController struct {
	service services.StakeholderNeedsAndExpectationEvaluationServiceInterface
	db      *gorm.DB
}

func NewStakeholderNeedsAndExpectationEvaluationController(service services.StakeholderNeedsAndExpectationEvaluationServiceInterface) *StakeholderNeedsAndExpectationEvaluationController {
	return &StakeholderNeedsAndExpectationEvaluationController{
		service: service,
	}
}

func (sneec *StakeholderNeedsAndExpectationEvaluationController) GetAllStakeholderNeedsAndExpectationEvaluations(c *gin.Context) {
	stakeholderIDStr := c.Param("stakeholderId")

	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	needsAndExpectationIDStr := c.Param("needsAndExpectationId")

	needsAndExpectationID, err := strconv.ParseUint(needsAndExpectationIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	complianceStatuses, err := sneec.service.GetAll(uint(stakeholderID), uint(needsAndExpectationID), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sneec.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, complianceStatuses)
}

func (sneec *StakeholderNeedsAndExpectationEvaluationController) CreateStakeholderNeedsAndExpectationEvaluation(c *gin.Context) {
	var evaluationFormDTO organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationFormDTO

	stakeholderIDStr := c.Param("stakeholderId")

	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	needsAndExpectationIDStr := c.Param("needsAndExpectationId")

	needsAndExpectationID, err := strconv.ParseUint(needsAndExpectationIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&evaluationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sneec.db, c, err.Error())
		return
	}

	createdID, err := sneec.service.Create(uint(stakeholderID), uint(needsAndExpectationID), currentUserID, &evaluationFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sneec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (sneec *StakeholderNeedsAndExpectationEvaluationController) UpdateStakeholderNeedsAndExpectationEvaluation(c *gin.Context) {
	var evaluationFormDTO organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationFormDTO

	stakeholderIDStr := c.Param("stakeholderId")

	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	needsAndExpectationIDStr := c.Param("needsAndExpectationId")

	needsAndExpectationID, err := strconv.ParseUint(needsAndExpectationIDStr, 10, 32)
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

	if err := c.ShouldBindJSON(&evaluationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sneec.db, c, err.Error())
		return
	}

	err = sneec.service.Update(uint(stakeholderID), uint(needsAndExpectationID), uint(id), currentUserID, &evaluationFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sneec.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (sneec *StakeholderNeedsAndExpectationEvaluationController) GetLatestStakeholderNeedsAndExpectationEvaluation(c *gin.Context) {
	stakeholderIDStr := c.Param("stakeholderId")

	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	needsAndExpectationIDStr := c.Param("needsAndExpectationId")

	needsAndExpectationID, err := strconv.ParseUint(needsAndExpectationIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceActionPlanStatusUpdate, err := sneec.service.GetLatestStakeholderNeedsAndExpectationEvaluation(uint(stakeholderID), uint(needsAndExpectationID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sneec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, complianceActionPlanStatusUpdate)
}
