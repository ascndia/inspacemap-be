package service

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"

	"github.com/google/uuid"
)

type venueService struct {
	venueRepo repository.VenueRepository
}

func NewVenueService(vRepo repository.VenueRepository) VenueService {
	return &venueService{
		venueRepo: vRepo,
	}
}

// =================================================================
// 1. WRITE OPERATIONS
// =================================================================

func (s *venueService) CreateVenue(ctx context.Context, req models.CreateVenueRequest) (*models.IDResponse, error) {
	venue := entity.Venue{
		Name:         req.Name,
		Slug:         req.Slug,
		Description:  req.Description,
		Address:      req.Address,
		City:         req.City,
		Province:     req.Province,
		PostalCode:   req.PostalCode,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		CoverImageID: req.CoverImageID,
		Visibility:   entity.VisibilityPrivate, // Default
	}

	if req.Visibility != "" {
		venue.Visibility = entity.VisibilityStatus(req.Visibility)
	}

	// Init Gallery jika ada
	if len(req.Gallery) > 0 {
		for _, item := range req.Gallery {
			venue.Gallery = append(venue.Gallery, entity.VenueGalleryItem{
				MediaAssetID: item.MediaAssetID,
				Caption:      item.Caption,
				SortOrder:    item.SortOrder,
				IsVisible:    item.IsVisible,
				IsFeatured:   item.IsFeatured,
			})
		}
	}

	if err := s.venueRepo.Create(ctx, &venue); err != nil {
		return nil, err
	}

	return &models.IDResponse{ID: venue.ID}, nil
}

func (s *venueService) UpdateVenue(ctx context.Context, id uuid.UUID, req models.UpdateVenueRequest) error {
	// 1. Get Existing
	venue, err := s.venueRepo.GetByID(ctx, id) // Pakai BaseRepo GetByID cukup utk update field dasar
	if err != nil {
		return errors.New("venue not found")
	}

	// 2. Partial Update Logic
	if req.Name != nil {
		venue.Name = *req.Name
	}
	if req.Slug != nil {
		venue.Slug = *req.Slug
	}
	if req.Description != nil {
		venue.Description = *req.Description
	}
	if req.Address != nil {
		venue.Address = *req.Address
	}
	if req.City != nil {
		venue.City = *req.City
	}
	if req.Province != nil {
		venue.Province = *req.Province
	}
	if req.PostalCode != nil {
		venue.PostalCode = *req.PostalCode
	}
	if req.Latitude != nil {
		venue.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		venue.Longitude = *req.Longitude
	}
	if req.CoverImageID != nil {
		venue.CoverImageID = req.CoverImageID
	}
	if req.Visibility != nil {
		venue.Visibility = entity.VisibilityStatus(*req.Visibility)
	}

	// 3. Save
	return s.venueRepo.Update(ctx, venue)
}

func (s *venueService) DeleteVenue(ctx context.Context, id uuid.UUID) error {
	// Soft Delete via Repository
	return s.venueRepo.Delete(ctx, id)
}

// =================================================================
// 2. READ OPERATIONS (ADMIN / DASHBOARD)
// =================================================================

func (s *venueService) GetVenueDetail(ctx context.Context, id uuid.UUID) (*models.VenueDetail, error) {
	venue, err := s.venueRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapEntityToDetail(venue), nil
}

func (s *venueService) GetVenueBySlug(ctx context.Context, slug string) (*models.VenueDetail, error) {
	// Menggunakan Repo GetBySlug (pastikan repo implement ini dengan Preload)
	venue, err := s.venueRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	// Jika Repo GetBySlug belum preload gallery/poi, mungkin perlu dipanggil ulang GetVenueByID
	// Tapi idealnya Repo GetBySlug sudah preload.

	// Mapping manual atau via helper yang sama
	// Untuk amannya kita bisa panggil mapEntityToDetail jika struct entity venue sudah terisi relasinya
	return s.mapEntityToDetail(venue), nil
}

