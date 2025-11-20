package service

import (
	"context"
	"errors"
	"fmt"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type mediaService struct {
	mediaRepo repository.MediaAssetRepository
	storage   StorageProvider // Injected Dependency (S3/MinIO)

	// Config (Sebaiknya dari Env)
	bucketName string
	cdnBaseURL string // e.g. "https://cdn.inspacemap.com"
}

func NewMediaService(
	repo repository.MediaAssetRepository,
	storage StorageProvider,
	bucketName string,
	cdnBaseURL string,
) MediaService {
	return &mediaService{
		mediaRepo:  repo,
		storage:    storage,
		bucketName: bucketName,
		cdnBaseURL: cdnBaseURL,
	}
}

// 1. Init Upload (Generate URL)
func (s *mediaService) InitDirectUpload(ctx context.Context, orgID uuid.UUID, req models.PresignedUploadRequest) (*models.PresignedUploadResponse, error) {
	// A. Generate File Key (Path Unik)
	// Format: org_id/category/uuid.ext
	// Contoh: 550e84.../panorama/abc-123.jpg
	fileExt := filepath.Ext(req.FileName)
	if fileExt == "" {
		fileExt = ".jpg"
	} // Default fallback

	newAssetID := uuid.New()
	key := fmt.Sprintf("%s/%s/%s%s", orgID.String(), req.Category, newAssetID.String(), fileExt)

	// B. Generate Presigned URL dari S3 Provider
	// URL ini hanya valid selama 15 menit
	uploadURL, err := s.storage.GetPresignedPutURL(ctx, s.bucketName, key, req.FileType, 15*time.Minute)
	if err != nil {
		return nil, errors.New("failed to generate upload url")
	}

	// C. Simpan Record 'Pending' ke Database
	// Kita simpan dulu agar punya ID. Status validasi dilakukan saat Confirm.
	asset := entity.MediaAsset{
		BaseEntity: entity.BaseEntity{
			ID: newAssetID,
		},
		OrganizationID: orgID,
		Bucket:         s.bucketName,
		Key:            key,
		FileName:       req.FileName,
		MimeType:       req.FileType,
		Type:           req.Category,
		SizeInBytes:    req.FileSize,

		// URL Public untuk akses (Read)
		// Asumsi: Bucket public atau dilayani via CloudFront/Nginx
		PublicURL: fmt.Sprintf("%s/%s", s.cdnBaseURL, key),

		// Visibility Default
		Visibility: entity.VisibilityPublic,
	}

	if err := s.mediaRepo.Create(ctx, &asset); err != nil {
		return nil, err
	}

	return &models.PresignedUploadResponse{
		UploadURL: uploadURL,
		AssetID:   newAssetID,
		Key:       key,
	}, nil
}

// 2. Confirm Upload
func (s *mediaService) ConfirmUpload(ctx context.Context, req models.ConfirmUploadRequest) error {
	// Ambil asset
	asset, err := s.mediaRepo.GetByID(ctx, req.AssetID)
	if err != nil {
		return errors.New("asset not found")
	}

	// Update metadata fisik (Width/Height) yang dikirim frontend
	// Frontend (JS/Flutter) bisa baca dimensi gambar sebelum/sesudah upload.
	// Ini penting untuk rendering map yang akurat.
	asset.Width = req.Width
	asset.Height = req.Height

	// Simpan update
	// Gunakan method update repository (perlu dibuat atau pakai BaseRepo.Update)
	// Di sini kita asumsikan BaseRepo.Update sudah cukup
	return s.mediaRepo.Update(ctx, asset)
}

// 3. Get Asset
func (s *mediaService) GetAsset(ctx context.Context, id uuid.UUID) (*entity.MediaAsset, error) {
	return s.mediaRepo.GetByID(ctx, id)
}

// 4. List Assets
func (s *mediaService) ListAssets(ctx context.Context, query models.MediaAssetQuery) ([]entity.MediaAsset, int64, error) {
	return s.mediaRepo.PagedMediaAssets(ctx, query)
}

// 5. Delete Asset
func (s *mediaService) DeleteAsset(ctx context.Context, id uuid.UUID) error {
	asset, err := s.mediaRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// A. Hapus dari Storage Fisik (S3)
	// Kita lakukan asynchronously atau biarkan error jika gagal (soft delete di DB lebih penting)
	_ = s.storage.DeleteObject(ctx, asset.Bucket, asset.Key)

	// B. Hapus dari DB (Soft Delete)
	return s.mediaRepo.Delete(ctx, id)
}
