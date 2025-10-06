package middleware

import (
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/presentation/controllers"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ? should have constants?
type RecoveryMiddleware struct {
}

func NewRecoveryMiddleware() *RecoveryMiddleware {
	return &RecoveryMiddleware{}
}

func (rm *RecoveryMiddleware) Recover(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				msgs, statusCode := rm.handleError(err)
				if len(msgs) == 1 {
					controllers.Respond(ctx, statusCode, msgs[0], nil)
				} else {
					controllers.Respond(ctx, statusCode, msgs, nil)
				}
			}
			ctx.Abort()
		}
	}()

}

func (rm *RecoveryMiddleware) handleError(err error) ([]controllers.Message, int) {
	if bindingErr, ok := err.(*exceptions.BindingError); ok {
		return rm.handleBindingError(bindingErr)
	} else if validationErrs, ok := err.(*exceptions.ValidationErrors); ok {
		return rm.handleValidationErrors(validationErrs)
	}
	return rm.unhandledErrors(err)
}

func (rm *RecoveryMiddleware) handleBindingError(bindingErr *exceptions.BindingError) ([]controllers.Message, int) {
	if numError, ok := bindingErr.Err.(*strconv.NumError); ok {
		msg := controllers.Message{
			Text:   "errors.numeric",
			Params: []string{numError.Num},
		}
		return []controllers.Message{msg}, 400
	}
	msg := controllers.Message{
		Text:   "errors.binding",
		Params: []string{},
	}
	return []controllers.Message{msg}, 400
}

func (rm *RecoveryMiddleware) handleValidationErrors(validationErrs *exceptions.ValidationErrors) ([]controllers.Message, int) {
	msgs := []controllers.Message{}
	for i, fieldErr := range validationErrs.FieldErrors {
		msgs[i] = controllers.Message{
			Text:   "errors." + fieldErr.Tag,
			Params: []string{fieldErr.Field},
		}
	}
	return msgs, 422
}

func (rm *RecoveryMiddleware) unhandledErrors(err error) ([]controllers.Message, int) {
	log.Println("an unhandled error occurred", err.Error())

	msg := controllers.Message{
		Text:   "errors.generic",
		Params: []string{},
	}
	return []controllers.Message{msg}, 500
}
