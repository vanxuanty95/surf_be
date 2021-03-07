package pfit_mgmt

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

var (
	validate *validator.Validate
)

func InitValidationService() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
}

func ValidateStruct(object interface{}) error {
	err := validate.Struct(object)
	if err != nil {

		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(fmt.Sprintf("field %s must be %s", err.Field(), err.ActualTag()))
		}
	}

	return nil
}
