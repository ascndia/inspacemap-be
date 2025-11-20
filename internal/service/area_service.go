package service

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"

	"github.com/google/uuid"
)

type areaService struct {
	areaRepo    repository.AreaRepository
	galleryRepo repository.AreaGalleryRepository
	nodeRepo    repository.GraphRepository // Butuh ini untuk cari Nearest Node
}

func NewAreaService(
	aRepo repository.AreaRepository,
	gRepo repository.AreaGalleryRepository,
	nRepo repository.GraphRepository, // Asumsi ada method find nearest
) AreaService {
	return &areaService{
		areaRepo:    aRepo,
		galleryRepo: gRepo,
		nodeRepo:    nRepo,
	}
}

func (s *areaService) CreateArea(ctx context.Context, req models.CreateAreaRequest) (*models.IDResponse, error) {
	// TODO: Validasi VenueID (biasanya dari URL param di handler, lalu inject ke struct req atau argumen terpisah)
	// Asumsi req.VenueID sudah diisi dari handler

	// 1. Mapping DTO -> Entity
	area := entity.Area{
		Name:         req.Name,
		Description:  req.Description,
		Category:     req.Category,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		MapX:         req.MapX,
		MapY:         req.MapY,
		CoverImageID: req.CoverImageID,
		// FloorID & VenueID diisi dari context/req
	}
	if req.FloorID != nil {
		area.FloorID = *req.FloorID // uuid pointer
	}

	// 2. Save
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

	// 2. Update Fields
	area.Name = req.Name
	area.Description = req.Description
	area.Category = req.Category
	area.Latitude = req.Latitude
	area.Longitude = req.Longitude
	area.MapX = req.MapX
	area.MapY = req.MapY
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
