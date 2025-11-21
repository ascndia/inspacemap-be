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

type VenueServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	venueRepo *MockVenueRepository
	service   service.VenueService
}

func (suite *VenueServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.venueRepo = NewMockVenueRepository(suite.ctrl)
	suite.service = service.NewVenueService(suite.venueRepo)
}

func (suite *VenueServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func TestVenueServiceTestSuite(t *testing.T) {
	suite.Run(t, new(VenueServiceTestSuite))
}

func (suite *VenueServiceTestSuite) TestCreateVenue_Success() {
	ctx := context.Background()
	req := models.CreateVenueRequest{
		Name:         "Test Venue",
		Slug:         "test-venue",
		Description:  "A test venue",
		Address:      "123 Test St",
		City:         "Test City",
		Province:     "Test Province",
		PostalCode:   "12345",
		Latitude:     -6.2088,
		Longitude:    106.8456,
		CoverImageID: nil, // nil UUID
		Visibility:   "public",
		Gallery: []models.VenueGalleryItemRequest{
			{
				MediaAssetID: uuid.New(),
				Caption:      "Test image",
				SortOrder:    1,
				IsVisible:    true,
				IsFeatured:   false,
			},
		},
	}

	venueID := uuid.New()
	suite.venueRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, venue *entity.Venue) error {
		venue.ID = venueID
		return nil
	})

	result, err := suite.service.CreateVenue(ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), venueID, result.ID)
}

func (suite *VenueServiceTestSuite) TestCreateVenue_DatabaseError() {
	ctx := context.Background()
	req := models.CreateVenueRequest{
		Name:        "Test Venue",
		Slug:        "test-venue",
		Description: "A test venue",
		Address:     "123 Test St",
		City:        "Test City",
		Province:    "Test Province",
		PostalCode:  "12345",
		Latitude:    -6.2088,
		Longitude:   106.8456,
		Visibility:  "private",
	}

	suite.venueRepo.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("database error"))

	result, err := suite.service.CreateVenue(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "database error", err.Error())
}

func (suite *VenueServiceTestSuite) TestUpdateVenue_Success() {
	ctx := context.Background()
	venueID := uuid.New()
	req := models.UpdateVenueRequest{
		Name:        stringPtr("Updated Venue"),
		Description: stringPtr("Updated description"),
		City:        stringPtr("Updated City"),
		Visibility:  stringPtr("public"),
	}

	existingVenue := &entity.Venue{
		BaseEntity: entity.BaseEntity{ID: venueID},
		Name:       "Old Name",
		City:       "Old City",
		Visibility: entity.VisibilityPrivate,
	}

	suite.venueRepo.EXPECT().GetByID(ctx, venueID).Return(existingVenue, nil)
	suite.venueRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	err := suite.service.UpdateVenue(ctx, venueID, req)

	assert.NoError(suite.T(), err)
}

func (suite *VenueServiceTestSuite) TestUpdateVenue_NotFound() {
	ctx := context.Background()
	venueID := uuid.New()
	req := models.UpdateVenueRequest{
		Name: stringPtr("Updated Venue"),
	}

	suite.venueRepo.EXPECT().GetByID(ctx, venueID).Return(nil, errors.New("not found"))

	err := suite.service.UpdateVenue(ctx, venueID, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "venue not found", err.Error())
}

func (suite *VenueServiceTestSuite) TestDeleteVenue_Success() {
	ctx := context.Background()
	venueID := uuid.New()

	suite.venueRepo.EXPECT().Delete(ctx, venueID).Return(nil)

	err := suite.service.DeleteVenue(ctx, venueID)

	assert.NoError(suite.T(), err)
}

func (suite *VenueServiceTestSuite) TestGetVenueDetail_Success() {
	ctx := context.Background()
	venueID := uuid.New()

	venue := &entity.Venue{
		BaseEntity:       entity.BaseEntity{ID: venueID},
		Name:             "Test Venue",
		Slug:             "test-venue",
		City:             "Test City",
		Visibility:       entity.VisibilityPublic,
		Gallery:          []entity.VenueGalleryItem{},
		PointsOfInterest: []entity.Area{},
	}

	suite.venueRepo.EXPECT().GetByID(ctx, venueID).Return(venue, nil)

	result, err := suite.service.GetVenueDetail(ctx, venueID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "Test Venue", result.Name)
	assert.Equal(suite.T(), "public", result.Visibility)
}

func (suite *VenueServiceTestSuite) TestGetVenueBySlug_Success() {
	ctx := context.Background()
	slug := "test-venue"

	venue := &entity.Venue{
		BaseEntity:       entity.BaseEntity{ID: uuid.New()},
		Name:             "Test Venue",
		Slug:             slug,
		City:             "Test City",
		Visibility:       entity.VisibilityPublic,
		Gallery:          []entity.VenueGalleryItem{},
		PointsOfInterest: []entity.Area{},
	}

	suite.venueRepo.EXPECT().GetBySlug(ctx, slug).Return(venue, nil)

	result, err := suite.service.GetVenueBySlug(ctx, slug)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "Test Venue", result.Name)
	assert.Equal(suite.T(), slug, result.Slug)
}

func (suite *VenueServiceTestSuite) TestListVenues_Success() {
	ctx := context.Background()
	limit := 10
	offset := 0
	query := models.VenueQuery{
		Limit:  &limit,
		Offset: &offset,
	}

	venues := []entity.Venue{
		{
			BaseEntity: entity.BaseEntity{ID: uuid.New()},
			Name:       "Venue 1",
			Slug:       "venue-1",
			City:       "City 1",
			Visibility: entity.VisibilityPublic,
		},
		{
			BaseEntity: entity.BaseEntity{ID: uuid.New()},
			Name:       "Venue 2",
			Slug:       "venue-2",
			City:       "City 2",
			Visibility: entity.VisibilityPrivate,
		},
	}

	suite.venueRepo.EXPECT().PagedVenues(ctx, query).Return(venues, int64(2), nil)

	result, total, err := suite.service.ListVenues(ctx, query)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), total)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), "Venue 1", result[0].Name)
	assert.Equal(suite.T(), "Venue 2", result[1].Name)
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
