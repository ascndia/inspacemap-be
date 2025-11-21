package unit

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
)

func TestGraphService_CreateFloor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGraphRepo := NewMockGraphRepository(ctrl)
	mockGraphRevisionRepo := NewMockGraphRevisionRepository(ctrl)
	mockFloorRepo := NewMockFloorRepository(ctrl)
	mockVenueRepo := NewMockVenueRepository(ctrl)

	graphService := service.NewGraphService(mockGraphRepo, mockGraphRevisionRepo, mockFloorRepo, mockVenueRepo)

	tests := []struct {
		name          string
		venueID       uuid.UUID
		req           models.CreateFloorRequest
		mockSetup     func()
		expectedError bool
		errorContains string
	}{
		{
			name:    "successful floor creation",
			venueID: uuid.New(),
			req: models.CreateFloorRequest{
				Name:           "Ground Floor",
				LevelIndex:     0,
				MapImageID:     &[]uuid.UUID{uuid.New()}[0],
				PixelsPerMeter: 100.0,
				MapWidth:       800,
				MapHeight:      600,
			},
			mockSetup: func() {
				draftRevision := &entity.GraphRevision{
					BaseEntity: entity.BaseEntity{
						ID: uuid.New(),
					},
					Status: "draft",
				}
				mockGraphRevisionRepo.EXPECT().
					GetDraftByVenueID(gomock.Any(), gomock.Any()).
					Return(draftRevision, nil)

				mockFloorRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "draft revision not found, auto-create draft",
			venueID: uuid.New(),
			req: models.CreateFloorRequest{
				Name:           "Ground Floor",
				LevelIndex:     0,
				MapImageID:     &[]uuid.UUID{uuid.New()}[0],
				PixelsPerMeter: 100.0,
				MapWidth:       800,
				MapHeight:      600,
			},
			mockSetup: func() {
				mockGraphRevisionRepo.EXPECT().
					GetDraftByVenueID(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("draft not found"))

				draftRevision := &entity.GraphRevision{
					BaseEntity: entity.BaseEntity{
						ID: uuid.New(),
					},
					Status: "draft",
				}
				mockGraphRevisionRepo.EXPECT().
					CreateDraft(gomock.Any(), gomock.Any()).
					Return(draftRevision, nil)

				mockFloorRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "floor creation fails",
			venueID: uuid.New(),
			req: models.CreateFloorRequest{
				Name:           "Ground Floor",
				LevelIndex:     0,
				MapImageID:     &[]uuid.UUID{uuid.New()}[0],
				PixelsPerMeter: 100.0,
				MapWidth:       800,
				MapHeight:      600,
			},
			mockSetup: func() {
				draftRevision := &entity.GraphRevision{
					BaseEntity: entity.BaseEntity{
						ID: uuid.New(),
					},
					Status: "draft",
				}
				mockGraphRevisionRepo.EXPECT().
					GetDraftByVenueID(gomock.Any(), gomock.Any()).
					Return(draftRevision, nil)

				mockFloorRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			expectedError: true,
			errorContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			_, err := graphService.CreateFloor(context.Background(), tt.venueID, tt.req)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGraphService_CreateNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGraphRepo := NewMockGraphRepository(ctrl)
	mockGraphRevisionRepo := NewMockGraphRevisionRepository(ctrl)
	mockFloorRepo := NewMockFloorRepository(ctrl)
	mockVenueRepo := NewMockVenueRepository(ctrl)

	graphService := service.NewGraphService(mockGraphRepo, mockGraphRevisionRepo, mockFloorRepo, mockVenueRepo)

	tests := []struct {
		name          string
		req           models.CreateNodeRequest
		mockSetup     func()
		expectedError bool
		errorContains string
	}{
		{
			name: "successful node creation",
			req: models.CreateNodeRequest{
				FloorID:         uuid.New(),
				X:               100.0,
				Y:               200.0,
				PanoramaAssetID: uuid.New(),
				Label:           "Entrance",
			},
			mockSetup: func() {
				draftRevision := &entity.GraphRevision{
					BaseEntity: entity.BaseEntity{
						ID: uuid.New(),
					},
					Status: "draft",
				}
				mockGraphRevisionRepo.EXPECT().
					GetDraftByFloorID(gomock.Any(), gomock.Any()).
					Return(draftRevision, nil)

				mockGraphRepo.EXPECT().
					CreateNode(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name: "draft revision not found",
			req: models.CreateNodeRequest{
				FloorID:         uuid.New(),
				X:               100.0,
				Y:               200.0,
				PanoramaAssetID: uuid.New(),
				Label:           "Entrance",
			},
			mockSetup: func() {
				mockGraphRevisionRepo.EXPECT().
					GetDraftByFloorID(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("draft not found"))
			},
			expectedError: true,
			errorContains: "cannot create node: target floor is not in draft mode",
		},
		{
			name: "invalid coordinates",
			req: models.CreateNodeRequest{
				FloorID:         uuid.New(),
				X:               -10.0,
				Y:               200.0,
				PanoramaAssetID: uuid.New(),
				Label:           "Entrance",
			},
			mockSetup: func() {
				draftRevision := &entity.GraphRevision{
					BaseEntity: entity.BaseEntity{
						ID: uuid.New(),
					},
					Status: "draft",
				}
				mockGraphRevisionRepo.EXPECT().
					GetDraftByFloorID(gomock.Any(), gomock.Any()).
					Return(draftRevision, nil)
			},
			expectedError: true,
			errorContains: "coordinates cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			_, err := graphService.CreateNode(context.Background(), tt.req)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGraphService_UpdateNodePosition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGraphRepo := NewMockGraphRepository(ctrl)
	mockGraphRevisionRepo := NewMockGraphRevisionRepository(ctrl)
	mockFloorRepo := NewMockFloorRepository(ctrl)
	mockVenueRepo := NewMockVenueRepository(ctrl)

	graphService := service.NewGraphService(mockGraphRepo, mockGraphRevisionRepo, mockFloorRepo, mockVenueRepo)

	tests := []struct {
		name          string
		nodeID        uuid.UUID
		req           models.UpdateNodePositionRequest
		mockSetup     func()
		expectedError bool
		errorContains string
	}{
		{
			name:   "successful position update",
			nodeID: uuid.New(),
			req: models.UpdateNodePositionRequest{
				X: 150.0,
				Y: 250.0,
			},
			mockSetup: func() {
				mockGraphRepo.EXPECT().
					UpdateNodePosition(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "position update fails",
			nodeID: uuid.New(),
			req: models.UpdateNodePositionRequest{
				X: 150.0,
				Y: 250.0,
			},
			mockSetup: func() {
				mockGraphRepo.EXPECT().
					UpdateNodePosition(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("update failed"))
			},
			expectedError: true,
			errorContains: "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := graphService.UpdateNodePosition(context.Background(), tt.nodeID, tt.req)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGraphService_UpdateNodeCalibration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGraphRepo := NewMockGraphRepository(ctrl)
	mockGraphRevisionRepo := NewMockGraphRevisionRepository(ctrl)
	mockFloorRepo := NewMockFloorRepository(ctrl)
	mockVenueRepo := NewMockVenueRepository(ctrl)

	graphService := service.NewGraphService(mockGraphRepo, mockGraphRevisionRepo, mockFloorRepo, mockVenueRepo)

	tests := []struct {
		name          string
		nodeID        uuid.UUID
		req           models.UpdateNodeCalibrationRequest
		mockSetup     func()
		expectedError bool
		errorContains string
	}{
		{
			name:   "successful calibration update",
			nodeID: uuid.New(),
			req: models.UpdateNodeCalibrationRequest{
				RotationOffset: 45.0,
			},
			mockSetup: func() {
				mockGraphRepo.EXPECT().
					UpdateNodeCalibration(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "calibration update fails",
			nodeID: uuid.New(),
			req: models.UpdateNodeCalibrationRequest{
				RotationOffset: 45.0,
			},
			mockSetup: func() {
				mockGraphRepo.EXPECT().
					UpdateNodeCalibration(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("update failed"))
			},
			expectedError: true,
			errorContains: "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := graphService.UpdateNodeCalibration(context.Background(), tt.nodeID, tt.req)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGraphService_DeleteNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGraphRepo := NewMockGraphRepository(ctrl)
	mockGraphRevisionRepo := NewMockGraphRevisionRepository(ctrl)
	mockFloorRepo := NewMockFloorRepository(ctrl)
	mockVenueRepo := NewMockVenueRepository(ctrl)

	graphService := service.NewGraphService(mockGraphRepo, mockGraphRevisionRepo, mockFloorRepo, mockVenueRepo)

	tests := []struct {
		name          string
		nodeID        uuid.UUID
		mockSetup     func()
		expectedError bool
		errorContains string
	}{
		{
			name:   "successful node deletion",
			nodeID: uuid.New(),
			mockSetup: func() {
				mockGraphRepo.EXPECT().
					DeleteNode(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "node deletion fails",
			nodeID: uuid.New(),
			mockSetup: func() {
				mockGraphRepo.EXPECT().
					DeleteNode(gomock.Any(), gomock.Any()).
					Return(errors.New("delete failed"))
			},
			expectedError: true,
			errorContains: "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := graphService.DeleteNode(context.Background(), tt.nodeID)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGraphService_ConnectNodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGraphRepo := NewMockGraphRepository(ctrl)
	mockGraphRevisionRepo := NewMockGraphRevisionRepository(ctrl)
	mockFloorRepo := NewMockFloorRepository(ctrl)
	mockVenueRepo := NewMockVenueRepository(ctrl)

	graphService := service.NewGraphService(mockGraphRepo, mockGraphRevisionRepo, mockFloorRepo, mockVenueRepo)

	tests := []struct {
		name          string
		req           models.ConnectNodesRequest
		mockSetup     func()
		expectedError bool
		errorContains string
	}{
		{
			name: "successful node connection",
			req: models.ConnectNodesRequest{
				FromNodeID: uuid.New(),
				ToNodeID:   uuid.New(),
			},
			mockSetup: func() {
				mockGraphRepo.EXPECT().
					ConnectNodes(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name: "self-connection not allowed",
			req: func() models.ConnectNodesRequest {
				sameID := uuid.New()
				return models.ConnectNodesRequest{
					FromNodeID: sameID,
					ToNodeID:   sameID,
				}
			}(),
			mockSetup: func() {
				// No mock calls expected since validation happens before repository call
			},
			expectedError: true,
			errorContains: "cannot connect node to itself",
		},
		{
			name: "connection creation fails",
			req: models.ConnectNodesRequest{
				FromNodeID: uuid.New(),
				ToNodeID:   uuid.New(),
			},
			mockSetup: func() {
				mockGraphRepo.EXPECT().
					ConnectNodes(gomock.Any(), gomock.Any()).
					Return(errors.New("connection failed"))
			},
			expectedError: true,
			errorContains: "connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := graphService.ConnectNodes(context.Background(), tt.req)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGraphService_DeleteConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGraphRepo := NewMockGraphRepository(ctrl)
	mockGraphRevisionRepo := NewMockGraphRevisionRepository(ctrl)
	mockFloorRepo := NewMockFloorRepository(ctrl)
	mockVenueRepo := NewMockVenueRepository(ctrl)

	graphService := service.NewGraphService(mockGraphRepo, mockGraphRevisionRepo, mockFloorRepo, mockVenueRepo)

	tests := []struct {
		name          string
		fromID        uuid.UUID
		toID          uuid.UUID
		mockSetup     func()
		expectedError bool
		errorContains string
	}{
		{
			name:   "successful connection deletion",
			fromID: uuid.New(),
			toID:   uuid.New(),
			mockSetup: func() {
				mockGraphRepo.EXPECT().
					DeleteEdge(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "connection deletion fails",
			fromID: uuid.New(),
			toID:   uuid.New(),
			mockSetup: func() {
				mockGraphRepo.EXPECT().
					DeleteEdge(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("delete failed"))
			},
			expectedError: true,
			errorContains: "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := graphService.DeleteConnection(context.Background(), tt.fromID, tt.toID)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGraphService_GetEditorData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGraphRepo := NewMockGraphRepository(ctrl)
	mockGraphRevisionRepo := NewMockGraphRevisionRepository(ctrl)
	mockFloorRepo := NewMockFloorRepository(ctrl)
	mockVenueRepo := NewMockVenueRepository(ctrl)

	graphService := service.NewGraphService(mockGraphRepo, mockGraphRevisionRepo, mockFloorRepo, mockVenueRepo)

	tests := []struct {
		name          string
		venueID       uuid.UUID
		mockSetup     func()
		expectedError bool
		errorContains string
	}{
		{
			name:    "successful editor data retrieval",
			venueID: uuid.New(),
			mockSetup: func() {
				draftRevision := &entity.GraphRevision{
					BaseEntity: entity.BaseEntity{
						ID: uuid.New(),
					},
					Status: "draft",
				}
				mockGraphRevisionRepo.EXPECT().
					GetDraftByVenueID(gomock.Any(), gomock.Any()).
					Return(draftRevision, nil)

				mockVenueRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&entity.Venue{Name: "Test Venue"}, nil)
			},
			expectedError: false,
		},
		{
			name:    "draft revision not found, auto-create",
			venueID: uuid.New(),
			mockSetup: func() {
				mockGraphRevisionRepo.EXPECT().
					GetDraftByVenueID(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("draft not found"))

				draftRevision := &entity.GraphRevision{
					BaseEntity: entity.BaseEntity{
						ID: uuid.New(),
					},
					Status: "draft",
				}
				mockGraphRevisionRepo.EXPECT().
					CreateDraft(gomock.Any(), gomock.Any()).
					Return(draftRevision, nil)

				mockVenueRepo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&entity.Venue{Name: "Test Venue"}, nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			_, err := graphService.GetEditorData(context.Background(), tt.venueID)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGraphService_PublishChanges(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGraphRepo := NewMockGraphRepository(ctrl)
	mockGraphRevisionRepo := NewMockGraphRevisionRepository(ctrl)
	mockFloorRepo := NewMockFloorRepository(ctrl)
	mockVenueRepo := NewMockVenueRepository(ctrl)

	graphService := service.NewGraphService(mockGraphRepo, mockGraphRevisionRepo, mockFloorRepo, mockVenueRepo)

	tests := []struct {
		name          string
		venueID       uuid.UUID
		req           models.PublishDraftRequest
		mockSetup     func()
		expectedError bool
		errorContains string
	}{
		{
			name:    "successful publish",
			venueID: uuid.New(),
			req: models.PublishDraftRequest{
				Note: "Initial draft publish",
			},
			mockSetup: func() {
				mockGraphRevisionRepo.EXPECT().
					PublishDraft(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "publish fails",
			venueID: uuid.New(),
			req: models.PublishDraftRequest{
				Note: "Initial draft publish",
			},
			mockSetup: func() {
				mockGraphRevisionRepo.EXPECT().
					PublishDraft(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("publish failed"))
			},
			expectedError: true,
			errorContains: "publish failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := graphService.PublishChanges(context.Background(), tt.venueID, tt.req)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
