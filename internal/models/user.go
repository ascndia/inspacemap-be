package models

import (
	"github.com/google/uuid"
)

type UserDetail struct {
	ID             uuid.UUID    `json:"id"`
	Email          string       `json:"email"`
	FullName       string       `json:"full_name"`
	AvatarURL      string       `json:"avatar_url"`
	Organizations  []OrgMemberDetail `json:"organizations"`
}

type UserFilter struct {
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	Email          *string    `json:"email,omitempty"`
	FullName       *string    `json:"full_name,omitempty"`
	RoleID         *uint      `json:"role_id,omitempty"`
	IsEmailVerified *bool      `json:"is_email_verified,omitempty"`
}

type UserQuery struct {
	UserFilter
	Limit          *int       `json:"limit,omitempty"`
	Offset         *int       `json:"offset,omitempty"`
	Sort 		*string    `json:"sort,omitempty"`
}

type UserQueryCursor struct {
	UserFilter
	Limit          *int       `json:"limit,omitempty"`
	Cursor         *string    `json:"cursor,omitempty"`
}
