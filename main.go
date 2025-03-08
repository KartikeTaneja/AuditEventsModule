// example/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	
	"AuditEventsModule/audit" // Import your audit package
)

func main() {
	// Initialize the audit logger
	err := audit.InitAuditLogger(audit.FileLoggerType, map[string]string{
		"filePath": "audit.log",
	})
	if err != nil {
		log.Fatalf("Failed to initialize audit logger: %v", err)
	}
	
	// Set up HTTP handlers
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/indices/delete", handleDeleteIndex)
	http.HandleFunc("/dashboard/create", handleCreateDashboard)
	
	// Start the server
	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Parse username from form or JWT token
	username := r.FormValue("username")
	if username == "" {
		username = "unknown"
	}
	
	// Log the audit event
	err := audit.CreateAuditEvent(
		username,
		audit.ActionUserLogin,
		fmt.Sprintf("Login from %s", r.RemoteAddr),
		time.Now().Unix(),
		123, // OrgID, in a real app you'd get this from the user's context
		map[string]interface{}{
			"userAgent": r.UserAgent(),
			"ipAddress": r.RemoteAddr,
		},
	)
	if err != nil {
		log.Printf("Failed to log audit event: %v", err)
	}
	
	// Normal login processing would happen here
	w.Write([]byte("Login successful"))
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	// Similar to login, but for logout
	username := r.FormValue("username")
	if username == "" {
		username = "unknown"
	}
	
	err := audit.CreateAuditEvent(
		username,
		audit.ActionUserLogout,
		"User logged out",
		time.Now().Unix(),
		123,
		nil,
	)
	if err != nil {
		log.Printf("Failed to log audit event: %v", err)
	}
	
	w.Write([]byte("Logout successful"))
}

func handleDeleteIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Get parameters
	username := r.FormValue("username")
	indexName := r.FormValue("indexName")
	
	// Log the audit event
	err := audit.CreateAuditEvent(
		username,
		audit.ActionIndexDelete,
		fmt.Sprintf("Deleted index '%s'", indexName),
		time.Now().Unix(),
		123,
		map[string]interface{}{
			"indexName": indexName,
		},
	)
	if err != nil {
		log.Printf("Failed to log audit event: %v", err)
	}
	
	// Process the index deletion
	w.Write([]byte(fmt.Sprintf("Index '%s' deleted", indexName)))
}

func handleCreateDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Get parameters
	username := r.FormValue("username")
	dashboardName := r.FormValue("dashboardName")
	dashboardID := fmt.Sprintf("dashboard_%d", time.Now().Unix())
	
	// Log the audit event
	err := audit.CreateAuditEvent(
		username,
		audit.ActionDashboardCreate,
		fmt.Sprintf("Created dashboard '%s'", dashboardName),
		time.Now().Unix(),
		123,
		map[string]interface{}{
			"dashboardId": dashboardID,
			"dashboardName": dashboardName,
		},
	)
	if err != nil {
		log.Printf("Failed to log audit event: %v", err)
	}
	
	// Process the dashboard creation
	w.Write([]byte(fmt.Sprintf("Dashboard '%s' created with ID %s", dashboardName, dashboardID)))
}