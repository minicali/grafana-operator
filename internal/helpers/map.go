package helpers

import (
	"reflect"
)

// MergeMaps merges two nested map[string]interface{} types.
// Values from the second map will overwrite values from the first map.
func MergeMaps(dst, src map[string]interface{}) map[string]interface{} {
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			// If the type of the value is a map, then recurse.
			if reflect.ValueOf(srcVal).Kind() == reflect.Map && reflect.ValueOf(dstVal).Kind() == reflect.Map {
				srcMap, srcMapOk := srcVal.(map[string]interface{})
				dstMap, dstMapOk := dstVal.(map[string]interface{})
				if srcMapOk && dstMapOk {
					srcVal = MergeMaps(dstMap, srcMap)
				}
			}
		}
		dst[key] = srcVal
	}
	return dst
}
