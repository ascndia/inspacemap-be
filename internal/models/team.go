package models

import (
	"time"

	"github.com/google/uuid"
)

type TeamMemberDetail struct {
	UserID    uuid.UUID `json:"user_id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url,omitempty"`

	RoleID   uuid.UUID `json:"role_id"`
	RoleName string    `json:"role_name"` // e.g. "Admin", "Editor"

	JoinedAt time.Time `json:"joined_at"`
}
