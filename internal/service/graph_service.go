package service

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"

	"github.com/google/uuid"
)

type graphService struct {
	graphRepo    repository.GraphRepository
	revisionRepo repository.GraphRevisionRepository
	floorRepo    repository.FloorRepository
	venueRepo    repository.VenueRepository
}

func NewGraphService(
	gRepo repository.GraphRepository,
	rRepo repository.GraphRevisionRepository,
	fRepo repository.FloorRepository,
	vRepo repository.VenueRepository,
) GraphService {
	return &graphService{
		graphRepo:    gRepo,
		revisionRepo: rRepo,
		floorRepo:    fRepo,
		venueRepo:    vRepo,
	}
}

// =================================================================
// 1. FLOOR MANAGEMENT
// =================================================================

func (s *graphService) CreateFloor(ctx context.Context, venueID uuid.UUID, req models.CreateFloorRequest) (*models.IDResponse, error) {
	// A. Pastikan kita bekerja di DRAFT Revision
	draft, err := s.revisionRepo.GetDraftByVenueID(ctx, venueID)
	if err != nil {
		// Auto-create draft jika belum ada (Lazy init)
		draft, err = s.revisionRepo.CreateDraft(ctx, venueID)
		if err != nil {
			return nil, err
		}
	}

	floor := entity.Floor{
		GraphRevisionID: draft.ID,
		VenueID:         venueID,
		Name:            req.Name,
		LevelIndex:      req.LevelIndex,
		MapImageID:      req.MapImageID,
		MapWidth:        req.MapWidth,
		MapHeight:       req.MapHeight,
		PixelsPerMeter:  req.PixelsPerMeter,
	}

	if err := s.floorRepo.Create(ctx, &floor); err != nil {
		return nil, err
	}

	return &models.IDResponse{ID: floor.ID}, nil
}

func (s *graphService) UpdateFloorMap(ctx context.Context, floorID uuid.UUID, req models.UpdateFloorRequest) error {
	// A. Security Check: Pastikan Floor ini milik DRAFT
	if _, err := s.revisionRepo.GetDraftByFloorID(ctx, floorID); err != nil {
		return errors.New("cannot edit floor: it belongs to a published version or does not exist")
	}

	// B. Update
	// Konversi DTO ke param repository
	return s.floorRepo.UpdateFloorMap(ctx, floorID, req.MapImageID, req.PixelsPerMeter)
}

// =================================================================
// 2. NODE OPERATIONS
// =================================================================

func (s *graphService) CreateNode(ctx context.Context, req models.CreateNodeRequest) (*models.IDResponse, error) {
	// A. Security Check: Floor harus ada di Draft
	if _, err := s.revisionRepo.GetDraftByFloorID(ctx, req.FloorID); err != nil {
		return nil, errors.New("cannot create node: target floor is not in draft mode")
	}

	// B. Validasi Koordinat
	if req.X < 0 || req.Y < 0 {
		return nil, errors.New("coordinates cannot be negative")
	}

	// C. Create Entity
	node := entity.GraphNode{
		FloorID:         req.FloorID,
		X:               req.X,
		Y:               req.Y,
		PanoramaAssetID: req.PanoramaAssetID,
		Label:           req.Label,
		// AreaID opsional, bisa null
	}

	if err := s.graphRepo.CreateNode(ctx, &node); err != nil {
		return nil, err
	}

	return &models.IDResponse{ID: node.ID}, nil
}

func (s *graphService) UpdateNodePosition(ctx context.Context, nodeID uuid.UUID, req models.UpdateNodePositionRequest) error {
	// Validasi ownership draft bisa dilakukan dengan query node -> floor -> revision
	// Untuk performa, kita asumsikan UI sudah membatasi, atau bisa tambah check di sini.
	return s.graphRepo.UpdateNodePosition(ctx, nodeID, req.X, req.Y)
}

func (s *graphService) UpdateNodeCalibration(ctx context.Context, nodeID uuid.UUID, req models.UpdateNodeCalibrationRequest) error {
	return s.graphRepo.UpdateNodeCalibration(ctx, nodeID, req.RotationOffset)
}

func (s *graphService) DeleteNode(ctx context.Context, nodeID uuid.UUID) error {
	return s.graphRepo.DeleteNode(ctx, nodeID)
}

// =================================================================
// 3. EDGE OPERATIONS (CONNECTING)
// =================================================================

