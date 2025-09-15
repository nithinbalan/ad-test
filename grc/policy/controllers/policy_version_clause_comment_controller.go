package controllers

import (
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/grc/policy/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PolicyVersionClauseCommentController struct {
	service services.PolicyVersionClauseCommentServiceInterface
	db      *gorm.DB
}

func NewPolicyVersionClauseCommentController(service services.PolicyVersionClauseCommentServiceInterface) *PolicyVersionClauseCommentController {
	return &PolicyVersionClauseCommentController{
		service: service,
	}
}

func (pvccc *PolicyVersionClauseCommentController) GetAllComments(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	clauseIDStr := c.Param("clauseId")
	clauseID, err := strconv.ParseUint(clauseIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comments, err := pvccc.service.GetAll(uint(documentID), uint(clauseID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(pvccc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (pvccc *PolicyVersionClauseCommentController) GetCommentDetails(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	clauseIDStr := c.Param("clauseId")
	clauseID, err := strconv.ParseUint(clauseIDStr, 10, 32)
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

	comment, err := pvccc.service.GetFullDetailsByID(uint(documentID), uint(clauseID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(pvccc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (pvccc *PolicyVersionClauseCommentController) CreateComment(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	clauseIDStr := c.Param("clauseId")
	clauseID, err := strconv.ParseUint(clauseIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO policyDTOs.PolicyVersionClauseCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pvccc.db, c, err.Error())
		return
	}

	id, err := pvccc.service.Create(uint(documentID), uint(clauseID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(pvccc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (pvccc *PolicyVersionClauseCommentController) UpdateComment(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	clauseIDStr := c.Param("clauseId")
	clauseID, err := strconv.ParseUint(clauseIDStr, 10, 32)
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

	var formDTO policyDTOs.PolicyVersionClauseCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pvccc.db, c, err.Error())
		return
	}

	err = pvccc.service.Update(uint(documentID), uint(clauseID), uint(commentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(pvccc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (pvccc *PolicyVersionClauseCommentController) DeleteComment(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	clauseIDStr := c.Param("clauseId")
	clauseID, err := strconv.ParseUint(clauseIDStr, 10, 32)
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
		utility.SendInternalErrResp(pvccc.db, c, err.Error())
		return
	}

	failedIDs, err := pvccc.service.Delete(uint(documentID), uint(clauseID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pvccc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
