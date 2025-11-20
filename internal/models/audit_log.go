package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLogFilter struct {
	OrganizationID string `json:"organization_id"`
	Entity         string `json:"entity"`
	EntityID       string `json:"entity_id"`
	Action         string `json:"action"`
	UserID         string `json:"user_id"`
	FromDate       string `json:"from_date"` // YYYY-MM-DD
	ToDate         string `json:"to_date"`
	UserAgent      string `json:"user_agent"`
	IPAddress      string `json:"ip_address"`
}

type AuditLogQuery struct {
	AuditLogFilter
	Limit  *int    `json:"limit"`
	Offset *int    `json:"offset"`
	Sort   *string `json:"sort"`
}

type AuditLogQueryCursor struct {
	AuditLogFilter
	Limit  *int    `json:"limit"`
	Cursor *string `json:"cursor,omitempty"`
}

type CreateAuditLogRequest struct {
	OrganizationID uuid.UUID              `json:"organization_id"`
	UserID         uuid.UUID              `json:"user_id"` // Pelaku
	Action         string                 `json:"action"`  // e.g. "NODE_CREATE"
	Entity         string                 `json:"entity"`  // e.g. "GraphNode"
	EntityID       string                 `json:"entity_id"`
	Details        map[string]interface{} `json:"details"` // Snapshot data JSON
	IPAddress      string                 `json:"ip_address"`
	UserAgent      string                 `json:"user_agent"`
}

type AuditLogResponse struct {
	ID             uint                   `json:"id"`
	CreatedAt      time.Time              `json:"created_at"`
	OrganizationID uuid.UUID              `json:"organization_id"`
	UserID         *uuid.UUID             `json:"user_id"`
	ActorName      string                 `json:"actor_name"`
	ActorEmail     string                 `json:"actor_email"`
	Action         string                 `json:"action"`
	Entity         string                 `json:"entity"`
	EntityID       string                 `json:"entity_id"`
	Details        map[string]interface{} `json:"details"`
	IPAddress      string                 `json:"ip_address"`
}

type AuditListResponse struct {
	Data       []AuditLogResponse `json:"data"`
	NextCursor string             `json:"next_cursor"`
	HasMore    bool               `json:"has_more"`
}
