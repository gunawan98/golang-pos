package exception

import "github.com/go-playground/validator/v10"

func RoleValidation(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	switch role {
	case "user", "cashier", "admin":
		return true
	}
	return false
}
