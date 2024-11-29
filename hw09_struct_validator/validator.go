package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrNotStruct  = errors.New("value is not struct")
	ErrInvalidTag = errors.New("invalid validate tag")
)

var ErrNotInSeq = errors.New("value not in sequency")

var (
	ErrTooBigInt   = errors.New("value is bigger than max")
	ErrTooSmallInt = errors.New("value is smaller than min")
)

var (
	ErrTooBigString         = errors.New("string length is bigger than max length")
	ErrInvalidByRegexString = errors.New("string doesn't match with regex")
)

var ErrNotIntInValidate = errors.New("expect int or []int value in validateInt")
var ErrNotStringInValidate = errors.New("expect string or []string value in validateString")

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	err := make([]string, 0)
	for _, val := range v {
		err = append(err, fmt.Sprintf("%s: %s", val.Field, val.Err))
	}
	return strings.Join(err, "\n")
}

func Validate(v interface{}) error {
	var ve ValidationErrors
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	for _, field := range reflect.VisibleFields(rv.Type()) {
		if field.IsExported() {
			value := rv.FieldByName(field.Name)
			err := validateValue(value, &ve, field)
			if err != nil {
				return err
			}
		}
	}

	if len(ve) != 0 {
		return ve
	}

	return nil
}

func joinErrors(ve *ValidationErrors, err error, field reflect.StructField) {
	for _, v := range strings.Split(err.Error(), "\n") {
		parsedErr := strings.SplitN(v, ": ", 2)
		*ve = append(*ve, ValidationError{Field: field.Name + "." + parsedErr[0], Err: errors.New(parsedErr[1])})
	}
}

