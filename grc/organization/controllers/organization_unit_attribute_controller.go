package controllers

import (
	"errors"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrganizationUnitAttributeController struct {
	service services.OrganizationUnitAttributeServiceInterface
	db      *gorm.DB
}

func NewOrganizationUnitAttributeController(service services.OrganizationUnitAttributeServiceInterface) *OrganizationUnitAttributeController {
	return &OrganizationUnitAttributeController{
		service: service,
	}
}

func (ouac *OrganizationUnitAttributeController) GetAllAttributes(c *gin.Context) {
	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "oua.created_at")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(ouac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(ouac.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"created_at": "oua.created_at",
		"updated_at": "oua.updated_at",
		"title":      "oua.title",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		sortField = "oua.created_at"
	}

	queryParams := organizationDTOs.OrganizationUnitAttributeListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	attributes, err := ouac.service.GetAll(uint(departmentID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(ouac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, attributes)
}

func (ouac *OrganizationUnitAttributeController) GetAttributeDetails(c *gin.Context) {
	idStr := c.Param("attributeId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	attribute, err := ouac.service.GetFullDetailsByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utility.SendNotFoundErrResp(c)
			return
		}
		utility.SendInternalErrResp(ouac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, attribute)
}

func (ouac *OrganizationUnitAttributeController) CreateAttribute(c *gin.Context) {
	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var attributeFormDTO organizationDTOs.OrganizationUnitAttributeFormDTO
	if err := c.ShouldBindJSON(&attributeFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ouac.db, c, err.Error())
		return
	}

	id, err := ouac.service.Create(uint(departmentID), currentUserID, &attributeFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ouac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (ouac *OrganizationUnitAttributeController) UpdateAttribute(c *gin.Context) {
	idStr := c.Param("attributeId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var attributeFormDTO organizationDTOs.OrganizationUnitAttributeFormDTO
	if err := c.ShouldBindJSON(&attributeFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ouac.db, c, err.Error())
		return
	}

	err = ouac.service.Update(uint(id), currentUserID, &attributeFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ouac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attribute updated successfully"})
}

func (ouac *OrganizationUnitAttributeController) DeleteAttributes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ouac.db, c, err.Error())
		return
	}

	failedIDs, err := ouac.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ouac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
