package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type validateParameter struct {
	Form              interface{}
	HasCustomValidate bool
	CustomValidation  func() error
}

func NewValidate(form interface{}) *validateParameter {
	return &validateParameter{
		Form: form,
	}
}

func (v *validateParameter) SetCustomValidation(hasCustom bool, CustomValidation func() error) {
	v.HasCustomValidate = hasCustom
	v.CustomValidation = CustomValidation
}

func checkPointer(param reflect.Type) bool {
	return param.Kind() == reflect.Ptr || param.Kind() == reflect.Interface
}

func doValidate(types reflect.StructField, values reflect.Value) error {
	validFuncs := types.Tag.Get("valid")
	if validFuncs == "" {
		return nil
	}

	validName := types.Tag.Get("validname")
	if validName == "" {
		validName = types.Name
	}

	valueStr := setValueStr(types, values)
	for _, validation := range strings.Split(validFuncs, ",") {
		options := strings.Split(validation, ":")
		validFunc := options[0]
		validOpt := []string{}

		if len(options) > 1 {
			validOpt = options[1:]
		}

		function, ok := cases[validFunc]
		if !ok {
			return fmt.Errorf("Validation %s not found", validFunc)
		}

		err := function(validName, valueStr, validOpt...)
		if err != nil {
			return err
		}
	}

	return nil
}

func setValueStr(types reflect.StructField, values reflect.Value) string {
	switch types.Type.Name() {
	case "int64":
		return strconv.FormatInt(values.Int(), 10)
	case "string":
		return values.String()
	case "float64":
		return strconv.FormatFloat(values.Float(), 'f', -1, 64)
	}

	return ""
}

func typeCheck(param interface{}) error {
	types := reflect.TypeOf(param)
	values := reflect.ValueOf(param)
	if checkPointer(types) {
		types = types.Elem()
		values = values.Elem()
	}

	for i := 0; i < types.NumField(); i++ {
		value := values.Field(i)
		field := types.Field(i)

		if value.Kind() == reflect.Struct {
			if !value.CanSet() {
				continue
			}

			err := typeCheck(value.Interface())
			if err != nil {
				return err
			}
			continue
		} else if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
			if value.IsNil() {
				continue
			}
			if !value.CanSet() {
				continue
			}

			err := typeCheck(value.Elem().Interface())
			if err != nil {
				return err
			}
			continue
		} else if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
			for i := 0; i < value.Len(); i++ {
				var err error
				if value.Index(i).Kind() == reflect.Struct {
					err = typeCheck(value.Index(i).Interface())
				} else {
					err = doValidate(field, value.Index(i))
				}

				if err != nil {
					return err
				}
			}
			continue
		} else if value.Kind() == reflect.Map {
			continue
		}

		err := doValidate(field, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *validateParameter) Validate() error {
	err := typeCheck(v.Form)
	if err != nil {
		return err
	}

	if v.HasCustomValidate {
		err = v.CustomValidation()
		if err != nil {
			return err
		}
	}

	return nil
}

var cases = map[string]func(string, string, ...string) error{
	"required":  checkRequired,
	"integer":   checkInteger,
	"float":     checkFloat,
	"non_zero":  checkNonZero,
	"bool":      checkBoolean,
	"alpha":     checkAlphabet,
	"alphanum":  checkAlphaNumeric,
	"ascii":     checkASCII,
	"email":     checkEmail,
	"month":     checkMonth,
	"year":      checkYear,
	"maxlength": checkMaxLength,
	"minlength": checkMinLength,
	"datetime":  checkDateTime,
	"currency":  checkCurrency,
	"language":  checkLanguage,
	"phoneno":   checkPhoneNo,
}

var patterns = map[string]string{
	"alpha":    `^[a-zA-Z]+$`,
	"alphanum": `^[0-9a-zA-Z]+$`,
	"email":    "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$",
	"phoneno":  `^(\+?62|0)[1-9][0-9]{7,11}$`,
}

var errors = map[string]string{
	"required":  "%s must not be empty",
	"integer":   "%s must numeric",
	"float":     "%s must numeric",
	"non_zero":  "%s must non-zero",
	"bool":      "%s must boolean",
	"alpha":     "%s must character only",
	"alphanum":  "%s must character and numeric",
	"ascii":     "%s must not have special character",
	"email":     "%s not valid",
	"month":     "%s not valid",
	"year":      "%s not valid",
	"maxlength": "%s must have max %d character(s)",
	"minlength": "%s must have min %d character(s)",
	"datetime":  "%s not valid",
	"currency":  "%s not valid",
	"language":  "%s tidak valid.",
	"phoneno":   "%s not valid",
}

