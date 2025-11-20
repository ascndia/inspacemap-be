package models

import "github.com/google/uuid"

type RevisionStatus string


type FilterGraphRevision struct {
	OrganizationID *uuid.UUID        `json:"organization_id,omitempty"`
	CreatedByID    *uuid.UUID        `json:"created_by_id,omitempty"`
	Status         *RevisionStatus  `json:"status,omitempty"`
	FloorID        *uuid.UUID        `json:"floor_id,omitempty"`
	Note		  *string           `json:"note,omitempty"`
	VenueID        *uuid.UUID        `json:"venue_id,omitempty"`
	CreatedAfter  *string           `json:"created_after,omitempty"`
	CreatedBefore *string           `json:"created_before,omitempty"`
}

type QueryGraphRevision struct {
	FilterGraphRevision
	Limit  *int    `json:"limit,omitempty"`
	Offset *int    `json:"offset,omitempty"`
	Sort   *string `json:"sort,omitempty"`
}

type CursorGraphRevisionQuery struct {
	FilterGraphRevision
	Limit  *int    `json:"limit,omitempty"`
	Cursor *string `json:"cursor,omitempty"`
}