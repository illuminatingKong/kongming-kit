package apolloconfig

import (
	"net/http"
	"strings"
)

func toCookieString(cookies []*http.Cookie) string {
	cookieList := make([]string, 0)
	for _, cookie := range cookies {
		cookieList = append(cookieList, cookie.String())

	}
	return strings.Join(cookieList, ";")
}
