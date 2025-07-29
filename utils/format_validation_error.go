package utils

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func FormatValidationError(err error, obj interface{}) []map[string]string {
	var validationErrors validator.ValidationErrors
	var errorsList []map[string]string

	if errors.As(err, &validationErrors) {
		objType := reflect.TypeOf(obj)
		if objType.Kind() == reflect.Ptr {
			objType = objType.Elem()
		}

		for _, e := range validationErrors {
			field, _ := objType.FieldByName(e.StructField())
			jsonTag := field.Tag.Get("json")
			jsonField := strings.Split(jsonTag, ",")[0] // ignore "omitempty", etc.

			if jsonField == "" {
				jsonField = strings.ToLower(e.StructField())
			}

			var msg string
			switch e.Tag() {
			case "required":
				msg = jsonField + " is required"
			case "email":
				msg = jsonField + " must be a valid email"
			case "min":
				msg = jsonField + " value is too short"
			case "strong_password":
				msg = jsonField + " must contain uppercase, lowercase, number, and special character"
			default:
				msg = jsonField + " is invalid"
			}

			errorsList = append(errorsList, map[string]string{
				"field":   jsonField,
				"message": msg,
			})
		}
	} else {
		errorsList = append(errorsList, map[string]string{
			"field":   "unknown",
			"message": err.Error(),
		})
	}

	return errorsList
}
