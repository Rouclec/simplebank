package api

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/rouclec/simplebank/util"
)

var validateCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

var validateEmail validator.Func = func(fl validator.FieldLevel) bool {
	if email, ok := fl.Field().Interface().(string); ok {
		// Regular expression to match valid email format
		// Minimum of 3 characters before @, a domain name after @
		re := regexp.MustCompile(`^.{3,}@.+\..+$`)
		return re.MatchString(email)
	}
	return false
}

var validatePassword validator.Func = func(fl validator.FieldLevel) bool {
	if password, ok := fl.Field().Interface().(string); ok {
		uppercaseRegex := regexp.MustCompile(`[A-Z]`)
		lowercaseRegex := regexp.MustCompile(`[a-z]`)
		digitRegex := regexp.MustCompile(`\d`)
		specialCharRegex := regexp.MustCompile(`[@$!%*?&]`)

		hasUppercase := uppercaseRegex.MatchString(password)
		hasLowercase := lowercaseRegex.MatchString(password)
		hasDigit := digitRegex.MatchString(password)
		hasSpecialChar := specialCharRegex.MatchString(password)

		isValid := hasUppercase && hasLowercase && hasDigit && hasSpecialChar && len(password) >= 6
		return isValid
	}
	return false
}
