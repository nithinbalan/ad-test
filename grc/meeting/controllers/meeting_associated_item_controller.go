package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/meeting/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MeetingAssociatedItemController struct {
	service services.MeetingAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewMeetingAssociatedItemController(service services.MeetingAssociatedItemServiceInterface) *MeetingAssociatedItemController {
	return &MeetingAssociatedItemController{
		service: service,
	}
}

func (maic *MeetingAssociatedItemController) GetAllMeetingControls(c *gin.Context) {

	meetingIDStr := c.Param("id")

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	meetings, err := maic.service.GetAllControls(uint(meetingID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, meetings)
}

func (maic *MeetingAssociatedItemController) GetAllMeetingComplianceRequirements(c *gin.Context) {

	meetingIDStr := c.Param("id")

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page_size parameter")
		return
	}
	queryParams := meetingDTOs.MeetingAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	meetings, err := maic.service.GetAllComplianceRequirements(uint(meetingID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, meetings)
}

func (maic *MeetingAssociatedItemController) GetAllMeetingRisks(c *gin.Context) {

	meetingIDStr := c.Param("id")

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	meetings, err := maic.service.GetAllRisks(uint(meetingID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, meetings)
}

func (maic *MeetingAssociatedItemController) GetAllMeetingProcesses(c *gin.Context) {

	meetingIDStr := c.Param("id")

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	meetings, err := maic.service.GetAllProcesses(uint(meetingID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, meetings)
}

func (maic *MeetingAssociatedItemController) GetAllMeetingIncidents(c *gin.Context) {

	meetingIDStr := c.Param("id")

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	meetings, err := maic.service.GetAllIncidents(uint(meetingID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, meetings)
}

func (maic *MeetingAssociatedItemController) GetAllMeetingDocuments(c *gin.Context) {

	mmetingIDStr := c.Param("id")

	meetingID, err := strconv.ParseUint(mmetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := maic.service.GetAllMeetingDocuments(uint(meetingID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}
func (maic *MeetingAssociatedItemController) CreateMeetingDocument(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	meetingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	err = maic.service.CreateMeetingDocument(uint(meetingID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (maic *MeetingAssociatedItemController) DeleteMeetingDocument(c *gin.Context) {
	meetingID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	documentID, _ := strconv.ParseUint(c.Param("documentId"), 10, 64)
	err := maic.service.DeleteMeetingDocument(uint(meetingID), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from meeting successfully"})
}

func (maic *MeetingAssociatedItemController) CreateMeetingControl(c *gin.Context) {
	var input commonDTOs.AssociatedControlMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	meetingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	err = maic.service.CreateMeetingControl(uint(meetingID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (maic *MeetingAssociatedItemController) DeleteMeetingControl(c *gin.Context) {
	meetingID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)
	err := maic.service.DeleteMeetingControl(uint(meetingID), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Control unmapped from meeting successfully"})
}

func (maic *MeetingAssociatedItemController) CreateMeetingProcess(c *gin.Context) {
	var input commonDTOs.AssociatedProcessMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	meetingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	err = maic.service.CreateMeetingProcess(uint(meetingID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (maic *MeetingAssociatedItemController) DeleteMeetingProcess(c *gin.Context) {
	meetingID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	processID, _ := strconv.ParseUint(c.Param("processId"), 10, 64)
	err := maic.service.DeleteMeetingProcess(uint(meetingID), uint(processID))
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Process unmapped from meeting successfully"})
}

func (maic *MeetingAssociatedItemController) GetAllMeetingStrategicObjectives(c *gin.Context) {

	mmetingIDStr := c.Param("id")

	meetingID, err := strconv.ParseUint(mmetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	objectives, err := maic.service.GetAllMeetingStrategicObjectives(uint(meetingID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, objectives)
}
func (maic *MeetingAssociatedItemController) CreateMeetingStrategicObjective(c *gin.Context) {
	var input commonDTOs.AssociatedStrategicObjectiveMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	meetingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	err = maic.service.CreateMeetingStrategicObjective(uint(meetingID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (maic *MeetingAssociatedItemController) DeleteMeetingStrategicObjective(c *gin.Context) {
	meetingIDStr := c.Param("id")
	strategicObjectiveIDStr := c.Param("strategicObjectiveId")

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategicObjectiveID, err := strconv.ParseUint(strategicObjectiveIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = maic.service.DeleteMeetingStrategicObjective(uint(meetingID), uint(strategicObjectiveID))
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Strategic Objective unmapped from meeting successfully"})
}

func (maic *MeetingAssociatedItemController) CreateMeetingRisk(c *gin.Context) {
	var input commonDTOs.AssociatedRiskMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	meetingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	err = maic.service.CreateMeetingRisk(uint(meetingID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (maic *MeetingAssociatedItemController) DeleteMeetingRisk(c *gin.Context) {
	meetingIDStr := c.Param("id")
	riskIDStr := c.Param("riskId")

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = maic.service.DeleteMeetingRisk(uint(meetingID), uint(riskID))
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Risk unmapped from meeting successfully"})
}

func (maic *MeetingAssociatedItemController) CreateMeetingIncident(c *gin.Context) {
	var input commonDTOs.AssociatedIncidentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	meetingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	err = maic.service.CreateMeetingIncident(uint(meetingID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (maic *MeetingAssociatedItemController) DeleteMeetingIncident(c *gin.Context) {
	meetingIDStr := c.Param("id")
	incidentIDStr := c.Param("incidentId")

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = maic.service.DeleteMeetingIncident(uint(meetingID), uint(incidentID))
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident unmapped from meeting successfully"})
}

func (maic *MeetingAssociatedItemController) GetAllMeetingAssets(c *gin.Context) {

	mmetingIDStr := c.Param("id")

	meetingID, err := strconv.ParseUint(mmetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	assets, err := maic.service.GetAllMeetingAssets(uint(meetingID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, assets)
}

func (maic *MeetingAssociatedItemController) CreateMeetingAsset(c *gin.Context) {
	var input commonDTOs.AssociatedAssetMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	meetingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	err = maic.service.CreateMeetingAsset(uint(meetingID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (maic *MeetingAssociatedItemController) DeleteMeetingAsset(c *gin.Context) {
	meetingIDStr := c.Param("id")
	assetIDStr := c.Param("assetId")

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = maic.service.DeleteMeetingAsset(uint(meetingID), uint(assetID))
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset unmapped from meeting successfully"})
}

func (maic *MeetingAssociatedItemController) GetAllMeetingKPI(c *gin.Context) {

	mmetingIDStr := c.Param("id")

	meetingID, err := strconv.ParseUint(mmetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	kpis, err := maic.service.GetAllMeetingKPI(uint(meetingID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, kpis)
}

func (maic *MeetingAssociatedItemController) CreateMeetingKPI(c *gin.Context) {
	var input commonDTOs.AssociatedKPIMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	meetingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	err = maic.service.CreateMeetingKPI(uint(meetingID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (maic *MeetingAssociatedItemController) DeleteMeetingKPI(c *gin.Context) {
	meetingIDStr := c.Param("id")
	kpiIDStr := c.Param("kpiRegisterId")

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kpiID, err := strconv.ParseUint(kpiIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = maic.service.DeleteMeetingKPI(uint(meetingID), uint(kpiID))
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KPI unmapped from meeting successfully"})
}

func (maic *MeetingAssociatedItemController) GetAllClauses(c *gin.Context) {
	mmetingIDStr := c.Param("id")

	meetingID, err := strconv.ParseUint(mmetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	clauses, err := maic.service.GetAllClauses(uint(meetingID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, clauses)
}
func (maic *MeetingAssociatedItemController) UpdateClauses(c *gin.Context) {
	var clauseFormDTO commonDTOs.AssociatedClauseFormDTO

	meetingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&clauseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}

	err = maic.service.UpdateClauses(uint(meetingID), currentUserID, &clauseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Clauses against the meeting updated successfully"})
}
