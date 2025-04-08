package controllers

import (
	"errors"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetOtherDocumentController struct {
	service services.AssetOtherDocumentServiceInterface
	db      *gorm.DB
}

func NewAssetOtherDocumentController(service services.AssetOtherDocumentServiceInterface) *AssetOtherDocumentController {
	return &AssetOtherDocumentController{
		service: service,
	}
}

func (aodc *AssetOtherDocumentController) GetAllDocuments(c *gin.Context) {
	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid asset id"))
		return
	}

	// Get keyword from query params
	keyword := c.Query("keyword")

	// Call service with keyword filter
	result, err := aodc.service.GetAll(uint(assetID), keyword)
	if err != nil {
		utility.SendInternalErrResp(aodc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (aodc *AssetOtherDocumentController) GetDocumentDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	documentIDStr := c.Param("documentId")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	result, err := aodc.service.GetByID(uint(id), uint(documentID))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (aodc *AssetOtherDocumentController) CreateDocument(c *gin.Context) {
	var form *dtos.AssetOtherDocumentFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	assetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid asset id"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aodc.db, c, err.Error())
		return
	}

	ids, err := aodc.service.Create(uint(assetID), form, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aodc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"ids": ids})
}

func (aodc *AssetOtherDocumentController) UpdateDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.AssetOtherDocumentFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	documentIDStr := c.Param("documentId")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aodc.db, c, err.Error())
		return
	}

	err = aodc.service.Update(uint(id), uint(documentID), &form, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aodc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (aodc *AssetOtherDocumentController) DeleteDocuments(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aodc.db, c, err.Error())
		return
	}

	failedIDs, err := aodc.service.Delete(uint(id), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aodc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
