package validation

import (
	"encoding/json"
	"testing"
)

// ValidateContract compares a JSON byte slice against a target structure or simple map.
// This is a lightweight alternative to full JSON Schema validation, focusing on presence and types.
func ValidateContract(t *testing.T, body []byte, schema map[string]interface{}) {
	t.Helper()

	var actual map[string]interface{}
	if err := json.Unmarshal(body, &actual); err != nil {
		t.Fatalf("failed to unmarshal actual body: %v", err)
	}

	for key, expectedType := range schema {
		val, ok := actual[key]
		if !ok {
			t.Errorf("contract violation: missing key '%s'", key)
			continue
		}

		// Simple type checking
		switch expectedType {
		case "string":
			if _, ok := val.(string); !ok {
				t.Errorf("contract violation: key '%s' expected string, got %T", key, val)
			}
		case "number", "float", "int":
			if _, ok := val.(float64); !ok {
				t.Errorf("contract violation: key '%s' expected number, got %T", key, val)
			}
		case "boolean", "bool":
			if _, ok := val.(bool); !ok {
				t.Errorf("contract violation: key '%s' expected boolean, got %T", key, val)
			}
		case "object", "map":
			if _, ok := val.(map[string]interface{}); !ok {
				t.Errorf("contract violation: key '%s' expected object, got %T", key, val)
			}
		case "array", "slice":
			if _, ok := val.([]interface{}); !ok {
				t.Errorf("contract violation: key '%s' expected array, got %T", key, val)
			}
		}
	}
}
