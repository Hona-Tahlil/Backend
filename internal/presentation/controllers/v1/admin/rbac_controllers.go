package admin

import (
	"hona/backend/internal/application/dto/rbac"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type AdminRBACController struct {
	rbacService usecase.RBACService
}

func NewAdminRBACController() *AdminRBACController {
	return &AdminRBACController{}
}

func (ac *AdminRBACController) GetAllRolesWithUsers(ctx *gin.Context) {
	res, err := ac.rbacService.GetAllRolesWithUsers()
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}

func (ac *AdminRBACController) GetRoleWithUsersByID(ctx *gin.Context) {
	type GetRoleWithUsersByIDParams struct {
		ID uint `uri:"id"`
	}
	params := controllers.Receive[GetRoleWithUsersByIDParams](ctx)

	GetRoleWithUsersByIDInfo := rbac.GetRoleByIDRequest{
		ID: params.ID,
	}
	res, err := ac.rbacService.GetRoleWithUsersByID(GetRoleWithUsersByIDInfo)
	if err != nil {
		panic(err)
	}

	controllers.Respond(ctx, 200, controllers.Message{}, *res)
}

func (ac *AdminRBACController) GetRoleWithUsersByType(ctx *gin.Context) {
	type GetRoleWithUsersByTypeParams struct {
		Type string `uri:"type"`
	}
	params := controllers.Receive[GetRoleWithUsersByTypeParams](ctx)

	GetRoleWithUsersByTypeInfo := rbac.GetRoleByTypeRequest{
		Type: params.Type,
	}
	res, err := ac.rbacService.GetRoleWithUsersByType(GetRoleWithUsersByTypeInfo)
	if err != nil {
		panic(err)
	}

	controllers.Respond(ctx, 200, controllers.Message{}, *res)
}

func (ac *AdminRBACController) GetRoleByID(ctx *gin.Context) {
	type GetRoleByIDParams struct {
		ID uint `uri:"id"`
	}
	params := controllers.Receive[GetRoleByIDParams](ctx)

	GetRoleByIDInfo := rbac.GetRoleByIDRequest{
		ID: params.ID,
	}
	res, err := ac.rbacService.GetRoleByID(GetRoleByIDInfo)
	if err != nil {
		panic(err)
	}
	controllers.Respond(ctx, 200, controllers.Message{}, *res)
}

func (ac *AdminRBACController) GetRoleByType(ctx *gin.Context) {
	type GetRoleByTypeParams struct {
		Type string `uri:"type"`
	}
	params := controllers.Receive[GetRoleByTypeParams](ctx)

	GetRoleByTypeInfo := rbac.GetRoleByTypeRequest{
		Type: params.Type,
	}
	res, err := ac.rbacService.GetRoleByType(GetRoleByTypeInfo)
	if err != nil {
		panic(err)
	}

	controllers.Respond(ctx, 200, controllers.Message{}, *res)
}

func (ac *AdminRBACController) GetAllRoles(ctx *gin.Context) {
	res, err := ac.rbacService.GetAllRoles()
	if err != nil {
		panic(err)
	}

	controllers.Respond(ctx, 200, controllers.Message{}, res)
}

func (ac *AdminRBACController) GetUserRolesByID(ctx *gin.Context) {
	type GetUserRolesByIDParams struct {
		ID uint `uri:"id"`
	}
	params := controllers.Receive[GetUserRolesByIDParams](ctx)
	GetUserRolesByIDInfo := rbac.GetUserRolesByIDRequest{
		ID: params.ID,
	}
	res, err := ac.rbacService.GetUserRolesByID(GetUserRolesByIDInfo)
	if err != nil {
		panic(err)
	}

	controllers.Respond(ctx, 200, controllers.Message{}, res)
}

func (ac *AdminRBACController) GetUserRolesByEmail(ctx *gin.Context) {
	type GetUserRolesByEmailParams struct {
		Email string `form:"email"`
	}
	params := controllers.Receive[GetUserRolesByEmailParams](ctx)
	GetUserRolesByEmailInfo := rbac.GetUserRolesByEmailRequest{
		Email: params.Email,
	}
	res, err := ac.rbacService.GetUserRolesByEmail(GetUserRolesByEmailInfo)
	if err != nil {
		panic(err)
	}

	controllers.Respond(ctx, 200, controllers.Message{}, res)
}

func (ac *AdminRBACController) RemoveRoleFromUserByID(ctx *gin.Context) {
	type RemoveRoleFromUserByIDParams struct {
		UserID uint `uri:"userID"`
		RoleID uint `uri:"roleID"`
	}
	params := controllers.Receive[RemoveRoleFromUserByIDParams](ctx)
	RemoveRoleFromUserByIDInfo := rbac.RemoveRoleFromUserByIDRequest{
		UserID: params.UserID,
		RoleID: params.RoleID,
	}
	err := ac.rbacService.RemoveRoleFromUserByID(RemoveRoleFromUserByIDInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{
		Text: "successMessage.generic",
	}
	controllers.Respond(ctx, 200, msg, nil)
}

func (ac *AdminRBACController) RemoveRoleFromUserByEmail(ctx *gin.Context) {
	type RemoveRoleFromUserByEmailParams struct {
		UserEmail string `uri:"userEmail"`
		RoleID    uint   `uri:"roleID"`
	}
	params := controllers.Receive[RemoveRoleFromUserByEmailParams](ctx)
	RemoveRoleFromUserByEmailInfo := rbac.RemoveRoleFromUserByEmailRequest{
		UserEmail: params.UserEmail,
		RoleID:    params.RoleID,
	}
	err := ac.rbacService.RemoveRoleFromUserByEmail(RemoveRoleFromUserByEmailInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{
		Text: "successMessage.generic",
	}
	controllers.Respond(ctx, 200, msg, nil)
}

func (ac *AdminRBACController) AddRoleToUserByID(ctx *gin.Context) {
	type AddRoleToUserByIDParams struct {
		UserID uint `uri:"userID"`
		RoleID uint `uri:"roleID"`
	}
	params := controllers.Receive[AddRoleToUserByIDParams](ctx)
	AddRoleToUserByIDInfo := rbac.AddRoleToUserByIDRequest{
		UserID: params.UserID,
		RoleID: params.RoleID,
	}
	err := ac.rbacService.AddRoleToUserByID(AddRoleToUserByIDInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{
		Text: "successMessage.generic",
	}
	controllers.Respond(ctx, 200, msg, nil)
}

func (ac *AdminRBACController) AddRoleToUserByEmail(ctx *gin.Context) {
	type AddRoleToUserByEmailParams struct {
		UserEmail string `uri:"userEmail"`
		RoleID    uint   `uri:"roleID"`
	}
	params := controllers.Receive[AddRoleToUserByEmailParams](ctx)
	AddRoleToUserByEmailInfo := rbac.AddRoleToUserByEmailRequest{
		UserEmail: params.UserEmail,
		RoleID:    params.RoleID,
	}
	err := ac.rbacService.AddRoleToUserByEmail(AddRoleToUserByEmailInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{
		Text: "successMessage.generic",
	}
	controllers.Respond(ctx, 200, msg, nil)
}

func (ac *AdminRBACController) AddRole(ctx *gin.Context) {
	type AddRoleParams struct {
		Type        string  `json:"type"`
		Description *string `json:"description"`
	}
	params := controllers.Receive[AddRoleParams](ctx)
	AddRoleInfo := rbac.AddRoleRequest{
		Type:        params.Type,
		Description: params.Description,
	}
	err := ac.rbacService.AddRole(AddRoleInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{
		Text: "successMessage.generic",
	}
	controllers.Respond(ctx, 200, msg, nil)
}

func (ac *AdminRBACController) RemoveRoleByID(ctx *gin.Context) {
	type RemoveRoleByIDParams struct {
		ID uint `uri:"id"`
	}
	params := controllers.Receive[RemoveRoleByIDParams](ctx)
	RemoveRoleByIDInfo := rbac.RemoveRoleByIDRequest{
		ID: params.ID,
	}
	err := ac.rbacService.RemoveRoleByID(RemoveRoleByIDInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{
		Text: "successMessage.generic",
	}
	controllers.Respond(ctx, 200, msg, nil)
}

func (ac *AdminRBACController) RemoveRoleByType(ctx *gin.Context) {
	type RemoveRoleByTypeParams struct {
		Type string `uri:"type"`
	}
	params := controllers.Receive[RemoveRoleByTypeParams](ctx)
	RemoveRoleByTypeInfo := rbac.RemoveRoleByTypeRequest{
		Type: params.Type,
	}
	err := ac.rbacService.RemoveRoleByType(RemoveRoleByTypeInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{
		Text: "successMessage.generic",
	}
	controllers.Respond(ctx, 200, msg, nil)
}

func (ac *AdminRBACController) AddPermissionToRole(ctx *gin.Context) {
	type AddPermissionToRoleParams struct {
		RoleID       uint `json:"roleID"`
		PermissionID uint `json:"permissionID"`
	}
	params := controllers.Receive[AddPermissionToRoleParams](ctx)
	AddPermissionToRoleInfo := rbac.AddPermissionToRoleRequest{
		RoleID:       params.RoleID,
		PermissionID: params.PermissionID,
	}
	err := ac.rbacService.AddPermissionToRole(AddPermissionToRoleInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{
		Text: "successMessage.generic",
	}
	controllers.Respond(ctx, 200, msg, nil)
}

func (ac *AdminRBACController) RemovePermissionFromRole(ctx *gin.Context) {
	type RemovePermissionFromRoleParams struct {
		RoleID       uint `json:"roleID"`
		PermissionID uint `json:"permissionID"`
	}
	params := controllers.Receive[RemovePermissionFromRoleParams](ctx)
	RemovePermissionFromRoleInfo := rbac.RemovePermissionFromRoleRequest{
		RoleID:       params.RoleID,
		PermissionID: params.PermissionID,
	}
	err := ac.rbacService.RemovePermissionFromRole(RemovePermissionFromRoleInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{
		Text: "successMessage.generic",
	}
	controllers.Respond(ctx, 200, msg, nil)
}

func (ac *AdminRBACController) GetPermissionRoles(ctx *gin.Context) {
	type GetPermissionRolesParams struct {
		PermissionID uint `json:"permissionID"`
	}
	params := controllers.Receive[GetPermissionRolesParams](ctx)
	GetPermissionRolesInfo := rbac.GetPermissionRolesRequest{
		PermissionID: params.PermissionID,
	}
	res, err := ac.rbacService.GetPermissionRoles(GetPermissionRolesInfo)
	if err != nil {
		panic(err)
	}

	controllers.Respond(ctx, 200, controllers.Message{}, res)
}
