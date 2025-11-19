package repository

import (
	"context"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"

	"github.com/google/uuid"
)

type BaseRepository[T any, K comparable] interface {
	GetByID(ctx context.Context, id K) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id K) error
}

type UserRepository interface {
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.User, error)
	FilterUsers(ctx context.Context, filter models.UserFilter) ([]entity.User, error)
	PagedUsers(ctx context.Context, query models.UserQuery) ([]entity.User, int64, error)
	CursorUsers(ctx context.Context, query models.UserQueryCursor) ([]entity.User, string, error)
	BaseRepository[entity.User, uuid.UUID]
}

type OrganizationRepository interface {
	BaseRepository[entity.Organization, uuid.UUID]
	GetByDomain(ctx context.Context, domain string) (*entity.Organization, error)
	FilterOrganizations(ctx context.Context, filter models.OrganizationFilter) ([]entity.Organization, error)
	PagedOrganizations(ctx context.Context, query models.OrganizationQuery) ([]entity.Organization, int64, error)
	CursorOrganizations(ctx context.Context, query models.OrganizationQueryCursor) ([]entity.Organization,  string, error)
}

type UserInvitationRepository interface {
	BaseRepository[entity.UserInvitation, uuid.UUID]
	GetByToken(ctx context.Context, token string) (*entity.UserInvitation, error)
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.UserInvitation, error)
	GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]entity.UserInvitation, error)
	GetByEmail(ctx context.Context, email string) ([]entity.UserInvitation, error)
	GetByInviterID(ctx context.Context, inviterID uuid.UUID) ([]entity.UserInvitation, error)
	GetByStatus(ctx context.Context, status string) ([]entity.UserInvitation, error)
}

type RoleRepository interface {
	BaseRepository[entity.Role, uuid.UUID]
    GetByName(ctx context.Context, name string) (*entity.Role, error)
    AttachPermission(ctx context.Context, roleID uuid.UUID, permID uuid.UUID) error
    DetachPermission(ctx context.Context, roleID uuid.UUID, permID uuid.UUID) error
    GetPermissions(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error)
}

type PermissionRepository interface {
	BaseRepository[entity.Permission, uuid.UUID]
	GetByName(ctx context.Context, name string) (*entity.Permission, error)
	GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error)
}

type AuthRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	FindUserByProvider(ctx context.Context, provider string, uid string) (*entity.User, error)
	GetOrganizationSSOConfig(ctx context.Context, domain string) (*entity.OrganizationSSO, error)
	ValidateAPIKey(ctx context.Context, keyHash string) (*entity.ApiKey, error)
}

type VenueRepository interface {
	BaseRepository[entity.Venue, uuid.UUID]
    FilterVenues(ctx context.Context, filter models.VenueFilter) ([]entity.Venue, error)
	PagedVenues(ctx context.Context, query models.VenueQuery) ([]entity.Venue, error)
	CursorVenues(ctx context.Context, query models.VenueQueryCursor) ([]entity.Venue, error)
}

type GraphRepository interface {
	CreateNode(ctx context.Context, node *entity.GraphNode) error
	UpdateNodePosition(ctx context.Context, id uuid.UUID, x, y float64) error
	UpdateNodeCalibration(ctx context.Context, id uuid.UUID, offset float64) error
	DeleteNode(ctx context.Context, id uuid.UUID) error
	ConnectNodes(ctx context.Context, edge *entity.GraphEdge) error
	DeleteEdge(ctx context.Context, fromID, toID uuid.UUID) error
}

type GraphRevisionRepository interface {
	BaseRepository[entity.GraphRevision, uuid.UUID]
	PublishDraft(ctx context.Context, revisionID uuid.UUID, note string) error
	GetDraftByFloorID(ctx context.Context, floorID uuid.UUID) (*entity.GraphRevision, error)
	GetDraftByVenueID(ctx context.Context, venueID uuid.UUID) ([]entity.GraphRevision, error)
	GetDraftByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.GraphRevision, error)
	GetLiveByFloorID(ctx context.Context, floorID uuid.UUID) (*entity.GraphRevision, error)
	GetLiveByVenueID(ctx context.Context, venueID uuid.UUID) ([]entity.GraphRevision, error)
	GetLiveByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.GraphRevision, error)

	GetByFloorID(ctx context.Context, floorID uuid.UUID) ([]entity.GraphRevision, error)
	GetByVenueID(ctx context.Context, venueID uuid.UUID) ([]entity.GraphRevision, error)
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.GraphRevision, error)

	FilterGraphRevisions(ctx context.Context, filter models.FilterGraphRevision) ([]entity.GraphRevision, error)
	PagedGraphRevisions(ctx context.Context, query models.QueryGraphRevision) ([]entity.GraphRevision, error)
	CursorGraphRevisions(ctx context.Context, query models.CursorGraphRevisionQuery) ([]entity.GraphRevision, error)
}

type AreaRepository interface {
	BaseRepository[entity.Area, uuid.UUID]
	GetByVenueID(ctx context.Context, venueID uuid.UUID) ([]entity.Area, error)
	GetByFloorID(ctx context.Context, floorID uuid.UUID) ([]entity.Area, error)
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.Area, error)
	FilterAreas(ctx context.Context, filter models.AreaFilter) ([]entity.Area, error)
	PagedAreas(ctx context.Context, query models.AreaQuery) ([]entity.Area, error)
	CursorAreas(ctx context.Context, query models.AreaQueryCursor) ([]entity.Area, error)
}

type FloorRepository interface {
	BaseRepository[entity.Floor, uuid.UUID]
	GetByVenueID(ctx context.Context, venueID uuid.UUID) ([]entity.Floor, error)
	GetByGraphRevisionID(ctx context.Context, revisionID uuid.UUID) ([]entity.Floor, error)
	UpdateFloorMap(ctx context.Context, id uuid.UUID, mapImageID *uuid.UUID, pixelsPerMeter float64) error
}

type MediaAssetRepository interface {
	BaseRepository[entity.MediaAsset, uuid.UUID]
	FilterMediaAssets(ctx context.Context, filter models.MediaAssetFilter) ([]entity.MediaAsset, error)
	PagedMediaAssets(ctx context.Context, query models.MediaAssetQuery) ([]entity.MediaAsset, error)
	CursorMediaAssets(ctx context.Context, query models.MediaAssetQueryCursor) ([]entity.MediaAsset, error)

}

type AuditLogRepository interface {
	BaseRepository[entity.AuditLog, uuid.UUID]
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.AuditLog, error)
}