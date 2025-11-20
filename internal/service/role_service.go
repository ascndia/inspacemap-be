package service

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"

	"github.com/google/uuid"
)

type roleService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

func NewRoleService(rRepo repository.RoleRepository, pRepo repository.PermissionRepository) RoleService {
	return &roleService{
		roleRepo:       rRepo,
		permissionRepo: pRepo,
	}
}

func (s *roleService) CreateCustomRole(ctx context.Context, orgID uuid.UUID, req models.CreateRoleRequest) (*models.IDResponse, error) {
	if existing, _ := s.roleRepo.GetByOrganizationIDAndName(ctx, &orgID, req.Name); existing != nil {
		return nil, errors.New("role with this name already exists in your organization")
	}
	newRole := entity.Role{
		OrganizationID: &orgID, // Set pointer ke Org ID
		Name:           req.Name,
		Description:    req.Description,
		IsSystem:       false,
	}

	if err := s.roleRepo.Create(ctx, &newRole); err != nil {
		return nil, err
	}
	for _, permID := range req.PermissionIDs {
		if err := s.roleRepo.AttachPermission(ctx, newRole.ID, permID); err != nil {
			return nil, errors.New("failed to attach permission settings")
		}
	}

	return &models.IDResponse{ID: newRole.ID}, nil
}

func (s *roleService) DeleteCustomRole(ctx context.Context, roleID uuid.UUID) error {
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return errors.New("role not found")
	}

	// Security Check: Tidak boleh menghapus System Role
	if role.IsSystem {
		return errors.New("cannot delete system role")
	}

	// TODO: Cek apakah ada user yang masih pakai role ini?
	// Jika ada, tolak delete atau pindahkan user ke role default.

	return s.roleRepo.Delete(ctx, roleID)
}

// 3. GetAvailablePermissions
// Mengembalikan daftar izin yang dikelompokkan (Grouping) agar UI mudah merender (e.g. Accordion)
func (s *roleService) GetPermissionsByOrganization(ctx context.Context) ([]models.PermissionNode, error) {
	perms, err := s.permissionRepo.GetByOrganization(ctx)
	if err != nil {
		return nil, err
	}

	// Grouping Logic: Map[GroupName][]Items
	groupedMap := make(map[string][]models.PermissionItem)

	for _, p := range perms {
		groupName := p.Group
		if groupName == "" {
			groupName = "Other"
		}

		groupedMap[groupName] = append(groupedMap[groupName], models.PermissionItem{
			ID:          p.ID, // ID UUID (dari BaseModel)
			Key:         p.Key,
			Description: p.Description,
		})
	}

	// Convert Map to Slice untuk JSON Response yang terurut
	var nodes []models.PermissionNode
	for groupName, items := range groupedMap {
		nodes = append(nodes, models.PermissionNode{
			Group: groupName,
			Items: items,
		})
	}

	return nodes, nil
}

// 4. GetOrgRoles
// Mengambil daftar Role System + Custom Role milik Org
func (s *roleService) GetRolesByOrganization(ctx context.Context, orgID uuid.UUID) ([]models.RoleDetail, error) {
	roles, err := s.roleRepo.GetByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var details []models.RoleDetail
	for _, r := range roles {
		// Mapping Permissions (Entity -> String Keys)
		var permKeys []string
		for _, p := range r.Permissions {
			permKeys = append(permKeys, p.Key)
		}

		details = append(details, models.RoleDetail{
			ID:          r.ID,
			Name:        r.Name,
			IsSystem:    r.IsSystem,
			Permissions: permKeys,
		})
	}

	return details, nil
}
