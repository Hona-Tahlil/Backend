package validation

import (
	"reflect"

	"hona/backend/internal/domain/exceptions"

	"github.com/go-playground/validator/v10"
)

func ValidateFields[T any](params T) {
	val := validator.New()
	if err := val.Struct(params); err != nil {
		validationErrors, _ := err.(validator.ValidationErrors)
		customValidationError := exceptions.NewValidationErrors()
		for _, err := range validationErrors {
			field := formatValidationError[T](err)
			customValidationError.AddError(field, err.Tag())
		}
		panic(customValidationError)
	}
}

func formatValidationError[T any](err validator.FieldError) string {
	var params T

	tagTypes := []string{"json", "uri", "form"}

	reflectType := reflect.TypeOf(params)

	field, _ := reflectType.FieldByName(err.StructField())
	tagValue := getAnyTag(field, tagTypes...)

	return tagValue
}

func getAnyTag(field reflect.StructField, tagNames ...string) string {
	for _, tagName := range tagNames {
		if tag := field.Tag.Get(tagName); tag != "" {
			return tag
		}
	}
	return field.Name
}
