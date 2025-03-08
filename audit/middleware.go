// audit/middleware.go
package audit

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ContextKey for storing audit information in request context
type ContextKey string

const (
	// AuditUserKey is the context key for the username
	AuditUserKey ContextKey = "audit_username"
	
	// AuditOrgIDKey is the context key for the organization ID
	AuditOrgIDKey ContextKey = "audit_org_id"
)

// WithAuditContext adds audit information to the request context
func WithAuditContext(r *http.Request, username string, orgID int64) *http.Request {
	ctx := context.WithValue(r.Context(), AuditUserKey, username)
	ctx = context.WithValue(ctx, AuditOrgIDKey, orgID)
	return r.WithContext(ctx)
}

// getUsernameFromContext extracts username from request context
func getUsernameFromContext(ctx context.Context) string {
	if username, ok := ctx.Value(AuditUserKey).(string); ok {
		return username
	}
	return "unknown"
}

// getOrgIDFromContext extracts organization ID from request context
func getOrgIDFromContext(ctx context.Context) int64 {
	if orgID, ok := ctx.Value(AuditOrgIDKey).(int64); ok {
		return orgID
	}
	return 0
}

// AuditMiddleware creates middleware that logs HTTP requests to the audit log
func AuditMiddleware(actionMapping map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get username and orgID from context (set by authentication middleware)
			username := getUsernameFromContext(r.Context())
			orgID := getOrgIDFromContext(r.Context())
			
			// Get path and method
			path := r.URL.Path
			method := r.Method
			
			// Find the matching action based on path patterns
			var actionString string
			for pattern, action := range actionMapping {
				parts := strings.Split(pattern, " ")
				if len(parts) != 2 {
					continue
				}
				
				methodPattern := parts[0]
				pathPattern := parts[1]
				
				// Check if method matches
				if methodPattern != "*" && methodPattern != method {
					continue
				}
				
				// Check if path matches (simple wildcard support)
				if strings.HasSuffix(pathPattern, "*") {
					prefix := strings.TrimSuffix(pathPattern, "*")
					if strings.HasPrefix(path, prefix) {
						actionString = action
						break
					}
				} else if pathPattern == path {
					actionString = action
					break
				}
			}
			
			// If no specific mapping, use a generic one
			if actionString == "" {
				actionString = fmt.Sprintf("%s request to %s", method, path)
			}

			// Create response wrapper to track status code
			rw := &responseWriter{ResponseWriter: w, statusCode: 200}
			
			// Process the request and record the time it takes
			startTime := time.Now()
			next.ServeHTTP(rw, r)
			duration := time.Since(startTime)
			
			// Build extra message including status code and duration
			extraMsg := fmt.Sprintf("Status: %d, Duration: %s", rw.statusCode, duration)
			
			// Log the audit event
			metadata := map[string]interface{}{
				"requestURI": r.RequestURI,
				"userAgent": r.UserAgent(),
				"remoteAddr": r.RemoteAddr,
				"statusCode": rw.statusCode,
				"durationMs": duration.Milliseconds(),
			}
			
			// Ignore errors here - we don't want to fail the request if logging fails
			_ = CreateAuditEvent(username, actionString, extraMsg, time.Now().Unix(), orgID, metadata)
		})
	}
}

// responseWriter is a wrapper around http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}