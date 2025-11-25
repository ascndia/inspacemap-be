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
	areaRepo    repository.AreaRepository
	galleryRepo repository.AreaGalleryRepository
	nodeRepo    repository.GraphRepository
	floorRepo   repository.FloorRepository // Need to get revision from floor
}

func NewAreaService(
	aRepo repository.AreaRepository,
	gRepo repository.AreaGalleryRepository,
	nRepo repository.GraphRepository,
	fRepo repository.FloorRepository,
) AreaService {
	return &areaService{
		areaRepo:    aRepo,
		galleryRepo: gRepo,
		nodeRepo:    nRepo,
		floorRepo:   fRepo,
	}
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
		CoverImageID:    req.CoverImageID,
	}

	// 4. Save
	if err := s.areaRepo.Create(ctx, &area); err != nil {
		return nil, err
	}

	return &models.IDResponse{ID: area.ID}, nil
}

func (s *areaService) UpdateArea(ctx context.Context, id uuid.UUID, req models.CreateAreaRequest) error {
	// 1. Get Existing
	area, err := s.areaRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("area not found")
	}

	// 2. Parse Boundary
	boundary := make(entity.Boundary, len(req.Boundary))
	for i, p := range req.Boundary {
		boundary[i] = entity.BoundaryPoint{X: p.X, Y: p.Y}
	}

	// 3. Update Fields
	area.Name = req.Name
	area.Description = req.Description
	area.Category = req.Category
	area.Latitude = req.Latitude
	area.Longitude = req.Longitude
	area.Boundary = boundary
	area.CoverImageID = req.CoverImageID
	if req.FloorID != nil {
		area.FloorID = *req.FloorID
	}

	return s.areaRepo.Update(ctx, area)
}

func (s *areaService) DeleteArea(ctx context.Context, id uuid.UUID) error {
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

	// 3. Cari Node Terdekat (Start Point untuk 360)
	// Logic: Cari node yang punya AreaID == id ini. Ambil yang pertama.
	// (Perlu implementasi query di GraphRepo: GetOneNodeByAreaID)
	var nearestNodeID *uuid.UUID
	// node, _ := s.nodeRepo.GetOneByAreaID(ctx, id)
	// if node != nil { nearestNodeID = &node.ID }

	return &models.AreaDetail{
		ID:            area.ID,
		Name:          area.Name,
		Description:   area.Description,
		Gallery:       galleryDTOs,
		NearestNodeID: nearestNodeID,
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
