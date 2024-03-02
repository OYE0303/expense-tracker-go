package efactory

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"
)

// copyValues copys non-zero values from src to dest
func copyValues[T any](dest *T, src T) error {
	destValue := reflect.ValueOf(dest).Elem()
	srcValue := reflect.ValueOf(src)

	if destValue.Kind() != reflect.Struct {
		return errDestValueNotStruct
	}

	if srcValue.Kind() != reflect.Struct {
		return errSourceValueNotStruct
	}

	if destValue.Type() != srcValue.Type() {
		return errDestAndSourceIsDiff
	}

	for i := 0; i < destValue.NumField(); i++ {
		destField := destValue.Field(i)
		srcField := srcValue.FieldByName(destValue.Type().Field(i).Name)

		if srcField.IsValid() && destField.Type() == srcField.Type() && !srcField.IsZero() {
			destField.Set(srcField)
		}
	}

	return nil
}

// genFinalError generates a final error message from the given errors
func genFinalError(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	errorMessages := make([]string, len(errs))
	for i, err := range errs {
		errorMessages[i] = err.Error()
	}

	return fmt.Errorf(strings.Join(errorMessages, "\n"))
}

// setNonZeroValues sets non-zero values to the given struct
func setNonZeroValues(i int, v interface{}) {
	val := reflect.ValueOf(v).Elem()
	typeOfVal := val.Type()

	for k := 0; k < val.NumField(); k++ {
		curVal := val.Field(k)
		curField := typeOfVal.Field(k)

		// handle time.Time
		if curField.Type == reflect.TypeOf(time.Time{}) {
			curVal.Set(reflect.ValueOf(time.Now()))
			continue
		}

		// skip unexported fields
		if curField.PkgPath != "" {
			fmt.Println("filed name", curField.Name, "is unexported")
			continue
		}

		// If the field is a struct, recursively set non-zero values for its fields
		if curField.Type.Kind() == reflect.Struct {
			setNonZeroValues(i, curVal.Addr().Interface())
			continue
		}

		// If the field is a pointer, create a new instance of the pointed-to struct type and set its non-zero values
		if curField.Type.Kind() == reflect.Ptr && curField.Type.Elem().Kind() == reflect.Struct {
			if curVal.IsNil() {
				newInstance := reflect.New(curField.Type.Elem()).Elem()
				setNonZeroValues(i, newInstance.Addr().Interface())
				curVal.Set(newInstance.Addr())
			} else {
				setNonZeroValues(i, curVal.Interface())
			}
			continue
		}

		// If the field is a slice, handle it as before
		if curField.Type.Kind() == reflect.Slice {
			if curVal.Len() == 0 {
				elemType := curField.Type.Elem()
				elemValue := reflect.New(elemType).Elem()
				curVal.Set(reflect.Append(curVal, elemValue))
			}

			for j := 0; j < curVal.Len(); j++ {
				elem := curVal.Index(j)
				if elem.Kind() == reflect.Struct {
					setNonZeroValues(i, elem.Addr().Interface())
				} else {
					elem.Set(reflect.ValueOf(genNonZeroValue(elem.Type(), i)))
				}
			}

			continue
		}

		// For other types, set non-zero values if the field is zero
		if curVal.IsZero() && curField.Name != "ID" {
			v := genNonZeroValue(curField.Type, i)
			curVal.Set(reflect.ValueOf(v))
		}
	}
}

// genNonZeroValue generates a non-zero value for the given type
func genNonZeroValue(t reflect.Type, i int) interface{} {
	switch t.Kind() {
	case reflect.Int:
		return int(i)
	case reflect.Int8:
		return int8(i)
	case reflect.Int16:
		return int16(i)
	case reflect.Int32:
		return int32(i)
	case reflect.Int64:
		return int64(i)
	case reflect.Uint:
		return uint(i)
	case reflect.Uint8:
		return uint8(i)
	case reflect.Uint16:
		return uint16(i)
	case reflect.Uint32:
		return uint32(i)
	case reflect.Uint64:
		return uint64(i)
	case reflect.Float32:
		return defaultFloat32
	case reflect.Float64:
		return defaultFloat64
	case reflect.Bool:
		return defaultBool
	case reflect.String:
		return fmt.Sprintf("%s%d", defaultString, i)
	case reflect.Pointer:
		v := genNonZeroValue(t.Elem(), i)
		ptr := reflect.New(t.Elem())
		ptr.Elem().Set(reflect.ValueOf(v))
		return ptr.Interface()

	default:
		return reflect.New(t).Elem().Interface()
	}
}

func camelToSnake(input string) string {
	var buf bytes.Buffer

	for i, r := range input {
		if unicode.IsUpper(r) {
			if i > 0 && unicode.IsLower(rune(input[i-1])) {
				buf.WriteRune('_')
			}
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(r)
		}
	}

	return buf.String()
}
