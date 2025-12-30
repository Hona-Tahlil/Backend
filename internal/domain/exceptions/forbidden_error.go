package exceptions

import "hona/backend/bootstrap"

var forbiddenErrTags = bootstrap.Run().Constants.ErrorTags

type ForbiddenError struct {
	Message  string
	Resource string
}

func (e ForbiddenError) Error() string {
	return e.Message
}

func NewNoPropertyAccessForbiddenError(resource string) *ForbiddenError {
	msg := forbiddenErrTags.ForbiddenStatus
	if resource != "" {
		return &ForbiddenError{Message: "errors.forbiddenError", Resource: resource}
	}
	return &ForbiddenError{Message: msg, Resource: resource}
}

func NewAccessDeniedForbiddenError(resource string) *ForbiddenError {
	if resource == "" {
		resource = forbiddenErrTags.AccessDenied
	}
	return &ForbiddenError{Message: forbiddenErrTags.AccessDenied, Resource: resource}
}

func NewForbiddenError(message, resource string) *ForbiddenError {
	if message == "" {
		message = forbiddenErrTags.ForbiddenStatus
	}
	return &ForbiddenError{Message: message, Resource: resource}
}
