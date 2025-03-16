package controllers

import (
	"jethur/grc/documents/dtos"
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentRequirementController struct {
	service services.DocumentRequirementServiceInterface
	db      *gorm.DB
}

func NewDocumentRequirementController(service services.DocumentRequirementServiceInterface) *DocumentRequirementController {
	return &DocumentRequirementController{
		service: service,
	}
}

func (drc *DocumentRequirementController) GetAllDocumentRequirements(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	documentIdStr := c.Param("id")
	documentId, err := strconv.ParseUint(documentIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	idStr := c.Param("requirementGroupId")
	requirementGroupId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	documentRequirements, err := drc.service.GetAll(uint(documentId), uint(requirementGroupId), sort, sortBy, keyword)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, documentRequirements)
}

func (drc *DocumentRequirementController) GetDocumentRequirementDetails(c *gin.Context) {
	documentIdStr := c.Param("id")
	documentId, err := strconv.ParseUint(documentIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	requirementGroupIdStr := c.Param("requirementGroupId")
	requirementGroupId, err := strconv.ParseUint(requirementGroupIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	idStr := c.Param("requirementId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	documentRequirement, err := drc.service.GetByID(uint(documentId), uint(requirementGroupId), uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, documentRequirement)
}
func (drc *DocumentRequirementController) CreateDocumentRequirement(c *gin.Context) {
	documentIdStr := c.Param("id")
	documentId, err := strconv.ParseUint(documentIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var dto dtos.DocumentRequirementFormDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get requirementGroup ID from URL param
	requirementGroupIdStr := c.Param("requirementGroupId")
	requirementGroupId, err := strconv.ParseUint(requirementGroupIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Get current user ID
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}

	// Call service method
	id, err := drc.service.Create(uint(documentId), uint(requirementGroupId), currentUserID, &dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(drc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, id)
}

func (drc *DocumentRequirementController) UpdateDocumentRequirement(c *gin.Context) {
	var dto dtos.DocumentRequirementFormDTO
	documentIdStr := c.Param("id")
	documentId, err := strconv.ParseUint(documentIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	// Get IDs from URL params
	requirementId, err := strconv.Atoi(c.Param("requirementId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	requirementGroupIdStr := c.Param("requirementGroupId")
	requirementGroupId, err := strconv.ParseUint(requirementGroupIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Bind request body to DTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get current user ID
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}

	// Call service method
	updatedRequirement, err := drc.service.Update(uint(documentId), uint(requirementGroupId), uint(requirementId), currentUserID, &dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(drc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedRequirement)
}

func (drc *DocumentRequirementController) DeleteDocumentRequirements(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	documentIdStr := c.Param("id")
	documentId, err := strconv.ParseUint(documentIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	requirementGroupIdStr := c.Param("requirementGroupId")
	requirementGroupId, err := strconv.ParseUint(requirementGroupIdStr, 10, 32)
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
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}

	failedIDs, err := drc.service.Delete(uint(documentId), uint(requirementGroupId), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
