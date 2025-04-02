// /internal/rbac/service.go
package rbac

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	// ErrUnauthorized icazə olmadıqda qaytarılan xəta
	ErrUnauthorized = errors.New("unauthorized: permission denied")

	// ErrNotFound entity tapılmadıqda qaytarılan xəta
	ErrNotFound = errors.New("entity not found")

	// ErrInvalidInput yanlış input parametrləri üçün xəta
	ErrInvalidInput = errors.New("invalid input parameters")
)

// Service RBAC servis interfeysi
type Service interface {
	// Role əməliyyatları
	GetAllRoles(ctx context.Context) ([]Role, error)
	GetRoleByID(ctx context.Context, id int) (*Role, error)
	CreateRole(ctx context.Context, input CreateRoleInput) (*Role, error)
	UpdateRole(ctx context.Context, id int, input UpdateRoleInput) (*Role, error)
	DeleteRole(ctx context.Context, id int) error

	// Permission əməliyyatları
	GetAllPermissions(ctx context.Context) ([]Permission, error)
	GetPermissionsByResourceType(ctx context.Context, resourceType string) ([]Permission, error)
	GetPermissionByID(ctx context.Context, id int) (*Permission, error)
	CreatePermission(ctx context.Context, input CreatePermissionInput) (*Permission, error)
	UpdatePermission(ctx context.Context, id int, input UpdatePermissionInput) (*Permission, error)
	DeletePermission(ctx context.Context, id int) error

	// Role-Permission əməliyyatları
	GetRolePermissions(ctx context.Context, roleID int) ([]Permission, error)
	AssignPermissionToRole(ctx context.Context, roleID, permissionID int) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID int) error
	UpdateRolePermissions(ctx context.Context, roleID int, permissionIDs []int) error

	// User-Permission əməliyyatları
	GetUserPermissions(ctx context.Context, userID int) ([]UserPermissionDetail, error)
	GetUserDirectPermissions(ctx context.Context, userID int) ([]UserPermissionDetail, error)
	GetUserRolePermissions(ctx context.Context, userID int) ([]PermissionWithRole, error)
	AssignPermissionToUser(ctx context.Context, input AssignUserPermissionInput) error
	RemovePermissionFromUser(ctx context.Context, userPermissionID int) error

	// İcazə yoxlama əməliyyatları
	CheckPermission(ctx context.Context, userID int, resourceType, action string, resourceID *int) (bool, error)
	EnforcePermission(ctx context.Context, userID int, resourceType, action string, resourceID *int) error

	// Keş əməliyyatları
	ClearCache() error
}

// serviceImpl RBAC servis implementasiyası
type serviceImpl struct {
	repo Repository
	// İcazə keşi - performans üçün
	permissionsCache   map[int]map[string]bool // userID -> "resourceType:action[:resourceID]" -> bool
	permissionsCacheMu sync.RWMutex
}

// NewService yeni RBAC servisi yaradır
func NewService(repo Repository) Service {
	return &serviceImpl{
		repo:             repo,
		permissionsCache: make(map[int]map[string]bool),
	}
}

// GetAllRoles bütün rolları gətirir
func (s *serviceImpl) GetAllRoles(ctx context.Context) ([]Role, error) {
	return s.repo.GetAllRoles(ctx)
}

// GetRoleByID ID ilə rolu gətirir
func (s *serviceImpl) GetRoleByID(ctx context.Context, id int) (*Role, error) {
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrNotFound
	}
	return role, nil
}

// CreateRole yeni rol yaradır
func (s *serviceImpl) CreateRole(ctx context.Context, input CreateRoleInput) (*Role, error) {
	// Input validasiyası burada ola bilər
	if input.Name == "" {
		return nil, fmt.Errorf("%w: role name is required", ErrInvalidInput)
	}

	role, err := s.repo.CreateRole(ctx, input)
	if err != nil {
		return nil, err
	}

	// Yeni rol yaradıldıqda keşi təmizləyirik
	s.ClearCache()

	return role, nil
}

