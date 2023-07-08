package utils

import "reflect"

func MapStructs(mapFrom interface{}, mapTo interface{}) interface{} {
	mapToValue := reflect.ValueOf(mapTo).Elem()
	mapToType := reflect.TypeOf(mapTo).Elem()
	mapFromValue := reflect.ValueOf(mapFrom).Elem()
	mapFromType := reflect.TypeOf(mapFrom).Elem()

	for i := 0; i < mapToValue.NumField(); i++ {
		mapToField := mapToValue.Field(i)
		mapToName := mapToType.Field(i)
		mapFromField := mapFromValue.FieldByName(mapToName.Name)
		mapFromName, found := mapFromType.FieldByName(mapToName.Name)
		if !found {
			continue
		}
		if mapToName.Name == mapFromName.Name {
			mapToField.Set(mapFromField)
		}
	}

	return mapTo
}
