// /internal/rbac/repository.go
package rbac

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repository rbac üçün verilənlər bazası əməliyyatlarını təmin edir
type Repository interface {
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

	// User-Permission əməliyyatları
	GetUserPermissions(ctx context.Context, userID int) ([]UserPermissionDetail, error)
	GetUserDirectPermissions(ctx context.Context, userID int) ([]UserPermissionDetail, error)
	GetUserRolePermissions(ctx context.Context, userID int) ([]PermissionWithRole, error)
	AssignPermissionToUser(ctx context.Context, input AssignUserPermissionInput) error
	RemovePermissionFromUser(ctx context.Context, userPermissionID int) error

	// İcazə yoxlama əməliyyatları
	CheckUserPermission(ctx context.Context, userID int, resourceType, action string, resourceID *int) (bool, error)
}

// PostgresRepository rbac üçün PostgreSQL repository implementasiyası
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository yeni PostgreSQL RBAC repository yaradır
func NewPostgresRepository(db *sqlx.DB) Repository {
	return &PostgresRepository{db: db}
}

// GetAllRoles bütün rolları gətirir
func (r *PostgresRepository) GetAllRoles(ctx context.Context) ([]Role, error) {
	roles := []Role{}
	query := "SELECT * FROM roles ORDER BY id"

	err := r.db.SelectContext(ctx, &roles, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}

	return roles, nil
}

// GetRoleByID ID ilə rolu gətirir
func (r *PostgresRepository) GetRoleByID(ctx context.Context, id int) (*Role, error) {
	role := Role{}
	query := "SELECT * FROM roles WHERE id = $1"

	err := r.db.GetContext(ctx, &role, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &role, nil
}

// CreateRole yeni rol yaradır
func (r *PostgresRepository) CreateRole(ctx context.Context, input CreateRoleInput) (*Role, error) {
	var role Role
	query := `
		INSERT INTO roles (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`

	err := r.db.QueryRowxContext(ctx, query, input.Name, input.Description).StructScan(&role)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &role, nil
}

// UpdateRole rolu yeniləyir
func (r *PostgresRepository) UpdateRole(ctx context.Context, id int, input UpdateRoleInput) (*Role, error) {
	var role Role
	query := `
		UPDATE roles
		SET name = COALESCE($2, name),
			description = COALESCE($3, description),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, name, description, created_at, updated_at
	`

	err := r.db.QueryRowxContext(ctx, query, id, input.Name, input.Description).StructScan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return &role, nil
}

// DeleteRole rolu silir
func (r *PostgresRepository) DeleteRole(ctx context.Context, id int) error {
	query := "DELETE FROM roles WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role with ID %d not found", id)
	}

	return nil
}

// GetAllPermissions bütün icazələri gətirir
func (r *PostgresRepository) GetAllPermissions(ctx context.Context) ([]Permission, error) {
	permissions := []Permission{}
	query := "SELECT * FROM permissions ORDER BY id"

	err := r.db.SelectContext(ctx, &permissions, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	return permissions, nil
}

// GetPermissionsByResourceType müəyyən resurs tipinə aid bütün icazələri gətirir
func (r *PostgresRepository) GetPermissionsByResourceType(ctx context.Context, resourceType string) ([]Permission, error) {
	permissions := []Permission{}
	query := "SELECT * FROM permissions WHERE resource_type = $1 ORDER BY id"

	err := r.db.SelectContext(ctx, &permissions, query, resourceType)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions by resource type: %w", err)
	}

	return permissions, nil
}

