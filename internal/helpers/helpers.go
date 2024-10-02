package helpers

import (
	"html/template"
)

func SliceHelper(args ...interface{}) []interface{} {
	return args
}

func AttrHelper(s string) template.HTMLAttr {
	return template.HTMLAttr(s)
}
