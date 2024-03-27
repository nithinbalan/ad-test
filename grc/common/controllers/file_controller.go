package controllers

import (
	"jethur/grc/common/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FileController struct {
	service services.FileServiceInterface
	db      *gorm.DB
}

func NewFileController(service services.FileServiceInterface, db *gorm.DB) *FileController {
	return &FileController{
		service: service,
		db:      db,
	}
}

func (fc *FileController) TempFileUpload(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	uploadedFile, err := fc.service.TempFileUpload(file)
	if err != nil {
		utility.SendInternalErrResp(fc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, uploadedFile)
}
