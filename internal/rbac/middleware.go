// /internal/middleware/rbac.go
package rbac

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Zam83-AZE/logistics_system/internal/rbac"
)

// Auth istifadəçi autentifikasiya məlumatları
type Auth struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name"`
}

// GetUserFromContext kontekstdən istifadəçi məlumatlarını əldə edir
func GetUserFromContext(c *gin.Context) (*Auth, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	auth, ok := user.(*Auth)
	if !ok {
		return nil, false
	}

	return auth, true
}

// RequirePermission icazə tələb edən middleware yaradır
func RequirePermission(service rbac.Service, resourceType, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, exists := GetUserFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		// ResourceID əldə etməyə cəhd edirik (URL-də varsa)
		var resourceID *int
		if idParam := c.Param("id"); idParam != "" {
			if id, err := strconv.Atoi(idParam); err == nil {
				resourceID = &id
			}
		}

		// İcazəni yoxlayırıq
		hasPermission, err := service.CheckPermission(c.Request.Context(), auth.UserID, resourceType, action, resourceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permission"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermissionFunc xüsusi resourceID funksiyası ilə icazə tələb edən middleware yaradır
func RequirePermissionFunc(service rbac.Service, resourceType, action string, resourceIDFunc func(*gin.Context) *int) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, exists := GetUserFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		// ResourceID əldə etmək üçün verilmiş funksiyanı çağırırıq
		resourceID := resourceIDFunc(c)

		// İcazəni yoxlayırıq
		hasPermission, err := service.CheckPermission(c.Request.Context(), auth.UserID, resourceType, action, resourceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permission"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole müəyyən rol tələb edən middleware yaradır
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, exists := GetUserFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		if auth.RoleName != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient role privileges"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin admin rol tələb edən middleware
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("Admin")
}

// CheckPermissionFromBodyJSON request body-dən icazə tələbi əldə edən middleware
func CheckPermissionFromBodyJSON(service rbac.Service, resourceTypeField, actionField, resourceIDField string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, exists := GetUserFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		// Body-dən məlumatları əldə etmək üçün map
		var requestData map[string]interface{}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			c.Abort()
			return
		}

		// ResourceType əldə edirik
		resourceTypeValue, exists := requestData[resourceTypeField]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "resource type not specified"})
			c.Abort()
			return
		}
		resourceType, ok := resourceTypeValue.(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource type"})
			c.Abort()
			return
		}

		// Action əldə edirik
		actionValue, exists := requestData[actionField]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "action not specified"})
			c.Abort()
			return
		}
		action, ok := actionValue.(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
			c.Abort()
			return
		}

		// ResourceID (əgər varsa) əldə edirik
		var resourceID *int
		if resourceIDField != "" {
			if resourceIDValue, exists := requestData[resourceIDField]; exists {
				if floatValue, ok := resourceIDValue.(float64); ok {
					intValue := int(floatValue)
					resourceID = &intValue
				}
			}
		}

		// Original request body-ni keşləyirik ki, sonra yenidən istifadə edə bilək
		c.Set("original_request", requestData)

		// İcazəni yoxlayırıq
		hasPermission, err := service.CheckPermission(c.Request.Context(), auth.UserID, resourceType, action, resourceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permission"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			c.Abort()
			return
		}

		// Original request body-ni yenidən bind edirik
		c.Set("request_data", requestData)
		c.Next()
	}
}

// EnforcePermission bir əməliyyat ərzində icazə yoxlamaq üçün helper funksiya
func EnforcePermission(c *gin.Context, service rbac.Service, resourceType, action string, resourceID *int) error {
	auth, exists := GetUserFromContext(c)
	if !exists {
		return errors.New("authentication required")
	}

	return service.EnforcePermission(c.Request.Context(), auth.UserID, resourceType, action, resourceID)
}