func validateValue(value reflect.Value, ve *ValidationErrors, field reflect.StructField) error {
	switch value.Kind() { //nolint:exhaustive
	case reflect.Int:
		err := validateInt(int(value.Int()), ve, field)
		if err != nil {
			return err
		}
	case reflect.String:
		err := validateString(value.String(), ve, field)
		if err != nil {
			return err
		}
	case reflect.Struct:
		tag := field.Tag.Get("validate")
		if tag == "nested" {
			err := Validate(value.Interface())
			if errors.As(err, &ValidationErrors{}) {
				joinErrors(ve, err, field)
			} else {
				return err
			}
		} else if tag != "" {
			return ErrInvalidTag
		}
	case reflect.Slice:
		switch values := value.Interface().(type) {
		case []int:
			err := validateInt(values, ve, field)
			if err != nil {
				return err
			}
		case []string:
			err := validateString(values, ve, field)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateInt(intValue any, ve *ValidationErrors, field reflect.StructField) error {
	tag := field.Tag.Get("validate")

	if tag == "" {
		return nil
	}

	rules := strings.Split(tag, "|")
	for _, rule := range rules {
		parsedRule := strings.Split(rule, ":")

		if len(parsedRule) != 2 {
			return ErrInvalidTag
		}

		switch parsedRule[0] {
		case "max":
			intMax, err := strconv.Atoi(parsedRule[1])
			if err != nil {
				return err
			}
			validateMax(intMax, intValue, ve, field)
		case "min":
			intMin, err := strconv.Atoi(parsedRule[1])
			if err != nil {
				return err
			}
			validateMin(intMin, intValue, ve, field)
		case "in":
			seq := strings.Split(strings.ReplaceAll(parsedRule[1], " ", ""), ",")
			validateIntSeq(seq, intValue, ve, field)
		default:
			return ErrInvalidTag
		}
	}

	return nil
}

func validateString(stringValue any, ve *ValidationErrors, field reflect.StructField) error {
	tag := field.Tag.Get("validate")

	if tag == "" {
		return nil
	}

	rules := strings.Split(tag, "|")
	for _, rule := range rules {
		parsedRule := strings.Split(rule, ":")

		if len(parsedRule) != 2 {
			return ErrInvalidTag
		}

		switch parsedRule[0] {
		case "len":
			maxLen, err := strconv.Atoi(parsedRule[1])
			if err != nil {
				return err
			}
			validateLen(maxLen, stringValue, ve, field)
		case "regexp":
			regex, err := regexp.Compile(parsedRule[1])
			if err != nil {
				return err
			}
			validateRegex(regex, stringValue, ve, field)
		case "in":
			seq := strings.Split(strings.ReplaceAll(parsedRule[1], " ", ""), ",")
			validateStringSeq(seq, stringValue, ve, field)
		default:
			return ErrInvalidTag
		}
	}

	return nil
}

func validateMax(maxValue int, value any, ve *ValidationErrors, field reflect.StructField) error {
	switch v := value.(type) {
	case int:
		if maxValue < v {
			*ve = append(*ve, ValidationError{Field: field.Name, Err: ErrTooBigInt})
		}
	case []int:
		for _, vv := range v {
			if maxValue < vv {
				*ve = append(*ve, ValidationError{Field: field.Name, Err: ErrTooBigInt})
			}
		}
	default:
		return ErrNotIntInValidate
	}
	
	return nil
}

func validateMin(minValue int, value any, ve *ValidationErrors, field reflect.StructField) error {
	switch v := value.(type) {
	case int:
		if minValue > v {
			*ve = append(*ve, ValidationError{Field: field.Name, Err: ErrTooSmallInt})
		}
	case []int:
		for _, vv := range v {
			if minValue > vv {
				*ve = append(*ve, ValidationError{Field: field.Name, Err: ErrTooSmallInt})
			}
		}
	default:
		return ErrNotIntInValidate
	}
	
	return nil
}

func validateIntSeq(seq []string, value any, ve *ValidationErrors, field reflect.StructField) error {
	switch v := value.(type) {
	case int:
		if !slices.Contains(seq, strconv.Itoa(v)) {
			*ve = append(*ve, ValidationError{Field: field.Name, Err: ErrNotInSeq})
		}
	case []int:
		for _, vv := range v {
			if !slices.Contains(seq, strconv.Itoa(vv)) {
				*ve = append(*ve, ValidationError{Field: field.Name, Err: ErrNotInSeq})
			}
		}
	default:
		return ErrNotIntInValidate
	}
	
	return nil
}

func validateLen(maxLen int, value any, ve *ValidationErrors, field reflect.StructField) error {
	switch v := value.(type) {
	case string:
		if maxLen < utf8.RuneCountInString(v) {
			*ve = append(*ve, ValidationError{Field: field.Name, Err: ErrTooBigString})
		}
	case []string:
		for _, vv := range v {
			if maxLen < utf8.RuneCountInString(vv) {
				*ve = append(*ve, ValidationError{Field: field.Name, Err: ErrTooBigString})
			}
		}
	default:
		return ErrNotStringInValidate
	}
	
	return nil
}

func validateRegex(regex *regexp.Regexp, value any, ve *ValidationErrors, field reflect.StructField) error {
	switch v := value.(type) {
	case string:
		if !regex.MatchString(v) {
			*ve = append(*ve, ValidationError{Field: field.Name, Err: ErrInvalidByRegexString})
		}
	case []string:
		for _, vv := range v {
			if !regex.MatchString(vv) {
				*ve = append(*ve, ValidationError{Field: field.Name, Err: ErrInvalidByRegexString})
			}
		}
	default:
		return ErrNotStringInValidate
	}
	
	return nil
}

func validateStringSeq(seq []string, value any, ve *ValidationErrors, field reflect.StructField) error {
	switch v := value.(type) {
	case string:
		if !slices.Contains(seq, v) {
			*ve = append(*ve, ValidationError{Field: field.Name, Err: ErrNotInSeq})
		}
	case []string:
		for _, vv := range v {
			if !slices.Contains(seq, vv) {
				*ve = append(*ve, ValidationError{Field: field.Name, Err: ErrNotInSeq})
			}
		}
	default:
		return ErrNotStringInValidate
	}
	
	return nil
}
