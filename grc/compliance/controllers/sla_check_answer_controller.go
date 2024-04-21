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

func (scqc *SLACheckAnswerController) CreateSLACheckAnswer(c *gin.Context) {

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
		utility.SendInternalErrResp(scqc.db, c, err.Error())
		return
	}

	createdID, err := scqc.service.CreateAnswer(uint(slaId), uint(questionId), currentUserID, &slaCheckAnswerFormDTO)

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