// UpdateRole rolu yeniləyir
func (s *serviceImpl) UpdateRole(ctx context.Context, id int, input UpdateRoleInput) (*Role, error) {
	// Əvvəlcə rolun mövcudluğunu yoxlayırıq
	existingRole, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingRole == nil {
		return nil, ErrNotFound
	}

	role, err := s.repo.UpdateRole(ctx, id, input)
	if err != nil {
		return nil, err
	}

	// Rol yeniləndikdə keşi təmizləyirik
	s.ClearCache()

	return role, nil
}

// DeleteRole rolu silir
func (s *serviceImpl) DeleteRole(ctx context.Context, id int) error {
	// Əvvəlcə rolun mövcudluğunu yoxlayırıq
	existingRole, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return err
	}
	if existingRole == nil {
		return ErrNotFound
	}

	err = s.repo.DeleteRole(ctx, id)
	if err != nil {
		return err
	}

	// Rol silindikdə keşi təmizləyirik
	s.ClearCache()

	return nil
}

// GetAllPermissions bütün icazələri gətirir
func (s *serviceImpl) GetAllPermissions(ctx context.Context) ([]Permission, error) {
	return s.repo.GetAllPermissions(ctx)
}

// GetPermissionsByResourceType resurs tipinə görə icazələri gətirir
func (s *serviceImpl) GetPermissionsByResourceType(ctx context.Context, resourceType string) ([]Permission, error) {
	return s.repo.GetPermissionsByResourceType(ctx, resourceType)
}

// GetPermissionByID ID ilə icazəni gətirir
func (s *serviceImpl) GetPermissionByID(ctx context.Context, id int) (*Permission, error) {
	permission, err := s.repo.GetPermissionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, ErrNotFound
	}
	return permission, nil
}

// CreatePermission yeni icazə yaradır
func (s *serviceImpl) CreatePermission(ctx context.Context, input CreatePermissionInput) (*Permission, error) {
	// Input validasiyası
	if input.Name == "" || input.ResourceType == "" || input.Action == "" {
		return nil, fmt.Errorf("%w: name, resource_type and action are required", ErrInvalidInput)
	}

	permission, err := s.repo.CreatePermission(ctx, input)
	if err != nil {
		return nil, err
	}

	// İcazə yaradıldıqda keşi təmizləyirik
	s.ClearCache()

	return permission, nil
}

// UpdatePermission icazəni yeniləyir
func (s *serviceImpl) UpdatePermission(ctx context.Context, id int, input UpdatePermissionInput) (*Permission, error) {
	// Əvvəlcə icazənin mövcudluğunu yoxlayırıq
	existingPermission, err := s.repo.GetPermissionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingPermission == nil {
		return nil, ErrNotFound
	}

	permission, err := s.repo.UpdatePermission(ctx, id, input)
	if err != nil {
		return nil, err
	}

	// İcazə yeniləndikdə keşi təmizləyirik
	s.ClearCache()

	return permission, nil
}

// DeletePermission icazəni silir
func (s *serviceImpl) DeletePermission(ctx context.Context, id int) error {
	// Əvvəlcə icazənin mövcudluğunu yoxlayırıq
	existingPermission, err := s.repo.GetPermissionByID(ctx, id)
	if err != nil {
		return err
	}
	if existingPermission == nil {
		return ErrNotFound
	}

	err = s.repo.DeletePermission(ctx, id)
	if err != nil {
		return err
	}

	// İcazə silindikdə keşi təmizləyirik
	s.ClearCache()

	return nil
}

// GetRolePermissions rola aid bütün icazələri gətirir
func (s *serviceImpl) GetRolePermissions(ctx context.Context, roleID int) ([]Permission, error) {
	// Əvvəlcə rolun mövcudluğunu yoxlayırıq
	existingRole, err := s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if existingRole == nil {
		return nil, ErrNotFound
	}

	return s.repo.GetRolePermissions(ctx, roleID)
}

