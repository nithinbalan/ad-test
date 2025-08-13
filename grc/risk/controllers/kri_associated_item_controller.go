package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KRIAssociatedItemController struct {
	service services.KRIAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewKRIAssociatedItemController(service services.KRIAssociatedItemServiceInterface) *KRIAssociatedItemController {
	return &KRIAssociatedItemController{
		service: service,
	}
}

func (paic *KRIAssociatedItemController) GetAllKRIDocuments(c *gin.Context) {

	kriIDStr := c.Param("id")

	kriID, err := strconv.ParseUint(kriIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := paic.service.GetAllKRIDocuments(uint(kriID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}
func (paic *KRIAssociatedItemController) CreateKRIDocument(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kriID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	err = paic.service.CreateKRIDocument(uint(kriID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (paic *KRIAssociatedItemController) DeleteKRIDocument(c *gin.Context) {
	kriID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	documentID, _ := strconv.ParseUint(c.Param("documentId"), 10, 64)
	err := paic.service.DeleteKRIDocument(uint(kriID), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from KRI successfully"})
}

func (paic *KRIAssociatedItemController) GetAllKRIControls(c *gin.Context) {

	kriIDStr := c.Param("id")

	kriID, err := strconv.ParseUint(kriIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := paic.service.GetAllKRIControls(uint(kriID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}
func (paic *KRIAssociatedItemController) CreateKRIControl(c *gin.Context) {
	var input commonDTOs.AssociatedControlMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	kriID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	err = paic.service.CreateKRIControl(uint(kriID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (paic *KRIAssociatedItemController) DeleteKRIControl(c *gin.Context) {
	kriID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)
	err := paic.service.DeleteKRIControl(uint(kriID), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Control unmapped from KRI successfully"})
}
