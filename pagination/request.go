package pagination

import (
	"encoding/json"
	"hash/crc32"
	"reflect"
)

type Request interface {
	GetPageToken() string
	GetMaxPageSize() int32
}

func calculateRequestChecksum(request Request) (uint32, error) {
	val := reflect.ValueOf(request).Elem()
	typ := val.Type()
	fieldsToHash := make(map[string]interface{})
	for i := range typ.NumField() {
		field := typ.Field(i)
		if field.Name != "MaxPageSize" && field.Name != "PageToken" {
			// Add this field to our map
			fieldsToHash[field.Name] = val.FieldByName(field.Name).String()
		}
	}
	data, err := json.Marshal(fieldsToHash)
	if err != nil {
		return 0, err
	}
	return crc32.ChecksumIEEE(data), nil
}
