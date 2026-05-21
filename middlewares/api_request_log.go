package middlewares

import (
	"context"
	"errors"
	"fmt"
	"jethur/grc/models"
	"jethur/utility"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var skippedAPILogPrefixes = []string{
	"/logs/api-request-logs",
	"/health",
	"/healthz",
	"/metrics",
}

func shouldSkipAPILog(c *gin.Context) bool {
	if c.Request.Method == http.MethodOptions {
		return true
	}
	path := c.Request.URL.Path
	for _, p := range skippedAPILogPrefixes {
		if strings.Contains(path, p) {
			return true
		}
	}
	return false
}

// SweepStalePendingAPILogs marks rows left pending after the previous process
// exited as failed (completed_at set), matching chk_api_logs_status (completed/failed/aborted).
func SweepStalePendingAPILogs(db *gorm.DB) {
	if db == nil {
		return
	}
	res := db.Exec(`
		UPDATE api_request_logs
		SET log_status = ?,
		    completed_at = GETDATE(),
		    error_message = CASE
		        WHEN error_message IS NULL OR error_message = N'' THEN N'request interrupted (previous process ended before completing)'
		        ELSE error_message
		    END
		WHERE log_status = ?
	`, "failed", "pending")
	if res.Error != nil {
		log.Println("API log startup sweep error:", res.Error)
		return
	}
	if res.RowsAffected > 0 {
		log.Printf("API log startup sweep: marked %d stale pending row(s) as failed", res.RowsAffected)
	}
}

func APILoggingMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if shouldSkipAPILog(c) {
			c.Next()
			return
		}

		start := time.Now()
		rid := CurrentRequestID(c)

		entry := &models.APIRequestLog{
			Method:      c.Request.Method,
			URL:         c.Request.URL.RequestURI(),
			RequestedAt: start,
			LogStatus:   "pending",
			RequestID:   rid,
		}
		if err := db.Create(entry).Error; err != nil {
			log.Println("API log insert error:", err)
		}

		defer func() {
			if rec := recover(); rec != nil {
				stack := string(debug.Stack())
				msg := fmt.Sprintf("panic: %v", rec)
				updateAPILogFinal(db, entry.ID, http.StatusInternalServerError, start,
					"failed", &msg, &stack, currentUserID(c))
				panic(rec)
			}
		}()

		c.Next()

		statusCode := c.Writer.Status()
		logStatus := resolveLogStatus(c, statusCode)
		errMsg := collectErrorMessage(c, statusCode)
		updateAPILogFinal(db, entry.ID, statusCode, start, logStatus, errMsg, nil, currentUserID(c))
	}
}

func resolveLogStatus(c *gin.Context, statusCode int) string {
	if cerr := c.Request.Context().Err(); cerr != nil {
		if errors.Is(cerr, context.Canceled) || errors.Is(cerr, context.DeadlineExceeded) {
			return "aborted"
		}
	}
	switch {
	case statusCode >= 400:
		return "failed"
	default:
		return "completed"
	}
}

func collectErrorMessage(c *gin.Context, statusCode int) *string {
	if v, ok := c.Get(utility.ContextErrorKey); ok {
		if s, ok := v.(string); ok && s != "" {
			return &s
		}
	}
	if len(c.Errors) > 0 {
		msg := c.Errors.String()
		if msg != "" {
			return &msg
		}
	}
	if cerr := c.Request.Context().Err(); cerr != nil {
		msg := cerr.Error()
		return &msg
	}
	if statusCode >= 400 {
		msg := fmt.Sprintf("HTTP %d (no error message captured)", statusCode)
		return &msg
	}
	return nil
}

func currentUserID(c *gin.Context) *uint {
	if v, ok := c.Get("currentUserID"); ok {
		if id, ok := v.(uint); ok && id > 0 {
			return &id
		}
	}
	return nil
}

func updateAPILogFinal(
	db *gorm.DB, logID uint, statusCode int, start time.Time,
	status string, errMsg, stack *string, userID *uint,
) {
	if logID == 0 {
		return
	}

	completedAt := time.Now()
	durationMs := int(completedAt.Sub(start).Milliseconds())

	fields := map[string]any{
		"status_code":  statusCode,
		"completed_at": completedAt,
		"duration_ms":  durationMs,
		"log_status":   status,
	}
	if errMsg != nil {
		fields["error_message"] = *errMsg
	}
	if stack != nil {
		fields["stack_trace"] = *stack
	}
	if userID != nil {
		fields["user_id"] = *userID
	}

	if err := db.Model(&models.APIRequestLog{}).
		Where("id = ?", logID).
		Updates(fields).Error; err != nil {
		log.Println("API log update error:", err)
	}
}
