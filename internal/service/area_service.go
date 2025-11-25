package service

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"

	"github.com/google/uuid"
)

type Point struct {
	X, Y float64
}

func IsPointInPolygon(p Point, polygon []Point) bool {
	intersectCount := 0
	j := len(polygon) - 1
	for i := 0; i < len(polygon); i++ {
		if (polygon[i].Y > p.Y) != (polygon[j].Y > p.Y) &&
			(p.X < (polygon[j].X-polygon[i].X)*(p.Y-polygon[i].Y)/(polygon[j].Y-polygon[i].Y)+polygon[i].X) {
			intersectCount++
		}
		j = i
	}
	return intersectCount%2 == 1
}

type areaService struct {
	areaRepo     repository.AreaRepository
	galleryRepo  repository.AreaGalleryRepository
	nodeRepo     repository.GraphRepository
	floorRepo    repository.FloorRepository
	revisionRepo repository.GraphRevisionRepository
}

func NewAreaService(
	aRepo repository.AreaRepository,
	gRepo repository.AreaGalleryRepository,
	nRepo repository.GraphRepository,
	fRepo repository.FloorRepository,
	rRepo repository.GraphRevisionRepository,
) AreaService {
	return &areaService{
		areaRepo:     aRepo,
		galleryRepo:  gRepo,
		nodeRepo:     nRepo,
		floorRepo:    fRepo,
		revisionRepo: rRepo,
	}
}

func (s *areaService) validateAreaInDraft(areaID uuid.UUID) error {
	area, err := s.areaRepo.GetByID(context.Background(), areaID)
	if err != nil {
		return err
	}

	// Get the floor to find the revision
	floor, err := s.floorRepo.GetByID(context.Background(), area.FloorID)
	if err != nil {
		return err
	}

	// Check if the revision is draft
	_, err = s.revisionRepo.GetDraftByFloorID(context.Background(), floor.ID)
	if err != nil {
		return errors.New("area is not in a draft revision")
	}

	return nil
}

func (s *areaService) CreateArea(ctx context.Context, req models.CreateAreaRequest) (*models.IDResponse, error) {
	// 1. Resolve GraphRevisionID from Floor
	if req.FloorID == nil {
		return nil, errors.New("floor_id is required")
	}
	floor, err := s.floorRepo.GetByID(ctx, *req.FloorID)
	if err != nil {
		return nil, errors.New("floor not found")
	}

	// 2. Parse Boundary points from DTO to Entity
	boundary := make(entity.Boundary, len(req.Boundary))
	for i, p := range req.Boundary {
		boundary[i] = entity.BoundaryPoint{X: p.X, Y: p.Y}
	}

	// 3. Mapping DTO -> Entity
	area := entity.Area{
		Name:            req.Name,
		Description:     req.Description,
		Category:        req.Category,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		Boundary:        boundary,
		GraphRevisionID: floor.GraphRevisionID,
		FloorID:         *req.FloorID,
		VenueID:         floor.VenueID, // Set venue ID from floor
		CoverImageID:    req.CoverImageID,
	}

	// 4. Save
	if err := s.areaRepo.Create(ctx, &area); err != nil {
		return nil, err
	}

	return &models.IDResponse{ID: area.ID}, nil
}

func (s *areaService) UpdateArea(ctx context.Context, id uuid.UUID, req models.UpdateAreaRequest) error {
	// Validate area is in draft revision
	if err := s.validateAreaInDraft(id); err != nil {
		return err
	}

	// 1. Get Existing
	area, err := s.areaRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("area not found")
	}

	// 2. Partial Update Logic - only update fields that are provided
	if req.Name != nil {
		area.Name = *req.Name
	}
	if req.Description != nil {
		area.Description = *req.Description
	}
	if req.Category != nil {
		area.Category = *req.Category
	}
	if req.Latitude != nil {
		area.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		area.Longitude = *req.Longitude
	}
	if req.Boundary != nil {
		// Parse Boundary points from DTO to Entity
		boundary := make(entity.Boundary, len(req.Boundary))
		for i, p := range req.Boundary {
			boundary[i] = entity.BoundaryPoint{X: p.X, Y: p.Y}
		}
		area.Boundary = boundary
	}
	if req.FloorID != nil {
		area.FloorID = *req.FloorID
	}
	if req.CoverImageID != nil {
		area.CoverImageID = req.CoverImageID
	}

	return s.areaRepo.Update(ctx, area)
}

