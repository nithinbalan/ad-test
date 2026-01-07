package controllers

import (
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BrandingController struct {
	service services.BrandingServiceInterface
	db      *gorm.DB
}

func NewBrandingController(service services.BrandingServiceInterface) *BrandingController {
	return &BrandingController{
		service: service,
	}
}

func (bc *BrandingController) GetBranding(c *gin.Context) {
	branding, err := bc.service.GetFullDetails()
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, branding)
}

func (bc *BrandingController) UpdateBranding(c *gin.Context) {
	var brandingFormDTO dtos.BrandingUpdateFormDTO

	brandingIDStr := c.Param("id")
	brandingID, err := strconv.ParseUint(brandingIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&brandingFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bc.db, c, err.Error())
		return
	}

	err = bc.service.Update(uint(brandingID), currentUserID, &brandingFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(bc.db, c, err.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}
