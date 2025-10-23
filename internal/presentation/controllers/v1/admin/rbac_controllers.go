package admin

import "github.com/gin-gonic/gin"

// TODO: receiver
type AdminRBACController struct {
}

func NewAdminRBACController() *AdminRBACController {
	return &AdminRBACController{}
}

// TODO: Get All Roles With their Permissions and Users
func (ac *AdminRBACController) GetAllRolesWithUsers(ctx gin.Context) {
	// TODO: Call Service
	// TODO: handle Error
	// TODO: Respond
}

// TODO: Get a Role's Users and Permissions

// TODO: Get a Role's Permissions

// TODO: Get ALl Roles and their Permissions

// TODO: Get a User's Roles

// TODO: Remove a Role from a User

// TODO: add a Role to a User

// TODO: Add a New Role

// TODO: Remove a Role

// TODO: Add a Permission To a Role

// TODO: Remove a Permission From a Role

// TODO: Get a Permissions Roles
