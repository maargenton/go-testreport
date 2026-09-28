//go:build sample

package builderror

import "github.com/maargenton/go-testreport/pkg/sample/builderrorindirect/subpkg"

func Foo(s string) string {
	return subpkg.Foo(s)
}

func Bar(s string) string {
	return s
}
