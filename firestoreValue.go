package p

import (
	"errors"
	"fmt"
	"strconv"
)

// FirestoreValue holds Firestore fields.
type FirestoreValue struct {
	Fields interface{} `json:"fields"`
}

// getStringValue extracts a string value from a Firestore value
func (v FirestoreValue) getStringValue(name string) (string, error) {
	fields, ok := v.Fields.(map[string]interface{})
	mapped, ok := fields[name].(map[string]interface{})
	if !ok {
		return "", errors.New(fmt.Sprintf("Error extracting value %s from %+v", name, fields))
	}
	value, ok := mapped["stringValue"].(string)
	if !ok {
		return "", errors.New(fmt.Sprintf("Error extracting value %s from %+v", name, fields))
	}
	return value, nil
}

// getIntegerValue extracts an integer value from a Firestore value
func (v FirestoreValue) getIntegerValue(name string) (int, error) {
	fields, ok := v.Fields.(map[string]interface{})
	mapped, ok := fields[name].(map[string]interface{})
	if !ok {
		return 0, errors.New(fmt.Sprintf("Error extracting value %s from %+v", name, fields))
	}
	strValue, ok := mapped["integerValue"].(string)
	if !ok {
		return 0, errors.New(fmt.Sprintf("Error extracting value %s from %+v", name, fields))
	}
	value, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, err
	}
	return value, nil
}
