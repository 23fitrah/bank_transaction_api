package utils

import (
	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.Field()
			var message string

			switch fieldErr.Tag() {
			case "required":
				message = fieldName + " is required"
			case "min":
				message = fieldName + " must be at least " + fieldErr.Param() + " characters"
			case "max":
				message = fieldName + " must be at most " + fieldErr.Param() + " characters"
			case "email":
				message = fieldName + " must be a valid email address"
			case "len":
				message = fieldName + " must be exactly " + fieldErr.Param() + " characters"
			case "numeric":
				message = fieldName + " must be a numeric value"
			default:
				message = fieldName + " is invalid"
			}

			errors["reason"] = message
		}
	}

	return errors
}
