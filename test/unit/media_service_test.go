package unit

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMediaService_InitDirectUpload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockMediaAssetRepository(ctrl)
	mockStorage := NewMockStorageProvider(ctrl)

	mediaSvc := service.NewMediaService(mockRepo, mockStorage, "test-bucket", "https://cdn.example.com")

	ctx := context.Background()
	orgID := uuid.New()
	req := models.PresignedUploadRequest{
		FileName: "test.jpg",
		FileType: "image/jpeg",
		Category: "panorama",
		FileSize: 1024,
	}

	expectedKey := "" // We'll capture this
	expectedURL := "https://presigned-url.example.com"

	mockStorage.EXPECT().
		GetPresignedPutURL(ctx, "test-bucket", gomock.Any(), "image/jpeg", 15*time.Minute).
		DoAndReturn(func(ctx context.Context, bucket, key, contentType string, expiry time.Duration) (string, error) {
			expectedKey = key
			return expectedURL, nil
		}).
		Times(1)

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, asset *entity.MediaAsset) error {
			assert.Equal(t, orgID, asset.OrganizationID)
			assert.Equal(t, "test-bucket", asset.Bucket)
			assert.Equal(t, expectedKey, asset.Key)
			assert.Equal(t, "test.jpg", asset.FileName)
			assert.Equal(t, "image/jpeg", asset.MimeType)
			assert.Equal(t, "panorama", asset.Type)
			assert.Equal(t, int64(1024), asset.SizeInBytes)
			assert.Contains(t, asset.PublicURL, "https://cdn.example.com/")
			assert.Equal(t, entity.VisibilityPublic, asset.Visibility)
			return nil
		}).
		Times(1)

	resp, err := mediaSvc.InitDirectUpload(ctx, orgID, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedURL, resp.UploadURL)
	assert.NotEqual(t, uuid.Nil, resp.AssetID)
	assert.Equal(t, expectedKey, resp.Key)
	assert.Contains(t, expectedKey, orgID.String())
	assert.Contains(t, expectedKey, "panorama")
	assert.Contains(t, expectedKey, ".jpg")
}

func TestMediaService_InitDirectUpload_StorageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockMediaAssetRepository(ctrl)
	mockStorage := NewMockStorageProvider(ctrl)

	mediaSvc := service.NewMediaService(mockRepo, mockStorage, "test-bucket", "https://cdn.example.com")

	ctx := context.Background()
	orgID := uuid.New()
	req := models.PresignedUploadRequest{
		FileName: "test.jpg",
		FileType: "image/jpeg",
		Category: "panorama",
		FileSize: 1024,
	}

	mockStorage.EXPECT().
		GetPresignedPutURL(ctx, "test-bucket", gomock.Any(), "image/jpeg", 15*time.Minute).
		Return("", errors.New("storage error")).
		Times(1)

	// Repo should not be called
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Times(0)

	resp, err := mediaSvc.InitDirectUpload(ctx, orgID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to generate upload url")
}

func TestMediaService_InitDirectUpload_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockMediaAssetRepository(ctrl)
	mockStorage := NewMockStorageProvider(ctrl)

	mediaSvc := service.NewMediaService(mockRepo, mockStorage, "test-bucket", "https://cdn.example.com")

	ctx := context.Background()
	orgID := uuid.New()
	req := models.PresignedUploadRequest{
		FileName: "test.jpg",
		FileType: "image/jpeg",
		Category: "panorama",
		FileSize: 1024,
	}

	mockStorage.EXPECT().
		GetPresignedPutURL(ctx, "test-bucket", gomock.Any(), "image/jpeg", 15*time.Minute).
		Return("https://presigned-url.example.com", nil).
		Times(1)

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(errors.New("repo error")).
		Times(1)

	resp, err := mediaSvc.InitDirectUpload(ctx, orgID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestMediaService_ConfirmUpload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockMediaAssetRepository(ctrl)
	mockStorage := NewMockStorageProvider(ctrl)

	mediaSvc := service.NewMediaService(mockRepo, mockStorage, "test-bucket", "https://cdn.example.com")

	ctx := context.Background()
	assetID := uuid.New()
	req := models.ConfirmUploadRequest{
		AssetID: assetID,
		Width:   1920,
		Height:  1080,
	}

	existingAsset := &entity.MediaAsset{
		BaseEntity: entity.BaseEntity{ID: assetID},
	}

	mockRepo.EXPECT().
		GetByID(ctx, assetID).
		Return(existingAsset, nil).
		Times(1)

	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, asset *entity.MediaAsset) error {
			assert.Equal(t, assetID, asset.ID)
			assert.Equal(t, 1920, asset.Width)
			assert.Equal(t, 1080, asset.Height)
			return nil
		}).
		Times(1)

	err := mediaSvc.ConfirmUpload(ctx, req)

	assert.NoError(t, err)
}

