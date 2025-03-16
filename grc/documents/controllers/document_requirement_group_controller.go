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

type DocumentRequirementGroupController struct {
	service services.DocumentRequirementGroupServiceInterface
	db      *gorm.DB
}

func NewDocumentRequirementGroupController(service services.DocumentRequirementGroupServiceInterface) *DocumentRequirementGroupController {
	return &DocumentRequirementGroupController{
		service: service,
	}
}

func (drgc *DocumentRequirementGroupController) GetAllDocumentRequirementGroups(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	idStr := c.Param("id")
	documentId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	documentRequirementGroups, err := drgc.service.GetAll(uint(documentId), sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(drgc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, documentRequirementGroups)
}

func (drgc *DocumentRequirementGroupController) GetDocumentRequirementGroupDetails(c *gin.Context) {
	requirementGroupIdStr := c.Param("requirementGroupId")
	id, err := strconv.ParseUint(requirementGroupIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	idStr := c.Param("id")
	documentId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	documentRequirementGroup, err := drgc.service.GetByID(uint(documentId), uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, documentRequirementGroup)
}
func (drgc *DocumentRequirementGroupController) CreateDocumentRequirementGroup(c *gin.Context) {
	var dto dtos.DocumentRequirementGroupFormDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get document ID from URL param
	idStr := c.Param("id")
	documentId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Get current user ID
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(drgc.db, c, err.Error())
		return
	}

	// Call service method
	id, err := drgc.service.Create(uint(documentId), currentUserID, &dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(drgc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, id)
}

func (drgc *DocumentRequirementGroupController) UpdateDocumentRequirementGroup(c *gin.Context) {
	var dto dtos.DocumentRequirementGroupFormDTO

	// Get IDs from URL params
	requirementGroupId, err := strconv.Atoi(c.Param("requirementGroupId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	idStr := c.Param("id")
	documentId, err := strconv.ParseUint(idStr, 10, 32)
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
		utility.SendInternalErrResp(drgc.db, c, err.Error())
		return
	}

	// Call service method
	updatedRequirementGroup, err := drgc.service.Update(uint(documentId), uint(requirementGroupId), currentUserID, &dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(drgc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedRequirementGroup)
}

func (drgc *DocumentRequirementGroupController) DeleteDocumentRequirementGroups(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	idStr := c.Param("id")
	documentId, err := strconv.ParseUint(idStr, 10, 32)
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
		utility.SendInternalErrResp(drgc.db, c, err.Error())
		return
	}

	failedIDs, err := drgc.service.Delete(uint(documentId), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(drgc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (drgc *DocumentRequirementGroupController) ActivateDocumentRequirementGroups(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	idStr := c.Param("id")
	documentId, err := strconv.ParseUint(idStr, 10, 32)
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
		utility.SendInternalErrResp(drgc.db, c, err.Error())
		return
	}

	err = drgc.service.Activate(uint(documentId), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(drgc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (drgc *DocumentRequirementGroupController) DeactivateDocumentRequirementGroups(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	idStr := c.Param("id")
	documentId, err := strconv.ParseUint(idStr, 10, 32)
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
		utility.SendInternalErrResp(drgc.db, c, err.Error())
		return
	}

	err = drgc.service.Deactivate(uint(documentId), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(drgc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
