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

type SLACheckQuestionController struct {
	service services.SLACheckQuestionServiceInterface
	db      *gorm.DB
}

func NewSLACheckQuestionController(service services.SLACheckQuestionServiceInterface) *SLACheckQuestionController {
	return &SLACheckQuestionController{
		service: service,
	}
}

func (scqc *SLACheckQuestionController) GetAllSLACheckQuestionGroups(c *gin.Context) {
	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	checkIdStr := c.Param("checkId")

	checkId, err := strconv.ParseUint(checkIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	slaCheckQuestionGroups, err := scqc.service.GetAllQuestionGroups(uint(slaID), uint(checkId))

	if err != nil {
		utility.SendInternalErrResp(scqc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, slaCheckQuestionGroups)
}

func (scqc *SLACheckQuestionController) CreateSLACheckQuestionGroup(c *gin.Context) {
	var slaCheckQuestionGroupFormDTO complianceDTOs.SLACheckQuestionGroupFormDTO

	slaIdStr := c.Param("slaId")

	slaId, err := strconv.ParseUint(slaIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	checkIdStr := c.Param("checkId")

	checkId, err := strconv.ParseUint(checkIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&slaCheckQuestionGroupFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(scqc.db, c, err.Error())
		return
	}

	createdID, err := scqc.service.CreateQuestionGroup(uint(slaId), uint(checkId), currentUserID, &slaCheckQuestionGroupFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(scqc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (scqc *SLACheckQuestionController) UpdateSLACheckQuestionGroup(c *gin.Context) {
	var slaCheckUpdateQuestionGroupDTO complianceDTOs.SLACheckUpdateQuestionGroupDTO

	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	checkIdStr := c.Param("checkId")

	checkId, err := strconv.ParseUint(checkIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	questionGroupIdStr := c.Param("id")

	questionGroupId, err := strconv.ParseUint(questionGroupIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&slaCheckUpdateQuestionGroupDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(scqc.db, c, err.Error())
		return
	}

	err = scqc.service.UpdateQuestionGroup(uint(slaID), uint(checkId), uint(questionGroupId), currentUserID, &slaCheckUpdateQuestionGroupDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(scqc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (scqc *SLACheckQuestionController) DeleteSLACheckQuestionGroup(c *gin.Context) {
	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	checkIdStr := c.Param("checkId")

	checkId, err := strconv.ParseUint(checkIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	questionGroupIdStr := c.Param("id")

	questionGroupId, err := strconv.ParseUint(questionGroupIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(scqc.db, c, err.Error())
		return
	}

	failedID, err := scqc.service.DeleteQuestionGroup(uint(slaID), uint(checkId), uint(questionGroupId), currentUserID)
	if err != nil {
		utility.SendInternalErrResp(scqc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}

func (scqc *SLACheckQuestionController) GetSLACheckQuestionsByGroup(c *gin.Context) {
	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	checkIdStr := c.Param("checkId")

	checkId, err := strconv.ParseUint(checkIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	groupIdStr := c.Param("id")

	groupId, err := strconv.ParseUint(groupIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	slaCheckQuestions, err := scqc.service.GetQuestionsByGroupID(uint(slaID), uint(checkId), uint(groupId))

	if err != nil {
		utility.SendInternalErrResp(scqc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, slaCheckQuestions)
}

func (scqc *SLACheckQuestionController) ChooseSLACheckQuestions(c *gin.Context) {
	var request struct {
		ChecklistQuestionIDs []uint `json:"checklist_question_ids"`
	}

	slaIdStr := c.Param("slaId")

	slaId, err := strconv.ParseUint(slaIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	checkIdStr := c.Param("checkId")

	checkId, err := strconv.ParseUint(checkIdStr, 10, 32)
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
		utility.SendInternalErrResp(scqc.db, c, err.Error())
		return
	}

	failedIDs, err := scqc.service.ChooseQuestions(uint(slaId), uint(checkId), currentUserID, request.ChecklistQuestionIDs)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(scqc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (scqc *SLACheckQuestionController) CreateSLACheckQuestion(c *gin.Context) {
	var slaCheckQuestionFormDTO complianceDTOs.SLACheckQuestionFormDTO

	slaIdStr := c.Param("slaId")

	slaId, err := strconv.ParseUint(slaIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	checkIdStr := c.Param("checkId")

	checkId, err := strconv.ParseUint(checkIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&slaCheckQuestionFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(scqc.db, c, err.Error())
		return
	}

	createdID, err := scqc.service.CreateQuestion(uint(slaId), uint(checkId), currentUserID, &slaCheckQuestionFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(scqc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (scqc *SLACheckQuestionController) UpdateSLACheckQuestion(c *gin.Context) {
	var slaCheckUpdateQuestionDTO complianceDTOs.SLACheckUpdateQuestionDTO

	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	checkIdStr := c.Param("checkId")

	checkId, err := strconv.ParseUint(checkIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	questionIdStr := c.Param("id")

	questionId, err := strconv.ParseUint(questionIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&slaCheckUpdateQuestionDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(scqc.db, c, err.Error())
		return
	}

	err = scqc.service.UpdateQuestion(uint(slaID), uint(checkId), uint(questionId), currentUserID, &slaCheckUpdateQuestionDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(scqc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (scqc *SLACheckQuestionController) DeleteSLACheckQuestion(c *gin.Context) {
	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	checkIdStr := c.Param("checkId")

	checkId, err := strconv.ParseUint(checkIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	questionIdStr := c.Param("id")

	questionId, err := strconv.ParseUint(questionIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(scqc.db, c, err.Error())
		return
	}

	failedID, err := scqc.service.DeleteQuestion(uint(slaID), uint(checkId), uint(questionId), currentUserID)
	if err != nil {
		utility.SendInternalErrResp(scqc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
