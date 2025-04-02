// /internal/rbac/handlers.go
package rbac

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler RBAC HTTP handler-lərini təmin edir
type Handler struct {
	service Service
}

// NewHandler yeni RBAC handler yaradır
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes API endpoint-lərini qeydiyyatdan keçirir
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	rbacGroup := router.Group("/rbac")
	{
		// Roles API
		rbacGroup.GET("/roles", h.GetAllRoles)
		rbacGroup.GET("/roles/:id", h.GetRoleByID)
		rbacGroup.POST("/roles", h.CreateRole)
		rbacGroup.PUT("/roles/:id", h.UpdateRole)
		rbacGroup.DELETE("/roles/:id", h.DeleteRole)

		// Role permissions API
		rbacGroup.GET("/roles/:id/permissions", h.GetRolePermissions)
		rbacGroup.POST("/roles/:id/permissions", h.UpdateRolePermissions)
		rbacGroup.POST("/roles/:id/permissions/:permissionId", h.AssignPermissionToRole)
		rbacGroup.DELETE("/roles/:id/permissions/:permissionId", h.RemovePermissionFromRole)

		// Permissions API
		rbacGroup.GET("/permissions", h.GetAllPermissions)
		rbacGroup.GET("/permissions/resource/:type", h.GetPermissionsByResourceType)
		rbacGroup.GET("/permissions/:id", h.GetPermissionByID)
		rbacGroup.POST("/permissions", h.CreatePermission)
		rbacGroup.PUT("/permissions/:id", h.UpdatePermission)
		rbacGroup.DELETE("/permissions/:id", h.DeletePermission)

		// User permissions API
		rbacGroup.GET("/users/:id/permissions", h.GetUserPermissions)
		rbacGroup.GET("/users/:id/permissions/direct", h.GetUserDirectPermissions)
		rbacGroup.GET("/users/:id/permissions/role", h.GetUserRolePermissions)
		rbacGroup.POST("/users/permissions", h.AssignPermissionToUser)
		rbacGroup.DELETE("/users/permissions/:id", h.RemovePermissionFromUser)

		// Permission check API
		rbacGroup.POST("/check", h.CheckPermission)
	}
}