func (s *venueService) ListVenues(ctx context.Context, query models.VenueQuery) ([]models.VenueListItem, int64, error) {
	venues, total, err := s.venueRepo.PagedVenues(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	var list []models.VenueListItem
	for _, v := range venues {
		coverURL := ""
		if v.CoverImage != nil {
			coverURL = v.CoverImage.ThumbnailURL
		}

		list = append(list, models.VenueListItem{
			ID:            v.ID,
			Name:          v.Name,
			Slug:          v.Slug,
			City:          v.City,
			CoverImageURL: coverURL,
			Visibility:    string(v.Visibility),
			IsLive:        v.LiveRevisionID != uuid.Nil,
		})
	}

	return list, total, nil
}

// =================================================================
// 3. MOBILE APP CONSUMER
// =================================================================

func (s *venueService) GetMobileManifest(ctx context.Context, slug string) (*models.ManifestResponse, error) {
	venueEntity, err := s.venueRepo.GetLiveManifestData(slug)
	if err != nil {
		return nil, err
	}

	// Tentukan Start Node
	var startNodeID uuid.UUID
	if venueEntity.LiveRevision.StartNodeID != nil {
		startNodeID = *venueEntity.LiveRevision.StartNodeID
	} else if len(venueEntity.LiveRevision.Floors) > 0 && len(venueEntity.LiveRevision.Floors[0].Nodes) > 0 {
		startNodeID = venueEntity.LiveRevision.Floors[0].Nodes[0].ID
	}

	var floorDTOs []models.FloorData
	for _, floor := range venueEntity.LiveRevision.Floors {
		var nodeDTOs []models.NodeData
		for _, node := range floor.Nodes {

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

			var areaName string
			if node.Area != nil {
				areaName = node.Area.Name
			}

			panoURL := ""
			if node.Panorama != nil {
				panoURL = node.Panorama.PublicURL
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

	return &models.ManifestResponse{
		VenueID:     venueEntity.ID,
		VenueName:   venueEntity.Name,
		LastUpdated: venueEntity.LiveRevision.CreatedAt,
		StartNodeID: startNodeID,
		Floors:      floorDTOs,
	}, nil
}

func (s *venueService) mapEntityToDetail(venue *entity.Venue) *models.VenueDetail {
	var galleryDTOs []models.VenueGalleryDetail
	for _, item := range venue.Gallery {
		url, thumb := "", ""
		if item.MediaAsset.ID != uuid.Nil {
			url = item.MediaAsset.PublicURL
			thumb = item.MediaAsset.ThumbnailURL
		}
		galleryDTOs = append(galleryDTOs, models.VenueGalleryDetail{
			MediaID:      item.MediaAssetID,
			URL:          url,
			ThumbnailURL: thumb,
			Caption:      item.Caption,
			SortOrder:    item.SortOrder,
			IsFeatured:   item.IsFeatured,
		})
	}

	var poiDTOs []models.AreaPinDetail
	for _, area := range venue.PointsOfInterest {
		thumb := ""
		if area.CoverImage != nil {
			thumb = area.CoverImage.ThumbnailURL
		}
		poiDTOs = append(poiDTOs, models.AreaPinDetail{
			ID:           area.ID,
			Name:         area.Name,
			Category:     area.Category,
			Coordinates:  models.GeoPoint{Latitude: area.Latitude, Longitude: area.Longitude},
			ThumbnailURL: thumb,
		})
	}

	coverURL := ""
	if venue.CoverImage != nil {
		coverURL = venue.CoverImage.PublicURL
	}

	return &models.VenueDetail{
		ID:               venue.ID,
		OrganizationID:   venue.OrganizationID,
		Name:             venue.Name,
		Slug:             venue.Slug,
		Description:      venue.Description,
		Address:          venue.Address,
		City:             venue.City,
		Province:         venue.Province,
		PostalCode:       venue.PostalCode,
		FullAddress:      venue.Address + ", " + venue.City,
		Coordinates:      models.GeoPoint{Latitude: venue.Latitude, Longitude: venue.Longitude},
		Visibility:       string(venue.Visibility),
		CoverImageURL:    coverURL,
		Gallery:          galleryDTOs,
		PointsOfInterest: poiDTOs,
		CreatedAt:        venue.CreatedAt,
		UpdatedAt:        venue.UpdatedAt,
	}
}
