package models

import (
	"encoding/json"
	"io"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
)

func UnmarshalAndValidate[T interface{}](body io.ReadCloser) (s T, err error) {
	var b []byte
	b, err = io.ReadAll(io.Reader(body))
	defer body.Close()
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &s)
	if err != nil {
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("datecustom", dateCustomValidator)
	validate.RegisterValidation("gtfieldcustom", gtFieldIgnoreNilValidator)
	validate.RegisterValidation("energycustom", enumFieldValidator[EnergyClass])
	validate.RegisterValidation("listingtypecustom", enumFieldValidator[ListingType])
	validate.RegisterValidation("sorttypecustom", enumFieldValidator[SortType])
	err = validate.Struct(s)
	return
}

func enumFieldValidator[T StringEnum](fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(T)
	return value.IsEnumValue()
}

func gtFieldIgnoreNilValidator(fl validator.FieldLevel) bool {
	otherField := fl.Parent().FieldByName(fl.Param())
	if !otherField.IsNil() {
		switch fl.Field().Kind() {
		case reflect.Int:
			return otherField.Elem().Int() <= fl.Field().Int()
		case reflect.Float32:
			return otherField.Elem().Float() <= fl.Field().Float()
		}
	}
	return true
}

func dateCustomValidator(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	return err == nil
}
