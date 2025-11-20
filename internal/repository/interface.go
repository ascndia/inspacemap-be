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
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
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

type OrganizationMemberRepository interface {
	BaseRepository[entity.OrganizationMember, uuid.UUID] // ID tabel ini biasanya uint (pivot)
	AddMember(ctx context.Context, member *entity.OrganizationMember) error
	RemoveMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error
	UpdateRole(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, roleID uuid.UUID) error
	GetMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (*entity.OrganizationMember, error)
	GetMembersByOrg(ctx context.Context, orgID uuid.UUID) ([]entity.OrganizationMember, error)
	GetMembersByUser(ctx context.Context, userID uuid.UUID) ([]entity.OrganizationMember, error)
}

type UserInvitationRepository interface {
	BaseRepository[entity.UserInvitation, uuid.UUID]
	GetByToken(ctx context.Context, token string) (*entity.UserInvitation, error)
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.UserInvitation, error)
	GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]entity.UserInvitation, error)
	GetByEmail(ctx context.Context, email string) ([]entity.UserInvitation, error)
	GetByInviterID(ctx context.Context, inviterID uuid.UUID) ([]entity.UserInvitation, error)
	GetByStatus(ctx context.Context, orgID uuid.UUID, status string) ([]entity.UserInvitation, error)
	RevokeInvitation(ctx context.Context, id uuid.UUID) error
}

type RoleRepository interface {
	BaseRepository[entity.Role, uuid.UUID]
    GetByName(ctx context.Context, name string) (*entity.Role, error)
	GetByOrganizationIDAndName(ctx context.Context, orgID *uuid.UUID, name string) (*entity.Role, error)
    AttachPermission(ctx context.Context, roleID uuid.UUID, permID uuid.UUID) error
    DetachPermission(ctx context.Context, roleID uuid.UUID, permID uuid.UUID) error
    GetPermissions(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error)
}

type PermissionRepository interface {
	BaseRepository[entity.Permission, uuid.UUID]
	GetByKey(ctx context.Context, key string) (*entity.Permission, error)
	GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error)
}

type AuthRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	ValidateAPIKey(ctx context.Context, keyHash string) (*entity.ApiKey, error)
}

type VenueRepository interface {
	BaseRepository[entity.Venue, uuid.UUID]
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.Venue, error)
    FilterVenues(ctx context.Context, filter models.VenueFilter) ([]entity.Venue, error)
	PagedVenues(ctx context.Context, query models.VenueQuery) ([]entity.Venue, int64, error)
	CursorVenues(ctx context.Context, query models.VenueQueryCursor) ([]entity.Venue, string, error)
}

type VenueGalleryRepository interface {
	BaseRepository[entity.VenueGalleryItem, uuid.UUID]
	GetByVenueID(ctx context.Context, venueID uuid.UUID) ([]entity.VenueGalleryItem, error)
	FilterVenueGalleries(ctx context.Context, filter models.VenueGalleryFilter) ([]entity.VenueGalleryItem, error)
	PagedVenueGalleries(ctx context.Context, query models.VenueGalleryQuery) ([]entity.VenueGalleryItem, int64, error)
	CursorVenueGalleries(ctx context.Context, query models.VenueGalleryCursor) ([]entity.VenueGalleryItem, string, error)
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
	PagedAreas(ctx context.Context, query models.AreaQuery) ([]entity.Area, int64, error)
	CursorAreas(ctx context.Context, query models.AreaQueryCursor) ([]entity.Area, string, error)
}

type AreaGalleryRepository interface {
	BaseRepository[entity.AreaGalleryItem, uuid.UUID]
	GetByAreaID(ctx context.Context, areaID uuid.UUID) ([]entity.AreaGalleryItem, error)
	GetByVenueID(ctx context.Context, venueID uuid.UUID) ([]entity.AreaGalleryItem, error)
	FilterAreaGalleries(ctx context.Context, filter models.AreaGalleryFilter) ([]entity.AreaGalleryItem, error)
	PagedAreaGalleries(ctx context.Context, query models.AreaGalleryQuery) ([]entity.AreaGalleryItem, int64, error)
	CursorAreaGalleries(ctx context.Context, query models.AreaGalleryCursor) ([]entity.AreaGalleryItem, string, error)
}

type FloorRepository interface {
	BaseRepository[entity.Floor, uuid.UUID]
	GetByVenueID(ctx context.Context, venueID uuid.UUID) ([]entity.Floor, error)
	GetByGraphRevisionID(ctx context.Context, revisionID uuid.UUID) ([]entity.Floor, error)
	UpdateFloorMap(ctx context.Context, id uuid.UUID, mapImageID *uuid.UUID, pixelsPerMeter float64) error
	FilterFloors(ctx context.Context, filter models.FloorFilter) ([]entity.Floor, error)
	PagedFloors(ctx context.Context, query models.FloorQuery) ([]entity.Floor, int64, error)
	CursorFloors(ctx context.Context, query models.FloorQueryCursor) ([]entity.Floor, string, error)
}

type MediaAssetRepository interface {
	BaseRepository[entity.MediaAsset, uuid.UUID]
	FilterMediaAssets(ctx context.Context, filter models.MediaAssetFilter) ([]entity.MediaAsset, error)
	PagedMediaAssets(ctx context.Context, query models.MediaAssetQuery) ([]entity.MediaAsset, int64, error)
	CursorMediaAssets(ctx context.Context, query models.MediaAssetQueryCursor) ([]entity.MediaAsset, string, error)
}

type AuditLogRepository interface {
	BaseRepository[entity.AuditLog, uint]
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.AuditLog, error)
	FilterAuditLogs(ctx context.Context, filter models.AuditLogFilter) ([]entity.AuditLog, error)
	PagedAuditLogs(ctx context.Context, query models.AuditLogQuery) ([]entity.AuditLog, int64, error)
	CursorAuditLogs(ctx context.Context, query models.AuditLogQueryCursor) ([]entity.AuditLog, string, error)
}