func (s *areaService) DeleteArea(ctx context.Context, id uuid.UUID) error {
	// Validate area is in draft revision
	if err := s.validateAreaInDraft(id); err != nil {
		return err
	}

	return s.areaRepo.Delete(ctx, id)
}

// GetAreaDetail: Dipanggil saat user klik Pin di Peta Mobile App
func (s *areaService) GetAreaDetail(ctx context.Context, id uuid.UUID) (*models.AreaDetail, error) {
	// 1. Ambil Info Dasar
	area, err := s.areaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. Ambil Gallery (Foto-foto ruangan)
	galleryItems, _ := s.galleryRepo.GetByAreaID(ctx, id)
	var galleryDTOs []models.AreaGalleryDetail
	for _, item := range galleryItems {
		galleryDTOs = append(galleryDTOs, models.AreaGalleryDetail{
			MediaID:      item.MediaAssetID,
			URL:          item.MediaAsset.PublicURL,
			ThumbnailURL: item.MediaAsset.ThumbnailURL,
			Caption:      item.Caption,
			SortOrder:    item.SortOrder,
		})
	}

	// 3. Get Start Node ID (for 360° navigation starting point)
	return &models.AreaDetail{
		ID:          area.ID,
		Name:        area.Name,
		Description: area.Description,
		Gallery:     galleryDTOs,
		StartNodeID: area.StartNodeID,
	}, nil
}

// GetAreaEditorDetail: Dipanggil oleh Editor untuk mendapatkan detail lengkap area
func (s *areaService) GetAreaEditorDetail(ctx context.Context, id uuid.UUID) (*models.AreaEditorDetail, error) {
	// 1. Ambil Info Dasar dengan preload
	area, err := s.areaRepo.GetAreaWithDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. Ambil Gallery
	galleryItems, _ := s.galleryRepo.GetByAreaID(ctx, id)
	var galleryDTOs []models.AreaGalleryDetail
	for _, item := range galleryItems {
		galleryDTOs = append(galleryDTOs, models.AreaGalleryDetail{
			MediaID:      item.MediaAssetID,
			URL:          item.MediaAsset.PublicURL,
			ThumbnailURL: item.MediaAsset.ThumbnailURL,
			Caption:      item.Caption,
			SortOrder:    item.SortOrder,
		})
	}

	// 4. Convert boundary to model type
	var boundaryDTOs []models.BoundaryPoint
	for _, bp := range area.Boundary {
		boundaryDTOs = append(boundaryDTOs, models.BoundaryPoint{
			X: bp.X,
			Y: bp.Y,
		})
	}

	// 5. Get floor name
	floorName := ""
	if area.FloorID != uuid.Nil {
		if floor, err := s.floorRepo.GetByID(ctx, area.FloorID); err == nil {
			floorName = floor.Name
		}
	}

	// 6. Get cover image URL
	coverURL := ""
	if area.CoverImage != nil {
		coverURL = area.CoverImage.ThumbnailURL
		if coverURL == "" {
			coverURL = area.CoverImage.PublicURL
		}
	}

	return &models.AreaEditorDetail{
		ID:           area.ID,
		Name:         area.Name,
		Description:  area.Description,
		Category:     area.Category,
		Latitude:     area.Latitude,
		Longitude:    area.Longitude,
		Boundary:     boundaryDTOs,
		StartNodeID:  area.StartNodeID,
		FloorID:      area.FloorID,
		FloorName:    floorName,
		RevisionID:   area.GraphRevisionID,
		CoverImageID: area.CoverImageID,
		CoverURL:     coverURL,
		Gallery:      galleryDTOs,
		CreatedAt:    area.CreatedAt,
		UpdatedAt:    area.UpdatedAt,
	}, nil
}

