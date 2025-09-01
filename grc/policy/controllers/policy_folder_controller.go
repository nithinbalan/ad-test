package controllers

import (
	"jethur/grc/models"
	"jethur/grc/policy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PolicyFolderController struct {
	service services.PolicyFolderServiceInterface
	db      *gorm.DB
}

func NewPolicyFolderController(service services.PolicyFolderServiceInterface) *PolicyFolderController {
	return &PolicyFolderController{
		service: service,
	}
}

func (cfc *PolicyFolderController) GetAllPolicyFolders(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	folders, err := cfc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(cfc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, folders)
}

func (cfc *PolicyFolderController) GetPolicyFolderDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	folder, err := cfc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, folder)
}

func (cfc *PolicyFolderController) CreatePolicyFolder(c *gin.Context) {
	var folder *models.PolicyFolder

	if err := c.ShouldBindJSON(&folder); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cfc.db, c, err.Error())
		return
	}

	createdFolder, err := cfc.service.Create(currentUserID, folder)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cfc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdFolder)
}

func (cfc *PolicyFolderController) UpdatePolicyFolder(c *gin.Context) {
	var updatedFolder *models.PolicyFolder

	folderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedFolder); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cfc.db, c, err.Error())
		return
	}

	updatedFolder, err = cfc.service.Update(uint(folderID), currentUserID, updatedFolder)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cfc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedFolder)
}

func (cfc *PolicyFolderController) DeletePolicyFolders(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cfc.db, c, err.Error())
		return
	}

	failedIDs, err := cfc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cfc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (cfc *PolicyFolderController) ActivatePolicyFolders(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cfc.db, c, err.Error())
		return
	}

	err = cfc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cfc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (cfc *PolicyFolderController) DeactivatePolicyFolders(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cfc.db, c, err.Error())
		return
	}

	err = cfc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cfc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