func checkRequired(name, str string, option ...string) error {
	if str == "" {
		errorStr := fmt.Sprintf(errors["required"], name)
		return fmt.Errorf(errorStr)
	}
	return nil
}
func checkInteger(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	if _, err := strconv.ParseInt(str, 10, 64); err != nil {
		errorStr := fmt.Sprintf(errors["integer"], name)
		return fmt.Errorf(errorStr)
	}
	return nil
}
func checkFloat(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	if _, err := strconv.ParseFloat(str, 64); err != nil {
		errorStr := fmt.Sprintf(errors["float"], name)
		return fmt.Errorf(errorStr)
	}
	return nil
}
func checkNonZero(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	if val, err := strconv.ParseFloat(str, 64); err == nil {
		if val == 0 {
			errorStr := fmt.Sprintf(errors["non_zero"], name)
			return fmt.Errorf(errorStr)
		}
	}
	return nil
}
func checkBoolean(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	if str != "true" && str != "false" {
		errorStr := fmt.Sprintf(errors["bool"], name)
		return fmt.Errorf(errorStr)
	}
	return nil
}
func checkAlphabet(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	pattern := patterns["alpha"]
	regex := regexp.MustCompile(pattern)

	if !regex.MatchString(str) {
		errorStr := fmt.Sprintf(errors["alpha"], name)
		return fmt.Errorf(errorStr)
	}
	return nil
}
func checkAlphaNumeric(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	pattern := patterns["alphanum"]
	regex := regexp.MustCompile(pattern)

	if !regex.MatchString(str) {
		errorStr := fmt.Sprintf(errors["alphanum"], name)
		return fmt.Errorf(errorStr)
	}
	return nil
}
func checkASCII(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	pattern := patterns["ascii"]
	regex := regexp.MustCompile(pattern)

	if !regex.MatchString(str) {
		errorStr := fmt.Sprintf(errors["ascii"], name)
		return fmt.Errorf(errorStr)
	}
	return nil
}
func checkEmail(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	pattern := patterns["email"]
	regex := regexp.MustCompile(pattern)

	if !regex.MatchString(str) {
		errorStr := fmt.Sprintf(errors["email"])
		return fmt.Errorf(errorStr)
	}
	return nil
}
func checkMonth(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	errorStr := fmt.Sprintf(errors["month"], name)
	month, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf(errorStr)
	}
	if month < 1 || month > 12 {
		return fmt.Errorf(errorStr)
	}
	return nil
}
func checkYear(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	errorStr := fmt.Sprintf(errors["year"], name)
	month, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf(errorStr)
	}
	if month < 1900 || month > 2200 {
		return fmt.Errorf(errorStr)
	}
	return nil
}
func checkMaxLength(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	if len(option) == 0 {
		return fmt.Errorf("Parameter maksimum empty")
	}

	max, err := strconv.Atoi(option[0])
	if err != nil {
		return fmt.Errorf("Parameter maksimum must be numeric")
	}

	if len(str) > max {
		errorStr := fmt.Sprintf(errors["maxlength"], name, max)
		return fmt.Errorf(errorStr)
	}
	return nil
}
func checkMinLength(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	if len(option) == 0 {
		return fmt.Errorf("Parameter minimum empty")
	}

	min, err := strconv.Atoi(option[0])
	if err != nil {
		return fmt.Errorf("Parameter minimum must be numeric")
	}

	if len(str) < min {
		errorStr := fmt.Sprintf(errors["minlength"], name, min)
		return fmt.Errorf(errorStr)
	}
	return nil
}

func checkDateTime(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	_, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		errorStr := fmt.Sprintf(errors["datetime"], name)
		return fmt.Errorf(errorStr)
	}

	return nil
}

func checkCurrency(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	validCurrency := []string{"IDR"}
	for _, v := range validCurrency {
		if v == str {
			return nil
		}
	}

	errorStr := fmt.Sprintf(errors["currency"], name)
	return fmt.Errorf(errorStr)
}

func checkLanguage(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	validLanguage := []string{"en-US", "id-ID"}
	for _, v := range validLanguage {
		if v == str {
			return nil
		}
	}

	errorStr := fmt.Sprintf(errors["language"])
	return fmt.Errorf(errorStr)
}

func checkPhoneNo(name, str string, option ...string) error {
	if str == "" {
		return nil
	}

	pattern := patterns["phoneno"]
	regex := regexp.MustCompile(pattern)

	if !regex.MatchString(str) {
		errorStr := fmt.Sprintf(errors["phoneno"])
		return fmt.Errorf(errorStr)
	}
	return nil
}

func CheckExpiredCC(month int64, year int64) bool {
	if year < int64(time.Now().UTC().Year()) {
		return true
	}
	if year == int64(time.Now().UTC().Year()) && month < int64(time.Now().UTC().Month()) {
		return true
	}
	return false
}
