package webapi

import (
	"fmt"
	"testing"
)

func TestOmitEmptyValues(t *testing.T) {
	response := map[string]interface{}{
		"code":  201,
		"data":  map[string]interface{}{"now": "2024-03-14T20:40:04.720266+08:00", "version": "1333"},
		"extra": nil, "httpCode": 200,
		"limit": 0, "message": "get version", "page": 1, "total": 0,
	}
	r := omitEmptyValues(response, "extra", "page", "limit", "total")
	fmt.Println(r)
}
