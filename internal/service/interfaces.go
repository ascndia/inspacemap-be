package service

import (
	"context"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error)
	Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error)
}
type TeamService interface {
	CreateUser(ctx context.Context, orgID uuid.UUID, req models.CreateUserRequest) error
	RemoveMember(ctx context.Context, orgID uuid.UUID, targetUserID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, orgID uuid.UUID, req models.UpdateUserRoleRequest) error
	GetMembersList(ctx context.Context, orgID uuid.UUID) ([]models.TeamMemberDetail, error)
}

type RoleService interface {
	GetSystemRoles(ctx context.Context) ([]models.RoleDetail, error)
	GetAvailablePermissions(ctx context.Context) ([]models.PermissionNode, error)
}

type GraphService interface {
	CreateNode(ctx context.Context, req models.CreateNodeRequest) (*models.CreateNodeResponse, error)
	CreateFloor(ctx context.Context, revisionID uuid.UUID, userID uuid.UUID, req models.CreateFloorRequest) (*models.IDResponse, error)
	UpdateFloor(ctx context.Context, floorID uuid.UUID, req models.UpdateFloorRequest) error
	DeleteFloor(ctx context.Context, floorID uuid.UUID) error
	GetFloor(ctx context.Context, floorID uuid.UUID) (*models.FloorAdminDetail, error)
	GetFloors(ctx context.Context, revisionID uuid.UUID) ([]models.FloorAdminDetail, error)
	ConnectNodes(ctx context.Context, req models.ConnectNodesRequest) (*models.CreateEdgeResponse, error)
	UpdateNodePosition(ctx context.Context, nodeID uuid.UUID, req models.UpdateNodePositionRequest) error
	UpdateNodeCalibration(ctx context.Context, nodeID uuid.UUID, req models.UpdateNodeCalibrationRequest) error
	UpdateNode(ctx context.Context, nodeID uuid.UUID, req models.UpdateNodeRequest) error
	DeleteNode(ctx context.Context, nodeID uuid.UUID) error
	DeleteConnection(ctx context.Context, fromID, toID uuid.UUID) error
	GetEditorDataByRevision(ctx context.Context, revisionID uuid.UUID) (*models.ManifestResponse, error)
	PublishChanges(ctx context.Context, req models.PublishDraftRequest) error

	// Graph Revision Management
	CreateDraftRevision(ctx context.Context, venueID uuid.UUID, userID uuid.UUID) (*models.IDResponse, error)
	ListRevisions(ctx context.Context, venueID uuid.UUID) ([]models.RevisionHistoryItem, error)
	GetRevisionDetail(ctx context.Context, revisionID uuid.UUID) (*models.GraphRevisionDetail, error)
	UpdateRevision(ctx context.Context, revisionID uuid.UUID, req models.UpdateRevisionRequest) error
	DeleteRevision(ctx context.Context, revisionID uuid.UUID) error
	DeepCopyRevision(ctx context.Context, sourceRevisionID, targetVenueID uuid.UUID, note string) (*models.DeepCopyRevisionResponse, error)
}

type OrganizationService interface {
	GetDetailByID(ctx context.Context, id uuid.UUID) (*models.OrganizationDetail, error)
	GetDetailBySlug(ctx context.Context, slug string) (*models.OrganizationDetail, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, req models.UpdateOrganizationRequest) error
	ListOrganizations(ctx context.Context, query models.OrganizationQuery) ([]models.OrganizationDetail, int64, error)
	DeactivateOrganization(ctx context.Context, id uuid.UUID) error
}

type StorageProvider interface {
	GetPresignedPutURL(ctx context.Context, bucket, key, contentType string, expiry time.Duration) (string, error)
	DeleteObject(ctx context.Context, bucket, key string) error
}

type MediaService interface {
	InitDirectUpload(ctx context.Context, orgID uuid.UUID, req models.PresignedUploadRequest) (*models.PresignedUploadResponse, error)
	ConfirmUpload(ctx context.Context, req models.ConfirmUploadRequest) error
	GetAsset(ctx context.Context, id uuid.UUID) (*entity.MediaAsset, error)
	ListAssets(ctx context.Context, query models.MediaAssetQuery) ([]entity.MediaAsset, int64, error)
	DeleteAsset(ctx context.Context, id uuid.UUID) error
}

type AreaService interface {
	CreateArea(ctx context.Context, req models.CreateAreaRequest) (*models.IDResponse, error)
	UpdateArea(ctx context.Context, id uuid.UUID, req models.UpdateAreaRequest) error
	DeleteArea(ctx context.Context, id uuid.UUID) error
	ListAreas(ctx context.Context, query models.AreaQuery) ([]models.AreaListItem, int64, error)
	GetAreaDetail(ctx context.Context, id uuid.UUID) (*models.AreaDetail, error)
	GetAreaEditorDetail(ctx context.Context, id uuid.UUID) (*models.AreaEditorDetail, error)
	GetVenueAreas(ctx context.Context, venueID uuid.UUID) ([]models.AreaPinDetail, error)
	SetAreaStartNode(ctx context.Context, areaID uuid.UUID, req models.SetStartNodeRequest) (*models.SetStartNodeResponse, error)
}

type VenueService interface {
	CreateVenue(ctx context.Context, req models.CreateVenueRequest) (*models.IDResponse, error)
	UpdateVenue(ctx context.Context, id uuid.UUID, req models.UpdateVenueRequest) error
	DeleteVenue(ctx context.Context, id uuid.UUID) error
	GetVenueDetail(ctx context.Context, id uuid.UUID) (*models.VenueDetail, error)
	GetVenueBySlug(ctx context.Context, slug string) (*models.VenueDetail, error)
	ListVenues(ctx context.Context, query models.VenueQuery) ([]models.VenueListItem, int64, error)
	GetMobileManifest(ctx context.Context, orgSlug, venueSlug string) (*models.MobileManifest, error)
}
type VenueGalleryService interface {
	ReorderGallery(ctx context.Context, req models.ReorderVenueGalleryRequest) error
	AddGalleryItems(ctx context.Context, req models.AddGalleryVenueItemsRequest) error
	UpdateGalleryItem(ctx context.Context, req models.UpdateVenueGalleryItemRequest) error
	RemoveGalleryItem(ctx context.Context, targetID, mediaID uuid.UUID) error
}

type AreaGalleryService interface {
	ReorderGallery(ctx context.Context, req models.ReorderAreaGalleryRequest) error
	AddGalleryItems(ctx context.Context, req models.AddAreaGalleryItemsRequest) error
	UpdateGalleryItem(ctx context.Context, req models.UpdateAreaGalleryItemRequest) error
	RemoveGalleryItem(ctx context.Context, targetID, mediaID uuid.UUID) error
	GetGalleryItems(ctx context.Context, areaID uuid.UUID) ([]models.AreaGalleryDetail, error)
}

type AuditService interface {
	GetActivityLogs(ctx context.Context, orgID uuid.UUID, query models.AuditLogQueryCursor) (*models.AuditListResponse, error)
	LogActivity(ctx context.Context, payload models.CreateAuditLogRequest)
}
