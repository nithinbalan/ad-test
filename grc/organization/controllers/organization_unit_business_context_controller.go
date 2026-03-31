package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrganizationUnitBusinessContextController struct {
	service services.OrganizationUnitBusinessContextServiceInterface
	db      *gorm.DB
}

func NewOrganizationUnitBusinessContextController(service services.OrganizationUnitBusinessContextServiceInterface) *OrganizationUnitBusinessContextController {
	return &OrganizationUnitBusinessContextController{
		service: service,
	}
}

func (oubcc *OrganizationUnitBusinessContextController) GetAllBusinessContexts(c *gin.Context) {
	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "oubc.created_at")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"created_at": "oubc.created_at",
		"updated_at": "oubc.updated_at",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		sortField = "oubc.created_at"
	}

	queryParams := organizationDTOs.OrganizationUnitBusinessContextListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortField,
		Keyword: keyword,
	}

	contexts, err := oubcc.service.GetAll(uint(departmentID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(oubcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, contexts)
}

func (oubcc *OrganizationUnitBusinessContextController) GetBusinessContextDetails(c *gin.Context) {
	departmentID, err := strconv.Atoi(c.Param("departmentId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	context, err := oubcc.service.GetFullDetailsByID(uint(departmentID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		utility.SendInternalErrResp(oubcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, context)
}

func (oubcc *OrganizationUnitBusinessContextController) CreateBusinessContext(c *gin.Context) {
	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var contextFormDTO organizationDTOs.OrganizationUnitBusinessContextFormDTO
	if err := c.ShouldBindJSON(&contextFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(oubcc.db, c, err.Error())
		return
	}

	id, err := oubcc.service.Create(uint(departmentID), currentUserID, &contextFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(oubcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (oubcc *OrganizationUnitBusinessContextController) UpdateBusinessContext(c *gin.Context) {
	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var contextFormDTO organizationDTOs.OrganizationUnitBusinessContextFormDTO
	if err := c.ShouldBindJSON(&contextFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(oubcc.db, c, err.Error())
		return
	}

	err = oubcc.service.Update(uint(departmentID), currentUserID, &contextFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(oubcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (oubcc *OrganizationUnitBusinessContextController) DeleteBusinessContext(c *gin.Context) {
	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(oubcc.db, c, err.Error())
		return
	}

	err = oubcc.service.Delete(uint(departmentID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(oubcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
