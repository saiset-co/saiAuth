package internal

import "github.com/go-playground/validator/v10"

func createErrorResponse(errs map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{
		"Status":    "NOK",
		"ErrorType": "validation",
	}
	resErrors := processErrors(errs)
	res["Errors"] = resErrors

	return res
}

func processErrors(errs map[string]interface{}) map[string]interface{} {
	resErrors := make(map[string]interface{})

	for field, e := range errs {
		fieldResErrors := extractFieldErrors(e.(validator.ValidationErrors))
		resErrors[field] = fieldResErrors
	}

	return resErrors
}

func extractFieldErrors(validationErrors validator.ValidationErrors) []string {
	fieldResErrors := make([]string, 0, len(validationErrors))

	for _, e := range validationErrors {
		fieldError := e.(validator.FieldError)
		fieldResErrors = append(fieldResErrors, fieldError.Tag())
	}

	return fieldResErrors
}
