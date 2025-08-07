package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectAssociatedItemController struct {
	service services.ProjectAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewProjectAssociatedItemController(service services.ProjectAssociatedItemServiceInterface) *ProjectAssociatedItemController {
	return &ProjectAssociatedItemController{
		service: service,
	}
}

func (paic *ProjectAssociatedItemController) GetAllProjectDocuments(c *gin.Context) {

	projectIDStr := c.Param("projectId")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	valueStr := c.Query("value")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(paic.db, c, "Invalid value parameter")
			return
		}
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Value:    value,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := paic.service.GetAllProjectDocuments(uint(projectID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}
func (paic *ProjectAssociatedItemController) CreateProjectDocument(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	err = paic.service.CreateProjectDocument(uint(projectID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (paic *ProjectAssociatedItemController) DeleteProjectDocument(c *gin.Context) {
	projectID, _ := strconv.ParseUint(c.Param("projectId"), 10, 64)
	documentID, _ := strconv.ParseUint(c.Param("documentId"), 10, 64)
	err := paic.service.DeleteProjectDocument(uint(projectID), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from project successfully"})
}

func (paic *ProjectAssociatedItemController) GetAllProjectControls(c *gin.Context) {

	projectIDStr := c.Param("projectId")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	valueStr := c.Query("value")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(paic.db, c, "Invalid value parameter")
			return
		}
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Value:    value,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := paic.service.GetAllProjectControls(uint(projectID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}

func (paic *ProjectAssociatedItemController) CreateProjectControl(c *gin.Context) {
	var input commonDTOs.AssociatedControlMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	err = paic.service.CreateProjectControl(uint(projectID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (paic *ProjectAssociatedItemController) DeleteProjectContol(c *gin.Context) {
	projectID, _ := strconv.ParseUint(c.Param("projectId"), 10, 64)
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)
	err := paic.service.DeleteProjectContol(uint(projectID), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from project successfully"})
}

func (paic *ProjectAssociatedItemController) GetAllProjectProcesses(c *gin.Context) {

	projectIDStr := c.Param("projectId")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	valueStr := c.Query("value")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(paic.db, c, "Invalid value parameter")
			return
		}
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Value:    value,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := paic.service.GetAllProjectProcesses(uint(projectID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}

func (paic *ProjectAssociatedItemController) CreateProjectProcess(c *gin.Context) {
	var input commonDTOs.AssociatedProcessMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	err = paic.service.CreateProjectProcess(uint(projectID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (paic *ProjectAssociatedItemController) DeleteProjectProcess(c *gin.Context) {
	projectID, _ := strconv.ParseUint(c.Param("projectId"), 10, 64)
	processID, _ := strconv.ParseUint(c.Param("processId"), 10, 64)
	err := paic.service.DeleteProjectProcess(uint(projectID), uint(processID))
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Process unmapped from project successfully"})
}
