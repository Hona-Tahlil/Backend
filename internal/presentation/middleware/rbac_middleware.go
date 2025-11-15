package middleware

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/infrastructure/persistence"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type RBACMiddleware struct {
	uintOfWork *persistence.UnitOfWork
}

func NewRBACMiddleware(uintOfWork *persistence.UnitOfWork) *RBACMiddleware {
	return &RBACMiddleware{
		uintOfWork: uintOfWork,
	}
}

func (rm *RBACMiddleware) NeedsPermission(allowedPermissions []enums.Permission) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID := controllers.GetID(ctx)
		user, err := rm.uintOfWork.Factory().UserRepository().FindUserByID(UserID)
		if err != nil {
			unauthorizedError := exceptions.NewUnauthorizedError("user not found")
			panic(unauthorizedError)
		}

		allowed := rm.isAllowed(allowedPermissions, user.Roles)

		if !allowed {
			accessDeniedErr := exceptions.NewAccessDeniedError("you don't have the required access")
			panic(accessDeniedErr)
		}
		ctx.Next()
	}
}

func (rm *RBACMiddleware) isAllowed(allowedPermissions []enums.Permission, roles []entities.Role) bool {
	var allowed bool = false
	for _, permission := range allowedPermissions {
		for _, role := range roles {
			for _, p := range role.Permissions {
				if p.Type == permission {
					allowed = true
					break
				}
			}
			if allowed {
				break
			}
		}
		if allowed {
			break
		}
	}
	return allowed
}
