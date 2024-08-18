package controllers

import (
	"jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExpenseCategoryController struct {
	service services.ExpenseCategoryServiceInterface
	db      *gorm.DB
}

func NewExpenseCategoryController(service services.ExpenseCategoryServiceInterface) *ExpenseCategoryController {
	return &ExpenseCategoryController{
		service: service,
	}
}

func (ecc *ExpenseCategoryController) GetAllExpenseCategories(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	expenseCategories, err := ecc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(ecc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, expenseCategories)
}

func (ecc *ExpenseCategoryController) GetExpenseCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	expenseCategory, err := ecc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, expenseCategory)
}

func (ecc *ExpenseCategoryController) CreateExpenseCategory(c *gin.Context) {

	var expenseCategory *models.ExpenseCategory

	if err := c.ShouldBindJSON(&expenseCategory); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ecc.db, c, err.Error())
		return
	}

	createdExpenseCategory, err := ecc.service.Create(currentUserID, expenseCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ecc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdExpenseCategory)
}

func (ecc *ExpenseCategoryController) UpdateExpenseCategory(c *gin.Context) {

	expenseCategoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedExpenseCategory *models.ExpenseCategory
	if err := c.ShouldBindJSON(&updatedExpenseCategory); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ecc.db, c, err.Error())
		return
	}

	updatedExpenseCategory, err = ecc.service.Update(uint(expenseCategoryID), currentUserID, updatedExpenseCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ecc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedExpenseCategory)
}

func (ecc *ExpenseCategoryController) DeleteExpenseCategories(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ecc.db, c, err.Error())
		return
	}

	failedIDs, err := ecc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ecc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ecc *ExpenseCategoryController) ActivateExpenseCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ecc.db, c, err.Error())
		return
	}

	err = ecc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ecc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (ecc *ExpenseCategoryController) DeactivateExpenseCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ecc.db, c, err.Error())
		return
	}

	err = ecc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ecc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
