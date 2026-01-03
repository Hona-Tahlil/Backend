package middleware

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/presentation/controllers"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
)

var errTags = bootstrap.Run().Constants.ErrorTags

type RecoveryMiddleware struct {
}

func NewRecoveryMiddleware() *RecoveryMiddleware {
	return &RecoveryMiddleware{}
}

func (rm *RecoveryMiddleware) Recover(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				msgs, statusCode := handleError(err)
				if statusCode != 422 && statusCode != 409 {
					controllers.Respond(ctx, statusCode, msgs[0], nil)
				} else {
					controllers.Respond(ctx, statusCode, msgs, nil)
				}
			}
			ctx.Abort()
		}
	}()
	ctx.Next()
}

func handleError(err error) ([]controllers.Message, int) {
	if bindingErr, ok := err.(*exceptions.BindingError); ok {
		return handleBindingError(bindingErr)
	} else if validationErrs, ok := err.(*exceptions.ValidationErrors); ok {
		return handleValidationErrors(validationErrs)
	} else if authErr, ok := err.(*exceptions.AuthError); ok {
		return handleAuthError(authErr)
	} else if notFoundErr, ok := err.(*exceptions.NotFoundError); ok {
		return handleNotFoundError(notFoundErr)
	} else if conflictErrs, ok := err.(*exceptions.ConflictErrors); ok {
		return handleConflictErrors(conflictErrs)
	} else if forbiddenErr, ok := err.(*exceptions.ForbiddenError); ok {
		return handleForbiddenError(forbiddenErr)
	}
	return unhandledErrors(err)
}

func handleBindingError(bindingErr *exceptions.BindingError) ([]controllers.Message, int) {
	msg := controllers.Message{
		Text: errTags.Binding,
	}
	return []controllers.Message{msg}, 400
}

func handleValidationErrors(validationErrs *exceptions.ValidationErrors) ([]controllers.Message, int) {
	msgs := make([]controllers.Message, len(validationErrs.FieldErrors))
	for i, fieldErr := range validationErrs.FieldErrors {
		var txt string
		if strings.HasPrefix(fieldErr.Tag, "errors.") {
			txt = fieldErr.Tag
		} else {
			txt = "errors." + fieldErr.Tag
		}
		msgs[i] = controllers.Message{
			Text:   txt,
			Params: []string{fieldErr.Field},
		}

	}
	return msgs, 422
}

func handleAuthError(authErr *exceptions.AuthError) ([]controllers.Message, int) {
	msg := controllers.Message{
		Text: authErr.Type,
	}
	return []controllers.Message{msg}, 401
}

func handleNotFoundError(notFoundErr *exceptions.NotFoundError) ([]controllers.Message, int) {
	msg := controllers.Message{
		Text:   errTags.NotFound,
		Params: []string{notFoundErr.Item},
	}
	return []controllers.Message{msg}, 404
}

func handleConflictErrors(conflictErrs *exceptions.ConflictErrors) ([]controllers.Message, int) {
	var msgs []controllers.Message
	for _, fieldErr := range conflictErrs.Errors {
		msgs = append(msgs, controllers.Message{
			Text:   fieldErr.Tag,
			Params: []string{fieldErr.Field},
		})

	}
	return msgs, 409
}

func unhandledErrors(err error) ([]controllers.Message, int) {
	slog.Error("an unhandled error occurred", "err", err)

	msg := controllers.Message{
		Text:   errTags.Generic,
		Params: []string{},
	}
	return []controllers.Message{msg}, 500
}


func handleForbiddenError(forbiddenError *exceptions.ForbiddenError) ([]controllers.Message, int) {
	msg := controllers.Message{
		Text:   forbiddenError.Message,
		Params: []string{forbiddenError.Resource},
	}
	return []controllers.Message{msg}, 403
}
