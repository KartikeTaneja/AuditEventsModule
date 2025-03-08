// audit/actions.go
package audit

// Action constants for common user actions
const (
	// Auth-related actions
	ActionUserLogin       = "User logged in"
	ActionUserLogout      = "User logged out"
	ActionUserPasswordReset = "User password reset"
	
	// Index-related actions
	ActionIndexCreate     = "Index created"
	ActionIndexDelete     = "Index deleted"
	ActionIndexUpdate     = "Index updated"
	
	// Organization-related actions
	ActionOrgSettingsUpdate = "Organization settings updated"
	
	// Dashboard-related actions
	ActionDashboardCreate = "Dashboard created"
	ActionDashboardUpdate = "Dashboard updated"
	ActionDashboardDelete = "Dashboard deleted"
	ActionDashboardFavorite = "Dashboard added to favorites"
	ActionDashboardUnfavorite = "Dashboard removed from favorites"
	
	// Saved Queries
	ActionSavedQueryCreate = "Saved query created"
	ActionSavedQueryUpdate = "Saved query updated"
	ActionSavedQueryDelete = "Saved query deleted"
	
	// Folder-related actions
	ActionFolderCreate    = "Folder created"
	ActionFolderUpdate    = "Folder updated"
	ActionFolderDelete    = "Folder deleted"
	
	// Alerts-related actions
	ActionAlertCreate     = "Alert created"
	ActionAlertUpdate     = "Alert updated"
	ActionAlertDelete     = "Alert deleted"
	ActionContactPointCreate = "Contact point created"
	ActionContactPointUpdate = "Contact point updated"
	ActionContactPointDelete = "Contact point deleted"
	
	// Lookup files
	ActionLookupFileCreate = "Lookup file created"
	ActionLookupFileDelete = "Lookup file deleted"
)