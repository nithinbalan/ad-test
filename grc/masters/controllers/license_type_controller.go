package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LicenseTypeController struct {
	service services.LicenseTypeServiceInterface
	db      *gorm.DB
}

func NewLicenseTypeController(service services.LicenseTypeServiceInterface) *LicenseTypeController {
	return &LicenseTypeController{
		service: service,
	}
}

func (ltc *LicenseTypeController) GetAllLicenseTypes(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	licenseTypes, err := ltc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(ltc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, licenseTypes)
}

func (ltc *LicenseTypeController) GetLicenseTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	licenseType, err := ltc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, licenseType)
}
