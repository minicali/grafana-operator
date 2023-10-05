package helpers

import (
	"bytes"

	"gopkg.in/ini.v1"
)

// ToINIConfig takes a map representing the INI configuration and returns a string in INI format.
func ToINIConfig(configMap map[string]interface{}) (string, error) {
	iniConfig := ini.Empty()

	// Populate the INI file structure
	for section, val := range configMap {
		sec, err := iniConfig.NewSection(section)
		if err != nil {
			return "", err
		}

		// Assuming val is a map of key-value pairs
		if kvPairs, ok := val.(map[string]interface{}); ok {
			for k, v := range kvPairs {
				if vStr, ok := v.(string); ok {
					_, err := sec.NewKey(k, vStr)
					if err != nil {
						return "", err
					}
				}
			}
		}
	}

	// Generate INI-formatted string
	var buffer bytes.Buffer
	_, err := iniConfig.WriteTo(&buffer)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
