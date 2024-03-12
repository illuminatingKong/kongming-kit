package util

// MergeLabels 新扩展的label 合并到之前的
func MergeLabels(from, to map[string]string) map[string]string {
	if to == nil {
		to = make(map[string]string)
	}

	for k, v := range from {
		if _, ok := to[k]; ok == true {
			if v == to[k] {
				continue
			}
			to[k] = v
		}
		to[k] = v
	}

	return to
}

func SetFieldValueIsNotExist(obj map[string]interface{}, value interface{}, fields ...string) map[string]interface{} {
	m := obj
	for _, field := range fields[:len(fields)-1] {
		if val, ok := m[field]; ok {
			if valMap, ok := val.(map[string]interface{}); ok {
				m = valMap
			} else {
				newVal := make(map[string]interface{})
				m[field] = newVal
				m = newVal
			}
		}
	}
	m[fields[len(fields)-1]] = value
	return obj
}