// GetVenueAreas: List Pin untuk Peta Google Maps
func (s *areaService) GetVenueAreas(ctx context.Context, venueID uuid.UUID) ([]models.AreaPinDetail, error) {
	areas, err := s.areaRepo.GetByVenueID(ctx, venueID)
	if err != nil {
		return nil, err
	}

	var pins []models.AreaPinDetail
	for _, a := range areas {
		thumb := ""
		if a.CoverImage != nil {
			thumb = a.CoverImage.ThumbnailURL
		}

		pins = append(pins, models.AreaPinDetail{
			ID:           a.ID,
			Name:         a.Name,
			Category:     a.Category,
			Coordinates:  models.GeoPoint{Latitude: a.Latitude, Longitude: a.Longitude},
			ThumbnailURL: thumb,
			// FloorName bisa diambil jika preload floor
		})
	}
	return pins, nil
}

func (s *areaService) SetAreaStartNode(ctx context.Context, areaID uuid.UUID, req models.SetStartNodeRequest) (*models.SetStartNodeResponse, error) {
	// Validate area is in draft revision
	if err := s.validateAreaInDraft(areaID); err != nil {
		return nil, err
	}

	// 1. Get Area
	area, err := s.areaRepo.GetAreaWithDetails(ctx, areaID)
	if err != nil {
		return nil, errors.New("area not found")
	}

	// 2. Get Node
	node, err := s.nodeRepo.GetNodeByID(ctx, req.NodeID) // Assuming this method exists
	if err != nil {
		return nil, errors.New("node not found")
	}

	// 3. Validation 1: Ensure both are on the same Floor
	if area.FloorID != node.FloorID {
		return nil, errors.New("area and node must be on the same floor")
	}

	// 4. Validation 2: Run Ray Casting Algorithm
	polygon := make([]Point, len(area.Boundary))
	for i, bp := range area.Boundary {
		polygon[i] = Point{X: bp.X, Y: bp.Y}
	}
	nodePoint := Point{X: node.X, Y: node.Y}
	isInside := IsPointInPolygon(nodePoint, polygon)

	// 5. Update start_node_id regardless
	area.StartNodeID = &req.NodeID
	err = s.areaRepo.Update(ctx, area)
	if err != nil {
		return nil, err
	}

	// 6. Return warning if outside
	response := &models.SetStartNodeResponse{}
	if !isInside {
		response.Warning = "The selected node is outside the area boundary"
	}

	return response, nil
}

func (s *areaService) ListAreas(ctx context.Context, query models.AreaQuery) ([]models.AreaListItem, int64, error) {
	// Get areas with filtering and pagination
	areas, total, err := s.areaRepo.PagedAreas(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	items := make([]models.AreaListItem, 0, len(areas)) // Pre-allocate with capacity
	for _, area := range areas {
		// Get floor name
		floorName := ""
		if area.FloorID != uuid.Nil {
			if floor, err := s.floorRepo.GetByID(ctx, area.FloorID); err == nil {
				floorName = floor.Name
			}
		}

		// Get cover image URL
		coverURL := ""
		if area.CoverImage != nil {
			coverURL = area.CoverImage.ThumbnailURL
			if coverURL == "" {
				coverURL = area.CoverImage.PublicURL
			}
		}

		items = append(items, models.AreaListItem{
			ID:           area.ID,
			Name:         area.Name,
			Description:  area.Description,
			Category:     area.Category,
			FloorID:      area.FloorID,
			FloorName:    floorName,
			RevisionID:   area.GraphRevisionID,
			CoverURL:     coverURL,
			GalleryCount: len(area.Gallery),
			CreatedAt:    area.CreatedAt,
			UpdatedAt:    area.UpdatedAt,
		})
	}

	return items, total, nil
}
