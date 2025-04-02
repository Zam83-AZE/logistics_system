// /internal/rbac/models.go
package rbac

import (
	"time"
)

// Role istifadəçi rolu modelidir
type Role struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Permission istifadəçi icazəsi modelidir
type Permission struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	ResourceType string    `json:"resource_type" db:"resource_type"`
	Action       string    `json:"action" db:"action"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// RolePermission rol və icazə arasındakı əlaqəni təmsil edir
type RolePermission struct {
	ID           int       `json:"id" db:"id"`
	RoleID       int       `json:"role_id" db:"role_id"`
	PermissionID int       `json:"permission_id" db:"permission_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// UserPermission istifadəçi və icazə arasındakı əlaqəni təmsil edir
type UserPermission struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	PermissionID int       `json:"permission_id" db:"permission_id"`
	ResourceID   *int      `json:"resource_id" db:"resource_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// PermissionWithRole icazə və rol adını birləşdirir
type PermissionWithRole struct {
	Permission
	RoleName string `json:"role_name"`
}

// PermissionWithResource icazə və resurs məlumatlarını birləşdirir
type PermissionWithResource struct {
	Permission
	ResourceID   int    `json:"resource_id"`
	ResourceName string `json:"resource_name"`
}

// UserPermissionDetail istifadəçi icazələrinin tam təfərrüatını təmsil edir
type UserPermissionDetail struct {
	ID             int     `json:"id" db:"id"`
	UserID         int     `json:"user_id" db:"user_id"`
	PermissionID   int     `json:"permission_id" db:"permission_id"`
	PermissionName string  `json:"permission_name" db:"permission_name"`
	ResourceType   string  `json:"resource_type" db:"resource_type"`
	Action         string  `json:"action" db:"action"`
	ResourceID     *int    `json:"resource_id" db:"resource_id"`
	ResourceName   *string `json:"resource_name"`
}

// CreateRoleInput yeni rol yaratmaq üçün input modeli
type CreateRoleInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateRoleInput rol yeniləmək üçün input modeli
type UpdateRoleInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreatePermissionInput yeni icazə yaratmaq üçün input modeli
type CreatePermissionInput struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	ResourceType string `json:"resource_type" binding:"required"`
	Action       string `json:"action" binding:"required"`
}

// UpdatePermissionInput icazə yeniləmək üçün input modeli
type UpdatePermissionInput struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ResourceType string `json:"resource_type"`
	Action       string `json:"action"`
}

// AssignRolePermissionInput rol-icazə təyin etmək üçün input modeli
type AssignRolePermissionInput struct {
	RoleID       int `json:"role_id" binding:"required"`
	PermissionID int `json:"permission_id" binding:"required"`
}

// AssignUserPermissionInput istifadəçi-icazə təyin etmək üçün input modeli
type AssignUserPermissionInput struct {
	UserID       int  `json:"user_id" binding:"required"`
	PermissionID int  `json:"permission_id" binding:"required"`
	ResourceID   *int `json:"resource_id"`
}

// ResourceType resurs tiplərini təmsil edir
const (
	ResourceTypeCustomer   = "customer"
	ResourceTypeContainer  = "container"
	ResourceTypeProduct    = "product"
	ResourceTypeInvoice    = "invoice"
	ResourceTypeUser       = "user"
	ResourceTypeRole       = "role"
	ResourceTypePermission = "permission"
	ResourceTypeReport     = "report"
)

// Action icazə əməliyyatlarını təmsil edir
const (
	ActionView   = "view"
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionAll    = "all" // Bütün əməliyyatlar
)
