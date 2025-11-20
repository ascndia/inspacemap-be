package service

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"

	"github.com/google/uuid"
)

type organizationService struct {
	orgRepo repository.OrganizationRepository
}

func NewOrganizationService(orgRepo repository.OrganizationRepository) OrganizationService {
	return &organizationService{
		orgRepo: orgRepo,
	}
}

// 1. Get Detail By ID
func (s *organizationService) GetDetailByID(ctx context.Context, id uuid.UUID) (*models.OrganizationDetail, error) {
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("organization not found")
	}
	return s.mapEntityToDetail(org), nil
}

// 2. Get Detail By Slug (Subdomain Resolution)
func (s *organizationService) GetDetailBySlug(ctx context.Context, slug string) (*models.OrganizationDetail, error) {
	// Repo GetByDomain mencari berdasarkan Slug atau Website
	org, err := s.orgRepo.GetByDomain(ctx, slug)
	if err != nil {
		return nil, errors.New("organization not found")
	}
	return s.mapEntityToDetail(org), nil
}

// 3. Update Profile
func (s *organizationService) UpdateProfile(ctx context.Context, id uuid.UUID, req models.UpdateOrganizationRequest) error {
	// Ambil data lama dulu
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("organization not found")
	}

	// Partial Update
	if req.Name != nil {
		org.Name = *req.Name
	}
	if req.Slug != nil {
		// TODO: Cek apakah slug baru sudah dipakai orang lain?
		// existing, _ := s.orgRepo.GetByDomain(ctx, *req.Slug)
		// if existing != nil && existing.ID != id { return error }
		org.Slug = *req.Slug
	}
	if req.LogoURL != nil {
		org.LogoURL = *req.LogoURL
	}
	if req.Website != nil {
		org.Website = *req.Website
	}

	// Update Settings (Merge atau Replace)
	if req.Settings != nil {
		// Strategi: Replace total settings (simple)
		org.Settings = req.Settings

		// Jika ingin Merge:
		// for k, v := range req.Settings { org.Settings[k] = v }
	}

	return s.orgRepo.Update(ctx, org)
}

// 4. List Organizations (Super Admin)
func (s *organizationService) ListOrganizations(ctx context.Context, query models.OrganizationQuery) ([]models.OrganizationDetail, int64, error) {
	// Panggil Repo Paged
	orgs, total, err := s.orgRepo.PagedOrganizations(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	var details []models.OrganizationDetail
	for _, org := range orgs {
		details = append(details, *s.mapEntityToDetail(&org))
	}

	return details, total, nil
}

// 5. Deactivate (Soft Delete / Suspend)
func (s *organizationService) DeactivateOrganization(ctx context.Context, id uuid.UUID) error {
	// Kita bisa pakai Delete (Soft Delete) atau set IsActive = false
	// Menggunakan IsActive = false lebih aman agar data relasi tidak hilang di query biasa
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	org.IsActive = false
	return s.orgRepo.Update(ctx, org)
}

// --- HELPER ---
func (s *organizationService) mapEntityToDetail(org *entity.Organization) *models.OrganizationDetail {
	return &models.OrganizationDetail{
		ID:        org.ID,
		Name:      org.Name,
		Slug:      org.Slug,
		LogoURL:   org.LogoURL,
		Website:   org.Website,
		IsActive:  org.IsActive,
		Settings:  org.Settings,
		CreatedAt: org.CreatedAt,
		// VenueCount & UserCount bisa diisi via query terpisah jika perlu
	}
}
