package sql

import (
	"fmt"
	"strings"
)

type Fields func(s []string) string
type Page func(i int) string
type PageSize func(i int) string
type Query func(key string) string

var BaseFields Fields = func(s []string) string {
	f := make([]string, 0)
	for _, i := range s {
		n := fmt.Sprintf("\"%s\"", i)
		f = append(f, n)
	}
	k := strings.Join(f, ", ")
	return fmt.Sprintf("\"fields\":[%s]", k)
}

var BasePage Page = func(i int) string {
	return fmt.Sprintf("\"page\":%d", i)
}

var BasePageSize PageSize = func(i int) string {
	return fmt.Sprintf("\"page_size\":%d", i)
}

func BuildFields(f Fields, s []string) string { return f(s) }

func BuildPage(p Page, i int) string { return p(i) }

func BuildPageSize(ps PageSize, i int) string { return ps(i) }

func BuildQuery(q Query, key string) string {
	return q(key)
}

func BuildSQL(fields, page, pageSize, query string) string {
	return fmt.Sprintf("{%s,%s,%s,%s}", fields, page, pageSize, query)
}