// AssignPermissionToRole rola icazə təyin edir
func (s *serviceImpl) AssignPermissionToRole(ctx context.Context, roleID, permissionID int) error {
	// Rolun və icazənin mövcudluğunu yoxlayırıq
	existingRole, err := s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		return err
	}
	if existingRole == nil {
		return fmt.Errorf("%w: role with ID %d not found", ErrNotFound, roleID)
	}

	existingPermission, err := s.repo.GetPermissionByID(ctx, permissionID)
	if err != nil {
		return err
	}
	if existingPermission == nil {
		return fmt.Errorf("%w: permission with ID %d not found", ErrNotFound, permissionID)
	}

	err = s.repo.AssignPermissionToRole(ctx, roleID, permissionID)
	if err != nil {
		return err
	}

	// İcazə təyin edildikdə keşi təmizləyirik
	s.ClearCache()

	return nil
}

// RemovePermissionFromRole roldan icazəni silir
func (s *serviceImpl) RemovePermissionFromRole(ctx context.Context, roleID, permissionID int) error {
	err := s.repo.RemovePermissionFromRole(ctx, roleID, permissionID)
	if err != nil {
		return err
	}

	// İcazə silindikdə keşi təmizləyirik
	s.ClearCache()

	return nil
}

// UpdateRolePermissions rolun bütün icazələrini yeniləyir
func (s *serviceImpl) UpdateRolePermissions(ctx context.Context, roleID int, permissionIDs []int) error {
	// Rolun mövcudluğunu yoxlayırıq
	existingRole, err := s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		return err
	}
	if existingRole == nil {
		return fmt.Errorf("%w: role with ID %d not found", ErrNotFound, roleID)
	}

	// Əvvəlcə rolun mövcud icazələrini əldə edirik
	currentPermissions, err := s.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return err
	}

	// Cari icazələri map-ə çeviririk (ID -> bool)
	currentPermissionMap := make(map[int]bool)
	for _, p := range currentPermissions {
		currentPermissionMap[p.ID] = true
	}

	// Əlavə ediləcək icazələri və silinəcək icazələri müəyyənləşdiririk
	permissionMap := make(map[int]bool)
	for _, id := range permissionIDs {
		permissionMap[id] = true

		// Əgər icazə artıq mövcud deyilsə, əlavə edirik
		if !currentPermissionMap[id] {
			if err := s.repo.AssignPermissionToRole(ctx, roleID, id); err != nil {
				return err
			}
		}
	}

	// Artıq mövcud olmayan icazələri silirik
	for _, p := range currentPermissions {
		if !permissionMap[p.ID] {
			if err := s.repo.RemovePermissionFromRole(ctx, roleID, p.ID); err != nil {
				return err
			}
		}
	}

	// İcazələr yeniləndikdə keşi təmizləyirik
	s.ClearCache()

	return nil
}

// GetUserPermissions istifadəçinin bütün icazələrini (rol və fərdi) gətirir
func (s *serviceImpl) GetUserPermissions(ctx context.Context, userID int) ([]UserPermissionDetail, error) {
	return s.repo.GetUserPermissions(ctx, userID)
}

// GetUserDirectPermissions istifadəçiyə birbaşa təyin edilmiş icazələri gətirir
func (s *serviceImpl) GetUserDirectPermissions(ctx context.Context, userID int) ([]UserPermissionDetail, error) {
	return s.repo.GetUserDirectPermissions(ctx, userID)
}

// GetUserRolePermissions istifadəçinin rolundan gələn icazələri gətirir
func (s *serviceImpl) GetUserRolePermissions(ctx context.Context, userID int) ([]PermissionWithRole, error) {
	return s.repo.GetUserRolePermissions(ctx, userID)
}