func (s *graphService) ConnectNodes(ctx context.Context, req models.ConnectNodesRequest) error {
	if req.FromNodeID == req.ToNodeID {
		return errors.New("cannot connect node to itself")
	}

	// A. Validasi Cross-Graph (Sudah ada di Repo level, tapi bisa double check di sini)
	// Repository graphRepo.ConnectNodes sudah kita pasang logic kalkulasi Heading & Distance.

	edge := entity.GraphEdge{
		FromNodeID: req.FromNodeID,
		ToNodeID:   req.ToNodeID,
		Type:       "walk", // Default type
	}

	// Buat koneksi satu arah
	if err := s.graphRepo.ConnectNodes(ctx, &edge); err != nil {
		return err
	}

	// Opsional: Jika ingin Bi-Directional (Dua arah otomatis)
	// edgeBack := entity.GraphEdge{ FromNodeID: req.ToNodeID, ToNodeID: req.FromNodeID, ... }
	// s.graphRepo.ConnectNodes(ctx, &edgeBack)

	return nil
}

func (s *graphService) DeleteConnection(ctx context.Context, fromID, toID uuid.UUID) error {
	return s.graphRepo.DeleteEdge(ctx, fromID, toID)
}

// =================================================================
// 4. WORKFLOW (EDITOR & PUBLISH)
// =================================================================

// GetEditorData: Mengambil data DRAFT lengkap untuk ditampilkan di Canvas Web Admin
func (s *graphService) GetEditorData(ctx context.Context, venueID uuid.UUID) (*models.ManifestResponse, error) {
	// 1. Ambil Draft (Auto-create jika tidak ada)
	draft, err := s.revisionRepo.GetDraftByVenueID(ctx, venueID)
	if err != nil {
		draft, err = s.revisionRepo.CreateDraft(ctx, venueID)
		if err != nil {
			return nil, err
		}
		// Draft baru pasti kosong, return struktur kosong
		// (Atau reload draft yang baru dibuat untuk memastikan relasi terload)
		// return &model.ManifestResponse{...}, nil
	}

	// 2. Mapping Entity GraphRevision -> DTO ManifestResponse
	// Logic mapping ini mirip dengan GetMobileManifest, tapi sumber datanya adalah DRAFT

	var startNodeID uuid.UUID
	if draft.StartNodeID != nil {
		startNodeID = *draft.StartNodeID
	}

	var floorDTOs []models.FloorData
	for _, floor := range draft.Floors {
		var nodeDTOs []models.NodeData
		for _, node := range floor.Nodes {
			// Mapping Neighbors
			var neighborDTOs []models.NeighborData
			for _, edge := range node.OutgoingEdges {
				neighborDTOs = append(neighborDTOs, models.NeighborData{
					TargetNodeID: edge.ToNodeID,
					Heading:      edge.Heading,
					Distance:     edge.Distance,
					Type:         edge.Type,
					IsActive:     edge.IsActive,
				})
			}

			// Resolve Area Name
			var areaName string
			if node.Area != nil {
				areaName = node.Area.Name
			}

			// Resolve Panorama URL (Safety check jika asset terhapus)
			panoURL := ""
			if node.Panorama != nil {
				panoURL = node.Panorama.ThumbnailURL // Gunakan thumbnail untuk editor agar ringan!
			}

			nodeDTOs = append(nodeDTOs, models.NodeData{
				ID:             node.ID,
				X:              int(node.X),
				Y:              int(node.Y),
				PanoramaURL:    panoURL,
				RotationOffset: node.RotationOffset,
				AreaID:         node.AreaID,
				AreaName:       areaName,
				Neighbors:      neighborDTOs,
			})
		}

		mapURL := ""
		if floor.MapImage != nil {
			mapURL = floor.MapImage.PublicURL
		}

		floorDTOs = append(floorDTOs, models.FloorData{
			ID:          floor.ID,
			LevelName:   floor.Name,
			LevelIndex:  floor.LevelIndex,
			MapImageURL: mapURL,
			MapWidth:    floor.MapWidth,
			MapHeight:   floor.MapHeight,
			Nodes:       nodeDTOs,
		})
	}

	// Ambil nama venue
	venueName := ""
	if venue, err := s.venueRepo.GetByID(ctx, venueID); err == nil {
		venueName = venue.Name
	}

	return &models.ManifestResponse{
		VenueID:     venueID,
		VenueName:   venueName,
		LastUpdated: draft.CreatedAt,
		StartNodeID: startNodeID,
		Floors:      floorDTOs,
	}, nil
}

func (s *graphService) PublishChanges(ctx context.Context, venueID uuid.UUID, req models.PublishDraftRequest) error {
	// Panggil Repository untuk melakukan Deep Copy Transaction
	return s.revisionRepo.PublishDraft(ctx, venueID, req.Note)
}
