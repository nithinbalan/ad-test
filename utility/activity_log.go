package utility

import "fmt"

// CompareFields compares old and new data, logging only the fields that have changed.
func CompareFieldChanges(oldData, newData map[string]interface{}) map[string]interface{} {
	changedData := make(map[string]interface{})

	for key, oldVal := range oldData {
		newVal := newData[key]

		// Convert both oldVal and newVal to strings, handling pointers
		var oldStr, newStr string

		if oldVal != nil {
			if str, ok := oldVal.(*string); ok && str != nil {
				oldStr = *str
			} else {
				oldStr = fmt.Sprintf("%v", oldVal)
			}
		}

		if newVal != nil {
			if str, ok := newVal.(*string); ok && str != nil {
				newStr = *str
			} else {
				newStr = fmt.Sprintf("%v", newVal)
			}
		}

		// Compare the actual values
		if oldStr != newStr {
			changedData[key] = map[string]interface{}{
				"old": oldStr,
				"new": newStr,
			}
		}
	}

	return changedData
}
