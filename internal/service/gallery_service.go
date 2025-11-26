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
		galleryItem := entity.VenueGalleryItem{
			VenueID:      req.VenueID,
			MediaAssetID: item.MediaAssetID,
			Caption:      item.Caption,
			SortOrder:    item.SortOrder,
		}
		if item.IsVisible != nil {
			galleryItem.IsVisible = *item.IsVisible
		} else {
			galleryItem.IsVisible = true // Default to visible
		}
		if item.IsFeatured != nil {
			galleryItem.IsFeatured = *item.IsFeatured
		} else {
			galleryItem.IsFeatured = false // Default to not featured
		}
		items = append(items, galleryItem)
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
	repo         repository.AreaGalleryRepository
	areaRepo     repository.AreaRepository
	revisionRepo repository.GraphRevisionRepository
}

func NewAreaGalleryService(repo repository.AreaGalleryRepository, areaRepo repository.AreaRepository, revisionRepo repository.GraphRevisionRepository) AreaGalleryService {
	return &areaGalleryService{
		repo:         repo,
		areaRepo:     areaRepo,
		revisionRepo: revisionRepo,
	}
}

func (s *areaGalleryService) AddGalleryItems(ctx context.Context, req models.AddAreaGalleryItemsRequest) error {
	// Validate that area belongs to a draft revision
	if err := s.validateAreaInDraft(ctx, req.AreaID); err != nil {
		return err
	}

	var items []entity.AreaGalleryItem
	for _, item := range req.Items {
		galleryItem := entity.AreaGalleryItem{
			AreaID:       req.AreaID,
			MediaAssetID: item.MediaAssetID,
			Caption:      item.Caption,
			SortOrder:    item.SortOrder,
		}
		if item.IsVisible != nil {
			galleryItem.IsVisible = *item.IsVisible
		} else {
			galleryItem.IsVisible = true // Default to visible
		}
		// Area doesn't have IsFeatured
		items = append(items, galleryItem)
	}
	return s.repo.AddAreaItems(ctx, items)
}

func (s *areaGalleryService) ReorderGallery(ctx context.Context, req models.ReorderAreaGalleryRequest) error {
	// Validate that area belongs to a draft revision
	if err := s.validateAreaInDraft(ctx, req.AreaID); err != nil {
		return err
	}

	return s.repo.ReorderAreaItems(ctx, req.AreaID, req.MediaAssetIDs)
}

func (s *areaGalleryService) UpdateGalleryItem(ctx context.Context, req models.UpdateAreaGalleryItemRequest) error {
	// Validate that area belongs to a draft revision
	if err := s.validateAreaInDraft(ctx, req.AreaID); err != nil {
		return err
	}

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
	// Validate that area belongs to a draft revision
	if err := s.validateAreaInDraft(ctx, targetID); err != nil {
		return err
	}

	return s.repo.RemoveAreaItem(ctx, targetID, mediaID)
}

func (s *areaGalleryService) GetGalleryItems(ctx context.Context, areaID uuid.UUID) ([]models.AreaGalleryDetail, error) {
	items, err := s.repo.GetByAreaID(ctx, areaID)
	if err != nil {
		return nil, err
	}

	var details []models.AreaGalleryDetail
	for _, item := range items {
		detail := models.AreaGalleryDetail{
			MediaID:      item.MediaAssetID,
			URL:          item.MediaAsset.PublicURL,
			ThumbnailURL: item.MediaAsset.ThumbnailURL,
			Caption:      item.Caption,
			SortOrder:    item.SortOrder,
			IsFeatured:   false, // Area galleries don't have featured items
		}
		details = append(details, detail)
	}

	return details, nil
}

// Helper method to validate area belongs to draft revision
func (s *areaGalleryService) validateAreaInDraft(ctx context.Context, areaID uuid.UUID) error {
	area, err := s.areaRepo.GetByID(ctx, areaID)
	if err != nil {
		return errors.New("area not found")
	}

	// Check if the area's graph revision is draft
	_, err = s.revisionRepo.GetDraftByFloorID(ctx, area.FloorID)
	if err != nil {
		return errors.New("area gallery can only be modified in draft revisions")
	}

	return nil
}