// GetPermissionByID ID ilə icazəni gətirir
func (r *PostgresRepository) GetPermissionByID(ctx context.Context, id int) (*Permission, error) {
	permission := Permission{}
	query := "SELECT * FROM permissions WHERE id = $1"

	err := r.db.GetContext(ctx, &permission, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return &permission, nil
}

// CreatePermission yeni icazə yaradır
func (r *PostgresRepository) CreatePermission(ctx context.Context, input CreatePermissionInput) (*Permission, error) {
	var permission Permission
	query := `
		INSERT INTO permissions (name, description, resource_type, action)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, resource_type, action, created_at, updated_at
	`

	err := r.db.QueryRowxContext(
		ctx, query, input.Name, input.Description, input.ResourceType, input.Action,
	).StructScan(&permission)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return &permission, nil
}

// UpdatePermission icazəni yeniləyir
func (r *PostgresRepository) UpdatePermission(ctx context.Context, id int, input UpdatePermissionInput) (*Permission, error) {
	var permission Permission
	query := `
		UPDATE permissions
		SET name = COALESCE($2, name),
			description = COALESCE($3, description),
			resource_type = COALESCE($4, resource_type),
			action = COALESCE($5, action),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, name, description, resource_type, action, created_at, updated_at
	`

	err := r.db.QueryRowxContext(
		ctx, query, id, input.Name, input.Description, input.ResourceType, input.Action,
	).StructScan(&permission)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update permission: %w", err)
	}

	return &permission, nil
}

// DeletePermission icazəni silir
func (r *PostgresRepository) DeletePermission(ctx context.Context, id int) error {
	query := "DELETE FROM permissions WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("permission with ID %d not found", id)
	}

	return nil
}

