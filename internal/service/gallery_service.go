package service

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"

	"github.com/google/uuid"
)

// =================================================================
// VENUE GALLERY SERVICE
// =================================================================

type venueGalleryService struct {
	repo repository.VenueGalleryRepository
}

func NewVenueGalleryService(repo repository.VenueGalleryRepository) VenueGalleryService {
	return &venueGalleryService{repo: repo}
}

func (s *venueGalleryService) AddGalleryItems(ctx context.Context, req models.AddGalleryVenueItemsRequest) error {
	var items []entity.VenueGalleryItem
	for _, item := range req.Items {
		items = append(items, entity.VenueGalleryItem{
			VenueID:      req.VenueID,
			MediaAssetID: item.MediaAssetID,
			Caption:      item.Caption,
			SortOrder:    item.SortOrder,
			IsVisible:    item.IsVisible,
			IsFeatured:   item.IsFeatured,
		})
	}
	return s.repo.AddVenueItems(ctx, items)
}

func (s *venueGalleryService) ReorderGallery(ctx context.Context, req models.ReorderVenueGalleryRequest) error {
	return s.repo.ReorderVenueItems(ctx, req.VenueID, req.MediaAssetIDs)
}

func (s *venueGalleryService) UpdateGalleryItem(ctx context.Context, req models.UpdateVenueGalleryItemRequest) error {
	// Strategi Fetch-Merge-Update untuk Partial Update yang aman
	// 1. Ambil semua item (biasanya gallery tidak terlalu banyak, jadi ini aman)
	existingItems, err := s.repo.GetByVenueID(ctx, req.VenueID)
	if err != nil {
		return err
	}

	// 2. Cari item yang mau diedit
	var targetItem *entity.VenueGalleryItem
	for i := range existingItems {
		if existingItems[i].MediaAssetID == req.MediaAssetID {
			targetItem = &existingItems[i]
			break
		}
	}

	if targetItem == nil {
		return errors.New("gallery item not found")
	}

	// 3. Merge perubahan (Hanya field yang dikirim user)
	if req.Caption != nil {
		targetItem.Caption = *req.Caption
	}
	if req.IsVisible != nil {
		targetItem.IsVisible = *req.IsVisible
	}
	if req.IsFeatured != nil {
		targetItem.IsFeatured = *req.IsFeatured
	}
	if req.SortOrder != nil {
		targetItem.SortOrder = *req.SortOrder
	}

	// 4. Save
	return s.repo.Update(ctx, targetItem)
}

func (s *venueGalleryService) RemoveGalleryItem(ctx context.Context, targetID, mediaID uuid.UUID) error {
	return s.repo.RemoveVenueItem(ctx, targetID, mediaID)
}

type areaGalleryService struct {
	repo repository.AreaGalleryRepository
}

func NewAreaGalleryService(repo repository.AreaGalleryRepository) AreaGalleryService {
	return &areaGalleryService{repo: repo}
}

func (s *areaGalleryService) AddGalleryItems(ctx context.Context, req models.AddAreaGalleryItemsRequest) error {
	var items []entity.AreaGalleryItem
	for _, item := range req.Items {
		items = append(items, entity.AreaGalleryItem{
			AreaID:       req.AreaID,
			MediaAssetID: item.MediaAssetID,
			Caption:      item.Caption,
			SortOrder:    item.SortOrder,
			IsVisible:    item.IsVisible,
			// Area tidak punya IsFeatured
		})
	}
	return s.repo.AddAreaItems(ctx, items)
}

func (s *areaGalleryService) ReorderGallery(ctx context.Context, req models.ReorderAreaGalleryRequest) error {
	return s.repo.ReorderAreaItems(ctx, req.AreaID, req.MediaAssetIDs)
}

func (s *areaGalleryService) UpdateGalleryItem(ctx context.Context, req models.UpdateAreaGalleryItemRequest) error {
	// Strategi Fetch-Merge-Update
	existingItems, err := s.repo.GetByAreaID(ctx, req.AreaID)
	if err != nil {
		return err
	}

	var targetItem *entity.AreaGalleryItem
	for i := range existingItems {
		if existingItems[i].MediaAssetID == req.MediaAssetID {
			targetItem = &existingItems[i]
			break
		}
	}

	if targetItem == nil {
		return errors.New("gallery item not found")
	}

	if req.Caption != nil {
		targetItem.Caption = *req.Caption
	}
	if req.IsVisible != nil {
		targetItem.IsVisible = *req.IsVisible
	}
	if req.SortOrder != nil {
		targetItem.SortOrder = *req.SortOrder
	}

	return s.repo.UpdateAreaItem(ctx, targetItem)
}

func (s *areaGalleryService) RemoveGalleryItem(ctx context.Context, targetID, mediaID uuid.UUID) error {
	return s.repo.RemoveAreaItem(ctx, targetID, mediaID)
}
