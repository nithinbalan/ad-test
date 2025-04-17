package controllers

import (
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KRILibraryController struct {
	service services.KRILibraryServiceInterface
	db      *gorm.DB
}

func NewKRILibraryController(service services.KRILibraryServiceInterface) *KRILibraryController {
	return &KRILibraryController{
		service: service,
	}
}

func (klc *KRILibraryController) GetAllKRILibrary(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(klc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(klc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := libraryDTOs.KRILibraryListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	kriLibrary, err := klc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(klc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, kriLibrary)
}

func (klc *KRILibraryController) GetKRILibraryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kriLibrary, err := klc.service.GetByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(klc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, kriLibrary)
}

func (klc *KRILibraryController) CreateKRILibrary(c *gin.Context) {
	var libraryFormDTOO libraryDTOs.KRILibraryFormDTO

	if err := c.ShouldBindJSON(&libraryFormDTOO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(klc.db, c, err.Error())
		return
	}

	id, err := klc.service.Create(currentUserID, &libraryFormDTOO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(klc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (klc *KRILibraryController) UpdateKRILibrary(c *gin.Context) {
	kriLibraryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var libraryFormDTOO libraryDTOs.KRILibraryFormDTO
	if err := c.ShouldBindJSON(&libraryFormDTOO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(klc.db, c, err.Error())
		return
	}

	err = klc.service.Update(uint(kriLibraryID), currentUserID, &libraryFormDTOO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(klc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (klc *KRILibraryController) DeleteKRILibraries(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(klc.db, c, err.Error())
		return
	}

	failedIDs, err := klc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(klc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (klc *KRILibraryController) ActivateKRILibraries(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(klc.db, c, err.Error())
		return
	}

	err = klc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(klc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (klc *KRILibraryController) DeactivateKRILibraries(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(klc.db, c, err.Error())
		return
	}

	err = klc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(klc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
