package controllers

import (
	"fmt"
	"jethur/grc/jethur-ai/dtos"
	"jethur/grc/jethur-ai/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentAIController struct {
	service services.DocumentAIServiceInterface
	db      *gorm.DB
}

func NewDocumentAIController(service services.DocumentAIServiceInterface) *DocumentAIController {
	return &DocumentAIController{
		service: service,
	}
}

func (ac *DocumentAIController) GetDocumentSummary(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	summary, err := ac.service.GetAISummary(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, summary)
}

func (ac *DocumentAIController) GetDocumentRecommendations(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ac.db, c, err.Error())
		return
	}

	summary, err := ac.service.GetDocumentRecommendations(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, summary)
}

func (dcc *DocumentAIController) ChatWithDocument(c *gin.Context) {
	var chatRequest dtos.DocumentChatRequest
	if err := c.ShouldBindJSON(&chatRequest); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		http.Error(c.Writer, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	streamChan, err := dcc.service.ProcessChat(c.Request.Context(), chatRequest.Content, chatRequest.Role, chatRequest.DocumentID, chatRequest.ContextID)
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}

	// Stream each full JSON line
	for msg := range streamChan {
		fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
		flusher.Flush()
	}
}

// func (dcc *DocumentAIController) ChatWithDocument(c *gin.Context) {
// 	var chatRequest dtos.DocumentChatRequest
// 	if err := c.ShouldBindJSON(&chatRequest); err != nil {
// 		utility.SendInvalidErrResp(c, err)
// 		return
// 	}

// 	// SSE headers
// 	c.Writer.Header().Set("Content-Type", "text/event-stream")
// 	c.Writer.Header().Set("Cache-Control", "no-cache")
// 	c.Writer.Header().Set("Connection", "keep-alive")
// 	c.Writer.Header().Set("Transfer-Encoding", "chunked")

// 	flusher, ok := c.Writer.(http.Flusher)
// 	if !ok {
// 		http.Error(c.Writer, "Streaming unsupported!", http.StatusInternalServerError)
// 		return
// 	}

// 	responseChan, err := dcc.service.ProcessChat(c.Request.Context(), chatRequest.Content, chatRequest.Role, chatRequest.DocumentID, chatRequest.ContextID)
// 	if err != nil {
// 		switch e := err.(type) {
// 		case *utility.ResourceNotFoundError:
// 			utility.SendNotFoundErrResp(c)
// 		default:
// 			utility.SendInternalErrResp(dcc.db, c, e.Error())
// 		}
// 		return
// 	}

// 	for msg := range responseChan {
// 		fmt.Fprintf(c.Writer, "data: %s\n\n", strings.TrimSpace(msg))
// 		flusher.Flush()
// 	}
// }

// func (dcc *DocumentAIController) ChatWithDocument(c *gin.Context) {
// 	var chatRequest dtos.DocumentChatRequest
// 	if err := c.ShouldBindJSON(&chatRequest); err != nil {
// 		utility.SendInvalidErrResp(c, err)
// 		return
// 	}

// 	// Set headers for SSE
// 	c.Header("Content-Type", "text/event-stream")
// 	c.Header("Cache-Control", "no-cache")
// 	c.Header("Connection", "keep-alive")
// 	c.Header("Transfer-Encoding", "chunked")

// 	// Get response channel from service
// 	responseChan, err := dcc.service.ProcessChat(c.Request.Context(), chatRequest.Content, chatRequest.Role, chatRequest.DocumentID, chatRequest.ContextID)
// 	if err != nil {
// 		switch e := err.(type) {
// 		case *utility.ResourceNotFoundError:
// 			utility.SendNotFoundErrResp(c)
// 		default:
// 			utility.SendInternalErrResp(dcc.db, c, e.Error())
// 		}
// 		return
// 	}

// 	// Stream responses to client
// 	c.Stream(func(w io.Writer) bool {
// 		if msg, ok := <-responseChan; ok {
// 			c.SSEvent("message", msg)
// 			return true
// 		}
// 		return false
// 	})
// }
