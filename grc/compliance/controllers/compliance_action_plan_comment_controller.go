package controllers

import (
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceActionPlanCommentController struct {
	service services.ComplianceActionPlanCommentServiceInterface
	db      *gorm.DB
}

func NewComplianceActionPlanCommentController(service services.ComplianceActionPlanCommentServiceInterface) *ComplianceActionPlanCommentController {
	return &ComplianceActionPlanCommentController{
		service: service,
	}
}

func (capc *ComplianceActionPlanCommentController) GetAllComments(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")
	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	actionPlanIDStr := c.Param("actionPlanId")
	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comments, err := capc.service.GetAll(uint(complianceID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(capc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (capc *ComplianceActionPlanCommentController) GetCommentDetails(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")
	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	actionPlanIDStr := c.Param("actionPlanId")
	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
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

	comment, err := capc.service.GetByID(uint(complianceID), uint(actionPlanID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(capc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (capc *ComplianceActionPlanCommentController) CreateComment(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")
	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	actionPlanIDStr := c.Param("actionPlanId")
	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO complianceDTOs.ComplianceActionPlanCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(capc.db, c, err.Error())
		return
	}

	id, err := capc.service.Create(uint(complianceID), uint(actionPlanID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(capc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (capc *ComplianceActionPlanCommentController) UpdateComment(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")
	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	actionPlanIDStr := c.Param("actionPlanId")
	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
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

	var formDTO complianceDTOs.ComplianceActionPlanCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(capc.db, c, err.Error())
		return
	}

	err = capc.service.Update(uint(complianceID), uint(actionPlanID), uint(commentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(capc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (capc *ComplianceActionPlanCommentController) DeleteComment(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")
	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	actionPlanIDStr := c.Param("actionPlanId")
	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
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
		utility.SendInternalErrResp(capc.db, c, err.Error())
		return
	}

	failedIDs, err := capc.service.Delete(uint(complianceID), uint(actionPlanID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(capc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
