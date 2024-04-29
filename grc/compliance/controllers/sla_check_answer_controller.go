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

type SLACheckAnswerController struct {
	service services.SLACheckAnswerServiceInterface
	db      *gorm.DB
}

func NewSLACheckAnswerController(service services.SLACheckAnswerServiceInterface) *SLACheckAnswerController {
	return &SLACheckAnswerController{
		service: service,
	}
}

func (scac *SLACheckAnswerController) CreateSLACheckAnswer(c *gin.Context) {
	var slaCheckAnswerFormDTO complianceDTOs.SLACheckAnswerFormDTO

	slaIdStr := c.Param("slaId")

	slaId, err := strconv.ParseUint(slaIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	questionIdStr := c.Param("questionId")

	questionId, err := strconv.ParseUint(questionIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&slaCheckAnswerFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(scac.db, c, err.Error())
		return
	}

	createdID, err := scac.service.CreateAnswer(uint(slaId), uint(questionId), currentUserID, &slaCheckAnswerFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(scac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (scac *SLACheckAnswerController) UpdateSLACheckAnswer(c *gin.Context) {
	var slaCheckAnswerFormDTO complianceDTOs.SLACheckAnswerFormDTO

	slaIdStr := c.Param("slaId")

	slaId, err := strconv.ParseUint(slaIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	questionIdStr := c.Param("questionId")

	questionId, err := strconv.ParseUint(questionIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	answerIdStr := c.Param("answerId")

	answerId, err := strconv.ParseUint(answerIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&slaCheckAnswerFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(scac.db, c, err.Error())
		return
	}

	err = scac.service.UpdateAnswer(uint(slaId), uint(questionId), uint(answerId), currentUserID, &slaCheckAnswerFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(scac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (scac *SLACheckAnswerController) DeleteSLACheckAnswer(c *gin.Context) {
	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	questionIdStr := c.Param("questionId")

	questionId, err := strconv.ParseUint(questionIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	answerIdStr := c.Param("answerId")

	answerId, err := strconv.ParseUint(answerIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(scac.db, c, err.Error())
		return
	}

	failedID, err := scac.service.DeleteAnswer(uint(slaID), uint(questionId), uint(answerId), currentUserID)
	if err != nil {
		utility.SendInternalErrResp(scac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
