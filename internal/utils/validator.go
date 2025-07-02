package utils

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s any) error {
	return validate.Struct(s)
}

func GetValidationErrors(err error) map[string]string {
	var errors = make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errors[fieldError.Field()] = getErrorMessage(fieldError)
		}
	}
	return errors
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Minimum length is " + err.Param()
	case "max":
		return "Maximum length is " + err.Param()
	default:
		return "Invalid value"
	}
}
