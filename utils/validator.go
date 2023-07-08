package utils

import "gopkg.in/go-playground/validator.v9"

type ValidationError struct {
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func FormatValidationErrors(err error) map[string]ValidationError {
	errors := make(map[string]ValidationError)

	for _, err := range err.(validator.ValidationErrors) {
		var el ValidationError
		el.Tag = err.Tag()
		el.Value = err.Param()
		errors[err.Field()] = el
	}

	return errors
}
