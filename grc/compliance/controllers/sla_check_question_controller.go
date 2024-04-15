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

// func (scc *SLACheckQuestionController) GetAllSLACheckQuestions(c *gin.Context) {

// 	slaIdStr := c.Param("slaId")

// 	slaId, err := strconv.ParseUint(slaIdStr, 10, 32)
// 	if err != nil {
// 		utility.SendNotFoundErrResp(c)
// 		return
// 	}

// 	checkIdStr := c.Param("checkId")

// 	checkId, err := strconv.ParseUint(checkIdStr, 10, 32)
// 	if err != nil {
// 		utility.SendNotFoundErrResp(c)
// 		return
// 	}

// 	sort := c.Query("sort")
// 	sortBy := c.Query("sort_by")
// 	keyword := c.Query("keyword")

// 	slaCheckQuestions, err := scc.service.GetAll(uint(slaId), uint(checkId), sort, sortBy, keyword)

// 	if err != nil {
// 		switch e := err.(type) {
// 		case *utility.ResourceNotFoundError:
// 			utility.SendNotFoundErrResp(c)
// 			return
// 		default:
// 			utility.SendInternalErrResp(scc.db, c, e.Error())
// 			return
// 		}
// 	}
// 	c.JSON(http.StatusOK, slaCheckQuestions)
// }

func (scc *SLACheckQuestionController) ChooseSLACheckQuestions(c *gin.Context) {

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
		utility.SendInternalErrResp(scc.db, c, err.Error())
		return
	}

	failedIDs, err := scc.service.ChooseQuestion(uint(slaId), uint(checkId), currentUserID, request.ChecklistQuestionIDs)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(scc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (scc *SLACheckQuestionController) CreateSLACheckQuestion(c *gin.Context) {

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
		utility.SendInternalErrResp(scc.db, c, err.Error())
		return
	}

	createdID, err := scc.service.Create(uint(slaId), uint(checkId), currentUserID, &slaCheckQuestionFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(scc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

// func (scc *SLACheckQuestionController) SLACheckQuestion(c *gin.Context) {
// 	var slaCheckQuestionFormDTO complianceDTOs.SLACheckQuestionFormDTO

// 	slaIDStr := c.Param("slaId")

// 	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
// 	if err != nil {
// 		utility.SendNotFoundErrResp(c)
// 		return
// 	}

// 	checkIdStr := c.Param("checkId")

// 	checkId, err := strconv.ParseUint(checkIdStr, 10, 32)
// 	if err != nil {
// 		utility.SendNotFoundErrResp(c)
// 		return
// 	}

// 	idStr := c.Param("id")

// 	id, err := strconv.ParseUint(idStr, 10, 64)
// 	if err != nil {
// 		utility.SendNotFoundErrResp(c)
// 		return
// 	}

// 	if err := c.ShouldBindJSON(&slaCheckQuestionFormDTO); err != nil {
// 		utility.SendInvalidErrResp(c, err)
// 		return
// 	}

// 	currentUserID, err := utility.GetCurrentUserID(c)
// 	if err != nil {
// 		utility.SendInternalErrResp(scc.db, c, err.Error())
// 		return
// 	}

// 	err = scc.service.Update(uint(slaID), uint(checkId), uint(id), currentUserID, &slaCheckQuestionFormDTO)
// 	if err != nil {
// 		switch e := err.(type) {
// 		case *utility.ResourceNotFoundError:
// 			utility.SendNotFoundErrResp(c)
// 			return
// 		case *utility.ValidationError:
// 			utility.SendInvalidErrResp(c, e)
// 			return
// 		default:
// 			utility.SendInternalErrResp(scc.db, c, e.Error())
// 			return
// 		}
// 	}

// 	c.Status(http.StatusOK)
// }