// GetAllRoles bütün rolları gətirir
// @Summary Get all roles
// @Description Get all roles in the system
// @Tags rbac
// @Accept json
// @Produce json
// @Success 200 {array} Role
// @Router /api/rbac/roles [get]
func (h *Handler) GetAllRoles(c *gin.Context) {
	roles, err := h.service.GetAllRoles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// GetRoleByID ID ilə rolu gətirir
// @Summary Get role by ID
// @Description Get a role by its ID
// @Tags rbac
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} Role
// @Failure 404 {object} object
// @Router /api/rbac/roles/{id} [get]
func (h *Handler) GetRoleByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	role, err := h.service.GetRoleByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// CreateRole yeni rol yaradır
// @Summary Create a new role
// @Description Create a new role in the system
// @Tags rbac
// @Accept json
// @Produce json
// @Param role body CreateRoleInput true "Role information"
// @Success 201 {object} Role
// @Failure 400 {object} object
// @Router /api/rbac/roles [post]
func (h *Handler) CreateRole(c *gin.Context) {
	var input CreateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.service.CreateRole(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// UpdateRole rolu yeniləyir
// @Summary Update a role
// @Description Update an existing role
// @Tags rbac
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param role body UpdateRoleInput true "Role information"
// @Success 200 {object} Role
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Router /api/rbac/roles/{id} [put]
func (h *Handler) UpdateRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	var input UpdateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.service.UpdateRole(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
			return
		}
		if errors.Is(err, ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// DeleteRole rolu silir
// @Summary Delete a role
// @Description Delete an existing role
// @Tags rbac
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 204 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Router /api/rbac/roles/{id} [delete]
func (h *Handler) DeleteRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	err = h.service.DeleteRole(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllPermissions bütün icazələri gətirir
// @Summary Get all permissions
// @Description Get all permissions in the system
// @Tags rbac
// @Accept json
// @Produce json
// @Success 200 {array} Permission
// @Router /api/rbac/permissions [get]
func (h *Handler) GetAllPermissions(c *gin.Context) {
	permissions, err := h.service.GetAllPermissions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// GetPermissionsByResourceType müəyyən resurs tipinə aid bütün icazələri gətirir
// @Summary Get permissions by resource type
// @Description Get all permissions for a specific resource type
// @Tags rbac
// @Accept json
// @Produce json
// @Param type path string true "Resource Type"
// @Success 200 {array} Permission
// @Router /api/rbac/permissions/resource/{type} [get]
func (h *Handler) GetPermissionsByResourceType(c *gin.Context) {
	resourceType := c.Param("type")

	permissions, err := h.service.GetPermissionsByResourceType(c.Request.Context(), resourceType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// GetPermissionByID ID ilə icazəni gətirir
// @Summary Get permission by ID
// @Description Get a permission by its ID
// @Tags rbac
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Success 200 {object} Permission
// @Failure 404 {object} object
// @Router /api/rbac/permissions/{id} [get]
func (h *Handler) GetPermissionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission ID"})
		return
	}

	permission, err := h.service.GetPermissionByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "permission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permission)
}

// CreatePermission yeni icazə yaradır
// @Summary Create a new permission
// @Description Create a new permission in the system
// @Tags rbac
// @Accept json
// @Produce json
// @Param permission body CreatePermissionInput true "Permission information"
// @Success 201 {object} Permission
// @Failure 400 {object} object
// @Router /api/rbac/permissions [post]
func (h *Handler) CreatePermission(c *gin.Context) {
	var input CreatePermissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission, err := h.service.CreatePermission(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, permission)
}

// UpdatePermission icazəni yeniləyir
// @Summary Update a permission
// @Description Update an existing permission
// @Tags rbac
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Param permission body UpdatePermissionInput true "Permission information"
// @Success 200 {object} Permission
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Router /api/rbac/permissions/{id} [put]
func (h *Handler) UpdatePermission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission ID"})
		return
	}

	var input UpdatePermissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission, err := h.service.UpdatePermission(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "permission not found"})
			return
		}
		if errors.Is(err, ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permission)
}

// DeletePermission icazəni silir
// @Summary Delete a permission
// @Description Delete an existing permission
// @Tags rbac
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Success 204 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Router /api/rbac/permissions/{id} [delete]
func (h *Handler) DeletePermission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission ID"})
		return
	}

	err = h.service.DeletePermission(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "permission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetRolePermissions rola aid bütün icazələri gətirir
// @Summary Get role permissions
// @Description Get all permissions assigned to a role
// @Tags rbac
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {array} Permission
// @Failure 404 {object} object
// @Router /api/rbac/roles/{id}/permissions [get]
func (h *Handler) GetRolePermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	permissions, err := h.service.GetRolePermissions(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// UpdateRolePermissions rolun bütün icazələrini yeniləyir
// @Summary Update role permissions
// @Description Update all permissions assigned to a role
// @Tags rbac
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param permissionIds body []int true "Permission IDs"
// @Success 200 {object} object
// @Failure 404 {object} object
// @Router /api/rbac/roles/{id}/permissions [post]
func (h *Handler) UpdateRolePermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	var permissionIDs []int
	if err := c.ShouldBindJSON(&permissionIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateRolePermissions(c.Request.Context(), id, permissionIDs)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role permissions updated successfully"})
}

// AssignPermissionToRole rola icazə təyin edir
// @Summary Assign permission to role
// @Description Assign a permission to a role
// @Tags rbac
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param permissionId path int true "Permission ID"
// @Success 200 {object} object
// @Failure 404 {object} object
// @Router /api/rbac/roles/{id}/permissions/{permissionId} [post]
func (h *Handler) AssignPermissionToRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	permissionID, err := strconv.Atoi(c.Param("permissionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission ID"})
		return
	}

	err = h.service.AssignPermissionToRole(c.Request.Context(), roleID, permissionID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission assigned to role successfully"})
}

// RemovePermissionFromRole roldan icazəni silir
// @Summary Remove permission from role
// @Description Remove a permission from a role
// @Tags rbac
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param permissionId path int true "Permission ID"
// @Success 200 {object} object
// @Failure 404 {object} object
// @Router /api/rbac/roles/{id}/permissions/{permissionId} [delete]
func (h *Handler) RemovePermissionFromRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	permissionID, err := strconv.Atoi(c.Param("permissionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission ID"})
		return
	}

	err = h.service.RemovePermissionFromRole(c.Request.Context(), roleID, permissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission removed from role successfully"})
}

