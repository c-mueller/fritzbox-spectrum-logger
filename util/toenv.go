package util

import (
	"fmt"
	"reflect"
)

const (
	dockerfileFormat = "ENV %s=\"%s\"\n"
	bashFormat       = "export %s=\"%s\"\n"
	fishFormat       = "set -x %s \"%s\"\ns"
)

func MapStructToDockerfileCommands(data interface{}) (string, error) {
	return handleAnalysis(data, dockerfileFormat)
}

func MapStructToBashCommands(data interface{}) (string, error) {
	return handleAnalysis(data, bashFormat)
}

func MapStructToFishCommands(data interface{}) (string, error) {
	return handleAnalysis(data, fishFormat)
}

func handleAnalysis(data interface{}, format string) (string, error) {
	value := reflect.ValueOf(data)
	return getEnvFields(value, format), nil
}

func getEnvFields(value reflect.Value, format string) string {
	kind := value.Kind()
	if kind == reflect.Uintptr || kind == reflect.Ptr || kind == reflect.UnsafePointer || kind == reflect.Invalid {
		// Skip Pointer types
		return ""
	}

	result := ""
	fieldCount := value.Type().NumField()
	fieldType := value.Type()
	for i := 0; i < fieldCount; i++ {
		v := fieldType.Field(i)
		if v.Type.Kind() == reflect.Struct {
			// Look at child structs
			result += getEnvFields(value.Field(i), format)
		} else if v.Type.Kind() == reflect.Array {
			// Arrays are not supported and get ignored
			continue
		} else {
			tagValue, tagFound := v.Tag.Lookup("env")
			// if the current field has a 'env' tag it gets added to the output string
			if tagFound {
				currentField := value.Field(i)
				variableValue := ""

				switch v.Type.Kind() {
				case reflect.String:
					variableValue = currentField.String()
				case reflect.Int:
					variableValue = fmt.Sprint(currentField.Int())
				case reflect.Uint:
					variableValue = fmt.Sprint(currentField.Uint())
				case reflect.Bool:
					variableValue = fmt.Sprint(currentField.Bool())
				case reflect.Float32:
					variableValue = fmt.Sprint(currentField.Float())
				case reflect.Float64:
					variableValue = fmt.Sprint(currentField.Float())

				}

				result += fmt.Sprintf(format, tagValue, variableValue)
			}
		}
	}
	return result
}
