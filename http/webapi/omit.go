package webapi

func omitEmptyValues(response map[string]interface{}, keysToOmit ...string) map[string]interface{} {
	filteredResponse := make(map[string]interface{})
	for key, value := range response {
		if shouldOmit(key, value, keysToOmit...) {
			continue
		}
		// 将键值对添加到新的map中
		filteredResponse[key] = value
	}
	return filteredResponse
}

func shouldOmit(key string, value interface{}, keysToOmit ...string) bool {
	// 检查键是否在需要省略的列表中
	for _, k := range keysToOmit {
		if key == k {
			if value == nil {
				return true
			}
			// 检查值是否为空或零值
			switch v := value.(type) {
			case string:
				return v == ""
			case int64:
				return v == int64(0)
			case int:
				return v == 0
			case nil:
				return true
			case map[string]interface{}:
				return len(v) == 0
			default:
				// 对于其他类型，你可能需要添加更多的检查
				return false
			}
		}
	}
	return false
}
