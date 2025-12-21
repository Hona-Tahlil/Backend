package httpv1

import (
	"hona/backend/wire"

	"github.com/gin-gonic/gin"
)

func SetUpAdminRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	rbacGroup := v1.Group("/rbac")
	{
		// ===== ROLES =====
		rbacGroup.GET("/roles", app.Controllers.AdminControllers.AdminRBACController.GetAllRoles)
		rbacGroup.GET("/roles/:id", app.Controllers.AdminControllers.AdminRBACController.GetRoleByID)
		rbacGroup.GET("/roles/type/:type", app.Controllers.AdminControllers.AdminRBACController.GetRoleByType)
		rbacGroup.POST("/roles", app.Controllers.AdminControllers.AdminRBACController.AddRole)
		rbacGroup.DELETE("/roles/:id", app.Controllers.AdminControllers.AdminRBACController.RemoveRoleByID)
		rbacGroup.DELETE("/roles/type/:type", app.Controllers.AdminControllers.AdminRBACController.RemoveRoleByType)

		// ===== ROLES WITH USERS =====
		rbacGroup.GET("/roles-with-users", app.Controllers.AdminControllers.AdminRBACController.ListRolesWithUsers)
		rbacGroup.GET("/roles-with-users/:id", app.Controllers.AdminControllers.AdminRBACController.GetRoleWithUsersByID)
		rbacGroup.GET("/roles-with-users/type/:type", app.Controllers.AdminControllers.AdminRBACController.GetRoleWithUsersByType)

		// ===== USERS & ROLES =====
		rbacGroup.GET("/users/:id/roles", app.Controllers.AdminControllers.AdminRBACController.GetUserRolesByID)
		rbacGroup.GET("/users/email/roles/", app.Controllers.AdminControllers.AdminRBACController.GetUserRolesByEmail)
		rbacGroup.POST("/users/roles/", app.Controllers.AdminControllers.AdminRBACController.AddRoleToUserByID)
		rbacGroup.POST("/users/email/roles", app.Controllers.AdminControllers.AdminRBACController.AddRoleToUserByEmail)
		rbacGroup.DELETE("/users/:userID/roles/:roleID", app.Controllers.AdminControllers.AdminRBACController.RemoveRoleFromUserByID)
		rbacGroup.DELETE("/users/email/:userEmail/roles/:roleID", app.Controllers.AdminControllers.AdminRBACController.RemoveRoleFromUserByEmail)

		// ===== PERMISSIONS =====
		rbacGroup.GET("/permissions/:permissionID/roles", app.Controllers.AdminControllers.AdminRBACController.GetPermissionRoles)
		rbacGroup.POST("/permissions/roles", app.Controllers.AdminControllers.AdminRBACController.AddPermissionToRole)
		rbacGroup.DELETE("/permissions/:roleID/:permissionID", app.Controllers.AdminControllers.AdminRBACController.RemovePermissionFromRole)

	}

	petSitters := v1.Group("/petsitters")
	{
		petSitters.POST("/search", app.Controllers.AdminControllers.AdminPetSitterController.SearchPetSitters)
	}
}