// GetUserPermissions istifadəçinin bütün icazələrini (rol və fərdi) gətirir
// @Summary Get user permissions
// @Description Get all permissions assigned to a user
// @Tags rbac
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} UserPermissionDetail
// @Router /api/rbac/users/{id}/permissions [get]
func (h *Handler) GetUserPermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	permissions, err := h.service.GetUserPermissions(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// GetUserDirectPermissions istifadəçiyə birbaşa təyin edilmiş icazələri gətirir
// @Summary Get user direct permissions
// @Description Get permissions directly assigned to a user
// @Tags rbac
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} UserPermissionDetail
// @Router /api/rbac/users/{id}/permissions/direct [get]
func (h *Handler) GetUserDirectPermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	permissions, err := h.service.GetUserDirectPermissions(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// GetUserRolePermissions istifadəçinin rolundan gələn icazələri gətirir
// @Summary Get user role permissions
// @Description Get permissions assigned to a user's role
// @Tags rbac
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} PermissionWithRole
// @Router /api/rbac/users/{id}/permissions/role [get]
func (h *Handler) GetUserRolePermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	permissions, err := h.service.GetUserRolePermissions(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// AssignPermissionToUser istifadəçiyə icazə təyin edir
// @Summary Assign permission to user
// @Description Assign a permission to a user
// @Tags rbac
// @Accept json
// @Produce json
// @Param input body AssignUserPermissionInput true "User permission information"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Router /api/rbac/users/permissions [post]
func (h *Handler) AssignPermissionToUser(c *gin.Context) {
	var input AssignUserPermissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AssignPermissionToUser(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission assigned to user successfully"})
}

// RemovePermissionFromUser istifadəçidən icazəni silir
// @Summary Remove permission from user
// @Description Remove a permission from a user
// @Tags rbac
// @Accept json
// @Produce json
// @Param id path int true "User Permission ID"
// @Success 200 {object} object
// @Failure 404 {object} object
// @Router /api/rbac/users/permissions/{id} [delete]
func (h *Handler) RemovePermissionFromUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user permission ID"})
		return
	}

	err = h.service.RemovePermissionFromUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission removed from user successfully"})
}

// CheckPermissionRequest icazə yoxlama sorğusu
type CheckPermissionRequest struct {
	UserID       int    `json:"user_id" binding:"required"`
	ResourceType string `json:"resource_type" binding:"required"`
	Action       string `json:"action" binding:"required"`
	ResourceID   *int   `json:"resource_id"`
}

// CheckPermissionResponse icazə yoxlama cavabı
type CheckPermissionResponse struct {
	HasPermission bool `json:"has_permission"`
}

// CheckPermission istifadəçi icazəsini yoxlayır
// @Summary Check user permission
// @Description Check if a user has a specific permission
// @Tags rbac
// @Accept json
// @Produce json
// @Param request body CheckPermissionRequest true "Permission check request"
// @Success 200 {object} CheckPermissionResponse
// @Failure 400 {object} object
// @Router /api/rbac/check [post]
func (h *Handler) CheckPermission(c *gin.Context) {
	var req CheckPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hasPermission, err := h.service.CheckPermission(
		c.Request.Context(),
		req.UserID,
		req.ResourceType,
		req.Action,
		req.ResourceID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, CheckPermissionResponse{HasPermission: hasPermission})
}