// AssignPermissionToUser istifadəçiyə icazə təyin edir
func (s *serviceImpl) AssignPermissionToUser(ctx context.Context, input AssignUserPermissionInput) error {
	// İcazənin mövcudluğunu yoxlayırıq
	existingPermission, err := s.repo.GetPermissionByID(ctx, input.PermissionID)
	if err != nil {
		return err
	}
	if existingPermission == nil {
		return fmt.Errorf("%w: permission with ID %d not found", ErrNotFound, input.PermissionID)
	}

	err = s.repo.AssignPermissionToUser(ctx, input)
	if err != nil {
		return err
	}

	// İcazə təyin edildikdə keşi təmizləyirik
	s.clearUserCache(input.UserID)

	return nil
}

// RemovePermissionFromUser istifadəçidən icazəni silir
func (s *serviceImpl) RemovePermissionFromUser(ctx context.Context, userPermissionID int) error {
	// İstifadəçi icazələrini əldə etmək üçün əvvəlcə user_id-ni əldə edirik
	// Sonra həmin user_id üçün keşi təmizləyirik

	err := s.repo.RemovePermissionFromUser(ctx, userPermissionID)
	if err != nil {
		return err
	}

	// Əgər hansı istifadəçinin icazəsi silindiyini bilmiriksə, bütün keşi təmizləyirik
	s.ClearCache()

	return nil
}

// generateCacheKey keş üçün açar yaradır
func generateCacheKey(resourceType, action string, resourceID *int) string {
	if resourceID != nil {
		return fmt.Sprintf("%s:%s:%d", resourceType, action, *resourceID)
	}
	return fmt.Sprintf("%s:%s", resourceType, action)
}

// CheckPermission istifadəçinin icazəsini yoxlayır
func (s *serviceImpl) CheckPermission(ctx context.Context, userID int, resourceType, action string, resourceID *int) (bool, error) {
	// Keşdə yoxlayırıq
	s.permissionsCacheMu.RLock()
	userCache, exists := s.permissionsCache[userID]
	if exists {
		cacheKey := generateCacheKey(resourceType, action, resourceID)
		if result, ok := userCache[cacheKey]; ok {
			s.permissionsCacheMu.RUnlock()
			return result, nil
		}
	}
	s.permissionsCacheMu.RUnlock()

	// Keşdə tapılmadıqda verilənlər bazasında yoxlayırıq
	result, err := s.repo.CheckUserPermission(ctx, userID, resourceType, action, resourceID)
	if err != nil {
		return false, err
	}

	// Nəticəni keşə əlavə edirik
	s.permissionsCacheMu.Lock()
	defer s.permissionsCacheMu.Unlock()

	if _, exists := s.permissionsCache[userID]; !exists {
		s.permissionsCache[userID] = make(map[string]bool)
	}

	cacheKey := generateCacheKey(resourceType, action, resourceID)
	s.permissionsCache[userID][cacheKey] = result

	return result, nil
}

// EnforcePermission istifadəçinin icazəsini yoxlayır və icazə yoxdursa xəta qaytarır
func (s *serviceImpl) EnforcePermission(ctx context.Context, userID int, resourceType, action string, resourceID *int) error {
	hasPermission, err := s.CheckPermission(ctx, userID, resourceType, action, resourceID)
	if err != nil {
		return err
	}

	if !hasPermission {
		return fmt.Errorf("%w: user %d does not have %s permission for %s",
			ErrUnauthorized, userID, action, resourceType)
	}

	return nil
}

// clearUserCache istifadəçi üçün keşi təmizləyir
func (s *serviceImpl) clearUserCache(userID int) {
	s.permissionsCacheMu.Lock()
	defer s.permissionsCacheMu.Unlock()

	delete(s.permissionsCache, userID)
}

// ClearCache bütün icazə keşini təmizləyir
func (s *serviceImpl) ClearCache() error {
	s.permissionsCacheMu.Lock()
	defer s.permissionsCacheMu.Unlock()

	s.permissionsCache = make(map[int]map[string]bool)
	return nil
}
