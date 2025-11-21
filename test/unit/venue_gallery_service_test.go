package unit

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type VenueGalleryServiceTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	galleryRepo *MockVenueGalleryRepository
	service     service.VenueGalleryService
}

func (suite *VenueGalleryServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.galleryRepo = NewMockVenueGalleryRepository(suite.ctrl)
	suite.service = service.NewVenueGalleryService(suite.galleryRepo)
}

func (suite *VenueGalleryServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func TestVenueGalleryServiceTestSuite(t *testing.T) {
	suite.Run(t, new(VenueGalleryServiceTestSuite))
}

func (suite *VenueGalleryServiceTestSuite) TestAddGalleryItems_Success() {
	ctx := context.Background()
	venueID := uuid.New()
	req := models.AddGalleryVenueItemsRequest{
		VenueID: venueID,
		Items: []models.VenueGalleryItemPayload{
			{
				MediaAssetID: uuid.New(),
				Caption:      "Test image 1",
				SortOrder:    1,
				IsVisible:    true,
				IsFeatured:   false,
			},
			{
				MediaAssetID: uuid.New(),
				Caption:      "Test image 2",
				SortOrder:    2,
				IsVisible:    true,
				IsFeatured:   true,
			},
		},
	}

	suite.galleryRepo.EXPECT().AddVenueItems(ctx, gomock.Any()).Return(nil)

	err := suite.service.AddGalleryItems(ctx, req)

	assert.NoError(suite.T(), err)
}

func (suite *VenueGalleryServiceTestSuite) TestAddGalleryItems_DatabaseError() {
	ctx := context.Background()
	venueID := uuid.New()
	req := models.AddGalleryVenueItemsRequest{
		VenueID: venueID,
		Items: []models.VenueGalleryItemPayload{
			{
				MediaAssetID: uuid.New(),
				Caption:      "Test image",
				SortOrder:    1,
				IsVisible:    true,
				IsFeatured:   false,
			},
		},
	}

	suite.galleryRepo.EXPECT().AddVenueItems(ctx, gomock.Any()).Return(errors.New("database error"))

	err := suite.service.AddGalleryItems(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
}

func (suite *VenueGalleryServiceTestSuite) TestReorderGallery_Success() {
	ctx := context.Background()
	req := models.ReorderVenueGalleryRequest{
		VenueID:       uuid.New(),
		MediaAssetIDs: []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
	}

	suite.galleryRepo.EXPECT().ReorderVenueItems(ctx, req.VenueID, req.MediaAssetIDs).Return(nil)

	err := suite.service.ReorderGallery(ctx, req)

	assert.NoError(suite.T(), err)
}

func (suite *VenueGalleryServiceTestSuite) TestReorderGallery_DatabaseError() {
	ctx := context.Background()
	req := models.ReorderVenueGalleryRequest{
		VenueID:       uuid.New(),
		MediaAssetIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	suite.galleryRepo.EXPECT().ReorderVenueItems(ctx, req.VenueID, req.MediaAssetIDs).Return(errors.New("database error"))

	err := suite.service.ReorderGallery(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
}

func (suite *VenueGalleryServiceTestSuite) TestUpdateGalleryItem_Success() {
	ctx := context.Background()
	venueID := uuid.New()
	mediaID := uuid.New()
	req := models.UpdateVenueGalleryItemRequest{
		VenueID:      venueID,
		MediaAssetID: mediaID,
		Caption:      stringPtr("Updated caption"),
		SortOrder:    intPtr(5),
		IsVisible:    boolPtr(false),
		IsFeatured:   boolPtr(true),
	}

	existingItems := []entity.VenueGalleryItem{
		{
			BaseEntity:   entity.BaseEntity{ID: uuid.New()},
			VenueID:      venueID,
			MediaAssetID: mediaID,
			Caption:      "Old caption",
			SortOrder:    1,
			IsVisible:    true,
			IsFeatured:   false,
		},
	}

	suite.galleryRepo.EXPECT().GetByVenueID(ctx, venueID).Return(existingItems, nil)
	suite.galleryRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	err := suite.service.UpdateGalleryItem(ctx, req)

	assert.NoError(suite.T(), err)
}

func (suite *VenueGalleryServiceTestSuite) TestUpdateGalleryItem_DatabaseError() {
	ctx := context.Background()
	venueID := uuid.New()
	mediaID := uuid.New()
	req := models.UpdateVenueGalleryItemRequest{
		VenueID:      venueID,
		MediaAssetID: mediaID,
		Caption:      stringPtr("Updated caption"),
	}

	suite.galleryRepo.EXPECT().GetByVenueID(ctx, venueID).Return([]entity.VenueGalleryItem{}, nil)

	err := suite.service.UpdateGalleryItem(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "gallery item not found", err.Error())
}

func (suite *VenueGalleryServiceTestSuite) TestRemoveGalleryItem_Success() {
	ctx := context.Background()
	venueID := uuid.New()
	mediaID := uuid.New()

	suite.galleryRepo.EXPECT().RemoveVenueItem(ctx, venueID, mediaID).Return(nil)

	err := suite.service.RemoveGalleryItem(ctx, venueID, mediaID)

	assert.NoError(suite.T(), err)
}

func (suite *VenueGalleryServiceTestSuite) TestRemoveGalleryItem_DatabaseError() {
	ctx := context.Background()
	venueID := uuid.New()
	mediaID := uuid.New()

	suite.galleryRepo.EXPECT().RemoveVenueItem(ctx, venueID, mediaID).Return(errors.New("database error"))

	err := suite.service.RemoveGalleryItem(ctx, venueID, mediaID)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
