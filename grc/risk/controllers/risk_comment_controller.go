package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskCommentController struct {
	service services.RiskCommentServiceInterface
	db      *gorm.DB
}

func NewRiskCommentController(service services.RiskCommentServiceInterface) *RiskCommentController {
	return &RiskCommentController{service: service}
}

func (rcc *RiskCommentController) GetAllComments(c *gin.Context) {
	riskIDStr := c.Param("riskId")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comments, err := rcc.service.GetAll(uint(riskID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (rcc *RiskCommentController) GetCommentDetails(c *gin.Context) {
	riskIDStr := c.Param("riskId")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
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

	comment, err := rcc.service.GetFullDetailsByID(uint(riskID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (rcc *RiskCommentController) CreateComment(c *gin.Context) {
	riskIDStr := c.Param("riskId")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.RiskCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	id, err := rcc.service.Create(uint(riskID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (rcc *RiskCommentController) UpdateComment(c *gin.Context) {
	riskIDStr := c.Param("riskId")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
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

	var formDTO dtos.RiskCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	err = rcc.service.Update(uint(riskID), uint(commentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rcc *RiskCommentController) DeleteComment(c *gin.Context) {
	riskIDStr := c.Param("riskId")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
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
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	failedIDs, err := rcc.service.Delete(uint(riskID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (rcc *RiskCommentController) SubmitRisk(c *gin.Context) {

	riskID, err := strconv.Atoi(c.Param("riskId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	err = rcc.service.SubmitRisk(uint(riskID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rcc *RiskCommentController) RevertRisk(c *gin.Context) {
	var commentForm dtos.RiskCommentFormDTO

	riskID, err := strconv.Atoi(c.Param("riskId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&commentForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	err = rcc.service.RevertRisk(uint(riskID), currentUserID, &commentForm)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rcc *RiskCommentController) RejectRisk(c *gin.Context) {
	var commentForm dtos.RiskCommentFormDTO

	riskID, err := strconv.Atoi(c.Param("riskId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&commentForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	err = rcc.service.RejectRisk(uint(riskID), currentUserID, &commentForm)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rcc *RiskCommentController) ApproveRisk(c *gin.Context) {
	var commentForm dtos.RiskCommentFormDTO

	riskID, err := strconv.Atoi(c.Param("riskId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&commentForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	err = rcc.service.ApproveRisk(uint(riskID), currentUserID, &commentForm)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rcc *RiskCommentController) ResubmitRisk(c *gin.Context) {
	var commentForm dtos.RiskCommentFormDTO

	riskID, err := strconv.Atoi(c.Param("riskId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&commentForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	err = rcc.service.ResubmitRisk(uint(riskID), currentUserID, &commentForm)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
