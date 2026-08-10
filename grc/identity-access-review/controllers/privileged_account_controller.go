package controllers

import (
	"net/http"
	"strconv"

	"jethur/grc/identity-access-review/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PrivilegedAccountController struct {
	service services.PrivilegedAccountServiceInterface
	db      *gorm.DB
}

func NewPrivilegedAccountController(service services.PrivilegedAccountServiceInterface) *PrivilegedAccountController {
	return &PrivilegedAccountController{service: service}
}

func (pac *PrivilegedAccountController) GetAllPrivilegedAccounts(c *gin.Context) {
	privilegedAccounts, err := pac.service.GetAll(
		c.Query("sort_order"),
		c.Query("sort_by"),
		c.Query("keyword"),
		c.Query("status"),
	)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, privilegedAccounts)
}

func (pac *PrivilegedAccountController) GetPrivilegedAccountDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	privilegedAccount, err := pac.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	c.JSON(http.StatusOK, privilegedAccount)
}

func (pac *PrivilegedAccountController) CreatePrivilegedAccount(c *gin.Context) {
	var privilegedAccount *models.PrivilegedAccount
	if err := c.ShouldBindJSON(&privilegedAccount); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	createdAccount, err := pac.service.Create(currentUserID, privilegedAccount)
	if err != nil {
		if validationErr, ok := err.(*utility.ValidationError); ok {
			utility.SendInvalidErrResp(c, validationErr)
			return
		}
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}
	c.JSON(http.StatusCreated, createdAccount)
}

func (pac *PrivilegedAccountController) UpdatePrivilegedAccount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedAccount *models.PrivilegedAccount
	if err := c.ShouldBindJSON(&updatedAccount); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	updatedAccount, err = pac.service.Update(uint(id), currentUserID, updatedAccount)
	if err != nil {
		switch typedErr := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, typedErr)
		default:
			utility.SendInternalErrResp(pac.db, c, typedErr.Error())
		}
		return
	}
	c.JSON(http.StatusOK, updatedAccount)
}

func (pac *PrivilegedAccountController) DeletePrivilegedAccounts(c *gin.Context) {
	ids, currentUserID, ok := pac.bindBulkRequest(c)
	if !ok {
		return
	}

	failedIDs, err := pac.service.Delete(ids, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (pac *PrivilegedAccountController) ActivatePrivilegedAccounts(c *gin.Context) {
	ids, currentUserID, ok := pac.bindBulkRequest(c)
	if !ok {
		return
	}
	if err := pac.service.Activate(ids, currentUserID); err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (pac *PrivilegedAccountController) DeactivatePrivilegedAccounts(c *gin.Context) {
	ids, currentUserID, ok := pac.bindBulkRequest(c)
	if !ok {
		return
	}
	if err := pac.service.Deactivate(ids, currentUserID); err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (pac *PrivilegedAccountController) bindBulkRequest(c *gin.Context) ([]uint, uint, bool) {
	var request struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return nil, 0, false
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return nil, 0, false
	}
	return request.IDs, currentUserID, true
}
