package controllers

import (
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JobDescriptionCompetencyController struct {
	service services.JobDescriptionCompetencyServiceInterface
	db      *gorm.DB
}

func NewJobDescriptionCompetencyController(service services.JobDescriptionCompetencyServiceInterface) *JobDescriptionCompetencyController {
	return &JobDescriptionCompetencyController{
		service: service,
	}
}

func (jdcc *JobDescriptionCompetencyController) MapCompetencies(c *gin.Context) {
	var form dtos.JobDescriptionMapCompetenciesDTO

	jdID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(jdcc.db, c, err.Error())
		return
	}

	if err := jdcc.service.MapCompetencies(uint(jdID), currentUserID, &form); err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(jdcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

// PUT /job-description-competencies/:id
func (jdcc *JobDescriptionCompetencyController) UpdateCompetencyDetails(c *gin.Context) {
	var form dtos.UpdateJDCompetencyDTO
	_, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	cId, err := strconv.Atoi(c.Param("cId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(jdcc.db, c, err.Error())
		return
	}

	if err := jdcc.service.UpdateCompetencyDetails(uint(cId), currentUserID, &form); err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(jdcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

// DELETE /job-description-competencies/:id
func (jdcc *JobDescriptionCompetencyController) DeleteCompetency(c *gin.Context) {
	_, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	cId, err := strconv.Atoi(c.Param("cId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(jdcc.db, c, err.Error())
		return
	}

	if err := jdcc.service.DeleteCompetency(uint(cId), currentUserID); err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(jdcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
