package controllers

import (
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetPolicyCommentController struct {
	service services.AssetPolicyCommentServiceInterface
	db      *gorm.DB
}

func NewAssetPolicyCommentController(service services.AssetPolicyCommentServiceInterface) *AssetPolicyCommentController {
	return &AssetPolicyCommentController{service: service}
}

func (apcc *AssetPolicyCommentController) GetAllComments(c *gin.Context) {
	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assetPolicyIdStr := c.Param("policyId")
	assetPolicyId, err := strconv.ParseUint(assetPolicyIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comments, err := apcc.service.GetAll(uint(assetID), uint(assetPolicyId))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (apcc *AssetPolicyCommentController) GetCommentDetails(c *gin.Context) {
	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assetPolicyIdStr := c.Param("policyId")
	assetPolicyId, err := strconv.ParseUint(assetPolicyIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comment, err := apcc.service.GetFullDetailsByID(uint(assetID), uint(assetPolicyId), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (apcc *AssetPolicyCommentController) CreateComment(c *gin.Context) {
	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assetPolicyIdStr := c.Param("policyId")
	assetPolicyId, err := strconv.ParseUint(assetPolicyIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.AssetPolicyCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apcc.db, c, err.Error())
		return
	}

	id, err := apcc.service.Create(uint(assetID), uint(assetPolicyId), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (apcc *AssetPolicyCommentController) UpdateComment(c *gin.Context) {
	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assetPolicyIdStr := c.Param("policyId")
	assetPolicyId, err := strconv.ParseUint(assetPolicyIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.AssetPolicyCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apcc.db, c, err.Error())
		return
	}

	err = apcc.service.Update(uint(assetID), uint(assetPolicyId), uint(commentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (apcc *AssetPolicyCommentController) DeleteComment(c *gin.Context) {
	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assetPolicyIdStr := c.Param("policyId")
	assetPolicyId, err := strconv.ParseUint(assetPolicyIdStr, 10, 32)
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
		utility.SendInternalErrResp(apcc.db, c, err.Error())
		return
	}

	failedIDs, err := apcc.service.Delete(uint(assetID), uint(assetPolicyId), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
