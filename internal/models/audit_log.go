package models

import "time"

type AuditFilter struct {
	OrganizationID string `json:"organization_id"`
	Entity         string `json:"entity"`
	EntityID       string `json:"entity_id"`
	Action         string `json:"action"`
	UserID         string `json:"user_id"`
	FromDate       string `json:"from_date"` // YYYY-MM-DD
	ToDate         string `json:"to_date"`
	UserAgent    string `json:"user_agent"`
	IPAddress    string `json:"ip_address"`
}

type AuditQuery struct {
	AuditFilter
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}

type AuditQueryCursor struct {
	AuditFilter
	Limit  *int    `json:"limit"`
	Cursor *string `json:"cursor"`
}

type AuditLogResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	
	ActorName string    `json:"actor_name"` 
	ActorEmail string   `json:"actor_email"`
	
	Action    string    `json:"action"`     
	Entity    string    `json:"entity"`
	EntityID  string    `json:"entity_id"`
	
	Details   map[string]interface{} `json:"details"`
	IPAddress string    `json:"ip_address"`
}

type AuditListResponse struct {
	Data       []AuditLogResponse `json:"data"`
	NextCursor uint               `json:"next_cursor"`
	HasMore    bool               `json:"has_more"`
}