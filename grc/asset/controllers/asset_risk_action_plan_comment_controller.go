package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"

	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetRiskActionPlanCommentController struct {
	service services.AssetRiskActionPlanCommentServiceInterface
	db      *gorm.DB
}

func NewAssetRiskActionPlanCommentController(service services.AssetRiskActionPlanCommentServiceInterface) *AssetRiskActionPlanCommentController {
	return &AssetRiskActionPlanCommentController{
		service: service,
	}
}

func (arapcc *AssetRiskActionPlanCommentController) GetAllAssetRiskActionPlanComments(c *gin.Context) {
	riskIDStr := c.Param("id")

	var riskID uint64
	var err error

	if riskIDStr != "" {
		riskID, err = strconv.ParseUint(riskIDStr, 10, 32)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		riskID = 0
	}
	actionPlanID, err := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comments, err := arapcc.service.GetAll(uint(riskID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(arapcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (arapcc *AssetRiskActionPlanCommentController) GetAssetRiskActionPlanCommentDetails(c *gin.Context) {
	riskIDStr := c.Param("id")

	var riskID uint64
	var err error

	if riskIDStr != "" {
		riskID, err = strconv.ParseUint(riskIDStr, 10, 32)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		riskID = 0
	}
	actionPlanID, err := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comment, err := arapcc.service.GetFullDetailsByID(uint(riskID), uint(actionPlanID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(arapcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (arapcc *AssetRiskActionPlanCommentController) CreateAssetRiskActionPlanComment(c *gin.Context) {
	riskIDStr := c.Param("id")

	var riskID uint64
	var err error

	if riskIDStr != "" {
		riskID, err = strconv.ParseUint(riskIDStr, 10, 32)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		riskID = 0
	}
	actionPlanID, err := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO assetDTOs.AssetRiskActionPlanCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arapcc.db, c, err.Error())
		return
	}

	id, err := arapcc.service.Create(uint(riskID), uint(actionPlanID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(arapcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (arapcc *AssetRiskActionPlanCommentController) UpdateAssetRiskActionPlanComment(c *gin.Context) {
	riskIDStr := c.Param("id")

	var riskID uint64
	var err error

	if riskIDStr != "" {
		riskID, err = strconv.ParseUint(riskIDStr, 10, 32)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		riskID = 0
	}
	actionPlanID, err := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO assetDTOs.AssetRiskActionPlanCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arapcc.db, c, err.Error())
		return
	}

	err = arapcc.service.Update(uint(riskID), uint(actionPlanID), uint(commentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(arapcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arapcc *AssetRiskActionPlanCommentController) DeleteAssetRiskActionPlanComments(c *gin.Context) {
	riskIDStr := c.Param("id")

	var riskID uint64
	var err error

	if riskIDStr != "" {
		riskID, err = strconv.ParseUint(riskIDStr, 10, 32)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		riskID = 0
	}
	actionPlanID, err := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
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
		utility.SendInternalErrResp(arapcc.db, c, err.Error())
		return
	}

	failedIDs, err := arapcc.service.Delete(uint(riskID), uint(actionPlanID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arapcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (arapcc *AssetRiskActionPlanCommentController) GetAssetRiskActionPlanCommentDocumentPreview(c *gin.Context) {
	riskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanID, err := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	id, err := strconv.ParseUint(c.Param("commentId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	docPath, err := arapcc.service.GetDocumentPath(uint(riskID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(arapcc.db, c, err.Error())
		}
		return
	}
	c.File(docPath)
}

func (arapcc *AssetRiskActionPlanCommentController) DownloadAssetRiskActionPlanCommentDocument(c *gin.Context) {
	riskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanID, err := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	id, err := strconv.ParseUint(c.Param("commentId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	docPath, err := arapcc.service.GetDocumentPath(uint(riskID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(arapcc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(docPath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(docPath)
}
