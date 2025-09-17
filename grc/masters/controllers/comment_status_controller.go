package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentStatusController struct {
	service services.CommentStatusServiceInterface
	db      *gorm.DB
}

func NewCommentStatusController(service services.CommentStatusServiceInterface) *CommentStatusController {
	return &CommentStatusController{
		service: service,
	}
}

func (apsc *CommentStatusController) GetAllCommentStatuses(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	commentStatuses, err := apsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(apsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, commentStatuses)
}

func (apsc *CommentStatusController) GetCommentStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendInternalErrResp(apsc.db, c, "Invalid id")
		return
	}

	commentStatus, err := apsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, commentStatus)
}
