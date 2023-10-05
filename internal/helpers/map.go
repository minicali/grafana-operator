package helpers

import (
	"encoding/json"
	"reflect"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
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

// UnmarshalJSONToMap takes an apiextensionsv1.JSON object and unmarshals it into a map[string]interface{}.
// It returns the unmarshaled map and any error encountered during unmarshaling.
func UnmarshalJSONToMap(jsonData apiextensionsv1.JSON) (map[string]interface{}, error) {
	var result map[string]interface{}
	var intermediate string

	err := json.Unmarshal(jsonData.Raw, &intermediate)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(intermediate), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
