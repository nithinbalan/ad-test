package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectSecurityGateChecklistTemplateController struct {
	service services.ProjectSecurityGateChecklistTemplateServiceInterface
	db      *gorm.DB
}

func NewProjectSecurityGateChecklistTemplateController(service services.ProjectSecurityGateChecklistTemplateServiceInterface) *ProjectSecurityGateChecklistTemplateController {
	return &ProjectSecurityGateChecklistTemplateController{
		service: service,
	}
}

func (psgctc *ProjectSecurityGateChecklistTemplateController) parseIDs(c *gin.Context) (templateID uint, versionID uint, ok bool) {
	tID, err := strconv.ParseUint(c.Param("templateId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return 0, 0, false
	}

	vID, err := strconv.ParseUint(c.Param("versionId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return 0, 0, false
	}

	return uint(tID), uint(vID), true
}

func (psgctc *ProjectSecurityGateChecklistTemplateController) GetAllProjectSecurityGateChecklistTemplates(c *gin.Context) {
	templates, err := psgctc.service.GetAll()
	if err != nil {
		utility.SendInternalErrResp(psgctc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, templates)
}

func (psgctc *ProjectSecurityGateChecklistTemplateController) GetProjectSecurityGateChecklistTemplateDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("templateId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	template, err := psgctc.service.GetDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(psgctc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, template)
}

func (psgctc *ProjectSecurityGateChecklistTemplateController) CreateProjectSecurityGateChecklistTemplateVersion(c *gin.Context) {
	templateID, err := strconv.ParseUint(c.Param("templateId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgctc.db, c, err.Error())
		return
	}

	id, err := psgctc.service.CreateVersion(uint(templateID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psgctc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (psgctc *ProjectSecurityGateChecklistTemplateController) PublishProjectSecurityGateChecklistTemplateVersion(c *gin.Context) {
	templateID, err := strconv.ParseUint(c.Param("templateId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	versionID, err := strconv.ParseUint(c.Param("versionId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgctc.db, c, err.Error())
		return
	}

	err = psgctc.service.PublishVersion(uint(templateID), uint(versionID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psgctc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (psgctc *ProjectSecurityGateChecklistTemplateController) CreateProjectSecurityGateChecklistTemplateRequirement(c *gin.Context) {
	templateID, versionID, ok := psgctc.parseIDs(c)
	if !ok {
		return
	}

	var form dtos.ProjectSecurityGateChecklistTemplateRequirementFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgctc.db, c, err.Error())
		return
	}

	id, err := psgctc.service.CreateRequirement(templateID, versionID, currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psgctc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (psgctc *ProjectSecurityGateChecklistTemplateController) UpdateProjectSecurityGateChecklistTemplateRequirement(c *gin.Context) {
	templateID, versionID, ok := psgctc.parseIDs(c)
	if !ok {
		return
	}

	requirementID, err := strconv.ParseUint(c.Param("requirementId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.ProjectSecurityGateChecklistTemplateRequirementFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgctc.db, c, err.Error())
		return
	}

	err = psgctc.service.UpdateRequirement(templateID, versionID, uint(requirementID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psgctc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (psgctc *ProjectSecurityGateChecklistTemplateController) DeleteProjectSecurityGateChecklistTemplateRequirement(c *gin.Context) {
	templateID, versionID, ok := psgctc.parseIDs(c)
	if !ok {
		return
	}

	requirementID, err := strconv.ParseUint(c.Param("requirementId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgctc.db, c, err.Error())
		return
	}

	err = psgctc.service.DeleteRequirement(templateID, versionID, uint(requirementID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psgctc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