func TestMediaService_ConfirmUpload_AssetNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockMediaAssetRepository(ctrl)
	mockStorage := NewMockStorageProvider(ctrl)

	mediaSvc := service.NewMediaService(mockRepo, mockStorage, "test-bucket", "https://cdn.example.com")

	ctx := context.Background()
	assetID := uuid.New()
	req := models.ConfirmUploadRequest{
		AssetID: assetID,
		Width:   1920,
		Height:  1080,
	}

	mockRepo.EXPECT().
		GetByID(ctx, assetID).
		Return(nil, errors.New("not found")).
		Times(1)

	// Update should not be called
	mockRepo.EXPECT().Update(ctx, gomock.Any()).Times(0)

	err := mediaSvc.ConfirmUpload(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "asset not found")
}

func TestMediaService_GetAsset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockMediaAssetRepository(ctrl)
	mockStorage := NewMockStorageProvider(ctrl)

	mediaSvc := service.NewMediaService(mockRepo, mockStorage, "test-bucket", "https://cdn.example.com")

	ctx := context.Background()
	assetID := uuid.New()
	expectedAsset := &entity.MediaAsset{
		BaseEntity: entity.BaseEntity{ID: assetID},
		FileName:   "test.jpg",
	}

	mockRepo.EXPECT().
		GetByID(ctx, assetID).
		Return(expectedAsset, nil).
		Times(1)

	asset, err := mediaSvc.GetAsset(ctx, assetID)

	assert.NoError(t, err)
	assert.Equal(t, expectedAsset, asset)
}

func TestMediaService_ListAssets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockMediaAssetRepository(ctrl)
	mockStorage := NewMockStorageProvider(ctrl)

	mediaSvc := service.NewMediaService(mockRepo, mockStorage, "test-bucket", "https://cdn.example.com")

	ctx := context.Background()
	orgID := uuid.New()
	query := models.MediaAssetQuery{
		MediaAssetFilter: models.MediaAssetFilter{
			OrganizationID: &orgID,
		},
		Limit:  &[]int{10}[0], // pointer to 10
		Offset: &[]int{0}[0],  // pointer to 0
	}

	expectedAssets := []entity.MediaAsset{
		{FileName: "test1.jpg"},
		{FileName: "test2.jpg"},
	}
	expectedTotal := int64(2)

	mockRepo.EXPECT().
		PagedMediaAssets(ctx, query).
		Return(expectedAssets, expectedTotal, nil).
		Times(1)

	assets, total, err := mediaSvc.ListAssets(ctx, query)

	assert.NoError(t, err)
	assert.Equal(t, expectedAssets, assets)
	assert.Equal(t, expectedTotal, total)
}

func TestMediaService_DeleteAsset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockMediaAssetRepository(ctrl)
	mockStorage := NewMockStorageProvider(ctrl)

	mediaSvc := service.NewMediaService(mockRepo, mockStorage, "test-bucket", "https://cdn.example.com")

	ctx := context.Background()
	assetID := uuid.New()
	existingAsset := &entity.MediaAsset{
		BaseEntity: entity.BaseEntity{ID: assetID},
		Bucket:     "test-bucket",
		Key:        "test-key",
	}

	mockRepo.EXPECT().
		GetByID(ctx, assetID).
		Return(existingAsset, nil).
		Times(1)

	mockStorage.EXPECT().
		DeleteObject(ctx, "test-bucket", "test-key").
		Return(nil).
		Times(1)

	mockRepo.EXPECT().
		Delete(ctx, assetID).
		Return(nil).
		Times(1)

	err := mediaSvc.DeleteAsset(ctx, assetID)

	assert.NoError(t, err)
}

func TestMediaService_DeleteAsset_AssetNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockMediaAssetRepository(ctrl)
	mockStorage := NewMockStorageProvider(ctrl)

	mediaSvc := service.NewMediaService(mockRepo, mockStorage, "test-bucket", "https://cdn.example.com")

	ctx := context.Background()
	assetID := uuid.New()

	mockRepo.EXPECT().
		GetByID(ctx, assetID).
		Return(nil, errors.New("not found")).
		Times(1)

	// Storage and delete should not be called
	mockStorage.EXPECT().DeleteObject(ctx, gomock.Any(), gomock.Any()).Times(0)
	mockRepo.EXPECT().Delete(ctx, assetID).Times(0)

	err := mediaSvc.DeleteAsset(ctx, assetID)

	assert.Error(t, err)
}
