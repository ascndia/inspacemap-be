package service

import (
	"context"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"
)

type roleService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository // [FIX] Tambahkan dependensi ini
}

// [FIX] Update Constructor untuk menerima PermissionRepo
func NewRoleService(rRepo repository.RoleRepository, pRepo repository.PermissionRepository) RoleService {
	return &roleService{
		roleRepo:       rRepo,
		permissionRepo: pRepo,
	}
}

func (s *roleService) GetSystemRoles(ctx context.Context) ([]models.RoleDetail, error) {
	roles, err := s.roleRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var details []models.RoleDetail
	for _, r := range roles {
		perms, _ := s.roleRepo.GetPermissions(ctx, r.ID)

		var permKeys []string
		for _, p := range perms {
			permKeys = append(permKeys, p.Key)
		}

		details = append(details, models.RoleDetail{
			ID:          r.ID,
			Name:        r.Name,
			Permissions: permKeys,
		})
	}

	return details, nil
}

// [FIX] Implementasi Method yang Hilang
func (s *roleService) GetAvailablePermissions(ctx context.Context) ([]models.PermissionNode, error) {
	// Ambil semua permission dari DB
	perms, err := s.permissionRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Grouping: Map[GroupName][]Items
	// Agar di UI bisa ditampilkan per kategori (misal: "Venue Management", "User Access")
	groupedMap := make(map[string][]models.PermissionItem)
	for _, p := range perms {
		groupName := p.Group
		if groupName == "" {
			groupName = "Other"
		}

		groupedMap[groupName] = append(groupedMap[groupName], models.PermissionItem{
			ID:          p.ID,
			Key:         p.Key,
			Description: p.Description,
		})
	}

	// Convert Map to Slice
	var nodes []models.PermissionNode
	for groupName, items := range groupedMap {
		nodes = append(nodes, models.PermissionNode{
			Group: groupName,
			Items: items,
		})
	}

	return nodes, nil
}
