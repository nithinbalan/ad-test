package controllers

import (
	"jethur/grc/models"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JobDescriptionController struct {
	service services.JobDescriptionServiceInterface
	db      *gorm.DB
}

func NewJobDescriptionController(service services.JobDescriptionServiceInterface) *JobDescriptionController {
	return &JobDescriptionController{
		service: service,
	}
}

func (jdc *JobDescriptionController) GetAllJobDescriptions(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	jobDescriptions, err := jdc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(jdc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, jobDescriptions)
}

func (jdc *JobDescriptionController) GetJobDescriptionDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	jobDescription, err := jdc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, jobDescription)
}

func (jdc *JobDescriptionController) CreateJobDescription(c *gin.Context) {
	var jobDescription *models.JobDescription

	if err := c.ShouldBindJSON(&jobDescription); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(jdc.db, c, err.Error())
		return
	}

	createdJobDescription, err := jdc.service.Create(currentUserID, jobDescription)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(jdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdJobDescription)
}

func (jdc *JobDescriptionController) UpdateJobDescription(c *gin.Context) {
	jobDescriptionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedJobDescription *models.JobDescription
	if err := c.ShouldBindJSON(&updatedJobDescription); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(jdc.db, c, err.Error())
		return
	}

	updatedJobDescription, err = jdc.service.Update(uint(jobDescriptionID), currentUserID, updatedJobDescription)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(jdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedJobDescription)
}

func (jdc *JobDescriptionController) DeleteJobDescription(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(jdc.db, c, err.Error())
		return
	}

	failedIDs, err := jdc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(jdc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (jdc *JobDescriptionController) ActivateJobDescriptions(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(jdc.db, c, err.Error())
		return
	}

	err = jdc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(jdc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (jdc *JobDescriptionController) DeactivateJobDescriptions(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(jdc.db, c, err.Error())
		return
	}

	err = jdc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(jdc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
