package validator

// import (
// 	"github.com/gin-gonic/gin/binding"
// 	"github.com/go-playground/validator/v10"
// )

// var not_empty validator.Func = func(fl validator.FieldLevel) bool {
// 	if str, ok := fl.Field().Interface().(string); ok {
// 		return len(str) != 0
// 	}
// 	return false
// }

// func init() {
// 	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
// 		v.RegisterValidation("not_empty", not_empty)
// 	}
// }
