package controllers

import (
	"jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectClassificationController struct {
	service services.ProjectClassificationServiceInterface
	db      *gorm.DB
}

func NewProjectClassificationController(service services.ProjectClassificationServiceInterface) *ProjectClassificationController {
	return &ProjectClassificationController{
		service: service,
	}
}

func (pcc *ProjectClassificationController) GetAllProjectClassifications(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	projectClassifications, err := pcc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(pcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, projectClassifications)
}

func (pcc *ProjectClassificationController) GetProjectClassificationDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	projectClassification, err := pcc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, projectClassification)
}

func (pcc *ProjectClassificationController) CreateProjectClassification(c *gin.Context) {
	var form dtos.ProjectClassificationFormDTO

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcc.db, c, err.Error())
		return
	}

	createdClassification, err := pcc.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdClassification)
}

func (pcc *ProjectClassificationController) UpdateProjectClassification(c *gin.Context) {
	var form dtos.ProjectClassificationFormDTO

	projectClassificationID, err := strconv.Atoi(c.Param("id"))
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
		utility.SendInternalErrResp(pcc.db, c, err.Error())
		return
	}

	updatedClassification, err := pcc.service.Update(uint(projectClassificationID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedClassification)
}

func (pcc *ProjectClassificationController) DeleteProjectClassifications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcc.db, c, err.Error())
		return
	}

	failedIDs, err := pcc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (pcc *ProjectClassificationController) ActivateProjectClassifications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcc.db, c, err.Error())
		return
	}

	err = pcc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (pcc *ProjectClassificationController) DeactivateProjectClassifications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcc.db, c, err.Error())
		return
	}

	err = pcc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}