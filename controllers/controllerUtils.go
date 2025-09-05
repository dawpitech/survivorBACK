package controllers

import (
	"reflect"
	"strings"
)

func IsBodyCorrectlyFormed(model any, updates map[string]interface{}) bool {
	allowedFields := getAllowedFields(model)
	for key := range updates {
		if !allowedFields[key] {
			return false
		}
	}
	return true
}

func getAllowedFields(model any) map[string]bool {
	allowed := make(map[string]bool)
	userType := reflect.TypeOf(model)
	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)

		franceDeveloppeTag := field.Tag.Get("fd")
		isAllowed := false
		if franceDeveloppeTag != "" {
			parts := strings.Split(franceDeveloppeTag, ",")
			for _, part := range parts {
				if part == "editable" {
					isAllowed = true
				}
			}
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			parts := strings.Split(jsonTag, ",")
			allowed[parts[0]] = isAllowed
		} else {
			allowed[field.Name] = isAllowed
		}
	}
	return allowed
}
