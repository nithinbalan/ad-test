package controllers

import (
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BCMBIAOutageBucketController struct {
	service services.BCMBIAOutageBucketServiceInterface
	db      *gorm.DB
}

func NewBCMBIAOutageBucketController(service services.BCMBIAOutageBucketServiceInterface) *BCMBIAOutageBucketController {
	return &BCMBIAOutageBucketController{
		service: service,
	}
}

func (bbobc *BCMBIAOutageBucketController) GetAllBCMBIAOutageBuckets(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	bcmBiaOutageBuckets, err := bbobc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(bbobc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, bcmBiaOutageBuckets)
}

func (bbobc *BCMBIAOutageBucketController) GetBCMBIAOutageBucketDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	bcmBiaOutageBucket, err := bbobc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, bcmBiaOutageBucket)
}