// GetRolePermissions rola aid bütün icazələri gətirir
func (r *PostgresRepository) GetRolePermissions(ctx context.Context, roleID int) ([]Permission, error) {
	permissions := []Permission{}
	query := `
		SELECT p.*
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.id
	`

	err := r.db.SelectContext(ctx, &permissions, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	return permissions, nil
}

// AssignPermissionToRole rola icazə təyin edir
func (r *PostgresRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionID int) error {
	query := `
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES ($1, $2)
		ON CONFLICT (role_id, permission_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to assign permission to role: %w", err)
	}

	return nil
}

// RemovePermissionFromRole roldan icazəni silir
func (r *PostgresRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID int) error {
	query := "DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2"

	result, err := r.db.ExecContext(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to remove permission from role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role permission combination not found")
	}

	return nil
}

// GetUserPermissions istifadəçinin bütün icazələrini (rol və fərdi) gətirir
func (r *PostgresRepository) GetUserPermissions(ctx context.Context, userID int) ([]UserPermissionDetail, error) {
	userPermissions := []UserPermissionDetail{}
	query := `
		SELECT DISTINCT ON (p.id, up.resource_id) 
			COALESCE(up.id, 0) as id, 
			$1 as user_id,
			p.id as permission_id, 
			p.name as permission_name,
			p.resource_type, 
			p.action, 
			up.resource_id
		FROM permissions p
		
		-- Rol icazələri
		LEFT JOIN role_permissions rp ON p.id = rp.permission_id
		LEFT JOIN roles r ON rp.role_id = r.id
		LEFT JOIN users u ON u.role_id = r.id
		
		-- Birbaşa istifadəçi icazələri
		LEFT JOIN user_permissions up ON (p.id = up.permission_id AND up.user_id = $1)
		
		WHERE (u.id = $1 OR up.user_id = $1)
		ORDER BY p.id, up.resource_id, id DESC
	`

	err := r.db.SelectContext(ctx, &userPermissions, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// Resurs adları əlavə edilə bilər (resource_id əsasında)
	// Bu funksionallıq buraya əlavə edilə bilər

	return userPermissions, nil
}

// GetUserDirectPermissions yalnız istifadəçiyə birbaşa təyin edilmiş icazələri gətirir
func (r *PostgresRepository) GetUserDirectPermissions(ctx context.Context, userID int) ([]UserPermissionDetail, error) {
	userPermissions := []UserPermissionDetail{}
	query := `
		SELECT 
			up.id, 
			up.user_id,
			p.id as permission_id, 
			p.name as permission_name,
			p.resource_type, 
			p.action, 
			up.resource_id
		FROM user_permissions up
		JOIN permissions p ON up.permission_id = p.id
		WHERE up.user_id = $1
		ORDER BY up.id
	`

	err := r.db.SelectContext(ctx, &userPermissions, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user direct permissions: %w", err)
	}

	return userPermissions, nil
}

// GetUserRolePermissions roldan gələn icazələri (rol əsasında) əldə edir
func (r *PostgresRepository) GetUserRolePermissions(ctx context.Context, userID int) ([]PermissionWithRole, error) {
	rolePermissions := []PermissionWithRole{}
	query := `
		SELECT 
			p.*,
			r.name as role_name
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN roles r ON rp.role_id = r.id
		JOIN users u ON u.role_id = r.id
		WHERE u.id = $1
		ORDER BY p.id
	`

	err := r.db.SelectContext(ctx, &rolePermissions, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role permissions: %w", err)
	}

	return rolePermissions, nil
}

// AssignPermissionToUser istifadəçiyə icazə təyin edir
func (r *PostgresRepository) AssignPermissionToUser(ctx context.Context, input AssignUserPermissionInput) error {
	query := `
		INSERT INTO user_permissions (user_id, permission_id, resource_id)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(ctx, query, input.UserID, input.PermissionID, input.ResourceID)
	if err != nil {
		return fmt.Errorf("failed to assign permission to user: %w", err)
	}

	return nil
}

// RemovePermissionFromUser istifadəçidən icazəni silir
func (r *PostgresRepository) RemovePermissionFromUser(ctx context.Context, userPermissionID int) error {
	query := "DELETE FROM user_permissions WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, userPermissionID)
	if err != nil {
		return fmt.Errorf("failed to remove permission from user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user permission with ID %d not found", userPermissionID)
	}

	return nil
}

// CheckUserPermission istifadəçinin icazəsini yoxlayır
func (r *PostgresRepository) CheckUserPermission(ctx context.Context, userID int, resourceType, action string, resourceID *int) (bool, error) {
	// Birinci addım: İstifadəçinin birbaşa özünə təyin edilmiş xüsusi (resurs ID ilə) icazələri yoxlanılır
	if resourceID != nil {
		var count int
		query := `
			SELECT COUNT(1) FROM user_permissions up
			JOIN permissions p ON up.permission_id = p.id
			WHERE up.user_id = $1
			AND p.resource_type = $2
			AND (p.action = $3 OR p.action = 'all')
			AND up.resource_id = $4
		`
		err := r.db.GetContext(ctx, &count, query, userID, resourceType, action, resourceID)
		if err != nil {
			return false, fmt.Errorf("failed to check specific resource permission: %w", err)
		}
		if count > 0 {
			return true, nil
		}
	}

	// İkinci addım: İstifadəçinin ümumi (resurs ID olmadan) icazələri yoxlanılır
	var count int
	query := `
		SELECT COUNT(1) FROM user_permissions up
		JOIN permissions p ON up.permission_id = p.id
		WHERE up.user_id = $1
		AND p.resource_type = $2
		AND (p.action = $3 OR p.action = 'all')
		AND up.resource_id IS NULL
	`
	err := r.db.GetContext(ctx, &count, query, userID, resourceType, action)
	if err != nil {
		return false, fmt.Errorf("failed to check general user permission: %w", err)
	}
	if count > 0 {
		return true, nil
	}

	// Üçüncü addım: İstifadəçinin rolundan gələn icazələr yoxlanılır
	query = `
		SELECT COUNT(1) FROM roles r
		JOIN users u ON u.role_id = r.id
		JOIN role_permissions rp ON r.id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		WHERE u.id = $1
		AND p.resource_type = $2
		AND (p.action = $3 OR p.action = 'all')
	`
	err = r.db.GetContext(ctx, &count, query, userID, resourceType, action)
	if err != nil {
		return false, fmt.Errorf("failed to check role permission: %w", err)
	}

	return count > 0, nil
}
