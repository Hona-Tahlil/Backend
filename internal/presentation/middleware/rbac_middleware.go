package middleware

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/infrastructure/persistence"

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
		userID, _ := ctx.Get(bootstrap.Run().Constants.Context.ID)
		user, err := rm.uintOfWork.Factory().UserRepository().FindUserByID(userID.(uint))
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
