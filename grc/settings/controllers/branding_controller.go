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

func (bc *BrandingController) GetAllBranding(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	allowedParams := []string{
		"created_by_ids",
		"updated_by_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := dtos.BrandingListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Filters: filters,
	}

	branding, err := bc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(bc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, branding)
}

func (bc *BrandingController) UpdateBranding(c *gin.Context) {
	var brandingFormDTO dtos.BrandingFormDTO

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

func (bc *BrandingController) GetBrandingPreviewImage(c *gin.Context) {
	brandingPreviewType := c.Query("type")
	if brandingPreviewType == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	brandingImagePath, err := bc.service.GetBrandingPreviewImage(brandingPreviewType)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, err)
			return
		default:
			utility.SendInternalErrResp(bc.db, c, err.Error())
		}
		return
	}

	c.File(brandingImagePath)
}
