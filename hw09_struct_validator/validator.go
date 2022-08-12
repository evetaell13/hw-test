package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrInterfaceNotStruct = "interface not struct"
	ErrNotStruct          = errors.New("value is not a struct")
	ErrInvalidType        = "strconv invalid type"
)

var (
	ErrValidateIntMin  = "error value int invalid min = %d"
	ErrValidateIntMax  = "error value int invalid max = %d"
	ErrInvalidIntGroup = "error value invalid in group = %d"
	ErrMin             = errors.New("value invalid min")
	ErrMax             = errors.New("value invalid max")
	ErrGroup           = errors.New("value invalid in group")

	ErrLen                  = errors.New("string invalid len")
	ErrInString             = errors.New("string invalid in group")
	ErrRegexp               = errors.New("string invalid regexp")
	ErrValidateStringLen    = "error string invalid len = %s"
	ErrValidateStringIn     = "error string invalid in group = %s"
	ErrValidateStringRegexp = "error string invalid regexp = %s"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var strErrors strings.Builder
	for _, v := range v {
		tmp := fmt.Sprintf("%s:%s;", v.Field, v.Err)
		strErrors.WriteString(tmp)
	}

	return strErrors.String()
}

func Validate(v interface{}) (sliceErrs []ValidationError) {
	val := reflect.ValueOf(v)

	// интерфейс не структура
	if val.Kind() != reflect.Struct {
		sliceErrs = append(sliceErrs, ValidationError{
			Field: "none",
			Err:   errors.Wrap(ErrNotStruct, ErrInterfaceNotStruct),
		})
		return sliceErrs
	}

	for i := 0; i < val.NumField(); i++ {
		tag := val.Type().Field(i).Tag.Get("validate")
		// Если у поля нет тэга `validate`, пропускаем его
		if tag == "" {
			continue
		}
		rules := getValidateRules(tag)

		field := val.Field(i)
		fieldName := val.Type().Field(i).Name // имя поля
		switch field.Kind() {                 //nolint:exhaustive
		case reflect.String:
			str := field.String()
			for tagKey, tagVal := range rules {
				err := ValidateString(str, tagKey, tagVal)
				if err != nil {
					vr := ValidationError{Field: fieldName, Err: err}
					sliceErrs = append(sliceErrs, vr)
				}
			}
		case reflect.Int:
			dig := field.Int()
			for tagKey, tagVal := range rules {
				err := ValidateInt(int(dig), tagKey, tagVal)
				if err != nil {
					vr := ValidationError{Field: fieldName, Err: err}
					sliceErrs = append(sliceErrs, vr)
				}
			}
		case reflect.Slice:
			errs := GetReflectSlice(field, fieldName, rules)
			sliceErrs = append(sliceErrs, errs...)
		default:
			continue
		}
	}

	return sliceErrs
}

// обработка полей типа string

func ValidateString(v string, key, val string) error {
	switch {
	case key == "regexp":
		re, err := regexp.Compile(val)
		if err != nil {
			return errors.Wrap(err, ErrInvalidType)
		}
		res := re.FindString(v)
		if res == "" {
			return errors.Wrap(ErrRegexp, fmt.Sprintf(ErrValidateStringRegexp, v))
		}
	case key == "len":
		b, _ := strconv.Atoi(val)
		if len([]rune(v)) != b {
			return errors.Wrap(ErrLen, fmt.Sprintf(ErrValidateStringLen, v))
		}
	case key == "in":
		count := 0
		group := strings.Split(val, ",")
		for _, word := range group {
			if word == v {
				count++
			}
		}
		if count == 0 {
			return errors.Wrap(ErrInString, fmt.Sprintf(ErrValidateStringIn, v))
		}
	}

	return nil
}

// обработка полей типа int

func ValidateInt(v int, key, val string) error {
	switch {
	case key == "max":
		b, err := strconv.Atoi(val)
		if err != nil {
			return errors.Wrap(err, ErrInvalidType)
		}
		if v > b {
			return errors.Wrap(ErrMax, fmt.Sprintf(ErrValidateIntMax, v))
		}
	case key == "min":
		b, err := strconv.Atoi(val)
		if err != nil {
			return errors.Wrap(err, ErrInvalidType)
		}
		if v < b {
			return errors.Wrap(ErrMin, fmt.Sprintf(ErrValidateIntMin, v))
		}
	case key == "in":
		count := 0
		group := strings.Split(val, ",")
		for _, word := range group {
			w, err := strconv.Atoi(word)
			if err != nil {
				return errors.Wrap(err, ErrInvalidType)
			}
			if w == v {
				count++
			}
		}
		if count == 0 {
			return errors.Wrap(ErrGroup, fmt.Sprintf(ErrInvalidIntGroup, v))
		}
	}

	return nil
}

func getValidateRules(s string) map[string]string {
	rules := make(map[string]string)

	sl := strings.FieldsFunc(s, func(r rune) bool {
		return string(r) == ":" || string(r) == "|"
	})

	for i := 0; i < len(sl); i += 2 {
		rules[sl[i]] = sl[i+1]
	}

	return rules
}

func GetReflectSlice(field reflect.Value, fieldName string, rules map[string]string) []ValidationError {
	sliceErrs := make([]ValidationError, 0)
	fieldKind := field.Index(0).Kind()

	switch fieldKind { //nolint:exhaustive
	case reflect.Int:
		sl := field.Interface().([]int)
		for _, item := range sl {
			for tagKey, tagVal := range rules {
				err := ValidateInt(item, tagKey, tagVal)
				if err != nil {
					vr := ValidationError{Field: fieldName, Err: err}
					sliceErrs = append(sliceErrs, vr)
				}
			}
		}
	case reflect.String:
		sl := field.Interface().([]string)
		for _, item := range sl {
			for tagKey, tagVal := range rules {
				err := ValidateString(item, tagKey, tagVal)
				if err != nil {
					vr := ValidationError{Field: fieldName, Err: err}
					sliceErrs = append(sliceErrs, vr)
				}
			}
		}
	}

	return sliceErrs
}
