package validator

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/team-scaletech/common/helpers"
)

// IAPIValidatorService is an interface defining the methods for API request validation
type IAPIValidatorService interface {
	ValidateStruct(req interface{}) (string, bool)
}

// APIValidator is a struct representing the API request validator service
type APIValidator struct{}

// NewAPIValidatorService creates a new instance of the APIValidator service
func NewAPIValidatorService() IAPIValidatorService {
	return &APIValidator{}
}

// ValidateStruct validates the structure of the API request using the provided validator
func (uv *APIValidator) ValidateStruct(req interface{}) (string, bool) {
	// Create a new instance of the validator
	validate := validator.New()
	key := helpers.GetStructName(req)
	// Validate the structure of the request
	err := validate.Struct(req)
	if err != nil {
		// If validation fails, extract and process validation errors
		valErrs := err.(validator.ValidationErrors)
		for _, v := range valErrs {
			// Extract and format field name from the validation error
			fieldName := strings.Replace(strings.Replace(v.Namespace(), key+".", "", 1), ".", " ", 3)
			reg, _ := regexp.Compile("[^A-Z`[]]+")
			fieldName = strings.Replace(reg.ReplaceAllString(fieldName, ""), "[", "", 2)

			// Get the error string using the GetError function (not provided in the code snippet)
			errorString := helpers.GetError(fieldName, v.Tag())

			// Return the error string and false to indicate validation failure
			return errorString, false
		}
	}

	// If validation succeeds, return an empty string and true to indicate validation success
	return "", true
}
